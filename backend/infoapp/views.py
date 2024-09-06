import base64
import dns.resolver
import json
import logging
import mysql.connector
import pandas as pd
from contextlib import contextmanager
from decouple import config
from django.contrib.auth.decorators import login_required
from django.db import transaction
from django.http import JsonResponse
from django.middleware.csrf import get_token
from django_ratelimit.decorators import ratelimit
from django.shortcuts import redirect, render
from django.views.decorators.cache import never_cache
from django.views.decorators.csrf import csrf_exempt, requires_csrf_token
from django.views.decorators.http import require_POST, require_GET
from io import BytesIO
from ldap3 import ALL, Connection, Server, SUBTREE
from mysql.connector import Error
from re import sub
import websockets
import asyncio
from django.test import Client

# Configuração básica de logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


@contextmanager
def get_database_connection():
    """Context manager for managing database connections."""
    connection = None
    try:
        connection = mysql.connector.connect(
            host=config("DB_HOST"),
            database=config("DB_NAME"),
            user=config("DB_USER"),
            password=config("DB_PASSWORD"),
        )
        yield connection
    except mysql.connector.Error as err:
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
def getInfoMainPanel(request):
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

    except mysql.connector.Error as e:
        logger.error(f"Database query error: {e}")

    try:
        with get_database_connection() as connection:
            cursor = connection.cursor()
            query = "SELECT COUNT(mac_address) FROM machines"
            cursor.execute(query)
            result = cursor.fetchone()

            # Converta os resultados para uma lista de dicionários
            totalMachines = result[0]

    except mysql.connector.Error as e:
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

    except mysql.connector.Error as e:
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
def getInfoSOChart(request):
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

    except mysql.connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)

    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)

    return JsonResponse(results_list, status=200, safe=False)


@csrf_exempt
@require_GET
@login_required(login_url="/login")
def getInfoLastUpdate(request):
    cursor = None
    query = ""
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

    except mysql.connector.Error as e:
        logger.error(f"Database query error: {e}")
        return JsonResponse({"error": "Erro ao consultar o banco de dados"}, status=500)

    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return JsonResponse({"error": "Erro inesperado"}, status=500)

    return JsonResponse(results_list, status=200, safe=False)


@requires_csrf_token
@login_required(login_url="/login")
def computers(request):
    if request.method == "POST":
        return redirect("/home")

    if request.method == "GET":
        return render(request, "index.html", {})


# Request csrf
@requires_csrf_token
# Necessita estar logado
@login_required(login_url="/login")
# Função que pega os dados dos computadores por quantidade conforme solicitado
def getDataComputers(request, quantity):
    if request.method == "POST":
        return
    if request.method == "GET":
        # Iniciando varaiveis
        connection = None
        cursor = None
        query = None
        results = None
        MAC_ADDRESS_INDEX = None
        try:
            # Conectando no BD
            connection = mysql.connector.connect(
                host=config("DB_HOST"),
                database=config("DB_NAME"),
                user=config("DB_USER"),
                password=config("DB_PASSWORD"),
            )
            # Verificando a conexão
            if connection.is_connected():
                cursor = connection.cursor()

            # Verificando a quantidade solicitando
            match quantity:
                case "10":
                    query = (
                        "SELECT * FROM machines ORDER BY insertion_date DESC LIMIT 10"
                    )
                case "50":
                    query = (
                        "SELECT * FROM machines ORDER BY insertion_date DESC LIMIT 50"
                    )
                case "100":
                    query = (
                        "SELECT * FROM machines ORDER BY insertion_date DESC LIMIT 100"
                    )
                case "all":
                    query = "SELECT * FROM machines ORDER BY insertion_date DESC"

            # Executandoa a busca
            cursor.execute(query)

            # Obtendo os resultados como listas
            results = [list(row) for row in cursor.fetchall()]
            # Fechando a conexão
            cursor.close()
            connection.close()

            # Presumindo que o índice da coluna mac_address é 0 (modifique conforme necessário)
            MAC_ADDRESS_INDEX = 0

            # Reverter o endereço MAC para cada item nos resultados
            for row in results:
                row[MAC_ADDRESS_INDEX] = revert_mac_address(row[MAC_ADDRESS_INDEX])

            return JsonResponse({"machines": results}, status=200, safe=True)

        except Exception as e:
            print(e)


