from base64 import b64encode
from io import BytesIO
from dns.resolver import resolve
from json import loads, JSONDecodeError
from logging import basicConfig, WARNING, getLogger
from mysql import connector
from pandas import DataFrame, ExcelWriter
from contextlib import contextmanager
from decouple import config
from django.contrib.auth.decorators import login_required
from django.db import transaction
from django.http import JsonResponse, FileResponse
from django.middleware.csrf import get_token
from django_ratelimit.decorators import ratelimit
from django.shortcuts import render
from django.views.decorators.cache import never_cache
from django.views.decorators.csrf import csrf_exempt, requires_csrf_token
from django.views.decorators.http import require_POST, require_GET
from io import BytesIO
from ldap3 import ALL, Connection, Server, SUBTREE
from re import sub
from django.test import Client
from ast import literal_eval
from json import loads
from os import path

# Configuração básica de logging
basicConfig(level=WARNING)
logger = getLogger(__name__)


@contextmanager
def get_database_connection():
    """Context manager for managing database connections."""
    connection = None
    try:
        connection = connector.connect(
            host=config("DB_HOST"),
            database=config("DB_NAME"),
            user=config("DB_USER"),
            password=config("DB_PASSWORD"),
        )
        yield connection
    except connector.Error as err:
        logger.error(f"Database connection error: {err}")
        raise
    finally:
        if connection and connection.is_connected():
            connection.close()


# Create your views here.
@requires_csrf_token
@login_required(login_url="/login")
@require_GET
def home(request):
    return render(request, "index.html", {})


@requires_csrf_token
@login_required(login_url="/login")
@require_GET
def getInfo_main_panel(request):
    connection = None
    cursor = None
    query = None
    result = None
    totalWindows = None
    totalUnix = None
    totalMachines = None
    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = "SELECT COUNT(*) FROM machines WHERE system_name LIKE '%windows%'"
            cursor.execute(query)
            result = cursor.fetchone()

            # Converta os resultados para uma lista de dicionários
            totalWindows = result[0]

    except connector.Error as e:
        logger.error(f"Database query error: {e}")

    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = "SELECT COUNT(mac_address) FROM machines"
            cursor.execute(query)
            result = cursor.fetchone()

            # Converta os resultados para uma lista de dicionários
            totalMachines = result[0]

    except connector.Error as e:
        logger.error(f"Database query error: {e}")

    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = """SELECT COUNT(*)
                        FROM machines
                        WHERE system_name LIKE '%linux%'
                        OR system_name LIKE '%freebsd%';"""
            cursor.execute(query)
            result = cursor.fetchone()

            # Converta os resultados para uma lista de dicionários
            totalUnix = result[0]

    except connector.Error as e:
        logger.error(f"Database query error: {e}")

    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()

    return JsonResponse(
        {"windows": totalWindows, "total": totalMachines, "unix": totalUnix},
        status=200,
        safe=True,
    )


@csrf_exempt
@require_GET
@login_required(login_url="/login")
def getInfo_so_chart(request):
    cursor = None
    query = ""
    results = None
    results_list = None
    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = """SELECT distribution, COUNT(*) as count 
                       FROM machines 
                       GROUP BY distribution;"""
            cursor.execute(query)
            results = cursor.fetchall()

            # Converta os resultados para uma lista de dicionários
            results_list = [{"system_name": row[0], "count": row[1]} for row in results]

    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()

    return JsonResponse(results_list, status=200, safe=False)


@csrf_exempt
@require_GET
@login_required(login_url="/login")
def getInfo_last_update(request):
    cursor = None
    results = None
    results_list = None
    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = """SELECT DISTINCT date_format(insertion_date, '%M %Y'), COUNT(*) AS computer_count
                        from machines
                        GROUP BY date_format(insertion_date, '%M %Y')
                        ORDER BY date_format(insertion_date, '%M %Y');"""
            cursor.execute(query)
            results = cursor.fetchall()

            # Converta os resultados para uma lista de dicionários
            results_list = [{"date": row[0], "count": row[1]} for row in results]

    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)

    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()

    return JsonResponse(results_list, status=200, safe=False)


@requires_csrf_token
@login_required(login_url="/login")
@require_GET
def computers(request):
    return render(request, "index.html", {})


@requires_csrf_token
@login_required(login_url="/login")
@never_cache
@transaction.atomic
@require_GET
def get_data_computers(request, quantity):
    # Iniciando varaiveis
    connection = None
    query = ""
    try:
        # Verificando a quantidade solicitando
        if quantity == "10":
            query = "SELECT * FROM machines ORDER BY insertion_date DESC LIMIT 10"
        elif quantity == "50":
            query = "SELECT * FROM machines ORDER BY insertion_date DESC LIMIT 50"
        elif quantity == "100":
            query = "SELECT * FROM machines ORDER BY insertion_date DESC LIMIT 100"
        elif quantity == "all":
            query = "SELECT * FROM machines ORDER BY insertion_date DESC"
        with get_database_connection() as connection:
            cursor = connection.cursor()
            cursor.execute(query)
            results = [list(row) for row in cursor.fetchall()]

        # Presumindo que o índice da coluna mac_address é 0 (modifique conforme necessário)
        MAC_ADDRESS_INDEX = 0

        # Reverter o endereço MAC para cada item nos resultados
        for row in results:
            row[MAC_ADDRESS_INDEX] = revert_mac_address(row[MAC_ADDRESS_INDEX])

    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)

    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)

    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()

    return JsonResponse({"machines": results}, status=200, safe=True)


