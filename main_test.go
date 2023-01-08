package main

import (
	"regexp"
	"testing"
)

func TestAuthorizationUrl(t *testing.T) {
	url := AuthorizationUrl()
	want := regexp.MustCompile(`^https://api\.shop-pro.jp/oauth/authorize\?.+$`)

	if !want.MatchString(url) {
		t.Fatalf(`AuthorizationUrl() = %q, want match for %#q`, url, want)
	}
}

func TestAuthorizationCode(t *testing.T) {
	authorizationCompleteUrl := "https://api.shop-pro.jp/oauth/authorize/b8370b810a4efea7bcbc3b71c90de261a25ee535f1df023ed448319cf4bd2e3f"
	want := "b8370b810a4efea7bcbc3b71c90de261a25ee535f1df023ed448319cf4bd2e3f"

	if actual := AuthorizationCode(authorizationCompleteUrl); actual != want {
		t.Fatalf(`AuthorizationCode(%q) = %q, want equal to %#q`, authorizationCompleteUrl, actual, want)
	}
}
