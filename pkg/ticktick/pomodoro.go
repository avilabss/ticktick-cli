package ticktick

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/avilabss/ticktick-cli/pkg/logger"
)

// PomodoroService handles pomodoro-related API calls.
type PomodoroService struct {
	client *Client
}

// GetTimeline returns pomodoros starting from the specified timestamp.
// If to is 0, it returns the latest pomodoros.
//
// The API sends 31 results by default.
func (s *PomodoroService) GetTimeline(to int64) (*Pomodoros, error) {
	endpoint := "/v2/pomodoros/timeline"
	if to > 0 {
		endpoint = fmt.Sprintf("%s?to=%d", endpoint, to)
	}

	slog.Debug("Fetching timeline page", "to", to)

	res, err := s.client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	var items []Pomodoro
	if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
		return nil, err
	}

	slog.Debug("Received timeline page", "count", len(items))

	return &Pomodoros{Items: items, service: s}, nil
}

// Next fetches the next batch of pomodoros based on the last item's start time.
func (p *Pomodoros) Next() (*Pomodoros, error) {
	if len(p.Items) == 0 {
		return nil, fmt.Errorf("no pomodoros to paginate from")
	}

	lastStartTime := p.Items[len(p.Items)-1].StartTime
	logger.Trace("Paginating", "lastStartTime", lastStartTime)

	to, err := time.Parse(TimeFormat, lastStartTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start time: %w", err)
	}

	return p.service.GetTimeline(to.UnixMilli())
}

// Stats returns general pomodoro statistics (today and total counts/durations).
func (s *PomodoroService) Stats() (*PomodoroStats, error) {
	res, err := s.client.Get("/v2/pomodoros/statistics/generalForDesktop")
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	var stats PomodoroStats
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// Create creates a new pomodoro record via the batch endpoint.
func (s *PomodoroService) Create(pomo Pomodoro) (*BatchResponse, error) {
	req := struct {
		Add    []Pomodoro `json:"add"`
		Update []any      `json:"update"`
		Delete []any      `json:"delete"`
	}{
		Add: []Pomodoro{pomo},
	}

	var result BatchResponse
	if err := s.client.PostJSON("/v2/batch/pomodoro", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeletePomo deletes a pomodoro by ID.
func (s *PomodoroService) DeletePomo(pomodoroID string) error {
	req := PomodoroDeleteRequest{
		PomodoroIDs: []string{pomodoroID},
		TimingIDs:   []string{},
	}
	return s.client.DeleteJSON("/v2/pomodoro", req, nil)
}

// ListTimers returns all custom focus timers.
func (s *PomodoroService) ListTimers() ([]FocusTimer, error) {
	res, err := s.client.Get("/v2/timer")
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	var timers []FocusTimer
	if err := json.NewDecoder(res.Body).Decode(&timers); err != nil {
		return nil, err
	}
	return timers, nil
}

// TimerOverview returns overview stats for a specific timer.
func (s *PomodoroService) TimerOverview(timerID string) (*FocusTimerOverview, error) {
	endpoint := fmt.Sprintf("/v2/timer/overview/%s", timerID)
	res, err := s.client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	var overview FocusTimerOverview
	if err := json.NewDecoder(res.Body).Decode(&overview); err != nil {
		return nil, err
	}
	return &overview, nil
}

// GetAll returns all pomodoros between start and end, paginating automatically.
func (s *PomodoroService) GetAll(start, end time.Time) (*Pomodoros, error) {
	var allItems []Pomodoro

	startUnix := start.UnixMilli()
	slog.Debug("GetAll", "start", start.Format(time.RFC3339), "end", end.Format(time.RFC3339))

	pageNum := 1
	currentPage, err := s.GetTimeline(end.UnixMilli())
	if err != nil {
		return nil, err
	}

	for len(currentPage.Items) > 0 {
		slog.Debug("Processing page", "page", pageNum, "items", len(currentPage.Items))

		reachedStart := false
		for _, p := range currentPage.Items {
			pTime, err := time.Parse(TimeFormat, p.StartTime)
			if err != nil {
				return nil, fmt.Errorf("failed to parse pomodoro start time: %w", err)
			}
			if pTime.UnixMilli() < startUnix {
				slog.Debug("Reached start boundary", "at", p.StartTime)
				reachedStart = true
				break
			}
			allItems = append(allItems, p)
		}

		if reachedStart {
			break
		}

		nextPage, err := currentPage.Next()
		if err != nil {
			return nil, err
		}

		if len(nextPage.Items) > 0 && nextPage.Items[0].ID == currentPage.Items[0].ID {
			slog.Debug("Duplicate page detected, stopping pagination")
			break
		}

		currentPage = nextPage
		pageNum++
	}

	slog.Debug("GetAll complete", "total", len(allItems), "pages", pageNum)
	return &Pomodoros{Items: allItems, service: s}, nil
}
