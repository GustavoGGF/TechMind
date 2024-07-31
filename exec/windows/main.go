package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/shirou/gopsutil/host"
	"github.com/yusufpapurcu/wmi"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var (
	modntdll          = syscall.NewLazyDLL("ntdll.dll")
	procRtlGetVersion = modntdll.NewProc("RtlGetVersion")
)

var (
	modkernel32         = syscall.NewLazyDLL("kernel32.dll")
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
	System                  string                   `json:"system"`
	Name                    string                   `json:"name"`
	Distribution            string                   `json:"distribution"`
	InterfaceInternet       string                   `json:"interfaceInternet"`
	MacAddress              string                   `json:"macAddress"`
	InsertionDate           string                   `json:"insertionDate"`
	CurrentUser             string                   `json:"currentUser"`
	PlatformVersion         string                   `json:"platformVersion"`
	Domain                  string                   `json:"domain"`
	IP                      string                   `json:"ip"`
	Manufacturer            string                   `json:"manufacturer"`
	Model                   string                   `json:"model"`
	SerialNumber            string                   `json:"serialNumber"`
	MaxCapacityMemory       string                   `json:"maxCapacityMemory"`
	NumberOfDevices         string                   `json:"numberOfDevices"`
	HardDiskModel           string                   `json:"hardDiskModel"`
	HardDiskSerialNumber    string                   `json:"hardDiskSerialNumber"`
	HardDiskUserCapacity    string                   `json:"hardDiskUserCapacity"`
	HardDiskSataVersion     string                   `json:"hardDiskSataVersion"`
	CPUArchitecture         string                   `json:"cpuArchitecture"`
	CPUOperationMode        string                   `json:"cpuOperationMode"`
	CPUS                    uint32                   `json:"cpus"`
	CPUVendorID             string                   `json:"cpuVendorID"`
	CPUModelName            string                   `json:"cpuModelName"`
	CPUThread               uint32                   `json:"cpuThread"`
	CPUCore                 uint32                   `json:"cpuCore"`
	CPUSocket               int                      `json:"cpuSocket"`
	CPUMaxMHz               uint32                   `json:"cpuMaxMHz"`
	CPUMinMHz               uint32                   `json:"cpuMinMHz"`
	GPUProduct              string                   `json:"gpuProduct"`
	GPUVendorID             string                   `json:"gpuVendorID"`
	GPUBusInfo              string                   `json:"gpuBusInfo"`
	GPULogicalName          string                   `json:"gpuLogicalName"`
	GPUClock                string                   `json:"gpuClock"`
	GPUConfiguration        string                   `json:"gpuConfiguration"`
	AudioDeviceProduct      string                   `json:"audioDeviceProduct"`
	AudioDeviceModel        string                   `json:"audioDeviceModel"`
	BiosVersion             string                   `json:"biosVersion"`
	MotherboardManufacturer string                   `json:"motherboardManufacturer"`
	MotherboardProductName  string                   `json:"motherboardProductName"`
	MotherboardVersion      string                   `json:"motherboardVersion"`
	MotherbaoardSerialName  string                   `json:"motherboardSerialName"`
	MotherboardAssetTag     string                   `json:"motherboardAssetTag"`
	SoftwareNames           []InstalledSoftware      `json:"installedPackages"`
	Memories                []map[string]interface{} `json:"memories"`
}

// Win32_ComputerSystem representa a classe WMI Win32_ComputerSystem
type Win32_ComputerSystem struct {
	Manufacturer string
	Model        string
}

// Win32_BIOS representa a classe WMI Win32_BIOS
type Win32_BIOS struct {
	SerialNumber string
}

// Win32_PhysicalMemoryArray representa a classe WMI Win32_PhysicalMemoryArray
type Win32_PhysicalMemoryArray struct {
	MaxCapacity uint32
}

// Win32_PhysicalMemoryArray representa a classe WMI Win32_PhysicalMemoryArray
type Win32_PhysicalMemoryArray2 struct {
	MemoryDevices uint32
}

// Win32_PhysicalMemory representa a classe WMI Win32_PhysicalMemory
type Win32_PhysicalMemory struct {
	BankLabel     string
	Capacity      uint64
	DeviceLocator string
	MemoryType    uint16
	TypeDetail    uint16
	Speed         uint32
	SerialNumber  string
}

