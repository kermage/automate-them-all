package srtf

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tempDir := t.TempDir()

	// Setup:
	// /tempDir/folder1/movie.srt
	// /tempDir/folder2/movie.srt
	// /tempDir/nested/folder3/movie.srt
	// /tempDir/deep/nested/folder4/movie.srt

	folders := []string{"folder1", "folder2", "nested/folder3", "deep/nested/folder4"}
	srtFilename := "movie.srt"

	for _, folder := range folders {
		path := filepath.Join(tempDir, folder)
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("failed to create directory %s: %v", path, err)
		}
		srtPath := filepath.Join(path, srtFilename)
		if err := os.WriteFile(srtPath, []byte("subtitle content for "+folder), 0644); err != nil {
			t.Fatalf("failed to write file %s: %v", srtPath, err)
		}
	}

	var out strings.Builder
	if err := Run(tempDir, srtFilename, &out); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Verify:
	// /tempDir/folder1.srt
	// /tempDir/folder2.srt
	// /tempDir/nested/folder3.srt
	// /tempDir/deep/nested/folder4.srt

	expectedFiles := map[string]string{
		"folder1.srt":             "subtitle content for folder1",
		"folder2.srt":             "subtitle content for folder2",
		"nested/folder3.srt":      "subtitle content for nested/folder3",
		"deep/nested/folder4.srt": "subtitle content for deep/nested/folder4",
	}

	for fileName, content := range expectedFiles {
		path := filepath.Join(tempDir, fileName)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("expected file %s to exist, but it doesn't: %v", path, err)
			continue
		}
		if string(data) != content {
			t.Errorf("expected content for %s to be %q, got %q", path, content, string(data))
		}
	}
}

// Verify that Run returns an error when a file cannot be written to the destination
func TestRunPropagatesCopyError(t *testing.T) {
	tempDir := t.TempDir()

	// Create source: /tempDir/movie/movie.srt
	srcDir := filepath.Join(tempDir, "movie")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "movie.srt"), []byte("sub"), 0644); err != nil {
		t.Fatal(err)
	}

	// Make tempDir read-only so os.Create for the destination fails.
	if err := os.Chmod(tempDir, 0555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(tempDir, 0755) })

	var out strings.Builder
	err := Run(tempDir, "movie.srt", &out)
	if err == nil {
		t.Fatal("expected Run to return an error when destination is not writable, got nil")
	}
}
