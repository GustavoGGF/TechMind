package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/host"
)

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
	FirstSlotDim          string   `json:"firstSlotDim"`
	SecondSlotDim         string   `json:"secondSlotDim"`
	ThirdSlotDim          string   `json:"thirdSlotDim"`
	FourthSlotDim         string   `json:"fourthSlotDim"`
	FirstSize             string   `json:"firstSize"`
	SecondSize            string   `json:"secondSize"`
	ThirdSize             string   `json:"thirdSize"`
	FourthSize            string   `json:"fourthSize"`
	FirstType             string   `json:"firstType"`
	SecondType            string   `json:"secondType"`
	ThirdType             string   `json:"thirdType"`
	FourthType            string   `json:"fourthType"`
	FirstTypeDetails      string   `json:"firstTypeDetails"`
	SecondTypeDetails     string   `json:"secondTypeDetails"`
	ThirdTypeDetails      string   `json:"thirdTypeDetails"`
	FourthTypeDetails     string   `json:"fourthTypeDetails"`
	FirstSpeedMemory      string   `json:"firstSpeedMemory"`
	SecondSpeedMemory     string   `json:"secondSpeedMemory"`
	ThirdSpeedMemory      string   `json:"thirdSpeedMemory"`
	FourthSpeedMemory     string   `json:"fourthSpeedMemory"`
	FirstSerialNumber     string   `json:"firstSerialNumber"`
	SecondSerialNumber    string   `json:"secondSerialNumber"`
	ThirdSerialNumber     string   `json:"thirdSerialNumber"`
	FourthSerialNumber    string   `json:"fourthSerialNumber"`
	HardDiskModel         string   `json:"hardDiskModel"`
	HardDiskSerialNumber  string   `json:"hardDiskSerialNumber"`
	HardDiskUserCapacity  string   `json:"hardDiskUserCapacity"`
	HardDiskSataVersion   string   `json:"hardDiskSataVersion"`
	CPUArchitecture       string   `json:"cpuArchitecture"`
	CPUOperationMode      string   `json:"cpuOperationMode"`
	CPUS                  string   `json:"cpus"`
	CPUVendorID           string   `json:"cpuVendorID"`
	CPUModelName          string   `json:"cpuModelName"`
	CPUThread             string   `json:"cpuThread"`
	CPUCore               string   `json:"cpuCore"`
	CPUSocket             string   `json:"cpuSocket"`
	CPUMaxMHz             string   `json:"cpuMaxMHz"`
	CPUMinMHz             string   `json:"cpuMinMHz"`
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
}

func bytesEqual(b []byte) bool {
    for _, v := range b {
        if v != 0 {
            return false
        }
    }
    return true
}

func getDomain() (string, error) {
    cmd := exec.Command("hostname", "--domain")
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("erro ao obter domínio: %w", err)
    }
    domain := strings.TrimSpace(string(output))
    return domain, nil
}

func getIPAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("erro ao obter interfaces de rede: %v", err)
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Printf("erro ao obter endereços para interface %s: %v", iface.Name, err)
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ipNet.IP.IsGlobalUnicast() {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("nenhum endereço IP global encontrado")
}

func getSystemManufacturer() (string, error) {
	cmd := exec.Command("sudo","dmidecode", "-s", "system-manufacturer")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
	}
	manufacturer := strings.TrimSpace(string(output))
	if manufacturer == "" {
		return "", fmt.Errorf("marca do sistema não encontrada")
	}
	return manufacturer, nil
}

func getSystemProduct() (string, error) {
	cmd := exec.Command("sudo","dmidecode", "-s", "system-product-name")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
	}
	manufacturer := strings.TrimSpace(string(output))
	if manufacturer == "" {
		return "", fmt.Errorf("marca do sistema não encontrada")
	}
	return manufacturer, nil
}

func getSerialNumber() (string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-s", "system-serial-number")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
	}
	serialNumber := strings.TrimSpace(string(output))
	if serialNumber == "" {
		return "", fmt.Errorf("número de série não encontrado")
	}
	return serialNumber, nil
}

func getMaximumCapacity() (string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Maximum Capacity:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", fmt.Errorf("maximum capacity não encontrada")
}


func getNumberOfDevices() (string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Number Of Devices:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", fmt.Errorf("number of devices não encontrada")
}

