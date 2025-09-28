# Project Specification

## Executive Summary

Granola CLI is a command-line tool that enables users to export their notes from
the Granola note-taking application to local Markdown files. The tool provides
a simple, efficient way to backup, migrate, or work with Granola notes outside
the application while preserving all metadata and formatting.

## Project Overview

### Problem Statement

Users of the Granola note-taking application need a way to export their notes
for backup, migration, or offline access. Currently, notes are only accessible
through the Granola application, limiting users' ability to work with their
content in other tools or preserve it locally.

### Solution

A command-line tool that connects to the Granola API, authenticates using
bearer tokens, fetches all user notes in JSON format, and converts them to
clean, readable Markdown files with preserved metadata.

### Goals and Objectives

- [x] Provide simple, secure authentication via Supabase tokens
- [x] Connect to Granola API with proper headers
- [x] Support configurable timeout for API requests
- [ ] Export all notes from Granola API to local Markdown files
- [ ] Preserve note metadata (creation date, tags, etc.) in exports
- [ ] Support batch export of all notes in a single command
- [ ] Create well-organized file structure for exported notes
- [ ] Incremental exports (only new notes)

### Success Criteria

- Successfully exports 100% of user notes from Granola API
- Maintains data integrity during conversion (no content loss)
- Preserves all metadata in a standard format (YAML frontmatter)
- Completes export process efficiently (< 1 minute for 1000 notes)
- Produces valid, readable Markdown files compatible with common editors

## Scope

### In Scope

- Bearer token authentication with Granola API
- Fetching all notes via API endpoints
- JSON to Markdown conversion
- Metadata preservation in YAML frontmatter
- Configurable output directory
- Error handling and retry logic
- Debug mode for troubleshooting

### Out of Scope

- Two-way synchronization with Granola
- Selective note export (filtering)
- Real-time export/watch mode
- Export to formats other than Markdown
- Note editing or modification capabilities

### Future Considerations

- Support for modified note exports
- Export filtering by date, tags, or categories
- Multiple export formats (HTML, PDF, etc.)
- Integration with other note-taking tools
- Scheduled/automated exports

## Requirements

### Functional Requirements

#### Core Features

1. **API Authentication** ✅
   - Secure authentication using Supabase tokens
   - Token extraction from supabase.json file
   - Token configuration via environment variables or config file
   - Priority: High

2. **Note Export**
   - Fetch all notes from Granola API
   - Parse JSON response data
   - Handle pagination if required
   - Priority: High

3. **Markdown Conversion**
   - Convert JSON note data to Markdown format
   - Preserve note content and formatting
   - Add metadata as YAML frontmatter
   - Priority: High

4. **File Management**
   - Create output directory structure
   - Generate appropriate filenames
   - Handle duplicate names
   - Write files to disk
   - Priority: High

### Non-Functional Requirements

#### Performance

- Export 1000 notes in under 60 seconds
- Efficient memory usage for large note collections
- Concurrent API requests where applicable
- Minimal CPU overhead during conversion

#### Security

- Secure storage of bearer tokens (environment variables)
- No token logging in debug output
- HTTPS-only API communication
- No sensitive data in configuration files

#### Usability

- Simple command-line interface
- Clear error messages and feedback
- Progress indicators for long operations
- Comprehensive help documentation
- Configuration via multiple sources (env, config, flags)

#### Reliability

- Graceful handling of API errors
- Automatic retry for transient failures
- Partial export recovery (resume from failure)
- Data integrity verification
- Clear error reporting with actionable messages

### Logging Strategy

#### Error Logging Best Practices

- **Log Once**: Log errors only where they are handled (typically in cmd package), not where they are created
- **Internal Packages**: Return errors without logging to prevent duplicate log entries
- **Command Functions**: Return errors to Cobra framework rather than logging directly - Cobra handles error display to users
- **Avoid Duplication**: Since Cobra prints returned errors to stderr, logging them in commands creates duplicate output

