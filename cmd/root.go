package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/sevi418/resume-cli/internal/ai"
	"github.com/sevi418/resume-cli/internal/util"
	"github.com/spf13/cobra"
)

var cfg config

type config struct {
	output  string
	mock    bool
	verbose bool
}

var rootCmd = &cobra.Command{
	Use:           "resume-cli",
	Short:         "AI resume parser and JD scorer",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupLogger(cfg.verbose)
		slog.Debug(
			"command started",
			"command", cmd.CommandPath(),
			"mock", cfg.mock,
			"output", outputTarget(cfg.output),
		)
	},
}

func Execute() {
	setupLogger(hasVerboseFlag(os.Args[1:]))

	if err := util.LoadDotEnv(".env"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setupLogger(verbose bool) {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: compactLogAttrs,
	})))
}

func hasVerboseFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--verbose" || arg == "-v" || strings.HasPrefix(arg, "--verbose=") {
			return true
		}
	}
	return false
}

func compactLogAttrs(groups []string, attr slog.Attr) slog.Attr {
	if attr.Key == slog.TimeKey {
		return slog.Attr{}
	}
	return attr
}

func outputTarget(path string) string {
	if path == "" {
		return "stdout"
	}
	return path
}

func newAIService() (ai.AIService, error) {
	if cfg.mock {
		return ai.NewMockAIService(), nil
	}
	return ai.NewRealAIServiceFromEnv()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfg.output, "output", "o", "", "write result to file")
	rootCmd.PersistentFlags().BoolVar(&cfg.mock, "mock", false, "use mock AI responses without an API key")
	rootCmd.PersistentFlags().BoolVarP(&cfg.verbose, "verbose", "v", false, "enable verbose logs")

	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(scoreCmd)
}
