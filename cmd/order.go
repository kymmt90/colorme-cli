package cmd

import (
	"errors"

	"github.com/kymmt90/colorme-cli/pkg/api"
	"github.com/kymmt90/colorme-cli/pkg/order"
	"github.com/spf13/cobra"
)

var orderCmd = &cobra.Command{
	Use:   "order",
	Short: "Manage orders",
	RunE: func(cmd *cobra.Command, args []string) error {
		if GetAccessToken() == "" {
			return errors.New("log in or set COLORME_ACCESS_TOKEN")
		}

		client, err := api.NewClient(apiBaseURL, GetAccessToken())
		if err != nil {
			return err
		}

		if err := order.ListOrders(client); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(orderCmd)
}
