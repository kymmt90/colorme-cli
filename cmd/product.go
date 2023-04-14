package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/kymmt90/colorme-cli/pkg/api"
	"github.com/spf13/cobra"
)

const productTemplate = `=== Product %d
Name: %s
Stocks: %d
Model Number: %s
Price: ¥%d
Description: %s
`

type ProductsResource struct {
	Products []Product `json:"products"`
}

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ModelNumber string `json:"model_number"`
	Price       int    `json:"sales_price"`
	Description string `json:"expl"`
	Stocks      int    `json:"stocks"`
}

type ShopResource struct {
	Shop Shop `json:"shop"`
}

type Shop struct {
	URL string `json:"url"`
}

var (
	productCmd = &cobra.Command{
		Use:   "product",
		Short: "Manage products",
		Run: func(cmd *cobra.Command, args []string) {
			err := ListProducts()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(productCmd)
}

func ListProducts() error {
	if accessToken == nil {
		fmt.Fprintln(os.Stderr, "Set COLORME_ACCESS_TOKEN")
		os.Exit(1)
	}

	client, err := api.NewClient("https://api.shop-pro.jp/v1", *accessToken)
	if err != nil {
		return fmt.Errorf("GetProducts: %w", err)
	}

	resShop, err := client.FetchShop()
	if err != nil {
		return fmt.Errorf("GetProducts: %w", err)
	}
	defer resShop.Close()

	var shop ShopResource
	err = json.NewDecoder(resShop).Decode(&shop)
	if err != nil {
		return fmt.Errorf("GetProducts: %w", err)
	}

	resProducts, err := client.FetchProducts()
	if err != nil {
		return fmt.Errorf("GetProducts: %w", err)
	}
	defer resProducts.Close()

	var products ProductsResource
	err = json.NewDecoder(resProducts).Decode(&products)
	if err != nil {
		return fmt.Errorf("GetProducts: %w", err)
	}

	for i, v := range products.Products {
		fmt.Printf(productTemplate+"\n", i+1, v.Name, v.Stocks, v.ModelNumber, v.Price, v.Description)
		fmt.Printf("View this product on the shop: %s/?pid=%d\n\n", shop.Shop.URL, v.ID)
	}

	return nil
}
