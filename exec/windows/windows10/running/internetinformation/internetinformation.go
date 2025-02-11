package internetinformation

import (
	"errors"
	"fmt"
	"net"
)

// Função para obter o macaddress
func GetMac() (string, error) {
	// Obtém todas as interfaces de rede do sistema
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("erro ao obter interfaces de rede: %v", err)
	}

	// Itera sobre as interfaces de rede para encontrar o endereço MAC
	for _, iface := range interfaces {
		// Ignora as interfaces loopback e desativadas
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		// Obtém o endereço MAC da interface
		mac := iface.HardwareAddr
		if mac != nil {
			return mac.String(), nil
		}
	}
	return "", fmt.Errorf("nenhum endereço MAC encontrado")
}

// Obtem o IP
func GetIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("erro ao obter endereços de rede: %v", err)
	}

	for _, addr := range addrs {
		// Verifica se o endereço é do tipo *net.IPNet e não é um endereço de loopback
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// Verifica se é um endereço IPv4
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("nenhum endereço IP encontrado")
}
