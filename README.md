# strava-cli

A command-line tool to track your zone 2 training progress by pulling running data from Strava. Built to make it easy to see trends in your Sunday long runs — distance, average heart rate, time, and pace — without manually tabulating data in the Strava UI.

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- A Strava account
- A Strava API application ([create one here](https://www.strava.com/settings/api))
  - Set the **Authorization Callback Domain** to `localhost`

## Installation

```bash
git clone https://github.com/<your-username>/strava-cli.git
cd strava-cli
go build -o strava-cli.exe .
```

## Setup

Run the auth command and follow the prompts:

```bash
./strava-cli auth
```

You'll be asked for your Strava Client ID and Client Secret, then redirected to Strava in your browser to authorize the app. Credentials and tokens are stored locally in `~/.strava-cli/`.

## Usage

### View your Sunday long runs

```bash
./strava-cli runs
```

### Customise the time range and day

```bash
# Last 24 weeks of Saturday runs
./strava-cli runs --weeks 24 --day saturday
```

### Output

```
DATE            DISTANCE (km)  TIME         AVG HR   PACE (/km)
────            ─────────────  ────         ──────   ──────────
23 Mar 2026     18.50          1h 42m 30s   142 bpm  5:32
16 Mar 2026     16.20          1h 28m 15s   145 bpm  5:27
09 Mar 2026     15.00          1h 22m 00s   140 bpm  5:28
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--weeks` | `-w` | `12` | Number of weeks to look back |
| `--day` | `-d` | `sunday` | Day of week to filter |

## Project Structure

```
strava-cli/
├── main.go                  # Entry point
├── cmd/
│   ├── root.go              # Root cobra command
│   ├── auth.go              # Strava OAuth2 authentication
│   └── runs.go              # Fetch and display runs
└── internal/
    ├── auth/
    │   ├── config.go         # Client ID/secret persistence
    │   ├── oauth.go          # OAuth2 flow and token refresh
    │   └── token.go          # Token storage
    └── strava/
        ├── client.go         # Strava API HTTP client
        └── filter.go         # Weekday filtering
```

## Roadmap

- [ ] Summary statistics (averages, trends over time)
- [ ] Terminal charts for HR and distance trends
- [ ] Local caching to reduce API calls
- [ ] Export to CSV
