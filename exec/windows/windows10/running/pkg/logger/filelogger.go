package logger

import (
	"log"
	"os"
	"path/filepath"
)

// Função que cria um arquivo de log

func LogToFile(msg string) {
	logFile := getSafeLogPath()

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Fallback: tentar logar em uma pasta temporária
		tempFile := filepath.Join(os.TempDir(), "techmind_log_fallback.txt")
		_ = os.WriteFile(tempFile, []byte("Erro ao abrir arquivo de log principal: "+err.Error()+"\n"), 0644)
		return
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	logger.Println(msg)
}

func getSafeLogPath() string {
	// Caminho seguro em AppData\Local
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		// Se falhar, usar pasta temporária
		return filepath.Join(os.TempDir(), "techmind_log_fallback.txt")
	}

	logDir := filepath.Join(appData, "techmind")
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return filepath.Join(os.TempDir(), "techmind_log_fallback.txt")
	}

	// Nome do log
	return filepath.Join(logDir, "log.txt")
}