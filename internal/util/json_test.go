package util

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepairJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "markdown fenced json",
			raw:  "```json\n{\"name\":\"张三\"}\n```",
			want: "{\"name\":\"张三\"}",
		},
		{
			name: "extracts object from prose",
			raw:  "result:\n{\"name\":\"John\"}\nthanks",
			want: "{\"name\":\"John\"}",
		},
		{
			name: "removes trailing commas",
			raw:  "{\"skills\":[\"Go\",],\"name\":\"张三\",}",
			want: "{\"skills\":[\"Go\"],\"name\":\"张三\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := RepairJSON(tt.raw)
			require.Equal(t, tt.want, got)
			require.True(t, json.Valid([]byte(got)))
		})
	}
}
