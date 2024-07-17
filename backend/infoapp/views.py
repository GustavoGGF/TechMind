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

        except Error as e:
            print("Erro ao conectar ao MySQL", e)

        finally:
            if connection.is_connected():
                cursor.close()
                connection.close()

        return JsonResponse(
            {"windows": totalWindows, "total": totalMachines}, status=200, safe=True
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
        data = ""
        system = ""
        name = ""
        distribution = ""
        insertionDate = ""
        macAddress = ""
        connection = ""
        cursor = ""
        results = ""
        select_query = ""
        update_query = ""
        currentUser = ""
        user = ""
        version = ""
        domain = ""
        ip = ""
        model = ""
        serial_number = ""
        max_capacity_memory = ""
        number_of_slots = ""
        first_slot_dim = ""
        second_slot_dim = ""
        third_slot_dim = ""
        fourth_slot_dim = ""
        first_size = ""
        second_size = ""
        third_size = ""
        fourth_size = ""
        first_type = ""
        second_type = ""
        third_type = ""
        fourth_type = ""
        first_type_details = ""
        second_type_details = ""
        third_type_details = ""
        fourth_type_details = ""
        first_speed_memory = ""
        second_speed_memory = ""
        third_speed_memory = ""
        fourth_speed_memory = ""
        first_serial_number = ""
        second_serial_number = ""
        third_serial_number = ""
        fourth_serial_number = ""
        hard_disk_model = ""
        hard_disk_serial_number = ""
        hard_disk_user_capacity = ""
        hard_disk_sata_version = ""
        cpu_architecture = ""
        cpu_operation_mode = ""
        cpus = ""
        cpu_vendor_id = ""
        cpu_model_name = ""
        cpu_thread = ""
        cpu_core = ""
        cpu_socket = ""
        cpu_max_mhz = ""
        cpu_min_mhz = ""
        gpu_product = ""
        gpu_vendor_id = ""
        gpu_bus_info = ""
        gpu_logical_name = ""
        gpu_clock = ""
        gpu_configuration = ""
        audio_device_product = ""
        audio_device_model = ""
        bios_version = ""
        motherboard_manufacturer = ""
        motherboard_product_name = ""
        motherboard_version = ""
        motherboard_serial_name = ""
        motherboard_asset_tag = ""
        softwares = ""

        try:
            data = json.loads(request.body.decode("utf-8"))
            system = data.get("system")
            name = data.get("name")
            distribution = data.get("distribution")
            insertionDate = data.get("insertionDate")
            macAddress = data.get("macAddress")
            user = data.get("currentUser")
            if contains_backslash(user):
                currentUser = user.split("\\")[-1]
            else:
                currentUser = user
            ver = data.get("platformVersion")
            if ver:
                version = ver.split(" ")[0]

            domain = data.get("domain")
            ip = data.get("ip")
            manufacturer = data.get("manufacturer")
            model = data.get("model")
            serial_number = data.get("serialNumber")
            max_capacity_memory = data.get("maxCapacityMemory")
            number_of_slots = data.get("numberOfDevices")
            first_slot_dim = data.get("firstSlotDim")
            second_slot_dim = data.get("secondSlotDim")
            third_slot_dim = data.get("thirdSlotDim")
            fourth_slot_dim = data.get("fourthSlotDim")
            first_size = data.get("firstSize")
            second_size = data.get("secondSize")
            third_size = data.get("thirdSize")
            fourth_size = data.get("fourthSize")
            first_type = data.get("firstType")
            second_type = data.get("secondType")
            third_type = data.get("thirdType")
            fourth_type = data.get("fourthType")
            first_type_details = data.get("firstTypeDetails")
            second_type_details = data.get("secondTypeDetails")
            third_type_details = data.get("thirdTypeDetails")
            first_type_details = data.get("fourthTypeDetails")
            first_speed_memory = data.get("firstSpeedMemory")
            second_speed_memory = data.get("secondSpeedMemory")
            third_speed_memory = data.get("thirdSpeedMemory")
            fourth_speed_memory = data.get("fourthSpeedMemory")
            first_serial_number = data.get("firstSerialNumber")
            second_serial_number = data.get("fecondSerialNumber")
            third_serial_number = data.get("fhirdSerialNumber")
            fourth_serial_number = data.get("fourthSerialNumber")
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
            audio_device_model = data.get("audioDeviceModel")
            bios_version = data.get("biosVersion")
            motherboard_manufacturer = data.get("motherboardManufacturer")
            motherboard_product_name = data.get("motherboardProductName")
            motherboard_version = data.get("motherboardVersion")
            motherboard_serial_name = data.get("motherbaoardSerialName")
            motherboard_asset_tag = data.get("motherboardAssetTag")
            softwares_list = data.get("installedPackages")
            softwares = ""
            for soft in softwares_list:
                softwares += soft + ","

            if macAddress == None:
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
                first_slot_dim = %s, second_slot_dim = %s, third_slot_dim = %s, fourth_slot_dim = %s,
                first_size = %s, second_size = %s, third_size = %s, fourth_size = %s, 
                first_type = %s, second_type = %s, third_type = %s, fourth_type = %s, first_type_details = %s
                , second_type_details = %s, third_type_details = %s, fourth_type_details = %s,
                first_speed_memory = %s, second_speed_memory = %s, third_speed_memory = %s, fourth_speed_memory = %s
                , first_serial_number = %s, second_serial_number = %s, third_serial_number =%s, fourth_serial_number =%s 
                , hard_disk_model = %s, hard_disk_serial_number = %s, hard_disk_user_capacity = %s,
                hard_disk_sata_version = %s, cpu_architecture = %s, cpu_operation_mode = %s, cpus = %s,
                cpu_vendor_id = %s, cpu_model_name = %s, cpu_thread = %s, cpu_core = %s, cpu_socket = %s,
                cpu_max_mhz = %s, cpu_min_mhz = %s, gpu_product = %s, gpu_vendor_id = %s, 
                gpu_bus_info = %s, gpu_logical_name = %s, gpu_clock = %s, gpu_configuration =%s 
                , audio_device_product = %s, audio_device_model = %s, bios_version = %s, 
                motherboard_manufacturer = %s, motherboard_product_name = %s,
                motherboard_version = %s, motherboard_serial_name = %s,
                motherboard_asset_tag = %s, softwares = %s WHERE mac_address = %s"""

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
                        first_slot_dim,
                        second_slot_dim,
                        third_slot_dim,
                        fourth_slot_dim,
                        first_size,
                        second_size,
                        third_size,
                        fourth_size,
                        first_type,
                        second_type,
                        third_type,
                        fourth_type,
                        first_type_details,
                        second_type_details,
                        third_type_details,
                        fourth_type_details,
                        first_speed_memory,
                        second_speed_memory,
                        third_speed_memory,
                        fourth_speed_memory,
                        first_serial_number,
                        second_serial_number,
                        third_serial_number,
                        fourth_serial_number,
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
                max_capacity_memory, number_of_slots, first_slot_dim, second_slot_dim, third_slot_dim,
                fourth_slot_dim, first_size, second_size, third_size, fourth_size, first_type, second_type,
                third_type, fourth_type, first_type_details, second_type_details, third_type_details, fourth_type_details,
                first_speed_memory, second_speed_memory, third_speed_memory, fourth_speed_memory, hard_disk_model,
                hard_disk_serial_number, hard_disk_user_capacity, hard_disk_sata_version, cpu_architecture,
                cpu_operation_mode, cpus, cpu_vendor_id, cpu_model_name, cpu_thread, cpu_core, cpu_socket,
                cpu_max_mhz, cpu_min_mhz, gpu_product, gpu_vendor_id, gpu_bus_info, gpu_logical_name, gpu_clock,
                gpu_configuration, audio_device_product, audio_device_model, bios_version, motherboard_manufacturer,
                motherboard_product_name, motherboard_version, motherboard_serial_name, motherboard_asset_tag,
                softwares) 
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s,
                ,%s , %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s,
                %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"""

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
                        first_slot_dim,
                        second_slot_dim,
                        third_slot_dim,
                        fourth_slot_dim,
                        first_size,
                        second_size,
                        third_size,
                        fourth_size,
                        first_type,
                        second_type,
                        third_type,
                        fourth_type,
                        first_type_details,
                        second_type_details,
                        third_type_details,
                        fourth_type_details,
                        first_speed_memory,
                        second_speed_memory,
                        third_speed_memory,
                        fourth_speed_memory,
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