type MemoryInfo struct {
	BankLabel     string `json:"BankLabel"`
	Capacity      uint64 `json:"Capacity"`
	DeviceLocator string `json:"DeviceLocator"`
	MemoryType    uint16 `json:"MemoryType"`
	TypeDetail    uint16 `json:"TypeDetail"`
	Speed         uint32 `json:"Speed"`
	SerialNumber  string `json:"SerialNumber"`
}

// Estrutura que representa os dados que queremos obter do WMI
type Win32_DiskDrive struct {
	Model string
}

// Estrutura para mapear a consulta WMI
type Win32_DiskDrive2 struct {
	SerialNumber string
}

// Estrutura para mapear a consulta WMI
type Win32_DiskDrive3 struct {
	Size string // O campo Size é retornado como string
}

type Win32_Processor struct {
	Architecture              uint16
	Manufacturer              string
	Name                      string
	NumberOfLogicalProcessors uint32
	NumberOfCores             uint32
	SocketDesignation         string
	MaxClockSpeed             uint32
	CurrentClockSpeed         uint32
}

type Win32_OperatingSystem struct {
	OSArchitecture string
}

type Win32_VideoController struct {
	Name                        string
	AdapterCompatibility        string
	DeviceID                    string
	DriverVersion               string
	CurrentHorizontalResolution uint32
	CurrentVerticalResolution   uint32
	AdapterRAM                  uint32
}

type Win32_SoundDevice struct {
	ProductName string
}

type InstalledSoftware struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Vendor  string `json:"vendor"`
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

	// Criar um logger para o arquivo
	logger := log.New(file, "", log.LstdFlags)

	// Escrever mensagem no arquivo de log
	logger.Println(msg)
}

func getMac(dis string) (string, error) {
	if dis == "Windows 10" {
		// Obtém todas as interfaces de rede do sistema
		interfaces, err := net.Interfaces()
		if err != nil {
			logToFile(fmt.Sprintf("erro ao obter interfaces de rede: %v", err))
			return "", nil
		}

		// Itera sobre as interfaces de rede para encontrar o endereço MAC
		for _, iface := range interfaces {
			// Ignora as interfaces loopback e desativadas
			if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
				continue
			}

			// Obtém o endereço MAC da interface
			mac := iface.HardwareAddr
			if mac != nil {
				return mac.String(), nil
			}
		}
		logToFile(fmt.Sprintln("nenhum endereço MAC encontrado"))
		return "", nil
	}
	logToFile(fmt.Sprintln("sistema operacional não suportado ou não especificado"))
	return "", nil
}

func getDomain() (string, error) {
	var buf [windows.MAX_PATH]uint16
	var size uint32 = windows.MAX_PATH

	// Obtém o nome do domínio do sistema
	err := windows.GetComputerNameEx(windows.ComputerNameDnsDomain, &buf[0], &size)
	if err != nil {
		logToFile(fmt.Sprintf("Erro ao obter o nome do domínio: %v\n", err))
		return "", nil
	}

	// Converte o buffer UTF16 para string UTF-8
	domain := syscall.UTF16ToString(buf[:])
	return domain, nil

}

func getIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logToFile(fmt.Sprintf("erro ao obter endereços de rede: %v", err))
		return "", nil
	}

	for _, addr := range addrs {
		// Verifica se o endereço é do tipo *net.IPNet e não é um endereço de loopback
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// Verifica se é um endereço IPv4
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	logToFile(fmt.Sprintln("nenhum endereço IP encontrado"))
	return "", errors.New("nenhum endereço IP encontrado")
}

func getDeviceBrand() (string, string, error) {
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

func getSerialNumber() (string, error) {
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

func getMaxMemoryCapacity() (string, error) {
	var arrays []Win32_PhysicalMemoryArray

	// Consulta para obter a capacidade máxima suportada
	err := wmi.Query("SELECT MaxCapacity FROM Win32_PhysicalMemoryArray", &arrays)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar WMI (PhysicalMemoryArray): %v", err)
	}

	// Assume que há apenas um PhysicalMemoryArray e pega a capacidade máxima suportada
	var maxCapacity uint32
	if len(arrays) > 0 {
		maxCapacity = arrays[0].MaxCapacity
	}

	maxCapacityGB := float64(maxCapacity) / (1024 * 1024)
	// Converte o número para string
	maxCapacityStr := fmt.Sprintf("%.0f GB", maxCapacityGB)

	return maxCapacityStr, nil
}

func getMemorySlots() (string, error) {
	var arrays []Win32_PhysicalMemoryArray2

	// Consulta para obter o número de slots de memória
	err := wmi.Query("SELECT MemoryDevices FROM Win32_PhysicalMemoryArray", &arrays)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar WMI (PhysicalMemoryArray): %v", err)
	}

	// Assume que há apenas um PhysicalMemoryArray e pega o número de slots de memória
	if len(arrays) > 0 {
		return fmt.Sprintf("%d", arrays[0].MemoryDevices), nil
	}

	return "", fmt.Errorf("nenhuma informação encontrada")
}

