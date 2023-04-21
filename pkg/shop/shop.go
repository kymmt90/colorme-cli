package shop

import (
	"encoding/json"
	"fmt"
	"io"
)

type ShopResource struct {
	Shop Shop `json:"shop"`
}

type Shop struct {
	LoginID string `json:"login_id"`
	URL     string `json:"url"`
}

func Deserialize(r io.Reader) (*ShopResource, error) {
	var shop ShopResource
	err := json.NewDecoder(r).Decode(&shop)
	if err != nil {
		return nil, fmt.Errorf("Deserialize: %w", err)
	}

	return &shop, nil
}
