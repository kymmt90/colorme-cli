package api_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kymmt90/colorme-cli/pkg/api"
)

func TestFetchShop(t *testing.T) {
	accessToken := "access-token"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get("Authorization"), fmt.Sprintf("Bearer %s", accessToken); got != want {
			t.Errorf("accessToken = %q; want %q", got, want)
		}

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
	}))

	c := &api.Client{
		Client:      &http.Client{},
		AccessToken: accessToken,
		BaseURL:     ts.URL,
	}
	res, err := c.FetchShop()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer res.Close()
}

func TestFetchProducts(t *testing.T) {
	accessToken := "access-token"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get("Authorization"), fmt.Sprintf("Bearer %s", accessToken); got != want {
			t.Errorf("accessToken = %q; want %q", got, want)
		}

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
	}))

	c := &api.Client{
		Client:      &http.Client{},
		AccessToken: accessToken,
		BaseURL:     ts.URL,
	}
	res, err := c.FetchProducts()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer res.Close()
}
