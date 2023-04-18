package order

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/kymmt90/colorme-cli/pkg/api"
	"github.com/kymmt90/colorme-cli/pkg/shop"
)

type OrdersResource struct {
	Orders []Order `json:"sales"`
}

type Order struct {
	ID         int  `json:"id"`
	Paid       bool `json:"paid"`
	Delivered  bool `json:"delivered"`
	Canceled   bool `json:"canceled"`
	TotalPrice int  `json:"total_price"`
}

const orderTemplate = `=== Order %d
Total Price: ¥%d
Paid: %v
Delivered: %v
Canceled: %v`

func ListOrders(client *api.Client) error {
	resShop, err := client.FetchShop()
	if err != nil {
		return fmt.Errorf("ListOrders: %w", err)
	}
	defer resShop.Close()

	var shop shop.ShopResource
	err = json.NewDecoder(resShop).Decode(&shop)
	if err != nil {
		return fmt.Errorf("ListOrders: %w", err)
	}

	resOrders, err := client.FetchOrders()
	if err != nil {
		return fmt.Errorf("ListOrders: %w", err)
	}
	defer resOrders.Close()

	orders, err := parseOrders(resOrders)
	if err != nil {
		return fmt.Errorf("ListOrders: %w", err)
	}

	for _, v := range orders.Orders {
		fmt.Printf(orderTemplate+"\n", v.ID, v.TotalPrice, v.Paid, v.Delivered, v.Canceled)
		fmt.Printf("View this order on the admin: https://admin.shop-pro.jp/?mode=shopcust_sales&sales_id=%d\n\n", v.ID)
	}

	return nil
}

func parseOrders(rawResp io.Reader) (OrdersResource, error) {
	var orders OrdersResource
	err := json.NewDecoder(rawResp).Decode(&orders)
	if err != nil {
		return OrdersResource{}, fmt.Errorf("parseOrders: %w", err)
	}

	return orders, nil
}
