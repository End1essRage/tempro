package generator

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"
)

type ProjectConfig struct {
	ModuleName string
}

const (
	GolangSimple = "templates/golang/simple"
)

func GenerateFiles(template string, config ProjectConfig) error {
	path := ""

	switch template {
	case "simple":
	case "":
	default:
		path = GolangSimple
	}

	return generateAll(path, &config)
}

func generateAll(tmplDirPath string, config *ProjectConfig) error {
	info, err := os.Stat(tmplDirPath)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	if !info.IsDir() {
		return generateFile(tmplDirPath, tmplDirPath, config)
	}

	templates, err := os.ReadDir(tmplDirPath)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	for _, tmpl := range templates {
		generateFile(path.Join(tmplDirPath, tmpl.Name()), tmpl.Name(), config)
	}

	return nil
}

func generateFile(tmplPath string, tmplName string, config *ProjectConfig) error {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	//for debug
	os.Mkdir("output", os.ModeDir)
	outputFile, err := os.Create("output/" + wipeTmplExt(tmplName))
	//for debug
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	defer outputFile.Close()

	return tmpl.Execute(outputFile, config)
}

func wipeTmplExt(path string) string {
	stringShards := strings.Split(path, ".")
	countShard := len(stringShards)

	if stringShards[countShard-1] == "tmpl" {
		stringShards[countShard-1] = ""
	}

	return strings.Join(stringShards, ".")
}
