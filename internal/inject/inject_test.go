package inject

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) (targetJSON, sourceFile, keyPath string)
		wantErr  bool
		wantMsg  string
		expected string
	}{
		{
			name: "missing JSON file (creates file, prints created)",
			setup: func(t *testing.T) (string, string, string) {
				tempDir := t.TempDir()
				sourceFile := filepath.Join(tempDir, "source.md")
				err := os.WriteFile(sourceFile, []byte("new content"), 0644)
				if err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}
				return filepath.Join(tempDir, "missing.json"), sourceFile, "mykey"
			},
			wantErr:  false,
			wantMsg:  "[created] content value",
			expected: `{"mykey":"new content"}`,
		},
		{
			name: "target flag omitted (uses stripped source filename)",
			setup: func(t *testing.T) (string, string, string) {
				tempDir := t.TempDir()
				sourceFile := filepath.Join(tempDir, "readme.md")
				err := os.WriteFile(sourceFile, []byte("content"), 0644)
				if err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}

				targetJSON := filepath.Join(tempDir, "target.json")
				err = os.WriteFile(targetJSON, []byte(`{}`), 0644)
				if err != nil {
					t.Fatalf("failed to create JSON file: %v", err)
				}

				return targetJSON, sourceFile, ""
			},
			wantErr:  false,
			wantMsg:  "[added] content value",
			expected: `{"readme":"content"}`,
		},
		{
			name: "exact match (prints skipped)",
			setup: func(t *testing.T) (string, string, string) {
				tempDir := t.TempDir()
				sourceFile := filepath.Join(tempDir, "source.md")
				err := os.WriteFile(sourceFile, []byte("content"), 0644)
				if err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}

				targetJSON := filepath.Join(tempDir, "target.json")
				err = os.WriteFile(targetJSON, []byte(`{"mykey":"content"}`), 0644)
				if err != nil {
					t.Fatalf("failed to create JSON file: %v", err)
				}

				return targetJSON, sourceFile, "mykey"
			},
			wantErr:  false,
			wantMsg:  "[skipped] content identical",
			expected: `{"mykey":"content"}`,
		},
		{
			name: "existing key differs (prints updated)",
			setup: func(t *testing.T) (string, string, string) {
				tempDir := t.TempDir()
				sourceFile := filepath.Join(tempDir, "source.md")
				err := os.WriteFile(sourceFile, []byte("new content"), 0644)
				if err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}

				targetJSON := filepath.Join(tempDir, "target.json")
				err = os.WriteFile(targetJSON, []byte(`{"mykey":"old content"}`), 0644)
				if err != nil {
					t.Fatalf("failed to create JSON file: %v", err)
				}

				return targetJSON, sourceFile, "mykey"
			},
			wantErr:  false,
			wantMsg:  "[updated] content value",
			expected: `{"mykey":"new content"}`,
		},
		{
			name: "new key in existing JSON (prints added)",
			setup: func(t *testing.T) (string, string, string) {
				tempDir := t.TempDir()
				sourceFile := filepath.Join(tempDir, "source.md")
				err := os.WriteFile(sourceFile, []byte("content"), 0644)
				if err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}

				targetJSON := filepath.Join(tempDir, "target.json")
				err = os.WriteFile(targetJSON, []byte(`{"otherkey":"value"}`), 0644)
				if err != nil {
					t.Fatalf("failed to create JSON file: %v", err)
				}

				return targetJSON, sourceFile, "mykey"
			},
			wantErr:  false,
			wantMsg:  "[added] content value",
			expected: `{"otherkey":"value","mykey":"content"}`,
		},
		{
			name: "missing source file (returns error)",
			setup: func(t *testing.T) (string, string, string) {
				tempDir := t.TempDir()

				return filepath.Join(tempDir, "target.json"), filepath.Join(tempDir, "missing.md"), "mykey"
			},
			wantErr:  true,
			wantMsg:  "",
			expected: "",
		},
		{
			name: "target json read error (is a directory)",
			setup: func(t *testing.T) (string, string, string) {
				tempDir := t.TempDir()
				sourceFile := filepath.Join(tempDir, "source.md")
				err := os.WriteFile(sourceFile, []byte("content"), 0644)
				if err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}

				targetJSON := filepath.Join(tempDir, "target.json")
				err = os.Mkdir(targetJSON, 0755)
				if err != nil {
					t.Fatalf("failed to create JSON dir: %v", err)
				}

				return targetJSON, sourceFile, "mykey"
			},
			wantErr:  true,
			wantMsg:  "",
			expected: "",
		},
		{
			name: "target json stat error (permission denied)",
			setup: func(t *testing.T) (string, string, string) {
				tempDir := t.TempDir()
				sourceFile := filepath.Join(tempDir, "source.md")
				err := os.WriteFile(sourceFile, []byte("content"), 0644)
				if err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}

				unreadableDir := filepath.Join(tempDir, "noperms")
				err = os.Mkdir(unreadableDir, 0000)
				if err != nil {
					t.Fatalf("failed to create noperms dir: %v", err)
				}

				return filepath.Join(unreadableDir, "target.json"), sourceFile, "mykey"
			},
			wantErr:  true,
			wantMsg:  "",
			expected: "",
		},
		{
			name: "unwritable target JSON",
			setup: func(t *testing.T) (string, string, string) {
				tempDir := t.TempDir()
				sourceFile := filepath.Join(tempDir, "source.md")
				err := os.WriteFile(sourceFile, []byte("content"), 0644)
				if err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}

				targetJSON := filepath.Join(tempDir, "target.json")
				err = os.WriteFile(targetJSON, []byte(`{}`), 0400) // Read-only
				if err != nil {
					t.Fatalf("failed to create JSON file: %v", err)
				}

				return targetJSON, sourceFile, "mykey"
			},
			wantErr:  true,
			wantMsg:  "",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			targetJSON, sourceFile, keyPath := tc.setup(t)
			var out bytes.Buffer

			err := Run(targetJSON, sourceFile, keyPath, &out)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			output := out.String()
			if !strings.Contains(output, tc.wantMsg) {
				t.Errorf("expected output to contain: %q, got: %q", tc.wantMsg, output)
			}

			if tc.expected != "" {
				b, err := os.ReadFile(targetJSON)
				if err != nil {
					t.Fatalf("failed to read target json: %v", err)
				}

				// Trim whitespaces to normalize formatting differences
				actualJSON := strings.TrimSpace(string(b))
				expectedJSON := strings.TrimSpace(tc.expected)
				if actualJSON != expectedJSON {
					t.Errorf("expected json: %q, got: %q", expectedJSON, actualJSON)
				}
			}
		})
	}
}
