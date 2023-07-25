package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	URL        *url.URL
	UserAgent  string
	httpClient *http.Client
}

func NewClient(URL, userAgent string) (*Client, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, fmt.Errorf("parsing URL %s: %w", URL, err)
	}

	c := &Client{
		URL:        u,
		UserAgent:  userAgent,
		httpClient: http.DefaultClient,
	}

	return c, nil
}

func (c *Client) NewRequest() (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, c.URL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating API request: %w", err)
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v any) error {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("invoking API: %w", err)
	}

	if res == nil {
		return fmt.Errorf("empty response from %s", req.URL.RequestURI())
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading API response: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("calling %s:\nstatus: %s\nresponseData: %s", req.URL.RequestURI(), res.Status, body)
	}

	err = json.Unmarshal(body, v)

	if err != nil {
		return fmt.Errorf("reading response from %s %s: %s", req.Method, req.URL.RequestURI(), err)
	}

	return nil
}
