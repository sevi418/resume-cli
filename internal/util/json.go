package util

import (
	"regexp"
	"strings"
)

var trailingCommaRE = regexp.MustCompile(`,\s*([}\]])`)

func RepairJSON(raw string) string {
	s := strings.TrimSpace(raw)
	s = stripMarkdownFence(s)

	if start := strings.Index(s, "{"); start >= 0 {
		if end := strings.LastIndex(s, "}"); end >= start {
			s = s[start : end+1]
		}
	}

	s = trailingCommaRE.ReplaceAllString(s, "$1")
	return strings.TrimSpace(s)
}

func stripMarkdownFence(s string) string {
	if !strings.HasPrefix(s, "```") {
		return s
	}

	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		return s
	}
	if strings.HasPrefix(strings.TrimSpace(lines[len(lines)-1]), "```") {
		lines = lines[1 : len(lines)-1]
		return strings.TrimSpace(strings.Join(lines, "\n"))
	}
	return s
}
