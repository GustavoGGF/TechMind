package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"time"
	"techmind/windows/pkg/sysinfo"
	"techmind/windows/pkg/network"
	"techmind/windows/pkg/memory"
	"techmind/windows/pkg/storage"
	"techmind/windows/pkg/cpu"

	"github.com/shirou/gopsutil/host"
)

// Struct onde que mandará o JSON
type Data struct {
	System                  string                                  `json:"system"`
	Name                    string                                  `json:"name"`
	Distribution            string                                  `json:"distribution"`
	InsertionDate           string                                  `json:"insertionDate"`
	MacAddress              string                                  `json:"macAddress"`
	CurrentUser             string                                  `json:"currentUser"`
	PlatformVersion         string                                  `json:"platformVersion"`
	Domain                  string                                  `json:"domain"`
	IP                      string                                  `json:"ip"`
	Manufacturer            string                                  `json:"manufacturer"`
	Model                   string                                  `json:"model"`
	SerialNumber            string                                  `json:"serialNumber"`
	MaxCapacityMemory       string                                  `json:"maxCapacityMemory"`
	NumberOfDevices         string                                  `json:"numberOfDevices"`
	HardDiskModel           string                                  `json:"hardDiskModel"`
	HardDiskSerialNumber    string                                  `json:"hardDiskSerialNumber"`
	HardDiskUserCapacity    string                                  `json:"hardDiskUserCapacity"`
	HardDiskSataVersion     string                                  `json:"hardDiskSataVersion"`
	CPUArchitecture         string                                  `json:"cpuArchitecture"`
	CPUOperationMode        string                                  `json:"cpuOperationMode"`
	CPUS                    uint32                                  `json:"cpus"`
	CPUVendorID             string                                  `json:"cpuVendorID"`
	CPUModelName            string                                  `json:"cpuModelName"`
	CPUThread               uint32                                  `json:"cpuThread"`
	CPUCore                 uint32                                  `json:"cpuCore"`
	CPUSocket               int                                     `json:"cpuSocket"`
	CPUMaxMHz               uint32                                  `json:"cpuMaxMHz"`
	CPUMinMHz               uint32                                  `json:"cpuMinMHz"`
	GPUProduct              string                                  `json:"gpuProduct"`
	GPUVendorID             string                                  `json:"gpuVendorID"`
	GPUBusInfo              string                                  `json:"gpuBusInfo"`
	GPULogicalName          string                                  `json:"gpuLogicalName"`
	GPUClock                string                                  `json:"gpuClock"`
	GPUConfiguration        string                                  `json:"gpuConfiguration"`
	AudioDeviceProduct      string                                  `json:"audioDeviceProduct"`
	AudioDeviceModel        string                                  `json:"audioDeviceModel"`
	BiosVersion             string                                  `json:"biosVersion"`
	MotherboardManufacturer string                                  `json:"motherboardManufacturer"`
	MotherboardProductName  string                                  `json:"motherboardProductName"`
	MotherboardVersion      string                                  `json:"motherboardVersion"`
	MotherbaoardSerialName  string                                  `json:"motherboardSerialName"`
	MotherboardAssetTag     string                                  `json:"motherboardAssetTag"`
	SoftwareNames           []map[string]interface{}  `json:"installedPackages"`
	Memories                []map[string]interface{}                `json:"memories"`
	License                 string                                  `json:"license"`
}
// 	SoftwareNames           []softwareinformation.InstalledSoftware `json:"installedPackages"`

var(
	MemoryArray          []map[string]interface{}
)


