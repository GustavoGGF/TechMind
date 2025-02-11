package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"

	"techmind/windows/audioinformation"
	"techmind/windows/cpuinformation"
	"techmind/windows/generalinformation"
	"techmind/windows/gpuinformation"
	"techmind/windows/harddiskinformation"
	"techmind/windows/internetinformation"
	"techmind/windows/memoryinformation"
	"techmind/windows/softwareinformation"

	"github.com/kardianos/service"
	"github.com/shirou/gopsutil/host"
)

var (
	Sys                  string
	Hostname             string
	Edition              string
	Formated_Date        string
	MacAddress           string
	CurrentUser          string
	Version              string
	Domain               string
	Ip                   string
	ManuFacturer         string
	Model                string
	SerialNumber         string
	MaxCapacityMemory    string
	MemorySlots          string
	Mem                  string
	MemoryArray          []map[string]interface{}
	ModelHardDisk        string
	SerialNumberHardDisk string
	CapacityHardDisk     string
	HardDiskSerialNumber []string
	Arch                 string
	OperationMode        string
	CPUCount             uint32
	VendorID             string
	ModelName            string
	Threads              uint32
	Cores                uint32
	Sockets              int
	MaxMHz               uint32
	MinMHz               uint32
	GPUProduct           string
	VendordIDGPU         string
	BusInfo              string
	GPULogicalName       string
	Clock                string
	ConfigurationGPU     string
	Product              string
	SMBiosInfo           string
	License              string
	VendorIDGPU          string
	CombinedSoftware     []softwareinformation.InstalledSoftware
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
	SoftwareNames           []softwareinformation.InstalledSoftware `json:"installedPackages"`
	Memories                []map[string]interface{}                `json:"memories"`
	License                 string                                  `json:"license"`
}

// Função que cria um arquivo de log
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

