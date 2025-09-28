# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with
 code in this repository.

## Project Overview

Granola CLI is a command-line tool for exporting notes from the Granola
note-taking application to Markdown files. It connects to the Granola API,
authenticates using bearer tokens, fetches all notes in JSON format, and
converts them to clean Markdown files while preserving metadata.

## Common Commands

### Build

```bash
go build
```

### Run

```bash
go run main.go export
# Or after building:
./granola export
```

### Test

```bash
go test ./...
go test -v ./...  # verbose output
```

### Module Management

```bash
go mod tidy       # clean up dependencies
go mod download   # download dependencies
```

### Releases

```bash
# Create a new release (automated via GitHub Actions)
git tag v0.1.0
git push origin v0.1.0

# Test release locally (requires GoReleaser)
goreleaser release --snapshot --clean

# Check GoReleaser configuration
goreleaser check
```

### Linting

```bash
# Markdown linting (runs in GitHub Actions, installed via brew)
markdownlint-cli2 "**/*.md"

# Go linting (if golangci-lint is installed)
golangci-lint run
```

## Architecture

The project follows a modular Go CLI application structure:

- **Entry Point**: `main.go` - Uses Charmbracelet's fang for execution context
- **Command Structure**: `cmd/` directory contains Cobra command definitions
  - `cmd/root.go` - Defines the root command with configuration initialization using constructor pattern
  - `cmd/export.go` - Implements the export command for fetching and converting notes
- **Internal Packages**:
  - `internal/api/` - Granola API client with Supabase token authentication
  - `internal/converter/` - JSON to Markdown conversion logic (to be implemented)
  - `internal/models/` - Data models for Granola notes and metadata (to be implemented)
- **Configuration**: Supports multiple configuration sources:
  - Environment variables via `.env` file (using godotenv)
  - Config file (`.granola.toml` in home directory or current directory)
  - Command-line flags (e.g., `--debug`, `--config`, `--supabase`, `--timeout`)
  - Environment variable mapping: `SUPABASE_FILE`, `DEBUG_MODE`
- **Logging**: Uses Charmbracelet's log package for structured logging
  - Debug mode can be enabled via `--debug` flag or config
  - Logger includes timestamp and caller information
  - Log levels: Debug, Info, Warn, Error (defaults to Warn, Debug with debug flag)
  - Logger is created in Execute() and injected via dependency injection
  - **Logging Best Practices**:
    - Log errors only at the command level (cmd package) where they are handled
    - Internal packages should return errors without logging to avoid duplicates
    - Commands return errors to Cobra rather than logging them (Cobra handles display)
    - Debug/Info logging can occur at any level for progress tracking

## Key Dependencies

- **cobra**: Command-line interface framework
- **viper**: Configuration management (env vars, config files, flags)
- **afero**: Filesystem abstraction for testable file operations
- **charmbracelet/fang**: Enhanced command execution with context
- **charmbracelet/log**: Structured logging with customizable output formats
- **godotenv**: .env file support
- **net/http**: Standard library HTTP client for API communication
- **encoding/json**: Standard library JSON parsing

## Build and Release

- **GoReleaser**: Automated release management configured in `.goreleaser.yaml`
  - Builds for Linux, macOS (Darwin), and Windows
  - Creates tar.gz archives (zip for Windows)
  - Automatically generates changelog from commit messages
  - Triggered by pushing version tags (e.g., v1.0.0)

## Development Notes

- The root command is "granola" with "export" as the primary subcommand
- Commands use constructor pattern (e.g., `NewRootCmd()`, `NewExportCmd()`)
- Debug logging is available via the `--debug` flag for troubleshooting API calls
- Configuration precedence: flags > env vars > config file > defaults
- Logger is created in Execute() and passed via dependency injection (no globals)
- Releases are automated via GoReleaser when tags are pushed to GitHub
- Binary builds have CGO disabled for maximum portability

## Testing Approach

- **Test Pattern**: Use sub-test pattern with `t.Run()` for better organization
- **Test Isolation**: All tests should be independent and run in parallel using `t.Parallel()` whenever possible
- **Test Focus**: Write one test per function testing the happy path
- **Mock Filesystem**: Use Afero for filesystem abstraction in tests
- **Test Output**: Always use `io.Discard` for logger output in tests
- **No Framework Testing**: Avoid testing third-party framework functionality (e.g., Cobra's Execute)
- **Dependency Injection**: Refactor code to use interfaces for better testability (planned for API client)

## Granola-Specific Implementation Notes

### API Integration

- Bearer token authentication required for all API calls
- Token extracted from supabase.json file (WorkOS tokens)
- Token file path configured via `--supabase` flag or `SUPABASE_FILE` environment variable
- API endpoint: `https://api.granola.ai/v2/get-documents`
- Required headers:
  - `Authorization: Bearer <token>`
  - `User-Agent: Granola/5.354.0`
  - `X-Client-Version: 5.354.0`
  - `Content-Type: application/json`
  - `Accept: */*`
- HTTP client with configurable timeout (default 2 minutes)
- Handle API rate limiting and error responses gracefully

### Export Process

1. Load supabase.json file from configured path
2. Extract access token from WorkOS tokens in supabase.json
3. Authenticate with Granola API using bearer token
4. Fetch all documents from the API (returns JSON with `docs` array)
5. Parse JSON response into Go structs
6. Convert each document to Markdown format (to be implemented)
7. Preserve metadata as YAML frontmatter (to be implemented)
8. Save files to specified output directory (to be implemented)

### File Naming

- Use note title or ID for filename
- Sanitize filenames for filesystem compatibility
- Handle duplicate names appropriately

### Error Handling

- Validate API token before making requests
- Handle network errors and timeouts
- Provide clear error messages for users
- Support retry logic for failed API calls
- Follow Go best practice: "log errors where they're handled, not where they're created"
- Internal packages return wrapped errors without logging
- Command functions return errors to Cobra for user display
