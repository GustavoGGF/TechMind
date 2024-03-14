package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/shirou/gopsutil/host"
)

type Data struct {
    System, Name, Distribution string
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

    jsonData := Data{System: sys, Name: name, Distribution: infos.Platform}

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

