package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	pdfparser "github.com/sevi418/resume-cli/internal/pdf"
	"github.com/sevi418/resume-cli/internal/util"
	"github.com/spf13/cobra"
)

var jdPath string

var scoreCmd = &cobra.Command{
	Use:   "score <pdf_path>",
	Short: "Score resume match against a JD",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(jdPath) == "" {
			return fmt.Errorf("--jd is required")
		}
		start := time.Now()
		slog.Debug("score started", "pdf", args[0], "jd", jdPath)

		resumeText, err := pdfparser.ParsePDF(args[0])
		if err != nil {
			return err
		}

		jdText, err := readJD(jdPath)
		if err != nil {
			return err
		}

		service, err := newAIService()
		if err != nil {
			return err
		}

		score, err := service.ScoreMatch(context.Background(), resumeText, jdText)
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(score, "", "  ")
		if err != nil {
			return fmt.Errorf("encode score JSON: %w", err)
		}
		data = append(data, '\n')

		logCompleted := func() {
			slog.Debug(
				"score completed",
				"resume_chars", len([]rune(resumeText)),
				"jd_chars", len([]rune(jdText)),
				"overall_score", score.OverallScore,
				"skill_score", score.SkillScore,
				"experience_score", score.ExperienceScore,
				"education_score", score.EducationScore,
				"interview_question_count", len(score.InterviewQuestions),
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
			fmt.Fprintf(cmd.ErrOrStderr(), "wrote score to %s\n", cfg.output)
		}
		return nil
	},
}

func readJD(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("JD file does not exist: %s", path)
		}
		return "", fmt.Errorf("read JD file %q: %w", path, err)
	}

	text := string(data)
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("JD file is empty: %s", path)
	}
	slog.Debug("JD loaded", "path", path, "bytes", len(data))
	return text, nil
}

func init() {
	scoreCmd.Flags().StringVar(&jdPath, "jd", "", "path to JD text file")
}
