package ticktick

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	APIToken   string
	HTTPClient *http.Client
}

type Option func(*Client) error

type Task struct {
	TaskID      string   `json:"taskId"`
	Title       string   `json:"title"`
	ProjectName string   `json:"projectName"`
	Tags        []string `json:"tags"`
	StartTime   string   `json:"startTime"`
	EndTime     string   `json:"endTime"`
}

type Pomodoro struct {
	ID            string `json:"id"`
	Tasks         []Task `json:"tasks"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	Status        int    `json:"status"`
	PauseDuration int    `json:"pauseDuration"`
	AdjustTime    int    `json:"adjustTime"`
	Etag          string `json:"etag"`
	Type          int    `json:"type"`
	Added         bool   `json:"added"`
}

// NewTicktickClient creates a new TickTick API client with sensible defaults.
//
// Options can be provided to override the default timeout and transport.
//
// The default timeout is 30 seconds. Use WithTimeout to set a custom timeout.
//
// The default transport is http.DefaultTransport. Use WithTransport to set a custom transport.
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
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", c.BaseURL, endpoint), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", fmt.Sprintf("t=%s", c.APIToken))
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return res, nil
}

// GetPomodorosTimeline returns pomodoros starting from the specified timestamp.
// If to is 0, it returns the latest pomodoros.
//
// The API sends 31 results by default.
// Use "startTime" converted in unix timestamp of the last record as "to" param for the request to get next set of results.
func (c *Client) GetPomodorosTimeline(to int64) ([]Pomodoro, error) {
	endpoint := "/v2/pomodoros/timeline"
	if to > 0 {
		endpoint = fmt.Sprintf("%s?to=%d", endpoint, to)
	}

	res, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var pomodoros []Pomodoro
	err = json.NewDecoder(res.Body).Decode(&pomodoros)
	if err != nil {
		return nil, err
	}

	return pomodoros, nil
}

// GetNextPomodorosTimeline returns the next set of pomodoros after the last one in previousPomodoros.
func (c *Client) GetNextPomodorosTimeline(previousPomodoros []Pomodoro) ([]Pomodoro, error) {
	if len(previousPomodoros) == 0 {
		return nil, fmt.Errorf("previousPomodoros cannot be empty")
	}

	lastPomodoro := previousPomodoros[len(previousPomodoros)-1]
	lastStartTime := lastPomodoro.StartTime

	to, err := time.Parse(TimeFormat, lastStartTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse last pomodoro start time: %w", err)
	}

	return c.GetPomodorosTimeline(to.UnixMilli())
}

func (c *Client) GetAllPomodorosTimeline(start time.Time, end time.Time) ([]Pomodoro, error) {
	var allPomodoros []Pomodoro

	startUnix := start.UnixMilli()
	endUnix := end.UnixMilli()
	currentTo := endUnix

	for {
		pomodoros, err := c.GetPomodorosTimeline(currentTo)
		if err != nil {
			return nil, err
		}

		if len(pomodoros) == 0 {
			break
		}

		reachedStart := false
		for _, p := range pomodoros {
			pTime, err := time.Parse(TimeFormat, p.StartTime)
			if err != nil {
				return nil, fmt.Errorf("failed to parse pomodoro start time: %w", err)
			}
			pUnix := pTime.UnixMilli()
			if pUnix < startUnix {
				reachedStart = true
				break
			}
			allPomodoros = append(allPomodoros, p)
			currentTo = pUnix
		}

		if reachedStart {
			break
		}
	}

	return allPomodoros, nil
}
