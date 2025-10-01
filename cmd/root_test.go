package cmd

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func TestNewRootCmd(t *testing.T) {
	t.Run("creates root command with correct configuration", func(t *testing.T) {
		logger := log.New(io.Discard)
		cmd := NewRootCmd(logger)

		if cmd == nil {
			t.Fatal("expected command to be created, got nil")
		}

		if cmd.Use != "granola" {
			t.Errorf("expected Use to be 'granola', got %s", cmd.Use)
		}

		if cmd.Short != "An application for exporting Granola meeting notes." {
			t.Errorf("unexpected Short description: %s", cmd.Short)
		}

		if cmd.Long != "An application for exporting Granola meeting notes to Markdown files." {
			t.Errorf("unexpected Long description: %s", cmd.Long)
		}

		// Check persistent flags are set
		configFlag := cmd.PersistentFlags().Lookup("config")
		if configFlag == nil {
			t.Error("expected config flag to be set")
		}

		debugFlag := cmd.PersistentFlags().Lookup("debug")
		if debugFlag == nil {
			t.Error("expected debug flag to be set")
		}

		supabaseFlag := cmd.PersistentFlags().Lookup("supabase")
		if supabaseFlag == nil {
			t.Error("expected supabase flag to be set")
		}

		// Check PreRunE is set
		if cmd.PreRunE == nil {
			t.Error("expected PreRunE to be set")
		}

		// Check notes subcommand is added
		found := false
		for _, subCmd := range cmd.Commands() {
			if subCmd.Use == "notes" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected notes subcommand to be added")
		}
	})
}

func TestInitConfig(t *testing.T) {
	t.Run("updates logger level when debug is enabled", func(t *testing.T) {
		logger := log.NewWithOptions(io.Discard, log.Options{
			ReportCaller:    false,
			ReportTimestamp: false,
			Level:           log.WarnLevel,
		})

		viper.Reset()

		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, ".granola.toml")
		configContent := `debug = true`

		if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
			t.Fatalf("failed to write to test configFile: %v", err)
		}

		viper.Set("config", configFile)

		initConfig(logger)

		if logger.GetLevel() != log.DebugLevel {
			t.Errorf("expected logger level to be DebugLevel, got %v", logger.GetLevel())
		}

		if !viper.GetBool("debug") {
			t.Error("expected debug mode to be enabled in viper")
		}
	})

	t.Run("loads environment variables from .env file", func(t *testing.T) {
		logger := log.New(io.Discard)

		viper.Reset()

		tmpDir := t.TempDir()
		envFile := filepath.Join(tmpDir, ".env")
		envContent := `DEBUG_MODE=true
SUPABASE_FILE=/path/to/supabase.json`

		if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
			t.Fatalf("failed to write to test .env file: %v", err)
		}

		oldWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get the current working directory: %v", err)
		}

		defer func() {
			if err := os.Chdir(oldWd); err != nil {
				t.Fatalf("failed to change to old working directory: %v", err)
			}
		}()

		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to change to temp directory: %v", err)
		}

		initConfig(logger)
		if !viper.GetBool("debug") {
			t.Error("expected DEBUG_MODE from .env to be loaded")
		}

		if viper.GetString("supabase") != "/path/to/supabase.json" {
			t.Errorf("expected SUPABASE_FILE from .env to be loaded, got %s", viper.GetString("supabase"))
		}
	})
}
