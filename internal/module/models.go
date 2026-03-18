package module

// Информация о зависимости
type Dependency struct {
	Path       string `json:"path"`
	Current    string `json:"current"`
	Latest     string `json:"latest"`
	UpdateType string `json:"update_type"` // версия v.X.Y.Z: X - major, Y - minor, Z - patch
	Indirect   bool   `json:"indirect"`
	HasUpdate  bool   `json:"has_update"`
}

// Информация о модуле
type ModuleInfo struct {
	Name         string       `json:"name"`
	GoVersion    string       `json:"go_version"`
	Dependencies []Dependency `json:"dependencies"`
	TempDir      string       `json:"-"`
}
