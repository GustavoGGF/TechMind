package network

import (
	"errors"
	"fmt"
	"net"
)

// GetMac obtém o endereço MAC da primeira interface de rede disponível
func GetMac() (string, error) {
	// Obtém a lista de todas as interfaces de rede disponíveis no sistema
	interfaces, err := net.Interfaces()
	if err != nil {
		// Retorna o erro caso ocorra algum problema ao obter as interfaces
		return "", err
	}

	// Itera sobre as interfaces de rede
	for _, iface := range interfaces {
		// Obtém o endereço MAC da interface
		mac := iface.HardwareAddr.String()
		// Verifica se a interface possui um endereço MAC
		if len(mac) > 0 {
			// Retorna o endereço MAC se encontrado
			return mac, nil
		}
	}
	// Retorna uma string vazia e sem erro caso nenhum endereço MAC seja encontrado
	return "", nil
}

// GetIP obtém o primeiro endereço IP não loopback do sistema.
func GetIP() (string, error) {
	// Obtém todos os endereços de rede associados às interfaces do sistema
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		// Retorna um erro formatado caso ocorra um problema ao obter os endereços de rede
		return "", fmt.Errorf("erro ao obter endereços de rede: %v", err)
	}

	// Itera sobre os endereços de rede encontrados
	for _, addr := range addrs {
		// Verifica se o endereço é do tipo *net.IPNet e se não é um endereço de loopback
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// Verifica se o endereço IP é do tipo IPv4
			if ipNet.IP.To4() != nil {
				// Retorna o endereço IP encontrado como uma string
				return ipNet.IP.String(), nil
			}
		}
	}
	// Retorna um erro se nenhum endereço IP válido for encontrado
	return "", errors.New("nenhum endereço IP encontrado")
}
