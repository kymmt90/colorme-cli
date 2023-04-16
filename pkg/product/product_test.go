package product_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kymmt90/colorme-cli/pkg/api"
	"github.com/kymmt90/colorme-cli/pkg/product"
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

	err := product.ListProducts(c)
	if err != nil {
		t.Errorf("%v", err)
	}
}
