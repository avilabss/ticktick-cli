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
// Options can be provided to override the default HTTP client.
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

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient HTTPClient) Option {
	if httpClient == nil {
		return func(c *Client) error {
			return fmt.Errorf("httpClient is required")
		}
	}

	return func(c *Client) error {
		c.HTTPClient = httpClient
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
