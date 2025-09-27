package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Logger *log.Logger

// RootCmd is the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "granola",
	Short: "An application for exporting Granola notes.",
	Long:  "An application for exporting Granola notes to Markdown files.",
}

// Execute adds initialization.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		Logger.Error("error running command", "error", err)
		os.Exit(1)
	}
}

// init sets and binds flags.
func init() {
	cobra.OnInitialize(initConfig)

	var configFile string
	var debug bool
	var supabase string

	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.config.toml)")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug mode")
	RootCmd.PersistentFlags().StringVar(&supabase, "supabase", "", "path to supabase.json")

	_ = viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("supabase", RootCmd.PersistentFlags().Lookup("supabase"))
}

// initConfig loads env variables and the config file.
func initConfig() {
	initLogger(false)

	if err := godotenv.Load(); err != nil {
		Logger.Debug(".env file not found, using environment variables")
	} else {
		Logger.Debug(".env file loaded successfully")
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
			Logger.Debug("config file not found")
		} else {
			Logger.Error("error loading config file", "error", err)
		}
	} else {
		Logger.Debug("using config file", "file", viper.ConfigFileUsed())
	}

	if viper.GetBool("debug") {
		initLogger(true)
	}
}

// initLogger initializes the logger.
func initLogger(debug bool) {
	Logger = log.New(os.Stderr)
	Logger.SetReportCaller(true)
	Logger.SetReportTimestamp(true)

	if debug {
		Logger.SetLevel(log.DebugLevel)
	} else {
		Logger.SetLevel(log.WarnLevel)
	}
}
