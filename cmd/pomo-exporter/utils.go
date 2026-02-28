package main

import (
	"slices"
	"time"
)

func filter(items []string, exclude []string) []string {
	var result []string
	for _, item := range items {
		if slices.Contains(exclude, item) {
			continue
		}
		result = append(result, item)
	}
	return result
}

func monthRange(year int, month time.Month) (time.Time, time.Time) {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return start, end
}
