# ticktick-cli

CLI tool for interacting with the TickTick API.

## Installation

### From GitHub Releases

Download the latest binary for your platform from the [releases page](https://github.com/avilabss/ticktick-cli/releases).

**macOS (Apple Silicon)**

```sh
curl -L -o tick https://github.com/avilabss/ticktick-cli/releases/latest/download/tick-darwin-arm64
chmod +x tick
sudo mv tick /usr/local/bin/
```

**macOS (Intel)**

```sh
curl -L -o tick https://github.com/avilabss/ticktick-cli/releases/latest/download/tick-darwin-amd64
chmod +x tick
sudo mv tick /usr/local/bin/
```

**Linux (x86_64)**

```sh
curl -L -o tick https://github.com/avilabss/ticktick-cli/releases/latest/download/tick-linux-amd64
chmod +x tick
sudo mv tick /usr/local/bin/
```

**Linux (ARM64)**

```sh
curl -L -o tick https://github.com/avilabss/ticktick-cli/releases/latest/download/tick-linux-arm64
chmod +x tick
sudo mv tick /usr/local/bin/
```

**Windows**

Download `tick-windows-amd64.exe` from the [releases page](https://github.com/avilabss/ticktick-cli/releases), rename to `tick.exe`, and add its location to your `PATH`.

### From source

```sh
go install github.com/avilabss/ticktick-cli/cmd/tick@latest
```

### Verify

```sh
tick pomodoro export --help
```

## Setup

Create a `.env` file in the project root:

```
TICKTICK_API_TOKEN=your_token_here
```

## Usage

```sh
tick <command> <subcommand> [flags]

# Or use just
just run pomodoro export
```

## Commands

### pomodoro export

Export pomodoro timeline data to CSV.

```sh
tick pomodoro export [flags]
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
tick pomodoro export

# Export January 2026
tick pomodoro export --year 2026 --month 1

# Include only specific project
tick pomodoro export --include-projects "TickTick"

# Exclude specific tags
tick pomodoro export --exclude-tags "fun"

# Combine filters
tick pomodoro export --include-projects "TickTick" --exclude-tags "fun"

# Custom output path
tick pomodoro export --output report.csv
```

## Development

### Testing

```sh
just test                # unit tests
just test-v              # unit tests (verbose)
just test-integration    # integration tests (requires .env)
```

### Linting

```sh
just lint
```

### All Justfile commands

```sh
just --list
```

### Releasing

Every push to `main` automatically:
1. Runs tests
2. Increments the patch version (e.g. v0.0.1 → v0.0.2)
3. Builds binaries for all platforms (linux/darwin/windows, amd64/arm64)
4. Publishes a GitHub release with SHA256 checksums

To bump minor or major version, create the tag manually before the next push:

```sh
git tag v0.1.0
git push origin v0.1.0
```

The next auto-release will increment from that tag (v0.1.0 → v0.1.1).
