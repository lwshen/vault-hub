# E2E Tests for VaultHub CLI

This directory contains end-to-end tests for the VaultHub CLI commands.

## Running Tests

### Run all E2E tests
```bash
go test -v ./e2e/...
```

### Run specific test
```bash
go test -v ./e2e/... -run TestUpdate_ByName_WithValue
```

### Run with debug output
```bash
go test -v ./e2e/... -debug
```

## Test Structure

- `setup_test.go` - Test infrastructure and utilities
- `update_test.go` - Tests for the `update` command
- `fixtures/` - Test data files

## Test Coverage

### Core Functionality
- Update by name
- Update by ID
- Update from value file
- Client-side encryption (default)
- Disable client-side encryption
- JSON output format

### Error Scenarios
- Missing name and ID
- Both name and ID provided
- Missing value
- Empty value file
- Non-existent vault
- Invalid value file path

### Edge Cases
- Special characters in value
- Multiline values

## Architecture

The E2E tests work by:

1. Building the server and CLI binaries
2. Starting a test server with SQLite in-memory database
3. Executing CLI commands against the test server
4. Verifying outputs and exit codes
5. Cleaning up resources after tests

## Adding New Tests

To add a new test:

1. Create a test function in `update_test.go` (or create new file)
2. Use `StartTestServer(t)` to get a test server instance
3. Use `RunCLI(t, args...)` to execute CLI commands
4. Use `result.MustSucceed(t)` or `result.MustFail(t, code)` to verify results
5. Use `result.ContainsStdout(t, "expected")` to verify output content