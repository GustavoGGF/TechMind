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
	"runtime"
	"strings"
	"time"

	"techmind/windows/general"
	"techmind/windows/memory"
	"techmind/windows/network"
	"techmind/windows/software"

	"github.com/kardianos/service"
	"golang.org/x/sys/windows/registry"
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
	SoftwareNames           []software.InstalledSoftware `json:"installedPackages"`
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

var (
	MemoryArray          []map[string]interface{}
)

type program struct {
	service.Service
}

func (p *program) Start(s service.Service) error {
	go func() {
		err := run()
		if err != nil {
			logToFile(fmt.Sprintln(err))
		}
	}()

	select{}
}

func run() error{
	// Pega informação do Sistema, como se é Windows, FreeBSD, Linux, etc...
	sys := runtime.GOOS
	sys = strings.TrimSpace(sys)

	// Pega o nome do computador
	hostname, err := general.GetComputerName()
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao obter o nome do computador:", err))
	}
	hostname = strings.TrimSpace(hostname)

	// Chama a função GetWindowsDistribution que retorna a distribuição e a versão do Windows
	distribution, version, err := general.GetWindowsDistribution()
	if err != nil {
		// Se ocorrer um erro ao obter a distribuição e versão, registra a mensagem de erro em um arquivo de log
		logToFile(fmt.Sprintln(err))
	}
	
	// Pega a data atual e formatada
	date_now := time.Now()
	Formated_Date := date_now.Format("2006-01-02 15:04")
			
	// Obtém o endereço MAC da máquina usando a função GetMac do pacote network
	mac, err := network.GetMac()
	if err != nil {
		// Se ocorrer um erro ao obter o MAC, registra a mensagem de erro em um arquivo de log
		logToFile(fmt.Sprintln("Erro ao obter o Mac", err))
		return err // Encerra a execução da função em caso de erro
	}
	
	// Pega o usuario que esta logado
	currentUser, err := general.GetCurrentUser()
	if err != nil {
		logToFile(fmt.Sprintf("Erro ao obter o usuário atual: %v", err))
	}

	// Chama a função GetDomain para obter o nome do domínio do sistema
	domain, err := general.GetDomain()
	if err != nil {
		// Se ocorrer um erro ao obter o domínio, registra a mensagem de erro em um arquivo de log
		logToFile(fmt.Sprintln(err))
	}
	
	// Chama a função GetIP da biblioteca network para obter o endereço IP do sistema
	ip, err := network.GetIP()
	if err != nil {
		// Se ocorrer um erro ao obter o endereço IP, registra a mensagem de erro em um arquivo de log
		logToFile(fmt.Sprintf("Erro: %v", err))
	}
	
	// Chama a função GetDeviceBrand para obter o fabricante e o modelo do dispositivo
	manufacturer, model, err := general.GetDeviceBrand()
	if err != nil {
		// Se ocorrer um erro ao obter o fabricante e o modelo, registra a mensagem de erro em um arquivo de log
		logToFile(fmt.Sprintf("Erro: %v", err))
	}

	// Chama a função GetSerialNumber para obter o número de série do dispositivo
	serialNumber, err := general.GetSerialNumber()
	if err != nil {
		// Se ocorrer um erro ao obter o número de série, registra a mensagem de erro em um arquivo de log
		logToFile(fmt.Sprintf("Erro: %v", err))
	}
	
	// Chama a função GetMaxMemoryCapacity para obter a capacidade máxima de memória
	maxCapacityMemory, err := memory.GetMaxMemoryCapacity()
	if err != nil {
		// Registra um erro caso a função retorne um erro
		logToFile(fmt.Sprintf("Erro: %v", err))
	}
	
	memorySlots, err := memory.GetMemorySlots()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %d", err))
	}
	
	memories, err := memory.GetMemoryDetails()
	if err != nil {
		fmt.Printf("Erro ao obter dados da memoria: %v\n", err)
	}
	
	for _, mem := range memories {
		capacityGB := mem.Capacity / (1024 * 1024 * 1024)
		memoryType := memory.GetMemoryType(mem.MemoryType)
	
		memoryInfo := map[string]interface{}{
			"BankLabel":     mem.BankLabel,
			"Capacity":      capacityGB,
			"DeviceLocator": mem.DeviceLocator,
			"MemoryType":    memoryType,
			"TypeDetail":    mem.TypeDetail, 
			"Speed":         mem.Speed,
			"SerialNumber":  mem.SerialNumber,
		}
		MemoryArray = append(MemoryArray, memoryInfo)
	}
	
	// Montando o Json
	jsonData := Data{
		System: sys,
		Name: hostname,
		MacAddress: mac,
		InsertionDate: Formated_Date,
		Distribution: distribution,
		CurrentUser: currentUser,
		PlatformVersion: version,
		Domain: domain,
		IP: ip,
		Manufacturer: manufacturer,
		Model: model,
		SerialNumber: serialNumber,
		MaxCapacityMemory: maxCapacityMemory,
		NumberOfDevices: memorySlots,
		Memories: MemoryArray,
	}
	
	requestBody, err := json.Marshal(jsonData)
	if err != nil {
		logToFile(fmt.Sprintf("Erro ao montar o json: %v", err))
		return err
	}
	
	// Acessar variáveis de ambiente
	url := "http://10.1.1.73:3000/home/computers/post-machines"
	
	resp, erro := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if erro != nil {
		logToFile(fmt.Sprintf("Erro ao fazer o post: %v", err))
		return err
	}
	
	// Erro response 400 gerar aviso na tela
	defer resp.Body.Close()
	
	// Ler o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logToFile(fmt.Sprintf("Erro ao ler o corpo da resposta: %v", err))
		return err
	}
	
	// Verificar o código de status da resposta
	if resp.StatusCode != http.StatusOK {
		logToFile(fmt.Sprintf("Resposta com status code %d: %s", resp.StatusCode, body))
		if resp.StatusCode == http.StatusBadRequest {
			logToFile(fmt.Sprintln("Erro: Resposta 400 - Solicitação inválida."))
		} else {
			logToFile(fmt.Sprintf("Erro: Resposta %d\n", resp.StatusCode))
		}
		return nil
	}
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Força o encerramento do processo após a execução
	os.Exit(0)
	return nil
}

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

