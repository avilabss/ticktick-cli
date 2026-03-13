package ticktick

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/avilabss/ticktick-cli/internal/logger"
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
	client.Task = &TaskService{client: client}
	client.Habit = &HabitService{client: client}

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

// do sends an HTTP request and returns the response.
// It sets the cookie header and content-type for requests with a body.
// Both 200 and 204 are accepted as success status codes.
func (c *Client) do(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	logger.Trace("HTTP request", "method", method, "url", url)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", fmt.Sprintf("t=%s", c.APIToken))
	req.Header.Set("Origin", "https://ticktick.com")
	req.Header.Set("x-device", `{"platform":"web","os":"","device":"","name":"","version":4531,"id":"","channel":"website","campaign":""}`)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	slog.Debug("HTTP response", "status", res.StatusCode, "url", url)

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()
		slog.Debug("Error response body", "status", res.StatusCode, "body", string(body))
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return res, nil
}

// Get sends a GET request to the specified endpoint.
func (c *Client) Get(endpoint string) (*http.Response, error) {
	return c.do(http.MethodGet, endpoint, nil)
}

// Post sends a POST request with a body to the specified endpoint.
func (c *Client) Post(endpoint string, body io.Reader) (*http.Response, error) {
	return c.do(http.MethodPost, endpoint, body)
}

// Put sends a PUT request with a body to the specified endpoint.
func (c *Client) Put(endpoint string, body io.Reader) (*http.Response, error) {
	return c.do(http.MethodPut, endpoint, body)
}

// Delete sends a DELETE request with an optional body to the specified endpoint.
func (c *Client) Delete(endpoint string, body io.Reader) (*http.Response, error) {
	return c.do(http.MethodDelete, endpoint, body)
}

// PostJSON marshals v to JSON, sends a POST request, and decodes the response into result.
// If result is nil, the response body is discarded.
func (c *Client) PostJSON(endpoint string, v any, result any) error {
	return c.doJSON(http.MethodPost, endpoint, v, result)
}

// PutJSON marshals v to JSON, sends a PUT request, and decodes the response into result.
func (c *Client) PutJSON(endpoint string, v any, result any) error {
	return c.doJSON(http.MethodPut, endpoint, v, result)
}

// DeleteJSON marshals v to JSON, sends a DELETE request, and decodes the response into result.
func (c *Client) DeleteJSON(endpoint string, v any, result any) error {
	return c.doJSON(http.MethodDelete, endpoint, v, result)
}

func (c *Client) doJSON(method, endpoint string, v any, result any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	logger.Trace("Request body", "endpoint", endpoint, "body", string(data))

	res, err := c.do(method, endpoint, bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	if result != nil {
		if err := json.NewDecoder(res.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}
