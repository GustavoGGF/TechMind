import json
from django.shortcuts import redirect, render
from django.views.decorators.csrf import requires_csrf_token, csrf_exempt
from django.contrib.auth.decorators import login_required
from django.http import JsonResponse
import mysql.connector
from mysql.connector import Error
from decouple import config
import logging
from re import sub
from django.middleware.csrf import get_token

# Configuração básica de logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


# Create your views here.
@requires_csrf_token
@login_required(login_url="/login")
def home(request):
    if request.method == "POST":
        return
    if request.method == "GET":
        return render(request, "index.html", {})


@requires_csrf_token
@login_required(login_url="/login")
def getInfoMainPanel(request):
    if request.method == "POST":
        return JsonResponse({"Corno": "Manso"}, status=301)
    if request.method == "GET":
        connection = None
        cursor = None
        query = None
        query2 = None
        query3 = None
        result = None
        totalWindows = None
        totalUnix = None
        totalMachines = None
        try:
            # Conectar ao banco de dados
            connection = mysql.connector.connect(
                host=config("DB_HOST"),
                database=config("DB_NAME"),
                user=config("DB_USER"),
                password=config("DB_PASSWORD"),
            )

            if connection.is_connected():
                cursor = connection.cursor()

                # Consulta SQL para contar os itens na coluna 'windows' da tabela 'machines'
            query = "SELECT COUNT(*) FROM machines WHERE system_name LIKE '%windows%'"
            cursor.execute(query)

            # Pegar o resultado da consulta
            result = cursor.fetchone()

            totalWindows = result[0]

            query2 = "SELECT COUNT(mac_address) FROM machines"
            cursor.execute(query2)

            # Pegar o resultado da consulta
            result = cursor.fetchone()

            totalMachines = result[0]

            query3 = """SELECT COUNT(*)
                        FROM machines
                        WHERE system_name LIKE '%linux%'
                        OR system_name LIKE '%freebsd%';"""

            cursor.execute(query3)

            # Pegar o resultado da consulta
            result = cursor.fetchone()

            totalUnix = result[0]

        except Error as e:
            print("Erro ao conectar ao MySQL", e)

        finally:
            if connection.is_connected():
                cursor.close()
                connection.close()

        return JsonResponse(
            {"windows": totalWindows, "total": totalMachines, "unix": totalUnix},
            status=200,
            safe=True,
        )


@requires_csrf_token
@login_required(login_url="/login")
def computers(request):
    if request.method == "POST":
        return redirect("/home")

    if request.method == "GET":
        return render(request, "index.html", {})


@requires_csrf_token
@login_required(login_url="/login")
def getDataComputers(request):
    if request.method == "POST":
        return
    if request.method == "GET":
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
            query = "SELECT * FROM machines ORDER BY insertion_date ASC LIMIT 10"
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


@csrf_exempt
def postMachines(request):
    if request.method == "POST":
        data = None
        system = None
        name = None
        distribution = None
        insertionDate = None
        macAddress = None
        connection = None
        cursor = None
        results = None
        select_query = None
        update_query = None
        currentUser = None
        user = None
        version = None
        domain = None
        ip = None
        model = None
        serial_number = None
        max_capacity_memory = None
        number_of_slots = None
        hard_disk_model = None
        hard_disk_serial_number = None
        hard_disk_user_capacity = None
        hard_disk_sata_version = None
        cpu_architecture = None
        cpu_operation_mode = None
        cpus = None
        cpu_vendor_id = None
        cpu_model_name = None
        cpu_thread = None
        cpu_core = None
        cpu_socket = None
        cpu_max_mhz = None
        cpu_min_mhz = None
        gpu_product = None
        gpu_vendor_id = None
        gpu_bus_info = None
        gpu_logical_name = None
        gpu_clock = None
        gpu_configuration = None
        audio_device_product = None
        audio_device_model = None
        bios_version = None
        motherboard_manufacturer = None
        motherboard_product_name = None
        motherboard_version = None
        motherboard_serial_name = None
        motherboard_asset_tag = None
        softwares = None
        softwares_list = None
        memories = None

        try:
            data = json.loads(request.body.decode("utf-8"))
            system = data.get("system")
            name = data.get("name")
            distribution = data.get("distribution")
            insertionDate = data.get("insertionDate")
            macAddress = data.get("macAddress")
            user = data.get("currentUser")

            if user != None:
                if contains_backslash(user):
                    currentUser = user.split("\\")[-1]
                else:
                    currentUser = user

            ver = data.get("platformVersion")
            if ver != None:
                version = ver.split(" ")[0]

            domain = data.get("domain")
            ip = data.get("ip")
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
            gpu_vendor_id = data.get("gpuVendorID")
            gpu_bus_info = data.get("gpuBusInfo")
            gpu_logical_name = data.get("gpuLogicalName")
            gpu_clock = data.get("gpuClock")
            gpu_configuration = data.get("gpuConfiguration")
            audio_device_product = data.get("audioDeviceProduct")
            logger.info("audio_device_product: ", audio_device_product)
            audio_device_model = data.get("audioDeviceModel")
            bios_version = data.get("biosVersion")
            motherboard_manufacturer = data.get("motherboardManufacturer")
            motherboard_product_name = data.get("motherboardProductName")
            motherboard_version = data.get("motherboardVersion")
            motherboard_serial_name = data.get("motherbaoardSerialName")
            motherboard_asset_tag = data.get("motherboardAssetTag")
            softwares_list = data.get("installedPackages")
            softwares = None
            if softwares_list != None:
                if distribution == "Windows10" or distribution == "Windows8.1":
                    softwares = str(softwares_list)
                else:
                    softwares = ""
                    for soft in softwares_list:
                        softwares += soft + ","

            if macAddress == None:
                logger.error("Mac Address is required")
                return JsonResponse(
                    {"error": "Mac Address is required"}, status=400, safe=False
                )

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
            normalized_mac = normalize_mac_address(macAddress)

            cursor.execute(select_query, (normalized_mac,))

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
                motherboard_asset_tag = %s, softwares = %s, memories = %s WHERE mac_address = %s"""

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
                        normalized_mac,
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
                softwares, memories) 
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s,
                %s , %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"""

                cursor.execute(
                    query,
                    (
                        normalized_mac,
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
                    ),
                )

                # Confirmando a inserção
                connection.commit()

                # Fechando a conexão
                cursor.close()
                connection.close()

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

    if request.method == "GET":
        return redirect("/computers")


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


