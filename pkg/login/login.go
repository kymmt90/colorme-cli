package login

import (
	"fmt"

	"github.com/kymmt90/colorme-cli/pkg/api"
	"github.com/kymmt90/colorme-cli/pkg/config"
	"github.com/kymmt90/colorme-cli/pkg/oauth"
	"github.com/kymmt90/colorme-cli/pkg/shop"
)

func Login() error {
	fmt.Println("Open the authorization URL...")

	token, err := oauth.DoAuthorizationCodeFlow()
	if err != nil {
		return fmt.Errorf("Login: %w", err)
	}

	client, err := api.NewClient("https://api.shop-pro.jp/v1", token.AccessToken)
	if err != nil {
		return fmt.Errorf("Login: %w", err)
	}

	resShop, err := client.FetchShop()
	if err != nil {
		return fmt.Errorf("Login: %w", err)
	}

	shopResource, err := shop.Deserialize(resShop)
	if err != nil {
		return fmt.Errorf("Login: %w", err)
	}

	cfg := &config.UserConfig{
		LoginID:     shopResource.Shop.LoginID,
		AccessToken: token.AccessToken,
	}
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("Login: %w", err)
	}

	fmt.Println("Login succeeded")

	return nil
}
