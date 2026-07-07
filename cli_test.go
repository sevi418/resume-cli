package main_test

import (
	"encoding/json"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCLIParseExtractScoreMock(t *testing.T) {
	parseOut := runCLI(t, "parse", "testdata/01_zh_fullstack_senior.pdf")
	require.NotEmpty(t, parseOut)

	extractOut := runCLI(t, "extract", "testdata/01_zh_fullstack_senior.pdf", "--mock")
	var resume struct {
		Name string `json:"name"`
	}
	require.NoError(t, json.Unmarshal([]byte(extractOut), &resume))
	require.Equal(t, "张三", resume.Name)

	scoreOut := runCLI(t, "score", "testdata/01_zh_fullstack_senior.pdf", "--jd", "testdata/sample_jd.txt", "--mock")
	var score struct {
		OverallScore int `json:"overall_score"`
	}
	require.NoError(t, json.Unmarshal([]byte(scoreOut), &score))
	require.GreaterOrEqual(t, score.OverallScore, 80)
}

func runCLI(t *testing.T, args ...string) string {
	t.Helper()

	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	return string(out)
}
