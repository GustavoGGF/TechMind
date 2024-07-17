package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"

	// "github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/host"
)

var (
    modntdll = syscall.NewLazyDLL("ntdll.dll")
    procRtlGetVersion = modntdll.NewProc("RtlGetVersion")
)

var (
    modkernel32        = syscall.NewLazyDLL("kernel32.dll")
    procGetComputerName = modkernel32.NewProc("GetComputerNameW")
)

type RTL_OSVERSIONINFOEX struct {
    dwOSVersionInfoSize uint32
    dwMajorVersion      uint32
    dwMinorVersion      uint32
    dwBuildNumber       uint32
    dwPlatformId        uint32
    szCSDVersion        [128]uint16
}

type Data struct {
    System, Name, Distribution, InsertionDate, MacAddress, CurrentUser,PlatformVersion string
}

func getWindowsVersion() (string, error) {
    var info RTL_OSVERSIONINFOEX
    info.dwOSVersionInfoSize = uint32(unsafe.Sizeof(info))

    ret, _, _ := procRtlGetVersion.Call(uintptr(unsafe.Pointer(&info)))
    if ret != 0 {
        return "", syscall.GetLastError()
    }

    major := info.dwMajorVersion
    minor := info.dwMinorVersion

    switch {
    case major == 10 && minor == 0:
        return "Windows 10", nil
    case major == 6 && minor == 3:
        return "Windows 8.1", nil
    case major == 6 && minor == 2:
        return "Windows 8", nil
    case major == 6 && minor == 1:
        return "Windows 7", nil
    default:
        return fmt.Sprintf("Windows (versão %d.%d)", major, minor), nil
    }
}

func getComputerName() (string, error) {
    var nSize uint32 = 256
    nameBuf := make([]uint16, nSize)
    ret, _, err := procGetComputerName.Call(uintptr(unsafe.Pointer(&nameBuf[0])), uintptr(unsafe.Pointer(&nSize)))
    if ret == 0 {
        return "", err
    }
    return string(utf16.Decode(nameBuf[:nSize-1])), nil
}

func logToFile(msg string) {
	// Abrir arquivo de log
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo de log:", err)
		return
	}
	defer file.Close()

	// Criar um logger para o arquivo
	logger := log.New(file, "", log.LstdFlags)

	// Escrever mensagem no arquivo de log
	logger.Println(msg)
}

func getMac(dis string) (string, error) {
    if dis == "Windows 10" {
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

    return "", fmt.Errorf("sistema operacional não suportado ou não especificado")
}

func main() {
    distribution := ""
    
    if ver, err := getWindowsVersion(); err == nil {
        
        distribution = ver

    } else {
        logToFile(fmt.Sprintf("Erro ao obter a versão do Windows: %v", err))
        return
    }

    sys := runtime.GOOS

    hostname:= ""

    if computerName, err := getComputerName(); err == nil {
        hostname = computerName
    } else {
        logToFile(fmt.Sprintf("Erro ao obter o nome da máquina: %v", err))
        return
    }

    sys = strings.ReplaceAll(sys, " ", "")
    hostname = strings.ReplaceAll(hostname, " ", "")
    distribution = strings.ReplaceAll(distribution, " ", "")
    
    url := "http://10.1.1.73:3000/home/computers/post-machines" //env

    date_now := time.Now()

    formated_date := date_now.Format("2006-01-02 15:04")

    macAddress, err := getMac("Windows 10")
    
    if err != nil {
        logToFile(fmt.Sprintf("Erro ao obter o endereço MAC: %v", err))
        return
    }

    currentUser, err := user.Current()
	if err != nil {
		logToFile(fmt.Sprintf("Erro ao obter o usuário atual: %v", err))
	}

    nameUser := string(currentUser.Username)

    info, err := host.Info()
    if err != nil {
        logToFile(fmt.Sprintf("Error:", err))
        return
    }

    version := info.PlatformVersion

    jsonData := Data{System: sys, Name: hostname, Distribution: distribution, InsertionDate: formated_date, MacAddress:macAddress, CurrentUser:nameUser, PlatformVersion:version}

    requestBody, err := json.Marshal(jsonData)
    if err != nil{
        logToFile(fmt.Sprintf("Erro ao montar o json: %v", err))
        return
    }

    resp, erro := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
    if erro != nil{
        logToFile(fmt.Sprintf("Erro ao fazer o post: %v", err))
    }

    // Erro response 400 gerar aviso na tela

    defer resp.Body.Close()  

}
