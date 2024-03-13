package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/gofiber/fiber/v2/log"
)

type SystemInfo struct {
    System   string
    Hostname string
}

type Data struct {
    System, Name string
}

func main() {
    sys := runtime.GOOS

    if sys == "windows" {

        info := sys

        url := "http://10.1.9.4:3000/api/machines"

        jsonData := Data{System: info}

        requestBody, err := json.Marshal(jsonData)
        if err != nil{
            log.Info(err)
        }
        
        resp, erro := http.Post(url, "application/json", bytes.NewBuffer(requestBody))

        if erro != nil {
	        archive, err := os.Create("dados.txt")
            if err != nil {
                log.Info("Erro ao criar o arquivo:", err)
                return
            }
        	defer archive.Close()

	        _, err = archive.WriteString(erro.Error())
            if err != nil {
                log.Info("Erro ao escrever no arquivo:", err)
                return
            }
        }

        defer resp.Body.Close()


    } else if sys == "linux" {
        url := "http://10.1.9.0:3000/api/machines" //env

        name, err := os.Hostname()

        if err != nil {
            fmt.Println(err)
        }

        jsonData := Data{System: sys, Name: name}

        fmt.Println(jsonData)

        requestBody, err := json.Marshal(jsonData)
        if err != nil{
            fmt.Println(err)
        }

        resp, erro := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
        if erro != nil{
            fmt.Println(erro)
        }

        defer resp.Body.Close()
        return
    }
}