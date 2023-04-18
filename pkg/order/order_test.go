package order

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

func TestListOrders(t *testing.T) {
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

		case "/sales":
			f, err := os.Open("fixtures/orders.json")
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			orders, err := io.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(orders)
		}
	}))

	c := &api.Client{
		Client:      &http.Client{},
		AccessToken: accessToken,
		BaseURL:     ts.URL,
	}

	err := ListOrders(c)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_parseOrders(t *testing.T) {
	f, err := os.Open("fixtures/orders.json")
	if err != nil {
		t.Fatal(err)
	}

	orders, err := parseOrders(f)
	if err != nil {
		t.Errorf("parseOrders returns error %v; want nil", err)
	}

	if got, want := len(orders.Orders), 1; got != want {
		t.Errorf("len(orders.Orders) = %d; want %d", got, want)
	}

	want := Order{
		ID:         4434233,
		Paid:       true,
		Delivered:  true,
		Canceled:   false,
		TotalPrice: 1930,
	}
	if got := orders.Orders[0]; !reflect.DeepEqual(got, want) {
		t.Errorf("orders[0] = %v; want %v", got, want)
	}
}
