package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

var csvHeader = []string{"Date", "Week", "Start Time", "End Time", "Duration", "Tags", "Description"}

func exportCSV(pomodoros []ticktick.Pomodoro, projectName string, filterTags []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(csvHeader); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	rowCount := 0
	for _, p := range pomodoros {
		for _, t := range p.Tasks {
			if projectName != "" && !strings.Contains(strings.ToLower(t.ProjectName), strings.ToLower(projectName)) {
				continue
			}

			startTime, err := time.Parse(ticktick.TimeFormat, t.StartTime)
			if err != nil {
				log.Printf("Skipping entry: error parsing start time: %v", err)
				continue
			}

			endTime, err := time.Parse(ticktick.TimeFormat, t.EndTime)
			if err != nil {
				log.Printf("Skipping entry: error parsing end time: %v", err)
				continue
			}

			row := []string{
				startTime.Format(DateFormat),
				fmt.Sprintf("%d", (startTime.Day()-1)/7+1),
				startTime.Format(TimeFormat),
				endTime.Format(TimeFormat),
				endTime.Sub(startTime).String(),
				strings.Join(filter(t.Tags, filterTags), ", "),
				t.Title,
			}

			if err := writer.Write(row); err != nil {
				return fmt.Errorf("failed to write row: %w", err)
			}
			rowCount++
		}
	}

	log.Printf("Exported %d rows to %s", rowCount, filename)
	return nil
}
