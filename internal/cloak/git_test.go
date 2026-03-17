package cloak

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetAbsoluteGitDir(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) string
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Success - Repo Root",
			setup: func(t *testing.T) string {
				return setupTestGitRepo(t)
			},
			wantErr: false,
		},
		{
			name: "Success - Subdirectory",
			setup: func(t *testing.T) string {
				tmpDir := setupTestGitRepo(t)
				subDir := filepath.Join(tmpDir, "nested", "dir")
				if err := os.MkdirAll(subDir, DirPerm); err != nil {
					t.Fatalf("failed to create subdir: %v", err)
				}
				return subDir
			},
			wantErr: false,
		},
		{
			name: "Failure - Not a Git Repository",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr:  true,
			errorMsg: "not a git repository",
		},
		{
			name: "Failure - Git Unavailable",
			setup: func(t *testing.T) string {
				oldPath := os.Getenv("PATH")
				if err := os.Setenv("PATH", ""); err != nil {
					t.Fatalf("failed to clear PATH: %v", err)
				}
				t.Cleanup(func() {
					_ = os.Setenv("PATH", oldPath)
				})
				return t.TempDir()
			},
			wantErr:  true,
			errorMsg: "git is not available",
		},
		{
			name: "Failure - Non-existent Directory",
			setup: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "non-existent")
			},
			wantErr:  true,
			errorMsg: "failed to resolve git directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitDir, err := getAbsoluteGitDir(tt.setup(t))

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errorMsg)
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Fatalf("expected error containing %q, got %v", tt.errorMsg, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gitDir == "" {
				t.Error("expected non-empty git directory")
			}

			if !filepath.IsAbs(gitDir) {
				t.Errorf("expected absolute path, got %s", gitDir)
			}
		})
	}
}

func TestHasIgnoreTarget(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		target   string
		expected bool
	}{
		{"Empty content", "", ".oops/", false},
		{"Exact match", ".oops/", ".oops/", true},
		{"Prefix match", ".oops/\nfoo", ".oops/", true},
		{"Suffix match", "foo\n.oops/", ".oops/", true},
		{"Middle match", "foo\n.oops/\nbar", ".oops/", true},
		{"No match substring", "not.oops/", ".oops/", false},
		{"No match line", "foo\nbar", ".oops/", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasIgnoreTarget(tt.content, tt.target); got != tt.expected {
				t.Errorf("hasIgnoreTarget(%q, %q) = %v, want %v", tt.content, tt.target, got, tt.expected)
			}
		})
	}
}

func TestAddIgnoreTarget(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		target   string
		expected string
	}{
		{"Add to empty", "", ".oops/", ".oops/\n"},
		{"Already exists", ".oops/\n", ".oops/", ".oops/\n"},
		{"Add to existing with newline", "foo\n", ".oops/", "foo\n.oops/\n"},
		{"Add to existing without newline", "foo", ".oops/", "foo\n.oops/\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addIgnoreTarget(tt.content, tt.target); got != tt.expected {
				t.Errorf("addIgnoreTarget(%q, %q) = %q, want %q", tt.content, tt.target, got, tt.expected)
			}
		})
	}
}

func TestRemoveIgnoreTarget(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		target   string
		expected string
	}{
		{"Remove only line", ".oops/", ".oops/", ""},
		{"Remove only line with newline", ".oops/\n", ".oops/", ""},
		{"Remove first line", ".oops/\nfoo", ".oops/", "foo"},
		{"Remove last line", "foo\n.oops/", ".oops/", "foo"},
		{"Remove middle line", "foo\n.oops/\nbar", ".oops/", "foo\nbar"},
		{"Do not remove substring in another line", "frontend/node_modules\nfoo\nnode_modules\n", "node_modules", "frontend/node_modules\nfoo\n"},
		{"Not existing", "foo\nbar", ".oops/", "foo\nbar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeIgnoreTarget(tt.content, tt.target); got != tt.expected {
				t.Errorf("removeIgnoreTarget(%q, %q) = %q, want %q", tt.content, tt.target, got, tt.expected)
			}
		})
	}
}
