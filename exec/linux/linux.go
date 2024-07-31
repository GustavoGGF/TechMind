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
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/host"
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
	Memories 			[]map[string]string `json:"memories"`
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
    output, err := cmd.CombinedOutput() // Captura a saída e erro juntos
    if err != nil {
        // Verifica se o erro é relacionado a uma opção inválida
        if strings.Contains(string(output), "illegal option") || strings.Contains(string(output), "invalid option") {
            return "", nil // Retorna uma string vazia e sem erro
        }
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
    cmd := exec.Command("sudo", "dmidecode", "-s", "system-manufacturer")
    output, err := cmd.CombinedOutput() // Captura a saída e erro juntos
    
    if err != nil {
        // Verifica se o erro é relacionado a um comando não encontrado
        if strings.Contains(string(output), "command not found") || strings.Contains(err.Error(), "executable file not found in $PATH") {
            return "", nil // Retorna uma string vazia e sem erro
        }
        return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
    }
    
    manufacturer := strings.TrimSpace(string(output))
    
    if manufacturer == "" {
        return "", fmt.Errorf("marca do sistema não encontrada")
    }
    
    return manufacturer, nil
}

func getSystemProduct() (string, error) {
    cmd := exec.Command("sudo", "dmidecode", "-s", "system-product-name")
    output, err := cmd.CombinedOutput() // Captura a saída e erro juntos
    
    if err != nil {
        // Verifica se o erro é relacionado a um comando não encontrado
        if strings.Contains(string(output), "command not found") || strings.Contains(err.Error(), "executable file not found in $PATH") {
            return "", nil // Retorna uma string vazia e sem erro
        }
        return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
    }
    
    product := strings.TrimSpace(string(output))
    
    if product == "" {
        return "", fmt.Errorf("produto do sistema não encontrado")
    }
    
    return product, nil
}

func getSerialNumber() (string, error) {
    cmd := exec.Command("sudo", "dmidecode", "-s", "system-serial-number")
    output, err := cmd.CombinedOutput() // Captura a saída e erro juntos
    
    if err != nil {
        // Verifica se o erro é relacionado a um comando não encontrado
        if strings.Contains(string(output), "command not found") || strings.Contains(err.Error(), "executable file not found in $PATH") {
            return "", nil // Retorna uma string vazia e sem erro
        }
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
    output, err := cmd.CombinedOutput() // Captura a saída e erro juntos

    if err != nil {
        // Verifica se o erro é relacionado a um comando não encontrado
        if strings.Contains(string(output), "command not found") || strings.Contains(err.Error(), "executable file not found in $PATH") {
            return "", nil // Retorna uma string vazia e sem erro
        }
        return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
    }

    lines := strings.Split(string(output), "\n")
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
    output, err := cmd.CombinedOutput() // Captura a saída e erro juntos

    if err != nil {
        // Verifica se o erro é relacionado a um comando não encontrado
        if strings.Contains(string(output), "command not found") || strings.Contains(err.Error(), "executable file not found in $PATH") {
            return "", nil // Retorna uma string vazia e sem erro
        }
        return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
    }

    lines := strings.Split(string(output), "\n")
    for _, line := range lines {
        if strings.Contains(line, "Number Of Devices:") {
            parts := strings.Split(line, ":")
            if len(parts) > 1 {
                return strings.TrimSpace(parts[1]), nil
            }
        }
    }

    // Retorna uma string vazia se "Number Of Devices:" não for encontrado
    return "", nil
}

func convertNumberOfDevices(numberOfDevices string) (int, error) {
    // Remove espaços em branco ao redor da string
    numberOfDevices = strings.TrimSpace(numberOfDevices)

    // Verifica se a string é vazia
    if numberOfDevices == "" {
        return 0, nil
    }

    // Converte a string para um inteiro
    numberDevices, err := strconv.Atoi(numberOfDevices)
    if err != nil {
        return 0, fmt.Errorf("erro ao converter a string para inteiro: %v", err)
    }

    return numberDevices, nil
}

func getMemorySlotNames(numberOfDevices int) ([]string, error) {
    // Executa o comando dmidecode
    cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out // Captura também a saída de erro
    output, err := cmd.CombinedOutput() // Captura a saída e erro juntos

    if err != nil {
        // Verifica se o erro é relacionado a um comando não encontrado
        if strings.Contains(string(output), "command not found") || strings.Contains(err.Error(), "executable file not found in $PATH") {
            return nil, nil // Retorna nil e sem erro
        }
        return nil, fmt.Errorf("erro ao executar dmidecode: %v", err)
    }

    lines := strings.Split(out.String(), "\n")
    slotCount := 0
    var slotNames []string

    for _, line := range lines {
        if strings.Contains(line, "Locator:") {
            if slotCount < numberOfDevices {
                parts := strings.Split(line, ":")
                if len(parts) > 1 {
                    // Adiciona uma vírgula antes do próximo valor
                    slotName := strings.TrimSpace(parts[1])
                    if slotCount > 0 {
                        slotNames = append(slotNames, ", "+slotName)
                    } else {
                        slotNames = append(slotNames, slotName)
                    }
                    slotCount++
                }
            }
        }
    }

    if slotCount == 0 {
        return nil, fmt.Errorf("nenhum slot de memória encontrado")
    }

    return slotNames, nil
}

func getMemorySizes(numberOfDevices int) ([]string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	slotCount := 0
	var sizes []string

	for _, line := range lines {
		if strings.Contains(line, "Size:") {
			if slotCount < numberOfDevices {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					// Adiciona uma vírgula antes do próximo valor, exceto para o primeiro valor
					size := strings.TrimSpace(parts[1])
					if slotCount > 0 {
						sizes = append(sizes, ", "+size)
					} else {
						sizes = append(sizes, size)
					}
					slotCount++
				}
			}
		}
	}

	if slotCount == 0 {
		return nil, fmt.Errorf("nenhum tamanho de memória encontrado")
	}

	return sizes, nil
}

