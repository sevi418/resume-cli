package cmd

import (
	"fmt"
	"log/slog"
	"time"

	pdfparser "github.com/sevi418/resume-cli/internal/pdf"
	"github.com/sevi418/resume-cli/internal/util"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse <pdf_path>",
	Short: "Extract plain text from a PDF resume",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		slog.Debug("parse started", "pdf", args[0])

		text, err := pdfparser.ParsePDF(args[0])
		if err != nil {
			return err
		}

		if text != "" && text[len(text)-1] != '\n' {
			text += "\n"
		}
		logCompleted := func() {
			slog.Debug(
				"parse completed",
				"chars", len([]rune(text)),
				"bytes", len([]byte(text)),
				"output", outputTarget(cfg.output),
				"elapsed", time.Since(start),
			)
		}
		if cfg.output == "" {
			logCompleted()
		}
		if err := util.WriteOutput(cfg.output, []byte(text)); err != nil {
			return err
		}
		if cfg.output != "" {
			logCompleted()
			fmt.Fprintf(cmd.ErrOrStderr(), "wrote parsed text to %s\n", cfg.output)
		}
		return nil
	},
}
