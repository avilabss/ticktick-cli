# AGENTS.md

## Project Overview

ticktick-cli (`tt`) is a CLI tool for interacting with the TickTick API. It uses cookie-based authentication (not OAuth). The API token is loaded from a `.env` file via `godotenv`.

The binary is invoked as `tt` with subcommands: `tt <service> <action> [flags]`.

Currently implemented: `tt pomodoro export` — exports pomodoro timeline data to CSV with filtering options.

## Tech Stack

- Go 1.25.6
- No CLI framework — uses stdlib `flag` with manual subcommand routing
- Logging: `log/slog` (stdlib) with a thin wrapper in `pkg/logger` for a custom Trace level
- Dependencies: `github.com/joho/godotenv` (only external dep)
- Task runner: `just` (Justfile)
- CI: GitHub Actions (lint, test, build)

## Project Structure

```
cmd/tt/
  main.go                       — entrypoint: .env loading, client creation, verbosity parsing, subcommand routing
  integration_test.go           — CLI integration tests (build tag: integration)
  pomodoro/
    pomodoro.go                 — "pomodoro" subcommand routing (dispatches to export, future commands)
    export.go                   — "pomodoro export": flag parsing, CSV export logic, runExport()
    types.go                    — exportArgs struct, format constants (dateFormat, timeFormat)
    utils.go                    — helpers: splitCSV, includeExclude, matchesFilter, monthRange
    utils_test.go               — unit tests for utils

pkg/ticktick/
  client.go                     — Client constructor (functional options pattern), HTTP Get helper
  client_test.go                — unit tests for client (mock HTTPClient)
  pomodoro.go                   — PomodoroService: GetTimeline, GetAll, Next (pagination)
  pomodoro_test.go              — unit tests for pomodoro service
  integration_test.go           — API integration tests (build tag: integration)
  types.go                      — Client, HTTPClient interface, Option, Pomodoro, PomodoroTask, Pomodoros types

pkg/logger/
  logger.go                     — slog setup: SetVerbosity maps -v/-vv/-vvv to slog levels, custom Trace level

.github/workflows/
  ci.yml                        — CI pipeline: lint, test, build
```

## Architecture & Patterns

### Client + Sub-services
`ticktick.Client` has service fields (e.g. `client.Pomodoro`) following the go-github pattern. Each service is a separate struct with a back-reference to the client for making HTTP calls. New services (habits, tasks) should follow this pattern:
1. Create the service struct in `pkg/ticktick/<service>.go`
2. Add the field to `Client` in `types.go`
3. Wire it in `NewTicktickClient` in `client.go`
4. Create `cmd/tt/<service>/` package for CLI commands

### Functional Options
`NewTicktickClient(apiToken, ...Option)` uses functional options (`WithHTTPClient`).

### HTTPClient Interface
`Client.HTTPClient` is typed as `HTTPClient` interface (just `Do(*http.Request) (*http.Response, error)`). `*http.Client` satisfies it. Tests inject a `mockHTTPClient` with a `DoFunc` field.

### Pagination
`Pomodoros` wraps `[]Pomodoro` with a service reference. Callers can use `.Next()` on a result set for manual pagination, or `GetAll()` which paginates automatically. The API returns results newest-first.

### Subcommand Routing
`main.go` parses global flags (`-v`), then dispatches `os.Args` to service packages. Each service package (e.g. `cmd/tt/pomodoro/`) handles its own subcommand routing and flag parsing via `flag.NewFlagSet`.

### Logging Levels
- 0 (default): Warn+ only (errors, warnings)
- 1 (`-v`): Info (progress, counts, results)
- 2 (`-vv`): Debug (parsed args, pagination details, API responses)
- 3 (`-vvv`): Trace (raw HTTP requests, individual row data, internal state)

Use `slog.Info`, `slog.Debug` directly. Use `logger.Trace` for the custom trace level.

### TickTick API Notes
- Auth: cookie-based, header `Cookie: t=<token>`
- Base URL: `https://api.ticktick.com/api`
- Timestamps: format `2006-01-02T15:04:05.000-0700` (NOT RFC3339 — no colon in timezone offset)
- Timeline API uses millisecond unix timestamps for the `to` parameter
- API returns ~31 results per page, newest first

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
6. New services should follow the existing package structure: `pkg/ticktick/<service>.go` for API logic, `cmd/tt/<service>/` for CLI commands.
7. No external CLI frameworks — use stdlib `flag` with `flag.NewFlagSet` for subcommand flag parsing.
8. Never pass `nil` as a context. Use `context.TODO()` when no context is available.
9. Write unit tests for new functionality. Use the `mockHTTPClient` pattern for HTTP mocking.
10. Integration tests must use `//go:build integration` build tag.
