package module

import (
	"fmt"
	"os"

	"golang.org/x/mod/modfile"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// Parse анализирует go.mod файл и возвращает информацию о модуле
func (p *Parser) Parse(goModPath string) (*ModuleInfo, error) {
	// Чтение файла go.mod
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, fmt.Errorf("чтение go.mod: %w", err)
	}

	// Парсинг go.mod
	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, fmt.Errorf("парсинг go.mod: %w", err)
	}

	info := &ModuleInfo{
		Name: modFile.Module.Mod.Path,
	}

	if modFile.Go != nil {
		info.GoVersion = modFile.Go.Version
	} else {
		info.GoVersion = "unknown version"
	}

	// Проходимся по всем зависимостям
	for _, req := range modFile.Require {
		dep := Dependency{
			Path:      req.Mod.Path,
			Current:   req.Mod.Version,
			Indirect:  req.Indirect, // Проверяем, это indirect-зависимости или прямые
			HasUpdate: false,
		}
		info.Dependencies = append(info.Dependencies, dep)
	}

	return info, nil
}
