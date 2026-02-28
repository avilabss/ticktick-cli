package ticktick

import "net/http"

const TimeFormat = "2006-01-02T15:04:05.000-0700"

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
