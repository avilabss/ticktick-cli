package pomodoro

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/avilabss/ticktick-cli/internal/logger"
	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

var csvHeader = []string{"Date", "Week", "Start Time", "End Time", "Duration", "Tags", "Description"}

func exportCmd(client **ticktick.Client) *cobra.Command {
	now := time.Now()
	var ea exportArgs
	var includeTags, excludeTags, includeProjects, excludeProjects string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export pomodoros to CSV",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ea.IncludeTags = splitCSV(includeTags)
			ea.ExcludeTags = splitCSV(excludeTags)
			ea.IncludeProjects = splitCSV(includeProjects)
			ea.ExcludeProjects = splitCSV(excludeProjects)

			if ea.Output == "" {
				ea.Output = fmt.Sprintf("pomodoros-%04d-%02d.csv", ea.Year, ea.Month)
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

			return runExport(*client, ea)
		},
	}
	cmd.Flags().IntVar(&ea.Year, "year", now.Year(), "year to fetch pomodoros for")
	cmd.Flags().IntVar(&ea.Month, "month", int(now.Month()), "month to fetch pomodoros for (1-12)")
	cmd.Flags().StringVar(&includeTags, "include-tags", "", "comma-separated tags to include")
	cmd.Flags().StringVar(&excludeTags, "exclude-tags", "", "comma-separated tags to exclude")
	cmd.Flags().StringVar(&includeProjects, "include-projects", "", "comma-separated project names to include")
	cmd.Flags().StringVar(&excludeProjects, "exclude-projects", "", "comma-separated project names to exclude")
	cmd.Flags().StringVar(&ea.Output, "output", "", "output CSV file path (default: pomodoros-YYYY-MM.csv)")
	return cmd
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

func runExport(client *ticktick.Client, ea exportArgs) error {
	start, end := monthRange(ea.Year, time.Month(ea.Month))
	slog.Info("Fetching pomodoros", "from", start.Format(time.RFC3339), "to", end.Format(time.RFC3339))

	result, err := client.Pomodoro.GetAll(start, end)
	if err != nil {
		slog.Error("Failed to fetch pomodoros", "error", err)
		return err
	}
	slog.Info("Fetched pomodoros", "count", len(result.Items))

	slices.Reverse(result.Items)

	if err := exportCSV(result.Items, ea, ea.Output); err != nil {
		slog.Error("Failed to export CSV", "error", err)
		return err
	}
	return nil
}