func getMemoryTypes(numberOfDevices int) ([]string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	slotCount := 0
	var types []string

	for _, line := range lines {
		if strings.Contains(line, "Type:") {
			if slotCount < numberOfDevices {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					// Adiciona uma vírgula antes do próximo valor, exceto para o primeiro valor
					typeMem := strings.TrimSpace(parts[1])
					if slotCount > 0 {
						types = append(types, ", "+typeMem)
					} else {
						types = append(types, typeMem)
					}
					slotCount++
				}
			}
		}
	}

	if slotCount == 0 {
		return nil, fmt.Errorf("nenhum tipo de memória encontrado")
	}

	return types, nil
}

func getMemoryTypeDetails(numberOfDevices int) ([]string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	slotCount := 0
	var details []string

	for _, line := range lines {
		if strings.Contains(line, "Type Detail:") {
			if slotCount < numberOfDevices {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					// Adiciona uma vírgula antes do próximo valor, exceto para o primeiro valor
					detailMem := strings.TrimSpace(parts[1])
					if slotCount > 0 {
						details = append(details, ", "+detailMem)
					} else {
						details = append(details, detailMem)
					}
					slotCount++
				}
			}
		}
	}

	if slotCount == 0 {
		return nil, fmt.Errorf("nenhum detalhe de tipo de memória encontrado")
	}

	return details, nil
}

func getMemorySpeeds(numberOfDevices int) ([]string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	slotCount := 0
	var speeds []string

	for _, line := range lines {
		if strings.Contains(line, "Speed:") {
			if slotCount < numberOfDevices {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					// Adiciona uma vírgula antes do próximo valor, exceto para o primeiro valor
					speedsMem := strings.TrimSpace(parts[1])
					if slotCount > 0 {
						speeds = append(speeds, ", "+speedsMem)
					} else {
						speeds = append(speeds, speedsMem)
					}
					slotCount++
				}
			}
		}
	}

	if slotCount == 0 {
		return nil, fmt.Errorf("nenhuma velocidade de memória encontrada")
	}

	return speeds, nil
}

func getMemorySerialNumbers(numberOfDevices int) ([]string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar dmidecode: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	slotCount := 0
	var serialNumbers []string

	for _, line := range lines {
		if strings.Contains(line, "Serial Number:") {
			if slotCount < numberOfDevices {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					// Adiciona uma vírgula antes do próximo valor, exceto para o primeiro valor
					serialNumbersMem := strings.TrimSpace(parts[1])
					if slotCount > 0 {
						serialNumbers = append(serialNumbers, ", "+serialNumbersMem)
					} else {
						serialNumbers = append(serialNumbers, serialNumbersMem)
					}
					slotCount++
				}
			}
		}
	}

	if slotCount == 0 {
		return nil, fmt.Errorf("nenhum número de série de memória encontrado")
	}

	return serialNumbers, nil
}

