package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/kymmt90/colorme-cli/auth"
	"github.com/kymmt90/colorme-cli/cmd"
)

var jsonOutputFlag bool

type Order struct {
	ID         int  `json:"id"`
	TotalPrice int  `json:"total_price"`
	Paid       bool `json:"paid"`
	Delivered  bool `json:"delivered"`
	Customer   struct {
		Name           string `json:"name"`
		PrefectureName string `json:"pref_name"`
		Address1       string `json:"address1"`
		Address2       string `json:"address2"`
	} `json:"customer"`
}

const orderTemplate = `=== Order %d
+Total Price: ¥%d
+Customer name: %s
+Customer address: %s`

var orderFlagSet = flag.NewFlagSet("order", flag.ExitOnError)

func init() {
	orderFlagSet.BoolVar(&jsonOutputFlag, "json", false, "output as JSON")
}

func main() {
	if len(os.Args[1:]) == 0 {
		fmt.Fprintf(os.Stderr, "$ colorme login\n")
		os.Exit(1)
	}

	command := os.Args[1]
	if command == "order" {
		orderFlagSet.Parse(os.Args[2:])
		accessToken := auth.GetAccessTokenFromEnv()
		if accessToken == nil {
			fmt.Fprintln(os.Stderr, "Set COLORME_ACCESS_TOKEN")
			os.Exit(1)
		}
		GetOrders(*accessToken, jsonOutputFlag)
	} else if command == "login" || command == "product" {
		cmd.Execute()
	} else {
		fmt.Fprintf(os.Stderr, "$ colorme login\n")
		os.Exit(1)
	}
}

func GetOrders(accessToken string, outputAsJson bool) {
	req, err := http.NewRequest("GET", "https://api.shop-pro.jp/v1/sales?after=2022-01-01", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

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

	ordersJson, err := json.Marshal(payload["sales"])
	if err != nil {
		log.Fatal(err)
	}

	var orders []Order
	if err := json.Unmarshal(ordersJson, &orders); err != nil {
		log.Fatal(err)
	}

	for _, v := range orders {
		fmt.Printf(orderTemplate+"\n", v.ID, v.TotalPrice, v.Customer.Name, v.Customer.PrefectureName+v.Customer.Address1+v.Customer.Address2)
	}
}