def normalize_mac_address(mac):
    return sub(r"[^0-9A-Fa-f]", "", mac).upper()


def revert_mac_address(normalized_mac):
    # Adiciona os separadores de volta ao endereço MAC normalizado
    return ":".join(normalized_mac[i : i + 2] for i in range(0, len(normalized_mac), 2))


def contains_backslash(s):
    return "\\" in s


# Função que recebe os dados do computador
@csrf_exempt
@require_POST
@transaction.atomic
@ratelimit(key="ip", rate="200/d", method="POST", block=True)
@never_cache
def post_machines(request):
    # Declarando varaiveis
    try:
        # Pegando os dados
        data = loads(request.body.decode("utf-8"))
        system = data.get("system")
        name = data.get("name")
        distribution = data.get("distribution")
        insertionDate = data.get("insertionDate")
        macAddress = data.get("macAddress")
        user = data.get("currentUser")

        # Validando o usuario
        if user != None:
            if contains_backslash(user):
                currentUser = user.split("\\")[-1]
            else:
                currentUser = user

        # Validando a versão do SO
        ver = data.get("platformVersion")
        if ver != None:
            version = ver.split(" ")[0]

        domain = data.get("domain")
        ip = data.get("ip")
        manufacturer = data.get("manufacturer")
        if len(manufacturer) > 20:
            logger.error("Manufacturer: ", manufacturer)
            logger.error("Lenv Manufacturer: ", len(manufacturer))
        model = data.get("model")
        serial_number = data.get("serialNumber")
        max_capacity_memory = data.get("maxCapacityMemory")
        number_of_slots = data.get("numberOfDevices")
        memories = data.get("memories")
        memories = str(memories)
        hard_disk_model = data.get("hardDiskModel")
        hard_disk_serial_number = data.get("hardDiskSerialNumber")
        hard_disk_user_capacity = data.get("hardDiskUserCapacity")
        hard_disk_sata_version = data.get("hardDiskSataVersion")
        cpu_architecture = data.get("cpuArchitecture")
        cpu_operation_mode = data.get("cpuOperationMode")
        cpu_vendor_id = data.get("cpuVendorID")
        cpu_model_name = data.get("cpuModelName")
        cpu_thread = data.get("cpuThread")
        cpu_core = data.get("cpuCore")
        cpu_max_mhz = data.get("cpuMaxMHz")
        cpu_min_mhz = data.get("cpuMinMHz")
        gpu_product = data.get("gpuProduct")
        gpu_vendor_id = data.get("gpuVendorID")
        gpu_bus_info = data.get("gpuBusInfo")
        gpu_logical_name = data.get("gpuLogicalName")
        gpu_clock = data.get("gpuClock")
        gpu_configuration = data.get("gpuConfiguration")
        audio_device_product = data.get("audioDeviceProduct")
        audio_device_model = data.get("audioDeviceModel")
        bios_version = data.get("biosVersion").strip()
        motherboard_manufacturer = data.get("motherboardManufacturer")
        motherboard_product_name = data.get("motherboardProductName")
        motherboard_version = data.get("motherboardVersion")
        motherboard_serial_name = data.get("motherbaoardSerialName")
        motherboard_asset_tag = data.get("motherboardAssetTag")
        # Ajustando a lista de softwares
        softwares_list = data.get("installedPackages")
        softwares = None
        try:
            if softwares_list != None:
                if (
                    distribution == "Windows 10"
                    or distribution == "Windows 8.1"
                    or distribution == "Windows Server 2012 R2"
                    or distribution == "Windows Server 2012"
                    or distribution == "Windows10"
                    or distribution == "Microsoft Windows 10 Pro"
                    or distribution == "Microsoft Windows 11 Pro"
                ):
                    softwares = str(softwares_list)
                else:
                    softwares = ",".join(str(soft) for soft in softwares_list)
        except TypeError as e:
            softwares = ""

        license = data.get("license")

        # Verificando a existencia do macAddress
        if macAddress == None:
            logger.error("Mac Address is required")
            return JsonResponse(
                {"error": "Mac Address is required"}, status=400, safe=False
            )

        # Conectando no banco de dados
        connection = connector.connect(
            host=config("DB_HOST"),
            database=config("DB_NAME"),
            user=config("DB_USER"),
            password=config("DB_PASSWORD"),
        )

        if connection.is_connected():
            cursor = connection.cursor()

        # Comando SQL para verificar se o endereço MAC existe na tabela
        select_query = "SELECT * FROM machines WHERE mac_address = %s"

        # Exectando a query
        cursor.execute(select_query, (normalize_mac_address(macAddress),))

        # Obtendo os resultados
        results = cursor.fetchall()

        if results:
            # Comando SQL para atualizar o nome do dispositivo
            update_query = """UPDATE machines SET name = %s, system_name = %s, 
                distribution = %s, insertion_date = %s, logged_user = %s, version = %s , 
                domain = %s, ip = %s, manufacturer= %s, model = %s,
                serial_number = %s, max_capacity_memory = %s, number_of_slots = %s,  
                hard_disk_model = %s, hard_disk_serial_number = %s, hard_disk_user_capacity = %s,
                hard_disk_sata_version = %s, cpu_architecture = %s, cpu_operation_mode = %s,
                cpu_vendor_id = %s, cpu_model_name = %s, cpu_thread = %s, cpu_core = %s,
                cpu_max_mhz = %s, cpu_min_mhz = %s, gpu_product = %s, gpu_vendor_id = %s, 
                gpu_bus_info = %s, gpu_logical_name = %s, gpu_clock = %s, gpu_configuration =%s ,
                audio_device_product = %s, audio_device_model = %s, bios_version = %s, 
                motherboard_manufacturer = %s, motherboard_product_name = %s,
                motherboard_version = %s, motherboard_serial_name = %s,
                motherboard_asset_tag = %s, softwares = %s, memories = %s, license = %s 
                WHERE mac_address = %s"""

            cursor.execute(
                update_query,
                (
                    name,
                    system,
                    distribution,
                    insertionDate,
                    currentUser,
                    version,
                    domain,
                    ip,
                    manufacturer,
                    model,
                    serial_number,
                    max_capacity_memory,
                    number_of_slots,
                    hard_disk_model,
                    hard_disk_serial_number,
                    hard_disk_user_capacity,
                    hard_disk_sata_version,
                    cpu_architecture,
                    cpu_operation_mode,
                    cpu_vendor_id,
                    cpu_model_name,
                    cpu_thread,
                    cpu_core,
                    cpu_max_mhz,
                    cpu_min_mhz,
                    gpu_product,
                    gpu_vendor_id,
                    gpu_bus_info,
                    gpu_logical_name,
                    gpu_clock,
                    gpu_configuration,
                    audio_device_product,
                    audio_device_model,
                    bios_version,
                    motherboard_manufacturer,
                    motherboard_product_name,
                    motherboard_version,
                    motherboard_serial_name,
                    motherboard_asset_tag,
                    softwares,
                    memories,
                    license,
                    (normalize_mac_address(macAddress)),
                ),
            )

            # Confirmando a inserção
            connection.commit()

            # Fechando a conexão
            cursor.close()
            connection.close()

            return JsonResponse({}, status=200, safe=False)
        else:
            query = """INSERT INTO machines (mac_address, name, system_name, distribution, 
                insertion_date, logged_user, version, domain, ip, manufacturer, model, serial_number,
                max_capacity_memory, number_of_slots, hard_disk_model, hard_disk_serial_number, 
                hard_disk_user_capacity, hard_disk_sata_version, cpu_architecture, cpu_operation_mode, 
                cpu_vendor_id, cpu_model_name, cpu_thread, cpu_core, cpu_max_mhz, cpu_min_mhz, gpu_product, 
                gpu_vendor_id, gpu_bus_info, gpu_logical_name, gpu_clock, gpu_configuration, audio_device_product, 
                audio_device_model, bios_version, motherboard_manufacturer, motherboard_product_name, motherboard_version, 
                motherboard_serial_name, motherboard_asset_tag, softwares, memories, license) 
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s,
                %s , %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s,
                %s)"""

            cursor.execute(
                query,
                (
                    (normalize_mac_address(macAddress)),
                    name,
                    system,
                    distribution,
                    insertionDate,
                    currentUser,
                    version,
                    domain,
                    ip,
                    manufacturer,
                    model,
                    serial_number,
                    max_capacity_memory,
                    number_of_slots,
                    hard_disk_model,
                    hard_disk_serial_number,
                    hard_disk_user_capacity,
                    hard_disk_sata_version,
                    cpu_architecture,
                    cpu_operation_mode,
                    cpu_vendor_id,
                    cpu_model_name,
                    cpu_thread,
                    cpu_core,
                    cpu_max_mhz,
                    cpu_min_mhz,
                    gpu_product,
                    gpu_vendor_id,
                    gpu_bus_info,
                    gpu_logical_name,
                    gpu_clock,
                    gpu_configuration,
                    audio_device_product,
                    audio_device_model,
                    bios_version,
                    motherboard_manufacturer,
                    motherboard_product_name,
                    motherboard_version,
                    motherboard_serial_name,
                    motherboard_asset_tag,
                    softwares,
                    memories,
                    license,
                ),
            )

            # Confirmando a inserção
            connection.commit()

            # Fechando a conexão
            cursor.close()
            connection.close()

            # Após realizar um novo POST busca pelos dados atualizados
            try:
                try:
                    with get_database_connection() as connection:
                        cursor = connection.cursor()
                        query = "SELECT COUNT(*) FROM machines WHERE system_name LIKE '%windows%'"
                        cursor.execute(query)
                        result = cursor.fetchone()

                        # Converta os resultados para uma lista de dicionários
                        totalWindows = result[0]

                except connector.Error as e:
                    logger.error(f"Database query error: {e}")

                try:
                    with get_database_connection() as connection:
                        cursor = connection.cursor()
                        query = "SELECT COUNT(mac_address) FROM machines"
                        cursor.execute(query)
                        result = cursor.fetchone()

                        # Converta os resultados para uma lista de dicionários
                        totalMachines = result[0]

                except connector.Error as e:
                    logger.error(f"Database query error: {e}")

                try:
                    with get_database_connection() as connection:
                        cursor = connection.cursor()
                        query = """SELECT COUNT(*)
                                FROM machines
                                WHERE system_name LIKE '%linux%'
                                OR system_name LIKE '%freebsd%';"""
                        cursor.execute(query)
                        result = cursor.fetchone()

                        # Converta os resultados para uma lista de dicionários
                        totalUnix = result[0]

                except connector.Error as e:
                    logger.error(f"Database query error: {e}")

                finally:
                    if connection.is_connected():
                        cursor.close()
                        connection.close()

                client = Client()
                response = client.get("/home/get-Info-Main-Panel/")

                return response
            except Exception as e:
                logger.error(e)

            return JsonResponse({}, status=200, safe=False)

    except JSONDecodeError:
        logger.error("Erro ao decodificar JSON")
        return JsonResponse({"error": "Invalid JSON"}, status=400, safe=False)

    except connector.Error as err:
        logger.error(f"Erro ao inserir dados: {err}")
        return JsonResponse({"error": "Invalid MYSQL"}, status=400, safe=False)

    except Exception as e:
        logger.error(f"Error: {str(e)}")
        return JsonResponse({"error": str(e)}, status=400)


