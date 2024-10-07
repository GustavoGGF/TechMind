package internetinformation

import (
	"fmt"
	"log"
	"net"
)

func bytesEqual(b []byte) bool {
    for _, v := range b {
        if v != 0 {
            return false
        }
    }
    return true
}

func GetMac()(string, error){
	interfaces, err := net.Interfaces()
    if err != nil {
        return "", err
    }
    imac := ""

    for _, iface := range interfaces {
        if iface.Flags&net.FlagUp != 0 && !bytesEqual(iface.HardwareAddr) {
            imac = iface.HardwareAddr.String()
            break
        }
    }

	return imac, nil
}

func GetIP() (string, error){
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Printf("%s: %v", iface.Name, err)
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ipNet.IP.IsGlobalUnicast() {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("nenhum endereço IP global encontrado")
}