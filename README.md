# ticktick-cli

CLI tools for interacting with the TickTick API.

## Setup

1. Create a `.env` file in the project root:

```
TICKTICK_API_TOKEN=your_token_here
```

2. Install dependencies:

```sh
go mod tidy
```

## pomo-exporter

Exports pomodoro timeline data from TickTick to CSV.

### Usage

```sh
go run ./cmd/pomo-exporter [flags]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--year` | current year | Year to fetch pomodoros for |
| `--month` | current month | Month to fetch pomodoros for (1-12) |
| `--filter-tags` | _(none)_ | Comma-separated tags to remove from output |
| `--project-name` | _(none)_ | Filter by project name (case-insensitive, partial match) |
| `--output` | `pomodoros-YYYY-MM.csv` | Output CSV file path |

### Examples

```sh
# Fetch current month's pomodoros
go run ./cmd/pomo-exporter

# Fetch January 2026
go run ./cmd/pomo-exporter --year 2026 --month 1

# Filter by project and remove specific tags from output
go run ./cmd/pomo-exporter --project-name "Whitebox" --filter-tags "freelancing,whitebox"

# Custom output path
go run ./cmd/pomo-exporter --output report.csv
```
