package cmd_go_env

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		//получаем текущую папку
		currentDir, err := os.Getwd()
		if err != nil {
			exitWithError("Ошибка получения текущей директории: %v", err)
		}

		//ищем файл main.go
		mainGoPath := filepath.Join(currentDir, "main.go")
		if !fileExists(mainGoPath) {
			exitWithError("Файл main.go не найден в %s", currentDir)
		}

		//вычитываем структуру Config, получаем список переменных окружения
		configFields, err := parseConfigStruct(mainGoPath)
		if err != nil {
			exitWithError("%v", err)
		}

		//генерим .env или .yaml
		isYml, err := cmd.Flags().GetBool("yml")
		if err != nil {
			exitWithError("ошибка чтения флага %v", err)
		}

		if isYml {
			ymlPath := filepath.Join(currentDir, "env.yml")
			if err := generateYmlFile(ymlPath, configFields); err != nil {
				exitWithError("Ошибка генерации .yml: %v", err)
			}
			fmt.Printf("✅ Файл .yml успешно создан в %s\n", ymlPath)
		} else {
			envPath := filepath.Join(currentDir, ".env")
			if err := generateEnvFile(envPath, configFields); err != nil {
				exitWithError("Ошибка генерации .env: %v", err)
			}
			fmt.Printf("✅ Файл .env успешно создан в %s\n", envPath)
		}

		//Дебаг message
		fmt.Println("Добавленные переменные:")
		for _, field := range configFields {
			fmt.Printf("- %s\n", field)
		}
	},
}

func init() {
	genCmd.Flags().BoolP("yml", "y", false, "для генерации файла env.yml")
}

// Парсим структуру Config и возвращаем список полей
func parseConfigStruct(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга файла: %v", err)
	}

	var fields []string
	ast.Inspect(file, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok || ts.Name.Name != "Config" {
			return true
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range st.Fields.List {
			for _, name := range field.Names {
				fields = append(fields, name.Name)
			}
		}
		return false
	})

	if len(fields) == 0 {
		return nil, fmt.Errorf("структура Config не найдена или не содержит полей")
	}
	return fields, nil
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

// Вспомогательные функции
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func exitWithError(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}