func getMemorySlotName(slotNumber int) (string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	slotCount := 0

	for _, line := range lines {
		if strings.Contains(line, "Locator:") {
			if slotCount == slotNumber {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
			slotCount++
		}
	}

	return "", fmt.Errorf("slot de memória número %d não encontrado", slotNumber)
}

func getMemorySize(slotNumber int) (string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	slotCount := 0

	for _, line := range lines {
		if strings.Contains(line, "Size:") {
			if slotCount == slotNumber {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
			slotCount++
		}
	}

	return "", fmt.Errorf("slot de memória número %d não encontrado", slotNumber)
}

func getMemoryType(slotNumber int) (string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	slotCount := 0

	for _, line := range lines {
		if strings.Contains(line, "Type:") {
			if slotCount == slotNumber {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
			slotCount++
		}
	}

	return "", fmt.Errorf("slot de memória número %d não encontrado", slotNumber)
}

func getMemoryTypeDetails(slotNumber int) (string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	slotCount := 0

	for _, line := range lines {
		if strings.Contains(line, "Type Detail:") {
			if slotCount == slotNumber {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
			slotCount++
		}
	}

	return "", fmt.Errorf("slot de memória número %d não encontrado", slotNumber)
}

func getMemorySpeed(slotNumber int) (string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	slotCount := 0

	for _, line := range lines {
		if strings.Contains(line, "Speed:") {
			if slotCount == slotNumber {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
			slotCount++
		}
	}

	return "", fmt.Errorf("slot de memória número %d não encontrado", slotNumber)
}

func getMemorySerialNumber(slotNumber int) (string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	slotCount := 0

	for _, line := range lines {
		if strings.Contains(line, "Serial Number:") {
			if slotCount == slotNumber {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
			slotCount++
		}
	}

	return "", fmt.Errorf("slot de memória número %d não encontrado", slotNumber)
}

func getHDModel(device string) (string, error) {
	// Executa o comando smartctl
	cmd := exec.Command("sudo", "smartctl", "-i", device)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar o comando: %w", err)
	}

	// Captura a saída e converte para string
	output := out.String()

	// Divide a saída em linhas
	lines := strings.Split(output, "\n")

	// Procura pela linha que contém "Device Model"
	for _, line := range lines {
		if strings.Contains(line, "Device Model") {
			// Extrai e retorna o modelo do dispositivo
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				deviceModel := strings.TrimSpace(parts[1])
				return deviceModel, nil
			}
		}
	}

	return "", fmt.Errorf("device model não encontrado")
}

func getHDSerialModel(device string) (string, error) {
	// Executa o comando smartctl
	cmd := exec.Command("sudo", "smartctl", "-i", device)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar o comando: %w", err)
	}

	// Captura a saída e converte para string
	output := out.String()

	// Divide a saída em linhas
	lines := strings.Split(output, "\n")

	// Procura pela linha que contém "Serial Number"
	for _, line := range lines {
		if strings.Contains(line, "Serial Number") {
			// Extrai e retorna o modelo do dispositivo
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				deviceModel := strings.TrimSpace(parts[1])
				return deviceModel, nil
			}
		}
	}

	return "", fmt.Errorf("device model não encontrado")
}

func getHDUserCapacity(device string) (string, error) {
	// Executa o comando smartctl
	cmd := exec.Command("sudo", "smartctl", "-i", device)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar o comando: %w", err)
	}

	// Captura a saída e converte para string
	output := out.String()

	// Divide a saída em linhas
	lines := strings.Split(output, "\n")

	// Procura pela linha que contém "Serial Number"
	for _, line := range lines {
		if strings.Contains(line, "User Capacity") {
			// Extrai e retorna o modelo do dispositivo
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				deviceModel := strings.TrimSpace(parts[1])
				return deviceModel, nil
			}
		}
	}

	return "", fmt.Errorf("device model não encontrado")
}

func getHDSataVersion(device string) (string, error) {
	// Executa o comando smartctl
	cmd := exec.Command("sudo", "smartctl", "-i", device)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar o comando: %w", err)
	}

	// Captura a saída e converte para string
	output := out.String()

	// Divide a saída em linhas
	lines := strings.Split(output, "\n")

	// Procura pela linha que contém "Serial Number"
	for _, line := range lines {
		if strings.Contains(line, "SATA Version is") {
			// Extrai e retorna o modelo do dispositivo
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				deviceModel := strings.TrimSpace(parts[1])
				return deviceModel, nil
			}
		}
	}

	return "", fmt.Errorf("device model não encontrado")
}

