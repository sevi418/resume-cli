package ai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockAIService_ExtractResume(t *testing.T) {
	t.Parallel()

	service := NewMockAIService()

	zh, err := service.ExtractResume(context.Background(), "张三 高级全栈工程师 Go React")
	require.NoError(t, err)
	require.Equal(t, "张三", zh.Name)
	require.Contains(t, zh.Skills, "Go")

	en, err := service.ExtractResume(context.Background(), "John Smith senior backend engineer")
	require.NoError(t, err)
	require.Equal(t, "John Smith", en.Name)
	require.Contains(t, en.Skills, "gRPC")
}

func TestMockAIService_ScoreMatch(t *testing.T) {
	t.Parallel()

	service := NewMockAIService()

	high, err := service.ScoreMatch(context.Background(), "张三 Go React 微服务", "高级全栈工程师 Go React")
	require.NoError(t, err)
	require.GreaterOrEqual(t, high.OverallScore, 80)

	low, err := service.ScoreMatch(context.Background(), "产品经理 用户增长", "高级全栈工程师 Go React")
	require.NoError(t, err)
	require.Less(t, low.OverallScore, high.OverallScore)
}
