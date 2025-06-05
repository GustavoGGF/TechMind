package logger

import (
	"fmt"
	"log"
	"os"
)

// Função que cria um arquivo de log


func LogToFile(msg string) {
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