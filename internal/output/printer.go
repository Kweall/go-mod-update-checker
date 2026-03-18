package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"go-mod-update-checker/internal/module"

	"github.com/fatih/color"
)

type Printer struct {
}

func NewPrinter() *Printer {
	return &Printer{}
}

// PrintResults выводит результаты анализа модуля
func (p *Printer) PrintResults(info *module.ModuleInfo) {
	// Цветные принтеры для разных элементов
	headerColor := color.New(color.FgCyan, color.Bold)
	subHeaderColor := color.New(color.FgYellow, color.Bold)
	successColor := color.New(color.FgGreen)

	// Заголовок
	headerColor.Println("\n" + strings.Repeat("═", 60))
	headerColor.Println(" РЕЗУЛЬТАТЫ АНАЛИЗА МОДУЛЯ")
	headerColor.Println(strings.Repeat("═", 60))

	// Основная информация
	fmt.Println(color.CyanString("  Имя модуля:"))
	fmt.Printf("  %s\n", info.Name)

	fmt.Println(color.CyanString("  Версия Go:"))
	fmt.Printf("  %s\n", info.GoVersion)

	// Количество прямых и косвенных зависимостей
	directCount := 0
	indirectCount := 0
	for _, dep := range info.Dependencies {
		if dep.Indirect {
			indirectCount++
		} else {
			directCount++
		}
	}

	fmt.Println(color.CyanString("  Статистика зависимостей:"))
	fmt.Println("  • Прямые:", directCount)
	fmt.Println("  • Косвенные:", indirectCount)
	fmt.Println("  • Всего:", len(info.Dependencies))

	// Зависимости с обновлениями
	updatableCount := 0
	for _, dep := range info.Dependencies {
		if dep.HasUpdate {
			updatableCount++
		}
	}

	subHeaderColor.Printf("\n ДОСТУПНЫЕ ОБНОВЛЕНИЯ: %d\n", updatableCount)
	fmt.Println("\n" + strings.Repeat("═", 60))

	if updatableCount == 0 {
		successColor.Println("\n Все зависимости актуальны!")
	} else {
		for _, dep := range info.Dependencies {
			if dep.HasUpdate {
				p.printDependencyUpdate(dep)
			}
		}
	}

	// Итоговая строка
	fmt.Println("\n" + strings.Repeat("═", 60) + "\n")

	summaryColor := color.New(color.FgHiWhite)
	if updatableCount > 0 {
		summaryColor.Printf("Итого: %d из %d зависимостей можно обновить\n",
			updatableCount, len(info.Dependencies))
	} else {
		successColor.Print("Все зависимости актуальны!  ")
	}
	fmt.Println()
}

// printDependencyUpdate выводит информацию об обновлении зависимости
func (p *Printer) printDependencyUpdate(dep module.Dependency) {
	var text string

	// Цвета для разных типов обновлений
	switch dep.UpdateType {
	case "major":
		text = color.RedString("🔴 MAJOR")
	case "minor":
		text = color.YellowString("🟡 MINOR")
	case "patch":
		text = color.GreenString("🟢 PATCH")
	default:
		text = "⚪ UNKNOWN"
	}

	// Для indirect-зависимости подписываем
	indirectMark := ""
	if dep.Indirect {
		indirectMark = color.HiBlackString(" // indirect")
	}

	fmt.Printf("\n    %s  %s\n",
		color.HiWhiteString(dep.Path), // Название пакета
		indirectMark)

	fmt.Printf("     %s %s → %s\n",
		color.HiBlackString("версия:"),
		color.YellowString(dep.Current),
		color.GreenString(dep.Latest))

	fmt.Printf("     %s %s\n",
		color.HiBlackString("тип:"),
		text)
}

// PrintJSON выводит результаты в формате JSON
func (p *Printer) PrintJSON(info *module.ModuleInfo) error {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка маршалинга JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
