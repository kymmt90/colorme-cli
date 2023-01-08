package main

import (
	"fmt"
	"log"
	"net/url"
)

func main() {
	fmt.Println(AuthorizationUrl())
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
