package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	clientId        = "9dc453241d4fba503b235912fab6b9c3a90dc9eae88006affcc9ccf515621432"
	clientSecret    = "f73a04e4ea904aaf1c0282263073ea06d7a1b6d64b751eadfa4d196c69530048"
	redirectUri     = "urn:ietf:wg:oauth:2.0:oob"
	productTemplate = `=== Product %d
Name: %s
Stocks: %d
Model Number: %s
Price: ¥%d
Description: %s`
)

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ModelNumber string `json:"model_number"`
	Price       int    `json:"sales_price"`
	Description string `json:"expl"`
	Stocks      int    `json:"stocks"`
}

func main() {
	if len(os.Args[1:]) == 0 {
		fmt.Fprintf(os.Stderr, "$ colorme login\n")
		os.Exit(1)
	}

	command := os.Args[1]
	if command == "login" {
		Login()
	} else if command == "product" {
		accessToken := getAccessTokenFromEnv()
		GetProducts(accessToken)
	} else {
		fmt.Fprintf(os.Stderr, "$ colorme login\n")
		os.Exit(1)
	}
}

func getAccessTokenFromEnv() string {
	accessToken, found := os.LookupEnv("COLORME_ACCESS_TOKEN")
	if !found {
		fmt.Fprintf(os.Stderr, "Set COLORME_ACCESS_TOKEN")
		os.Exit(1)
	}

	return accessToken
}

func GetProducts(accessToken string) {
	req, err := http.NewRequest("GET", "https://api.shop-pro.jp/v1/products?limit=1", nil)
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

func Login() {
	fmt.Println("Access to this URL and authorize this app")
	fmt.Println(AuthorizationUrl())
	fmt.Println()
	fmt.Println("Paste the \"Authorization Complete\" page's URL")
	fmt.Printf("URL: ")

	authorizationCompleteUrl, err := scanFromStdin()
	if err != nil {
		log.Fatal(err)
	}

	code := AuthorizationCode(authorizationCompleteUrl)
	log.Printf("Authorization Code: %s\n", code)

	tokenEndpointRawResponse := GetTokenEndpointRawResponse(code)
	log.Println(tokenEndpointRawResponse)

	fmt.Println("Login succeeded")
}

func AuthorizationUrl() string {
	url, err := url.Parse("https://api.shop-pro.jp/oauth/authorize")
	if err != nil {
		log.Fatal(err)
	}

	q := url.Query()
	q.Set("client_id", clientId)
	q.Set("response_type", "code")
	q.Set("scope", "read_products write_products read_sales write_sales read_shop_coupons")
	q.Set("redirect_uri", redirectUri)

	url.RawQuery = q.Encode()

	return url.String()
}

func scanFromStdin() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text(), nil
	} else if err := scanner.Err(); err != nil {
		return "", err
	} else {
		return "", nil
	}
}

func AuthorizationCode(authorizationCompleteUrl string) string {
	url, err := url.Parse(authorizationCompleteUrl)
	if err != nil {
		log.Fatal(err)
	}

	splitted := strings.Split(url.Path, "/")

	return splitted[len(splitted)-1]
}

func GetTokenEndpointRawResponse(code string) string {
	v := url.Values{}
	v.Set("client_id", clientId)
	v.Set("client_secret", clientSecret)
	v.Set("code", code)
	v.Set("grant_type", "authorization_code")
	v.Set("redirect_uri", redirectUri)

	resp, err := http.PostForm("https://api.shop-pro.jp/oauth/token", v)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}