def normalize_mac_address(mac):
    return sub(r"[^0-9A-Fa-f]", "", mac).upper()


def revert_mac_address(normalized_mac):
    # Adiciona os separadores de volta ao endereço MAC normalizado
    return ":".join(normalized_mac[i : i + 2] for i in range(0, len(normalized_mac), 2))


def contains_backslash(s):
    return "\\" in s


async def send_json():
    async with websockets.connect("ws://localhost:3000/home/") as websocket:
        json_data = {"key": "value", "another_key": "another_value"}
        await websocket.send(json.dumps(json_data))


# Função que recebe os dados do computador
@csrf_exempt
@require_POST
@transaction.atomic
@ratelimit(key="ip", rate="200/d", method="POST", block=True)
@never_cache
def postMachines(request):
    # Declarando varaiveis
    audio_device_model = None
    audio_device_product = None
    bios_version = None
    connection = None
    cpu_architecture = None
    cpu_core = None
    cpu_max_mhz = None
    cpu_min_mhz = None
    cpu_model_name = None
    cpu_operation_mode = None
    cpu_socket = None
    cpu_thread = None
    cpu_vendor_id = None
    cpus = None
    currentUser = None
    cursor = None
    data = None
    distribution = None
    domain = None
    gpu_bus_info = None
    gpu_clock = None
    gpu_configuration = None
    gpu_logical_name = None
    gpu_product = None
    gpu_vendor_id = None
    hard_disk_model = None
    hard_disk_sata_version = None
    hard_disk_serial_number = None
    hard_disk_user_capacity = None
    insertionDate = None
    ip = None
    license = None
    macAddress = None
    max_capacity_memory = None
    memories = None
    model = None
    motherboard_asset_tag = None
    motherboard_manufacturer = None
    motherboard_product_name = None
    motherboard_serial_name = None
    motherboard_version = None
    name = None
    number_of_slots = None
    results = None
    select_query = None
    serial_number = None
    softwares = None
    softwares_list = None
    system = None
    totalWindows = None
    update_query = None
    user = None
    version = None
    try:
        # Pegando os dados
        data = json.loads(request.body.decode("utf-8"))
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
        logger.info(ip)
        manufacturer = data.get("manufacturer")
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
        cpus = data.get("cpus")
        cpu_vendor_id = data.get("cpuVendorID")
        cpu_model_name = data.get("cpuModelName")
        cpu_thread = data.get("cpuThread")
        cpu_core = data.get("cpuCore")
        cpu_socket = data.get("cpuSocket")
        cpu_max_mhz = data.get("cpuMaxMHz")
        cpu_min_mhz = data.get("cpuMinMHz")
        gpu_product = data.get("gpuProduct")

        logger.info(gpu_product)

        gpu_vendor_id = data.get("gpuVendorID")
        gpu_bus_info = data.get("gpuBusInfo")
        gpu_logical_name = data.get("gpuLogicalName")

        logger.info(gpu_logical_name)

        gpu_clock = data.get("gpuClock")
        gpu_configuration = data.get("gpuConfiguration")
        audio_device_product = data.get("audioDeviceProduct")

        logger.info(audio_device_product)

        audio_device_model = data.get("audioDeviceModel")
        bios_version = data.get("biosVersion").strip()

        logger.info(bios_version)

        motherboard_manufacturer = data.get("motherboardManufacturer")
        motherboard_product_name = data.get("motherboardProductName")
        motherboard_version = data.get("motherboardVersion")
        motherboard_serial_name = data.get("motherbaoardSerialName")
        motherboard_asset_tag = data.get("motherboardAssetTag")
        # Ajustando a lista de softwares
        softwares_list = data.get("installedPackages")
        softwares = None
        logger.info(distribution)
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
        connection = mysql.connector.connect(
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
                hard_disk_sata_version = %s, cpu_architecture = %s, cpu_operation_mode = %s, cpus = %s,
                cpu_vendor_id = %s, cpu_model_name = %s, cpu_thread = %s, cpu_core = %s, cpu_socket = %s,
                cpu_max_mhz = %s, cpu_min_mhz = %s, gpu_product = %s, gpu_vendor_id = %s, 
                gpu_bus_info = %s, gpu_logical_name = %s, gpu_clock = %s, gpu_configuration =%s 
                , audio_device_product = %s, audio_device_model = %s, bios_version = %s, 
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
                    cpus,
                    cpu_vendor_id,
                    cpu_model_name,
                    cpu_thread,
                    cpu_core,
                    cpu_socket,
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
                cpus, cpu_vendor_id, cpu_model_name, cpu_thread, cpu_core, cpu_socket, cpu_max_mhz, 
                cpu_min_mhz, gpu_product, gpu_vendor_id, gpu_bus_info, gpu_logical_name, gpu_clock,
                gpu_configuration, audio_device_product, audio_device_model, bios_version, motherboard_manufacturer,
                motherboard_product_name, motherboard_version, motherboard_serial_name, motherboard_asset_tag,
                softwares, memories, license) 
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s,
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
                    cpus,
                    cpu_vendor_id,
                    cpu_model_name,
                    cpu_thread,
                    cpu_core,
                    cpu_socket,
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

                except mysql.connector.Error as e:
                    logger.error(f"Database query error: {e}")

                try:
                    with get_database_connection() as connection:
                        cursor = connection.cursor()
                        query = "SELECT COUNT(mac_address) FROM machines"
                        cursor.execute(query)
                        result = cursor.fetchone()

                        # Converta os resultados para uma lista de dicionários
                        totalMachines = result[0]

                except mysql.connector.Error as e:
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

                except mysql.connector.Error as e:
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

    except json.JSONDecodeError:
        logger.error("Erro ao decodificar JSON")
        return JsonResponse({"error": "Invalid JSON"}, status=400, safe=False)

    except mysql.connector.Error as err:
        logger.error(f"Erro ao inserir dados: {err}")
        return JsonResponse({"error": "Invalid MYSQL"}, status=400, safe=False)

    except Exception as e:
        logger.error(f"Error: {str(e)}")
        return JsonResponse({"error": str(e)}, status=400)


