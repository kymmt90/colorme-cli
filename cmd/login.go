package cmd

import (
	"github.com/kymmt90/colorme-cli/pkg/login"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Colorme Shop",
	Run: func(cmd *cobra.Command, args []string) {
		login.Login()
	},
}
