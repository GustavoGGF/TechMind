package software

// Struct que armazena informações dos softwares
type InstalledSoftware struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Vendor  string `json:"vendor"`
}