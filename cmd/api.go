package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kymmt90/colorme-cli/pkg/api"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Call APIs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if GetAccessToken() == "" {
			return errors.New("log in or set COLORME_ACCESS_TOKEN")
		}

		if len(args) == 0 {
			return errors.New("specify a path")
		}

		client, err := api.NewClient(apiBaseURL, GetAccessToken())
		if err != nil {
			return err
		}

		path := args[0]
		resp, err := client.Get(path, "")
		if err != nil {
			return err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var out bytes.Buffer
		if err := json.Indent(&out, body, "", "  "); err != nil {
			return err
		}
		if _, err := out.WriteTo(os.Stdout); err != nil {
			return err
		}
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
