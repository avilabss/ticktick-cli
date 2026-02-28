package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
	"github.com/joho/godotenv"
)

func parseArgs() Args {
	now := time.Now()
	year := flag.Int("year", now.Year(), "year to fetch pomodoros for")
	month := flag.Int("month", int(now.Month()), "month to fetch pomodoros for (1-12)")
	filterTagsStr := flag.String("filter-tags", "", "comma-separated tags to remove from output")
	projectName := flag.String("project-name", "", "filter by project name (case-insensitive)")
	output := flag.String("output", "", "output CSV file path (default: pomodoros-YYYY-MM.csv)")
	flag.Parse()

	var filterTags []string
	if *filterTagsStr != "" {
		filterTags = strings.Split(*filterTagsStr, ",")
	}

	if *output == "" {
		*output = fmt.Sprintf("pomodoros-%04d-%02d.csv", *year, *month)
	}

	return Args{
		Year:        *year,
		Month:       *month,
		FilterTags:  filterTags,
		ProjectName: *projectName,
		Output:      *output,
	}
}

func main() {
	args := parseArgs()

	if err := godotenv.Load(); err != nil {
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

	start, end := monthRange(args.Year, time.Month(args.Month))
	log.Printf("Fetching pomodoros from %s to %s", start.Format(time.RFC3339), end.Format(time.RFC3339))

	result, err := client.Pomodoro.GetAll(start, end)
	if err != nil {
		log.Fatalf("Error fetching pomodoros: %v", err)
	}
	log.Printf("Fetched %d pomodoros", len(result.Items))

	slices.Reverse(result.Items)

	if err := exportCSV(result.Items, args.ProjectName, args.FilterTags, args.Output); err != nil {
		log.Fatalf("Error exporting CSV: %v", err)
	}
}
