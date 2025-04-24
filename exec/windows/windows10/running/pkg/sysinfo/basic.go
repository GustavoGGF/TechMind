package sysinfo

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"github.com/yusufpapurcu/wmi"
	"golang.org/x/sys/windows"
)

// Struct que pega o nome do computador
var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procGetComputerName = kernel32.NewProc("GetComputerNameW")
)

// Win32_ComputerSystem representa a classe WMI Win32_ComputerSystem
type Win32_ComputerSystem struct {
	Manufacturer string
	Model        string
}

// Win32_BIOS representa a classe WMI Win32_BIOS
type Win32_BIOS struct {
	SerialNumber string
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

// Função para extrair o valor da Edição do Windows
func ExtractWindowsEdition(output string) (string, error) {
	// Expressão regular para encontrar a linha da Edição do Windows
	re := regexp.MustCompile(`Edicao do Windows\s*:\s*(.*)`)
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return "", fmt.Errorf("não foi possível encontrar a Edição do Windows")
	}
	edition := match[1]
	// Verifique se a edição é inválida (por exemplo, "Versao do Service Pack")
	if edition == "Versao do Service Pack" || edition == "" {
		// Segunda tentativa com a expressão regular alternativa
		re = regexp.MustCompile(`(?i)Edição do Windows\s*:\s*(.*)`)
		match = re.FindStringSubmatch(output)
		if len(match) < 2 {
			return "", fmt.Errorf("não foi possível encontrar a Edição do Windows na segunda tentativa")
		}
		edition = match[1]
	}
	return edition, nil
}

// Função que gera muitas informações do windows
func GetWindowsInfo() (string, error) {
	var script = `
function Get-WindowsKey {
    param ($targets = ".")
    $hklm = 2147483650
    $regPath = "Software\Microsoft\Windows NT\CurrentVersion"
    $regValue = "DigitalProductId4"
    Foreach ($target in $targets) {
        $productKey = $null
        $win32os = $null
        $wmi = [WMIClass]"\\$target\root\default:stdRegProv"
        $data = $wmi.GetBinaryValue($hklm,$regPath,$regValue)
        $binArray = ($data.uValue)[52..66]
        $charsArray = "B","C","D","F","G","H","J","K","M","P","Q","R","T","V","W","X","Y","2","3","4","6","7","8","9"
        For ($i = 24; $i -ge 0; $i--) {
            $k = 0
            For ($j = 14; $j -ge 0; $j--) {
                $k = $k * 256 -bxor $binArray[$j]
                $binArray[$j] = [math]::truncate($k / 24)
                $k = $k % 24
            }
            $productKey = $charsArray[$k] + $productKey
            If (($i % 5 -eq 0) -and ($i -ne 0)) {
                $productKey = "-" + $productKey
            }
        }
        $win32os = Get-WmiObject Win32_OperatingSystem -computer $target
        $obj = New-Object Object
        $obj | Add-Member Noteproperty "Nome do Computador" -value $target
        $obj | Add-Member Noteproperty "Edicao do Windows" -value $win32os.Caption
        $obj | Add-Member Noteproperty "Versao do Service Pack" -value $win32os.CSDVersion
        $obj | Add-Member Noteproperty "Arquitetura" -value $win32os.OSArchitecture
        $obj | Add-Member Noteproperty "Número da Versão" -value $win32os.BuildNumber
        $obj | Add-Member Noteproperty "Registado para" -value $win32os.RegisteredUser
        $obj | Add-Member Noteproperty "Canal de origem (ProductID)" -value $win32os.SerialNumber
        $obj | Add-Member Noteproperty "Chave do produto (ProductKey)" -value $productkey
        $obj
    }
}
Get-WindowsKey
`
	// Tenta executar o comando PowerShell padrão
	out, err := TryCommand("powershell", "-WindowStyle", "Hidden", "-Command", script)
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found in %PATH%") {
			// Se falhar com erro específico, tenta com o caminho absoluto
			out, err = TryCommand("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "-WindowStyle", "Hidden", "-Command", script)
		}
		if err != nil {
			return "", err
		}
	}

	return out, nil
}

// Função para tentar executar um comando e verificar se ocorre um erro específico
func TryCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// Função para obter o dominio
func GetDomain() (string, error) {
	var buf [windows.MAX_PATH]uint16
	var size uint32 = windows.MAX_PATH

	// Obtém o nome do domínio do sistema
	err := windows.GetComputerNameEx(windows.ComputerNameDnsDomain, &buf[0], &size)
	if err != nil {
		return "", fmt.Errorf("Erro ao obter o nome do domínio: %v", err)
	}

	// Converte o buffer UTF16 para string UTF-8
	domain := syscall.UTF16ToString(buf[:])
	return domain, nil
}

// Descobre a Marca e o Modelo do equipamento
func GetDeviceBrand() (string, string, error) {
	var dst []Win32_ComputerSystem
	query := "SELECT Manufacturer, Model FROM Win32_ComputerSystem"
	err := wmi.Query(query, &dst)
	if err != nil {
		return "", "", fmt.Errorf("erro ao consultar WMI: %v", err)
	}
	if len(dst) == 0 {
		return "", "", fmt.Errorf("nenhuma informação encontrada")
	}
	return dst[0].Manufacturer, dst[0].Model, nil
}

// Obtem o Serial Number
func GetSerialNumber() (string, error) {
	var dst []Win32_BIOS
	query := "SELECT SerialNumber FROM Win32_BIOS"
	err := wmi.Query(query, &dst)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar WMI: %v", err)
	}
	if len(dst) == 0 {
		return "", fmt.Errorf("nenhuma informação encontrada")
	}
	return dst[0].SerialNumber, nil
}

// Função para extrair o valor da Edição do Windows
func ExtractWindowsLicense(output string) (string, error) {
	re := regexp.MustCompile(`Chave do produto \(ProductKey\)\s*:\s*(.*)`)
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return "", fmt.Errorf("não foi possível encontrar a Chave do Produto")
	}
	return match[1], nil
}

// Função para executar um comando e retornar a saída como string
func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	// Configura o processo para não abrir uma janela de console
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out // Captura erros na mesma saída

	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func GetSMBIOS() (string, error) {
	// Executa o comando WMIC para obter informações SMBIOS
	output, err := runCommand("wmic", "bios", "get", "serialnumber,version,manufacturer")
	if err != nil {
		return "", err
	}

	// Processa a saída para remover linhas em branco e espaços extras
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n"), nil
}

func GetSys()(string){
	sys := runtime.GOOS
	return sys
}