type program struct {
	service.Service
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	// Pega informação do Sistema, como se é Windows, FreeBSD, Linux, etc...
	sys := runtime.GOOS
	Sys = strings.TrimSpace(sys)

	// Pega o nome do computador
	hostname, err := generalinformation.GetComputerName()
	if err != nil {
		fmt.Println("Erro ao obter o nome do computador:", err)
	}
	Hostname = strings.TrimSpace(hostname)

	// Varaivel que armazena informações gerais do windows
	output, err := generalinformation.GetWindowsInfo()
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao obter informações do Windows:", err))
	}

	// Extrai e exibe a Edição do Windows
	edition, err := generalinformation.ExtractWindowsEdition(output)
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao extrair a Edição do Windows:", err))
	}
	Edition = strings.TrimSpace(edition)

	// Pega a data atual e formatada
	date_now := time.Now()
	Formated_Date = date_now.Format("2006-01-02 15:04")

	// Obtem macAddress
	macAddress, err := internetinformation.GetMac()
	if err != nil {
		logToFile(fmt.Sprintln(err))
		return
	}
	MacAddress = macAddress

	// Pega o usuario que esta logado
	currentUser, err := user.Current()
	if err != nil {
		logToFile(fmt.Sprintf("Erro ao obter o usuário atual: %v", err))
	}

	CurrentUser = currentUser.Username

	// Pega diversas informações do computador
	info, err := host.Info()
	if err != nil {
		logToFile(fmt.Sprintf("Error: %v", err))
	}
	// Obtem a versão do SO
	Version = info.PlatformVersion

	// Obtem o dominio como nt-lupatech.com.br
	domain, err := generalinformation.GetDomain()
	if err != nil {
		logToFile(fmt.Sprintln(err))
	}
	Domain = domain

	// Armazena o IP
	ip, err := internetinformation.GetIP()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %v", err))
	}
	Ip = ip

	// Armazena o Manufacturer e o Model
	manufacturer, model, err := generalinformation.GetDeviceBrand()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %v", err))
	}
	ManuFacturer = manufacturer
	Model = model

	// Armazena o Serial Number
	serialNumber, err := generalinformation.GetSerialNumber()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %v", err))
	}
	SerialNumber = serialNumber

	// Armazena a quantidade Máxima de memoria RAM suportada
	maxCapacityMemory, err := memoryinformation.GetMaxMemoryCapacity()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %v", err))
	}
	MaxCapacityMemory = maxCapacityMemory

	// Armazena informações sobre os slot's de memoria
	memorySlots, err := memoryinformation.GetMemorySlots()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %d", err))
	}
	MemorySlots = memorySlots

	// Armazena informações detalhadas sobnre cada memoria
	mem, err := memoryinformation.GetMemoryDetails()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %d", err))
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

	models, err := harddiskinformation.GetHardDiskModel()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %v", err))
	}
	var modelHardDisk string
	for _, model := range models {
		modelHardDisk = model
	}
	ModelHardDisk = modelHardDisk

	hardDiskSerialNumber, err := harddiskinformation.GetHardDiskSerialNumber()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %v", err))
	}
	HardDiskSerialNumber = hardDiskSerialNumber

	for _, serialNumber := range hardDiskSerialNumber {
		SerialNumberHardDisk = serialNumber
	}

	capacities, err := harddiskinformation.GetHardDiskCapacity()
	if err != nil {
		logToFile(fmt.Sprint("Erro:", err))
	}

	for _, capacity := range capacities {
		CapacityHardDisk = fmt.Sprintf("%.2f", capacity)
	}

	arch, err := cpuinformation.GetCPUArchitecture()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU architecture: %v", err))
	}
	Arch = arch

	operationMode, err := cpuinformation.GetCPUOperationMode()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU operation mode: %v", err))
	}
	OperationMode = operationMode

	cpuCount, err := cpuinformation.GetCPUCount()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU count: %v", err))
	}
	CPUCount = cpuCount

	vendorID, err := cpuinformation.GetCPUVendorID()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU Vendor ID: %v", err))
	}
	VendorID = vendorID

	modelName, err := cpuinformation.GetCPUModelName()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU Model Name: %v", err))
	}
	ModelName = modelName

	threads, err := cpuinformation.GetCPUThreads()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU threads: %v", err))
	}
	Threads = threads

	cores, err := cpuinformation.GetCPUCores()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU cores: %v", err))
	}
	Cores = cores

	sockets, err := cpuinformation.GetCPUSockets()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU sockets: %v", err))
	}
	Sockets = sockets

	maxMHz, err := cpuinformation.GetCPUMaxMHz()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU Max MHz: %v", err))
	}
	MaxMHz = maxMHz

	minMHz, err := cpuinformation.GetCPUMinMHz()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU Min MHz: %v", err))
	}
	MinMHz = minMHz

	gpuProduct, err := gpuinformation.GetGPUProduct()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get GPU product: %v", err))
	}
	GPUProduct = gpuProduct

	vendorIDGPU, err := gpuinformation.GetGPUVendorID()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get GPU Vendor ID: %v", err))
	}
	VendorIDGPU = vendorIDGPU

	busInfo, err := gpuinformation.GetGPUBusInfo()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get GPU Bus Info: %v", err))
	}
	BusInfo = busInfo

	gpuLogicalName, err := gpuinformation.GetGPULogicalName()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get GPU Logical Name: %v", err))
	}
	GPULogicalName = gpuLogicalName

	clock, err := gpuinformation.GetGPUClock()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get GPU Clock: %v", err))
	}
	Clock = clock

	horizRes, vertRes, ram, err := gpuinformation.GetGPUConfiguration()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get GPU configuration: %v", err))
	}

	// Formata a string com as informações da GPU
	configurationGPU := fmt.Sprintf("Resolution %dx%d, RAM %d MB", horizRes, vertRes, ram/1024/1024)
	ConfigurationGPU = configurationGPU

	product, err := audioinformation.GetAudioDeviceProduct()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get audio device product: %v", err))
	}
	Product = product

	smbiosInfo, err := generalinformation.GetSMBIOS()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get SMBIOS information: %v", err))
	}
	SMBiosInfo = smbiosInfo

	wmiSoftware, err := softwareinformation.GetInstalledSoftwareFromWMI()
	if err != nil {
		logToFile(fmt.Sprintln("Error querying WMI:", err))
	}

	registrySoftware, err := softwareinformation.GetInstalledSoftwareFromRegistry()
	if err != nil {
		logToFile(fmt.Sprintln("Error querying Registry:", err))
	}

	// Combinar as listas de software, removendo duplicatas
	softwareMap := make(map[string]softwareinformation.InstalledSoftware)
	for _, software := range append(wmiSoftware, registrySoftware...) {
		key := strings.ToLower(software.Name)
		softwareMap[key] = software
	}

	for _, software := range softwareMap {
		CombinedSoftware = append(CombinedSoftware, software)
	}

	license, err := generalinformation.ExtractWindowsLicense(output)
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao extrair a licença do Windows:", err))
	}
	License = license

	// Montando o Json
	jsonData := Data{
		System:               Sys,
		Name:                 Hostname,
		Distribution:         Edition,
		InsertionDate:        Formated_Date,
		MacAddress:           MacAddress,
		CurrentUser:          CurrentUser,
		PlatformVersion:      Version,
		Domain:               Domain,
		IP:                   Ip,
		Manufacturer:         ManuFacturer,
		Model:                Model,
		SerialNumber:         SerialNumber,
		MaxCapacityMemory:    MaxCapacityMemory,
		NumberOfDevices:      MemorySlots,
		Memories:             MemoryArray,
		HardDiskModel:        ModelHardDisk,
		HardDiskSerialNumber: SerialNumberHardDisk,
		HardDiskUserCapacity: CapacityHardDisk,
		CPUArchitecture:      Arch,
		CPUOperationMode:     OperationMode,
		CPUS:                 CPUCount,
		CPUVendorID:          VendorID,
		CPUModelName:         ModelName,
		CPUThread:            Threads,
		CPUCore:              Cores,
		CPUSocket:            Sockets,
		CPUMaxMHz:            MaxMHz,
		CPUMinMHz:            MinMHz,
		GPUProduct:           GPUProduct,
		GPUVendorID:          VendorIDGPU,
		GPUBusInfo:           BusInfo,
		GPULogicalName:       GPULogicalName,
		GPUClock:             Clock,
		GPUConfiguration:     ConfigurationGPU,
		BiosVersion:          SMBiosInfo,
		AudioDeviceProduct:   Product,
		SoftwareNames:        CombinedSoftware,
		License:              License,
	}

	requestBody, err := json.Marshal(jsonData)
	if err != nil {
		logToFile(fmt.Sprintf("Erro ao montar o json: %v", err))
		return
	}

	// Acessar variáveis de ambiente
	url := "http://10.1.1.73:3000/home/computers/post-machines"

	resp, erro := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if erro != nil {
		logToFile(fmt.Sprintf("Erro ao fazer o post: %v", err))
		return
	}

	// Erro response 400 gerar aviso na tela

	defer resp.Body.Close()

	// Ler o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logToFile(fmt.Sprintf("Erro ao ler o corpo da resposta: %v", err))
		return
	}

	// Verificar o código de status da resposta
	if resp.StatusCode != http.StatusOK {
		logToFile(fmt.Sprintf("Resposta com status code %d: %s", resp.StatusCode, body))
		if resp.StatusCode == http.StatusBadRequest {
			logToFile(fmt.Sprintln("Erro: Resposta 400 - Solicitação inválida."))
		} else {
			logToFile(fmt.Sprintf("Erro: Resposta %d\n", resp.StatusCode))
		}
		return
	}
}