// Função que cria um arquivo de log
func LogToFile(msg string) {
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

func GetGeneralInformation()(string, string, string, string, string, string, string, string, string){
	// Obtem o SO do equipamento
	sys := sysinfo.GetSys()
	// Pega o nome do computador
	hostname, err := sysinfo.GetComputerName()
	if err != nil {
		LogToFile(fmt.Sprintln("Erro ao obter o nome do computador:", err))
	}

	// Varaivel que armazena informações gerais do windows
	output, err := sysinfo.GetWindowsInfo()
	if err != nil {
		LogToFile(fmt.Sprintln("Erro ao obter informações do Windows:", err))
	}

	// Extrai a Edição do Windows
	edition, err := sysinfo.ExtractWindowsEdition(output)
	if err != nil {
		LogToFile(fmt.Sprintln("Erro ao extrair a Edição do Windows:", err))
	}

	// Pega o usuario que esta logado
	currentUser, err := user.Current()
	if err != nil {
		LogToFile(fmt.Sprintf("Erro ao obter o usuário atual: %v", err))
	}

	// Pega diversas informações do computador
	info, err := host.Info()
	if err != nil {
		LogToFile(fmt.Sprintf("Error ao obter host.info: %v", err))
	}
	 
	// Obtem a versão do SO
	version := info.PlatformVersion

	// Obtem o dominio como nt-lupatech.com.br
	domain, err := sysinfo.GetDomain()
	if err != nil {
		LogToFile(fmt.Sprintln("Erro ao obter o dominio: %v", err))
	}

	// obtem o Manufacturer e o Model
	manufacturer, model, err := sysinfo.GetDeviceBrand()
	if err != nil {
		LogToFile(fmt.Sprintf("Erro ao obter o modelo e fabricante: %v", err))
	}

	// Obtem o Serial Number
	serialNumber, err := sysinfo.GetSerialNumber()
	if err != nil {
		LogToFile(fmt.Sprintf("Erro ao obter o SerialNumber do equipamento: %v", err))
	}

	return sys, hostname, edition, currentUser.Username, version, domain, manufacturer, model, serialNumber
}

func GetNetWorkInformation()(string, string){
	// Obtem macAddress
	macAddress, err := network.GetMac()
	if err != nil {
		LogToFile(fmt.Sprintf("Erro fatal!! Não foi possivel obter o MAC ADDRESS: %s", err))
		return "", ""
	}

	// Obtem o IP
	ip, err := network.GetIP()
	if err != nil {
		LogToFile(fmt.Sprintf("Erro ao tentar obter o ip: %v", err))
	}

	return macAddress, ip
}

func GetMemoryInformation()(string, string,  []map[string]interface{}){
	// Obtem a quantidade Máxima de memoria RAM suportada
	maxCapacityMemory, err := memory.GetMaxMemoryCapacity()
	if err != nil {
		LogToFile(fmt.Sprintf("Erro ao obter a capacidade máxima da memoria: %v", err))
	}

	// Obtem informações sobre os slot's de memoria
	memorySlots, err := memory.GetMemorySlots()
	if err != nil {
		LogToFile(fmt.Sprintf("Erro ao obter a quantidade de slot's da memoria: %d", err))
	}

	// Armazena informações detalhadas sobnre cada memoria
	mem, err := memory.GetMemoryDetails()
	if err != nil {
		LogToFile(fmt.Sprintf("Erro ao tentar obter informações detalhadas da memoria RAM: %d", err))
	}

	for _, memory := range mem {
		memoryInfo := map[string]interface{}{
			"BankLabel":     memory.BankLabel,
			"Capacity":      memory.Capacity,
			"DeviceLocator": memory.DeviceLocator,
			"MemoryType":    memory.MemoryType,
			"TypeDetail":    memory.TypeDetail,
			"Speed":         memory.Speed,
			"SerialNumber":  memory.SerialNumber,
		}
		MemoryArray = append(MemoryArray, memoryInfo)
	}

	return maxCapacityMemory, memorySlots, MemoryArray
}

func GetHardDiskInformatin()(string, string, string){
	modelHardDisk, err := storage.GetHardDiskModel()
	if err != nil {
		LogToFile(fmt.Sprintf("Erro ao tentar obter o modelo do HD: %v", err))
	}

	var hdModel string
	var hdSerialNumber string
	var hdCapacity string

	for _, model := range modelHardDisk {
		hdModel = model
	}

	hardDiskSerialNumber, err := storage.GetHardDiskSerialNumber()
	if err != nil {
		LogToFile(fmt.Sprintf("Erro ao tentar obter o SN do HD: %v", err))
	}

	
	for _, serialNumber := range hardDiskSerialNumber {
		hdSerialNumber = serialNumber
	}

	capacities, err := storage.GetHardDiskCapacity()
	if err != nil {
		LogToFile(fmt.Sprint("Erro ao tentar obter a capacidade do HD:", err))
	}

	for _, capacity := range capacities {
		hdCapacity = fmt.Sprintf("%.2f", capacity)
	}

	return hdModel, hdSerialNumber, hdCapacity
}

func GetCpuInformation()(string, string, uint32, string, string, uint32, uint32, int, uint32, uint32){
	arch, err := cpu.GetCPUArchitecture()
	if err != nil {
		LogToFile(fmt.Sprintf("Falha ao obter arquitetura do processador: %v", err))
	}

	operationMode, err := cpu.GetCPUOperationMode()
	if err != nil {
		LogToFile(fmt.Sprintf("Falha ao obter modo de operação do CPU: %v", err))
	}

	cpuCount, err := cpu.GetCPUCount()
	if err != nil {
		LogToFile(fmt.Sprintf("Falha em obter a quantidade de CPU: %v", err))
	}

	vendorID, err := cpu.GetCPUVendorID()
	if err != nil {
		LogToFile(fmt.Sprintf("Falha ao Obter o Fabricante do CPU: %v", err))
	}

	modelName, err := cpu.GetCPUModelName()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get CPU Model Name: %v", err))
	}

	threads, err := cpu.GetCPUThreads()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get CPU threads: %v", err))
	}

	cores, err := cpu.GetCPUCores()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get CPU cores: %v", err))
	}

	sockets, err := cpu.GetCPUSockets()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get CPU sockets: %v", err))
	}

	maxMHz, err := cpu.GetCPUMaxMHz()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get CPU Max MHz: %v", err))
	}

	minMHz, err := cpu.GetCPUMinMHz()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get CPU Min MHz: %v", err))
	}

	return arch, operationMode, cpuCount, vendorID, modelName, threads, cores, sockets, maxMHz, minMHz
}

