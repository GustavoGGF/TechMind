package cpuinformation

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func GetCPUInfo() (string,error) {
	// Executa o comando sysctl para obter o Vendor ID
	cmd := exec.Command("sysctl", "-n", "hw.model")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "",err
	}

	// Exibe o Vendor ID
	return  out.String(), nil
}

func GetThread()(int,error){
		// Executa o comando sysctl para obter o número de threads (CPUs lógicos)
		cmd := exec.Command("sysctl", "-n", "hw.ncpu")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return 0, fmt.Errorf("erro ao obter o número de threads do CPU: %v", err)
		}
	
		// Converte a saída para inteiro
		var numThreads int
		_, err = fmt.Sscanf(out.String(), "%d", &numThreads)
		if err != nil {
			return 0, fmt.Errorf("erro ao converter a saída para número: %v", err)
		}
	
		return numThreads, nil
}

func GetMaxMHz()(int, error) {
	// Executa o comando sysctl para obter os níveis de frequência da CPU
	cmd := exec.Command("sysctl", "-n", "dev.cpu.0.freq_levels")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter os níveis de frequência da CPU: %v", err)
	}

	// A primeira parte da saída contém a frequência máxima em MHz
	// Exemplo de saída: "4000/1 3000/1 2000/1 ..."
	freqs := strings.Fields(out.String())
	if len(freqs) == 0 {
		return 0, fmt.Errorf("nenhuma frequência encontrada")
	}

	// Extrai o valor da primeira frequência (a máxima)
	maxFreqStr := strings.Split(freqs[0], "/")[0]
	maxFreq, err := strconv.Atoi(maxFreqStr)
	if err != nil {
		return 0, fmt.Errorf("erro ao converter frequência máxima para número: %v", err)
	}

	return maxFreq, nil
}

func GetCPUMinMHz() (int, error) {
	// Executa o comando sysctl para obter os níveis de frequência da CPU
	cmd := exec.Command("sysctl", "-n", "dev.cpu.0.freq_levels")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter os níveis de frequência da CPU: %v", err)
	}

	// A última parte da saída contém a frequência mínima em MHz
	// Exemplo de saída: "4000/1 3000/1 2000/1 ..."
	freqs := strings.Fields(out.String())
	if len(freqs) == 0 {
		return 0, fmt.Errorf("nenhuma frequência encontrada")
	}

	// Extrai o valor da última frequência (a mínima)
	minFreqStr := strings.Split(freqs[len(freqs)-1], "/")[0]
	minFreq, err := strconv.Atoi(minFreqStr)
	if err != nil {
		return 0, fmt.Errorf("erro ao converter frequência mínima para número: %v", err)
	}

	return minFreq, nil
}