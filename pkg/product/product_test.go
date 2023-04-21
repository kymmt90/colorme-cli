package product

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/kymmt90/colorme-cli/pkg/api"
)

func TestListProducts(t *testing.T) {
	accessToken := "access-token"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get("Authorization"), fmt.Sprintf("Bearer %s", accessToken); got != want {
			t.Errorf("accessToken = %q; want %q", got, want)
		}

		switch r.URL.Path {
		case "/shop":
			f, err := os.Open("fixtures/shop.json")
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			shop, err := io.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(shop)

		case "/products":
			f, err := os.Open("fixtures/products.json")
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			products, err := io.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(products)
		}
	}))
	defer ts.Close()

	c := &api.Client{
		Client:      &http.Client{},
		AccessToken: accessToken,
		BaseURL:     ts.URL,
	}

	err := ListProducts(c)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func Test_parseProducts(t *testing.T) {
	f, err := os.Open("fixtures/products.json")
	if err != nil {
		t.Fatal(err)
	}

	products, err := parseProducts(f)
	if err != nil {
		t.Errorf("parseProducts returns error %v; want nil", err)
	}

	if got, want := len(products.Products), 1; got != want {
		t.Errorf("len(products.Products) = %d; want %d", got, want)
	}

	want := Product{
		ID:          1342332,
		Name:        "Tシャツ",
		ModelNumber: "T-223",
		Price:       1980,
		Description: "綿100%のこだわりTシャツです。\n\n肌触りや吸水性の良さにみなさま驚かれます。弊社の人気商品です。\n",
		Stocks:      20,
	}
	if got := products.Products[0]; !reflect.DeepEqual(got, want) {
		t.Errorf("products[0] = %v; want %v", got, want)
	}
}
