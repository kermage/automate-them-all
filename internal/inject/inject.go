package inject

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func Run(targetJSON, sourceFile, keyPath string, out io.Writer) error {
	if out == nil {
		out = io.Discard
	}

	sourceContent, keyPath, err := readSource(sourceFile, keyPath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	jsonContentBytes, isNew, err := readJSONFile(targetJSON)
	if err != nil {
		return fmt.Errorf("failed to handle target JSON file: %w", err)
	}

	jsonContent := string(jsonContentBytes)
	existingValue := gjson.Get(jsonContent, keyPath)
	if existingValue.String() == sourceContent {
		fmt.Fprintln(out, "[skipped] content identical")
		return nil
	}

	var action string
	switch {
	case isNew:
		action = "created"
	case existingValue.Raw == "":
		action = "added"
	default:
		action = "updated"
	}

	newJSONContent, err := sjson.Set(jsonContent, keyPath, sourceContent)
	if err != nil {
		return fmt.Errorf("failed to set value in JSON: %w", err)
	}

	if err := os.WriteFile(targetJSON, []byte(newJSONContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated JSON: %w", err)
	}

	fmt.Fprintf(out, "[%s] content value\n", action)

	return nil
}

func readSource(path, key string) (string, string, error) {
	srcContentBytes, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
	}

	sourceContent := string(srcContentBytes)
	if key == "" {
		base := filepath.Base(path)
		key = strings.TrimSuffix(base, filepath.Ext(base))
	}

	return sourceContent, key, nil
}

func readJSONFile(path string) ([]byte, bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []byte("{}"), true, nil
		}

		return nil, false, fmt.Errorf("stat failed: %w", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("read failed: %w", err)
	}

	return content, false, nil
}
