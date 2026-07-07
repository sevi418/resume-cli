package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResumeValidate(t *testing.T) {
	t.Parallel()

	require.NoError(t, Resume{Name: "张三"}.Validate())
	require.Error(t, Resume{}.Validate())
}

func TestScoreValidate(t *testing.T) {
	t.Parallel()

	valid := Score{
		OverallScore:       80,
		SkillScore:         81,
		ExperienceScore:    82,
		EducationScore:     83,
		Comment:            "good match",
		InterviewQuestions: []string{"question"},
	}
	require.NoError(t, valid.Validate())

	invalid := valid
	invalid.OverallScore = 101
	require.Error(t, invalid.Validate())

	invalid = valid
	invalid.Comment = ""
	require.Error(t, invalid.Validate())
}
