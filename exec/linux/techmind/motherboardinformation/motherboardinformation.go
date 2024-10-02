package motherboardinformation

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func GetManufacturer()(string,error){
	cmd := exec.Command("sudo", "dmidecode", "-s", "system-manufacturer")
    output, err := cmd.CombinedOutput() // Captura a saída e erro juntos
    
    if err != nil {
        // Verifica se o erro é relacionado a um comando não encontrado
        if strings.Contains(string(output), "command not found") || strings.Contains(err.Error(), "executable file not found in $PATH") {
            cmd := exec.Command("sudo", "pkg", "install", "-y", "dmidecode")

			// Captura a saída do comando
			output, err := cmd.CombinedOutput()
			if err != nil {
				return string(output),  err
			}else{
				return GetManufacturer()
			}
        }
        return "", fmt.Errorf("erro ao executar dmidecode: %v", err)
    }
    
    manufacturer := strings.TrimSpace(string(output))
    
    if manufacturer == "" {
        return "", fmt.Errorf("marca do sistema não encontrada")
    }
    
    return manufacturer, nil
}

func GetModel()(string, error){
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

func GetSerialNumber()(string,error){
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

func GetMaxMem()(string,error){
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

func GetSlotDim()(string, error){
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

func ConvertNumberOfDevices(numberOfDevices string) (int, error) {
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

func GetMemorySlotNames(numberOfDevices int) ([]string, error) {
    // Executa o comando dmidecode
    cmd := exec.Command( "dmidecode", "-t", "memory")
    var out bytes.Buffer
	var stderr bytes.Buffer
    cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run() // Executa o comando
	
	if err != nil {
		if strings.Contains(stderr.String(), "command not found") || strings.Contains(err.Error(), "executable file not found in $PATH") {
			return nil, nil // Retorna nil e sem erro
		}
		return nil, fmt.Errorf("erro ao executar dmidecode: %v, stderr: %s", err, stderr.String())
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

func GetMemorySizes(numberOfDevices int) ([]string, error) {
	cmd := exec.Command( "dmidecode", "-t", "memory")
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

func GetMemoryTypes(numberOfDevices int) ([]string, error) {
	cmd := exec.Command( "dmidecode", "-t", "memory")
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

func GetMemoryTypeDetails(numberOfDevices int) ([]string, error) {
	cmd := exec.Command("dmidecode", "-t", "memory")
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

func GetMemorySpeeds(numberOfDevices int) ([]string, error) {
	cmd := exec.Command("dmidecode", "-t", "memory")
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

func GetMemorySerialNumbers(numberOfDevices int) ([]string, error) {
	cmd := exec.Command("dmidecode", "-t", "memory")
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

