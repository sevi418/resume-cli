package pdf

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	ledongpdf "github.com/ledongthuc/pdf"
)

func ParsePDF(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("pdf path is required")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("pdf file does not exist: %s", path)
		}
		return "", fmt.Errorf("stat pdf file %q: %w", path, err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("pdf path is a directory: %s", path)
	}
	if strings.ToLower(filepath.Ext(path)) != ".pdf" {
		return "", fmt.Errorf("file is not a PDF: %s", path)
	}

	start := time.Now()
	slog.Debug(
		"parsing pdf",
		"path", path,
		"bytes", info.Size(),
	)

	text, err := parseWithGoPDF(path)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(text) != "" {
		slog.Debug(
			"parsed pdf with Go parser",
			"path", path,
			"chars", len([]rune(text)),
			"elapsed", time.Since(start),
		)
		return text, nil
	}

	slog.Debug("Go PDF parser returned empty text, trying pdftotext", "path", path)
	fallbackText, fallbackErr := parseWithPDFToText(path)
	if fallbackErr == nil && strings.TrimSpace(fallbackText) != "" {
		slog.Debug(
			"parsed pdf with pdftotext fallback",
			"path", path,
			"chars", len([]rune(fallbackText)),
			"elapsed", time.Since(start),
		)
		return fallbackText, nil
	}
	if fallbackErr != nil {
		slog.Debug("pdftotext fallback failed", "error", fallbackErr)
	}

	return "", fmt.Errorf("pdf text is empty; this file may be a scanned PDF and require OCR")
}

func parseWithGoPDF(path string) (string, error) {
	file, reader, err := ledongpdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("read pdf %q: %w", path, err)
	}
	defer file.Close()

	stream, err := reader.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("extract pdf text %q: %w", path, err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, stream); err != nil {
		return "", fmt.Errorf("copy pdf text %q: %w", path, err)
	}
	return buf.String(), nil
}

func parseWithPDFToText(path string) (string, error) {
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return "", fmt.Errorf("pdftotext is not available: %w", err)
	}

	cmd := exec.Command("pdftotext", "-layout", path, "-")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return "", fmt.Errorf("pdftotext failed: %s: %w", msg, err)
		}
		return "", fmt.Errorf("pdftotext failed: %w", err)
	}
	return string(out), nil
}
