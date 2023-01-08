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
