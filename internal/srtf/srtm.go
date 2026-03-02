package srtf

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Run(rootPath string, srtFilename string, out io.Writer) error {
	if !strings.HasSuffix(srtFilename, ".srt") {
		srtFilename += ".srt"
	}

	return filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if d.Name() == srtFilename {
			parentDir := filepath.Dir(path)
			grandParentDir := filepath.Dir(parentDir)
			newPath := filepath.Join(grandParentDir, filepath.Base(parentDir)+".srt")

			fmt.Fprintf(out, "Copying: %s\n", path)
			fmt.Fprintf(out, "To: %s\n", newPath)

			if err := copyFile(path, newPath); err != nil {
				return fmt.Errorf("failed to copy %s to %s: %w", path, newPath, err)
			}
		}

		return nil
	})
}

func copyFile(source string, destination string) (err error) {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}
