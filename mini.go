// Package mini is a tiny, dependency-free HTTP client.
package mini

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

// Response wraps the parts of an HTTP response callers usually need.
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// Client issues HTTP requests with a default timeout and shared headers.
type Client struct {
	Timeout time.Duration
	Headers map[string]string
}

// New returns a Client with sensible defaults, tuned by opts.
func New(opts ...Option) *Client {
	c := &Client{Timeout: 30 * time.Second}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Option configures a Client.
type Option func(*Client)

// WithTimeout sets the per-request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) { c.Timeout = timeout }
}

// WithHeader sets a default header on every request.
func WithHeader(key, value string) Option {
	return func(c *Client) {
		if c.Headers == nil {
			c.Headers = make(map[string]string)
		}
		c.Headers[key] = value
	}
}

// Do sends a request with the given method, URL, and optional body.
func (c *Client) Do(method, url string, body []byte) (*Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "mini/0.1.0")
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/octet-stream")
	}
	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{StatusCode: resp.StatusCode, Header: resp.Header, Body: data}, nil
}

// Get performs a GET request.
func (c *Client) Get(url string) (*Response, error) {
	return c.Do(http.MethodGet, url, nil)
}

// Post performs a POST request with the given body.
func (c *Client) Post(url string, body []byte) (*Response, error) {
	return c.Do(http.MethodPost, url, body)
}

// Get performs a GET request with a default client.
func Get(url string) (*Response, error) {
	return New().Get(url)
}

// Post performs a POST request with a default client.
func Post(url string, body []byte) (*Response, error) {
	return New().Post(url, body)
}
