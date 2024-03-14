package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
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
    System, Name, Distribution string
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

	// Redirecionar a saída padrão para o arquivo
	os.Stdout = file

	// Imprimir mensagem de log
	fmt.Println(msg)
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
    
    url := "http://10.1.9.0:3000/api/machines" //env

    jsonData := Data{System: sys, Name: hostname, Distribution: distribution}

    requestBody, err := json.Marshal(jsonData)
    if err != nil{
        logToFile(fmt.Sprintf("Erro ao montar o json: %v", err))
        return
    }

    resp, erro := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
    if erro != nil{
        logToFile(fmt.Sprintf("Erro ao fazer o post: %v", err))
    }

    defer resp.Body.Close()  
}
