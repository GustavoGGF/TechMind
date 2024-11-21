package gpuinformation

import (
	"fmt"

	"github.com/yusufpapurcu/wmi"
)

type Win32_VideoController struct {
	Name                        string
	AdapterCompatibility        string
	DeviceID                    string
	DriverVersion               string
	CurrentHorizontalResolution uint32
	CurrentVerticalResolution   uint32
	AdapterRAM                  uint32
}

func GetGPUProduct() (string, error) {
	var dst []Win32_VideoController
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no video controllers found")
	}

	// Escolhe o primeiro controlador de vídeo
	return dst[0].Name, nil
}

func GetGPUVendorID() (string, error) {
	var dst []Win32_VideoController
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no video controllers found")
	}

	// Escolhe o primeiro controlador de vídeo
	return dst[0].AdapterCompatibility, nil
}

func GetGPUBusInfo() (string, error) {
	var dst []Win32_VideoController
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no video controllers found")
	}

	// Escolhe o primeiro controlador de vídeo
	return dst[0].DeviceID, nil
}

func GetGPULogicalName() (string, error) {
	var dst []Win32_VideoController
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no video controllers found")
	}

	// Itera sobre todos os controladores de vídeo para obter nomes lógicos
	for _, controller := range dst {
		fmt.Printf("GPU Logical Name: %s\n", controller.Name)
	}

	// Retorna o nome da primeira GPU encontrada como exemplo
	return dst[0].Name, nil
}

func GetGPUClock() (string, error) {
	var dst []Win32_VideoController
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no video controllers found")
	}

	// Exibe a versão do driver como uma aproximação para a frequência do clock
	return dst[0].DriverVersion, nil
}

func GetGPUConfiguration() (uint32, uint32, uint32, error) {
	var dst []Win32_VideoController
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, 0, 0, err
	}

	if len(dst) == 0 {
		return 0, 0, 0, fmt.Errorf("no video controllers found")
	}

	// Escolhe o primeiro controlador de vídeo
	controller := dst[0]
	return controller.CurrentHorizontalResolution, controller.CurrentVerticalResolution, controller.AdapterRAM, nil
}
