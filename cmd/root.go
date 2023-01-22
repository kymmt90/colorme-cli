package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "colorme",
	Short: "Colormeshop CLI",
	Run: func(cmd *cobra.Command, args []string) {
		println("run cobra")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