func getMemoryDetails() ([]Win32_PhysicalMemory, error) {
	var memories []Win32_PhysicalMemory

	// Consulta para obter informações detalhadas da memória RAM
	err := wmi.Query("SELECT BankLabel, Capacity, DeviceLocator, MemoryType, TypeDetail, Speed, SerialNumber FROM Win32_PhysicalMemory", &memories)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar WMI (PhysicalMemory): %v", err)
	}

	return memories, nil
}

// Função para obter o modelo do disco rígido
func getHardDiskModel() ([]string, error) {
	var dst []Win32_DiskDrive
	query := "SELECT Model FROM Win32_DiskDrive"
	err := wmi.Query(query, &dst)
	if err != nil {
		return nil, fmt.Errorf("falha ao executar a consulta WMI: %v", err)
	}

	var models []string
	for _, disk := range dst {
		models = append(models, disk.Model)
	}

	return models, nil
}

// Função para obter o número de série do disco rígido
func getHardDiskSerialNumber() ([]string, error) {
	var dst []Win32_DiskDrive2
	query := "SELECT SerialNumber FROM Win32_DiskDrive"
	err := wmi.Query(query, &dst)
	if err != nil {
		return nil, fmt.Errorf("falha ao executar a consulta WMI: %v", err)
	}

	var serialNumbers []string
	for _, disk := range dst {
		serialNumbers = append(serialNumbers, disk.SerialNumber)
	}

	return serialNumbers, nil
}

// Função para obter a capacidade do disco rígido em GB
func getHardDiskCapacity() ([]float64, error) {
	var dst []Win32_DiskDrive3
	query := "SELECT Size FROM Win32_DiskDrive"
	err := wmi.Query(query, &dst)
	if err != nil {
		return nil, fmt.Errorf("falha ao executar a consulta WMI: %v", err)
	}

	var capacities []float64
	for _, disk := range dst {
		size, err := strconv.ParseUint(disk.Size, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("falha ao converter o tamanho do disco: %v", err)
		}
		capacities = append(capacities, float64(size)/(1024*1024*1024)) // Convertendo bytes para gigabytes
	}

	return capacities, nil
}

func getCPUArchitecture() (string, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no processors found")
	}

	architecture := dst[0].Architecture
	var arch string
	switch architecture {
	case 0:
		arch = "x86"
	case 1:
		arch = "MIPS"
	case 2:
		arch = "Alpha"
	case 3:
		arch = "PowerPC"
	case 5:
		arch = "ARM"
	case 6:
		arch = "ia64"
	case 9:
		arch = "x64"
	default:
		arch = "Unknown"
	}

	return arch, nil
}

func getCPUOperationMode() (string, error) {
	var dst []Win32_OperatingSystem
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no operating systems found")
	}

	return dst[0].OSArchitecture, nil
}

func getCPUCount() (uint32, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, err
	}

	if len(dst) == 0 {
		return 0, fmt.Errorf("no processors found")
	}

	// Somando o número de núcleos de todos os processadores
	var totalCores uint32
	for _, processor := range dst {
		totalCores += processor.NumberOfCores
	}

	return totalCores, nil
}

func getCPUVendorID() (string, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no processors found")
	}

	// Escolhe o primeiro processador
	return dst[0].Manufacturer, nil
}

func getCPUModelName() (string, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no processors found")
	}

	// Escolhe o primeiro processador
	return dst[0].Name, nil
}

func getCPUThreads() (uint32, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, err
	}

	if len(dst) == 0 {
		return 0, fmt.Errorf("no processors found")
	}

	// Escolhe o primeiro processador
	return dst[0].NumberOfLogicalProcessors, nil
}

