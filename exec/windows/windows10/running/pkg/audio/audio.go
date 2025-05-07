package audio

import (
	"fmt"

	"github.com/yusufpapurcu/wmi"
)

type Win32_SoundDevice struct {
	ProductName string
}

func GetAudioDeviceProduct() (string, error) {
	var dst []Win32_SoundDevice
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no audio devices found")
	}

	// Escolhe o primeiro dispositivo de áudio encontrado
	return dst[0].ProductName, nil
}