func getCPUInfo(lineOption string) (string, error) {
	// Executa o comando smartctl
	cmd := exec.Command("sudo", "lscpu")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar o comando: %w", err)
	}

	// Captura a saída e converte para string
	output := out.String()

	// Divide a saída em linhas
	lines := strings.Split(output, "\n")

	// Procura pela linha que contém "Serial Number"
	for _, line := range lines {
		if strings.Contains(line, lineOption) {
			// Extrai e retorna o modelo do dispositivo
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				deviceModel := strings.TrimSpace(parts[1])
				return deviceModel, nil
			}
		}
	}

	return "", fmt.Errorf("device model não encontrado")
}

func getGPUProduct(lineOption string) (string, error) {
	// Executa o comando smartctl
	cmd := exec.Command("sudo", "lshw", "-c", "video")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar o comando: %w", err)
	}

	// Captura a saída e converte para string
	output := out.String()

	// Divide a saída em linhas
	lines := strings.Split(output, "\n")

	// Procura pela linha que contém "Serial Number"
	for _, line := range lines {
		if strings.Contains(line, lineOption) {
			// Extrai e retorna o modelo do dispositivo
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				deviceModel := strings.TrimSpace(parts[1])
				return deviceModel, nil
			}
		}
	}

	return "", fmt.Errorf("device model não encontrado")
} 

func getAudioDevices(lineOption string) (string, error) {
	// Executa o comando smartctl
	cmd := exec.Command("aplay", "-l")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar o comando: %w", err)
	}

	// Captura a saída e converte para string
	output := out.String()

	// Divide a saída em linhas
	lines := strings.Split(output, "\n")

	// Procura pela linha que contém "Serial Number"
	for _, line := range lines {
		if strings.Contains(line, lineOption) {
			// Extrai e retorna o modelo do dispositivo
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				deviceModel := strings.TrimSpace(parts[1])
				return deviceModel, nil
			}
		}
	}

	return "", fmt.Errorf("device model não encontrado")
}

func getAudioDeviceModel(lineOption string) (string, error) {
	// Executa o comando aplay -l para listar os dispositivos de áudio
	cmd := exec.Command("aplay", "-l")
	var out bytes.Buffer
	cmd.Stdout = &out

	// Executa o comando e verifica por erros
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("erro ao executar o comando aplay -l: %w", err)
	}

	// Captura a saída e converte para string
	output := out.String()

	// Divide a saída em linhas
	lines := strings.Split(output, "\n")

	// Expressão regular para capturar o modelo do dispositivo de áudio
	re := regexp.MustCompile(`card \d+: .*device \d+: (.*)`)

	for _, line := range lines {
		// Procura pelo padrão de linha desejado
		if strings.Contains(line, lineOption) {
			// Extrai o modelo do dispositivo usando a expressão regular
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				deviceModel := matches[1]
				return deviceModel, nil
			}
		}
	}

	return "", fmt.Errorf("modelo de dispositivo não encontrado para '%s'", lineOption)
}

func extractValueFromBrackets(input string) (string, error) {
	// Expressão regular para encontrar o valor entre colchetes
	re := regexp.MustCompile(`\[([^]]+)\]`)

	// Encontra todas as correspondências na string de entrada
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		return "", fmt.Errorf("não foi possível encontrar valor entre colchetes")
	}

	// Retorna o valor encontrado entre colchetes
	return matches[1], nil
}

// Função que executa o comando dmidecode e retorna a linha contendo "SMBIOS"
func getSMBIOSVersion() (string, error) {
	// Executa o comando dmidecode
	cmd := exec.Command("sudo", "dmidecode", "-t", "baseboard")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "SMBIOS") {
			return line, nil
		}
	}

	return "", fmt.Errorf("slot de memória número não encontrado")
}

func getSMBIOSInfo(info string) (string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "baseboard")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, info) {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", fmt.Errorf("number of devices não encontrada")
}

// Função para listar os pacotes instalados no sistema Linux usando dpkg --get-selections
func listInstalledPackages() (string, error) {
	// Comando a ser executado
	cmd := exec.Command("dpkg", "--get-selections")

	// Captura a saída do comando
	var out bytes.Buffer
	cmd.Stdout = &out

	// Executa o comando
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar dpkg --get-selections: %v", err)
	}

	// Filtra os pacotes instalados e retorna como uma string separada por vírgulas
	installedPackages := filterInstalledPackages(out.String())

	return installedPackages, nil
}

// Função para filtrar os pacotes instalados (ignorando os desinstalados) e retornar como uma string separada por vírgulas
func filterInstalledPackages(output string) string {
	var installedPackages []string 
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		// Verifica se a linha não está vazia
		if line != "" {
			// Separa a linha pelos espaços
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				packageName := fields[0]
				installedPackages = append(installedPackages, packageName)
			}
		}
	}

	// Retorna os pacotes instalados como uma string separada por vírgulas
	return strings.Join(installedPackages, ", ")
}

