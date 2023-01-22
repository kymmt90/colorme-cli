package auth

import (
	"os"
)

func GetAccessTokenFromEnv() *string {
	accessToken, found := os.LookupEnv("COLORME_ACCESS_TOKEN")
	if !found {
		return nil
	}

	return &accessToken
}
