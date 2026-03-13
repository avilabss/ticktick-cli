package ticktick

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

// HabitService handles habit-related API calls.
type HabitService struct {
	client *Client
}

// List returns all habits.
func (s *HabitService) List() ([]Habit, error) {
	slog.Debug("Fetching habits")

	res, err := s.client.Get("/v2/habits")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch habits: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	var habits []Habit
	if err := json.NewDecoder(res.Body).Decode(&habits); err != nil {
		return nil, fmt.Errorf("failed to decode habits: %w", err)
	}

	slog.Debug("Fetched habits", "count", len(habits))
	return habits, nil
}

// GetCheckins returns check-in history for the given habit IDs.
func (s *HabitService) GetCheckins(habitIDs []string) (*HabitCheckinQueryResponse, error) {
	req := struct {
		HabitIDs []string `json:"habitIds"`
	}{HabitIDs: habitIDs}

	var result HabitCheckinQueryResponse
	if err := s.client.PostJSON("/v2/habitCheckins/query", req, &result); err != nil {
		return nil, fmt.Errorf("failed to query checkins: %w", err)
	}
	return &result, nil
}

// Checkin creates a check-in for a habit.
func (s *HabitService) Checkin(checkin HabitCheckin) (*BatchResponse, error) {
	req := BatchHabitCheckinRequest{
		Add:    []HabitCheckin{checkin},
		Update: []HabitCheckin{},
		Delete: []string{},
	}

	var result BatchResponse
	if err := s.client.PostJSON("/v2/habitCheckins/batch", req, &result); err != nil {
		return nil, fmt.Errorf("failed to checkin habit: %w", err)
	}

	slog.Info("Habit checked in", "habitId", checkin.HabitID)
	return &result, nil
}

// GetRecords returns habit records for a date range.
func (s *HabitService) GetRecords(afterStamp int, habitIDs []string) (map[string][]HabitCheckin, error) {
	req := struct {
		AfterStamp int      `json:"afterStamp"`
		HabitIDs   []string `json:"habitIds"`
	}{AfterStamp: afterStamp, HabitIDs: habitIDs}

	var result struct {
		HabitRecords map[string][]HabitCheckin `json:"habitRecords"`
	}
	if err := s.client.PostJSON("/v2/getHabitRecords", req, &result); err != nil {
		return nil, fmt.Errorf("failed to get habit records: %w", err)
	}
	return result.HabitRecords, nil
}
