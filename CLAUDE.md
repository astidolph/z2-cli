# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

z2-cli is a CLI tool and web dashboard for tracking Zone 2 running training via Strava. It measures efficiency factor (EF = speed / heart rate) to track aerobic fitness over time. Go backend with a SvelteKit (Svelte 5) frontend.

## Build & Run Commands

### Go Backend
```bash
go build -o z2-cli.exe .        # Build binary
go install .                     # Install to GOPATH/bin
go vet ./...                     # Check for issues
z2-cli auth                      # OAuth2 setup with Strava
z2-cli runs                      # Show runs table
z2-cli chart                     # Generate HTML charts
z2-cli serve                     # Start API server on :8080
```

### Web Frontend (from /web directory)
```bash
npm install                      # Install dependencies
npm run dev                      # Dev server on :5173 (proxies /api → :8080)
npm run build                    # Production static build → web/build/
npm run check                    # TypeScript/Svelte type checking
npm run preview                  # Preview production build
```

### Development Workflow
Run the Go API server (`z2-cli serve`) and the Vite dev server (`cd web && npm run dev`) simultaneously. The Vite config proxies `/api` requests to localhost:8080.

## Architecture

### Data Flow
```
CLI commands / Web UI
  → Service layer (internal/service)
    → Cache check (internal/cache, 15-min TTL, ~/.z2-cli/cache.json)
    → Strava API client (internal/strava)
    → Filters & Stats (internal/strava/filter.go, internal/stats)
  → Output: table, HTML chart, or JSON API response
```

### Backend Packages
- **cmd/** — Cobra CLI commands (auth, config, runs, chart, serve)
- **internal/api/** — REST API: handlers, CORS middleware, SPA file serving. Routes under `/api/` (health, auth/status, config, runs, chart-data, refresh)
- **internal/service/** — Core orchestrator. `FetchRuns()` returns structured `RunsResult` with current/prior period data for trend comparison
- **internal/strava/** — Strava v3 API client with pagination; filters by weekday, max HR, min distance
- **internal/stats/** — EF calculation and summary aggregation (avg EF, HR, pace, total km, trend %)
- **internal/cache/** — File-based JSON cache with TTL and coverage-aware freshness
- **internal/chart/** — go-echarts chart generation (EF, pace, distance, HR, combined)
- **internal/auth/** — OAuth2 flow (callback on :8089), token refresh, config persistence

### Frontend Structure (web/src/)
- **lib/types.ts** — TypeScript interfaces mirroring Go structs (Activity, Summary, RunsResponse, ChartDataResponse)
- **lib/api.ts** — Typed HTTP client with query parameter builders
- **lib/components/** — NavBar, SummaryCard, LineChart (Chart.js), RunsTable, FilterBar
- **routes/** — SvelteKit pages: dashboard (/), runs (/runs), charts (/charts), settings (/settings)

### Key Design Decisions
- No external web framework — uses `net/http` directly
- Frontend uses `@sveltejs/adapter-static` for static site generation, served by Go's SPA handler in production
- EF trend compares current N-week window against the prior N-week window
- Config/tokens/cache stored in `~/.z2-cli/` as JSON files
- Svelte 5 runes mode — Chart.js canvas refs use `$state` and `$effect` (not legacy reactive statements)