#### Appropriate Logging Locations

- **Debug/Info Logging**: Can occur at any level for progress tracking and debugging
- **Error Logging**: Only at the top of the call stack where errors are actually handled
- **Library Code**: Never log errors in internal packages - just wrap and return them
- **Background Tasks**: Exception where error logging is required since there's no caller to return to

## Technical Architecture

### System Architecture

```text
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   CLI       │────▶│ API Client   │────▶│ Granola API │
│  (Cobra)    │     │ (HTTP/Auth)  │     │   (JSON)    │
└─────────────┘     └──────────────┘     └─────────────┘
       │                    │
       │                    ▼
       │            ┌──────────────┐
       │            │ JSON Parser  │
       │            │   (Models)   │
       │            └──────────────┘
       │                    │
       │                    ▼
       │            ┌──────────────┐
       └───────────▶│  Converter   │
                    │ (MD Generator)│
                    └──────────────┘
                            │
                            ▼
                    ┌──────────────┐
                    │ File System  │
                    │   (Export)   │
                    └──────────────┘
```

### Technology Stack

- **Language**: Go 1.23.1+
- **Framework**: Cobra (CLI framework)
- **Configuration**: Viper (TOML config files, environment variables, flags)
- **Logging**: Charmbracelet/log (structured logging with caller and timestamp)
- **Testing**: Afero for filesystem abstraction
- **Build & Release**: GoReleaser (automated multi-platform builds and releases)
- **Key Libraries**:
  - charmbracelet/fang: Enhanced command execution with context
  - godotenv: .env file support for local development
  - spf13/afero: Filesystem abstraction for testable file operations
  - net/http: HTTP client for API communication
  - encoding/json: JSON parsing and serialization
  - gopkg.in/yaml.v3: YAML frontmatter generation (planned)

### Data Model

#### Note Structure (JSON from API)

```json
{
  "id": "string",
  "title": "string",
  "content": "string",
  "created_at": "ISO8601 timestamp",
  "updated_at": "ISO8601 timestamp",
  "tags": ["string"],
  "metadata": {
    "key": "value"
  }
}
```

#### Markdown Output Format

```markdown
---
id: note-id
created: 2024-01-01T00:00:00Z
updated: 2024-01-01T00:00:00Z
tags: [tag1, tag2]
---

# Note Title

Note content in Markdown format...
```

### API Design

#### Granola API Endpoints

Authentication

- Header: `Authorization: Bearer <token>` (extracted from supabase.json)
- Additional Headers:
  - `User-Agent: Granola/5.354.0`
  - `X-Client-Version: 5.354.0`
  - `Content-Type: application/json`
  - `Accept: */*`

Get All Documents

- Endpoint: `GET https://api.granola.ai/v2/get-documents`
- Response: JSON object with `docs` array containing document objects
- HTTP Client: Configurable timeout (default 2 minutes)

Error Responses

- 401: Invalid or missing authentication token
- 429: Rate limit exceeded
- 500: Server error

## User Interface

### Command-Line Interface

```bash
# Main command
granola [global-flags] <command> [command-flags]

# Global Flags
--config string   Config file path (default: $HOME/.config.toml)
--debug          Enable debug logging
--help           Show help

# Export Command
granola export [flags]

# Export Flags
--supabase string    Path to supabase.json file (overrides env/config)
--timeout duration   HTTP timeout for API requests (default: 2m)
```

### Terminal User Interface

Not applicable - this is a CLI-only tool with no interactive TUI components.

## Testing Strategy

### Testing Approach

- **Unit testing**: Test individual components with dependency injection
  - Using sub-test pattern for better test organization
  - Happy path testing implemented
  - Mock filesystem using Afero for file operations
  - All tests should be independent and run in parallel using `t.Parallel()` whenever possible
