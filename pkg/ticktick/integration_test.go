//go:build integration

package ticktick

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func integrationClient(t *testing.T) *Client {
	t.Helper()
	godotenv.Load("../../.env")

	token := os.Getenv("TICKTICK_API_TOKEN")
	if token == "" {
		t.Skip("TICKTICK_API_TOKEN not set, skipping integration test")
	}

	client, err := NewTicktickClient(token)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

func TestIntegration_GetTimeline(t *testing.T) {
	client := integrationClient(t)

	result, err := client.Pomodoro.GetTimeline(0)
	if err != nil {
		t.Fatalf("GetTimeline failed: %v", err)
	}

	t.Logf("Got %d pomodoros from latest timeline", len(result.Items))

	if len(result.Items) == 0 {
		t.Log("Warning: no pomodoros returned, account may be empty")
	}
}

func TestIntegration_GetAll(t *testing.T) {
	client := integrationClient(t)

	start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC)

	result, err := client.Pomodoro.GetAll(start, end)
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	t.Logf("Got %d pomodoros for Feb 2026", len(result.Items))

	// Verify all items are within range
	for _, p := range result.Items {
		pTime, err := time.Parse(TimeFormat, p.StartTime)
		if err != nil {
			t.Errorf("failed to parse start time %q: %v", p.StartTime, err)
			continue
		}
		if pTime.Before(start) {
			t.Errorf("pomodoro %s has start time %s before range start %s", p.ID, p.StartTime, start)
		}
	}
}

func TestIntegration_Next(t *testing.T) {
	client := integrationClient(t)

	first, err := client.Pomodoro.GetTimeline(0)
	if err != nil {
		t.Fatalf("GetTimeline failed: %v", err)
	}

	if len(first.Items) == 0 {
		t.Skip("No pomodoros to paginate from")
	}

	next, err := first.Next()
	if err != nil {
		t.Fatalf("Next failed: %v", err)
	}

	t.Logf("First page: %d items, Next page: %d items", len(first.Items), len(next.Items))

	// Next page should have different items (if there are any)
	if len(next.Items) > 0 && next.Items[0].ID == first.Items[0].ID {
		t.Log("Warning: next page returned same first item, may indicate end of data")
	}
}
