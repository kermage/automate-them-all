package cloak

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const (
	DirPerm  = 0o755
	FilePerm = 0o644
)

func getAbsoluteGitDir(baseDir string) (string, error) {
	cmdGit := exec.Command("git", "rev-parse", "--absolute-git-dir")
	cmdGit.Dir = baseDir
	out, err := cmdGit.CombinedOutput()
	if err != nil {
		var execErr *exec.Error
		if errors.As(err, &execErr) {
			return "", fmt.Errorf("git is not available")
		}

		msg := strings.TrimSpace(string(out))
		if strings.Contains(msg, "not a git repository") {
			return "", fmt.Errorf("not a git repository")
		}

		return "", fmt.Errorf("failed to resolve git directory: %s", msg)
	}

	return strings.TrimSpace(string(out)), nil
}

func hasIgnoreTarget(content, target string) bool {
	return content == target ||
		strings.HasPrefix(content, target+"\n") ||
		strings.HasSuffix(content, "\n"+target) ||
		strings.Contains(content, "\n"+target+"\n")
}

func addIgnoreTarget(content, target string) string {
	if hasIgnoreTarget(content, target) {
		return content
	}

	newContent := content
	if len(newContent) > 0 && !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}
	newContent += target + "\n"
	return newContent
}

func removeIgnoreTarget(content, target string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if line == target {
			return strings.Join(append(lines[:i], lines[i+1:]...), "\n")
		}
	}
	return content
}
