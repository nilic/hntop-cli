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

type Options struct {
	UserAgent string
	Headers   map[string]string
	Body      io.Reader
}

func MakeHTTPRequest[T any](httpMethod, URL string, opts *Options, responseType T) (T, error) {
	c, err := newClient(URL, opts.UserAgent)
	if err != nil {
		return responseType, fmt.Errorf("creating API client: %w", err)
	}

	req, err := c.newRequest(httpMethod, opts.Headers, opts.Body)
	if err != nil {
		return responseType, fmt.Errorf("creating API request: %w", err)
	}

	var response T
	if err := c.do(req, &response); err != nil {
		return responseType, fmt.Errorf("calling API: %w", err)
	}

	return response, nil
}

type client struct {
	url        *url.URL
	userAgent  string
	httpClient *http.Client
}

func newClient(URL, userAgent string) (*client, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, fmt.Errorf("parsing URL %q: %w", URL, err)
	}

	c := &client{
		url:       u,
		userAgent: userAgent,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	return c, nil
}

func (c *client) newRequest(httpMethod string, headers map[string]string, body io.Reader) (*http.Request, error) {
	regex := regexp.MustCompile(`^(GET|POST|PUT|PATCH|DELETE)$`)
	if !regex.MatchString(httpMethod) {
		return nil, fmt.Errorf("invalid HTTP method: %q", httpMethod)
	}

	req, err := http.NewRequest(httpMethod, c.url.String(), body)
	if err != nil {
		return nil, fmt.Errorf("creating API request: %w", err)
	}

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (c *client) do(req *http.Request, v any) error {
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
		return fmt.Errorf(`calling "%s %s":\nstatus: %q\nresponse body: %q`, req.Method, req.URL.RequestURI(), res.Status, responseBody)
	}

	if err := json.Unmarshal(responseBody, v); err != nil {
		return fmt.Errorf(`reading response from "%s %s": %w`, req.Method, req.URL.RequestURI(), err)
	}

	return nil
}