func getHDModel(device string) (string, error) {
	// Executa o comando smartctl
	cmd := exec.Command("sudo", "smartctl", "-i", device)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar o comando, (verifique se smartmontools esta funcionando), codigo de erro: %w", err)
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

	// Procura pela linha que contém "Vendor" se "Device Model" não for encontrado
	for _, line := range lines {
		if strings.Contains(line, "Vendor") {
			// Extrai e retorna o nome do vendedor
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				vendor := strings.TrimSpace(parts[1])
				return vendor, nil
			}
		}
	}

	return "", fmt.Errorf("device model ou vendor não encontrado")
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
			// Extrai e retorna o número de série do dispositivo
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				serialNumber := strings.TrimSpace(parts[1])
				return serialNumber, nil
			}
		}
	}

	// Retorna uma string vazia se o número de série não for encontrado
	return "", nil
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

	return "", nil
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

	// Procura pela linha que contém "SATA Version is"
	for _, line := range lines {
		if strings.Contains(line, "SATA Version is") {
			// Extrai e retorna a versão SATA do dispositivo
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				sataVersion := strings.TrimSpace(parts[1])
				return sataVersion, nil
			}
		}
	}

	// Retorna uma string vazia se a versão SATA não for encontrada
	return "", nil
}

func getCPUInfo(lineOption string) (string, error) {
	// Executa o comando lscpu
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

	// Procura pela linha que contém o texto especificado em lineOption
	for _, line := range lines {
		if strings.Contains(line, lineOption) {
			// Extrai e retorna o valor associado ao texto especificado
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				info := strings.TrimSpace(parts[1])
				return info, nil
			}
		}
	}

	// Retorna uma string vazia se a informação não for encontrada
	return "", nil
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

	return "", fmt.Errorf("GPU Product não encontrado")
} 

func getAudioDevices(lineOption string) (string, error) {
	// Executa o comando aplay para listar dispositivos de áudio
	cmd := exec.Command("aplay", "-l")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao executar o comando (verifique alsa-utils): %w", err)
	}

	// Captura a saída e converte para string
	output := out.String()

	// Divide a saída em linhas
	lines := strings.Split(output, "\n")

	// Procura pela linha que contém o texto especificado em lineOption
	for _, line := range lines {
		if strings.Contains(line, lineOption) {
			// Extrai e retorna a parte relevante da linha
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				deviceInfo := strings.TrimSpace(parts[1])
				return deviceInfo, nil
			}
		}
	}

	// Retorna uma string vazia se o dispositivo de áudio não for encontrado
	return "", nil
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

	// Retorna uma string vazia se o modelo do dispositivo não for encontrado
	return "", nil
}