def postMachinesWithMac(request, mac_address):
    if request.method == "GET":
        return render(request, "index.html", {})


def infoMachine(request, mac_address):
    if request.method == "GET":
        connection = None
        cursor = None
        select_query = None
        try:
            connection = mysql.connector.connect(
                host=config("DB_HOST"),
                database=config("DB_NAME"),
                user=config("DB_USER"),
                password=config("DB_PASSWORD"),
            )

            if connection.is_connected():
                cursor = connection.cursor()

            # Comando SQL para verificar se o endereço MAC existe na tabela
            select_query = "SELECT * FROM machines WHERE mac_address = %s"
            # Normalizando o endereço MAC
            normalized_mac = normalize_mac_address(mac_address)

            cursor.execute(select_query, (normalized_mac,))

            # Obtendo os resultados
            results = cursor.fetchall()

            # Fechando a conexão
            cursor.close()
            connection.close()

            return JsonResponse({"data": results}, status=200, safe=True)
        except Exception as e:
            print(e)
            return JsonResponse({}, status=403, safe=True)


def devices(request):
    if request.method == "GET":
        return render(request, "index.html", {})
    if request.method == "POST":
        return


@login_required(login_url="/login")
def devices_post(request):
    if request.method == "GET":
        return redirect("/devices")
    if request.method == "POST":
        data = None
        equipament = ""
        model = ""
        serial_number = ""
        imob = ""
        connection = None
        cursor = None
        try:
            data = json.loads(request.body)
            equipament = data.get("device")
            model = data.get("model")
            serial_number = data.get("serial_number")
            imob = data.get("imob")
            brand = data.get("brand")

            logger.info(equipament)

            connection = mysql.connector.connect(
                host=config("DB_HOST"),
                database=config("DB_NAME"),
                user=config("DB_USER"),
                password=config("DB_PASSWORD"),
            )

            if connection.is_connected():
                cursor = connection.cursor()

            query = """INSERT INTO devices (equipament, model, serial_number, imob, brand) 
                VALUES (%s, %s, %s, %s, %s)"""

            cursor.execute(
                query,
                (equipament, model, serial_number, imob, brand),
            )

            # Confirmando a inserção
            connection.commit()

            # Fechando a conexão
            cursor.close()
            connection.close()

            return JsonResponse({}, status=200, safe=True)

        except mysql.connector.Error as err:
            if err.errno == 1062:
                update_query = """UPDATE devices SET equipament = %s, model = %s, 
                imob = %s, brand = %s WHERE serial_number = %s"""

                logger.info(equipament)

                cursor.execute(
                    update_query, (equipament, model, imob, brand, serial_number)
                )

                # Confirmando a inserção
                connection.commit()

                # Fechando a conexão
                cursor.close()
                connection.close()

                return JsonResponse({}, status=200, safe=False)
            else:
                logger.error(f"Erro ao inserir dados: {err}")
                return JsonResponse({"error": "Invalid MYSQL"}, status=400, safe=False)

        except Exception as e:
            print(e)
            return JsonResponse({"error": "Invalid MYSQL"}, status=400, safe=False)


