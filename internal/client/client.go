package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type Client struct {
	URL        *url.URL
	UserAgent  string
	httpClient *http.Client
}

func NewClient(URL, userAgent string) (*Client, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, fmt.Errorf("parsing URL %q: %w", URL, err)
	}

	c := &Client{
		URL:       u,
		UserAgent: userAgent,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	return c, nil
}

func (c *Client) NewRequest(httpMethod string, headers map[string]string, body io.Reader) (*http.Request, error) {
	regex := regexp.MustCompile(`^(GET|POST|PUT|PATCH|DELETE)$`)
	if !regex.MatchString(httpMethod) {
		return nil, fmt.Errorf("invalid HTTP method: %q", httpMethod)
	}

	req, err := http.NewRequest(httpMethod, c.URL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("creating API request: %w", err)
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v any) error {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("invoking API: %w", err)
	}

	if res == nil {
		return fmt.Errorf("empty response from %q", req.URL.RequestURI())
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading API response: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("calling %q:\nstatus: %q\nresponseData: %q", req.URL.RequestURI(), res.Status, responseBody)
	}

	err = json.Unmarshal(responseBody, v)

	if err != nil {
		return fmt.Errorf(`reading response from "%s %s": %w`, req.Method, req.URL.RequestURI(), err)
	}

	return nil
}

func MakeHTTPRequest[T any](httpMethod, URL, userAgent string, headers map[string]string, body io.Reader, responseType T) (T, error) {
	c, err := NewClient(URL, userAgent)
	if err != nil {
		return responseType, fmt.Errorf("creating API client: %w", err)
	}

	req, err := c.NewRequest(httpMethod, headers, body)
	if err != nil {
		return responseType, fmt.Errorf("creating API request: %w", err)
	}

	var response T
	err = c.Do(req, &response)
	if err != nil {
		return responseType, fmt.Errorf("calling API: %w", err)
	}

	return response, nil
}
