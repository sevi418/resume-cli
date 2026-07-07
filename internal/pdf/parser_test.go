package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePDF_TextFixtures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
	}{
		{
			name: "chinese resume",
			path: "../../testdata/01_zh_fullstack_senior.pdf",
		},
		{
			name: "english resume",
			path: "../../testdata/04_en_senior_backend.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			text, err := ParsePDF(tt.path)
			require.NoError(t, err)
			require.NotEmpty(t, strings.TrimSpace(text))
		})
	}
}

func TestParsePDF_Errors(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	nonPDF := filepath.Join(tempDir, "resume.txt")
	require.NoError(t, os.WriteFile(nonPDF, []byte("not pdf"), 0o644))

	brokenPDF := filepath.Join(tempDir, "broken.pdf")
	require.NoError(t, os.WriteFile(brokenPDF, []byte("not a real pdf"), 0o644))

	emptyPDF := filepath.Join(tempDir, "empty.pdf")
	require.NoError(t, writeBlankPDF(emptyPDF))

	tests := []struct {
		name    string
		path    string
		wantErr string
	}{
		{
			name:    "missing file",
			path:    filepath.Join(tempDir, "missing.pdf"),
			wantErr: "does not exist",
		},
		{
			name:    "non pdf file",
			path:    nonPDF,
			wantErr: "not a PDF",
		},
		{
			name:    "broken pdf",
			path:    brokenPDF,
			wantErr: "read pdf",
		},
		{
			name:    "empty text pdf",
			path:    emptyPDF,
			wantErr: "text is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := ParsePDF(tt.path)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func writeBlankPDF(path string) error {
	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] >>",
	}

	var b strings.Builder
	b.WriteString("%PDF-1.4\n")
	offsets := []int{0}
	for i, obj := range objects {
		offsets = append(offsets, b.Len())
		fmt.Fprintf(&b, "%d 0 obj\n%s\nendobj\n", i+1, obj)
	}

	xrefOffset := b.Len()
	fmt.Fprintf(&b, "xref\n0 %d\n", len(objects)+1)
	b.WriteString("0000000000 65535 f \n")
	for _, offset := range offsets[1:] {
		fmt.Fprintf(&b, "%010d 00000 n \n", offset)
	}
	fmt.Fprintf(&b, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, xrefOffset)

	return os.WriteFile(path, []byte(b.String()), 0o644)
}
