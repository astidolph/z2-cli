# z2-cli

In my experience starting out zone 2 running can be painfully slow paced and a painfully slow progression. I wanted a quick way to visualise progress by pulling together stats I was regularly putting together myself.

This is a command-line tool and web dashboard to track your zone 2 training progress by pulling running data from Strava. See your runs filtered by heart rate zone, track efficiency factor (EF) over time and measure aerobic fitness progression.

## Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/) (for the web frontend)
- A Strava account
- A Strava API application ([create one here](https://www.strava.com/settings/api))
  - Set the **Authorization Callback Domain** to `localhost`

## Installation

```bash
git clone https://github.com/<your-username>/z2-cli.git
cd z2-cli
go install .
```

This places the binary in your Go bin directory so you can run `z2-cli` from anywhere.

## Setup

### 1. Authenticate with Strava

```bash
z2-cli auth
```

You'll be asked for your Strava Client ID and Client Secret, then redirected to Strava in your browser to authorize the app. Credentials and tokens are stored locally in `~/.z2-cli/`.

### 2. Set your zone 2 heart rate ceiling

```bash
# Set directly
z2-cli config --zone2-hr 150

# Or calculate from age using the Maffetone formula (180 - age)
z2-cli config --age 32
```

This is used to filter runs — only runs with an average HR at or below this value are shown by default.

## Usage

### View your zone 2 runs sorted by efficiency

```bash
z2-cli runs --sort ef
```

This is the most useful command for tracking progress — it shows your zone 2 runs ranked by efficiency factor (EF), so your most aerobically efficient runs are at the top.

### View your zone 2 runs

```bash
z2-cli runs
```

### View all runs (skip zone 2 filtering)

```bash
z2-cli runs --all
```

### Filter for long runs

```bash
z2-cli runs --min-distance 12
```

### Combine filters

```bash
# Zone 2 long runs on Sundays from the last 24 weeks
z2-cli runs --min-distance 12 --day sunday --weeks 24
```

### Output

```
Zone 2 runs (avg HR ≤ 148 bpm) from the last 12 weeks:

DATE          DIST (km)  DIST (mi)  TIME         AVG HR    PACE (/km)  PACE (/mi)  EF
────          ─────────  ─────────  ────         ──────    ──────────  ──────────  ──
23 Mar 2026   18.50      11.49      1h 42m 30s   142 bpm   5:32        8:54        0.0211
16 Mar 2026   12.10      7.52       1h 06m 15s   144 bpm   5:27        8:46        0.0213
09 Mar 2026   15.00      9.32       1h 22m 00s   140 bpm   5:28        8:48        0.0215

Summary (last 12 weeks, 3 runs, 45.6 km / 28.3 mi total):
  Avg EF:   0.0213 ↑ (+2.1% vs prior 12 weeks)
  Avg HR:   142 bpm
  Avg Pace: 5:29/km (8:49/mi)
```

### Runs flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--weeks` | `-w` | `12` | Number of weeks to look back |
| `--day` | `-d` | | Day of week to filter (e.g. sunday, monday) |
| `--min-distance` | | | Minimum distance in km (e.g. 12 for long runs) |
| `--sort` | | `date` | Sort by: date, distance, time, hr, pace, ef |
| `--asc` | | `false` | Sort in ascending order (default is descending) |
| `--all` | `-a` | `false` | Show all runs, skip zone 2 filtering |

### Charts

Generate interactive charts that open in your browser:

```bash
z2-cli chart                        # EF trend (default)
z2-cli chart --type pace            # pace over time (km + mi)
z2-cli chart --type distance        # distance over time (km + mi)
z2-cli chart --type hr              # heart rate over time
z2-cli chart --type all             # all charts on one page
z2-cli chart --weeks 24 --type all  # last 24 weeks, all charts
```

Charts support the same filtering flags as the `runs` command (`--weeks`, `--day`, `--min-distance`, `--all`).

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--type` | `-t` | `ef` | Chart type: ef, pace, distance, hr, all |
| `--weeks` | `-w` | `12` | Number of weeks to look back |
| `--day` | `-d` | | Day of week to filter |
| `--min-distance` | | | Minimum distance in km |
| `--all` | `-a` | `false` | Show all runs, skip zone 2 filtering |

### Efficiency Factor (EF)

EF is calculated as speed (m/s) divided by average heart rate. A higher EF means you're running faster at the same effort — the key indicator that zone 2 training is working. The summary compares your current period's EF against the prior equivalent period to show your trend.

## Web Dashboard

The web UI provides the same data as the CLI in a browser-based dashboard you can access from your phone or any device. Built with SvelteKit (Svelte 5) and Chart.js with a dark theme.

### Pages

- **Dashboard** — Summary cards (EF trend, avg HR, avg pace, total distance) with a dual-axis EF vs Heart Rate chart
- **Runs** — Filterable table with per-run and cumulative average EF, sortable by any column
- **Charts** — EF, pace (km + mi), distance (km + mi), and heart rate charts with configurable lookback period
- **Settings** — Strava connection status, web-based OAuth login, and zone 2 HR configuration

### Development mode

Run the Go API server and SvelteKit dev server separately:

```bash
# Terminal 1 — API server
z2-cli serve

# Terminal 2 — frontend dev server (hot reload)
cd web
npm install
npm run dev
```

Open `http://localhost:5173`. The Vite dev server proxies `/api` requests to the Go backend on port 8080.

### API server

```bash
z2-cli serve              # default port 8080
z2-cli serve --port 3000  # custom port
```

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/health` | Health check |
| `GET` | `/api/auth/status` | Strava connection status |
| `GET` | `/api/auth/login` | Initiate OAuth2 flow (redirects to Strava) |
| `GET` | `/api/auth/callback` | OAuth2 callback (handles token exchange) |
| `GET` | `/api/config` | Get zone 2 HR setting |
| `PUT` | `/api/config` | Update zone 2 HR (by value or age) |
| `GET` | `/api/runs` | Runs and stats (query params: `weeks`, `day`, `minDistance`, `all`, `sort`, `asc`, `refresh`) |
| `GET` | `/api/chart-data` | Chart data arrays (same query params as runs) |
| `POST` | `/api/refresh` | Clear the Strava API cache |

Config, runs, chart-data, and refresh endpoints require an authenticated session.

### Caching

Strava API responses are cached locally in `~/.z2-cli/cache.json` with a 15-minute TTL. This keeps the web dashboard fast and avoids hitting Strava's rate limits (100 requests per 15 minutes, 1000 per day). Use the refresh button in the dashboard or `POST /api/refresh` to clear the cache after a new run.

## Deployment

### Docker

The project includes a multi-stage Dockerfile that builds a single self-contained binary with the frontend embedded:

```bash
docker build -t z2-cli .
docker run -p 8080:8080 -v z2-data:/home/z2user/.z2-cli z2-cli
```

The build process:
1. Builds the SvelteKit frontend to static files
2. Compiles the Go binary with the frontend embedded via `//go:embed`
3. Produces a minimal Alpine image running as a non-root user

### Fly.io

The project is configured for deployment on [Fly.io](https://fly.io) via `fly.toml`:

- Region: `lhr` (London)
- Force HTTPS with auto-redirect
- Persistent volume for config, tokens, and cache (`~/.z2-cli/`)
- Auto-stop/start machines to minimise costs when idle

```bash
fly deploy
```

### Security

- **Session auth** — HMAC-SHA256 signed session cookies (HttpOnly, SameSite=Lax, Secure over HTTPS) with 7-day TTL
- **CSRF protection** — OAuth state parameter signed and verified via cookies
- **Security headers** — HSTS, X-Content-Type-Options: nosniff, X-Frame-Options: DENY, Referrer-Policy
- **CORS** — Scoped to the local dev origin only (`localhost:5173`)
- **Non-root container** — Docker image runs as an unprivileged user

## Project Structure

```
z2-cli/
├── main.go                  # Entry point
├── frontend_embed.go        # Embeds built frontend into binary (production builds)
├── Dockerfile               # Multi-stage build (Node → Go → Alpine runtime)
├── fly.toml                 # Fly.io deployment config
├── cmd/
│   ├── root.go              # Root cobra command
│   ├── auth.go              # Strava OAuth2 authentication
│   ├── chart.go             # Interactive HTML chart generation
│   ├── config.go            # Training settings (zone 2 HR)
│   ├── runs.go              # Table display and formatting
│   └── serve.go             # Web API server command
├── internal/
│   ├── api/
│   │   ├── handlers.go      # REST API route handlers
│   │   ├── middleware.go     # CORS, security headers, session auth
│   │   ├── response.go      # JSON response helpers
│   │   └── server.go        # HTTP server and SPA file serving
│   ├── auth/
│   │   ├── config.go        # Config persistence (API creds + zone 2 HR)
│   │   ├── oauth.go         # OAuth2 flow and token refresh
│   │   └── token.go         # Token storage
│   ├── cache/
│   │   └── cache.go         # File-based Strava API response cache
│   ├── chart/
│   │   └── chart.go         # go-echarts chart rendering (EF, pace, distance, HR)
│   ├── service/
│   │   └── runs.go          # Core data logic (fetch, filter, sort, summarise)
│   ├── stats/
│   │   ├── efficiency.go    # Efficiency factor calculation
│   │   └── summary.go       # Period summaries and trend comparison
│   └── strava/
│       ├── client.go        # Strava API HTTP client
│       └── filter.go        # Weekday, HR, and distance filters
└── web/                     # SvelteKit frontend (Svelte 5, dark theme)
    ├── src/
    │   ├── lib/
    │   │   ├── api.ts       # Typed API client
    │   │   ├── types.ts     # TypeScript interfaces matching Go types
    │   │   ├── format.ts    # Display formatting helpers
    │   │   └── components/  # NavBar, SummaryCard, LineChart, RunsTable, FilterBar
    │   └── routes/
    │       ├── +page.svelte         # Dashboard (summary + dual-axis EF/HR chart)
    │       ├── runs/+page.svelte    # Runs table with filters and cumulative avg EF
    │       ├── charts/+page.svelte  # All chart types with configurable lookback
    │       └── settings/+page.svelte # Strava login + zone 2 HR config
    ├── svelte.config.js     # adapter-static for single-binary embedding
    └── vite.config.ts       # Dev proxy to Go API server
```
