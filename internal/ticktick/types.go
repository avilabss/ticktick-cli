package ticktick

import "net/http"

const TimeFormat = "2006-01-02T15:04:05.000-0700"

// HTTPClient is the interface for making HTTP requests.
// *http.Client satisfies this interface.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	BaseURL    string
	APIToken   string
	HTTPClient HTTPClient

	Pomodoro *PomodoroService
	Task     *TaskService
	Habit    *HabitService
}

type Option func(*Client) error

type PomodoroTask struct {
	TaskID      string   `json:"taskId"`
	Title       string   `json:"title"`
	ProjectName string   `json:"projectName"`
	Tags        []string `json:"tags"`
	StartTime   string   `json:"startTime"`
	EndTime     string   `json:"endTime"`
}

type Pomodoro struct {
	ID            string         `json:"id"`
	Tasks         []PomodoroTask `json:"tasks"`
	StartTime     string         `json:"startTime"`
	EndTime       string         `json:"endTime"`
	Status        int            `json:"status"`
	PauseDuration int            `json:"pauseDuration"`
	AdjustTime    int            `json:"adjustTime"`
	Etag          string         `json:"etag"`
	Type          int            `json:"type"`
	Added         bool           `json:"added"`
}

// Pomodoros is a result set that supports pagination.
type Pomodoros struct {
	Items   []Pomodoro
	service *PomodoroService
}

// BatchResponse is the common response type for all TickTick batch endpoints.
type BatchResponse struct {
	ID2Etag  map[string]string `json:"id2etag"`
	ID2Error map[string]string `json:"id2error"`
}

// Task represents a TickTick task.
type Task struct {
	ID           string    `json:"id"`
	ProjectID    string    `json:"projectId"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Desc         string    `json:"desc,omitempty"`
	Tags         []string  `json:"tags"`
	Priority     int       `json:"priority"`        // 0=none, 1=low, 3=medium, 5=high
	Status       int       `json:"status"`           // 0=active, 2=completed
	StartDate    string    `json:"startDate,omitempty"`
	DueDate      string    `json:"dueDate,omitempty"`
	TimeZone     string    `json:"timeZone"`
	IsAllDay     bool      `json:"isAllDay"`
	IsFloating   bool      `json:"isFloating"`
	Reminders    []any     `json:"reminders"`
	Items        []SubTask `json:"items"`
	SortOrder    int64     `json:"sortOrder"`
	RepeatFlag   string    `json:"repeatFlag,omitempty"`
	ExDate       []string  `json:"exDate"`
	Etag         string    `json:"etag,omitempty"`
	CreatedTime  string    `json:"createdTime"`
	ModifiedTime string    `json:"modifiedTime"`
}

// SubTask represents a checklist item within a task.
type SubTask struct {
	ID        string `json:"id"`
	Status    int    `json:"status"`
	Title     string `json:"title"`
	SortOrder int    `json:"sortOrder"`
}

// Project represents a TickTick project (list).
type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	SortOrder int64  `json:"sortOrder"`
	ViewMode  string `json:"viewMode"`
	GroupID   string `json:"groupId"`
	IsOwner   bool   `json:"isOwner"`
	Etag      string `json:"etag"`
}

// Tag represents a TickTick tag.
type Tag struct {
	Name      string `json:"name"`
	Label     string `json:"label"`
	SortOrder int64  `json:"sortOrder"`
	Color     string `json:"color"`
	Parent    string `json:"parent"`
	Etag      string `json:"etag"`
}

// SyncResponse represents the response from GET /v3/batch/check/0.
type SyncResponse struct {
	SyncTaskBean struct {
		Update []Task `json:"update"`
	} `json:"syncTaskBean"`
	ProjectProfiles []Project `json:"projectProfiles"`
	Tags            []Tag     `json:"tags"`
	InboxID         string    `json:"inboxId"`
}

// BatchTaskRequest is the request body for POST /v2/batch/task.
type BatchTaskRequest struct {
	Add               []Task `json:"add"`
	Update            []Task `json:"update"`
	Delete            []any  `json:"delete"`
	AddAttachments    []any  `json:"addAttachments"`
	UpdateAttachments []any  `json:"updateAttachments"`
	DeleteAttachments []any  `json:"deleteAttachments"`
}

// TaskDeleteRequest identifies a task for deletion.
type TaskDeleteRequest struct {
	TaskID    string `json:"taskId"`
	ProjectID string `json:"projectId"`
}

// PomodoroStats represents pomodoro statistics.
type PomodoroStats struct {
	TodayPomoCount    int `json:"todayPomoCount"`
	TotalPomoCount    int `json:"totalPomoCount"`
	TodayPomoDuration int `json:"todayPomoDuration"`
	TotalPomoDuration int `json:"totalPomoDuration"`
}

// PomodoroDeleteRequest is the request body for DELETE /v2/pomodoro.
type PomodoroDeleteRequest struct {
	PomodoroIDs []string `json:"pomodoroIds"`
	TimingIDs   []string `json:"timingIds"`
}

// FocusTimer represents a custom focus timer.
type FocusTimer struct {
	ID           string `json:"id"`
	Icon         string `json:"icon"`
	Color        string `json:"color"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	PomodoroTime int    `json:"pomodoroTime"`
	Status       int    `json:"status"`
	SortOrder    int64  `json:"sortOrder"`
	Etag         string `json:"etag"`
}

// FocusTimerOverview represents overview stats for a timer.
type FocusTimerOverview struct {
	Days  int `json:"days"`
	Today int `json:"today"`
	Total int `json:"total"`
}

// Habit represents a TickTick habit.
type Habit struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	IconRes       string   `json:"iconRes"`
	Color         string   `json:"color"`
	SortOrder     int64    `json:"sortOrder"`
	Status        int      `json:"status"`
	Type          string   `json:"type"` // "Boolean" or "Number"
	Goal          float64  `json:"goal"`
	Step          float64  `json:"step"`
	Unit          string   `json:"unit"`
	RepeatRule    string   `json:"repeatRule"`
	Reminders     []string `json:"reminders"`
	TotalCheckIns int      `json:"totalCheckIns"`
	SectionID     string   `json:"sectionId"`
	Etag          string   `json:"etag"`
}

// HabitCheckin represents a single check-in for a habit.
type HabitCheckin struct {
	ID           string  `json:"id"`
	HabitID      string  `json:"habitId"`
	CheckinStamp int     `json:"checkinStamp"` // YYYYMMDD as int
	CheckinTime  string  `json:"checkinTime"`
	OpTime       string  `json:"opTime"`
	Value        float64 `json:"value"`
	Goal         float64 `json:"goal"`
	Status       int     `json:"status"` // 2=done
}

// HabitCheckinQueryResponse is the response from POST /v2/habitCheckins/query.
type HabitCheckinQueryResponse struct {
	Checkins map[string][]HabitCheckin `json:"checkins"`
}

// BatchHabitCheckinRequest is the request for POST /v2/habitCheckins/batch.
type BatchHabitCheckinRequest struct {
	Add    []HabitCheckin `json:"add"`
	Update []HabitCheckin `json:"update"`
	Delete []string       `json:"delete"`
}
