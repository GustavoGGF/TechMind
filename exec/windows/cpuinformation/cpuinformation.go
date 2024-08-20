package cpuinformation

import (
	"fmt"

	"github.com/yusufpapurcu/wmi"
)

// Struc para armazenar informações sobre o CPU
type Win32_Processor struct {
	Architecture              uint16
	Manufacturer              string
	Name                      string
	NumberOfLogicalProcessors uint32
	NumberOfCores             uint32
	SocketDesignation         string
	MaxClockSpeed             uint32
	CurrentClockSpeed         uint32
}

// Struc para armazenar a arquitetura do SO
type Win32_OperatingSystem struct {
	OSArchitecture string
}

func GetCPUArchitecture() (string, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no processors found")
	}

	architecture := dst[0].Architecture
	var arch string
	switch architecture {
	case 0:
		arch = "x86"
	case 1:
		arch = "MIPS"
	case 2:
		arch = "Alpha"
	case 3:
		arch = "PowerPC"
	case 5:
		arch = "ARM"
	case 6:
		arch = "ia64"
	case 9:
		arch = "x64"
	default:
		arch = "Unknown"
	}

	return arch, nil
}

func GetCPUOperationMode() (string, error) {
	var dst []Win32_OperatingSystem
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no operating systems found")
	}

	return dst[0].OSArchitecture, nil
}

func GetCPUCount() (uint32, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, err
	}

	if len(dst) == 0 {
		return 0, fmt.Errorf("no processors found")
	}

	// Somando o número de núcleos de todos os processadores
	var totalCores uint32
	for _, processor := range dst {
		totalCores += processor.NumberOfCores
	}

	return totalCores, nil
}

func GetCPUVendorID() (string, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no processors found")
	}

	// Escolhe o primeiro processador
	return dst[0].Manufacturer, nil
}

func GetCPUModelName() (string, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no processors found")
	}

	// Escolhe o primeiro processador
	return dst[0].Name, nil
}

func GetCPUThreads() (uint32, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, err
	}

	if len(dst) == 0 {
		return 0, fmt.Errorf("no processors found")
	}

	// Escolhe o primeiro processador
	return dst[0].NumberOfLogicalProcessors, nil
}

func GetCPUCores() (uint32, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, err
	}

	if len(dst) == 0 {
		return 0, fmt.Errorf("no processors found")
	}

	// Escolhe o primeiro processador
	return dst[0].NumberOfCores, nil
}

func GetCPUSockets() (int, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, err
	}

	if len(dst) == 0 {
		return 0, fmt.Errorf("no processors found")
	}

	return len(dst), nil
}

func GetCPUMaxMHz() (uint32, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, err
	}

	if len(dst) == 0 {
		return 0, fmt.Errorf("no processors found")
	}

	// Escolhe o primeiro processador
	return dst[0].MaxClockSpeed, nil
}

func GetCPUMinMHz() (uint32, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, err
	}

	if len(dst) == 0 {
		return 0, fmt.Errorf("no processors found")
	}

	// Escolhe o primeiro processador
	return dst[0].CurrentClockSpeed, nil
}