# Função que pega os dispositivos em quantidade conforme solicitação
def devices_get(request, quantity):
    if request.method == "GET":
        # Declarando algumas variaveis
        connection = None
        cursor = None
        results = None
        select_query = None
        results = None
        try:
            # Conectando no BD
            connection = mysql.connector.connect(
                host=config("DB_HOST"),
                database=config("DB_NAME"),
                user=config("DB_USER"),
                password=config("DB_PASSWORD"),
            )
            # Confirmando conexão
            if connection.is_connected():
                cursor = connection.cursor()
            # Verificando quantidade selecionada
            match quantity:
                case "10":
                    select_query = "SELECT * FROM devices LIMIT 10;"
                case "50":
                    select_query = "SELECT * FROM devices LIMIT 50;"
                case "100":
                    select_query = "SELECT * FROM devices LIMIT 100;"
                case "all":
                    select_query = "SELECT * FROM devices;"

            cursor.execute(select_query)
            # Obtendo os resultados como listas
            results = [list(row) for row in cursor.fetchall()]
            # Fechando a conexão
            cursor.close()
            connection.close()

            return JsonResponse({"devices": results}, status=200, safe=True)
        except Exception as e:
            logger.info(e)


def devices_details(request, sn):
    if request.method == "GET":
        return render(request, "index.html", {})


def infoDevice(request, sn):
    if request.method == "GET":
        connection = None
        cursor = None
        select_query = None
        try:
            connection = mysql.connector.connect(
                host=config("DB_HOST"),
                database=config("DB_NAME"),
                user=config("DB_USER"),
                password=config("DB_PASSWORD"),
            )

            if connection.is_connected():
                cursor = connection.cursor()

            # Comando SQL para verificar se o endereço MAC existe na tabela
            select_query = "SELECT * FROM devices WHERE serial_number = %s"

            cursor.execute(select_query, (sn,))

            # Obtendo os resultados
            results = cursor.fetchall()

            # Fechando a conexão
            cursor.close()
            connection.close()

            return JsonResponse({"data": results}, status=200, safe=True)
        except Exception as e:
            print(e)
            return JsonResponse({}, status=403, safe=True)


