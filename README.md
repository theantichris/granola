# Granola CLI

[![Go Version](https://img.shields.io/github/go-mod/go-version/theantichris/granola)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/theantichris/granola.svg)](https://pkg.go.dev/github.com/theantichris/granola)
[![Go Report Card](https://goreportcard.com/badge/github.com/theantichris/granola)](https://goreportcard.com/report/github.com/theantichris/granola)
[![Go](https://github.com/theantichris/granola/actions/workflows/go.yml/badge.svg)](https://github.com/theantichris/granola/actions/workflows/go.yml)
[![Markdown Lint](https://github.com/theantichris/granola/actions/workflows/markdown.yml/badge.svg)](https://github.com/theantichris/granola/actions/workflows/markdown.yml)
[![License](https://img.shields.io/github/license/theantichris/granola)](LICENSE)
[![Release](https://img.shields.io/github/v/release/theantichris/granola)](https://github.com/theantichris/granola/releases)

A CLI tool for exporting your Granola notes and transcripts.

## Features

- 📝 **Export Granola Notes** - Export AI-generated notes from the Granola API to local Markdown files
- 🎙️ **Export Raw Transcripts** - Export verbatim meeting transcripts with timestamps from local cache
- 🔄 **JSON to Markdown** - Automatic conversion from ProseMirror JSON to clean Markdown with YAML frontmatter
- 🏷️ **Metadata Preservation** - Maintains note metadata including creation dates, update dates, and tags
- ⏱️ **Timestamp Tracking** - Transcripts include precise timestamps and speaker identification
- 🔐 **Secure Access** - API authentication using bearer tokens from Supabase, local cache file reading for transcripts
- ⚙️ **Flexible Configuration** - Configure via environment variables, config files, or flags
- 📁 **Batch Export** - Export all notes or transcripts in a single command to specified directories
- ⚡ **Incremental Updates** - Only updates files when content is modified (compares timestamps)
- 🚀 **Fast and Efficient** - Built with Go for optimal performance

## Installation

### From Release

Download the latest release from the [releases page](https://github.com/theantichris/granola/releases).

### From Source

```bash
git clone https://github.com/theantichris/granola.git
cd granola
go build -o granola
```

### Using Go Install

```bash
go install github.com/theantichris/granola@latest
```

## Quick Start

### Exporting AI-Generated Notes

Notes are AI-processed summaries and formatted content exported from the Granola API.

#### Finding Your Supabase Credentials

Granola stores authentication credentials in a `supabase.json` file. The location
depends on your operating system:

- **macOS**: `~/Library/Application Support/Granola/supabase.json`
- **Linux**: `~/.config/Granola/supabase.json` or `~/.local/share/Granola/supabase.json`
- **Windows**: `%APPDATA%\Granola\supabase.json`

#### Setup

1. **Configure the path to your supabase.json file:**

   ```bash
   # Via environment variable
   export SUPABASE_FILE="$HOME/Library/Application Support/Granola/supabase.json"

   # Or via .env file
   echo "SUPABASE_FILE=$HOME/Library/Application Support/Granola/supabase.json" >> .env

   # Or via command flag
   granola notes --supabase "$HOME/Library/Application Support/Granola/supabase.json"
   ```

2. **Export all your notes:**

   ```bash
   granola notes
   # Exports to ./notes directory by default
   ```

3. **Export with custom output directory:**

   ```bash
   granola notes --output /path/to/output
   ```

4. **Export with custom timeout:**

   ```bash
   granola notes --timeout 5m
   ```

### Exporting Raw Transcripts

Transcripts are verbatim meeting dialogue with timestamps, exported from the local cache file.

**Note**: Raw transcripts are only available for meetings where audio recording was enabled.

1. **Export all transcripts:**

   ```bash
   granola transcripts
   # Exports to ./transcripts directory by default
   # Reads from ~/Library/Application Support/Granola/cache-v3.json (macOS)
   ```

2. **Export with custom output directory:**

   ```bash
   granola transcripts --output /path/to/output
   ```

3. **Specify custom cache file location:**

   ```bash
   granola transcripts --cache /path/to/cache-v3.json
   ```

### What Gets Exported

#### Notes (Markdown files)

Each note is exported as a separate Markdown file with:

- **YAML frontmatter** containing metadata (ID, created/updated timestamps, tags)
- **Note title** as a top-level heading
- **Note content** converted from ProseMirror JSON to Markdown format
  - Supports headings, paragraphs, bullet lists, and nested lists

**Example note output:**

```markdown
---
id: abc-123
created: "2024-01-01T00:00:00Z"
updated: "2024-01-02T00:00:00Z"
tags:
  - work
  - planning
---

# Meeting Notes

## Key Points

- First important point
- Second important point
  - Nested detail

Action items discussed...
```

#### Transcripts (Text files)

Each transcript is exported as a plain text file with:

- **Header section** containing metadata (title, ID, created/updated timestamps, segment count)
- **Transcript segments** with timestamps and speaker identification
  - **[HH:MM:SS] format** for timestamps
  - **Speaker labels**: "System" (other participants) or "You" (user's microphone)
  - **Verbatim dialogue** including filler words and pauses

**Example transcript output:**

```text
================================================================================
🤖 Team Sync Meeting
ID: abc-123
Created: 2024-01-01T14:00:00.000Z
Updated: 2024-01-01T15:30:00.000Z
Segments: 142
================================================================================

[14:00:04] System: Good morning everyone, how's it going?
[14:00:06] You: Good morning! Ready to start.
[14:00:09] System: Great! Let's dive into the agenda.
[14:00:12] You: Sounds good to me.
```

### Incremental Exports

The CLI intelligently handles repeated exports:

- **First run**: All notes/transcripts are exported
- **Subsequent runs**: Only new or updated files are written
- **Timestamp comparison**: Files are only updated if the document's `updated_at`
  timestamp is newer than the existing file's modification time
- **Performance**: Saves time by skipping unchanged files

This means you can safely run `granola notes` or `granola transcripts` multiple times without worrying
about unnecessary file writes.

## Usage

### Basic Commands

```bash
# Export all notes (requires supabase file path to be configured)
granola notes

# Export notes with specific supabase file
granola notes --supabase /path/to/supabase.json

# Export notes to custom output directory
granola notes --output /path/to/notes

# Export notes with custom timeout
granola notes --timeout 5m

# Export notes with debug logging
granola notes --debug

# Export all transcripts (uses default cache file location)
granola transcripts

# Export transcripts with custom cache file
granola transcripts --cache /path/to/cache-v3.json

# Export transcripts to custom output directory
granola transcripts --output /path/to/transcripts

# Use custom config file
granola --config /path/to/config.toml notes

# Display help
granola --help
granola notes --help
granola transcripts --help
```

### Configuration

The application supports multiple configuration sources with the following
precedence:

1. Command-line flags
2. Environment variables
3. Configuration file
4. Default values

#### Configuration File

Create a `.granola.toml` file in your home directory or current directory:

```toml
debug = true
supabase = "/path/to/supabase.json"
timeout = "2m"
output = "/path/to/notes"
```

#### Environment Variables

Create a `.env` file for local development:

```bash
SUPABASE_FILE=/path/to/supabase.json
DEBUG_MODE=true
```

Or set environment variables directly:

```bash
export SUPABASE_FILE="/path/to/supabase.json"
export DEBUG_MODE=true
```

## Project Structure

```text
granola/
├── cmd/
│   ├── root.go         # Root command and configuration
│   ├── notes.go        # Notes export command (API-based)
│   └── transcripts.go  # Transcripts export command (cache-based)
├── internal/
│   ├── api/            # Granola API client and document models
│   ├── cache/          # Cache file reader for transcripts
│   ├── converter/      # Document to Markdown converter
│   ├── prosemirror/    # ProseMirror JSON to Markdown converter
│   ├── transcript/     # Transcript formatter and writer
│   └── writer/         # File system writer for Markdown files
├── main.go             # Application entry point
├── go.mod              # Go module dependencies
├── go.sum              # Dependency checksums
├── README.md           # Project documentation
├── CLAUDE.md           # Claude AI assistant guide
├── SPEC.md             # Project specification
└── LICENSE             # License file
```

## Development

### Prerequisites

- Go 1.23.1 or higher
- Git

### Building

```bash
# Build for current platform
go build

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o granola-linux
GOOS=darwin GOARCH=amd64 go build -o granola-darwin
GOOS=windows GOARCH=amd64 go build -o granola.exe
```

### Releasing

This project uses [GoReleaser](https://goreleaser.com/) for automated releases.

```bash
# Create a new tag
git tag v0.1.0
git push origin v0.1.0

# For local testing (requires GoReleaser installed)
goreleaser release --snapshot --clean
```

Releases are automatically built and published when a new tag is pushed to
GitHub.

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Linting

```bash
# Install Go linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run Go linter
golangci-lint run

# Install Markdown linter
brew install markdownlint-cli2

# Run Markdown linter
markdownlint-cli2 "**/*.md" "#notes" "#transcripts"
```

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Afero](https://github.com/spf13/afero) - Filesystem abstraction for testing
- [Charmbracelet Log](https://github.com/charmbracelet/log) - Structured logging
- [Charmbracelet Fang](https://github.com/charmbracelet/fang) - Enhanced
  command execution
- [Godotenv](https://github.com/joho/godotenv) - .env file support

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE)
file for details.

## Acknowledgments

- The Go team for the amazing language and tools
- The Cobra and Viper teams for excellent CLI libraries
- The Charmbracelet team for beautiful terminal tools

## Support

For issues, questions, or suggestions, please [open an issue](https://github.com/theantichris/granola/issues).

---

Built with ❤️ using Go
