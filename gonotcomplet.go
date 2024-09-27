package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)
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
func getsd() { 




	hard_disk_user_capacity, err := getHDUserCapacity("/dev/sda")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	hard_disk_sata_version, err := getHDSataVersion("/dev/sda")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	cpu_architecture, err := getCPUInfo("Architecture:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	cpu_operation_mode, err := getCPUInfo("CPU op-mode(s):")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	} 

	cpus, err := getCPUInfo("CPU(s):")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	cpu_vendor_id, err := getCPUInfo("Vendor ID:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	cpu_model_name, err := getCPUInfo("Model name:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	cpu_thread, err := getCPUInfo("Thread(s) per core:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	cpu_core, err := getCPUInfo("Core(s) per socket:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	cpu_socket, err := getCPUInfo("Socket(s):")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	cpu_max_mhz, err := getCPUInfo("CPU max MHz:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	cpu_min_mhz, err := getCPUInfo("CPU min MHz:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	gpu_product, err := getGPUProduct("product:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	gpu_vendor_id, err := getGPUProduct("vendor:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	gpu_bus_info, err := getGPUProduct("bus info:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	gpu_logical_name, err := getGPUProduct("logical name:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	gpu_clock, err := getGPUProduct("clock:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	gpu_configuration, err := getGPUProduct("configuration:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}  

	audio_device_product_first_value, err := getAudioDevices("card 0:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	} 

	audio_device_product, err := extractValueFromBrackets(audio_device_product_first_value)
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	audio_device_model_first_value, err := getAudioDeviceModel("device 0:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	audio_device_model, err := extractValueFromBrackets(audio_device_model_first_value)
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}
	
	bios_version, err := getSMBIOSVersion()
	if err != nil {
		fmt.Sprintln("Erro:", err)
	} 

	motherboard_manufacturer, err := getSMBIOSInfo("Manufacturer:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	motherboard_product_name, err := getSMBIOSInfo("Product Name:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	motherboard_version, err := getSMBIOSInfo("Version:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	motherboard_serial_name, err := getSMBIOSInfo("Serial Number:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	motherboard_asset_tag, err := getSMBIOSInfo("Asset Tag:")
	if err != nil {
		fmt.Sprintln("Erro:", err)
	}

	// Chama a função para listar os pacotes instalados
	installedPackages, err := listInstalledPackages()
	if err != nil {
		fmt.Sprintln("Erro ao listar os pacotes instalados: ", err)
	}

	// Convertendo a string para um slice de strings
    packageSlice := strings.Split(installedPackages, ",")

		// CPUArchitecture: cpu_architecture, 
		// CPUOperationMode:cpu_operation_mode, 
		// CPUS: cpus,
		// CPUVendorID: cpu_vendor_id, 
		// CPUModelName: cpu_model_name, 
		// CPUThread: cpu_thread,
		// CPUCore: cpu_core,
		// CPUSocket:cpu_socket, 
		// CPUMaxMHz: cpu_max_mhz, 
		// CPUMinMHz: cpu_min_mhz,
		// GPUProduct: gpu_product, 
		// GPUVendorID: gpu_vendor_id, 
		// GPUBusInfo: gpu_bus_info,
		// GPULogicalName: gpu_logical_name, 
		// GPUClock: gpu_clock, 
		// GPUConfiguration: gpu_configuration,
		// AudioDeviceProduct: audio_device_product, 
		// AudioDeviceModel: audio_device_model,
		// BiosVersion: bios_version,
		// MotherboardManufacturer:motherboard_manufacturer,
		// MotherboardProductName:motherboard_product_name, 
		// MotherboardVersion:motherboard_version,
		// MotherbaoardSerialName: motherboard_serial_name,
		// MotherboardAssetTag:motherboard_asset_tag,
		// SoftwareNames:packageSlice,

    requestBody, err := json.Marshal(jsonData)
    if err != nil{
        fmt.Println(err)
    }

    resp, erro := http.Post("url", "application/json", bytes.NewBuffer(requestBody))
    if erro != nil{
        fmt.Println(erro)
    }

    defer resp.Body.Close()  
}

