package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"glens/tools/glens/pkg/logging"
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
		fmt.Fprintln(os.Stderr, "failed to bind debug flag:", err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("log_format", rootCmd.PersistentFlags().Lookup("log-format")); err != nil {
		fmt.Fprintln(os.Stderr, "failed to bind log-format flag:", err)
		os.Exit(1)
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

	level := logging.LevelInfo
	if debug {
		level = logging.LevelDebug
	}

	format := logging.FormatJSON
	if logFormat == "console" {
		format = logging.FormatConsole
	}

	logging.Setup(logging.Config{Level: level, Format: format})
}
