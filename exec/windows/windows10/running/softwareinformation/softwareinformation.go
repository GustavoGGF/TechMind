package softwareinformation

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"golang.org/x/sys/windows/registry"
)

// Struct que armazena informações dos softwares
type InstalledSoftware struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Vendor  string `json:"vendor"`
}

func GetInstalledSoftwareFromWMI() ([]InstalledSoftware, error) {
	// Inicializa o OLE
	if err := ole.CoInitialize(0); err != nil {
		return nil, fmt.Errorf("falha ao inicializar OLE: %v", err)
	}
	defer ole.CoUninitialize()

	// Cria objeto SWbemLocator
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return nil, fmt.Errorf("falha ao criar objeto SWbemLocator: %v", err)
	}
	defer unknown.Release()

	// Obtém interface IDispatch
	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter interface IDispatch: %v", err)
	}
	defer wmi.Release()

	// Conecta ao namespace WMI root\cimv2
	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer", nil, `root\cimv2`)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao servidor WMI: %v", err)
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// Executa consulta WMI para obter software instalado
	resultRaw, err := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM Win32_Product")
	if err != nil {
		return nil, fmt.Errorf("falha ao executar consulta WMI: %v", err)
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	// Obtém contagem de itens
	countVar, err := oleutil.GetProperty(result, "Count")
	if err != nil {
		return nil, fmt.Errorf("falha ao obter contagem de itens: %v", err)
	}
	count := int(countVar.Val)
	if count == 0 {
		return nil, fmt.Errorf("nenhum software instalado encontrado")
	}

	var softwares []InstalledSoftware
	for i := 0; i < count; i++ {
		itemRaw, err := oleutil.CallMethod(result, "ItemIndex", i)
		if err != nil {
			return nil, fmt.Errorf("falha ao obter item %d: %v", i, err)
		}
		item := itemRaw.ToIDispatch()
		defer item.Release()

		nameVar, err := oleutil.GetProperty(item, "Name")
		if err != nil {
			return nil, fmt.Errorf("falha ao obter nome do software para item %d: %v", i, err)
		}
		versionVar, err := oleutil.GetProperty(item, "Version")
		if err != nil {
			return nil, fmt.Errorf("falha ao obter versão do software para item %d: %v", i, err)
		}
		vendorVar, err := oleutil.GetProperty(item, "Vendor")
		if err != nil {
			return nil, fmt.Errorf("falha ao obter fornecedor do software para item %d: %v", i, err)
		}

		software := InstalledSoftware{
			Name:    nameVar.ToString(),
			Version: versionVar.ToString(),
			Vendor:  vendorVar.ToString(),
		}
		softwares = append(softwares, software)
	}

	return softwares, nil
}

func GetInstalledSoftwareFromRegistry() ([]InstalledSoftware, error) {
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
