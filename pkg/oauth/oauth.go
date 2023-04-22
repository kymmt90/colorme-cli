package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

const (
	authorizationBaseURL = "https://api.shop-pro.jp"
	clientId             = "9dc453241d4fba503b235912fab6b9c3a90dc9eae88006affcc9ccf515621432"
	clientSecret         = "f73a04e4ea904aaf1c0282263073ea06d7a1b6d64b751eadfa4d196c69530048"
	redirectUri          = "http://127.0.0.1:8080/callback"
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
	u, err := buildAuthorizationUrl()
	if err != nil {
		return nil, err
	}

	if err := openInBrowser(u); err != nil {
		return nil, err
	}

	token, err := acceptCallbackFromAuthorizationServer()
	if err != nil {
		return nil, err
	}

	return token, nil
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

func buildAuthorizationUrl() (string, error) {
	u, err := url.Parse(authorizationBaseURL + "/oauth/authorize")
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

func acceptCallbackFromAuthorizationServer() (*TokenEndpointResponse, error) {
	var resp io.Reader
	gotTokenResponse := make(chan struct{})
	failed := make(chan error)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code, err := scanAuthorizationCodeFromCallback(r.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			failed <- err
			return
		}

		c, err := newClient(authorizationBaseURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			failed <- err
			return
		}

		resp, err = c.fetchAccessToken(code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			failed <- err
			return
		}

		w.Write([]byte("Authorization completed. You can close this page and return to your CLI."))

		gotTokenResponse <- struct{}{}
	})

	s := http.Server{
		Addr:         "127.0.0.1:8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  5 * time.Minute,
		Handler:      mux,
	}

	var tokenError error
	go func() {
		select {
		case <-gotTokenResponse:
			if err := s.Shutdown(context.Background()); err != nil {
				log.Println(err)
			}
		case err := <-failed:
			tokenError = err
		}

		close(gotTokenResponse)
		close(failed)

		if err := s.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
	}()

	if err := s.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return nil, err
		}
	}
	defer s.Close()

	if tokenError != nil {
		return nil, tokenError
	}

	var token TokenEndpointResponse
	if err := json.NewDecoder(resp).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func scanAuthorizationCodeFromCallback(url *url.URL) (string, error) {
	code := url.Query().Get("code")

	if code == "" {
		return "", fmt.Errorf("scanAuthorizationCodeFromCallback: code is empty")
	}

	return code, nil
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
		return nil, fmt.Errorf("fetchAccessToken: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetchAccessToken: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetchAccessToken: status code is %d", resp.StatusCode)
	}

	return resp.Body, nil
}
