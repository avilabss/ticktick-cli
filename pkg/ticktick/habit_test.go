package ticktick

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHabitList_Success(t *testing.T) {
	habits := []Habit{
		{ID: "h1", Name: "Meditate", Type: "Boolean", Goal: 1, TotalCheckIns: 30},
		{ID: "h2", Name: "Read", Type: "Number", Goal: 30, Unit: "pages", TotalCheckIns: 15},
	}
	client := mockClientWithResponse(http.StatusOK, habits)

	result, err := client.Habit.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 habits, got %d", len(result))
	}
	if result[0].Name != "Meditate" {
		t.Errorf("expected name 'Meditate', got %q", result[0].Name)
	}
	if result[1].Unit != "pages" {
		t.Errorf("expected unit 'pages', got %q", result[1].Unit)
	}
}

func TestHabitGetCheckins_Success(t *testing.T) {
	var capturedBody map[string]any
	resp := HabitCheckinQueryResponse{
		Checkins: map[string][]HabitCheckin{
			"h1": {{ID: "c1", HabitID: "h1", CheckinStamp: 20260301, Value: 1, Status: 2}},
		},
	}
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			_ = json.Unmarshal(body, &capturedBody)
			jsonResp, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(jsonResp))),
			}, nil
		},
	}
	client, _ := NewTicktickClient("token", WithHTTPClient(mock))

	result, err := client.Habit.GetCheckins([]string{"h1"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Checkins["h1"]) != 1 {
		t.Fatalf("expected 1 checkin for h1, got %d", len(result.Checkins["h1"]))
	}
	if result.Checkins["h1"][0].CheckinStamp != 20260301 {
		t.Errorf("expected stamp 20260301, got %d", result.Checkins["h1"][0].CheckinStamp)
	}

	// Verify request body
	ids, ok := capturedBody["habitIds"].([]any)
	if !ok || len(ids) != 1 {
		t.Fatalf("expected habitIds with 1 item, got %v", capturedBody["habitIds"])
	}
}

func TestHabitCheckin_Success(t *testing.T) {
	var capturedBody BatchHabitCheckinRequest
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			_ = json.Unmarshal(body, &capturedBody)
			resp := BatchResponse{ID2Etag: map[string]string{"c1": "etag1"}}
			jsonResp, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(jsonResp))),
			}, nil
		},
	}
	client, _ := NewTicktickClient("token", WithHTTPClient(mock))

	checkin := HabitCheckin{
		ID:           "c1",
		HabitID:      "h1",
		CheckinStamp: 20260301,
		Value:        1,
		Goal:         1,
		Status:       2,
	}

	result, err := client.Habit.Checkin(checkin)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID2Etag["c1"] != "etag1" {
		t.Errorf("expected etag 'etag1', got %q", result.ID2Etag["c1"])
	}
	if len(capturedBody.Add) != 1 {
		t.Fatalf("expected 1 add item, got %d", len(capturedBody.Add))
	}
	if capturedBody.Add[0].HabitID != "h1" {
		t.Errorf("expected habitId 'h1', got %q", capturedBody.Add[0].HabitID)
	}
}

func TestHabitGetRecords_Success(t *testing.T) {
	var capturedBody map[string]any
	resp := struct {
		HabitRecords map[string][]HabitCheckin `json:"habitRecords"`
	}{
		HabitRecords: map[string][]HabitCheckin{
			"h1": {
				{ID: "c1", HabitID: "h1", CheckinStamp: 20260301, Value: 1},
				{ID: "c2", HabitID: "h1", CheckinStamp: 20260302, Value: 1},
			},
		},
	}
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			_ = json.Unmarshal(body, &capturedBody)
			jsonResp, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(jsonResp))),
			}, nil
		},
	}
	client, _ := NewTicktickClient("token", WithHTTPClient(mock))

	result, err := client.Habit.GetRecords(20260301, []string{"h1"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result["h1"]) != 2 {
		t.Fatalf("expected 2 records for h1, got %d", len(result["h1"]))
	}

	// Verify request body
	stamp, ok := capturedBody["afterStamp"].(float64)
	if !ok || int(stamp) != 20260301 {
		t.Errorf("expected afterStamp=20260301, got %v", capturedBody["afterStamp"])
	}
}
