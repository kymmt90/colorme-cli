package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
			j, err := cmd.Flags().GetBool("json")
			if err != nil {
				log.Fatal(err)
			}

			GetProducts(j)
		},
	}
	outputAsJson bool
)

func init() {
	rootCmd.AddCommand(productCmd)
	productCmd.Flags().BoolVarP(&outputAsJson, "json", "j", false, "output as JSON")
}

func getShop() (*ShopResource, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", "https://api.shop-pro.jp/v1/shop", nil)
	if err != nil {
		return nil, fmt.Errorf("getShop: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+*accessToken)
	client := http.Client{
		Timeout: 30 * time.Second,
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

func GetProducts(outputAsJson bool) {
	if accessToken == nil {
		fmt.Fprintln(os.Stderr, "Set COLORME_ACCESS_TOKEN")
		os.Exit(1)
	}

	req, err := http.NewRequest("GET", "https://api.shop-pro.jp/v1/products?limit=1", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+*accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if outputAsJson {
		fmt.Println(string(body))
		return
	}

	shop, err := getShop()
	if err != nil {
		log.Fatal(err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Fatal(err)
	}

	productsJson, err := json.Marshal(payload["products"])
	if err != nil {
		log.Fatal(err)
	}

	var products []Product
	if err := json.Unmarshal(productsJson, &products); err != nil {
		log.Fatal(err)
	}

	for i, v := range products {
		fmt.Printf(productTemplate+"\n", i+1, v.Name, v.Stocks, v.ModelNumber, v.Price, v.Description)
		fmt.Printf("View this product on the shop: %s/?pid=%d\n", shop.Shop.URL, v.ID)
	}
}
