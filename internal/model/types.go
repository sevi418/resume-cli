package model

import "fmt"

type Education struct {
	School         string `json:"school"`
	Major          string `json:"major"`
	Degree         string `json:"degree"`
	GraduationTime string `json:"graduation_time"`
}

type Resume struct {
	Name      string      `json:"name"`
	Phone     string      `json:"phone"`
	Email     string      `json:"email"`
	City      string      `json:"city"`
	Education []Education `json:"education"`
	Skills    []string    `json:"skills"`
}

func (r Resume) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("resume name is empty")
	}
	return nil
}

type Score struct {
	OverallScore       int      `json:"overall_score"`
	SkillScore         int      `json:"skill_score"`
	ExperienceScore    int      `json:"experience_score"`
	EducationScore     int      `json:"education_score"`
	Comment            string   `json:"comment"`
	InterviewQuestions []string `json:"interview_questions"`
}

func (s Score) Validate() error {
	scores := map[string]int{
		"overall_score":    s.OverallScore,
		"skill_score":      s.SkillScore,
		"experience_score": s.ExperienceScore,
		"education_score":  s.EducationScore,
	}

	for name, score := range scores {
		if score < 0 || score > 100 {
			return fmt.Errorf("%s must be between 0 and 100", name)
		}
	}
	if s.Comment == "" {
		return fmt.Errorf("score comment is empty")
	}
	return nil
}
