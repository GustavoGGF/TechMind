package memory

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

// Win32_PhysicalMemoryArray representa informações sobre o array de memória física
// disponível no sistema, incluindo a capacidade máxima do array.
type Win32_PhysicalMemoryArray struct {
	MaxCapacity uint32 // Capacidade máxima do array de memória física em megabytes
}

type Win32_PhysicalMemoryArray2 struct {
	MemoryDevices uint32
}


type Win32_PhysicalMemory struct {
	BankLabel     string
	Capacity      uint64
	DeviceLocator string
	MemoryType    uint16
	TypeDetail    uint16
	Speed         uint32
	SerialNumber  string
}

// GetMaxMemoryCapacity obtém a capacidade máxima de memória suportada pelo sistema
// consultando a classe Win32_PhysicalMemoryArray via WMI.
func GetMaxMemoryCapacity() (string, error) {
	// Declara uma variável para armazenar os resultados da consulta WMI
	var arrays []Win32_PhysicalMemoryArray

	// Executa a consulta WMI para obter a capacidade máxima do array de memória física
	err := wmi.Query("SELECT MaxCapacity FROM Win32_PhysicalMemoryArray", &arrays)
	if err != nil {
		// Retorna um erro formatado se a consulta falhar
		return "", fmt.Errorf("erro ao consultar WMI (PhysicalMemoryArray): %v", err)
	}

	// Declara uma variável para armazenar a capacidade máxima
	var maxCapacity uint32
	// Verifica se a consulta retornou algum resultado
	if len(arrays) > 0 {
		// Atribui a capacidade máxima do primeiro (e supostamente único) array
		maxCapacity = arrays[0].MaxCapacity
	}

	// Converte a capacidade máxima de megabytes para gigabytes
	maxCapacityGB := float64(maxCapacity) / (1024 * 1024)
	// Formata a capacidade máxima como uma string em gigabytes, arredondando para o inteiro mais próximo
	maxCapacityStr := fmt.Sprintf("%.0f GB", maxCapacityGB)

	// Retorna a capacidade máxima como string
	return maxCapacityStr, nil
}

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

func GetMemoryDetails() ([]Win32_PhysicalMemory, error) {
	var memories []Win32_PhysicalMemory

	// Consulta WMI
	err := wmi.Query("SELECT BankLabel, Capacity, DeviceLocator, MemoryType, TypeDetail, Speed, SerialNumber FROM Win32_PhysicalMemory", &memories)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar WMI (PhysicalMemory): %v", err)
	}

	return memories, nil
}

func GetMemoryType(memoryType uint16) string {
	// Mapeamento dos códigos de tipo de memória para valores descritivos
	switch memoryType {
	case 0:
		return "Unknown"
	case 1:
		return "Other"
	case 2:
		return "DRAM"
	case 3:
		return "Synchronous DRAM"
	case 4:
		return "Cache DRAM"
	case 5:
		return "EDO"
	case 6:
		return "EDRAM"
	case 7:
		return "VRAM"
	case 8:
		return "SRAM"
	case 9:
		return "RAM"
	case 10:
		return "ROM"
	case 11:
		return "Flash"
	case 12:
		return "EEPROM"
	case 13:
		return "FEPROM"
	case 14:
		return "EPROM"
	case 15:
		return "CDRAM"
	case 16:
		return "3DRAM"
	case 17:
		return "SDRAM"
	case 18:
		return "SGRAM"
	default:
		return "Unknown"
	}
}