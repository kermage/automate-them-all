package cloak

import (
	"os/exec"
	"testing"
)

func setupTestGitRepo(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	initCmd := exec.Command("git", "init", tmpDir)
	if err := initCmd.Run(); err != nil {
		t.Fatalf("failed to init git repository: %v", err)
	}
	return tmpDir
}
