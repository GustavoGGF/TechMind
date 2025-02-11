package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"techmind/linux/cpuinformation"
	"techmind/linux/generalinformation"
	"techmind/linux/hdinformation"
	"techmind/linux/internetinformation"
	"techmind/linux/mbinformation"
	"techmind/linux/motherboardinformation"
)

func logToFile(msg string) {
	// Abrir arquivo de log
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo de log:", err)
		return
	}
	defer file.Close()

	// Criar um logger para o arquivo
	logger := log.New(file, "", log.LstdFlags)

	// Escrever mensagem no arquivo de log
	logger.Println(msg)
}

type Data struct {
	System                string   `json:"system"`
	Name                  string   `json:"name"`
	Distribution          string   `json:"distribution"`
	InterfaceInternet     string   `json:"interfaceInternet"`
	MacAddress            string   `json:"macAddress"`
	InsertionDate         string   `json:"insertionDate"`
	CurrentUser           string   `json:"currentUser"`
	PlatformVersion       string   `json:"platformVersion"`
	Domain                string   `json:"domain"`
	IP                    string   `json:"ip"`
	Manufacturer          string   `json:"manufacturer"`
	Model                 string   `json:"model"`
	SerialNumber          string   `json:"serialNumber"`
	MaxCapacityMemory     string   `json:"maxCapacityMemory"`
	NumberOfDevices       string   `json:"numberOfDevices"`
	HardDiskModel         string   `json:"hardDiskModel"`
	HardDiskSerialNumber  string   `json:"hardDiskSerialNumber"`
	HardDiskUserCapacity  string   `json:"hardDiskUserCapacity"`
	HardDiskSataVersion   string   `json:"hardDiskSataVersion"`
	CPUArchitecture       string   `json:"cpuArchitecture"`
	CPUOperationMode      string   `json:"cpuOperationMode"`
	CPUVendorID           string   `json:"cpuVendorID"`
	CPUModelName          string   `json:"cpuModelName"`
	CPUThread             int   `json:"cpuThread"`
	CPUCore               int   `json:"cpuCore"`
	CPUMaxMHz             int   `json:"cpuMaxMHz"`
	CPUMinMHz             int   `json:"cpuMinMHz"`
	GPUProduct            string   `json:"gpuProduct"`
	GPUVendorID           string   `json:"gpuVendorID"`
	GPUBusInfo            string   `json:"gpuBusInfo"`
	GPULogicalName        string   `json:"gpuLogicalName"`
	GPUClock              string   `json:"gpuClock"`
	GPUConfiguration      string   `json:"gpuConfiguration"`
	AudioDeviceProduct    string   `json:"audioDeviceProduct"`
	AudioDeviceModel      string   `json:"audioDeviceModel"`
	BiosVersion           string   `json:"biosVersion"`
	MotherboardManufacturer string `json:"motherboardManufacturer"`
	MotherboardProductName string `json:"motherboardProductName"`
	MotherboardVersion    string   `json:"motherboardVersion"`
	MotherbaoardSerialName string `json:"motherboardSerialName"`
	MotherboardAssetTag   string   `json:"motherboardAssetTag"`
	SoftwareNames         []string `json:"installedPackages"`
	Memories 			[]map[string]string `json:"memories"`
}

func safeGet(slice []string, index int) string {
	if index < len(slice) {
		return slice[index]
	}
	return "N/A"
}

