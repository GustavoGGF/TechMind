package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/host"
)

type Data struct {
    System, Name, Distribution,InterfaceInternet, MacAddress,InsertionDate string
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

    url := "http://10.1.9.19:3000/api/machines" //env

    fmt.Println("Conectando")

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

    date_now := time.Now()

    formated_date := date_now.Format("2006-01-02 15:04")

    jsonData := Data{System: sys, Name: name, Distribution: infos.Platform, InterfaceInternet: ifaceInt, MacAddress: imac, InsertionDate:formated_date}

    fmt.Println("Montando")

    requestBody, err := json.Marshal(jsonData)
    if err != nil{
        fmt.Println(err)
    }

    resp, erro := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
    if erro != nil{
        fmt.Println(erro)
    }

    fmt.Println("Enviado")

    defer resp.Body.Close()  
}

