package prompt

import (
	"io/fs"
	"os"
	"path/filepath"

	prompts "server/prompts"
)

// Read returns the content of the named prompt file.
// It reads from the embedded filesystem first, then falls back to filesystem search
// (for development convenience when iterating on prompt wording without rebuilding).
func Read(name string) ([]byte, error) {
	// 1. Embedded FS (bundled into binary at compile time)
	data, err := fs.ReadFile(prompts.FS, name)
	if err == nil {
		return data, nil
	}

	// 2. Filesystem fallback for development
	candidates := make([]string, 0, 3)
	if exePath, exeErr := os.Executable(); exeErr == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(candidates,
			filepath.Join(exeDir, "prompts", name),
			filepath.Join(exeDir, "..", "prompts", name),
		)
	}
	candidates = append(candidates, filepath.Join("prompts", name))

	var lastErr error
	for _, p := range candidates {
		data, err := os.ReadFile(p)
		if err == nil {
			return data, nil
		}
		lastErr = err
	}
	return nil, lastErr
}