- **Integration testing**: Test API client with mock server (planned)
- **End-to-end testing**: Full export workflow with test data (planned)
- **Performance testing**: Benchmark large note exports (1000+ notes) (planned)
- **Error scenario testing**: Invalid tokens, network failures, malformed data (planned)

### Current Test Coverage

- `cmd/root.go`: Unit tests for command creation and configuration
- `cmd/export.go`: Basic test structure (awaiting API interface refactoring)
- `internal/api/`: Token extraction and document fetching tests

## Documentation

### User Documentation

- **Installation guide**: Platform-specific installation instructions
- **Quick start guide**: Getting started with first export
- **CLI reference**: Complete command and flag documentation
- **Configuration guide**: Setting up API tokens and preferences
- **Troubleshooting**: Common issues and solutions

### Developer Documentation

- **Architecture documentation**: System design and component interaction
- **API client documentation**: Granola API integration details
- **Converter documentation**: Markdown conversion logic
- **Contributing guidelines**: How to contribute to the project
- **Testing guide**: How to run and write tests

## Risks and Mitigation

| Risk                      | Impact | Probability | Mitigation Strategy                           |
|---------------------------|--------|-------------|-----------------------------------------------|
| API changes/deprecation   | High   | Low         | Version API calls, maintain compatibility     |
| Rate limiting             | Medium | Medium      | Implement retry logic with exponential backoff|
| Large note collections    | Medium | Medium      | Stream processing, pagination support         |
| Token security breach     | High   | Low         | Secure storage, no logging of sensitive data  |
| Data loss during export   | High   | Low         | Validation, checksums, transaction-like saves |

## Dependencies

### External Dependencies

- **Granola API**: Source of note data, requires active service
- **Internet connectivity**: Required for API communication
- **Go standard library**: Core functionality (net/http, encoding/json)
- **Cobra/Viper**: CLI framework and configuration management
- **Charmbracelet libraries**: Enhanced CLI experience

### Internal Dependencies

- Export command depends on API client
- API client depends on configuration (token, URL)
- Converter depends on models package
- File writer depends on converter output

## Constraints

### Technical Constraints

- Requires Go 1.23.1 or higher
- API token required for authentication
- Internet connection required for API access
- Local filesystem write permissions needed

### Business Constraints

- Must work with existing Granola API (no modifications)
- Single developer/maintainer resource
- Open-source project (no licensing fees)

## Maintenance and Support

### Maintenance Plan

- Regular dependency updates via Dependabot
- API compatibility monitoring
- Bug fixes based on user reports
- Feature additions based on community feedback
- Security patches as needed

### Support Strategy

- GitHub Issues for bug reports and feature requests
- README documentation for common use cases
- GitHub Discussions for community support
- Release notes for version changes

### Update and Release Process

1. **Development**: Work on features/fixes in feature branches
2. **Testing**: Run tests with `go test ./...`
3. **Version Tagging**: Tag releases following semantic versioning (v0.1.0, v1.0.0)
4. **Automated Release**: GoReleaser automatically builds and publishes
 releases when tags are pushed
5. **Distribution**: Binaries available for Linux, macOS, and Windows via
 GitHub Releases

## Glossary

| Term          | Definition                                                     |
|---------------|----------------------------------------------------------------|
| Bearer Token  | Authentication token used to access Granola API                |
| Frontmatter   | YAML metadata block at the beginning of Markdown files         |
| Export        | Process of downloading and converting notes to local files     |
| Granola       | The note-taking application from which notes are exported      |
| Markdown      | Plain text formatting syntax for creating formatted documents  |

## References

- [Cobra Documentation](https://cobra.dev/)
- [Viper Documentation](https://github.com/spf13/viper)
- [Markdown Specification](https://spec.commonmark.org/)
- [YAML Frontmatter Standard](https://jekyllrb.com/docs/front-matter/)
- [GoReleaser Documentation](https://goreleaser.com/)
