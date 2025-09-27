package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	Debug      bool
	Logger     *log.Logger
	EnvVar     string
	Supabase   string
)

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
		Logger.Error("error running command")
		os.Exit(1)
	}
}

// init sets and binds flags.
func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.config.toml)")
	RootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "enable debug mode")
	RootCmd.PersistentFlags().StringVar(&Supabase, "supabase", "", "path to supabase.json")

	_ = viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("supabase", RootCmd.PersistentFlags().Lookup("supabase"))
}

// initConfig loads env variables and the config file.
func initConfig() {
	initLogger()

	if err := godotenv.Load(); err != nil {
		Logger.Debug(".env file not found, using environment variables")
	} else {
		Logger.Debug(".env file loaded successfully")
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".config")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv()
	_ = viper.BindEnv("envVar", "ENV_VAR")

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
		Debug = true
		initLogger()
	}
}

// initLogger initializes the logger.
func initLogger() {
	Logger = log.New(os.Stderr)
	Logger.SetReportCaller(true)
	Logger.SetReportTimestamp(true)

	if Debug {
		Logger.SetLevel(log.DebugLevel)
	} else {
		Logger.SetLevel(log.WarnLevel)
	}
}