def devices_token(request):
    if request.method == "GET":
        csrf = get_token(request)
        return JsonResponse({"token": csrf}, status=200, safe=True)


def devices_get(request):
    if request.method == "GET":
        connection = None
        cursor = None
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

            # Comando SQL para verificar se o endereço MAC existe na tabela
            select_query = "SELECT * FROM devices LIMIT 10;"
            cursor.execute(select_query)
            # Obtendo os resultados como listas
            results = [list(row) for row in cursor.fetchall()]
            # Fechando a conexão
            cursor.close()
            connection.close()

            return JsonResponse({"Dispositivos": results}, status=200, safe=True)
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
            query = "SELECT * FROM machines ORDER BY insertion_date ASC LIMIT 10"
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


def computersModify(request, mac_address):
    if request.method == "GET":
        return
    if request.method == "POST":
        data = None
        imob = None
        connection = None
        cursor = None
        location = None
        update_query = None
        result = None
        result2 = None
        try:
            data = json.loads(request.body)
            imob = data.get("imob")
            location = data.get("location")
            if imob:
                connection = mysql.connector.connect(
                    host=config("DB_HOST"),
                    database=config("DB_NAME"),
                    user=config("DB_USER"),
                    password=config("DB_PASSWORD"),
                )

                if connection.is_connected():
                    cursor = connection.cursor()

                update_query = "UPDATE machines SET imob =%s WHERE mac_address =%s"

                cursor.execute(update_query, (imob, mac_address))

                # Confirmando a inserção
                connection.commit()

                update_query = "select imob from machines WHERE mac_address =%s"

                cursor.execute(update_query, (mac_address,))

                # Obtendo o resultado
                result = (
                    cursor.fetchone()
                )  # Use fetchall() se esperar mais de um resultado

                # Fechando a conexão
                cursor.close()
                connection.close()

                if location:
                    connection = mysql.connector.connect(
                        host=config("DB_HOST"),
                        database=config("DB_NAME"),
                        user=config("DB_USER"),
                        password=config("DB_PASSWORD"),
                    )

                    if connection.is_connected():
                        cursor = connection.cursor()

                    update_query = (
                        "UPDATE machines SET location =%s WHERE mac_address =%s"
                    )

                    cursor.execute(update_query, (location, mac_address))

                    # Confirmando a inserção
                    connection.commit()

                    update_query = "select imob from machines WHERE mac_address =%s"

                    cursor.execute(update_query, (mac_address,))

                    # Obtendo o resultado
                    result2 = (
                        cursor.fetchone()
                    )  # Use fetchall() se esperar mais de um resultado

                    # Fechando a conexão
                    cursor.close()
                    connection.close()

                    return JsonResponse(
                        {"imob": result[0], "location": result2[0]},
                        status=200,
                        safe=True,
                    )
                else:
                    return JsonResponse({"imob": result[0]}, status=200, safe=True)
            if location:
                connection = mysql.connector.connect(
                    host=config("DB_HOST"),
                    database=config("DB_NAME"),
                    user=config("DB_USER"),
                    password=config("DB_PASSWORD"),
                )

                if connection.is_connected():
                    cursor = connection.cursor()

                update_query = "UPDATE machines SET location =%s WHERE mac_address =%s"

                cursor.execute(update_query, (location, mac_address))

                # Confirmando a inserção
                connection.commit()

                update_query = "select imob from machines WHERE mac_address =%s"

                cursor.execute(update_query, (mac_address,))

                # Obtendo o resultado
                result = (
                    cursor.fetchone()
                )  # Use fetchall() se esperar mais de um resultado
                logger.info(result)

                # Fechando a conexão
                cursor.close()
                connection.close()

                return JsonResponse({"location": result[0]}, status=200, safe=True)
            else:
                return JsonResponse(
                    {"message": "Imobilizado ou Localização Obrigatorio"}, status=310
                )
        except Exception as e:
            logger.error(e)
            return JsonResponse({}, status=420)


def getToken(request):
    if request.method == "GET":
        csrf = get_token(request)
        return JsonResponse({"token": csrf}, status=200, safe=True)
