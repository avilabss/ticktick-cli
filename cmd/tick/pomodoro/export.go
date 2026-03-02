package pomodoro

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/avilabss/ticktick-cli/pkg/logger"
	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

var csvHeader = []string{"Date", "Week", "Start Time", "End Time", "Duration", "Tags", "Description"}

func parseExportArgs(args []string) exportArgs {
	now := time.Now()
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	year := fs.Int("year", now.Year(), "year to fetch pomodoros for")
	month := fs.Int("month", int(now.Month()), "month to fetch pomodoros for (1-12)")
	includeTags := fs.String("include-tags", "", "comma-separated tags to include")
	excludeTags := fs.String("exclude-tags", "", "comma-separated tags to exclude")
	includeProjects := fs.String("include-projects", "", "comma-separated project names to include")
	excludeProjects := fs.String("exclude-projects", "", "comma-separated project names to exclude")
	output := fs.String("output", "", "output CSV file path (default: pomodoros-YYYY-MM.csv)")
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if *output == "" {
		*output = fmt.Sprintf("pomodoros-%04d-%02d.csv", *year, *month)
	}

	ea := exportArgs{
		Year:            *year,
		Month:           *month,
		IncludeTags:     splitCSV(*includeTags),
		ExcludeTags:     splitCSV(*excludeTags),
		IncludeProjects: splitCSV(*includeProjects),
		ExcludeProjects: splitCSV(*excludeProjects),
		Output:          *output,
	}

	slog.Debug("Parsed export args",
		"year", ea.Year,
		"month", ea.Month,
		"output", ea.Output,
		"includeTags", ea.IncludeTags,
		"excludeTags", ea.ExcludeTags,
		"includeProjects", ea.IncludeProjects,
		"excludeProjects", ea.ExcludeProjects,
	)

	return ea
}

func exportCSV(pomodoros []ticktick.Pomodoro, args exportArgs, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(csvHeader); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	slog.Debug("Processing pomodoros for export", "count", len(pomodoros))

	rowCount := 0
	skippedProject := 0
	for _, p := range pomodoros {
		for _, t := range p.Tasks {
			if !matchesFilter(t.ProjectName, args.IncludeProjects, args.ExcludeProjects) {
				skippedProject++
				logger.Trace("Skipping task (project filtered)",
					"title", t.Title,
					"project", t.ProjectName,
				)
				continue
			}

			startTime, err := time.Parse(ticktick.TimeFormat, t.StartTime)
			if err != nil {
				slog.Warn("Skipping entry: error parsing start time", "error", err)
				continue
			}

			endTime, err := time.Parse(ticktick.TimeFormat, t.EndTime)
			if err != nil {
				slog.Warn("Skipping entry: error parsing end time", "error", err)
				continue
			}

			tags := includeExclude(t.Tags, args.IncludeTags, args.ExcludeTags)
			duration := endTime.Sub(startTime)

			logger.Trace("Writing row",
				"date", startTime.Format(dateFormat),
				"start", startTime.Format(timeFormat),
				"end", endTime.Format(timeFormat),
				"duration", duration,
				"tags", tags,
				"title", t.Title,
			)

			row := []string{
				startTime.Format(dateFormat),
				fmt.Sprintf("%d", (startTime.Day()-1)/7+1),
				startTime.Format(timeFormat),
				endTime.Format(timeFormat),
				duration.String(),
				strings.Join(tags, ", "),
				t.Title,
			}

			if err := writer.Write(row); err != nil {
				return fmt.Errorf("failed to write row: %w", err)
			}
			rowCount++
		}
	}

	slog.Debug("Export filtering complete", "exported", rowCount, "skippedByProject", skippedProject)
	slog.Info("Exported to CSV", "rows", rowCount, "file", filename)
	return nil
}

func runExport(client *ticktick.Client, args []string) {
	ea := parseExportArgs(args)

	start, end := monthRange(ea.Year, time.Month(ea.Month))
	slog.Info("Fetching pomodoros", "from", start.Format(time.RFC3339), "to", end.Format(time.RFC3339))

	result, err := client.Pomodoro.GetAll(start, end)
	if err != nil {
		slog.Error("Failed to fetch pomodoros", "error", err)
		os.Exit(1)
	}
	slog.Info("Fetched pomodoros", "count", len(result.Items))

	slices.Reverse(result.Items)

	if err := exportCSV(result.Items, ea, ea.Output); err != nil {
		slog.Error("Failed to export CSV", "error", err)
		os.Exit(1)
	}
}
