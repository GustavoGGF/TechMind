package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"sync"
	"time"

	"techmind/windows/pkg/audio"
	"techmind/windows/pkg/cpu"
	"techmind/windows/pkg/gpu"
	"techmind/windows/pkg/memory"
	"techmind/windows/pkg/network"
	"techmind/windows/pkg/software"
	"techmind/windows/pkg/storage"
	"techmind/windows/pkg/sysinfo"

	"github.com/mitchellh/go-ps"
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
	SoftwareNames           []software.InstalledSoftware  `json:"installedPackages"`
	Memories                []map[string]interface{}                `json:"memories"`
	License                 string                                  `json:"license"`
}

var(
	MemoryArray          []map[string]interface{}
	CombinedSoftware     []software.InstalledSoftware
)

var (
	startingInfoMutex sync.Mutex
	logToFileMutex    sync.Mutex
	sendDataMutex     sync.Mutex
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

func GetGeneralInformation()(string, string, string, string, string, string, string, string, string, string, string){
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
		LogToFile(fmt.Sprintf("Erro ao obter o dominio: %v", err))
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

	smbiosInfo, err := sysinfo.GetSMBIOS()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get SMBIOS information: %v", err))
	}

	license, err := sysinfo.ExtractWindowsLicense(output)
	if err != nil {
		LogToFile(fmt.Sprintln("Erro ao extrair a licença do Windows:", err))
	}

	return sys, hostname, edition, currentUser.Username, version, domain, manufacturer, model, serialNumber, smbiosInfo, license
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

func GetGPUInformation()(string, string, string, string, string, string){
	gpuProduct, err := gpu.GetGPUProduct()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get GPU product: %v", err))
	}

	gpuVendorID, err := gpu.GetGPUVendorID()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get GPU Vendor ID: %v", err))
	}

	busInfo, err := gpu.GetGPUBusInfo()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get GPU Bus Info: %v", err))
	}

	gpuLogicalName, err := gpu.GetGPULogicalName()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get GPU Logical Name: %v", err))
	}

	clock, err := gpu.GetGPUClock()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get GPU Clock: %v", err))
	}

	horizRes, vertRes, ram, err := gpu.GetGPUConfiguration()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get GPU configuration: %v", err))
	}

	// Formata a string com as informações da GPU
	configurationGPU := fmt.Sprintf("Resolution %dx%d, RAM %d MB", horizRes, vertRes, ram/1024/1024)

	return gpuProduct, gpuVendorID, busInfo, gpuLogicalName, clock, configurationGPU
}

func GetAudioInformation()(string){
	product, err := audio.GetAudioDeviceProduct()
	if err != nil {
		LogToFile(fmt.Sprintf("Failed to get audio device product: %v", err))
	}

	return product
}

func GetSoftwareInformation()([]software.InstalledSoftware){
	wmiSoftware, err := software.GetInstalledSoftwareFromWMI()
	if err != nil {
		LogToFile(fmt.Sprintln("Error querying WMI:", err))
	}

	registrySoftware, err := software.GetInstalledSoftwareFromRegistry()
	if err != nil {
		LogToFile(fmt.Sprintln("Error querying Registry:", err))
	}

	// Combinar as listas de software, removendo duplicatas
	softwareMap := make(map[string]software.InstalledSoftware)
	for _, software := range append(wmiSoftware, registrySoftware...) {
		key := strings.ToLower(software.Name)
		softwareMap[key] = software
	}

	for _, software := range softwareMap {
		CombinedSoftware = append(CombinedSoftware, software)
	}

	return CombinedSoftware
}

