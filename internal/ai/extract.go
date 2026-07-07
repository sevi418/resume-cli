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

func (s *RealAIService) ExtractResume(ctx context.Context, text string) (*model.Resume, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("resume text is empty")
	}

	content, err := s.chatJSON(ctx, extractSystemPrompt, fmt.Sprintf(extractUserPrompt, text))
	if err != nil {
		return nil, err
	}

	var resume model.Resume
	if err := json.Unmarshal([]byte(util.RepairJSON(content)), &resume); err != nil {
		return nil, fmt.Errorf("parse AI resume JSON: %w", err)
	}
	if err := resume.Validate(); err != nil {
		return nil, err
	}
	slog.Debug(
		"AI resume JSON validated",
		"response_chars", len([]rune(content)),
		"name_present", resume.Name != "",
		"education_count", len(resume.Education),
		"skills_count", len(resume.Skills),
	)
	return &resume, nil
}

const extractSystemPrompt = `You are a resume parsing assistant.
Return only a valid JSON object. Do not include markdown or explanations.
The JSON keys must be exactly: name, phone, email, city, education, skills.
Keep extracted values in the same language as the resume.`

const extractUserPrompt = `Extract structured information from this resume text.

Schema:
{
  "name": "string",
  "phone": "string",
  "email": "string",
  "city": "string",
  "education": [
    {
      "school": "string",
      "major": "string",
      "degree": "string",
      "graduation_time": "string"
    }
  ],
  "skills": ["string"]
}

If a field is unavailable, use an empty string or empty array.

Resume text:
%s`
