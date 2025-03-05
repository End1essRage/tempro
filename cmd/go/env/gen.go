package cmd_go_env

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		//считываем из main.go текущей папки
		existing, err := getExisting()
		if err != nil {
			exitWithError("Ошибка считывания из main.go файла %s", err)
		}

		//генерим .env или .yaml
		isYml, err := cmd.Flags().GetBool("yml")
		if err != nil {
			exitWithError("ошибка чтения флага %v", err)
		}

		//получаем текущую папку
		currentDir, err := os.Getwd()
		if err != nil {
			exitWithError("Ошибка получения текущей директории: %v", err)
		}

		if isYml {
			ymlPath := filepath.Join(currentDir, "env.yml")
			if err := generateYmlFile(ymlPath, existing); err != nil {
				exitWithError("Ошибка генерации .yml: %v", err)
			}
			fmt.Printf("✅ Файл .yml успешно создан в %s\n", ymlPath)
		} else {
			envPath := filepath.Join(currentDir, ".env")
			if err := generateEnvFile(envPath, existing); err != nil {
				exitWithError("Ошибка генерации .env: %v", err)
			}
			fmt.Printf("✅ Файл .env успешно создан в %s\n", envPath)
		}
	},
}

func init() {
	genCmd.Flags().BoolP("yml", "y", false, "для генерации файла env.yml")
}

// Генерируем .yml файл
func generateYmlFile(path string, fields []string) error {
	//overrides
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "env:\n"); err != nil {
		return err
	}

	for _, field := range fields {
		if _, err := fmt.Fprintf(file, "- name: %s\n  value: 0\n", field); err != nil {
			return err
		}
	}
	return nil
}

// Генерируем .env файл
func generateEnvFile(path string, fields []string) error {
	//overrides
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, field := range fields {
		if _, err := fmt.Fprintf(file, "%s=\n", field); err != nil {
			return err
		}
	}
	return nil
}
