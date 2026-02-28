package main

import (
	"flag"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
	"github.com/joho/godotenv"
)

func printPomodoros(pomodoros []ticktick.Pomodoro, projectName string, filterTags []string) {
	for _, p := range pomodoros {
		for _, t := range p.Tasks {
			if projectName != "" && !strings.Contains(strings.ToLower(t.ProjectName), strings.ToLower(projectName)) {
				continue
			}

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
			tags := filter(t.Tags, filterTags)
			workType := strings.Join(tags, ", ")
			workDescription := t.Title

			log.Printf("\t%s | %d | %s | %s | %s | %s | %s\n", dateStr, weekNumOfMonth, startTimeStr, endTimeStr, duration, workType, workDescription)
		}
	}
}

func main() {
	now := time.Now()
	year := flag.Int("year", now.Year(), "year to fetch pomodoros for")
	month := flag.Int("month", int(now.Month()), "month to fetch pomodoros for (1-12)")
	filterTagsStr := flag.String("filter-tags", "", "comma-separated tags to remove from output")
	projectName := flag.String("project-name", "", "filter by project name (case-insensitive)")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	apiToken := os.Getenv("TICKTICK_API_TOKEN")
	if apiToken == "" {
		log.Fatal("TICKTICK_API_TOKEN is required")
	}

	var filterTags []string
	if *filterTagsStr != "" {
		filterTags = strings.Split(*filterTagsStr, ",")
	}

	client, err := ticktick.NewTicktickClient(apiToken)
	if err != nil {
		log.Fatalf("Error creating TickTick client: %v", err)
	}

	start, end := monthRange(*year, time.Month(*month))
	log.Printf("Fetching pomodoros from %s to %s\n", start.Format(time.RFC3339), end.Format(time.RFC3339))

	pomodoros, err := client.GetAllPomodorosTimeline(start, end)
	if err != nil {
		log.Fatalf("Error fetching all pomodoros timeline: %v", err)
	}

	slices.Reverse(pomodoros)
	printPomodoros(pomodoros, *projectName, filterTags)
}
