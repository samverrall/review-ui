package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GitClient defines the interface for git operations
type GitClient interface {
	IsGitRepo() (bool, error)
	GetChangedFiles() ([]string, error)
	GetFileDiff(filename string) (string, error)
}

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

// GetChangedFiles returns a list of unstaged changed files and untracked files
func GetChangedFiles() ([]string, error) {
	// Get modified/staged files from git status
	statusFiles, err := getStatusFiles()
	if err != nil {
		return nil, err
	}

	// Get untracked files from git ls-files
	untrackedFiles, err := getUntrackedFiles()
	if err != nil {
		return nil, err
	}

	// Combine and deduplicate
	allFiles := make(map[string]bool)
	for _, file := range statusFiles {
		allFiles[file] = true
	}
	for _, file := range untrackedFiles {
		allFiles[file] = true
	}

	// Convert map to slice
	var files []string
	for file := range allFiles {
		files = append(files, file)
	}

	return files, nil
}

// getStatusFiles returns modified/staged files from git status --porcelain
func getStatusFiles() ([]string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to get status files: %w", err)
	}

	output := strings.TrimSpace(out.String())
	if output == "" {
		return []string{}, nil
	}

	lines := strings.Split(output, "\n")
	var files []string
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		// Trim leading whitespace
		line = strings.TrimLeft(line, " \t")
		if len(line) < 3 {
			continue
		}
		// Skip status codes (first 1-2 characters) and find the filename
		// Status codes are typically 1-2 characters (XY format)
		var filenameStart int
		if len(line) > 1 && line[1] == ' ' {
			// XY file format (e.g., "M file")
			filenameStart = 2
		} else if len(line) > 2 && line[2] == ' ' {
			// XYY file format (though rare, being safe)
			filenameStart = 3
		} else {
			// Fallback: find first space
			spaceIndex := strings.Index(line, " ")
			if spaceIndex == -1 {
				continue
			}
			filenameStart = spaceIndex + 1
		}

		if filenameStart >= len(line) {
			continue
		}

		filename := strings.TrimSpace(line[filenameStart:])
		if filename != "" && !strings.HasSuffix(filename, "/") {
			// Skip directories (entries ending with /)
			files = append(files, filename)
		}
	}

	return files, nil
}

// getUntrackedFiles returns untracked files from git ls-files --others
func getUntrackedFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to get untracked files: %w", err)
	}

	output := strings.TrimSpace(out.String())
	if output == "" {
		return []string{}, nil
	}

	return strings.Split(output, "\n"), nil
}

// isFileTracked checks if a file is tracked by git
func isFileTracked(filename string) (bool, error) {
	cmd := exec.Command("git", "ls-files", "--error-unmatch", filename)
	err := cmd.Run()
	if err == nil {
		return true, nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		// Exit code 1 means file is not tracked (normal case)
		return false, nil
	}
	// Other errors (like git not available, permissions, etc.) should be returned
	return false, fmt.Errorf("failed to check if file is tracked: %w", err)
}

// GetFileDiff returns the unified diff for a specific file
func GetFileDiff(filename string) (string, error) {
	// Check if file is tracked
	tracked, err := isFileTracked(filename)
	if err != nil {
		return "", fmt.Errorf("failed to check if file is tracked: %w", err)
	}

	if tracked {
		// Use git diff for tracked files
		cmd := exec.Command("git", "diff", filename)
		var out bytes.Buffer
		cmd.Stdout = &out

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to get diff for %s: %w", filename, err)
		}

		return out.String(), nil
	} else {
		// For untracked files, read the entire file and format as additions
		content, err := os.ReadFile(filename)
		if err != nil {
			return "", fmt.Errorf("failed to read untracked file %s: %w", filename, err)
		}

		// Format as git diff for new file
		lines := strings.Split(string(content), "\n")
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1] // Remove trailing empty line if present
		}

		var diff strings.Builder
		diff.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", filename, filename))
		diff.WriteString("new file mode 100644\n")
		diff.WriteString("index 0000000..e69de29\n")
		diff.WriteString("--- /dev/null\n")
		diff.WriteString(fmt.Sprintf("+++ b/%s\n", filename))
		diff.WriteString(fmt.Sprintf("@@ -0,0 +1,%d @@\n", len(lines)))

		for _, line := range lines {
			diff.WriteString("+" + line + "\n")
		}

		return diff.String(), nil
	}
}
