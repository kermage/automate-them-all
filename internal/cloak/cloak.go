package cloak

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Run(baseDir string, args []string, isCloak bool, out io.Writer) error {
	if out == nil {
		out = io.Discard
	}

	gitDir, err := getAbsoluteGitDir(baseDir)
	if err != nil {
		return fmt.Errorf("Unable to resolve git directory: %v", err)
	}

	infoDir := filepath.Join(gitDir, "info")
	if isCloak {
		if err := os.MkdirAll(infoDir, DirPerm); err != nil {
			return fmt.Errorf("Unable to create .git/info: %v", err)
		}
	}

	excludeFile := filepath.Join(infoDir, "exclude")
	content, err := os.ReadFile(excludeFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Unable to read exclude file: %v", err)
	}

	if len(args) == 0 {
		return fmt.Errorf("at least one target is required")
	}

	s := string(content)
	changed := false

	for _, target := range args {
		var targetChanged bool
		s, targetChanged = toggleIgnoreTarget(s, target, isCloak)
		if targetChanged {
			changed = true
		}
		printCloakResult(out, target, isCloak, targetChanged)
	}

	if changed {
		if err := os.WriteFile(excludeFile, []byte(s), FilePerm); err != nil {
			return fmt.Errorf("Unable to write to exclude file: %v", err)
		}
	}

	return nil
}

func toggleIgnoreTarget(content, target string, isCloak bool) (string, bool) {
	has := hasIgnoreTarget(content, target)

	if isCloak {
		if has {
			return content, false
		}
		return addIgnoreTarget(content, target), true
	}

	if !has {
		return content, false
	}
	return removeIgnoreTarget(content, target), true
}

func printCloakResult(out io.Writer, target string, isCloak, changed bool) {
	switch {
	case isCloak && changed:
		fmt.Fprintln(out, "Cloaked "+target)
	case isCloak && !changed:
		fmt.Fprintln(out, target+" is already cloaked")
	case !isCloak && changed:
		fmt.Fprintln(out, "Uncloaked "+target)
	default:
		fmt.Fprintln(out, target+" is not cloaked")
	}
}
