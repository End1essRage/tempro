/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd_go

import (
	e "github.com/end1essrage/tempro/cmd/go/env"
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
	//Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	GoCmd.AddCommand(e.EnvCmd)
}
