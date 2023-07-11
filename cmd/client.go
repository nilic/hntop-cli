package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	apiBaseURL = "https://hn.algolia.com/api/v1/"
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
}

func NewClient() *Client {
	baseURL, _ := url.Parse(apiBaseURL)
	c := &Client{
		BaseURL:    baseURL,
		UserAgent:  appName + "/" + getVersion(),
		httpClient: http.DefaultClient,
	}
	return c
}

func (c *Client) NewRequest(path string) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parsing API base URL: %w", err)
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

func (c *Client) Do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("invoking HN API: %w", err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)

	if err != nil {
		return fmt.Errorf("reading response from %s %s: %s", req.Method, req.URL.RequestURI(), err)
	}

	return nil
}
