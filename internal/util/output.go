package util

import (
	"fmt"
	"os"
)

func WriteOutput(path string, data []byte) error {
	if path == "" {
		_, err := os.Stdout.Write(data)
		return err
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write output %q: %w", path, err)
	}
	return nil
}
