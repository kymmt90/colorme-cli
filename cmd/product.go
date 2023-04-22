package cmd

import (
	"errors"

	"github.com/kymmt90/colorme-cli/pkg/api"
	"github.com/kymmt90/colorme-cli/pkg/product"
	"github.com/spf13/cobra"
)

var (
	productCmd = &cobra.Command{
		Use:   "product",
		Short: "Manage products",
		RunE: func(cmd *cobra.Command, args []string) error {
			if GetAccessToken() == "" {
				return errors.New("log in or set COLORME_ACCESS_TOKEN")
			}

			client, err := api.NewClient(apiBaseURL, GetAccessToken())
			if err != nil {
				return err
			}

			err = product.ListProducts(client)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(productCmd)
}
