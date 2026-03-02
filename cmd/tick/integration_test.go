//go:build integration

package main

import (
	"encoding/csv"
	"os"
	"os/exec"
	"testing"
)

func TestIntegration_CLIExport(t *testing.T) {
	outputFile := "test-export.csv"
	defer os.Remove(outputFile)

	cmd := exec.Command("go", "run", "./", "pomodoro", "export",
		"--year", "2026",
		"--month", "2",
		"--output", outputFile,
	)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed: %v\nOutput: %s", err, out)
	}

	// Verify CSV file exists
	file, err := os.Open(outputFile)
	if err != nil {
		t.Fatalf("failed to open output file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	// Verify header
	if len(records) == 0 {
		t.Fatal("CSV file is empty")
	}
	expectedHeader := []string{"Date", "Week", "Start Time", "End Time", "Duration", "Tags", "Description"}
	for i, col := range expectedHeader {
		if records[0][i] != col {
			t.Errorf("header column %d: expected %q, got %q", i, col, records[0][i])
		}
	}

	t.Logf("CSV has %d rows (including header)", len(records))
}

func TestIntegration_CLIExportWithFilters(t *testing.T) {
	outputFile := "test-export-filtered.csv"
	defer os.Remove(outputFile)

	cmd := exec.Command("go", "run", "./", "pomodoro", "export",
		"--year", "2026",
		"--month", "2",
		"--exclude-tags", "break",
		"--output", outputFile,
	)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed: %v\nOutput: %s", err, out)
	}

	file, err := os.Open(outputFile)
	if err != nil {
		t.Fatalf("failed to open output file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	t.Logf("Filtered CSV has %d rows (including header)", len(records))
}

func TestIntegration_CLINoArgs(t *testing.T) {
	cmd := exec.Command("go", "run", "./")
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit code when no args provided")
	}

	if len(out) == 0 {
		t.Error("expected usage output")
	}
}
