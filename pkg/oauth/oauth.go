package oauth

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	clientId     = "9dc453241d4fba503b235912fab6b9c3a90dc9eae88006affcc9ccf515621432"
	clientSecret = "f73a04e4ea904aaf1c0282263073ea06d7a1b6d64b751eadfa4d196c69530048"
	redirectUri  = "urn:ietf:wg:oauth:2.0:oob"
)

type Client struct {
	*http.Client
	BaseURL string
}

type TokenEndpointResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	CreatedAt   int    `json:"created_at"`
}

func DoAuthorizationCodeFlow() (*TokenEndpointResponse, error) {
	c, err := newClient("https://api.shop-pro.jp")
	if err != nil {
		return nil, err
	}

	u, err := c.buildAuthorizationUrl()
	if err != nil {
		return nil, err
	}

	if err := openInBrowser(u); err != nil {
		return nil, err
	}

	fmt.Println("\nPaste the \"Authorization Complete\" page's URL")
	fmt.Printf("URL: ")

	code, err := scanAuthorizationCode()
	if err != nil {
		return nil, err
	}

	resToken, err := c.fetchAccessToken(code)
	if err != nil {
		return nil, err
	}

	var token TokenEndpointResponse
	if err = json.NewDecoder(resToken).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func newClient(baseURL string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}

	return &Client{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		BaseURL: u.String(),
	}, nil
}

func (c *Client) buildAuthorizationUrl() (string, error) {
	u, err := url.Parse(c.BaseURL + "/oauth/authorize")
	if err != nil {
		return "", fmt.Errorf("BuildAuthorizationUrl: %w", err)
	}

	q := u.Query()
	q.Set("client_id", clientId)
	q.Set("response_type", "code")
	q.Set("scope", "read_products write_products read_sales write_sales read_shop_coupons")
	q.Set("redirect_uri", redirectUri)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func openInBrowser(url string) error {
	args := []string{"open", url}
	cmd := exec.Command(args[0], args[1])
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("openInBrowser: %w", err)
	}

	return nil
}

func scanAuthorizationCode() (string, error) {
	s := bufio.NewScanner(os.Stdin)
	if s.Scan() {
		u, err := url.Parse(s.Text())
		if err != nil {
			return "", fmt.Errorf("scanAuthorizationCode: %w", err)
		}

		segments := strings.Split(u.Path, "/")

		return segments[len(segments)-1], nil
	}

	if err := s.Err(); err != nil {
		return "", fmt.Errorf("scanAuthorizationCode: %w", err)
	}

	return "", errors.New("scanAuthorizationCode: failed to scan")
}

func (c *Client) fetchAccessToken(code string) (io.ReadCloser, error) {
	v := url.Values{}
	v.Set("client_id", clientId)
	v.Set("client_secret", clientSecret)
	v.Set("code", code)
	v.Set("grant_type", "authorization_code")
	v.Set("redirect_uri", redirectUri)

	req, err := http.NewRequestWithContext(context.Background(), "POST", c.BaseURL+"/oauth/token", strings.NewReader(v.Encode()))
	if err != nil {
		return nil, fmt.Errorf("FetchAccessToken: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("FetchAccessToken: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FetchAccessToken: status code is %d", resp.StatusCode)
	}

	return resp.Body, nil
}
