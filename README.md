# ticktick-cli

CLI tool for interacting with the TickTick API.

## Setup

1. Create a `.env` file in the project root:

```
TICKTICK_API_TOKEN=your_token_here
```

2. Install dependencies:

```sh
go mod tidy
```

## Usage

```sh
go run ./cmd/tt <command> <subcommand> [flags]
```

## Commands

### pomodoro export

Export pomodoro timeline data to CSV.

```sh
go run ./cmd/tt pomodoro export [flags]
```

#### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--year` | current year | Year to fetch pomodoros for |
| `--month` | current month | Month to fetch pomodoros for (1-12) |
| `--include-tags` | _(none)_ | Comma-separated tags to include |
| `--exclude-tags` | _(none)_ | Comma-separated tags to exclude |
| `--include-projects` | _(none)_ | Comma-separated project names to include |
| `--exclude-projects` | _(none)_ | Comma-separated project names to exclude |
| `--output` | `pomodoros-YYYY-MM.csv` | Output CSV file path |

#### Examples

```sh
# Export current month's pomodoros
go run ./cmd/tt pomodoro export

# Export January 2026
go run ./cmd/tt pomodoro export --year 2026 --month 1

# Include only specific project
go run ./cmd/tt pomodoro export --include-projects "Whitebox"

# Exclude specific tags
go run ./cmd/tt pomodoro export --exclude-tags "freelancing,whitebox"

# Combine filters
go run ./cmd/tt pomodoro export --include-projects "Whitebox" --exclude-tags "freelancing"

# Custom output path
go run ./cmd/tt pomodoro export --output report.csv
```