func installService() {
	svcConfig := &service.Config{
		Name:        "TechMind",
		DisplayName: "TechMind Inventory",
		Description: "Ferramenta de inventário da empresa Lupatech",
	}

	prg := &program{}
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		logToFile(fmt.Sprintln(err))
	}

	err = svc.Install()
	if err != nil {
		logToFile(fmt.Sprintln("Falha ao instalar o serviço:", err))
	}

	// Configura o serviço para iniciar automaticamente
	err = configureServiceStartType("TechMind", "auto")
	if err != nil {
		logToFile(fmt.Sprintln("Falha ao configurar o tipo de início do serviço:", err))
	}

	err = svc.Start()
	if err != nil {
		logToFile(fmt.Sprintln("Falha ao iniciar o serviço: ", err))
	}

	// Configura o serviço para iniciar automaticamente
	cmd := exec.Command("sc", "config", "TechMind", "start=", "auto")
	err = cmd.Run()
	if err != nil {
		logToFile(fmt.Sprintln("Falha ao configurar o serviço para iniciar automaticamente: ", err))
	}

	log.Println("Servico instalado com sucesso.")

	select{}
}

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

	// Nome do aplicativo a ser removido do registro
	appName := "TechMind"

	// Abrir a chave do registro onde as entradas de autostart estão
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.WRITE)
	if err != nil {
		fmt.Println("Erro ao abrir o registro:", err)
		return
	}
	defer key.Close()
	
	// Remover a entrada do registro
	err = key.DeleteValue(appName)
	if err != nil {
		fmt.Println("Erro ao excluir o valor do registro:", err)
		return
	} else {
		log.Println("Regedit Excluido.")
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

	p := &program{}
	err := p.Start(nil)
	if err != nil{
		logToFile(fmt.Sprintln("Erro ao iniciar o serviço: ", err))
	}
}