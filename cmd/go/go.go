/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd_go

import (
	"fmt"

	"github.com/end1essrage/tempro/generator"
	"github.com/spf13/cobra"
)

var (
	Name string
)

// goCmd represents the go command
var GoCmd = &cobra.Command{
	Use:   "go",
	Short: "A brief description of your command",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("go called")

		fmt.Println("Name")

		// checks mod
		// runs go mod init

		if err := generator.GenerateFiles(generator.GolangSimple, generator.ProjectConfig{ModuleName: Name}); err != nil {
			return fmt.Errorf("error generating file : " + err.Error())
		}

		return nil
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// goCmd.PersistentFlags().String("foo", "", "A help for foo")

	GoCmd.PersistentFlags().StringVar(&Name, "name", "", "mod name")
	GoCmd.MarkPersistentFlagRequired("name")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// goCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
