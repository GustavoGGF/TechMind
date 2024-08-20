package harddiskinformation

import (
	"fmt"
	"strconv"

	"github.com/yusufpapurcu/wmi"
)

// Estrutura que representa os dados que queremos obter do WMI
type Win32_DiskDrive struct {
	Model string
}

// Estrutura para mapear a consulta WMI
type Win32_DiskDrive2 struct {
	SerialNumber string
}

// Estrutura para mapear a consulta WMI
type Win32_DiskDrive3 struct {
	Size string // O campo Size é retornado como string
}

// Função para obter o modelo do disco rígido
func GetHardDiskModel() ([]string, error) {
	var dst []Win32_DiskDrive
	query := "SELECT Model FROM Win32_DiskDrive"
	err := wmi.Query(query, &dst)
	if err != nil {
		return nil, fmt.Errorf("falha ao executar a consulta WMI: %v", err)
	}

	var models []string
	for _, disk := range dst {
		models = append(models, disk.Model)
	}

	return models, nil
}

// Função para obter o número de série do disco rígido
func GetHardDiskSerialNumber() ([]string, error) {
	var dst []Win32_DiskDrive2
	query := "SELECT SerialNumber FROM Win32_DiskDrive"
	err := wmi.Query(query, &dst)
	if err != nil {
		return nil, fmt.Errorf("falha ao executar a consulta WMI: %v", err)
	}

	var serialNumbers []string
	for _, disk := range dst {
		serialNumbers = append(serialNumbers, disk.SerialNumber)
	}

	return serialNumbers, nil
}

// Função para obter a capacidade do disco rígido em GB
func GetHardDiskCapacity() ([]float64, error) {
	var dst []Win32_DiskDrive3
	query := "SELECT Size FROM Win32_DiskDrive"
	err := wmi.Query(query, &dst)
	if err != nil {
		return nil, fmt.Errorf("falha ao executar a consulta WMI: %v", err)
	}

	var capacities []float64
	for _, disk := range dst {
		size, err := strconv.ParseUint(disk.Size, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("falha ao converter o tamanho do disco: %v", err)
		}
		capacities = append(capacities, float64(size)/(1024*1024*1024)) // Convertendo bytes para gigabytes
	}

	return capacities, nil
}
