# Biathlon Competition Tracker

A Go application for tracking biathlon competitions, processing events, and generating reports.

## Prerequisites

- Go 1.21 or later
- Git

## Project Structure

```
biathlon-tracker/
├── cmd/
│   └── biathlon-tracker/    # Main application entry point
├── internal/
│   ├── domain/             # Domain models and business logic
│   └── service/            # Business logic implementation
├── config/                 # Configuration files
└── sunny_5_skiers/        # Example competition data
```

## Building the Application

To build the application, run:

```bash
go build -o biathlon-tracker cmd/biathlon-tracker/main.go
```

## Running the Application

The application requires two command-line arguments:
1. Path to the configuration file (JSON)
2. Path to the events file

Example:
```bash
./biathlon-tracker config/config.json config/events
```

Or using `go run`:
```bash
go run cmd/biathlon-tracker/main.go config/config.json config/events
```

## Configuration File Format

The configuration file (`config.json`) should contain the following parameters:

```json
{
    "laps": 2,              // Number of laps in the race
    "lapLen": 3500,         // Length of each lap in meters
    "penaltyLen": 150,      // Length of penalty loop in meters
    "firingLines": 2,       // Number of firing lines
    "start": "10:00:00.000", // Race start time
    "startDelta": "00:01:30.000" // Time interval between competitors
}
```

## Events File Format

The events file contains one event per line in the following format:
```
[HH:MM:SS.sss] EVENT_TYPE COMPETITOR_ID [ADDITIONAL_DATA]
```

Example:
```
[10:00:00.000] REGISTERED 1
[10:00:00.000] START_TIME_SET 1 10:00:00.000
[10:00:01.000] STARTED 1
[10:10:00.000] ENDED_MAIN_LAP 1
[10:20:00.000] ENDED_MAIN_LAP 1
```

## Running Tests

To run all tests:
```bash
go test ./...
```

To run tests with coverage:
```bash
go test ./... -cover
```

To run tests in a specific package:
```bash
go test ./internal/domain/...
go test ./internal/service/...
```

## Test Coverage

The project includes comprehensive test coverage for:
- Domain models (Competitor, Config)
- Service layer (CompetitionService)
- Event processing
- Time calculations
- Status tracking

## Output Format

The application generates a report with the following information for each competitor:
- ID
- Status (NotStarted, NotFinished, Disqualified, Finished)
- Total time
- Total penalty time
- Number of hits
- Number of shots
- Average speed

Example output:
```
1. NotStarted 00:00:00.000 00:00:00.000 0 0 0.00
2. NotStarted 00:00:00.000 00:00:00.000 0 0 0.00
3. NotStarted 00:00:00.000 00:00:00.000 0 0 0.00
4. NotStarted 00:00:00.000 00:00:00.000 0 0 0.00
5. NotStarted 00:00:00.000 00:00:00.000 0 0 0.00
```

## Event Types

The application supports the following event types:
- REGISTERED
- START_TIME_SET
- STARTED
- ENDED_MAIN_LAP
- ENDED_PENALTY_LAP
- TARGET_HIT
- TARGET_MISSED
- DISQUALIFIED

## Error Handling

The application handles various error conditions:
- Invalid configuration
- Invalid event format
- Invalid timestamps
- Invalid competitor IDs
- Invalid event types

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
