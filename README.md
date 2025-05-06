# Biathlon Competition Tracker

A Go application for tracking biathlon competitions, processing events, and generating final reports.

## Requirements

- Go 1.20 or newer

## Project Structure

```
.
├── cmd/
│   └── biathlon-tracker/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── config.go
│   │   ├── competitor.go
│   │   ├── errors.go
│   │   └── event.go
│   └── service/
│       └── competition.go
├── go.mod
└── README.md
```

## Configuration

The application requires a JSON configuration file with the following structure:

```json
{
    "laps": 2,
    "lapLen": 3651,
    "penaltyLen": 50,
    "firingLines": 1,
    "start": "09:30:00",
    "startDelta": "00:00:30"
}
```

## Events Format

Events should be provided in a text file with the following format:

```
[HH:MM:SS.sss] eventID competitorID [extraParams]
```

### Event Types

#### Incoming Events
- 1: Competitor registered
- 2: Start time set by draw
- 3: Competitor on start line
- 4: Competitor started
- 5: Competitor on firing range
- 6: Target hit
- 7: Competitor left firing range
- 8: Competitor entered penalty laps
- 9: Competitor left penalty laps
- 10: Competitor ended main lap
- 11: Competitor cannot continue

#### Outgoing Events
- 32: Competitor disqualified
- 33: Competitor finished

## Building and Running

1. Build the application:
```bash
go build -o biathlon-tracker ./cmd/biathlon-tracker
```

2. Run the application:
```bash
./biathlon-tracker <config_file> <events_file>
```

## Output

The application generates a final report with the following information for each competitor:
- Status (Finished, NotStarted, NotFinished, Disqualified)
- Total time
- Time taken for each lap
- Average speed for each lap
- Time taken for penalty laps
- Average speed over penalty laps
- Number of hits/number of shots

Example output:
```
[NotFinished] 1 [{00:29:03.872, 2.093}, {,}] {00:01:44.296, 0.481} 4/5
```
