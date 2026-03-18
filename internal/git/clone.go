package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Cloner struct{}

func NewCloner() *Cloner {
	return &Cloner{}
}

// Clone клонирует репозиторий по URL в указанную директорию
func (c *Cloner) Clone(repoURL, destDir string) error {
	// Проверяем и нормализуем URL
	if !strings.HasSuffix(repoURL, ".git") && !strings.Contains(repoURL, "github.com") {
		if strings.Contains(repoURL, "github.com") && !strings.HasSuffix(repoURL, ".git") {
			repoURL = repoURL + ".git"
		}
	}

	// Клонируем репозиторий без истории
	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone: %w", err)
	}

	return nil
}
