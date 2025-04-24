package memory

import (
	"fmt"

	"github.com/yusufpapurcu/wmi"
)

// Win32_PhysicalMemoryArray representa a classe WMI Win32_PhysicalMemoryArray
type Win32_PhysicalMemoryArray struct {
	MaxCapacity uint32
}

// Win32_PhysicalMemoryArray representa a classe WMI Win32_PhysicalMemoryArray
type Win32_PhysicalMemoryArray2 struct {
	MemoryDevices uint32
}

// Win32_PhysicalMemory representa a classe WMI Win32_PhysicalMemory
type Win32_PhysicalMemory struct {
	BankLabel     string
	Capacity      uint64
	DeviceLocator string
	MemoryType    uint16
	TypeDetail    uint16
	Speed         uint32
	SerialNumber  string
}

// Obtem a quantidade Máxima de memoria RAM suportada
func GetMaxMemoryCapacity() (string, error) {
	var arrays []Win32_PhysicalMemoryArray

	// Consulta para obter a capacidade máxima suportada
	err := wmi.Query("SELECT MaxCapacity FROM Win32_PhysicalMemoryArray", &arrays)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar WMI (PhysicalMemoryArray): %v", err)
	}

	// Assume que há apenas um PhysicalMemoryArray e pega a capacidade máxima suportada
	var maxCapacity uint32
	if len(arrays) > 0 {
		maxCapacity = arrays[0].MaxCapacity
	}

	maxCapacityGB := float64(maxCapacity) / (1024 * 1024)
	// Converte o número para string
	maxCapacityStr := fmt.Sprintf("%.0f GB", maxCapacityGB)

	return maxCapacityStr, nil
}

// Obtem informação sobre os slots
func GetMemorySlots() (string, error) {
	var arrays []Win32_PhysicalMemoryArray2

	// Consulta para obter o número de slots de memória
	err := wmi.Query("SELECT MemoryDevices FROM Win32_PhysicalMemoryArray", &arrays)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar WMI (PhysicalMemoryArray2): %v", err)
	}

	// Assume que há apenas um PhysicalMemoryArray e pega o número de slots de memória
	if len(arrays) > 0 {
		return fmt.Sprintf("%d", arrays[0].MemoryDevices), nil
	}

	return "", fmt.Errorf("nenhuma informação encontrada")
}

// Obtem informações detalhadas sobre a memoria
func GetMemoryDetails() ([]Win32_PhysicalMemory, error) {
	var memories []Win32_PhysicalMemory

	// Consulta para obter informações detalhadas da memória RAM
	err := wmi.Query("SELECT BankLabel, Capacity, DeviceLocator, MemoryType, TypeDetail, Speed, SerialNumber FROM Win32_PhysicalMemory", &memories)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar WMI (PhysicalMemory): %v", err)
	}

	return memories, nil
}
