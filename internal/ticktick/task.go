package ticktick

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

// TaskService handles task and project related API calls.
type TaskService struct {
	client    *Client
	syncCache *SyncResponse
}

// generateID creates a random 24-character hex ID matching TickTick's format.
func generateID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Sync performs a full sync and returns all tasks, projects, and tags.
// Results are cached for the lifetime of the service instance.
func (s *TaskService) Sync() (*SyncResponse, error) {
	if s.syncCache != nil {
		return s.syncCache, nil
	}

	slog.Debug("Performing full sync")

	res, err := s.client.Get("/v3/batch/check/0")
	if err != nil {
		return nil, fmt.Errorf("sync failed: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	var result SyncResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode sync response: %w", err)
	}

	slog.Debug("Sync complete",
		"tasks", len(result.SyncTaskBean.Update),
		"projects", len(result.ProjectProfiles),
		"tags", len(result.Tags))

	s.syncCache = &result
	return s.syncCache, nil
}

// List returns all active tasks from the sync.
func (s *TaskService) List() ([]Task, error) {
	sync, err := s.Sync()
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for _, t := range sync.SyncTaskBean.Update {
		if t.Status == 0 {
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}

// Get returns a single task by ID.
// If projectID is empty, it will be resolved from the sync data.
func (s *TaskService) Get(taskID, projectID string) (*Task, error) {
	pid, err := s.resolveProjectID(taskID, projectID)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("/v2/task/%s?projectId=%s", taskID, pid)
	res, err := s.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	var task Task
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode task: %w", err)
	}
	return &task, nil
}

// Create creates a new task.
func (s *TaskService) Create(task Task) (*BatchResponse, error) {
	if task.ID == "" {
		task.ID = generateID()
	}

	if task.ProjectID == "" {
		sync, err := s.Sync()
		if err != nil {
			return nil, fmt.Errorf("failed to resolve inbox: %w", err)
		}
		task.ProjectID = sync.InboxID
	}

	now := time.Now().UTC().Format(TimeFormat)
	if task.CreatedTime == "" {
		task.CreatedTime = now
	}
	if task.ModifiedTime == "" {
		task.ModifiedTime = now
	}

	// TickTick API requires empty arrays, not null.
	if task.Tags == nil {
		task.Tags = []string{}
	}
	if task.Items == nil {
		task.Items = []SubTask{}
	}
	if task.Reminders == nil {
		task.Reminders = []any{}
	}
	if task.ExDate == nil {
		task.ExDate = []string{}
	}

	req := BatchTaskRequest{
		Add:               []Task{task},
		Update:            []Task{},
		Delete:            []any{},
		AddAttachments:    []any{},
		UpdateAttachments: []any{},
		DeleteAttachments: []any{},
	}

	var result BatchResponse
	if err := s.client.PostJSON("/v2/batch/task", req, &result); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	slog.Info("Task created", "id", task.ID, "title", task.Title)
	return &result, nil
}

// Complete marks a task as completed.
// If projectID is empty, it will be resolved from the sync data.
func (s *TaskService) Complete(taskID, projectID string) (*BatchResponse, error) {
	task, err := s.Get(taskID, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task for completion: %w", err)
	}

	task.Status = 2
	task.ModifiedTime = time.Now().UTC().Format(TimeFormat)

	req := BatchTaskRequest{
		Add:               []Task{},
		Update:            []Task{*task},
		Delete:            []any{},
		AddAttachments:    []any{},
		UpdateAttachments: []any{},
		DeleteAttachments: []any{},
	}

	var result BatchResponse
	if err := s.client.PostJSON("/v2/batch/task", req, &result); err != nil {
		return nil, fmt.Errorf("failed to complete task: %w", err)
	}

	slog.Info("Task completed", "id", taskID, "title", task.Title)
	return &result, nil
}

// Delete permanently deletes a task.
// If projectID is empty, it will be resolved from the sync data.
func (s *TaskService) Delete(taskID, projectID string) error {
	pid, err := s.resolveProjectID(taskID, projectID)
	if err != nil {
		return err
	}

	body := []TaskDeleteRequest{{TaskID: taskID, ProjectID: pid}}
	return s.client.DeleteJSON("/v2/task?deleteforever=true", body, nil)
}

// ListProjects returns all projects from the sync.
func (s *TaskService) ListProjects() ([]Project, error) {
	sync, err := s.Sync()
	if err != nil {
		return nil, err
	}
	return sync.ProjectProfiles, nil
}

// resolveProjectID returns projectID if non-empty, otherwise looks it up from sync data.
func (s *TaskService) resolveProjectID(taskID, projectID string) (string, error) {
	if projectID != "" {
		return projectID, nil
	}

	sync, err := s.Sync()
	if err != nil {
		return "", err
	}

	for _, t := range sync.SyncTaskBean.Update {
		if t.ID == taskID {
			return t.ProjectID, nil
		}
	}
	return "", fmt.Errorf("task %s not found in sync data", taskID)
}