func (p *program) Stop(s service.Service) error {
	// Código para parar o serviço
	return nil
}

// Configura o tipo de início do serviço
func configureServiceStartType(serviceName, startType string) error {
	var startTypeFlag string
	switch startType {
	case "auto":
		startTypeFlag = "auto"
	case "manual":
		startTypeFlag = "demand"
	case "disabled":
		startTypeFlag = "disabled"
	default:
		return fmt.Errorf("tipo de início inválido: %s", startType)
	}

	// Usa o comando sc.exe para configurar o tipo de início
	cmd := exec.Command("sc", "config", serviceName, "start=", startTypeFlag)
	return cmd.Run()
}

// Instala o serviço
func installService() {
	svcConfig := &service.Config{
		Name:        "TechMind",
		DisplayName: "TechMind Inventory",
		Description: "Ferramenta de inventário da empresa Lupatech",
	}

	prg := &program{}
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = svc.Install()
	if err != nil {
		log.Fatal("Falha ao instalar o serviço:", err)
	}

	// Configura o serviço para iniciar automaticamente
	err = configureServiceStartType("TechMind", "auto")
	if err != nil {
		log.Fatal("Falha ao configurar o tipo de início do serviço:", err)
	}

	err = svc.Start()
	if err != nil {
		log.Fatal("Falha ao iniciar o serviço: ", err)
	}

	// Configura o serviço para iniciar automaticamente
	cmd := exec.Command("sc", "config", "TechMind", "start=", "auto")
	err = cmd.Run()
	if err != nil {
		log.Fatal("Falha ao configurar o serviço para iniciar automaticamente: ", err)
	}

	log.Println("Servico instalado com sucesso.")
}

// Desinstala o serviço
func uninstallService() {
	svcConfig := &service.Config{
		Name: "TechMind",
	}

	prg := &program{}
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = svc.Uninstall()
	if err != nil {
		log.Fatal("Falha ao desinstalar o serviço:", err)
	}

	log.Println("Servico desinstalado com sucesso.")
}

func main() {

	install := flag.Bool("install", false, "Instalar o serviço")
	uninstall := flag.Bool("uninstall", false, "Desinstalar o serviço")
	flag.Parse()

	if *install {
		installService()
		return
	}

	if *uninstall {
		uninstallService()
		return
	}

	// Cria e executa o serviço
	svcConfig := &service.Config{
		Name:        "TechMind",
		DisplayName: "TechMind Inventory",
		Description: "Ferramenta de inventário da empresa Lupatech",
	}

	prg := &program{}
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	go func(){
		err = svc.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

}
