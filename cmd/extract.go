package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	pdfparser "github.com/sevi418/resume-cli/internal/pdf"
	"github.com/sevi418/resume-cli/internal/util"
	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract <pdf_path>",
	Short: "Extract structured resume information with AI",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		slog.Debug("extract started", "pdf", args[0])

		text, err := pdfparser.ParsePDF(args[0])
		if err != nil {
			return err
		}

		service, err := newAIService()
		if err != nil {
			return err
		}

		resume, err := service.ExtractResume(context.Background(), text)
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(resume, "", "  ")
		if err != nil {
			return fmt.Errorf("encode resume JSON: %w", err)
		}
		data = append(data, '\n')

		logCompleted := func() {
			slog.Debug(
				"extract completed",
				"resume_chars", len([]rune(text)),
				"education_count", len(resume.Education),
				"skills_count", len(resume.Skills),
				"bytes", len(data),
				"output", outputTarget(cfg.output),
				"elapsed", time.Since(start),
			)
		}
		if cfg.output == "" {
			logCompleted()
		}
		if err := util.WriteOutput(cfg.output, data); err != nil {
			return err
		}
		if cfg.output != "" {
			logCompleted()
			fmt.Fprintf(cmd.ErrOrStderr(), "wrote extracted resume to %s\n", cfg.output)
		}
		return nil
	},
}
