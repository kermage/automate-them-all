package cloak

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCloakCmd(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T) string
		targets      []string
		cloak        bool
		wantErr      bool
		errorMsg     string
		wantMsg      []string
		wantInConfig []string
		notInConfig  []string
	}{
		{
			name: "Failure - Not a git repository",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			targets:  []string{"test"},
			cloak:    true,
			wantErr:  true,
			errorMsg: "not a git repository",
		},
		{
			name: "Failure - No targets provided",
			setup: func(t *testing.T) string {
				return setupTestGitRepo(t)
			},
			targets:  nil,
			cloak:    true,
			wantErr:  true,
			errorMsg: "at least one target is required",
		},
		{
			name: "Success - Cloak single target",
			setup: func(t *testing.T) string {
				return setupTestGitRepo(t)
			},
			targets:      []string{"node_modules"},
			cloak:        true,
			wantMsg:      []string{"Cloaked node_modules"},
			wantInConfig: []string{"node_modules"},
		},
		{
			name: "Success - Cloak multiple targets",
			setup: func(t *testing.T) string {
				return setupTestGitRepo(t)
			},
			targets:      []string{"a", "b"},
			cloak:        true,
			wantMsg:      []string{"Cloaked a", "Cloaked b"},
			wantInConfig: []string{"a", "b"},
		},
		{
			name: "Success - Already cloaked",
			setup: func(t *testing.T) string {
				tmpDir := setupTestGitRepo(t)
				excludePath := filepath.Join(tmpDir, ".git", "info", "exclude")
				os.WriteFile(excludePath, []byte("node_modules\n"), 0o644)
				return tmpDir
			},
			targets:      []string{"node_modules"},
			cloak:        true,
			wantMsg:      []string{"node_modules is already cloaked"},
			wantInConfig: []string{"node_modules"},
		},
		{
			name: "Success - Uncloak target",
			setup: func(t *testing.T) string {
				tmpDir := setupTestGitRepo(t)
				excludePath := filepath.Join(tmpDir, ".git", "info", "exclude")
				os.WriteFile(excludePath, []byte("node_modules\n"), 0o644)
				return tmpDir
			},
			targets:     []string{"node_modules"},
			cloak:       false,
			wantMsg:     []string{"Uncloaked node_modules"},
			notInConfig: []string{"node_modules"},
		},
		{
			name: "Success - Uncloak preserves other entries",
			setup: func(t *testing.T) string {
				tmpDir := setupTestGitRepo(t)
				excludePath := filepath.Join(tmpDir, ".git", "info", "exclude")
				os.WriteFile(excludePath, []byte("node_modules\nother_entry\n"), 0o644)
				return tmpDir
			},
			targets:      []string{"node_modules"},
			cloak:        false,
			wantMsg:      []string{"Uncloaked node_modules"},
			wantInConfig: []string{"other_entry"},
			notInConfig:  []string{"node_modules"},
		},
		{
			name: "Success - Uncloak warning if not cloaked",
			setup: func(t *testing.T) string {
				return setupTestGitRepo(t)
			},
			targets: []string{"node_modules"},
			cloak:   false,
			wantMsg: []string{"node_modules is not cloaked"},
		},
		{
			name: "Success - Uncloak warning if exclude file does not exist",
			setup: func(t *testing.T) string {
				tmpDir := setupTestGitRepo(t)
				excludePath := filepath.Join(tmpDir, ".git", "info", "exclude")
				// Ensure it doesn't exist (some git versions might create it)
				_ = os.Remove(excludePath)
				return tmpDir
			},
			targets: []string{"node_modules"},
			cloak:   false,
			wantMsg: []string{"node_modules is not cloaked"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := tt.setup(t)
			var out strings.Builder

			err := Run(tmpDir, tt.targets, tt.cloak, &out)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errorMsg)
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %v", tt.errorMsg, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			output := out.String()
			for _, msg := range tt.wantMsg {
				if !strings.Contains(output, msg) {
					t.Errorf("expected output to contain %q, got %q", msg, output)
				}
			}

			excludePath := filepath.Join(tmpDir, ".git", "info", "exclude")
			if _, err := os.Stat(excludePath); err == nil {
				content, err := os.ReadFile(excludePath)
				if err != nil {
					t.Fatalf("failed to read exclude file: %v", err)
				}
				sContent := string(content)

				for _, target := range tt.wantInConfig {
					if !strings.Contains(sContent, target) {
						t.Errorf("expected exclude file to contain %q, got %q", target, sContent)
					}
				}
				for _, target := range tt.notInConfig {
					if strings.Contains(sContent, target) {
						t.Errorf("expected exclude file NOT to contain %q, got %q", target, sContent)
					}
				}
			}
		})
	}
}