func main() {
	url := "http://10.1.1.73:3000/home/computers/post-machines" //env

   sys := generalinformation.GetSO()

   version,err := generalinformation.GetVersion()
   if err != nil{
	logToFile(fmt.Sprintln("Erro ao Obter a Versão, função 'GetVersion': ", err))
   }
   
   name := generalinformation.GetHostName()

   distribution, err:= generalinformation.GetDistribution()
   if err != nil{
	logToFile(fmt.Sprintln("Erro ao obter a Distribuioção do SO:", err))
   }

   imac, err:=internetinformation.GetMac()
   if err != nil{
	logToFile(fmt.Sprintln("Erro ao obter o mac: ", err))
   }

   date:= generalinformation.GetTime()

   username := generalinformation.GetUser()

   domain,err:=generalinformation.GetDomain()
   if err !=nil{
	logToFile(fmt.Sprintln("Erro ao obter o dominio: ", err))
   }

   ip,err :=internetinformation.GetIP()
   if err!=nil{
	logToFile(fmt.Sprintln("Erro ao obter o IP: ",err))
   }

   manufacturer,err := motherboardinformation.GetManufacturer()
   if err!=nil{
	logToFile(fmt.Sprintln("Erro ao obter a Marca do Equipamento: ",err))
	}

	model,err := motherboardinformation.GetModel()
	if err != nil{
		logToFile(fmt.Sprintln("Erro ao obter o Modelo do Equipamento:", err))
	}

	serialNumber, err := motherboardinformation.GetSerialNumber()
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao obter o número de série: ", err))
	}

	maxCapacity, err := motherboardinformation.GetMaxMem()
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao obter o número Máximo da Memoria RAM: ", err))
	}

	numberOfDevices, err := motherboardinformation.GetSlotDim()
	if err != nil {
		logToFile(fmt.Sprintln("erro ao obter number of devices: ", err))
	}

	numberDevices, err := motherboardinformation.ConvertNumberOfDevices(numberOfDevices)
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao converter a string para inteiro:", err))
		return
	}

	slotNames, err := motherboardinformation.GetMemorySlotNames(numberDevices)
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao obter a quantidade de slot's: ", err))
	}
	// Remove colchetes e divide a string em partes
	cleanedInput := strings.Trim(fmt.Sprint(slotNames), "[]")
	partsNames := strings.Split(cleanedInput, ",")

	memorySizes, err := motherboardinformation.GetMemorySizes(numberDevices)
	if err != nil {
		logToFile(fmt.Sprintln("erro ao obter o tamanho de memória: ", err))
	}
	// Remove colchetes e divide a string em partes
	cleanedInput2 := strings.Trim(fmt.Sprint(memorySizes), "[]")
	partsSizes := strings.Split(cleanedInput2, ",")

	memoriesTypes, err := motherboardinformation.GetMemoryTypes(numberDevices)
	if err != nil {
		logToFile(fmt.Sprintln("erro ao obter o tipo de memória: ", err))
	}
	// Remove colchetes e divide a string em partes
	cleanedInput3 := strings.Trim(fmt.Sprint(memoriesTypes), "[]")
	partsTypes := strings.Split(cleanedInput3, ",")

	memoriesTypeDetails, err := motherboardinformation.GetMemoryTypeDetails(numberDevices)
	if err != nil {
		logToFile(fmt.Sprintln("erro ao obter tipo detalhado de memória: ", err))
	}
	// Remove colchetes e divide a string em partes
	cleanedInput4 := strings.Trim(fmt.Sprint(memoriesTypeDetails), "[]")
	partsTypeDetails := strings.Split(cleanedInput4, ",")

	memoriesSpeedMemory, err := motherboardinformation.GetMemorySpeeds(numberDevices)
	if err != nil {
		logToFile(fmt.Sprintln("erro ao obter a velocidade de memória: ", err))
	}
	// Remove colchetes e divide a string em partes
	cleanedInput5 := strings.Trim(fmt.Sprint(memoriesSpeedMemory), "[]")
	partsSpeed := strings.Split(cleanedInput5, ",")

	memoriesSerialNumber, err := motherboardinformation.GetMemorySerialNumbers(numberDevices)
	if err != nil {
		logToFile(fmt.Sprintln("erro ao obter o serial number de memória: ", err))
	}
	// Remove colchetes e divide a string em partes
	cleanedInput6 := strings.Trim(fmt.Sprint(memoriesSerialNumber), "[]")
	partsSerialNumber := strings.Split(cleanedInput6, ",")

	// Remove colchetes e divide a string em partes
	var memoriesList []map[string]string

	for i := 0; i < len(partsNames); i++ {
		func(i int) {
			defer func() {
				if r := recover(); r != nil {
					logToFile(fmt.Sprintf("Recuperado de um panic: %v", r))
				}
			}()
	
			obj := map[string]string{
				"BankLabel":   safeGet(partsNames, i),
				"Capacity":    safeGet(partsSizes, i),
				"MemoryType":  safeGet(partsTypes, i),
				"TypeDetail":  safeGet(partsTypeDetails, i),
				"Speed":       safeGet(partsSpeed, i),
				"SerialNumber": safeGet(partsSerialNumber, i),
			}
			memoriesList = append(memoriesList, obj)
		}(i)
	}

    values, err := hdinformation.GetHDModel()
	if err != nil{
		logToFile(fmt.Sprintln(err))
	}

	// Extrai os valores entre <>
	 valuesString := hdinformation.ExtractValues(values)

	 hard_disk_model := strings.Join(valuesString, ", ")

	// Obtém a lista de dispositivos ada conectados
	disks, err := hdinformation.DevicesListADA()

	if err != nil {
		logToFile(fmt.Sprintln("Erro ao listar dispositivos:", err))
	}

	var hard_disk_serial_numbers []string

	// Itera sobre cada disco encontrado e obtém o número de série
	for _, disk := range disks {
		hard_disk_serial_number, err := hdinformation.GetSerialNumber(disk)
		if err != nil {
			logToFile(fmt.Sprintf("Erro ao obter número de série do %s: %v\n", disk, err))
		} else {
			hard_disk_serial_numbers = append(hard_disk_serial_numbers, hard_disk_serial_number)
		}
	}

	var hard_disk_sizes [] string
	sectorSize := 512
	// Itera sobre cada disco encontrado e obtém o tamanho do disco
	for _, disk := range disks {
		bytes_size, err := hdinformation.GetLBA48(disk)
		if err != nil {
			logToFile(fmt.Sprintf("Erro ao obter número de série do %s: %v\n", disk, err))
		} else {
			sizeBytes := bytes_size * sectorSize
			sizeGB := float64(sizeBytes) / (1024 * 1024 * 1024)
			str := strconv.FormatFloat(sizeGB, 'f', -1, 64)
			hard_disk_sizes = append(hard_disk_sizes, str)
		}
	}

	var sata_versions [] string
	// Itera sobre cada disco encontrado e obtém o tamanho do disco
	for _, disk := range disks {
		sata_version, err := hdinformation.SataVersion(disk)
		if err != nil {
			logToFile(fmt.Sprintf("Erro ao obter número de série do %s: %v\n", disk, err))
		} else {
			sata_versions = append(sata_versions, sata_version)
		}
	}

	// Obtém a arquitetura do CPU
	arch := runtime.GOARCH

	var cpu_operation_mode string

	if arch == "amd64" || arch == "arm64" {
		cpu_operation_mode = "64 bits"
	}

	numCPUs := runtime.NumCPU()

	cpu_model_name, err := cpuinformation.GetCPUInfo()
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	cpu_thread, err := cpuinformation.GetThread()
	if err != nil {
		fmt.Sprintln(err)
	}

	cpu_max_mhz, err := cpuinformation.GetMaxMHz()
	if err != nil {
		fmt.Sprintln(err)
	}

	cpu_min_mhz, err := cpuinformation.GetCPUMinMHz()
	if err != nil {
		fmt.Sprintln(err)
	}

	motherboard_manufacturer, err := mbinformation.GetMotherboardManufacturer()
	if err !=nil{
		fmt.Sprintln(err)
	}

	motherboard_pd, err := mbinformation.GetMotherboardPD()
	if err !=nil{
		fmt.Sprintln(err)
	}

	motherboard_version,err := mbinformation.GetMotherboardVersion()
	if err !=nil{
		fmt.Sprintln(err)
	}

	motherboard_sn, err := mbinformation.GetMotherSN()
	if err !=nil{
		fmt.Sprintln(err)
	}

	motherboard_asset_tag, err := mbinformation.GetMotherboardAssetTag()
	if err !=nil{
		fmt.Sprintln(err)
	}

    jsonData := Data{
		System: sys, 
		Name: name,
		Distribution:distribution,
		PlatformVersion:version,
		MacAddress: imac, 
		InsertionDate: date, 
		CurrentUser:username, 
		Domain:domain, 
		IP:ip, 
		Manufacturer:manufacturer, 
		Model: model, 
		SerialNumber: serialNumber, 
		MaxCapacityMemory: maxCapacity, 
		NumberOfDevices: numberOfDevices, 
		Memories: memoriesList,
		HardDiskModel:hard_disk_model, 
		HardDiskSerialNumber: strings.Join(hard_disk_serial_numbers, ", "),
		HardDiskUserCapacity: strings.Join(hard_disk_sizes, ", "),
		HardDiskSataVersion: strings.Join(sata_versions, " | "),
		CPUArchitecture: arch, 
		CPUOperationMode:cpu_operation_mode, 
		CPUModelName: cpu_model_name, 
		CPUThread:cpu_thread,
		CPUCore: numCPUs,
		CPUMaxMHz:cpu_max_mhz,
		CPUMinMHz: cpu_min_mhz,
		MotherboardManufacturer:motherboard_manufacturer,
		MotherboardProductName: motherboard_pd,
		MotherboardVersion:motherboard_version ,
		MotherbaoardSerialName: motherboard_sn,
		MotherboardAssetTag: motherboard_asset_tag,
	}

    requestBody, err := json.Marshal(jsonData)
    if err != nil{
        fmt.Println(err)
    }

    resp, erro := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
    if erro != nil{
        fmt.Println(erro)
    }

    defer resp.Body.Close()  
}