func getCPUCores() (uint32, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, err
	}

	if len(dst) == 0 {
		return 0, fmt.Errorf("no processors found")
	}

	// Escolhe o primeiro processador
	return dst[0].NumberOfCores, nil
}

func getCPUSockets() (int, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, err
	}

	if len(dst) == 0 {
		return 0, fmt.Errorf("no processors found")
	}

	return len(dst), nil
}

func getCPUMaxMHz() (uint32, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, err
	}

	if len(dst) == 0 {
		return 0, fmt.Errorf("no processors found")
	}

	// Escolhe o primeiro processador
	return dst[0].MaxClockSpeed, nil
}

func getCPUMinMHz() (uint32, error) {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, err
	}

	if len(dst) == 0 {
		return 0, fmt.Errorf("no processors found")
	}

	// Escolhe o primeiro processador
	return dst[0].CurrentClockSpeed, nil
}

func getGPUProduct() (string, error) {
	var dst []Win32_VideoController
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no video controllers found")
	}

	// Escolhe o primeiro controlador de vídeo
	return dst[0].Name, nil
}

func getGPUVendorID() (string, error) {
	var dst []Win32_VideoController
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no video controllers found")
	}

	// Escolhe o primeiro controlador de vídeo
	return dst[0].AdapterCompatibility, nil
}

func getGPUBusInfo() (string, error) {
	var dst []Win32_VideoController
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no video controllers found")
	}

	// Escolhe o primeiro controlador de vídeo
	return dst[0].DeviceID, nil
}

func getGPULogicalName() (string, error) {
	var dst []Win32_VideoController
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no video controllers found")
	}

	// Itera sobre todos os controladores de vídeo para obter nomes lógicos
	for _, controller := range dst {
		fmt.Printf("GPU Logical Name: %s\n", controller.Name)
	}

	// Retorna o nome da primeira GPU encontrada como exemplo
	return dst[0].Name, nil
}

func getGPUClock() (string, error) {
	var dst []Win32_VideoController
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no video controllers found")
	}

	// Exibe a versão do driver como uma aproximação para a frequência do clock
	return dst[0].DriverVersion, nil
}

func getGPUConfiguration() (uint32, uint32, uint32, error) {
	var dst []Win32_VideoController
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return 0, 0, 0, err
	}

	if len(dst) == 0 {
		return 0, 0, 0, fmt.Errorf("no video controllers found")
	}

	// Escolhe o primeiro controlador de vídeo
	controller := dst[0]
	return controller.CurrentHorizontalResolution, controller.CurrentVerticalResolution, controller.AdapterRAM, nil
}

func getAudioDeviceProduct() (string, error) {
	var dst []Win32_SoundDevice
	query := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(query, &dst); err != nil {
		return "", err
	}

	if len(dst) == 0 {
		return "", fmt.Errorf("no audio devices found")
	}

	// Escolhe o primeiro dispositivo de áudio encontrado
	return dst[0].ProductName, nil
}

