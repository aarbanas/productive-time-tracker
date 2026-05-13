# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build -o productive-time-tracker

# Run (previous month by default, -c for current month)
./productive-time-tracker
./productive-time-tracker -c

# Test all
go test ./...

# Run a single test
go test ./utilities/ -run TestPreviousMonthBounds
```

## Architecture

This is a Go CLI tool that queries the [Productive.io](https://productive.io) API to report whether the user has correctly tracked their hours for a given month.

**Flow:** `main.go` loads credentials → creates an API client → calls `utilities.ReportMinutes` → prints whether hours are missing, over, or correct.

### Packages

- **`api/`** — HTTP client for the Productive.io Reports API (`/reports/time_reports`). Handles auth headers (`X-Auth-Token`, `X-Organization-Id`) and JSON:API response parsing. Only supports GET requests via `Client.GetTimeReport(after, before string)`.

- **`appinit/`** — Credential bootstrapping. Reads `~/.productive_token` and `~/.productive_org_id` from disk; prompts the user and saves them (mode 0600) if not found.

- **`utilities/`** — Business logic. `ReportMinutes` sums `worked_time + event_time` (absence) across all report rows and subtracts from `scheduled_time`, returning the delta in minutes. `PreviousMonthBounds` / `CurrentMonthBounds` return `YYYY-MM-DD` date range strings for the API filter.

### CI

GitHub Actions (`build.yml`) triggers on push to `main`: builds linux/windows/darwin amd64 binaries and publishes them as a GitHub Release tagged `latest`.