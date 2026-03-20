package app

import (
	"fmt"
	"go-mod-update-checker/internal/git"
	"go-mod-update-checker/internal/module"
	"go-mod-update-checker/internal/output"
	"go-mod-update-checker/internal/update"
	"os"
	"path/filepath"
)

func Run(repoURL string, jsonOutput bool) error {
	// Создание временной директории
	tempDir, err := os.MkdirTemp("", "go-mod-check-*")
	if err != nil {
		return fmt.Errorf("не удалось создать временную директорию: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			fmt.Fprintf(os.Stderr, "не удалось удалить временную директорию %s: %v\n", tempDir, err)
		}
	}()

	if !jsonOutput {
		fmt.Printf("Клонирование репозитория %s...\n", repoURL)
	}

	// Клонирование репозитория
	cloner := git.NewCloner()
	if err := cloner.Clone(repoURL, tempDir); err != nil {
		return fmt.Errorf("не удалось клонировать репозиторий: %w", err)
	}

	// Поиск go.mod файла
	goModPath, err := findGoMod(tempDir)
	if err != nil {
		return fmt.Errorf("не удалось найти go.mod: %w", err)
	}

	if !jsonOutput {
		fmt.Println("Анализ модуля...")
	}

	// Получение информации о модуле
	parser := module.NewParser()
	moduleInfo, err := parser.Parse(goModPath)
	if err != nil {
		return fmt.Errorf("не удалось проанализировать модуль: %w", err)
	}

	// Проверка доступных обновлений
	checker := update.NewChecker()
	if err := checker.CheckUpdates(moduleInfo, tempDir); err != nil {
		return fmt.Errorf("не удалось проверить обновления: %w", err)
	}

	// Вывод результатов
	printer := output.NewPrinter()
	if jsonOutput {
		return printer.PrintJSON(moduleInfo)
	}

	printer.PrintResults(moduleInfo)
	return nil
}

func findGoMod(rootDir string) (string, error) {
	var goModPath string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "go.mod" {
			goModPath = path
			return filepath.SkipAll
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if goModPath == "" {
		return "", fmt.Errorf("файл go.mod не найден в репозитории")
	}

	return goModPath, nil
}

func PrintUsage() {
	fmt.Println(`Go Module Update Checker - анализирует Go модуль и показывает доступные обновления зависимостей

Использование:
  go run main.go -repo <URL_репозитория> [опции]

Опции:
  -repo    URL Git репозитория (обязательно)
  -json    Вывод в JSON формате
  -help    Показать справку

Пример:
  go run main.go -help
  go run main.go -repo https://github.com/Kweall/chat-api
  go run main.go -repo https://github.com/Kweall/chat-api -json`)
}