@require_GET
def post_machines_with_mac(request, mac_address):
    return render(request, "index.html", {})


@require_GET
@never_cache
@requires_csrf_token
@login_required(login_url="/login")
def info_machine(request, mac_address):
    connection = None
    cursor = None
    query = None
    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = "SELECT * FROM machines WHERE mac_address = %s"
            # Normalizando o endereço MAC
            normalized_mac = normalize_mac_address(mac_address)
            cursor.execute(query, (normalized_mac,))
            results = cursor.fetchall()

        return JsonResponse({"data": results}, status=200, safe=True)
    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)

    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


@require_GET
@requires_csrf_token
@login_required(login_url="/login")
def devices(request):
    return render(request, "index.html", {})


@login_required(login_url="/login")
@require_POST
@requires_csrf_token
def devices_post(request):
    data = None
    equipament = ""
    model = ""
    serial_number = ""
    imob = ""
    connection = None
    cursor = None
    try:
        data = loads(request.body)
        equipament = data.get("device")
        model = data.get("model")
        serial_number = data.get("serial_number")
        imob = data.get("imob")
        brand = data.get("brand")

        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = """INSERT INTO devices (equipament, model, serial_number, imob, brand) 
                VALUES (%s, %s, %s, %s, %s)"""
            cursor.execute(
                query,
                (equipament, model, serial_number, imob, brand),
            )

        return JsonResponse({}, status=200, safe=True)

    except connector.Error as err:
        if err.errno == 1062:
            with get_database_connection() as connection:
                cursor = connection.cursor()
                query = """UPDATE devices SET equipament = %s, model = %s, 
                imob = %s, brand = %s WHERE serial_number = %s"""
                cursor.execute(
                    query,
                    (equipament, model, serial_number, imob, brand),
                )

            return JsonResponse({}, status=200, safe=False)
        else:
            logger.error(f"Erro ao inserir dados: {err}")
            return JsonResponse({"error": "Invalid MYSQL"}, status=400, safe=False)

    except Exception as e:
        logger.error(e)
        return JsonResponse({"error": "Invalid MYSQL"}, status=400, safe=False)

    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


