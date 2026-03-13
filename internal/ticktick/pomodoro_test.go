package ticktick

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func mockClientWithResponse(statusCode int, body any) *Client {
	jsonBody, _ := json.Marshal(body)
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader(string(jsonBody))),
			}, nil
		},
	}
	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	return client
}

func makePomodoro(id string, startTime string) Pomodoro {
	return Pomodoro{
		ID:        id,
		StartTime: startTime,
		EndTime:   startTime,
		Tasks:     []PomodoroTask{},
	}
}

func TestGetTimeline_Success(t *testing.T) {
	pomodoros := []Pomodoro{
		makePomodoro("1", "2026-02-10T03:27:47.000+0000"),
		makePomodoro("2", "2026-02-09T12:00:00.000+0000"),
	}

	client := mockClientWithResponse(http.StatusOK, pomodoros)
	result, err := client.Pomodoro.GetTimeline(0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Items))
	}
}

func TestGetTimeline_WithTimestamp(t *testing.T) {
	var capturedURL string
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedURL = req.URL.String()
			body, _ := json.Marshal([]Pomodoro{})
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(body))),
			}, nil
		},
	}
	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	_, _ = client.Pomodoro.GetTimeline(1234567890000)

	if !strings.Contains(capturedURL, "?to=1234567890000") {
		t.Errorf("expected URL to contain '?to=1234567890000', got %q", capturedURL)
	}
}

func TestGetTimeline_EmptyResponse(t *testing.T) {
	client := mockClientWithResponse(http.StatusOK, []Pomodoro{})
	result, err := client.Pomodoro.GetTimeline(0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.Items))
	}
}

func TestGetTimeline_APIError(t *testing.T) {
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}
	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	_, err := client.Pomodoro.GetTimeline(0)
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestNext_Success(t *testing.T) {
	callCount := 0
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			callCount++
			var body []Pomodoro
			if callCount == 1 {
				body = []Pomodoro{
					makePomodoro("1", "2026-02-10T03:27:47.000+0000"),
					makePomodoro("2", "2026-02-09T12:00:00.000+0000"),
				}
			} else {
				body = []Pomodoro{
					makePomodoro("3", "2026-02-08T10:00:00.000+0000"),
				}
			}
			jsonBody, _ := json.Marshal(body)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(jsonBody))),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	first, _ := client.Pomodoro.GetTimeline(0)
	next, err := first.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(next.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(next.Items))
	}
	if next.Items[0].ID != "3" {
		t.Errorf("expected ID '3', got %q", next.Items[0].ID)
	}
}

func TestNext_EmptyItems(t *testing.T) {
	result := &Pomodoros{Items: []Pomodoro{}}
	_, err := result.Next()
	if err == nil {
		t.Fatal("expected error for empty items")
	}
}

func TestGetAll_StopsAtStartBoundary(t *testing.T) {
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := []Pomodoro{
				makePomodoro("1", "2026-02-15T10:00:00.000+0000"),
				makePomodoro("2", "2026-02-10T10:00:00.000+0000"),
				makePomodoro("3", "2026-01-25T10:00:00.000+0000"), // before start
			}
			jsonBody, _ := json.Marshal(body)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(jsonBody))),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC)

	result, err := client.Pomodoro.GetAll(start, end)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Items) != 2 {
		t.Errorf("expected 2 items (before start excluded), got %d", len(result.Items))
	}
}

func TestGetAll_DuplicatePageDetection(t *testing.T) {
	callCount := 0
	samePomodoros := []Pomodoro{
		makePomodoro("1", "2026-02-15T10:00:00.000+0000"),
	}

	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			callCount++
			jsonBody, _ := json.Marshal(samePomodoros)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(jsonBody))),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC)

	result, err := client.Pomodoro.GetAll(start, end)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should stop after detecting duplicate, not loop forever
	if callCount > 3 {
		t.Errorf("expected loop to stop early, but made %d API calls", callCount)
	}
	if len(result.Items) == 0 {
		t.Error("expected at least some items")
	}
}

func TestGetAll_EmptyResponse(t *testing.T) {
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			jsonBody, _ := json.Marshal([]Pomodoro{})
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(jsonBody))),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC)

	result, err := client.Pomodoro.GetAll(start, end)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.Items))
	}
}

