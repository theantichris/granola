// Package cmd implements the CLI commands for the Granola application.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ErrRootCmd is used when the root command fails to execute.
var ErrRootCmd = errors.New("failed to run granola")

// NewRootCmd creates a new root command with the provided logger and binds flags.
func NewRootCmd(logger *log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "granola",
		Short: "An application for exporting Granola meeting notes.",
		Long:  "An application for exporting Granola meeting notes to Markdown files.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config")); err != nil {
				return fmt.Errorf("%w: %s", ErrRootCmd, err)
			}

			if err := viper.BindPFlag("debug", cmd.PersistentFlags().Lookup("debug")); err != nil {
				return fmt.Errorf("%w: %s", ErrRootCmd, err)
			}

			if err := viper.BindPFlag("supabase", cmd.PersistentFlags().Lookup("supabase")); err != nil {
				return fmt.Errorf("%w: %s", ErrRootCmd, err)
			}

			return nil
		},
	}

	var configFile string
	var debug bool
	var supabaseFile string

	cmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.config.toml)")
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug mode")
	cmd.PersistentFlags().StringVar(&supabaseFile, "supabase", "", "supabase.json file")

	cmd.AddCommand(NewNotesCmd(logger))
	cmd.AddCommand(NewTranscriptsCmd(logger))

	return cmd
}

// Execute creates the logger, initializes configuration, and returns the root command.
func Execute() *cobra.Command {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		Level:           log.WarnLevel,
	})

	cobra.OnInitialize(func() {
		initConfig(logger)
	})

	cmd := NewRootCmd(logger)

	return cmd
}

// initConfig loads env variables and the config file, then updates the logger level if debug mode is enabled.
func initConfig(logger *log.Logger) {
	if err := godotenv.Load(); err != nil {
		logger.Debug(".env file not found, using environment variables")
	} else {
		logger.Debug(".env file loaded successfully")
	}

	configFile := viper.GetString("config")

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".granola")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv()
	_ = viper.BindEnv("debug", "DEBUG_MODE")
	_ = viper.BindEnv("supabase", "SUPABASE_FILE")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Debug("config file not found")
		} else {
			logger.Error("error loading config file", "error", err)
		}
	} else {
		logger.Debug("using config file", "file", viper.ConfigFileUsed())
	}

	if viper.GetBool("debug") {
		logger.SetLevel(log.DebugLevel)
	}
}
