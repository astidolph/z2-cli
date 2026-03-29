# z2-cli

A command-line tool to track your zone 2 training progress by pulling running data from Strava. See your runs filtered by heart rate zone, track efficiency factor (EF) over time, and measure aerobic fitness progression — all from the terminal.

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
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

DATE          DISTANCE (km)  TIME         AVG HR    PACE (/km)  PACE (/mi)  EF
────          ─────────────  ────         ──────    ──────────  ──────────  ──
23 Mar 2026   18.50          1h 42m 30s   142 bpm   5:32        8:54        0.0211
16 Mar 2026   12.10          1h 06m 15s   144 bpm   5:27        8:46        0.0213
09 Mar 2026   15.00          1h 22m 00s   140 bpm   5:28        8:48        0.0215

Summary (last 12 weeks, 3 runs, 45.6 km total):
  Avg EF:   0.0213 ↑ (+2.1% vs prior 12 weeks)
  Avg HR:   142 bpm
  Avg Pace: 5:29/km (8:49/mi)
```

### Efficiency Factor (EF)

EF is calculated as speed (m/s) divided by average heart rate. A higher EF means you're running faster at the same effort — the key indicator that zone 2 training is working. The summary compares your current period's EF against the prior equivalent period to show your trend.

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--weeks` | `-w` | `12` | Number of weeks to look back |
| `--day` | `-d` | | Day of week to filter (e.g. sunday, monday) |
| `--min-distance` | | | Minimum distance in km (e.g. 12 for long runs) |
| `--sort` | | `date` | Sort by: date, distance, time, hr, pace, ef |
| `--asc` | | `false` | Sort in ascending order (default is descending) |
| `--all` | `-a` | `false` | Show all runs, skip zone 2 filtering |

## Project Structure

```
z2-cli/
├── main.go                  # Entry point
├── cmd/
│   ├── root.go              # Root cobra command
│   ├── auth.go              # Strava OAuth2 authentication
│   ├── config.go            # Training settings (zone 2 HR)
│   └── runs.go              # Fetch, filter, and display runs
└── internal/
    ├── auth/
    │   ├── config.go         # Config persistence (API creds + zone 2 HR)
    │   ├── oauth.go          # OAuth2 flow and token refresh
    │   └── token.go          # Token storage
    ├── stats/
    │   ├── efficiency.go     # Efficiency factor calculation
    │   └── summary.go        # Period summaries and trend comparison
    └── strava/
        ├── client.go         # Strava API HTTP client
        └── filter.go         # Weekday, HR, and distance filters
```
