package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args[1:]) == 0 {
		fmt.Fprintf(os.Stderr, "$ colorme login\n")
		os.Exit(1)
	}

	command := os.Args[1]
	if command == "login" {
		Login()
	} else {
		fmt.Fprintf(os.Stderr, "$ colorme login\n")
		os.Exit(1)
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

	log.Printf("Authorization Code: %s\n", AuthorizationCode(authorizationCompleteUrl))
}

func AuthorizationUrl() string {
	url, err := url.Parse("https://api.shop-pro.jp/oauth/authorize")
	if err != nil {
		log.Fatal(err)
	}

	q := url.Query()
	q.Set("client_id", "9dc453241d4fba503b235912fab6b9c3a90dc9eae88006affcc9ccf515621432")
	q.Set("response_type", "code")
	q.Set("scope", "read_products write_products read_sales write_sales read_shop_coupons")
	q.Set("redirect_uri", "urn:ietf:wg:oauth:2.0:oob")

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