# Função que pega os dispositivos em quantidade conforme solicitação
@require_GET
@requires_csrf_token
@login_required(login_url="/login")
def devices_get(request, quantity):
    # Declarando algumas variaveis
    connection = None
    cursor = None
    results = None
    query = None
    results = None
    try:
        # Verificando quantidade selecionada
        match quantity:
            case "10":
                query = "SELECT * FROM devices LIMIT 10;"
            case "50":
                query = "SELECT * FROM devices LIMIT 50;"
            case "100":
                query = "SELECT * FROM devices LIMIT 100;"
            case "all":
                query = "SELECT * FROM devices;"

        with get_database_connection() as connection:
            cursor = connection.cursor()
            cursor.execute(query)
            results = [list(row) for row in cursor.fetchall()]

        return JsonResponse({"devices": results}, status=200, safe=True)

    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)

    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


@require_GET
@login_required(login_url="/login")
@requires_csrf_token
def devices_details(request, sn):
    return render(request, "index.html", {})


@require_GET
@login_required(login_url="/login")
@requires_csrf_token
def info_device(request, sn):
    connection = None
    cursor = None
    query = None
    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = "SELECT * FROM devices WHERE serial_number = %s"
            # Normalizando o endereço MAC
            cursor.execute(query, (sn,))
            results = cursor.fetchall()

        return JsonResponse({"data": results}, status=200, safe=True)

    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)
    except Exception as e:
        logger.error(e)
        return JsonResponse({}, status=403, safe=True)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


