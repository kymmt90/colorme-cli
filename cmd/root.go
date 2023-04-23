package cmd

import (
	"fmt"
	"os"

	"github.com/kymmt90/colorme-cli/pkg/config"
	"github.com/spf13/cobra"
)

var (
	apiBaseURL = "https://api.shop-pro.jp/v1"
	userConfig *config.UserConfig

	rootCmd = &cobra.Command{
		Use:   "colorme",
		Short: "Color Me Shop CLI",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	var err error
	userConfig, err = config.LoadUserConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func GetAccessToken() string {
	token, ok := os.LookupEnv("COLORME_ACCESS_TOKEN")
	if ok {
		return token
	}

	return userConfig.AccessToken
}