func extractValueFromBrackets(input string) (string, error) {
	// Expressão regular para encontrar o valor entre colchetes
	re := regexp.MustCompile(`\[([^]]+)\]`)

	// Encontra todas as correspondências na string de entrada
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		// Retorna uma string vazia e nil para o erro se não encontrar valor entre colchetes
		return "", nil
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

func safeGet(slice []string, index int) string {
	if index < len(slice) {
		return slice[index]
	}
	return "N/A"
}

func main() {
    sys := runtime.GOOS

    infos, err:= host.Info()
    if err != nil{
        logToFile(fmt.Sprintln("Erro ao pegar info: ", err))
    }

    url := "http://10.1.1.73:3000/home/computers/post-machines" //env

    name, err := os.Hostname()

    if err != nil {
        fmt.Println(err)
    }
    interfaces, err := net.Interfaces()
    if err != nil {
		logToFile(fmt.Sprintln("Erro ao obter interfaces de rede:", err))
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
        logToFile(fmt.Sprintln("Erro ao obter o usuário atual: ", err))
    }

    username := user.Username

    version := infos.PlatformVersion

    domain, err := getDomain()
    if err != nil {
        logToFile(fmt.Sprintln("Erro ao obter domínio: ", err))
    }

    ip, err := getIPAddress()
	if err != nil {
		logToFile(fmt.Sprintln("Erro: ao obter o ip ", err))
	}

    manufacturer, err := getSystemManufacturer()
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao obter marca do sistema: ", err))
	}

    model, err := getSystemProduct()
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao obter o modelo do sistema: ", err))
	}

    serialNumber, err := getSerialNumber()
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao obter o número de série: ", err))
	}

    maxCapacity, err := getMaximumCapacity()
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao obter Maximum Capacity Memory: ", err))
	}

    numberOfDevices, err := getNumberOfDevices()
	if err != nil {
		logToFile(fmt.Sprintln("erro ao obter number of devices: ", err))
	}

	numberDevices, err := convertNumberOfDevices(numberOfDevices)
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao converter a string para inteiro:", err))
		return
	}

	slotNames, err := getMemorySlotNames(numberDevices)
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao obter a quantidade de slot's: ", err))
	}

	memorySizes, err := getMemorySizes(numberDevices)
	if err != nil {
		logToFile(fmt.Sprintln("erro ao obter o tamanho de memória: ", err))
	}

	// Remove colchetes e divide a string em partes
	cleanedInput := strings.Trim(fmt.Sprint(slotNames), "[]")
	partsNames := strings.Split(cleanedInput, ",")

	// Remove colchetes e divide a string em partes
	cleanedInput2 := strings.Trim(fmt.Sprint(memorySizes), "[]")
	partsSizes := strings.Split(cleanedInput2, ",")

	memoriesTypes, err := getMemoryTypes(numberDevices)
	if err != nil {
		logToFile(fmt.Sprintln("erro ao obter o tipo de memória: ", err))
	}

	// Remove colchetes e divide a string em partes
	cleanedInput3 := strings.Trim(fmt.Sprint(memoriesTypes), "[]")
	partsTypes := strings.Split(cleanedInput3, ",")

	memoriesTypeDetails, err := getMemoryTypeDetails(numberDevices)
	if err != nil {
		logToFile(fmt.Sprintln("erro ao obter tipo detalhado de memória: ", err))
	}

	// Remove colchetes e divide a string em partes
	cleanedInput4 := strings.Trim(fmt.Sprint(memoriesTypeDetails), "[]")
	partsTypeDetails := strings.Split(cleanedInput4, ",")

	memoriesSpeedMemory, err := getMemorySpeeds(numberDevices)
	if err != nil {
		logToFile(fmt.Sprintln("erro ao obter a velocidade de memória: ", err))
	}

	// Remove colchetes e divide a string em partes
	cleanedInput5 := strings.Trim(fmt.Sprint(memoriesSpeedMemory), "[]")
	partsSpeed := strings.Split(cleanedInput5, ",")

	memoriesSerialNumber, err := getMemorySerialNumbers(numberDevices)
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

	hard_disk_model, err := getHDModel("/dev/sda")
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao obter o modelo do HD:", err))
	}

	hard_disk_serial_number, err := getHDSerialModel("/dev/sda")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	} 

	hard_disk_user_capacity, err := getHDUserCapacity("/dev/sda")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	hard_disk_sata_version, err := getHDSataVersion("/dev/sda")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	cpu_architecture, err := getCPUInfo("Architecture:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	cpu_operation_mode, err := getCPUInfo("CPU op-mode(s):")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	} 

	cpus, err := getCPUInfo("CPU(s):")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	cpu_vendor_id, err := getCPUInfo("Vendor ID:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	cpu_model_name, err := getCPUInfo("Model name:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	cpu_thread, err := getCPUInfo("Thread(s) per core:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	cpu_core, err := getCPUInfo("Core(s) per socket:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	cpu_socket, err := getCPUInfo("Socket(s):")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	cpu_max_mhz, err := getCPUInfo("CPU max MHz:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	cpu_min_mhz, err := getCPUInfo("CPU min MHz:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	gpu_product, err := getGPUProduct("product:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	gpu_vendor_id, err := getGPUProduct("vendor:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	gpu_bus_info, err := getGPUProduct("bus info:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	gpu_logical_name, err := getGPUProduct("logical name:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	gpu_clock, err := getGPUProduct("clock:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	gpu_configuration, err := getGPUProduct("configuration:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}  

	audio_device_product_first_value, err := getAudioDevices("card 0:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	} 

	audio_device_product, err := extractValueFromBrackets(audio_device_product_first_value)
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	audio_device_model_first_value, err := getAudioDeviceModel("device 0:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	audio_device_model, err := extractValueFromBrackets(audio_device_model_first_value)
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}
	
	bios_version, err := getSMBIOSVersion()
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	} 

	motherboard_manufacturer, err := getSMBIOSInfo("Manufacturer:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	motherboard_product_name, err := getSMBIOSInfo("Product Name:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	motherboard_version, err := getSMBIOSInfo("Version:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	motherboard_serial_name, err := getSMBIOSInfo("Serial Number:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	motherboard_asset_tag, err := getSMBIOSInfo("Asset Tag:")
	if err != nil {
		logToFile(fmt.Sprintln("Erro:", err))
	}

	// Chama a função para listar os pacotes instalados
	installedPackages, err := listInstalledPackages()
	if err != nil {
		logToFile(fmt.Sprintln("Erro ao listar os pacotes instalados: ", err))
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
		Memories: memoriesList,
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

