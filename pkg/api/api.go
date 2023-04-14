package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	*http.Client
	AccessToken string
	BaseURL     string
}

var productFields = []string{"id", "name", "stocks", "model_number", "sales_price", "expl"}

func NewClient(baseURL string, accessToken string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}

	return &Client{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		AccessToken: accessToken,
		BaseURL:     u.String(),
	}, nil
}

func (c *Client) FetchShop() (io.ReadCloser, error) {
	res, err := c.get("/shop", "")
	if err != nil {
		return nil, fmt.Errorf("FetchShop: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FetchShop: status code is %d", res.StatusCode)
	}

	return res.Body, nil
}

func (c *Client) FetchProducts() (io.ReadCloser, error) {
	q := url.Values{}
	q.Set("fields", strings.Join(productFields, ","))
	q.Set("limit", "30")

	res, err := c.get("/products?", q.Encode())
	if err != nil {
		return nil, fmt.Errorf("FetchProducts: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FetchProducts: status code is %d", res.StatusCode)
	}

	return res.Body, nil
}

func (c *Client) get(path string, query string) (*http.Response, error) {
	u := c.BaseURL + path
	if query != "" {
		u += "?" + query
	}
	req, err := http.NewRequestWithContext(context.Background(), "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.AccessToken)

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return res, nil
}
