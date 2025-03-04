package cmd_go_env

import "github.com/spf13/cobra"

var EnvCmd = &cobra.Command{
	Use:   "env",
	Short: "A brief description of your command",
	Long:  "",
}

func init() {
	EnvCmd.AddCommand(genCmd)
}