func StartingInformationGathering() (Data, string){
	sys, hostname, edition, currentUser, version, domain, manuFacturer, model,serialNumber, smbiosInfo, license := GetGeneralInformation()
	
	// Pega a data atual e formatada
	dateNow := time.Now()
	formatedDate := dateNow.Format("2006-01-02 15:04")

	macAddress, ip := GetNetWorkInformation()

	maxCapacityMemory, memorySlots, MemoryArray := GetMemoryInformation()

	hdModel, hdSerialNumber, hdCapacity := GetHardDiskInformatin()

	arch, operationMode, cpuCount, vendorID, modelName, threads, cores, sockets, maxMHz, minMHz := GetCpuInformation()

	gpuProduct, gpuVendorID, busInfo, gpuLogicalName, clock, configurationGPU := GetGPUInformation()

	productAudio := GetAudioInformation()

	combinedSoftware := GetSoftwareInformation()

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
		GPUProduct:           gpuProduct,
		GPUVendorID:          gpuVendorID,
		GPUBusInfo:           busInfo,
		GPULogicalName:       gpuLogicalName,
		GPUClock:             clock,
		GPUConfiguration:     configurationGPU,
		BiosVersion:          smbiosInfo,
		AudioDeviceProduct:   productAudio,
		SoftwareNames:        combinedSoftware,
		License:	license,
	}

	return jsonData, ""
}

func SendSystemData(jsonPost Data){
	requestBody, err := json.Marshal(jsonPost)
	if err != nil {
		LogToFile(fmt.Sprintf("Erro ao montar o json: %v", err))
		return
	}

	url := "https://techmind.lupatech.com.br/home/computers/post-machines"

	LogToFile(fmt.Sprintln("URL: ", url))

	// Cria um transporte com verificação desabilitada
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	
	// Cria um cliente HTTP com esse transporte customizado
	client := &http.Client{Transport: transport}
	
	resp, erro := client.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if erro != nil {
		LogToFile(fmt.Sprintf("Erro ao fazer o post: %v", erro))
		return
	}

	defer resp.Body.Close()

	// Ler o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LogToFile(fmt.Sprintf("Erro ao ler o corpo da resposta: %v", err))
		return
	}

	if resp.StatusCode != http.StatusOK {
		LogToFile(fmt.Sprintf("Resposta com status code %d: %s", resp.StatusCode, body))
		if resp.StatusCode == http.StatusBadRequest {
			LogToFile(fmt.Sprintln("Erro: Resposta 400 - Solicitação inválida."))
		} else {
			LogToFile(fmt.Sprintf("Erro: Resposta %d\n", resp.StatusCode))
		}
		return
	}
}

func StartUpdateListener(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		LogToFile(fmt.Sprintf("Erro ao abrir porta %s: %v", port, err))
		return
	}
	LogToFile(fmt.Sprintf("Porta %s aberta aguardando conexões...", port))
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			LogToFile(fmt.Sprintf("Erro ao aceitar conexão: %v", err))
			continue
		}
		LogToFile(fmt.Sprintf("Conexão recebida de %s", conn.RemoteAddr().String()))
		conn.Close() // fecha logo após aceitar
	}
}

func KillExistingTechmind() {
	currentPid := os.Getpid()

	processes, err := ps.Processes()
	if err != nil {
		LogToFile("Erro ao listar processos: " + err.Error())
		return
	}

	for _, proc := range processes {
		if strings.EqualFold(proc.Executable(), "techmind.exe") && proc.Pid() != currentPid {
			// Finaliza o processo encontrado
			cmd := exec.Command("taskkill", "/PID", fmt.Sprint(proc.Pid()), "/F")
			err := cmd.Run()
			if err != nil {
				LogToFile(fmt.Sprintf("Erro ao finalizar processo %d: %v", proc.Pid(), err))
			} else {
				LogToFile(fmt.Sprintf("Finalizado processo duplicado: PID %d", proc.Pid()))
			}
		}
	}
}

func main() {
	// 1ª goroutine: coleta e envio de dados
	go func() {
		// Verificando se tem mais instancias da apilicação  rodando e finalizando
		startingInfoMutex.Lock()
		KillExistingTechmind()
		startingInfoMutex.Unlock()

		// Coleta de dados
		startingInfoMutex.Lock()
		dataJson, errorMac := StartingInformationGathering()
		startingInfoMutex.Unlock()

		if errorMac != "" {
			logToFileMutex.Lock()
			LogToFile(errorMac)
			logToFileMutex.Unlock()
			return
		}

		// Montando dados em formato json e mandando para a aplicação web
		sendDataMutex.Lock()
		SendSystemData(dataJson)
		sendDataMutex.Unlock()
	}()

	// 2ª goroutine: abre porta e escuta permanentemente
	go func() {
		StartUpdateListener("9090")
	}()

	select {}
}