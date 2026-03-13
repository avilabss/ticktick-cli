package pomodoro

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func includeExclude(items []string, include []string, exclude []string) []string {
	var result []string
	for _, item := range items {
		if len(include) > 0 && !slices.Contains(include, item) {
			continue
		}
		if slices.Contains(exclude, item) {
			continue
		}
		result = append(result, item)
	}
	return result
}

func matchesFilter(value string, include []string, exclude []string) bool {
	if len(include) > 0 && !slices.Contains(include, value) {
		return false
	}
	if slices.Contains(exclude, value) {
		return false
	}
	return true
}

func monthRange(year int, month time.Month) (time.Time, time.Time) {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return start, end
}

// formatDuration converts seconds to a human-readable string like "1h 25m" or "45m".
func formatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	h := seconds / 3600
	m := (seconds % 3600) / 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}
