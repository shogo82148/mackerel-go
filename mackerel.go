package mackerel

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

var defaultBaseURL *url.URL

func init() {
	var err error
	defaultBaseURL, err = url.Parse("https://api.mackerelio.com/")
	if err != nil {
		panic(err)
	}
}

// Client is yet another client for mackerel.io
type Client struct {
	// BaseURL is a base url for the mackerel API.
	BaseURL *url.URL

	// APIKey is an API key for mackerelio.
	APIKey string

	// UserAgent is a value of User-Agent in a api request.
	// If it is empty, "shogo82148-mackerel-go" is used.
	UserAgent string

	// HTTPClient is a client of http.
	// If it is nil, http.DefaultClient is used.
	HTTPClient *http.Client

	// MaxRetries is a maximum count for retries.
	// If it is zero, th client continues retry forever.
	MaxRetries int
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}

func (c *Client) urlfor(path string) string {
	base := c.BaseURL
	if base == nil {
		base = defaultBaseURL
	}

	// shallow copy
	u := new(url.URL)
	*u = *base

	u.Path = path
	return u.String()
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	u := c.urlfor(path)
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("X-Api-Key", c.APIKey)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	} else {
		req.Header.Set("User-Agent", "shogo82148-mackerel-go")
	}

	return req, nil
}
