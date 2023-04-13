package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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
	client = &http.Client{
		Timeout: 30 * time.Second,
	}
	productFields = []string{"id", "name", "stocks", "model_number", "sales_price", "expl"}
	productCmd    = &cobra.Command{
		Use:   "product",
		Short: "Manage products",
		Run: func(cmd *cobra.Command, args []string) {
			err := GetProducts()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(productCmd)
}

func GetProducts() error {
	if accessToken == nil {
		fmt.Fprintln(os.Stderr, "Set COLORME_ACCESS_TOKEN")
		os.Exit(1)
	}

	u, err := buildProductsURL()
	if err != nil {
		return fmt.Errorf("GetProducts: %w", err)
	}

	req, err := buildGetRequest(u, *accessToken)
	if err != nil {
		return fmt.Errorf("GetProducts: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("GetProducts: %w", err)
	}
	defer res.Body.Close()

	shop, err := getShop()
	if err != nil {
		return fmt.Errorf("GetProducts: %w", err)
	}

	var products ProductsResource
	err = json.NewDecoder(res.Body).Decode(&products)
	if err != nil {
		return fmt.Errorf("GetProducts: %w", err)
	}

	for i, v := range products.Products {
		fmt.Printf(productTemplate+"\n", i+1, v.Name, v.Stocks, v.ModelNumber, v.Price, v.Description)
		fmt.Printf("View this product on the shop: %s/?pid=%d\n\n", shop.Shop.URL, v.ID)
	}

	return nil
}

func buildProductsURL() (string, error) {
	u, err := url.Parse("https://api.shop-pro.jp/v1/products?limit=30")
	if err != nil {
		return "", fmt.Errorf("buildProductsURL: %w", err)
	}

	q := u.Query()
	q.Set("fields", strings.Join(productFields, ","))
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func buildGetRequest(url string, token string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("buildGetRequest: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	return req, nil
}

func getShop() (*ShopResource, error) {
	req, err := buildGetRequest("https://api.shop-pro.jp/v1/shop", *accessToken)
	if err != nil {
		return nil, fmt.Errorf("getShop: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getShop: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getShop: status code is %d", res.StatusCode)
	}

	var shop ShopResource
	err = json.NewDecoder(res.Body).Decode(&shop)
	if err != nil {
		return nil, fmt.Errorf("getShop: %w", err)
	}

	return &shop, nil
}
