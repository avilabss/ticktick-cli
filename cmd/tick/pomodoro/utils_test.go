package pomodoro

import (
	"slices"
	"testing"
	"time"
)

func TestSplitCSV(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty string", "", nil},
		{"single value", "foo", []string{"foo"}},
		{"multiple values", "foo,bar,baz", []string{"foo", "bar", "baz"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitCSV(tt.input)
			if !slices.Equal(result, tt.expected) {
				t.Errorf("splitCSV(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIncludeExclude(t *testing.T) {
	tests := []struct {
		name     string
		items    []string
		include  []string
		exclude  []string
		expected []string
	}{
		{
			name:     "no filters returns all",
			items:    []string{"a", "b", "c"},
			include:  nil,
			exclude:  nil,
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "include only",
			items:    []string{"a", "b", "c"},
			include:  []string{"a", "c"},
			exclude:  nil,
			expected: []string{"a", "c"},
		},
		{
			name:     "exclude only",
			items:    []string{"a", "b", "c"},
			include:  nil,
			exclude:  []string{"b"},
			expected: []string{"a", "c"},
		},
		{
			name:     "include and exclude",
			items:    []string{"a", "b", "c"},
			include:  []string{"a", "b"},
			exclude:  []string{"b"},
			expected: []string{"a"},
		},
		{
			name:     "empty items",
			items:    []string{},
			include:  []string{"a"},
			exclude:  nil,
			expected: nil,
		},
		{
			name:     "no matches",
			items:    []string{"a", "b"},
			include:  []string{"x"},
			exclude:  nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := includeExclude(tt.items, tt.include, tt.exclude)
			if !slices.Equal(result, tt.expected) {
				t.Errorf("includeExclude(%v, %v, %v) = %v, want %v",
					tt.items, tt.include, tt.exclude, result, tt.expected)
			}
		})
	}
}

func TestMatchesFilter(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		include  []string
		exclude  []string
		expected bool
	}{
		{"no filters", "anything", nil, nil, true},
		{"include match", "a", []string{"a", "b"}, nil, true},
		{"include miss", "c", []string{"a", "b"}, nil, false},
		{"exclude match", "a", nil, []string{"a"}, false},
		{"exclude miss", "b", nil, []string{"a"}, true},
		{"include and exclude match", "a", []string{"a"}, []string{"a"}, false},
		{"include match exclude miss", "a", []string{"a"}, []string{"b"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesFilter(tt.value, tt.include, tt.exclude)
			if result != tt.expected {
				t.Errorf("matchesFilter(%q, %v, %v) = %v, want %v",
					tt.value, tt.include, tt.exclude, result, tt.expected)
			}
		})
	}
}

func TestMonthRange(t *testing.T) {
	tests := []struct {
		name      string
		year      int
		month     time.Month
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name:      "february 2026",
			year:      2026,
			month:     time.February,
			wantStart: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
		},
		{
			name:      "leap year february",
			year:      2024,
			month:     time.February,
			wantStart: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
		},
		{
			name:      "december wraps to next year",
			year:      2026,
			month:     time.December,
			wantStart: time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := monthRange(tt.year, tt.month)
			if !start.Equal(tt.wantStart) {
				t.Errorf("start = %v, want %v", start, tt.wantStart)
			}
			if !end.Equal(tt.wantEnd) {
				t.Errorf("end = %v, want %v", end, tt.wantEnd)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int
		expected string
	}{
		{"zero", 0, "0s"},
		{"seconds only", 45, "45s"},
		{"one minute", 60, "1m"},
		{"minutes only", 1500, "25m"},
		{"one hour", 3600, "1h 0m"},
		{"hours and minutes", 5100, "1h 25m"},
		{"multiple hours", 9000, "2h 30m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.seconds)
			if result != tt.expected {
				t.Errorf("formatDuration(%d) = %q, want %q", tt.seconds, result, tt.expected)
			}
		})
	}
}
