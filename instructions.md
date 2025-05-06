
### Prerequisites

- Go 1.20 or later

## Testing

### Running Tests

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