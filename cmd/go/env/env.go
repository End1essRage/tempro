package cmd_go

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var EnvCmd = &cobra.Command{
	Use:   "env",
	Short: "A brief description of your command",
	Long:  "",
}

func init() {
	EnvCmd.AddCommand(genCmd)
	EnvCmd.AddCommand(addCmd)
}

// получаем все переменные из Config в main.go
func getExisting() ([]string, error) {
	//получаем текущую папку
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения текущей директории: %v", err)
	}

	//ищем файл main.go
	mainGoPath := filepath.Join(currentDir, "main.go")
	if !fileExists(mainGoPath) {
		return nil, fmt.Errorf("Файл main.go не найден в %s", currentDir)
	}

	//парсим
	configFields, err := parseConfigStruct(mainGoPath)
	if err != nil {
		return nil, err
	}

	return configFields, nil
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

// Вспомогательные функции
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
func exitWithError(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}
