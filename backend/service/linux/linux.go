package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/shirou/gopsutil/host"
)

type Data struct {
    System, Name, Distribution,InterfaceInternet, MacAddress string
     }

func bytesEqual(b []byte) bool {
    for _, v := range b {
        if v != 0 {
            return false
        }
    }
    return true
}

func main() {
    sys := runtime.GOOS

    infos, err:= host.Info()
    if err != nil{
        fmt.Println("Error")
    }

    url := "http://10.1.9.0:3000/api/machines" //env

    name, err := os.Hostname()

    if err != nil {
        fmt.Println(err)
    }
    interfaces, err := net.Interfaces()
    if err != nil {
        fmt.Println("Erro ao obter interfaces de rede:", err)
        return
    }

    ifaceInt:= ""
    imac := ""

    for _, iface := range interfaces {
        if iface.Flags&net.FlagUp != 0 && !bytesEqual(iface.HardwareAddr) {
            ifaceInt = iface.Name
            imac = iface.HardwareAddr.String()
            break
        }
    }

    jsonData := Data{System: sys, Name: name, Distribution: infos.Platform, InterfaceInternet: ifaceInt, MacAddress: imac}

    requestBody, err := json.Marshal(jsonData)
    if err != nil{
        fmt.Println(err)
    }

    resp, erro := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
    if erro != nil{
        fmt.Println(erro)
    }

    defer resp.Body.Close()  
}

