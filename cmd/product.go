package cmd

import (
	"fmt"
	"os"

	"github.com/kymmt90/colorme-cli/pkg/api"
	"github.com/kymmt90/colorme-cli/pkg/product"
	"github.com/spf13/cobra"
)

var (
	productCmd = &cobra.Command{
		Use:   "product",
		Short: "Manage products",
		Run: func(cmd *cobra.Command, args []string) {
			if accessToken == nil {
				fmt.Fprintln(os.Stderr, "Set COLORME_ACCESS_TOKEN")
				os.Exit(1)
			}

			client, err := api.NewClient("https://api.shop-pro.jp/v1", *accessToken)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			err = product.ListProducts(client)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(productCmd)
}
