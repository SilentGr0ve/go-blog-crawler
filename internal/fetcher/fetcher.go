package fetcher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Options struct {
	UserAgent string
	Timeout   time.Duration
	MaxBytes  int64
}

type Client struct {
	client    *http.Client
	userAgent string
	maxBytes  int64
}

func NewClient(options Options) *Client {
	if options.UserAgent == "" {
		options.UserAgent = "go-blog-crawler/0.1 (+https://github.com/SilentGr0ve/go-blog-crawler)"
	}
	if options.Timeout == 0 {
		options.Timeout = 10 * time.Second
	}
	if options.MaxBytes == 0 {
		options.MaxBytes = 5 * 1024 * 1024
	}

	return &Client{
		client:    &http.Client{Timeout: options.Timeout},
		userAgent: options.UserAgent,
		maxBytes:  options.MaxBytes,
	}
}

func (c *Client) Fetch(ctx context.Context, url string) (body []byte, statusCode int, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("fetcher: create request for %s: %w", url, err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			return nil, 0, fmt.Errorf("fetcher: request cancelled (ctx): %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return nil, 0, fmt.Errorf("fetcher: request timed out (ctx): %w", err)
		default:
			return nil, 0, fmt.Errorf("fetcher: do request %s: %w", url, err)
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	body, err = io.ReadAll(io.LimitReader(resp.Body, c.maxBytes))
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response body: %w", err)
	}

	return body, http.StatusOK, nil
}
