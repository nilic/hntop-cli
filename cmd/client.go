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
		UserAgent:  "hntop/" + getVersion(),
		httpClient: http.DefaultClient,
	}
	return c
}

func (c *Client) NewRequest(path string) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)

	if err != nil {
		return nil, fmt.Errorf("error reading response from %s %s: %s", req.Method, req.URL.RequestURI(), err)
	}

	return resp, nil
}
