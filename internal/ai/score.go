package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/sevi418/resume-cli/internal/model"
	"github.com/sevi418/resume-cli/internal/util"
)

func (s *RealAIService) ScoreMatch(ctx context.Context, resumeText, jdText string) (*model.Score, error) {
	if strings.TrimSpace(resumeText) == "" {
		return nil, fmt.Errorf("resume text is empty")
	}
	if strings.TrimSpace(jdText) == "" {
		return nil, fmt.Errorf("JD text is empty")
	}

	content, err := s.chatJSON(ctx, scoreSystemPrompt, fmt.Sprintf(scoreUserPrompt, jdText, resumeText))
	if err != nil {
		return nil, err
	}

	var score model.Score
	if err := json.Unmarshal([]byte(util.RepairJSON(content)), &score); err != nil {
		return nil, fmt.Errorf("parse AI score JSON: %w", err)
	}
	if err := score.Validate(); err != nil {
		return nil, err
	}
	slog.Debug(
		"AI score JSON validated",
		"response_chars", len([]rune(content)),
		"overall_score", score.OverallScore,
		"interview_question_count", len(score.InterviewQuestions),
	)
	return &score, nil
}

const scoreSystemPrompt = `You are a recruiting evaluation assistant.
Return only a valid JSON object. Do not include markdown or explanations.
All scores must be integers from 0 to 100.
Use the main language of the JD for comment and interview_questions.`

const scoreUserPrompt = `Compare the resume against the JD and return this JSON schema:
{
  "overall_score": 0,
  "skill_score": 0,
  "experience_score": 0,
  "education_score": 0,
  "comment": "string",
  "interview_questions": ["string"]
}

JD:
%s

Resume text:
%s`
