package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// IsGitRepo checks if the current directory is inside a git repository
func IsGitRepo() (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 128 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetChangedFiles returns a list of unstaged changed files
func GetChangedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	output := strings.TrimSpace(out.String())
	if output == "" {
		return []string{}, nil
	}

	files := strings.Split(output, "\n")
	return files, nil
}

// GetFileDiff returns the unified diff for a specific file
func GetFileDiff(filename string) (string, error) {
	cmd := exec.Command("git", "diff", filename)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get diff for %s: %w", filename, err)
	}

	return out.String(), nil
}
