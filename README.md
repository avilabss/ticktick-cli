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
tick --help
```

## Setup

Create a `.env` file in the project root:

```
TICKTICK_API_TOKEN=your_token_here
```

## Usage

```sh
tick <command> <subcommand> [flags]
```

## Commands

### Tasks

```sh
tick task list [--project NAME] [--tag TAG] [--priority N]
tick task add --title "..." [--project NAME] [--tags "a,b"] [--priority N] [--due DATE]
tick task complete ID [--project NAME]
tick task delete ID [--project NAME]
tick task get ID [--project NAME]
tick project list
```

### Pomodoro

```sh
tick pomodoro export [--year N] [--month N] [--output FILE] [--include-tags ...] [--exclude-tags ...]
tick pomodoro stats
tick pomodoro create [--task ID] [--duration MINS]
tick pomodoro delete ID
tick pomodoro timer list
tick pomodoro timer stats NAME
```

### Habits

```sh
tick habit list
tick habit checkin HABIT_NAME [--value N]
tick habit status [--date YYYY-MM-DD]
```

### Global Flags

```sh
tick -v ...     # Info level logging
tick -vv ...    # Debug level logging
tick -vvv ...   # Trace level logging
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
