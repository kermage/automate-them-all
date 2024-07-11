package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	HOOKS_DIR = ".hooks"
	HOOKS_PERM = 0755
)

func main() {
	var path string
	args := os.Args[1:]

	if len(args) == 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the full path to process: ")
		path, _ = reader.ReadString('\n')
		path = strings.TrimSpace(path)
	} else {
		path = args[0]
	}

	path, err := gitDir(path)

	if err != nil {
		log.Fatalln(err)
	}

	err = initDir(path)

	if err != nil {
		log.Println(err)
	}
}

func gitDir(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--absolute-git-dir")
	cmd.Dir = path
	res, err := cmd.Output()

	return string(res), err
}

func initDir(path string) error {
	basePath := filepath.Join(path, "..")
	cmd := exec.Command("git", "config", "core.hooksPath", HOOKS_DIR)
	cmd.Dir = basePath
	_, err := cmd.Output()

	if (err != nil) {
		return err
	}

	hooksPath := filepath.Join(basePath, HOOKS_DIR)
	_, err = os.Stat(hooksPath)

	if os.IsNotExist(err) {
		return os.Mkdir(hooksPath, HOOKS_PERM)
	}

	return err
}
