package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"go-mod-update-checker/internal/git"
	"go-mod-update-checker/internal/module"
	"go-mod-update-checker/internal/output"
	"go-mod-update-checker/internal/update"
)

func main() {
	// Флаги
	var repoURL string  // -repo - ссылка на репозиторий github
	var help bool       // -help - для справки
	var jsonOutput bool // -json - вывод зависимостей в формате JSON

	flag.StringVar(&repoURL, "repo", "", "URL Git репозитория (например, https://github.com/user/repo)")
	flag.BoolVar(&help, "help", false, "Показать справку")
	flag.BoolVar(&jsonOutput, "json", false, "Вывод в JSON формате")
	flag.Parse()

	if help || repoURL == "" {
		printUsage()
		return
	}

	if err := run(repoURL, jsonOutput); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}

func run(repoURL string, jsonOutput bool) error {
	// Создание временной директории
	tempDir, err := os.MkdirTemp("", "go-mod-check-*")
	if err != nil {
		return fmt.Errorf("не удалось создать временную директорию: %w", err)
	}
	defer os.RemoveAll(tempDir)

	fmt.Printf("Клонирование репозитория %s...\n", repoURL)

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

	fmt.Println("Анализ модуля...")

	// Получение информации о модуле
	parser := module.NewParser()
	moduleInfo, err := parser.Parse(goModPath)
	if err != nil {
		return fmt.Errorf("не удалось проанализировать модуль: %w", err)
	}
	moduleInfo.TempDir = tempDir

	// Проверка доступных обновлений
	checker := update.NewChecker()
	if err := checker.CheckUpdates(moduleInfo); err != nil {
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

func printUsage() {
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
