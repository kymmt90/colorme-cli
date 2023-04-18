package cmd

import (
	"fmt"
	"os"

	"github.com/kymmt90/colorme-cli/pkg/api"
	"github.com/kymmt90/colorme-cli/pkg/order"
	"github.com/spf13/cobra"
)

var orderCmd = &cobra.Command{
	Use:   "order",
	Short: "Manage orders",
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

		if err := order.ListOrders(client); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(orderCmd)
}