@require_GET
@login_required(login_url="/login")
@requires_csrf_token
@never_cache
def last_machines(request):
    connection = None
    cursor = None
    query = None
    results = None
    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = "SELECT * FROM machines ORDER BY insertion_date DESC"
            cursor.execute(query)
            results = [list(row) for row in cursor.fetchall()]

        # Presumindo que o índice da coluna mac_address é 0 (modifique conforme necessário)
        MAC_ADDRESS_INDEX = 0

        # Reverter o endereço MAC para cada item nos resultados
        for row in results:
            row[MAC_ADDRESS_INDEX] = revert_mac_address(row[MAC_ADDRESS_INDEX])

        return JsonResponse({"machines": results}, status=200, safe=True)

    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)

    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


@require_GET
def added_devices(request):
    device = None
    data = None
    computer = None
    connection = None
    cursor = None
    try:
        data = loads(request.body)
        device = data.get("device")
        computer = data.get("computer")

        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = "UPDATE devices SET linked_computer = %s WHERE serial_number = %s"
            cursor.execute(
                query,
                (
                    computer,
                    device,
                ),
            )

        return JsonResponse({}, status=200, safe=True)

    except connector.Error as err:
        logger.error(f"Erro ao inserir dados: {err}")
        return JsonResponse({"error": "Invalid MYSQL"}, status=400, safe=False)
    except Exception as e:
        logger.error(e)
        return JsonResponse({}, status=410, safe=True)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


@require_GET
@login_required(login_url="/login")
@requires_csrf_token
@never_cache
def computers_devices(request, mac_address):
    connection = None
    cursor = None
    query = None
    results = None
    mac = None
    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = "SELECT * FROM devices WHERE linked_computer =%s LIMIT 10"
            mac = revert_mac_address(mac_address)
            cursor.execute(query, (mac,))
            results = [list(row) for row in cursor.fetchall()]

        return JsonResponse({"data": results}, status=200, safe=True)
    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)

    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


# Função que salva os dados da aba outros
@requires_csrf_token
@require_POST
@transaction.atomic
@never_cache
def computers_modify(request, mac_address):
    # Pegando os dados json
    data = loads(request.body)
    imob = data.get("imob")
    location = data.get("location")
    note = data.get("note")
    alocate = data.get("alocate")
    new_imob = None
    new_location = None
    new_note = None
    new_alocate = None
    if len(imob) > 1:
        try:
            with get_database_connection() as connection:
                cursor = connection.cursor()
                query = "UPDATE machines SET imob = %s WHERE mac_address = %s"

                cursor.execute(query, (imob, mac_address))

                connection.commit()

        except connector.Error as e:
            logger.error(f"Database query error: {e}")
            return
        finally:
            if connection.is_connected():
                cursor.close()
                connection.close()

        try:
            with get_database_connection() as connection:
                cursor = connection.cursor()
                query = "select imob from machines where mac_address = %s"

                cursor.execute(query, (mac_address,))
                new_imob = cursor.fetchone()

        except connector.Error as e:
            logger.error(f"Database query error: {e}")
            return
        finally:
            if connection.is_connected():
                cursor.close()
                connection.close()

    if len(location) > 1:
        try:
            with get_database_connection() as connection:
                cursor = connection.cursor()
                query = "UPDATE machines SET location = %s WHERE mac_address = %s"

                cursor.execute(query, (location, mac_address))

                connection.commit()

        except connector.Error as e:
            logger.error(f"Database query error: {e}")
            return
        finally:
            if connection.is_connected():
                cursor.close()
                connection.close()

        try:
            with get_database_connection() as connection:
                cursor = connection.cursor()
                query = "select location from machines where mac_address = %s"

                cursor.execute(query, (mac_address,))
                new_location = cursor.fetchone()

        except connector.Error as e:
            logger.error(f"Database query error: {e}")
            return
        finally:
            if connection.is_connected():
                cursor.close()
                connection.close()

    if len(note) > 1:
        try:
            with get_database_connection() as connection:
                cursor = connection.cursor()
                query = "UPDATE machines SET note = %s WHERE mac_address = %s"

                cursor.execute(query, (note, mac_address))

                connection.commit()

        except connector.Error as e:
            logger.error(f"Database query error: {e}")
            return
        finally:
            if connection.is_connected():
                cursor.close()
                connection.close()

        try:
            with get_database_connection() as connection:
                cursor = connection.cursor()
                query = "select note from machines where mac_address = %s"

                cursor.execute(query, (mac_address,))
                new_note = cursor.fetchone()

        except connector.Error as e:
            logger.error(f"Database query error: {e}")
            return
        finally:
            if connection.is_connected():
                cursor.close()
                connection.close()
    if alocate:
        try:
            with get_database_connection() as connection:
                cursor = connection.cursor()
                query = "UPDATE machines SET alocate = 0 WHERE mac_address = %s"

                cursor.execute(query, (mac_address,))

                connection.commit()

        except connector.Error as e:
            logger.error(f"Database query error: {e}")
            return
        finally:
            if connection.is_connected():
                cursor.close()
                connection.close()

        try:
            with get_database_connection() as connection:
                cursor = connection.cursor()
                query = "select alocate from machines where mac_address = %s"

                cursor.execute(query, (mac_address,))
                new_alocate = cursor.fetchone()

        except connector.Error as e:
            logger.error(f"Database query error: {e}")
            return
        finally:
            if connection.is_connected():
                cursor.close()
                connection.close()

    return JsonResponse(
        {
            "imob": new_imob,
            "location": new_location,
            "note": new_note,
            "alocate": new_alocate,
        },
        status=200,
        safe=True,
    )


