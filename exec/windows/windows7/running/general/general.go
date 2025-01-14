package general

import (
	"fmt"
	"os/user"
	"syscall"
	"unsafe"

	"github.com/StackExchange/wmi"
	"golang.org/x/sys/windows"
)

// Declaração de variáveis para interagir com a biblioteca kernel32.dll do Windows
var (
	// kernel32 é um manipulador para a DLL kernel32.dll, que contém funções do sistema operacional Windows
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	// procGetComputerName é um ponteiro para a função GetComputerNameW na DLL kernel32,
	// que é usada para obter o nome do computador.
	procGetComputerName = kernel32.NewProc("GetComputerNameW")
)

// Win32_OperatingSystem representa as informações básicas do sistema operacional Windows,
// obtidas via consultas WMI (Windows Management Instrumentation).
type Win32_OperatingSystem struct {
	Name        string // Nome completo do sistema operacional (por exemplo, "Microsoft Windows 10 Pro")
	Version     string // Versão do sistema operacional (por exemplo, "10.0.19042")
	Caption     string // Descrição legível do sistema (por exemplo, "Windows 10 Pro")
	BuildNumber string // Número de build do sistema operacional (por exemplo, "19042")
}

// Win32_ComputerSystem representa informações sobre o sistema de computador,
// incluindo o fabricante e o modelo.
type Win32_ComputerSystem struct {
	Manufacturer string // Nome do fabricante do computador
	Model        string // Modelo do computador
}

// Win32_BIOS representa informações sobre o BIOS do sistema,
// especificamente o número de série do BIOS.
type Win32_BIOS struct {
	SerialNumber string // Número de série do BIOS
}

// Função que pega o nome do computador
func GetComputerName() (string, error) {
	var nSize uint32 = 256
	nameBuf := make([]uint16, nSize)
	ret, _, err := procGetComputerName.Call(uintptr(unsafe.Pointer(&nameBuf[0])), uintptr(unsafe.Pointer(&nSize)))
	if ret == 0 {
		return "", err
	}
	return syscall.UTF16ToString(nameBuf[:nSize]), nil
}

// GetWindowsDistribution consulta informações do sistema operacional Windows 
// usando WMI e retorna uma string contendo a descrição (Caption) do sistema operacional.
func GetWindowsDistribution() (string, string, error) {
	var osInfo []Win32_OperatingSystem
	// Executa a consulta WMI para obter o nome, versão, legenda e número de build do sistema operacional
	err := wmi.Query("SELECT Name, Version, Caption, BuildNumber FROM Win32_OperatingSystem", &osInfo)
	if err != nil {
		// Retorna o erro se houver problemas durante a execução da consulta WMI
		return "", "", err
	}

	// Verifica se a consulta retornou ao menos um resultado
	if len(osInfo) > 0 {
		// Retorna a legenda (Caption) do primeiro item, que descreve o sistema operacional
		os := fmt.Sprintf(osInfo[0].Caption)
		version := osInfo[0].Version
		return os, version, nil
	}

	// Retorna um erro caso as informações do sistema operacional não sejam encontradas
	return "","", fmt.Errorf("não foi possível encontrar informações do sistema operacional")
}

// GetCurrentUser obtém o nome de usuário do usuário atual do sistema.
func GetCurrentUser() (string, error) {
	// Obtém as informações do usuário atual usando o pacote "os/user"
	usr, err := user.Current()
	if err != nil {
		// Retorna o erro caso ocorra algum problema ao obter o usuário
		return "", err
	}
	// Retorna o nome de usuário do usuário atual
	return usr.Username, nil
}

// GetDomain obtém o nome do domínio DNS ao qual o computador está associado.
func GetDomain() (string, error) {
	// Cria um buffer para armazenar o nome do domínio com o tamanho máximo definido em windows.MAX_PATH
	var buf [windows.MAX_PATH]uint16
	var size uint32 = windows.MAX_PATH

	// Chama a função GetComputerNameEx para obter o nome do domínio do computador (DNS)
	err := windows.GetComputerNameEx(windows.ComputerNameDnsDomain, &buf[0], &size)
	if err != nil {
		// Retorna um erro formatado se a função falhar ao obter o nome do domínio
		return "", fmt.Errorf("Erro ao obter o nome do domínio: %v", err)
	}

	// Converte o buffer UTF16 para uma string UTF-8 legível
	domain := syscall.UTF16ToString(buf[:])
	// Retorna o nome do domínio obtido
	return domain, nil
}

// GetDeviceBrand obtém o fabricante e o modelo do sistema usando consultas WMI.
func GetDeviceBrand() (string, string, error) {
	// Declara uma variável para armazenar os resultados da consulta WMI
	var dst []Win32_ComputerSystem
	
	// Define a consulta WMI para obter o fabricante e o modelo do sistema
	query := "SELECT Manufacturer, Model FROM Win32_ComputerSystem"
	
	// Executa a consulta WMI e armazena o resultado em dst
	err := wmi.Query(query, &dst)
	if err != nil {
		// Retorna um erro formatado se a consulta falhar
		return "", "", fmt.Errorf("erro ao consultar WMI: %v", err)
	}
	
	// Verifica se a consulta retornou algum resultado
	if len(dst) == 0 {
		// Retorna um erro se não houver informações disponíveis
		return "", "", fmt.Errorf("nenhuma informação encontrada")
	}
	
	// Retorna o fabricante e o modelo do sistema encontrados
	return dst[0].Manufacturer, dst[0].Model, nil
}

// GetSerialNumber obtém o número de série do BIOS do sistema usando consultas WMI.
func GetSerialNumber() (string, error) {
	// Declara uma variável para armazenar os resultados da consulta WMI
	var dst []Win32_BIOS
	
	// Define a consulta WMI para obter o número de série do BIOS
	query := "SELECT SerialNumber FROM Win32_BIOS"
	
	// Executa a consulta WMI e armazena o resultado em dst
	err := wmi.Query(query, &dst)
	if err != nil {
		// Retorna um erro formatado se a consulta falhar
		return "", fmt.Errorf("erro ao consultar WMI: %v", err)
	}
	
	// Verifica se a consulta retornou algum resultado
	if len(dst) == 0 {
		// Retorna um erro se não houver informações disponíveis
		return "", fmt.Errorf("nenhuma informação encontrada")
	}
	
	// Retorna o número de série do BIOS encontrado
	return dst[0].SerialNumber, nil
}
