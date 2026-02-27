package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
	"github.com/joho/godotenv"
)

func PrintPomodoros(pomodoros []ticktick.Pomodoro) {
	for pi, p := range pomodoros {
		fmt.Printf("Pomodoro %d: %s\n", pi, p.ID)

		for ti, t := range p.Tasks {
			fmt.Printf("\tTask %d: %s\n", ti, t.TaskID)
			fmt.Printf("\t\tTitle: %s\n", t.Title)
			fmt.Printf("\t\tProject: %s\n", t.ProjectName)
			fmt.Printf("\t\tTags: %v\n", strings.Join(t.Tags, ", "))
			fmt.Printf("\t\tStart Time: %s\n", t.StartTime)
			fmt.Printf("\t\tEnd Time: %s\n", t.EndTime)
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

	var allPomodoros []ticktick.Pomodoro

	pomodoros, err := client.GetPomodorosTimeline(0)
	if err != nil {
		log.Fatalf("Error fetching pomodoros timeline: %v", err)
	}

	log.Printf("Fetched %d pomodoros\n", len(pomodoros))
	allPomodoros = append(allPomodoros, pomodoros...)

	nextPomodoros, err := client.GetNextPomodorosTimeline(pomodoros)
	if err != nil {
		log.Fatalf("Error fetching next pomodoros timeline: %v", err)
	}

	log.Printf("Fetched %d next pomodoros\n", len(nextPomodoros))
	allPomodoros = append(allPomodoros, nextPomodoros...)

	PrintPomodoros(allPomodoros)
}