// Função para executar um comando e retornar a saída como string
func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func getSMBIOS() (string, error) {
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

func getInstalledSoftwareFromWMI() ([]InstalledSoftware, error) {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return nil, err
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, err
	}
	defer wmi.Release()

	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer", nil, `root\cimv2`)
	if err != nil {
		return nil, err
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	resultRaw, err := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM Win32_Product")
	if err != nil {
		return nil, err
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	var softwares []InstalledSoftware
	countVar := oleutil.MustGetProperty(result, "Count")
	count := int(countVar.Val)

	for i := 0; i < count; i++ {
		itemRaw := oleutil.MustCallMethod(result, "ItemIndex", i)
		item := itemRaw.ToIDispatch()
		defer item.Release()

		name := oleutil.MustGetProperty(item, "Name").ToString()
		version := oleutil.MustGetProperty(item, "Version").ToString()
		vendor := oleutil.MustGetProperty(item, "Vendor").ToString()

		software := InstalledSoftware{
			Name:    name,
			Version: version,
			Vendor:  vendor,
		}
		softwares = append(softwares, software)
	}

	return softwares, nil
}

func getInstalledSoftwareFromRegistry() ([]InstalledSoftware, error) {
	var softwares []InstalledSoftware

	// Localização das chaves de registro
	registryPaths := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}

	for _, path := range registryPaths {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.READ)
		if err != nil {
			return nil, err
		}
		defer k.Close()

		subkeys, err := k.ReadSubKeyNames(-1)
		if err != nil {
			return nil, err
		}

		for _, subkey := range subkeys {
			subk, err := registry.OpenKey(k, subkey, registry.READ)
			if err != nil {
				continue
			}
			defer subk.Close()

			name, _, err := subk.GetStringValue("DisplayName")
			if err != nil || name == "" {
				continue
			}

			version, _, _ := subk.GetStringValue("DisplayVersion")
			vendor, _, _ := subk.GetStringValue("Publisher")

			software := InstalledSoftware{
				Name:    name,
				Version: version,
				Vendor:  vendor,
			}
			softwares = append(softwares, software)
		}
	}

	return softwares, nil
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

	hostname := ""

	if computerName, err := getComputerName(); err == nil {
		hostname = computerName
	} else {
		logToFile(fmt.Sprintf("Erro ao obter o nome da máquina: %v", err))
		return
	}

	sys = strings.ReplaceAll(sys, " ", "")
	hostname = strings.ReplaceAll(hostname, " ", "")
	distribution = strings.ReplaceAll(distribution, " ", "")

	url := "http://10.1.1.73:3000/home/computers/post-machines" //env

	date_now := time.Now()

	formated_date := date_now.Format("2006-01-02 15:04")

	macAddress, err := getMac("Windows 10")

	if err != nil {
		logToFile(fmt.Sprintf("Erro ao obter o endereço MAC: %v", err))
		return
	}

	currentUser, err := user.Current()
	if err != nil {
		logToFile(fmt.Sprintf("Erro ao obter o usuário atual: %v", err))
	}

	nameUser := string(currentUser.Username)

	info, err := host.Info()
	if err != nil {
		logToFile(fmt.Sprintf("Error: %v", err))
		return
	}

	version := info.PlatformVersion

	domain, err := getDomain()

	if err != nil {
		logToFile(fmt.Sprintf("Erro ao obter o dominio: %v", err))
		return
	}

	ip, err := getIP()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %v", err))
	}

	manufacturer, model, err := getDeviceBrand()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %v", err))
	}

	serialNumber, err := getSerialNumber()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %v", err))
	}

	maxCapacityMemory, err := getMaxMemoryCapacity()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %v", err))
	}

	memorySlots, err := getMemorySlots()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %d", err))
		return
	}

	mem, err := getMemoryDetails()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %d", err))
		return
	}

	// Armazena informações em um array
	var memoryArray []map[string]interface{}

	for _, memory := range mem {
		memoryInfo := map[string]interface{}{
			"BankLabel":     memory.BankLabel,
			"Capacity":      memory.Capacity,
			"DeviceLocator": memory.DeviceLocator,
			"MemoryType":    memory.MemoryType,
			"TypeDetail":    memory.TypeDetail,
			"Speed":         memory.Speed,
			"SerialNumber":  memory.SerialNumber,
		}
		memoryArray = append(memoryArray, memoryInfo)
	}

	models, err := getHardDiskModel()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %v", err))
	}
	var modelHardDisk string
	for _, model := range models {
		modelHardDisk = model
	}

	hardDiskSerialNumber, err := getHardDiskSerialNumber()
	if err != nil {
		logToFile(fmt.Sprintf("Erro: %v", err))
	}

	var serialNumberHardDisk string
	for _, serialNumber := range hardDiskSerialNumber {
		serialNumberHardDisk = serialNumber
	}

	capacities, err := getHardDiskCapacity()
	if err != nil {
		logToFile(fmt.Sprint("Erro:", err))
		return
	}
	var capacityHardDisk string
	for _, capacity := range capacities {
		capacityHardDisk = fmt.Sprintf("%.2f", capacity)
	}

	arch, err := getCPUArchitecture()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU architecture: %v", err))
	}

	operationMode, err := getCPUOperationMode()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU operation mode: %v", err))
	}

	cpuCount, err := getCPUCount()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU count: %v", err))
	}

	vendorID, err := getCPUVendorID()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU Vendor ID: %v", err))
	}

	modelName, err := getCPUModelName()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU Model Name: %v", err))
	}

	threads, err := getCPUThreads()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU threads: %v", err))
	}

	cores, err := getCPUCores()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU cores: %v", err))
	}

	sockets, err := getCPUSockets()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU sockets: %v", err))
	}

	maxMHz, err := getCPUMaxMHz()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU Max MHz: %v", err))
	}

	minMHz, err := getCPUMinMHz()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get CPU Min MHz: %v", err))
	}

	gpuProduct, err := getGPUProduct()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get GPU product: %v", err))
	}

	vendorIDGPU, err := getGPUVendorID()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get GPU Vendor ID: %v", err))
	}

	busInfo, err := getGPUBusInfo()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get GPU Bus Info: %v", err))
	}

	gpuLogicalName, err := getGPULogicalName()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get GPU Logical Name: %v", err))
	}

	clock, err := getGPUClock()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get GPU Clock: %v", err))
	}

	horizRes, vertRes, ram, err := getGPUConfiguration()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get GPU configuration: %v", err))
	}

	// Formata a string com as informações da GPU
	configurationGPU := fmt.Sprintf("Resolution %dx%d, RAM %d MB", horizRes, vertRes, ram/1024/1024)

	product, err := getAudioDeviceProduct()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get audio device product: %v", err))
	}

	smbiosInfo, err := getSMBIOS()
	if err != nil {
		logToFile(fmt.Sprintf("Failed to get SMBIOS information: %v", err))
	}

	fmt.Sprintln(smbiosInfo)

	wmiSoftware, err := getInstalledSoftwareFromWMI()
	if err != nil {
		log.Printf("Error querying WMI: %v", err)
	}

	registrySoftware, err := getInstalledSoftwareFromRegistry()
	if err != nil {
		log.Printf("Error querying Registry: %v", err)
	}

	// Combinar as listas de software, removendo duplicatas
	softwareMap := make(map[string]InstalledSoftware)
	for _, software := range append(wmiSoftware, registrySoftware...) {
		key := strings.ToLower(software.Name)
		softwareMap[key] = software
	}

	var combinedSoftware []InstalledSoftware
	for _, software := range softwareMap {
		combinedSoftware = append(combinedSoftware, software)
	}

	jsonData := Data{
		System:               sys,
		Name:                 hostname,
		Distribution:         distribution,
		InsertionDate:        formated_date,
		MacAddress:           macAddress,
		CurrentUser:          nameUser,
		PlatformVersion:      version,
		Domain:               domain,
		IP:                   ip,
		Manufacturer:         manufacturer,
		Model:                model,
		SerialNumber:         serialNumber,
		MaxCapacityMemory:    maxCapacityMemory,
		NumberOfDevices:      memorySlots,
		Memories:             memoryArray,
		HardDiskModel:        modelHardDisk,
		HardDiskSerialNumber: serialNumberHardDisk,
		HardDiskUserCapacity: capacityHardDisk,
		CPUArchitecture:      arch,
		CPUOperationMode:     operationMode,
		CPUS:                 cpuCount,
		CPUVendorID:          vendorID,
		CPUModelName:         modelName,
		CPUThread:            threads,
		CPUCore:              cores,
		CPUSocket:            sockets,
		CPUMaxMHz:            maxMHz,
		CPUMinMHz:            minMHz,
		GPUProduct:           gpuProduct,
		GPUVendorID:          vendorIDGPU,
		GPUBusInfo:           busInfo,
		GPULogicalName:       gpuLogicalName,
		GPUClock:             clock,
		GPUConfiguration:     configurationGPU,
		AudioDeviceProduct:   product,
		SoftwareNames:        combinedSoftware,
	}

	requestBody, err := json.Marshal(jsonData)
	if err != nil {
		logToFile(fmt.Sprintf("Erro ao montar o json: %v", err))
		return
	}

	resp, erro := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if erro != nil {
		logToFile(fmt.Sprintf("Erro ao fazer o post: %v", err))
	}

	// Erro response 400 gerar aviso na tela

	defer resp.Body.Close()

	// Ler o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logToFile(fmt.Sprintf("Erro ao ler o corpo da resposta: %v", err))
		return
	}

	// Verificar o código de status da resposta
	if resp.StatusCode != http.StatusOK {
		logToFile(fmt.Sprintf("Resposta com status code %d: %s", resp.StatusCode, body))
		if resp.StatusCode == http.StatusBadRequest {
			fmt.Println("Erro: Resposta 400 - Solicitação inválida.")
		} else {
			fmt.Printf("Erro: Resposta %d\n", resp.StatusCode)
		}
		return
	}

	// Exibir a resposta no console
	fmt.Printf("Resposta: %s\n", body)

}
