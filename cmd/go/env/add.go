package cmd_go_env

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Добавить переменную с типом string",
	Long:  "",
	PreRun: func(cmd *cobra.Command, args []string) {
		// Проверяем наличие флагов
		fileFlag, _ := cmd.Flags().GetBool("file")
		namesFlag, _ := cmd.Flags().GetStringSlice("names")

		if !fileFlag && len(namesFlag) == 0 {
			cmd.Help()
			fmt.Printf("\nP.S. невозможно создать пустой список переменных))")
			os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		//получаем список переменных для создания
		toCreate, err := getVariablesToCreate(cmd)
		if err != nil {
			exitWithError("ошибка при формировании списка переменных %v", err)
		}

		currentDir, err := os.Getwd()
		if err != nil {
			exitWithError("Ошибка получения текущей директории: %v", err)
		}

		//ищем файл main.go
		mainGoPath := filepath.Join(currentDir, "main.go")
		if !fileExists(mainGoPath) {
			exitWithError("Файл main.go не найден в %s", currentDir)
		}

		if err := updateCfg(mainGoPath, toCreate); err != nil {
			exitWithError("Ошибка обновления списка переменных %v", err)
		}

		fmt.Print("Список переменных успешно обновлен")
	},
}

func init() {
	addCmd.Flags().BoolP("file", "f", false, "для генерации из файла .env в корне приложения")
	addCmd.Flags().StringSliceP("names", "n", nil, "перечислите названия переменных, разделенные запятой")
}

func updateCfg(filePath string, varNames []string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var configFound bool
	ast.Inspect(file, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok || ts.Name.Name != "Config" {
			return true
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		for _, fieldName := range varNames {
			// Добавление нового поля
			newField := &ast.Field{
				Names: []*ast.Ident{ast.NewIdent(fieldName)},
				Type:  ast.NewIdent("string"),
			}
			st.Fields.List = append(st.Fields.List, newField)
		}

		configFound = true
		return false
	})

	if !configFound {
		return fmt.Errorf("Config struct not found")
	}

	ast.Inspect(file, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "loadCfg" {
			return true
		}

		newStmts := make([]*ast.AssignStmt, 0)
		// Добавление новой строки
		for _, varName := range varNames {
			newStmt := &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.SelectorExpr{
					X:   ast.NewIdent("res"),
					Sel: ast.NewIdent(varName),
				}},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("os"),
							Sel: ast.NewIdent("Getenv"),
						},
						Args: []ast.Expr{&ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"" + varName + "\"",
						}},
					},
				},
			}

			newStmts = append(newStmts, newStmt)
		}

		// Вставка перед return
		for _, newStmt := range newStmts {
			for i, stmt := range fn.Body.List {
				if ret, ok := stmt.(*ast.ReturnStmt); ok && len(ret.Results) > 0 {
					fn.Body.List = append(fn.Body.List[:i], append([]ast.Stmt{newStmt}, fn.Body.List[i:]...)...)
					break
				}
			}
		}

		return false
	})

	return saveModifiedFile(fset, file, filePath)
}

func saveModifiedFile(fset *token.FileSet, file *ast.File, path string) error {
	output, err := os.Create(path)
	if err != nil {
		return err
	}
	defer output.Close()

	return format.Node(output, fset, file)
}

func parseEnvFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(".env file not found")
	}
	defer file.Close()

	var variables []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) > 0 && parts[0] != "" {
			variables = append(variables, parts[0])
		}
	}

	if len(variables) == 0 {
		return nil, fmt.Errorf("no variables found in .env")
	}

	return variables, nil
}

func getVariablesToCreate(cmd *cobra.Command) ([]string, error) {
	toCreate := make([]string, 0)
	//получаем набор из флага
	names, err := cmd.Flags().GetStringSlice("names")
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении флага names: %v", err)
	}
	if len(names) > 0 {
		toCreate = append(toCreate, names...)
		//debug
		fmt.Print("from names: " + strings.Join(names, ", "))
	}

	//получаем набор из .env файла при наличии флага
	fromFile, err := cmd.Flags().GetBool("file")
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении флага file: %v", err)
	}
	if fromFile {
		//получаем текущую папку
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("ошибка получения текущей директории: %v", err)
		}

		//ищем файл .env
		envPath := filepath.Join(currentDir, ".env")
		if !fileExists(envPath) {
			return nil, fmt.Errorf("Файл .env не найден в %s", currentDir)
		}

		envVars, err := parseEnvFile(envPath)
		if err != nil {
			return nil, fmt.Errorf("Ошибка считывания .env файла: %v", err)
		}
		//debug
		fmt.Print("from .env: " + strings.Join(envVars, ", "))

		toCreate = append(toCreate, envVars...)
	}
	//вытянуть существующие
	existing, err := getExisting()
	if err != nil {
		return nil, fmt.Errorf("Ошибка считывания из main.go файла %v", err)
	}

	//оставить только уникальные через мапу
	result := make(map[string]bool)
	for _, v := range toCreate {
		result[v] = true
	}

	//пометить существующие
	for _, v := range existing {
		result[v] = false
	}

	keys := make([]string, 0, len(result))
	for k, v := range result {
		// пропускает только новые
		if v {
			keys = append(keys, k)
		}
	}

	return keys, nil
}
