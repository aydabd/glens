package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "glens",
	Short: "OpenAPI Integration Test Generator with AI Models",
	Long: `A powerful tool that analyzes OpenAPI specifications and generates
integration tests using multiple AI models (OpenAI GPT, Anthropic Sonnet, Google Flash).
Creates GitHub issues for each endpoint and generates comprehensive test reports.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.glens.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug logging")
	rootCmd.PersistentFlags().String("log-format", "console", "log format (console or json)")

	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind debug flag")
	}
	if err := viper.BindPFlag("log_format", rootCmd.PersistentFlags().Lookup("log-format")); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind log-format flag")
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".glens")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("")

	// Bind environment variables explicitly for GitHub
	_ = viper.BindEnv("github.token", "GITHUB_TOKEN")
	_ = viper.BindEnv("github.repository", "GITHUB_REPOSITORY")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	setupLogging()
}

func setupLogging() {
	logFormat := viper.GetString("log_format")
	debug := viper.GetBool("debug")

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if logFormat == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled")
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
