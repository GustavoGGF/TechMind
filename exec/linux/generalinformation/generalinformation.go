package generalinformation

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/host"
)

func GetSO() (string){
	sys := runtime.GOOS

	return sys
}

func GetVersion()(string, error){
	infos, err:= host.Info()
    if err != nil{
        return "", err
    }

	version := infos.PlatformVersion

	return version, nil
}

func GetHostName()(string){
	name, err := os.Hostname()
	if err != nil{
		return ""
	}
	
	return name
}

func GetDistribution()(string, error){
	infos, err:= host.Info()
    if err != nil{
        return "", err
    }

	dis := infos.Platform

	return dis, nil
}

func GetTime()(string){
	date_now := time.Now()
   formated_date := date_now.Format("2006-01-02 15:04")
   
   return formated_date
}

func GetUser()(string){
	user, err := user.Current()
    if err != nil {
        return ""
    }

    return user.Username
}

func GetDomain()(string, error){
    cmd := exec.Command("hostname", "--domain")
    output, err := cmd.CombinedOutput() // Captura a saída e erro juntos
    if err != nil {
        // Verifica se o erro é relacionado a uma opção inválida
        if strings.Contains(string(output), "illegal option") || strings.Contains(string(output), "invalid option") {
            return "", nil // Retorna uma string vazia e sem erro
        }
        return "", fmt.Errorf("erro ao obter domínio: %w", err)
    }
    
    domain := strings.TrimSpace(string(output))
    
    return domain, nil
}