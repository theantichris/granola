package writer

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/theantichris/granola/internal/api"
)

func TestWrite(t *testing.T) {
	t.Run("writes documents to markdown files", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		outputDir := "/test-output"

		docs := []api.Document{
			{
				ID:        "doc-1",
				Title:     "First Meeting",
				Content:   "Meeting content here",
				CreatedAt: "2024-01-01T00:00:00Z",
				UpdatedAt: "2024-01-01T00:00:00Z",
				Tags:      []string{"work"},
			},
			{
				ID:        "doc-2",
				Title:     "Second Meeting",
				Content:   "More content",
				CreatedAt: "2024-01-02T00:00:00Z",
				UpdatedAt: "2024-01-02T00:00:00Z",
				Tags:      []string{"personal"},
			},
		}

		err := Write(docs, outputDir, fs)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Check that directory was created
		exists, err := afero.DirExists(fs, outputDir)
		if err != nil {
			t.Fatalf("failed to check directory: %v", err)
		}
		if !exists {
			t.Error("expected output directory to be created")
		}

		// Check that files were created
		file1, err := afero.ReadFile(fs, filepath.Join(outputDir, "First Meeting.md"))
		if err != nil {
			t.Fatalf("failed to read first file: %v", err)
		}

		if !strings.Contains(string(file1), "# First Meeting") {
			t.Error("expected first file to contain title")
		}

		file2, err := afero.ReadFile(fs, filepath.Join(outputDir, "Second Meeting.md"))
		if err != nil {
			t.Fatalf("failed to read second file: %v", err)
		}

		if !strings.Contains(string(file2), "# Second Meeting") {
			t.Error("expected second file to contain title")
		}
	})

	t.Run("handles duplicate filenames", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		outputDir := "/test-duplicates"

		docs := []api.Document{
			{
				ID:        "doc-1",
				Title:     "Meeting",
				Content:   "First meeting",
				CreatedAt: "2024-01-01T00:00:00Z",
				UpdatedAt: "2024-01-01T00:00:00Z",
			},
			{
				ID:        "doc-2",
				Title:     "Meeting",
				Content:   "Second meeting",
				CreatedAt: "2024-01-02T00:00:00Z",
				UpdatedAt: "2024-01-02T00:00:00Z",
			},
		}

		err := Write(docs, outputDir, fs)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Check that both files exist with different names
		file1Exists, _ := afero.Exists(fs, filepath.Join(outputDir, "Meeting.md"))
		file2Exists, _ := afero.Exists(fs, filepath.Join(outputDir, "Meeting_2.md"))

		if !file1Exists {
			t.Error("expected first file to exist")
		}
		if !file2Exists {
			t.Error("expected second file with _2 suffix to exist")
		}
	})

	t.Run("sanitizes invalid filename characters", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		outputDir := "/test-sanitize"

		docs := []api.Document{
			{
				ID:        "doc-1",
				Title:     "Meeting: Notes/Ideas?",
				Content:   "Content",
				CreatedAt: "2024-01-01T00:00:00Z",
				UpdatedAt: "2024-01-01T00:00:00Z",
			},
		}

		err := Write(docs, outputDir, fs)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// List files in directory to see what was created
		files, err := afero.ReadDir(fs, outputDir)
		if err != nil {
			t.Fatalf("failed to read directory: %v", err)
		}

		t.Logf("Created files: %v", files)
		for _, f := range files {
			t.Logf("  - %s", f.Name())
		}

		// Check that file was created with sanitized name
		fileExists, _ := afero.Exists(fs, filepath.Join(outputDir, "Meeting_ Notes_Ideas.md"))
		if !fileExists {
			t.Error("expected sanitized filename to exist")
		}
	})

	t.Run("uses ID when title is empty", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		outputDir := "/test-no-title"

		docs := []api.Document{
			{
				ID:        "doc-abc-123",
				Title:     "",
				Content:   "Content without title",
				CreatedAt: "2024-01-01T00:00:00Z",
				UpdatedAt: "2024-01-01T00:00:00Z",
			},
		}

		err := Write(docs, outputDir, fs)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Check that file was created with ID as filename
		fileExists, _ := afero.Exists(fs, filepath.Join(outputDir, "doc-abc-123.md"))
		if !fileExists {
			t.Error("expected file with ID as name to exist")
		}
	})
}

func TestSanitizeFilename(t *testing.T) {
	t.Run("removes invalid characters", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			title    string
			id       string
			expected string
		}{
			{"Simple Title", "id-1", "Simple Title"},
			{"Title: With Colon", "id-2", "Title_ With Colon"},
			{"Title/With/Slashes", "id-3", "Title_With_Slashes"},
			{"Title?With?Questions", "id-4", "Title_With_Questions"},
			{"Title*With*Stars", "id-5", "Title_With_Stars"},
			{"", "id-6", "id-6"},
			{"   ", "id-7", "id-7"},
		}

		for _, tt := range tests {
			result := sanitizeFilename(tt.title, tt.id)
			if result != tt.expected {
				t.Errorf("sanitizeFilename(%q, %q) = %q, want %q", tt.title, tt.id, result, tt.expected)
			}
		}
	})

	t.Run("limits length to 100 characters", func(t *testing.T) {
		t.Parallel()

		longTitle := strings.Repeat("a", 150)
		result := sanitizeFilename(longTitle, "id-1")

		if len(result) != 100 {
			t.Errorf("expected length 100, got %d", len(result))
		}
	})
}

func TestMakeUnique(t *testing.T) {
	t.Run("returns original filename if not used", func(t *testing.T) {
		t.Parallel()

		used := make(map[string]int)
		result := makeUnique("test", used)

		if result != "test" {
			t.Errorf("expected 'test', got %q", result)
		}
	})

	t.Run("appends number if filename exists", func(t *testing.T) {
		t.Parallel()

		used := map[string]int{"test": 1}
		result := makeUnique("test", used)

		if result != "test_2" {
			t.Errorf("expected 'test_2', got %q", result)
		}
	})
}