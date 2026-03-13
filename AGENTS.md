# AGENTS.md

## Project Overview

ticktick-cli (`tick`) is a CLI tool for interacting with the TickTick API. It uses cookie-based authentication (not OAuth). The API token is loaded from a `.env` file via `godotenv`.

The binary is invoked as `tick` with subcommands: `tick <service> <action> [flags]`.

Implemented commands:
- **Tasks**: `tick task list|add|complete|delete|get`, `tick project list`
- **Pomodoro**: `tick pomodoro export|stats|create|delete`, `tick pomodoro timer list|stats`
- **Habits**: `tick habit list|checkin|status`

## Tech Stack

- Go 1.25.6
- CLI framework: `github.com/spf13/cobra`
- Logging: `log/slog` (stdlib) with a thin wrapper in `internal/logger` for a custom Trace level
- Dependencies: `github.com/spf13/cobra`, `github.com/joho/godotenv`
- Task runner: `just` (Justfile)
- CI: GitHub Actions (lint, test, build)

## Project Structure

```
cmd/tick/
  main.go                       — entrypoint: rootCmd.Execute()
  root.go                       — root cobra command, PersistentPreRunE (.env, client, verbosity), wires subcommands
  integration_test.go           — CLI integration tests (build tag: integration)
  task/
    task.go                     — NewCmd/NewProjectCmd: returns cobra command groups
    list.go                     — "task list": --project, --tag, --priority filters, tabwriter output
    add.go                      — "task add": --title, --project, --tags, --priority, --due
    complete.go                 — "task complete ID": optional --project
    delete.go                   — "task delete ID": optional --project
    get.go                      — "task get ID": optional --project, detailed output
    project.go                  — "project list": cobra subcommand, tabwriter output
  pomodoro/
    pomodoro.go                 — NewCmd: returns cobra command group
    export.go                   — "pomodoro export": CSV export with filters
    stats.go                    — "pomodoro stats": today/total counts and durations
    create.go                   — "pomodoro create": --task, --duration
    delete.go                   — "pomodoro delete ID"
    timer.go                    — "pomodoro timer list|stats": nested cobra commands
    types.go                    — exportArgs struct, format constants
    utils.go                    — helpers: splitCSV, includeExclude, matchesFilter, monthRange, formatDuration
    utils_test.go               — unit tests for utils
  habit/
    habit.go                    — NewCmd: returns cobra command group
    list.go                     — "habit list": tabwriter output
    checkin.go                  — "habit checkin NAME": --value for Number types
    status.go                   — "habit status": --date, shows done/pending per habit

internal/ticktick/
  client.go                     — Client constructor, HTTP helpers (Get, Post, Put, Delete, PostJSON, etc.)
  client_test.go                — unit tests for client HTTP methods
  task.go                       — TaskService: Sync, List, Get, Create, Complete, Delete, ListProjects
  task_test.go                  — unit tests for task service
  pomodoro.go                   — PomodoroService: GetTimeline, GetAll, Stats, Create, DeletePomo, ListTimers, TimerOverview
  pomodoro_test.go              — unit tests for pomodoro service
  habit.go                      — HabitService: List, GetCheckins, Checkin, GetRecords
  habit_test.go                 — unit tests for habit service
  integration_test.go           — API integration tests (build tag: integration)
  types.go                      — All types: Client, Task, Project, Habit, Pomodoro, BatchResponse, etc.

internal/logger/
  logger.go                     — slog setup: SetVerbosity maps -v/-vv/-vvv to slog levels, custom Trace level

.github/workflows/
  ci.yml                        — CI pipeline: lint, test, build
  release.yml                   — Auto-release on push to main: test → auto-tag (patch bump) → cross-platform build → GitHub release
```

## Architecture & Patterns

### Client + Sub-services
`ticktick.Client` has service fields (`client.Task`, `client.Pomodoro`, `client.Habit`) following the go-github pattern. Each service is a separate struct with a back-reference to the client for making HTTP calls. New services should follow this pattern:
1. Create the service struct in `internal/ticktick/<service>.go`
2. Add the field to `Client` in `types.go`
3. Wire it in `NewTicktickClient` in `client.go`
4. Create `cmd/tick/<service>/` package for CLI commands

### Functional Options
`NewTicktickClient(apiToken, ...Option)` uses functional options (`WithHTTPClient`).

