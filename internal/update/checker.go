package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"go-mod-update-checker/internal/module"

	"golang.org/x/mod/semver"
)

type Checker struct {
	httpClient *http.Client
}

func NewChecker() *Checker {
	return &Checker{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CheckUpdates проверяет обновления для всех зависимостей в модуле
func (c *Checker) CheckUpdates(info *module.ModuleInfo) error {
	// Временная директория с репозиторием
	originalDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(info.TempDir); err != nil {
		return err
	}

	fmt.Println("Проверка доступных обновлений...")

	for i, dep := range info.Dependencies {
		latest, err := c.getLatestVersion(dep.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Предупреждение: не удалось получить последнюю версию для %s: %v\n", dep.Path, err)
			continue
		}

		if latest != "" && latest != dep.Current {
			info.Dependencies[i].Latest = latest
			info.Dependencies[i].HasUpdate = true
			info.Dependencies[i].UpdateType = c.getUpdateType(dep.Current, latest)
		}
	}

	return nil
}

// getLatestVersion получает последнюю версию модуля
func (c *Checker) getLatestVersion(modulePath string) (string, error) {
	// Используем go list для получения последней версии
	cmd := exec.Command("go", "list", "-m", "--versions", modulePath)
	output, err := cmd.Output()
	if err != nil {
		// Если не получилось, то пробуем через proxy
		return c.getLatestVersionWithProxy(modulePath)
	}

	versions := strings.Fields(string(output))
	if len(versions) <= 1 {
		return "", nil // Нет иных версий
	}

	// Новейшая версия  лежит в самом конце
	return versions[len(versions)-1], nil
}

// getLatestVersionWithProxy получает последнюю версию через GOPROXY
func (c *Checker) getLatestVersionWithProxy(modulePath string) (string, error) {
	// Используем GOPROXY для получения информации
	proxyURL := os.Getenv("GOPROXY")
	if proxyURL == "" {
		proxyURL = "https://proxy.golang.org"
	}

	// Формируем URL для запроса к proxy
	url := fmt.Sprintf("%s/%s/@v/list", proxyURL, modulePath)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	versions := strings.Split(strings.TrimSpace(string(body)), "\n")
	if len(versions) == 0 {
		return "", nil
	}

	// Сортируем версии и возвращаем последнюю
	semver.Sort(versions)
	return versions[len(versions)-1], nil
}

// getUpdateType определяет тип обновления
func (c *Checker) getUpdateType(current, latest string) string {
	// Определяем тип обновления
	if !semver.IsValid(current) || !semver.IsValid(latest) {
		return "unknown"
	}

	currentMajor := semver.Major(current)
	latestMajor := semver.Major(latest)

	if currentMajor != latestMajor {
		return "major"
	}

	currentMinor := semver.MajorMinor(current)
	latestMinor := semver.MajorMinor(latest)

	if currentMinor != latestMinor {
		return "minor"
	}

	return "patch"
}
