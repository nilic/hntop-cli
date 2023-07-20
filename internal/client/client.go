package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
}

func NewClient(baseURL, userAgent string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing API base URL: %w", err)
	}

	c := &Client{
		BaseURL:    u,
		UserAgent:  userAgent,
		httpClient: http.DefaultClient,
	}

	return c, nil
}

func (c *Client) NewRequest(path string) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parsing API relative path: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating API request: %w", err)
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("invoking API: %w", err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)

	if err != nil {
		return fmt.Errorf("reading response from %s %s: %s", req.Method, req.URL.RequestURI(), err)
	}

	return nil
}
