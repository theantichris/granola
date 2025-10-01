# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with
 code in this repository.

## Project Overview

Granola CLI is a command-line tool for exporting notes and transcripts from the Granola
note-taking application. It provides two distinct export capabilities:

1. **Notes Export**: Connects to the Granola API, authenticates using bearer tokens,
   fetches AI-generated notes in JSON format, and converts them to clean Markdown files
2. **Transcripts Export**: Reads the local Granola cache file, extracts raw meeting transcripts
   with timestamps and speaker identification, and exports them to plain text files

## Common Commands

### Build

```bash
go build
```

### Run

```bash
# Export notes (AI-generated from API)
go run main.go notes
# Or after building:
./granola notes

# Export transcripts (raw from cache file)
go run main.go transcripts
# Or after building:
./granola transcripts
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
markdownlint-cli2 "**/*.md" "#notes" "#transcripts"

# Go linting (if golangci-lint is installed)
golangci-lint run
```

## Architecture

The project follows a modular Go CLI application structure:

- **Entry Point**: `main.go` - Uses Charmbracelet's fang for execution context
- **Command Structure**: `cmd/` directory contains Cobra command definitions
  - `cmd/root.go` - Defines the root command with configuration initialization using constructor pattern
  - `cmd/notes.go` - Implements the notes command for fetching and converting AI-generated notes from API
  - `cmd/transcripts.go` - Implements the transcripts command for reading and exporting raw transcripts from cache
- **Internal Packages**:
  - `internal/api/` - Granola API client with Supabase token authentication and document models (including ProseMirror structures)
  - `internal/cache/` - Cache file reader for extracting transcript data from local Granola cache
  - `internal/converter/` - Document to Markdown conversion with YAML frontmatter
  - `internal/prosemirror/` - ProseMirror JSON to Markdown conversion and plain text extraction
  - `internal/transcript/` - Transcript formatter for converting segments to readable text with timestamps
  - `internal/writer/` - File system writer for Markdown files with sanitization
- **Configuration**: Supports multiple configuration sources:
  - Environment variables via `.env` file (using godotenv)
  - Config file (`.granola.toml` in home directory or current directory)
  - Command-line flags:
    - Global: `--debug`, `--config`
    - Notes command: `--supabase`, `--timeout`, `--output`
    - Transcripts command: `--cache`, `--output`
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

- The root command is "granola" with "notes" and "transcripts" as the primary subcommands
- Commands use constructor pattern (e.g., `NewRootCmd()`, `NewNotesCmd()`, `NewTranscriptsCmd()`)
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
- API endpoint: `https://api.granola.ai/v2/get-documents` (POST request)
- Request body includes: `limit: 100`, `offset: 0`, `include_last_viewed_panel: true`
- Required headers:
  - `Authorization: Bearer <token>`
  - `User-Agent: Granola/5.354.0`
  - `X-Client-Version: 5.354.0`
  - `Content-Type: application/json`
  - `Accept: */*`
- HTTP client with configurable timeout (default 2 minutes)
- Handle API rate limiting and error responses gracefully
- Content is returned in ProseMirror JSON format in `last_viewed_panel.content`

### Notes Export Process (API-Based)

1. Load supabase.json file from configured path
2. Extract access token from WorkOS tokens in supabase.json
3. Authenticate with Granola API using bearer token (POST request with include_last_viewed_panel)
4. Fetch all documents from the API (returns JSON with `docs` array)
5. Parse JSON response into Go structs (Document model with ID, Title, LastViewedPanel, CreatedAt, UpdatedAt, Tags)
6. For each document:
   - Check if file exists and compare `updated_at` timestamp with file modification time
   - Skip if file is up-to-date (incremental export)
   - Convert ProseMirror JSON content to Markdown (supports headings, paragraphs, bullet lists, nested lists)
   - Add YAML frontmatter with metadata
   - Sanitize filenames and handle duplicates
   - Save/update file in output directory (default: ./notes)

### Transcript Export Process (Cache-Based)

1. Locate Granola cache file (default: `~/Library/Application Support/Granola/cache-v3.json`)
2. Read and parse double-JSON encoded cache structure:
   - Outer JSON: `{"cache": "<json-string>"}`
   - Inner JSON: Contains `state.documents` and `state.transcripts` maps
3. Extract document metadata (ID, Title, CreatedAt, UpdatedAt) from documents map
4. Extract transcript segments from transcripts map (keyed by document ID)
5. For each document with transcript segments:
   - Check if file exists and compare `updated_at` timestamp with file modification time
   - Skip if file is up-to-date (incremental export)
   - Format segments with timestamps and speaker identification:
     - Parse RFC3339 timestamps to HH:MM:SS format
     - Map "system" source to "System" speaker
     - Map "microphone" source to "You" speaker
   - Create text file with metadata header and formatted transcript
   - Sanitize filenames and handle duplicates
   - Save/update file in output directory (default: ./transcripts)

**Important Notes:**

- Notes are available for ALL meetings (AI-generated from various sources)
- Transcripts are ONLY available for meetings where audio recording was enabled
- The Granola API does NOT provide raw transcript data - it must be read from the local cache
- Speaker identification is limited to "System" (other participants/system audio) and "You" (user's microphone)
- Granola does not provide named speaker labels as it doesn't join meetings as a participant

### File Naming

- Use note title as filename, fallback to ID if title is empty
- Sanitize filenames by removing invalid characters (regex: `[<>:"/\\|?*\x00-\x1f]`)
- Handle duplicate names by appending `_2`, `_3`, etc.
- Limit filename length to 100 characters for compatibility

### Error Handling

- Validate API token before making requests
- Handle network errors and timeouts
- Provide clear error messages for users
- Support retry logic for failed API calls
- Follow Go best practice: "log errors where they're handled, not where they're created"
- Internal packages return wrapped errors without logging
- Command functions return errors to Cobra for user display