def lastMachines(request):
    if request.method == "GET":
        connection = None
        cursor = None
        query = None
        results = None
        try:
            connection = mysql.connector.connect(
                host=config("DB_HOST"),
                database=config("DB_NAME"),
                user=config("DB_USER"),
                password=config("DB_PASSWORD"),
            )

            if connection.is_connected():
                cursor = connection.cursor()

            # Consulta SQL para contar os itens na coluna 'windows' da tabela 'machines'
            query = "SELECT * FROM machines ORDER BY insertion_date DESC"
            cursor.execute(query)

            # Obtendo os resultados como listas
            results = [list(row) for row in cursor.fetchall()]
            # Fechando a conexão
            cursor.close()
            connection.close()

            # Presumindo que o índice da coluna mac_address é 0 (modifique conforme necessário)
            MAC_ADDRESS_INDEX = 0

            # Reverter o endereço MAC para cada item nos resultados
            for row in results:
                row[MAC_ADDRESS_INDEX] = revert_mac_address(row[MAC_ADDRESS_INDEX])

            return JsonResponse({"machines": results}, status=200, safe=True)

        except Exception as e:
            print(e)
    if request.method == "POST":
        return


def addedDevices(request):
    if request.method == "GET":
        return
    if request.method == "POST":
        device = None
        data = None
        computer = None
        connection = None
        cursor = None
        try:
            data = json.loads(request.body)
            device = data.get("device")
            computer = data.get("computer")

            connection = mysql.connector.connect(
                host=config("DB_HOST"),
                database=config("DB_NAME"),
                user=config("DB_USER"),
                password=config("DB_PASSWORD"),
            )

            if connection.is_connected():
                cursor = connection.cursor()

            select_query = (
                "UPDATE devices SET linked_computer = %s WHERE serial_number = %s"
            )

            cursor.execute(select_query, (computer, device))

            # Confirmando a inserção
            connection.commit()
            # Fechando a conexão
            cursor.close()
            connection.close()

            return JsonResponse({}, status=200, safe=True)

        except mysql.connector.Error as err:
            logger.error(f"Erro ao inserir dados: {err}")
            return JsonResponse({"error": "Invalid MYSQL"}, status=400, safe=False)

        except Exception as e:
            logger.error(e)
            return JsonResponse({}, status=410, safe=True)


def computersDevices(request, mac_address):
    if request.method == "GET":
        connection = None
        cursor = None
        query = None
        results = None
        mac = None
        try:
            connection = mysql.connector.connect(
                host=config("DB_HOST"),
                database=config("DB_NAME"),
                user=config("DB_USER"),
                password=config("DB_PASSWORD"),
            )

            if connection.is_connected():
                cursor = connection.cursor()

            # Consulta SQL para contar os itens na coluna 'windows' da tabela 'machines'
            query = "SELECT * FROM devices WHERE linked_computer =%s LIMIT 10"

            mac = revert_mac_address(mac_address)
            logger.info("revert_mac_address:", mac)

            cursor.execute(query, (mac,))

            # Obtendo os resultados como listas
            results = [list(row) for row in cursor.fetchall()]
            # Fechando a conexão
            cursor.close()
            connection.close()

            return JsonResponse({"data": results}, status=200, safe=True)
        except Exception as e:
            logger.info(e)
    if request.method == "POST":
        return


# Função que salva os dados da aba outros
@requires_csrf_token
@require_POST
@transaction.atomic
@never_cache
def computersModify(request, mac_address):
    # Pegando os dados json
    data = json.loads(request.body)
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

        except mysql.connector.Error as e:
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

        except mysql.connector.Error as e:
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

        except mysql.connector.Error as e:
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

        except mysql.connector.Error as e:
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

        except mysql.connector.Error as e:
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

        except mysql.connector.Error as e:
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

        except mysql.connector.Error as e:
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

        except mysql.connector.Error as e:
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
def getToken(request):
    if request.method == "GET":
        csrf = get_token(request)
        return JsonResponse({"token": csrf}, status=200, safe=True)