# Função que libera o token CSRF
@require_GET
def get_new_token(request):
    csrf = get_token(request)
    return JsonResponse({"token": csrf}, status=200, safe=True)


# Requer que esteja logado
@login_required(login_url="/login")
@require_GET
@never_cache
@requires_csrf_token
# Função que retorna as máquinas conforme quantidade solicitada
def get_quantity(request, quantity):
    # Inicia uma variavel e prepara a query conforme quantidade solicitada
    query = None
    match quantity:
        case "10":
            query = "SELECT * FROM machines ORDER BY insertion_date DESC LIMIT 10"
        case "50":
            query = "SELECT * FROM machines ORDER BY insertion_date DESC LIMIT 50"
        case "100":
            query = "SELECT * FROM machines ORDER BY insertion_date DESC LIMIT 100"
        case "all":
            query = "SELECT * FROM machines ORDER BY insertion_date DESC"

    # Iniciando demais varaiveis
    connection = None
    cursor = None
    results = None
    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()
            cursor.execute(query, ())
            results = [list(row) for row in cursor.fetchall()]

        return JsonResponse({"machines": results}, status=200, safe=True)
    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


# Obtem os computadores para o filtro SO
@require_GET
@never_cache
def get_data_so(request):
    # Declarando algumas variaveis
    connection = None
    cursor = None
    query = None
    results = None
    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = "SELECT DISTINCT system_name FROM machines;"
            cursor.execute(query)
            results = [list(row) for row in cursor.fetchall()]

        return JsonResponse({"SO": results}, status=200, safe=True)

    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)

    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


# Obtem os computadores para o filtro de distribution
@require_GET
@never_cache
def get_data_dis(request):
    connection = None
    cursor = None
    query = None
    results = None
    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = "SELECT DISTINCT distribution FROM machines;"
            cursor.execute(query, ())
            results = [list(row) for row in cursor.fetchall()]

        return JsonResponse({"DIS": results}, status=200, safe=True)

    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


# Função executada ao selecionar o SO no filtro
@require_GET
@never_cache
def get_data_so_filter(request, quantity, so):
    # Declarando algumas variaveis
    query = None
    connection = None
    cursor = None
    results = None
    try:
        # Caso a opção escolhida de Distribution seja para mostrar todos
        if so == "all":
            # Então so organiza pela quantidade que deseja ser exibida
            match quantity:
                case "10":
                    query = """SELECT * 
                                    FROM machines 
                                    ORDER BY insertion_date DESC 
                                    LIMIT 10;
                                    """
                case "50":
                    query = """SELECT *
                                    FROM machines 
                                    ORDER BY insertion_date DESC 
                                    LIMIT 50;
                                    """
                case "100":
                    query = """SELECT *
                                    FROM machines 
                                    ORDER BY insertion_date DESC 
                                    LIMIT 100;
                                    """
                case "all":
                    query = """SELECT *
                                    FROM machines 
                                    ORDER BY insertion_date DESC;
                                    """

            with get_database_connection() as connection:
                cursor = connection.cursor()
                cursor.execute(query, ())
                # Obtendo o resultado
                results = [list(row) for row in cursor.fetchall()]

            return JsonResponse({"machines": results}, status=200, safe=True)
        # Caso o SO seja algum especifico então ele busca pelo SO e pela quantidade
        else:
            match quantity:
                case "10":
                    query = """SELECT * 
                                    FROM machines 
                                    WHERE system_name = %s 
                                    ORDER BY insertion_date DESC 
                                    LIMIT 10;
                                    """
                case "50":
                    query = """SELECT * 
                                    FROM machines 
                                    WHERE system_name = %s 
                                    ORDER BY insertion_date DESC 
                                    LIMIT 50;
                                    """
                case "100":
                    query = """SELECT * 
                                    FROM machines 
                                    WHERE system_name = %s 
                                    ORDER BY insertion_date DESC 
                                    LIMIT 100;
                                    """
                case "all":
                    query = """SELECT * 
                                FROM machines 
                                WHERE system_name = %s 
                                ORDER BY insertion_date DESC;
                                """

            with get_database_connection() as connection:
                cursor = connection.cursor()
                cursor.execute(query, (so,))
                results = [list(row) for row in cursor.fetchall()]

            return JsonResponse({"machines": results}, status=200, safe=True)

    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