func TestGetAll_MultiplePages(t *testing.T) {
	callCount := 0
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			callCount++
			var body []Pomodoro
			switch callCount {
			case 1:
				body = []Pomodoro{
					makePomodoro("1", "2026-02-20T10:00:00.000+0000"),
					makePomodoro("2", "2026-02-15T10:00:00.000+0000"),
				}
			case 2:
				body = []Pomodoro{
					makePomodoro("3", "2026-02-10T10:00:00.000+0000"),
					makePomodoro("4", "2026-02-05T10:00:00.000+0000"),
				}
			case 3:
				body = []Pomodoro{
					makePomodoro("5", "2026-01-30T10:00:00.000+0000"), // before start
				}
			default:
				body = []Pomodoro{}
			}
			jsonBody, _ := json.Marshal(body)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(jsonBody))),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC)

	result, err := client.Pomodoro.GetAll(start, end)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Items) != 4 {
		t.Errorf("expected 4 items across 2 pages, got %d", len(result.Items))
	}
}

func TestGetAll_APIError(t *testing.T) {
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("network error")
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC)

	_, err := client.Pomodoro.GetAll(start, end)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
}

func TestStats_Success(t *testing.T) {
	stats := PomodoroStats{
		TodayPomoCount:    3,
		TotalPomoCount:    150,
		TodayPomoDuration: 4500,
		TotalPomoDuration: 225000,
	}
	client := mockClientWithResponse(http.StatusOK, stats)

	result, err := client.Pomodoro.Stats()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.TodayPomoCount != 3 {
		t.Errorf("expected TodayPomoCount=3, got %d", result.TodayPomoCount)
	}
	if result.TotalPomoCount != 150 {
		t.Errorf("expected TotalPomoCount=150, got %d", result.TotalPomoCount)
	}
}

func TestCreate_Success(t *testing.T) {
	var capturedBody map[string]any
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			_ = json.Unmarshal(body, &capturedBody)
			resp := BatchResponse{ID2Etag: map[string]string{"p1": "etag1"}}
			jsonResp, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(jsonResp))),
			}, nil
		},
	}
	client, _ := NewTicktickClient("token", WithHTTPClient(mock))

	pomo := Pomodoro{
		ID:        "test-pomo-id",
		StartTime: "2026-02-10T10:00:00.000+0000",
		EndTime:   "2026-02-10T10:25:00.000+0000",
		Status:    1,
	}

	result, err := client.Pomodoro.Create(pomo)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID2Etag["p1"] != "etag1" {
		t.Errorf("expected etag 'etag1', got %q", result.ID2Etag["p1"])
	}

	addList, ok := capturedBody["add"].([]any)
	if !ok || len(addList) != 1 {
		t.Fatalf("expected 1 item in add list, got %v", capturedBody["add"])
	}
}

func TestDeletePomo_Success(t *testing.T) {
	var capturedBody PomodoroDeleteRequest
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			body, _ := io.ReadAll(req.Body)
			_ = json.Unmarshal(body, &capturedBody)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}, nil
		},
	}
	client, _ := NewTicktickClient("token", WithHTTPClient(mock))

	err := client.Pomodoro.DeletePomo("pomo-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(capturedBody.PomodoroIDs) != 1 || capturedBody.PomodoroIDs[0] != "pomo-123" {
		t.Errorf("expected pomodoroIds=[pomo-123], got %v", capturedBody.PomodoroIDs)
	}
}

func TestListTimers_Success(t *testing.T) {
	timers := []FocusTimer{
		{ID: "t1", Name: "Focus", Type: "pomodoro", PomodoroTime: 25},
		{ID: "t2", Name: "Deep Work", Type: "pomodoro", PomodoroTime: 50},
	}
	client := mockClientWithResponse(http.StatusOK, timers)

	result, err := client.Pomodoro.ListTimers()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 timers, got %d", len(result))
	}
	if result[0].Name != "Focus" {
		t.Errorf("expected name 'Focus', got %q", result[0].Name)
	}
}

func TestTimerOverview_Success(t *testing.T) {
	var capturedURL string
	overview := FocusTimerOverview{Days: 30, Today: 1500, Total: 90000}
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedURL = req.URL.String()
			jsonBody, _ := json.Marshal(overview)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(jsonBody))),
			}, nil
		},
	}
	client, _ := NewTicktickClient("token", WithHTTPClient(mock))

	result, err := client.Pomodoro.TimerOverview("timer-abc")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(capturedURL, "/v2/timer/overview/timer-abc") {
		t.Errorf("expected URL to contain timer ID, got %q", capturedURL)
	}
	if result.Days != 30 {
		t.Errorf("expected Days=30, got %d", result.Days)
	}
	if result.Total != 90000 {
		t.Errorf("expected Total=90000, got %d", result.Total)
	}
}
