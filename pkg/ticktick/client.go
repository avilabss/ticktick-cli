package ticktick

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/avilabss/ticktick-cli/pkg/logger"
)

// NewTicktickClient creates a new TickTick API client with sensible defaults.
//
// Options can be provided to override the default timeout and transport.
func NewTicktickClient(apiToken string, options ...Option) (*Client, error) {
	if apiToken == "" {
		return nil, fmt.Errorf("apiToken is required")
	}

	client := &Client{
		BaseURL:  "https://api.ticktick.com/api",
		APIToken: apiToken,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	client.Pomodoro = &PomodoroService{client: client}

	for _, option := range options {
		if err := option(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) Option {
	if timeout <= 0 {
		return func(c *Client) error {
			return fmt.Errorf("timeout must be greater than 0")
		}
	}

	return func(c *Client) error {
		c.HTTPClient.Timeout = timeout
		return nil
	}
}

// WithTransport sets the HTTP client transport.
func WithTransport(transport http.RoundTripper) Option {
	if transport == nil {
		return func(c *Client) error {
			return fmt.Errorf("transport is required")
		}
	}

	return func(c *Client) error {
		c.HTTPClient.Transport = transport
		return nil
	}
}

// Get sends a GET request to the specified endpoint and returns the response.
//
// The endpoint should be the path after the base URL, e.g. "/v2/pomodoros/timeline".
func (c *Client) Get(endpoint string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	logger.Trace("HTTP request", "method", "GET", "url", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", fmt.Sprintf("t=%s", c.APIToken))
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	slog.Debug("HTTP response", "status", res.StatusCode, "url", url)

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return res, nil
}