@require_GET
@never_cache
def get_data_dis_filter(request, quantity, dis):
    # Declarando algumas variaveis
    query = None
    connection = None
    cursor = None
    results = None
    try:
        # Caso a opção escolhida de distribution seja para mostrar todos
        if dis == "all":
            # Então so organiza pela quantidade que deseja ser exibida
            match quantity:
                case "10":
                    query = """SELECT * 
                                    FROM machines 
                                    ORDER BY insertion_date DESC 
                                    LIMIT 10;
                                    """
                case "50":
                    query = """SELECT *
                                    FROM machines 
                                    ORDER BY insertion_date DESC 
                                    LIMIT 50;
                                    """
                case "100":
                    query = """SELECT *
                                    FROM machines 
                                    ORDER BY insertion_date DESC 
                                    LIMIT 100;
                                    """
                case "all":
                    query = """SELECT *
                                    FROM machines 
                                    ORDER BY insertion_date DESC;
                                    """

            with get_database_connection() as connection:
                cursor = connection.cursor()
                cursor.execute(query, ())
                results = [list(row) for row in cursor.fetchall()]

            return JsonResponse({"machines": results}, status=200, safe=True)
        # Caso o SO seja algum especifico então ele busca pelo distribution e pela quantidade
        else:
            match quantity:
                case "10":
                    query = """SELECT * 
                                    FROM machines 
                                    WHERE distribution = %s 
                                    ORDER BY insertion_date DESC 
                                    LIMIT 10;
                                    """
                case "50":
                    query = """SELECT * 
                                    FROM machines 
                                    WHERE distribution = %s 
                                    ORDER BY insertion_date DESC 
                                    LIMIT 50;
                                    """
                case "100":
                    query = """SELECT * 
                                    FROM machines 
                                    WHERE distribution = %s 
                                    ORDER BY insertion_date DESC 
                                    LIMIT 100;
                                    """
                case "all":
                    query = """SELECT * 
                                FROM machines 
                                WHERE distribution = %s 
                                ORDER BY insertion_date DESC;
                                """

            with get_database_connection() as connection:
                cursor = connection.cursor()
                cursor.execute(query, (dis,))
                results = [list(row) for row in cursor.fetchall()]

            return JsonResponse({"machines": results}, status=200, safe=True)

    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


# Função que busca o nome de computar que equivale ao que o usuario esta escrevendo
@require_GET
@never_cache
def get_data_varchar(request, quantity, name):
    query = None
    connection = None
    cursor = None
    results = None
    try:
        match quantity:
            case "10":
                query = """SELECT *
                                FROM machines
                                WHERE name LIKE %s LIMIT 10;
                                """
            case "50":
                query = """SELECT *
                                FROM machines
                                WHERE name LIKE %s LIMIT 50;
                                """
            case "100":
                query = """SELECT *
                                FROM machines
                                WHERE name LIKE %s LIMIT 100;
                                """
            case "all":
                query = """SELECT *
                                FROM machines
                                WHERE name LIKE %s;
                                """

        with get_database_connection() as connection:
            cursor = connection.cursor()
            cursor.execute(query, (f"{name}%",))
            results = [list(row) for row in cursor.fetchall()]

        return JsonResponse({"machines": results}, status=200, safe=True)

    except connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


# Função que gera relatorio DNS mostrando ip Identicos
@require_POST
@never_cache
def get_report_dns(request):
    if request.method == "POST":
        try:
            data = loads(request.body)
            username = data.get("username")
            pwd = data.get("pwd")

            # Conectar ao servidor LDAP
            server = Server(config("SERVER1"), get_info=ALL)
            conn = Connection(
                server,
                user=f"nt-lupatech\{username}",
                password=pwd,
                auto_bind=True,
                read_only=True,
            )

            # Realizar a pesquisa
            conn.search(
                search_base=config("LDAP_BASE"),
                search_filter="(objectClass=computer)",
                attributes=["dnsHostName"],
                search_scope=SUBTREE,
                types_only=False,
            )

            # Processar e imprimir os resultados
            ip_to_hostnames = {}
            for entry in conn.entries:
                if "dnsHostName" in entry:
                    hostname = entry.dnsHostName.value
                    if hostname:
                        ips = get_ip_from_dns(hostname)
                        if ips is not None:
                            for ip in ips:
                                if ip not in ip_to_hostnames:
                                    ip_to_hostnames[ip] = []
                                ip_to_hostnames[ip].append(hostname)
                        else:
                            logger.error(
                                f"get_ip_from_dns returned None for hostname: {hostname}"
                            )

            # Filtrar IPs com múltiplos hostnames e preparar os dados para o DataFrame
            dataPD = []
            for ip, hostnames in ip_to_hostnames.items():
                if len(hostnames) > 1:
                    dataPD.append({"ip": ip, "hostnames": ", ".join(hostnames)})

            # Criar DataFrame
            df = DataFrame(dataPD)

            # Salvar DataFrame em um buffer de memória
            buffer = BytesIO()
            with ExcelWriter(buffer, engine="openpyxl") as writer:
                df.to_excel(writer, index=False)

            # Obter o conteúdo do buffer e codificá-lo em base64
            buffer.seek(0)
            excel_base64 = b64encode(buffer.read()).decode("utf-8")

            # Criar a resposta JSON
            response_data = {
                "filename": "duplicated_ips.xlsx",
                "filedata": excel_base64,
            }

            return JsonResponse(response_data, status=200)

        except Exception as e:
            logger.error(f"An error occurred: {e}")
            return JsonResponse({}, status=420)


# Função que gera os ip's do DNS
def get_ip_from_dns(hostname):
    try:
        # Consulta DNS para registros A (IPv4)
        answers = resolve(hostname, "A")
        ips = [rdata.address for rdata in answers]
        return ips
    except Exception as e:
        pass


