package cmd

import (
	"fmt"
	"os"

	"github.com/kymmt90/colorme-cli/pkg/auth"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "colorme",
	Short: "Colormeshop CLI",
	Run: func(cmd *cobra.Command, args []string) {
		println("run cobra")
	},
}

var accessToken *string

func init() {
	cobra.OnInitialize(initConfig)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	accessToken = auth.GetAccessTokenFromEnv()
}