# Requer que esteja logado
@login_required(login_url="/login")
# Função que retorna as máquinas conforme quantidade solicitada
def getQuantity(request, quantity):
    if request.method == "GET":
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
            # Conectando ao banco
            connection = mysql.connector.connect(
                host=config("DB_HOST"),
                database=config("DB_NAME"),
                user=config("DB_USER"),
                password=config("DB_PASSWORD"),
            )

            # Verificando conexão
            if connection.is_connected():
                cursor = connection.cursor()

            # Executando a query
            cursor.execute(query, ())

            # Obtendo o resultado
            results = [list(row) for row in cursor.fetchall()]

            # Fechando a conexão
            cursor.close()
            connection.close()

            return JsonResponse({"machines": results}, status=200, safe=True)
        except Exception as e:
            logger.error(e)
            return JsonResponse({}, status=420)

    if request.method == "POST":
        return


# Obtem os computadores para o filtro SO
def getDataSO(request):
    if request.method == "GET":
        # Declarando algumas variaveis
        connection = None
        cursor = None
        query = None
        results = None
        try:
            # Conectando ao banco
            connection = mysql.connector.connect(
                host=config("DB_HOST"),
                database=config("DB_NAME"),
                user=config("DB_USER"),
                password=config("DB_PASSWORD"),
            )

            # Verificando conexão
            if connection.is_connected():
                cursor = connection.cursor()

            # Query SQL buscando todos os SO diferentes
            query = "SELECT DISTINCT system_name FROM machines;"

            # Executando a query
            cursor.execute(query, ())

            # Obtendo o resultado
            results = [list(row) for row in cursor.fetchall()]

            # Fechando a conexão
            cursor.close()
            connection.close()

            return JsonResponse({"SO": results}, status=200, safe=True)

        except Exception as e:
            logger.info(e)
            return JsonResponse({}, status=314)
    if request.method == "POST":
        return


# Obtem os computadores para o filtro de distribution
def getDataDIS(request):
    if request.method == "GET":
        connection = None
        cursor = None
        query = None
        results = None
        try:
            # Conectando ao banco
            connection = mysql.connector.connect(
                host=config("DB_HOST"),
                database=config("DB_NAME"),
                user=config("DB_USER"),
                password=config("DB_PASSWORD"),
            )

            # Verificando conexão
            if connection.is_connected():
                cursor = connection.cursor()

            # Query SQL buscando todos os SO diferentes
            query = "SELECT DISTINCT distribution FROM machines;"

            # Executando a query
            cursor.execute(query, ())

            # Obtendo o resultado
            results = [list(row) for row in cursor.fetchall()]

            # Fechando a conexão
            cursor.close()
            connection.close()

            return JsonResponse({"DIS": results}, status=200, safe=True)

        except Exception as e:
            logger.info(e)
            return JsonResponse({}, status=314)
    if request.method == "POST":
        return


# Função executada ao selecionar o SO no filtro
def getDataSoFilter(request, quantity, so):
    if request.method == "GET":
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

                # Conectando ao banco
                connection = mysql.connector.connect(
                    host=config("DB_HOST"),
                    database=config("DB_NAME"),
                    user=config("DB_USER"),
                    password=config("DB_PASSWORD"),
                )

                # Verificando conexão
                if connection.is_connected():
                    cursor = connection.cursor()

                # Executando a query
                cursor.execute(query, ())

                # Obtendo o resultado
                results = [list(row) for row in cursor.fetchall()]

                # Fechando a conexão
                cursor.close()
                connection.close()

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

                # Conectando ao banco
                connection = mysql.connector.connect(
                    host=config("DB_HOST"),
                    database=config("DB_NAME"),
                    user=config("DB_USER"),
                    password=config("DB_PASSWORD"),
                )

                # Verificando conexão
                if connection.is_connected():
                    cursor = connection.cursor()

                # Executando a query
                cursor.execute(query, (so,))

                # Obtendo o resultado
                results = [list(row) for row in cursor.fetchall()]

                # Fechando a conexão
                cursor.close()
                connection.close()

                return JsonResponse({"machines": results}, status=200, safe=True)

        except Exception as e:
            logger.error(e)
            return JsonResponse({}, status=420)
    if request.method == "POST":
        return


