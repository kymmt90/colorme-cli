package product

import (
	"encoding/json"
	"fmt"

	"github.com/kymmt90/colorme-cli/pkg/api"
	"github.com/kymmt90/colorme-cli/pkg/shop"
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

func ListProducts(client *api.Client) error {
	resShop, err := client.FetchShop()
	if err != nil {
		return fmt.Errorf("GetProducts: %w", err)
	}
	defer resShop.Close()

	var shop shop.ShopResource
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
