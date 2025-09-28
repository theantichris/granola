# Granola CLI

[![Go Version](https://img.shields.io/github/go-mod/go-version/theantichris/granola)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/theantichris/granola.svg)](https://pkg.go.dev/github.com/theantichris/granola)
[![Go Report Card](https://goreportcard.com/badge/github.com/theantichris/granola)](https://goreportcard.com/report/github.com/theantichris/granola)
[![Go](https://github.com/theantichris/granola/actions/workflows/go.yml/badge.svg)](https://github.com/theantichris/granola/actions/workflows/go.yml)
[![Markdown Lint](https://github.com/theantichris/granola/actions/workflows/markdown.yml/badge.svg)](https://github.com/theantichris/granola/actions/workflows/markdown.yml)
[![License](https://img.shields.io/github/license/theantichris/granola)](LICENSE)
[![Release](https://img.shields.io/github/v/release/theantichris/granola)](https://github.com/theantichris/granola/releases)

A CLI tool for exporting your Granola notes to Markdown files.

## Features

- 📝 **Export Granola Notes** - Export all your notes from the Granola API
- 🔄 **JSON to Markdown** - Automatic conversion from JSON format to clean Markdown
- 🏷️ **Metadata Preservation** - Maintains note metadata including creation
 dates and tags
- 🔐 **Bearer Token Auth** - Secure API authentication using bearer tokens
- ⚙️ **Flexible Configuration** - Configure via environment variables, config
 files, or flags
- 📁 **Batch Export** - Export all notes in a single command
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

1. **Configure the path to your supabase.json file:**

   ```bash
   # Via environment variable
   export SUPABASE_FILE="/path/to/supabase.json"

   # Or via .env file
   echo "SUPABASE_FILE=/path/to/supabase.json" >> .env

   # Or via command flag
   granola export --supabase /path/to/supabase.json
   ```

2. **Export all your notes:**

   ```bash
   granola export
   # Currently prints to stdout (file export coming soon)
   ```

3. **Export with custom timeout:**

   ```bash
   granola export --timeout 5m
   ```

## Usage

### Basic Commands

```bash
# Export all notes (requires supabase file path to be configured)
granola export

# Export with specific supabase file
granola export --supabase /path/to/supabase.json

# Export with custom timeout
granola export --timeout 5m

# Export with debug logging
granola export --debug

# Use custom config file
granola --config /path/to/config.toml export

# Display help
granola --help
granola export --help
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
│   └── export.go       # Export command implementation
├── internal/
│   ├── api/            # Granola API client
│   ├── converter/      # JSON to Markdown converter
│   └── models/         # Data models for notes
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
markdownlint-cli2 "**/*.md"
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
