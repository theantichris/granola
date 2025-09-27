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
  - `cmd/root.go` - Defines the root command with configuration initialization
  - `cmd/export.go` - Implements the export command for fetching and converting notes
- **Internal Packages**:
  - `internal/api/` - Granola API client with bearer token authentication
  - `internal/converter/` - JSON to Markdown conversion logic
  - `internal/models/` - Data models for Granola notes and metadata
- **Configuration**: Supports multiple configuration sources:
  - Environment variables via `.env` file (using godotenv)
  - Config file (`.config.toml` in home directory or current directory)
  - Command-line flags (e.g., `--debug`, `--config`, `--output`)
  - Granola-specific config: `GRANOLA_API_TOKEN`, `GRANOLA_API_URL`, `GRANOLA_OUTPUT_DIR`
- **Logging**: Uses Charmbracelet's log package for structured logging
  - Debug mode can be enabled via `--debug` flag or config
  - Logger includes timestamp and caller information
  - Log levels: Debug, Info, Warn, Error (defaults to Warn, Debug with debug flag)

## Key Dependencies

- **cobra**: Command-line interface framework
- **viper**: Configuration management (env vars, config files, flags)
- **charmbracelet/fang**: Enhanced command execution with context
- **charmbracelet/log**: Structured logging with customizable output formats
- **godotenv**: .env file support
- **Additional dependencies** (to be added):
  - HTTP client for API communication
  - JSON parsing libraries
  - File I/O utilities for Markdown export

## Build and Release

- **GoReleaser**: Automated release management configured in `.goreleaser.yaml`
  - Builds for Linux, macOS (Darwin), and Windows
  - Creates tar.gz archives (zip for Windows)
  - Automatically generates changelog from commit messages
  - Triggered by pushing version tags (e.g., v1.0.0)

## Development Notes

- The root command is "granola" with "export" as the primary subcommand
- Debug logging is available via the `--debug` flag for troubleshooting API calls
- Configuration precedence: flags > env vars > config file > defaults
- Logger instance is globally available as `cmd.Logger`
- Releases are automated via GoReleaser when tags are pushed to GitHub
- Binary builds have CGO disabled for maximum portability

## Granola-Specific Implementation Notes

### API Integration

- Bearer token authentication required for all API calls
- Token should be stored securely (environment variable recommended)
- API base URL is configurable for different environments
- Handle API rate limiting and error responses gracefully

### Export Process

1. Authenticate with Granola API using bearer token
2. Fetch all notes from the API (returns JSON)
3. Parse JSON response into Go structs
4. Convert each note to Markdown format
5. Preserve metadata as YAML frontmatter
6. Save files to specified output directory

### File Naming

- Use note title or ID for filename
- Sanitize filenames for filesystem compatibility
- Handle duplicate names appropriately

### Error Handling

- Validate API token before making requests
- Handle network errors and timeouts
- Provide clear error messages for users
- Support retry logic for failed API calls
