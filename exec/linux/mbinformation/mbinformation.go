package mbinformation

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

func GetMotherboardManufacturer()(string,error){
		// Executa o comando `dmidecode -t baseboard`
		cmd := exec.Command("dmidecode", "-t", "baseboard")
		var out bytes.Buffer
		cmd.Stdout = &out
	
		err := cmd.Run()
		if err != nil {
			return "", err
		}
	
		// Lê a saída linha por linha
		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			line := scanner.Text()
			// Procura pela linha que contém "Manufacturer:"
			if strings.Contains(line, "Manufacturer:") {
				// Extrai o valor após "Manufacturer:"
				manufacturer := strings.TrimSpace(strings.Split(line, ":")[1])
				return manufacturer, nil
			}
		}
	
		// Verifica se houve erros ao ler a saída
		if err := scanner.Err(); err != nil {
			return "", err
		}

		return "", nil
}

func GetMotherboardPD()(string,error){
	// Executa o comando `dmidecode -t baseboard`
	cmd := exec.Command("dmidecode", "-t", "baseboard")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Lê a saída linha por linha
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		// Procura pela linha que contém "Product Name:"
		if strings.Contains(line, "Product Name:") {
			// Extrai o valor após "Manufacturer:"
			manufacturer := strings.TrimSpace(strings.Split(line, ":")[1])
			return manufacturer, nil
		}
	}

	// Verifica se houve erros ao ler a saída
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func GetMotherboardVersion()(string,error){
	// Executa o comando `dmidecode -t baseboard`
	cmd := exec.Command("dmidecode", "-t", "baseboard")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Lê a saída linha por linha
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Version:") {
			// Extrai o valor após "Manufacturer:"
			manufacturer := strings.TrimSpace(strings.Split(line, ":")[1])
			return manufacturer, nil
		}
	}

	// Verifica se houve erros ao ler a saída
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func GetMotherSN()(string,error){
	// Executa o comando `dmidecode -t baseboard`
	cmd := exec.Command("dmidecode", "-t", "baseboard")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Lê a saída linha por linha
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Serial Number:") {
			// Extrai o valor após "Manufacturer:"
			manufacturer := strings.TrimSpace(strings.Split(line, ":")[1])
			return manufacturer, nil
		}
	}

	// Verifica se houve erros ao ler a saída
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func GetMotherboardAssetTag()(string,error){
	// Executa o comando `dmidecode -t baseboard`
	cmd := exec.Command("dmidecode", "-t", "baseboard")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Lê a saída linha por linha
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Asset Tag:") {
			// Extrai o valor após "Manufacturer:"
			manufacturer := strings.TrimSpace(strings.Split(line, ":")[1])
			return manufacturer, nil
		}
	}

	// Verifica se houve erros ao ler a saída
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}