### HTTPClient Interface
`Client.HTTPClient` is typed as `HTTPClient` interface (just `Do(*http.Request) (*http.Response, error)`). `*http.Client` satisfies it. Tests inject a `mockHTTPClient` with a `DoFunc` field.

### Pagination
`Pomodoros` wraps `[]Pomodoro` with a service reference. Callers can use `.Next()` on a result set for manual pagination, or `GetAll()` which paginates automatically. The API returns results newest-first.

### Subcommand Routing (Cobra)
`root.go` defines the root cobra command with `PersistentPreRunE` for .env loading, client creation, and verbosity setup. Each service package exports a `NewCmd(client **ticktick.Client) *cobra.Command` that returns a command group with subcommands. The `**ticktick.Client` (pointer-to-pointer) pattern lets root's `PersistentPreRunE` initialize the client before any subcommand `RunE` executes. All commands return errors instead of calling `os.Exit(1)` directly.

### Logging Levels
- 0 (default): Warn+ only (errors, warnings)
- 1 (`-v`): Info (progress, counts, results)
- 2 (`-vv`): Debug (parsed args, pagination details, API responses)
- 3 (`-vvv`): Trace (raw HTTP requests, individual row data, internal state)

Use `slog.Info`, `slog.Debug` directly. Use `logger.Trace` for the custom trace level.

### TickTick API Notes
- Auth: cookie-based, header `Cookie: t=<token>`. Token is the `t` cookie from ticktick.com.
- Required headers: `Origin: https://ticktick.com` and `x-device` (JSON) — POST requests fail with `access_forbidden` without these
- Base URL: `https://api.ticktick.com/api`
- Timestamps: format `2006-01-02T15:04:05.000-0700` (NOT RFC3339 — no colon in timezone offset)
- Timeline API uses millisecond unix timestamps for the `to` parameter
- API returns ~31 results per page, newest first
- Batch endpoints use `{"add":[], "update":[], "delete":[]}` → `{"id2etag":{}, "id2error":{}}`
- Sync endpoint (`GET /v3/batch/check/0`) returns all tasks, projects, tags in one call
- IDs are 24-char hex strings generated with `crypto/rand`
- Habit checkin stamps are `YYYYMMDD` integers (e.g. `20260301`)

## Testing

### Unit tests
- Use `mockHTTPClient` struct with `DoFunc` field to mock HTTP calls
- Run: `just test` or `go test ./...`
- Test files live alongside source files (`*_test.go`)

### Integration tests
- Guarded by `//go:build integration` build tag
- Require `TICKTICK_API_TOKEN` env var (loaded from `.env`)
- Run: `just test-integration` or `go test -tags integration -v ./...`
- NOT run in CI — manual only

### CI Pipeline
- Triggers on push to `main` and PRs to `main`
- Jobs: `lint` (golangci-lint), `test` (unit tests), `build` (compilation check)
- Build job depends on lint + test passing

## Common Commands (Justfile)

- `just test` — run unit tests
- `just test-v` — run unit tests (verbose)
- `just test-integration` — run integration tests
- `just lint` — run linter
- `just build` — verify compilation
- `just run <args>` — run the CLI

## Rules

1. **Do not build after changes.** If you do build, remove the generated binary immediately. Do not leave binaries in the repo.
2. **Update AGENTS.md** when you add new services, change architecture, modify patterns, or add dependencies.
3. **Update README.md** when you add new commands, flags, or change usage.
4. Keep Go conventions: unexported names for internal-only identifiers, doc comments starting with the function name, no `SCREAMING_SNAKE_CASE`.
5. All filtering (tags, projects) is case-sensitive and exact match.
6. New services should follow the existing package structure: `internal/ticktick/<service>.go` for API logic, `cmd/tick/<service>/` for CLI commands.
7. Use Cobra for CLI commands. Each service package exports `NewCmd(client **ticktick.Client) *cobra.Command`. Use `cobra.ExactArgs`, `cobra.NoArgs`, etc. for arg validation. Define flags on the command, not via `flag.NewFlagSet`.
8. Never pass `nil` as a context. Use `context.TODO()` when no context is available.
9. Write unit tests for new functionality. Use the `mockHTTPClient` pattern for HTTP mocking.
10. Integration tests must use `//go:build integration` build tag.