# @requires_csrf_token
@require_POST
@csrf_exempt
# @transaction.atomic
# @never_cache
def get_report_xls(request):
    selected_values = None
    selected_values_list = None
    try:
        data = loads(request.body)
        selected_values = data.get("selectedValues")
        if selected_values:
            # Converte a string de volta para uma lista
            selected_values_list = selected_values.split(",")
        # Fazendo o processamento necessário com selected_values_list
        # Itera sobre os valores em selected_values_list
        results = []
        for value in selected_values_list:
            value = normalize_mac_address(value)
            cursor = None
            query = None
            result = None
            try:
                with get_database_connection() as connection:
                    cursor = connection.cursor()

                    # Monta a query substituindo o placeholder pelo valor
                    query = "SELECT * FROM machines WHERE mac_address = %s;"
                    cursor.execute(query, (value,))
                    result = cursor.fetchone()

                    # Adiciona o resultado ao array de resultados
                    results.append(result)
            except connector.Error as e:
                logger.error(f"Database query error for system {value}: {e}")
                return JsonResponse({"status": "fail"}, safe=True, status=312)

            finally:
                if connection.is_connected():
                    cursor.close()
                    connection.close()
    except Exception as e:
        logger.error(e)
        return JsonResponse({"status": "fail"}, safe=True, status=312)

    hostnames = None
    so = None
    dis = None
    version = None
    manufacturer = None
    model = None
    sn = None
    hd_capacity = None
    cpu_model = None
    cpu_manufacturer = None
    gpu_model = None
    softwares_list = None
    software_data = None
    software_names = None
    software_names_str = None
    df = None
    buffer = None
    encoded_file = None
    response_data = None
    try:
        # Extraindo apenas a primeira coluna (MAC Address)
        hostnames = [row[1] for row in results]
        so = [row[2] for row in results]
        dis = [row[3] for row in results]
        version = [row[6] for row in results]
        manufacturer = [row[9] for row in results]
        model = [row[10] for row in results]
        sn = [row[11] for row in results]
        hd_capacity = [row[16] for row in results]
        cpu_model = [row[22] for row in results]
        cpu_manufacturer = [row[21] for row in results]
        gpu_model = [row[32] for row in results]
        softwares_list = [row[40] for row in results]
        software_data = literal_eval(softwares_list[0])
        # Extraindo apenas os nomes dos programas
        software_names = [
            software["name"] for software in software_data if software["name"]
        ]
        # Juntando os nomes separados por vírgula
        software_names_str = ", ".join(software_names)

        buffer = BytesIO()

        # Criar DataFrame com a coluna desejada
        df = DataFrame(
            {
                "HostName": hostnames,
                "SO": so,
                "Distribuição": dis,
                "Versão SO": version,
                "Marca": manufacturer,
                "Modelo": model,
                "SN": sn,
                "HD Capacidade": hd_capacity,
                "CPU Modelo": cpu_model,
                "CPU Fabricante": cpu_manufacturer,
                "GPU Modelo": gpu_model,
                "Programas Instalados": software_names_str,
            }
        )

        df.to_excel(buffer, index=False, sheet_name="Report")

        buffer.seek(0)

        # Codificar o buffer em base64
        encoded_file = b64encode(buffer.read()).decode("utf-8")

        # Criar a resposta JSON com o arquivo codificado
        response_data = {
            "file_name": "report.xlsx",
            "file_content": encoded_file,
        }

        return JsonResponse(response_data, status=200, safe=True)
    except Exception as e:
        logger.error(e)
        return JsonResponse({"status": "fail"}, safe=True, status=312)


def process_results(results):
    processed_data = []

    # Supondo que você queira processar os dados da primeira tupla de results
    if len(results) > 0:
        for i in range(len(results[0])):
            value = results[0][i]

            # Exemplo de tratamento: normalizar o endereço MAC se for o primeiro valor
            if i == 0:
                value = normalize_mac_address(value)

            # Adicionar o valor processado ao array
            processed_data.append(value)

    return processed_data


@require_GET
@csrf_exempt
def get_image(request, model):
    # Define o diretório correto para as imagens
    resultado = model.lower().replace(" ", "")
    base_dir = path.dirname(
        path.dirname(path.abspath(__file__))
    )  # Obtém o diretório base do projeto
    images_dir = path.join(
        base_dir, "static", "assets", "images", "models", f"{resultado}.png"
    )
    logger.error(images_dir)

    if path.isfile(images_dir):
        # Mudando para logger.info(), pois isso não é um erro
        logger.error(f"Arquivo encontrado: {images_dir}")
        return FileResponse(
            open(images_dir, "rb"), as_attachment=True, filename=f"{model}.png"
        )

    # Se o arquivo não for encontrado, retorna um erro
    logger.error(f"Arquivo não encontrado para o modelo: {model}")
    return JsonResponse({"error": "Arquivo não encontrado"}, status=404)

@csrf_exempt
@require_GET
def panel_administrator(request):
    return render(request, "index.html", {})

@require_GET
def panel_get_machines(request):
    results = []
    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()

            # Monta a query substituindo o placeholder pelo valor
            query = "SELECT name, ip, logged_user, insertion_date FROM machines;"
            cursor.execute(query, )
            result = cursor.fetchall()

            # Adiciona o resultado ao array de resultados
            results.append(result)
            
            return JsonResponse({"machines":results}, status=200, safe=True)
    except connector.Error as e:
        logger.error(f"Database query error for system: {e}")
        return JsonResponse({"status": "fail"}, safe=True, status=312)

    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()
