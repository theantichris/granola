package cmd

import (
	"testing"

	"github.com/spf13/afero"
)

func TestRunExport(t *testing.T) {
	t.Run("exports Granola documents", func(t *testing.T) {
		appFS = afero.NewMemMapFs()

		afero.WriteFile(appFS, "/test/spabase.json", []byte(`{"workos_tokens": {}}`), 0644)
	})
}