func main() {
    sys := runtime.GOOS

    infos, err:= host.Info()
    if err != nil{
        fmt.Println("Error")
    }

    url := "http://10.1.1.73:3000/home/computers/post-machines" //env

    name, err := os.Hostname()

    if err != nil {
        fmt.Println(err)
    }
    interfaces, err := net.Interfaces()
    if err != nil {
        fmt.Println("Erro ao obter interfaces de rede:", err)
        return
    }

    ifaceInt:= ""
    imac := ""

    for _, iface := range interfaces {
        if iface.Flags&net.FlagUp != 0 && !bytesEqual(iface.HardwareAddr) {
            ifaceInt = iface.Name
            imac = iface.HardwareAddr.String()
            break
        }
    }

    date_now := time.Now()

    formated_date := date_now.Format("2006-01-02 15:04")

    user, err := user.Current()
    if err != nil {
        log.Fatalf("Erro ao obter o usuário atual: %v", err)
    }

    username := user.Username

    version := infos.PlatformVersion

    domain, err := getDomain()
    if err != nil {
        log.Printf("Erro ao obter domínio: %v", err)
    }

    ip, err := getIPAddress()
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

    manufacturer, err := getSystemManufacturer()
	if err != nil {
		log.Fatalf("Erro ao obter marca do sistema: %v", err)
	}

    model, err := getSystemProduct()
	if err != nil {
		log.Fatalf("Erro ao obter marca do sistema: %v", err)
	}

    serialNumber, err := getSerialNumber()
	if err != nil {
		log.Fatalf("Erro ao obter o número de série: %v", err)
	}

    maxCapacity, err := getMaximumCapacity()
	if err != nil {
		log.Fatalf("Erro ao obter Maximum Capacity: %v", err)
	}

    numberOfDevices, err := getNumberOfDevices()
	if err != nil {
		log.Fatalf("erro ao obter number of devices: %v", err)
	}

	firstSlotName, err := getMemorySlotName(0)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	secondSlotName, err := getMemorySlotName(2)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	thirdSlotName, err := getMemorySlotName(4)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	fourSlotName, err := getMemorySlotName(6)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	firstSize, err := getMemorySize(0)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	secondSize, err := getMemorySize(1)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	thirdSize, err := getMemorySize(2)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	fourthSize, err := getMemorySize(3)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	firstType, err := getMemoryType(0)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	secondType, err := getMemoryType(3)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	thirdType, err := getMemoryType(3)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	fourthType, err := getMemoryType(3)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}
	
	firstTypeDetails, err := getMemoryTypeDetails(0)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	secondTypeDetails, err := getMemoryTypeDetails(1)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	thirdTypeDetails, err := getMemoryTypeDetails(2)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	fourthTypeDetails, err := getMemoryTypeDetails(3)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	firstSpeedMemory, err := getMemorySpeed(0)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	secondSpeedMemory, err := getMemorySpeed(1)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	thirdSpeedMemory, err := getMemorySpeed(2)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	fourthSpeedMemory, err := getMemorySpeed(3)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	first_serial_number, err := getMemorySerialNumber(0)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	second_serial_number, err := getMemorySerialNumber(1)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	third_serial_number, err := getMemorySerialNumber(2)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	}

	fourth_serial_number, err := getMemorySerialNumber(3)
	if err != nil {
		log.Fatalf("erro ao obter o nome do primeiro slot de memória: %v", err)
	} 

	hard_disk_model, err := getHDModel("/dev/sda")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	hard_disk_serial_number, err := getHDSerialModel("/dev/sda")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	} 

	hard_disk_user_capacity, err := getHDUserCapacity("/dev/sda")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	hard_disk_sata_version, err := getHDSataVersion("/dev/sda")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	cpu_architecture, err := getCPUInfo("Architecture:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	cpu_operation_mode, err := getCPUInfo("CPU op-mode(s):")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	} 

	cpus, err := getCPUInfo("CPU(s):")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	cpu_vendor_id, err := getCPUInfo("Vendor ID:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	cpu_model_name, err := getCPUInfo("Model name:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	cpu_thread, err := getCPUInfo("Thread(s) per core:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	cpu_core, err := getCPUInfo("Core(s) per socket:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	cpu_socket, err := getCPUInfo("Socket(s):")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	cpu_max_mhz, err := getCPUInfo("CPU max MHz:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	cpu_min_mhz, err := getCPUInfo("CPU min MHz:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	gpu_product, err := getGPUProduct("product:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	gpu_vendor_id, err := getGPUProduct("vendor:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	gpu_bus_info, err := getGPUProduct("bus info:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	gpu_logical_name, err := getGPUProduct("logical name:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	gpu_clock, err := getGPUProduct("clock:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	gpu_configuration, err := getGPUProduct("configuration:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}  

	audio_device_product_first_value, err := getAudioDevices("card 0:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	} 

	audio_device_product, err := extractValueFromBrackets(audio_device_product_first_value)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	audio_device_model_first_value, err := getAudioDeviceModel("device 0:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	audio_device_model, err := extractValueFromBrackets(audio_device_model_first_value)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}
	
	bios_version, err := getSMBIOSVersion()
	if err != nil {
		fmt.Println("Erro:", err)
		return
	} 

	motherboard_manufacturer, err := getSMBIOSInfo("Manufacturer:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	motherboard_product_name, err := getSMBIOSInfo("Product Name:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	motherboard_version, err := getSMBIOSInfo("Version:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	motherboard_serial_name, err := getSMBIOSInfo("Serial Number:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	motherboard_asset_tag, err := getSMBIOSInfo("Asset Tag:")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	// Chama a função para listar os pacotes instalados
	installedPackages, err := listInstalledPackages()
	if err != nil {
		log.Fatalf("Erro ao listar os pacotes instalados: %v", err)
	}

	// Convertendo a string para um slice de strings
    packageSlice := strings.Split(installedPackages, ",")

	// Converte a string de pacotes instalados para um slice de strings
    jsonData := Data{
		System: sys, 
		Name: name, 
		Distribution: infos.Platform, 
        InterfaceInternet: ifaceInt, 
		MacAddress: imac, 
		InsertionDate: formated_date, 
        CurrentUser:username, 
		PlatformVersion:version, 
		Domain:domain, 
		IP:ip, 
        Manufacturer:manufacturer, 
		Model: model, 
		SerialNumber: serialNumber, 
        MaxCapacityMemory: maxCapacity, 
		NumberOfDevices: numberOfDevices, 
		FirstSlotDim: firstSlotName, 
		SecondSlotDim: secondSlotName, 
		ThirdSlotDim: thirdSlotName, 
		FourthSlotDim: fourSlotName, 
		FirstSize: firstSize,
		SecondSize: secondSize,
		ThirdSize: thirdSize,
		FourthSize: fourthSize,
		FirstType: firstType, 
		SecondType: secondType, 
		ThirdType: thirdType, 
		FourthType: fourthType, 
		FirstTypeDetails: firstTypeDetails, 
		SecondTypeDetails: secondTypeDetails, 
		ThirdTypeDetails: thirdTypeDetails, 
		FourthTypeDetails: fourthTypeDetails,
		FirstSpeedMemory: firstSpeedMemory, 
		SecondSpeedMemory: secondSpeedMemory, 
		ThirdSpeedMemory: thirdSpeedMemory, 
		FourthSpeedMemory: fourthSpeedMemory, 
		FirstSerialNumber: first_serial_number, 
		SecondSerialNumber: second_serial_number,
		ThirdSerialNumber: third_serial_number, 
		FourthSerialNumber: fourth_serial_number, 
		HardDiskModel:hard_disk_model, 
		HardDiskSerialNumber: hard_disk_serial_number,
		HardDiskUserCapacity: hard_disk_user_capacity, 
		HardDiskSataVersion: hard_disk_sata_version,
		CPUArchitecture: cpu_architecture, 
		CPUOperationMode:cpu_operation_mode, 
		CPUS: cpus,
		CPUVendorID: cpu_vendor_id, 
		CPUModelName: cpu_model_name, 
		CPUThread: cpu_thread,
		CPUCore: cpu_core,
		CPUSocket:cpu_socket, 
		CPUMaxMHz: cpu_max_mhz, 
		CPUMinMHz: cpu_min_mhz,
		GPUProduct: gpu_product, 
		GPUVendorID: gpu_vendor_id, 
		GPUBusInfo: gpu_bus_info,
		GPULogicalName: gpu_logical_name, 
		GPUClock: gpu_clock, 
		GPUConfiguration: gpu_configuration,
		AudioDeviceProduct: audio_device_product, 
		AudioDeviceModel: audio_device_model,
		BiosVersion: bios_version,
		MotherboardManufacturer:motherboard_manufacturer,
		MotherboardProductName:motherboard_product_name, 
		MotherboardVersion:motherboard_version,
		MotherbaoardSerialName: motherboard_serial_name,
		MotherboardAssetTag:motherboard_asset_tag,
		SoftwareNames:packageSlice,
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

