package hdinformation

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func GetHDModel() (string, error) {
    // Comando que será executado
    cmd := exec.Command("camcontrol", "devlist")

    // Executa o comando e captura a saída
    output, err := cmd.Output()
    if err != nil {
        return "",  fmt.Errorf("erro ao obter a lista de HD: %w", err)
    }
	return string(output), nil
}

func ExtractValues(input string) []string {
    // Expressão regular para capturar valores entre <>
    re := regexp.MustCompile(`<([^>]+)>`)
    // FindAllStringSubmatch retorna todos os grupos correspondentes
    matches := re.FindAllStringSubmatch(input, -1)

    var result []string
    for _, match := range matches {
        if len(match) > 1 {
            result = append(result, match[1])
        }
    }
    return result
}

func DevicesListADA() ([]string, error) {
	// Executa o comando para listar dispositivos
	cmd := exec.Command("camcontrol", "devlist")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// Usa uma expressão regular para extrair os dispositivos ada (ex: ada0, ada1)
	re := regexp.MustCompile(`(ada\d+)`)
	matches := re.FindAllString(out.String(), -1)

	// Remove duplicatas (no caso de múltiplas entradas para o mesmo dispositivo)
	uniqueDisks := make(map[string]bool)
	for _, match := range matches {
		uniqueDisks[match] = true
	}

	// Converte o mapa para um slice
	var discos []string
	for disco := range uniqueDisks {
		discos = append(discos, disco)
	}

	return discos, nil
}

func GetSerialNumber(disco string) (string, error) {
	// Executa o comando camcontrol identify para o dispositivo específico
	cmd := exec.Command("camcontrol", "identify", disco)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Analisa a saída para encontrar o número de série
	return extractSerialNumber(out.String()), nil
}

// Função para extrair o número de série da saída do camcontrol
func extractSerialNumber(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// Procura pela linha que contém "serial number"
		if strings.Contains(line, "serial number") {
			parts := strings.Split(line, "serial number")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

func GetLBA48(disco string) (int, error){
		// Executa o comando camcontrol identify para o dispositivo específico
		cmd := exec.Command("camcontrol", "identify", disco)
		var out bytes.Buffer
		cmd.Stdout = &out
	
		err := cmd.Run()
		if err != nil {
			return 0, err
		}
	
		// Analisa a saída para encontrar o número de série
		return extractSize(out.String()), nil
}

// Função para extrair o número de série da saída do camcontrol
func extractSize(output string) int {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// Procura pela linha que contém "serial number"
		if strings.Contains(line, "LBA48 supported") {
			parts := strings.Split(line, "LBA48 supported")
			if len(parts) > 1{
				nw_parts := strings.TrimSpace(parts[1])
				nw_bytes := strings.Split(nw_parts, "sectors")
				if len(nw_bytes) > 1{
					byts := strings.TrimSpace(nw_bytes[0])
					cnvt_to_int, err := strconv.Atoi(byts)
					if err != nil{
						return 0
					}
					return cnvt_to_int
				}	
			}
		}
	}
	return 0
}

func SataVersion(disco string) (string, error){
	// Executa o comando camcontrol identify para o dispositivo específico
	cmd := exec.Command("camcontrol", "identify", disco)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Analisa a saída para encontrar o número de série
	return extractSV(out.String()), nil
}

// Função para extrair o valor da segunda linha "pass0"
func extractSV(output string) string {
	lines := strings.Split(output, "\n")
	pass0Count := 0 // Contador para rastrear ocorrências de "pass0"

	for _, line := range lines {
		// Procura pela linha que contém "pass0"
		if strings.Contains(line, "pass") {
			pass0Count++ // Incrementa o contador
			if pass0Count == 2 { // Quando encontrar a segunda ocorrência
				parts := strings.Split(line, ":") // Divide a linha em torno de "pass0:"
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1]) // Remove espaços em branco
				}
			}
		}
	}
	return "" // Retorna uma string vazia se não encontrar a segunda ocorrência
}
