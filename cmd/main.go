package main

import (
	"flag"
	"fmt"
	"go-mod-update-checker/internal/app"
	"os"
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
		app.PrintUsage()
		return
	}

	if err := app.Run(repoURL, jsonOutput); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}