def getDataDisFilter(request, quantity, dis):
    if request.method == "GET":
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

                # Conectando ao banco
                connection = mysql.connector.connect(
                    host=config("DB_HOST"),
                    database=config("DB_NAME"),
                    user=config("DB_USER"),
                    password=config("DB_PASSWORD"),
                )

                # Verificando conexão
                if connection.is_connected():
                    cursor = connection.cursor()

                # Executando a query
                cursor.execute(query, ())

                # Obtendo o resultado
                results = [list(row) for row in cursor.fetchall()]

                # Fechando a conexão
                cursor.close()
                connection.close()

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

                # Conectando ao banco
                connection = mysql.connector.connect(
                    host=config("DB_HOST"),
                    database=config("DB_NAME"),
                    user=config("DB_USER"),
                    password=config("DB_PASSWORD"),
                )

                # Verificando conexão
                if connection.is_connected():
                    cursor = connection.cursor()

                # Executando a query
                cursor.execute(query, (dis,))

                # Obtendo o resultado
                results = [list(row) for row in cursor.fetchall()]

                # Fechando a conexão
                cursor.close()
                connection.close()

                return JsonResponse({"machines": results}, status=200, safe=True)
        except Exception as e:
            logger.error(e)
            return JsonResponse({}, status=320)
    if request.method == "POST":
        return


# Função que busca o nome de computar que equivale ao que o usuario esta escrevendo
def getDataVarchar(request, quantity, name):
    if request.method == "GET":
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
            # Conectando ao banco
            connection = mysql.connector.connect(
                host=config("DB_HOST"),
                database=config("DB_NAME"),
                user=config("DB_USER"),
                password=config("DB_PASSWORD"),
            )

            # Verificando conexão
            if connection.is_connected():
                cursor = connection.cursor()

            # Executando a query
            cursor.execute(query, (f"{name}%",))

            # Obtendo o resultado
            results = [list(row) for row in cursor.fetchall()]

            # Fechando a conexão
            cursor.close()
            connection.close()

            return JsonResponse({"machines": results}, status=200, safe=True)

        except Exception as e:
            logger.error(e)
            return JsonResponse({}, status=420)

    if request.method == "POST":
        return


# Função que gera relatorio DNS mostrando ip Identicos
def getReportDNS(request):
    if request.method == "POST":
        try:
            data = json.loads(request.body)
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
            df = pd.DataFrame(dataPD)

            # Salvar DataFrame em um buffer de memória
            buffer = BytesIO()
            with pd.ExcelWriter(buffer, engine="openpyxl") as writer:
                df.to_excel(writer, index=False)

            # Obter o conteúdo do buffer e codificá-lo em base64
            buffer.seek(0)
            excel_base64 = base64.b64encode(buffer.read()).decode("utf-8")

            # Criar a resposta JSON
            response_data = {
                "filename": "duplicated_ips.xlsx",
                "filedata": excel_base64,
            }

            return JsonResponse(response_data, status=200)

        except Exception as e:
            logger.error(f"An error occurred: {e}")
            return JsonResponse({}, status=420)
    elif request.method == "GET":
        return


# Função que gera os ip's do DNS
def get_ip_from_dns(hostname):
    try:
        # Consulta DNS para registros A (IPv4)
        answers = dns.resolver.resolve(hostname, "A")
        ips = [rdata.address for rdata in answers]
        return ips
    except Exception as e:
        pass
