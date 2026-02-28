package main

import (
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
	"github.com/joho/godotenv"
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

func PrintPomodoros(pomodoros []ticktick.Pomodoro) {
	for _, p := range pomodoros {
		for _, t := range p.Tasks {
			startTime, err := time.Parse(ticktick.TimeFormat, t.StartTime)
			if err != nil {
				log.Printf("Error parsing start time: %v", err)
				continue
			}

			endTime, err := time.Parse(ticktick.TimeFormat, t.EndTime)
			if err != nil {
				log.Printf("Error parsing end time: %v", err)
				continue
			}

			dateStr := startTime.Format(DateFormat)
			weekNumOfMonth := (startTime.Day()-1)/7 + 1
			startTimeStr := startTime.Format(TimeFormat)
			endTimeStr := endTime.Format(TimeFormat)
			duration := endTime.Sub(startTime)
			tags := filter(t.Tags, []string{"freelancing", "whitebox"})
			workType := strings.Join(tags, ", ")
			workDescription := t.Title

			log.Printf("\t%s | %d | %s | %s | %s | %s | %s\n", dateStr, weekNumOfMonth, startTimeStr, endTimeStr, duration, workType, workDescription)
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	apiToken := os.Getenv("TICKTICK_API_TOKEN")
	if apiToken == "" {
		log.Fatal("TICKTICK_API_TOKEN is required")
	}

	client, err := ticktick.NewTicktickClient(apiToken)
	if err != nil {
		log.Fatalf("Error creating TickTick client: %v", err)
	}

	start, end := monthRange(2026, time.February)
	log.Printf("Fetching pomodoros from %s to %s\n", start.Format(time.RFC3339), end.Format(time.RFC3339))

	pomodoros, err := client.GetAllPomodorosTimeline(start, end)
	if err != nil {
		log.Fatalf("Error fetching all pomodoros timeline: %v", err)
	}

	slices.Reverse(pomodoros)
	PrintPomodoros(pomodoros)
}