func StartingInformationGathering() (Data, string){
	sys, hostname, edition, currentUser, version, domain, manuFacturer, model,serialNumber := GetGeneralInformation()
	
	// Pega a data atual e formatada
	dateNow := time.Now()
	formatedDate := dateNow.Format("2006-01-02 15:04")

	macAddress, ip := GetNetWorkInformation()

	maxCapacityMemory, memorySlots, MemoryArray := GetMemoryInformation()

	hdModel, hdSerialNumber, hdCapacity := GetHardDiskInformatin()

	arch, operationMode, cpuCount, vendorID, modelName, threads, cores, sockets, maxMHz, minMHz := GetCpuInformation()

	if macAddress == ""{
		return Data{}, fmt.Sprintln("Codigo cancelado, falta de macAddress para dar andamento")
	}
	// Montando o Json
	jsonData := Data{
		System:               sys,
		Name:                 hostname,
		Distribution:         edition,
		InsertionDate:        formatedDate,
		MacAddress:           macAddress,
		CurrentUser:          currentUser,
		PlatformVersion:      version,
		Domain:               domain,
		IP:                   ip,
		Manufacturer:         manuFacturer,
		Model:                model,
		SerialNumber:         serialNumber,
		MaxCapacityMemory:    maxCapacityMemory,
		NumberOfDevices:      memorySlots,
		Memories:             MemoryArray,
		HardDiskModel:        hdModel,
		HardDiskSerialNumber: hdSerialNumber,
		HardDiskUserCapacity: hdCapacity,
		CPUArchitecture:      arch,
		CPUOperationMode:     operationMode,
		CPUS:                 cpuCount,
		CPUVendorID:          vendorID,
		CPUModelName:         modelName,
		CPUThread:            threads,
		CPUCore:              cores,
		CPUSocket:            sockets,
		CPUMaxMHz:            maxMHz,
		CPUMinMHz:            minMHz,
			// GPUProduct:           GPUProduct,
			// GPUVendorID:          VendorIDGPU,
			// GPUBusInfo:           BusInfo,
			// GPULogicalName:       GPULogicalName,
			// GPUClock:             Clock,
			// GPUConfiguration:     ConfigurationGPU,
			// BiosVersion:          SMBiosInfo,
			// AudioDeviceProduct:   Product,
			// SoftwareNames:        CombinedSoftware,
			// License:              License,
	}

	return jsonData, ""
}

func main() {
	dataJson, errorMac := StartingInformationGathering()
	if errorMac != ""{
		LogToFile(errorMac)
		return
	}

	fmt.Sprintln(dataJson) 
}