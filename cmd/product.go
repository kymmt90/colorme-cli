package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

const productTemplate = `=== Product %d
Name: %s
Stocks: %d
Model Number: %s
Price: ¥%d
Description: %s`

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ModelNumber string `json:"model_number"`
	Price       int    `json:"sales_price"`
	Description string `json:"expl"`
	Stocks      int    `json:"stocks"`
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
	}
}
