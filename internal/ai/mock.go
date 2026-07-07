package ai

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/sevi418/resume-cli/internal/model"
)

type MockAIService struct{}

func NewMockAIService() *MockAIService {
	slog.Debug("configured mock AI service")
	return &MockAIService{}
}

func (s *MockAIService) ExtractResume(ctx context.Context, text string) (*model.Resume, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("resume text is empty")
	}
	if hasCJK(text) {
		return &model.Resume{
			Name:  "张三",
			Phone: "13800138000",
			Email: "zhangsan@example.com",
			City:  "北京",
			Education: []model.Education{
				{
					School:         "清华大学",
					Major:          "计算机科学与技术",
					Degree:         "本科",
					GraduationTime: "2020-06",
				},
			},
			Skills: []string{"Go", "React", "PostgreSQL", "Docker", "Kubernetes", "微服务"},
		}, nil
	}

	return &model.Resume{
		Name:  "John Smith",
		Phone: "+1-555-0100",
		Email: "john.smith@example.com",
		City:  "San Francisco",
		Education: []model.Education{
			{
				School:         "University of California, Berkeley",
				Major:          "Computer Science",
				Degree:         "Bachelor",
				GraduationTime: "2018-05",
			},
		},
		Skills: []string{"Go", "Python", "Kubernetes", "gRPC", "PostgreSQL", "AWS"},
	}, nil
}

func (s *MockAIService) ScoreMatch(ctx context.Context, resumeText, jdText string) (*model.Score, error) {
	if strings.TrimSpace(resumeText) == "" {
		return nil, fmt.Errorf("resume text is empty")
	}
	if strings.TrimSpace(jdText) == "" {
		return nil, fmt.Errorf("JD text is empty")
	}

	if hasCJK(resumeText) && strings.Contains(resumeText, "产品") {
		return &model.Score{
			OverallScore:    42,
			SkillScore:      35,
			ExperienceScore: 45,
			EducationScore:  70,
			Comment:         "候选人履历偏产品管理，与高级全栈工程师岗位的 Go、React、微服务和工程实践要求匹配度较低。",
			InterviewQuestions: []string{
				"请说明你是否有实际编码和系统设计经验。",
				"你在过往项目中是否直接负责过 Go 或 React 的开发工作？",
			},
		}, nil
	}

	if hasCJK(jdText) {
		return &model.Score{
			OverallScore:    86,
			SkillScore:      90,
			ExperienceScore: 85,
			EducationScore:  82,
			Comment:         "候选人的 Go、React、微服务和容器化经验与岗位要求高度匹配，适合进入下一轮技术面试。",
			InterviewQuestions: []string{
				"请介绍一个你主导的 Go 微服务项目，并说明服务拆分依据。",
				"你如何在 React 项目中处理复杂状态和性能问题？",
				"请描述一次 CI/CD 或 Kubernetes 落地中的关键问题。",
			},
		}, nil
	}

	return &model.Score{
		OverallScore:    84,
		SkillScore:      88,
		ExperienceScore: 83,
		EducationScore:  80,
		Comment:         "The candidate matches the backend/full-stack requirements well and has relevant distributed systems experience.",
		InterviewQuestions: []string{
			"Describe a Go service you designed and operated in production.",
			"How have you handled API compatibility across service versions?",
		},
	}, nil
}

func hasCJK(s string) bool {
	for _, r := range s {
		if r >= '\u4e00' && r <= '\u9fff' {
			return true
		}
	}
	return false
}
