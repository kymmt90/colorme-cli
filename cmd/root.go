package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/kymmt90/colorme-cli/auth"
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
	accessToken = auth.GetAccessTokenFromEnv()
	if accessToken == nil {
		fmt.Fprintln(os.Stderr, "Set COLORME_ACCESS_TOKEN")
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
