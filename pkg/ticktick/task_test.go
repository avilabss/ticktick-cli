package ticktick

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func mockSyncResponse() SyncResponse {
	return SyncResponse{
		SyncTaskBean: struct {
			Update []Task `json:"update"`
		}{
			Update: []Task{
				{ID: "t1", ProjectID: "p1", Title: "Active Task", Status: 0, Priority: 3, Tags: []string{"work"}},
				{ID: "t2", ProjectID: "p1", Title: "Completed Task", Status: 2},
				{ID: "t3", ProjectID: "p2", Title: "Another Task", Status: 0, Priority: 5},
			},
		},
		ProjectProfiles: []Project{
			{ID: "p1", Name: "Work", Color: "#FF0000"},
			{ID: "p2", Name: "Personal", Color: "#00FF00"},
		},
		Tags: []Tag{
			{Name: "work", Label: "work"},
			{Name: "personal", Label: "personal"},
		},
	}
}

func TestSync_Success(t *testing.T) {
	syncResp := mockSyncResponse()
	client := mockClientWithResponse(http.StatusOK, syncResp)

	result, err := client.Task.Sync()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.SyncTaskBean.Update) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(result.SyncTaskBean.Update))
	}
	if len(result.ProjectProfiles) != 2 {
		t.Errorf("expected 2 projects, got %d", len(result.ProjectProfiles))
	}
	if len(result.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(result.Tags))
	}
}

func TestSync_Cached(t *testing.T) {
	callCount := 0
	syncResp := mockSyncResponse()
	jsonBody, _ := json.Marshal(syncResp)

	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			callCount++
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(jsonBody))),
			}, nil
		},
	}
	client, _ := NewTicktickClient("token", WithHTTPClient(mock))

	_, _ = client.Task.Sync()
	_, _ = client.Task.Sync()

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

func TestList_FiltersActiveTasks(t *testing.T) {
	syncResp := mockSyncResponse()
	client := mockClientWithResponse(http.StatusOK, syncResp)

	tasks, err := client.Task.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 active tasks, got %d", len(tasks))
	}
	for _, task := range tasks {
		if task.Status != 0 {
			t.Errorf("expected only active tasks (status=0), got status=%d", task.Status)
		}
	}
}

func TestListProjects_Success(t *testing.T) {
	syncResp := mockSyncResponse()
	client := mockClientWithResponse(http.StatusOK, syncResp)

	projects, err := client.Task.ListProjects()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(projects))
	}
	if projects[0].Name != "Work" {
		t.Errorf("expected first project 'Work', got %q", projects[0].Name)
	}
}

func TestTaskCreate_Success(t *testing.T) {
	var capturedBody BatchTaskRequest
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			_ = json.Unmarshal(body, &capturedBody)
			resp := BatchResponse{
				ID2Etag:  map[string]string{"test-id": "etag1"},
				ID2Error: map[string]string{},
			}
			respBody, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(respBody))),
			}, nil
		},
	}
	client, _ := NewTicktickClient("token", WithHTTPClient(mock))

	task := Task{
		ID:        "test-id",
		Title:     "New Task",
		ProjectID: "p1",
		Priority:  3,
	}
	result, err := client.Task.Create(task)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(capturedBody.Add) != 1 {
		t.Fatalf("expected 1 task in add, got %d", len(capturedBody.Add))
	}
	if capturedBody.Add[0].Title != "New Task" {
		t.Errorf("expected title 'New Task', got %q", capturedBody.Add[0].Title)
	}
	if result.ID2Etag["test-id"] != "etag1" {
		t.Errorf("expected etag 'etag1', got %q", result.ID2Etag["test-id"])
	}
}

func TestTaskDelete_Success(t *testing.T) {
	var capturedBody []TaskDeleteRequest
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			_ = json.Unmarshal(body, &capturedBody)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}, nil
		},
	}
	client, _ := NewTicktickClient("token", WithHTTPClient(mock))

	err := client.Task.Delete("t1", "p1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(capturedBody) != 1 {
		t.Fatalf("expected 1 delete request, got %d", len(capturedBody))
	}
	if capturedBody[0].TaskID != "t1" {
		t.Errorf("expected taskId 't1', got %q", capturedBody[0].TaskID)
	}
	if capturedBody[0].ProjectID != "p1" {
		t.Errorf("expected projectId 'p1', got %q", capturedBody[0].ProjectID)
	}
}

func TestGenerateID(t *testing.T) {
	id := generateID()
	if len(id) != 24 {
		t.Errorf("expected 24-char ID, got %d chars: %q", len(id), id)
	}

	id2 := generateID()
	if id == id2 {
		t.Error("expected unique IDs")
	}
}
