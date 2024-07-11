package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var fullPath string
	var srtFilename string

	getInput(&fullPath, "Enter the full path to process: ", 0)
	if fullPath == "" {
		log.Fatal("Error unspecified path to process")
	}

	getInput(&srtFilename, "Enter the wanted SRT file: ", 1)
	if srtFilename == "" {
		log.Fatal("Error unspecified SRT file")
	}

	handleContents(fullPath, strings.TrimSuffix(srtFilename, ".srt")+".srt")
}

func getInput(variable *string, prompt string, index int) {
	args := os.Args[1:]

	if index < len(args) {
		*variable = args[index]
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	temp, _ := reader.ReadString('\n')
	*variable = strings.TrimSpace(temp)
}

func getContents(directory string) []fs.DirEntry {
	contents, err := os.ReadDir(directory)

	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	return contents
}

func handleContents(fullPath string, srtFilename string) {
	for _, item := range getContents(fullPath) {
		var currentPath = filepath.Join(fullPath, item.Name())

		if item.IsDir() {
			handleContents(currentPath, srtFilename)
		} else {
			if item.Name() != srtFilename {
				continue
			}

			var newPath = filepath.Join(fullPath, "..", filepath.Base(fullPath)+".srt")

			log.Printf("Copying: %s\n", currentPath)
			log.Printf("To: %s\n", newPath)

			err := copyFile(currentPath, newPath)

			if err != nil {
				log.Printf("Error: %s\n", err)
			}
		}
	}
}

func copyFile(source string, destination string) (err error) {
	in, err := os.Open(source)

	if err != nil {
		return
	}

	defer in.Close()
	out, err := os.Create(destination)

	if err != nil {
		return
	}

	defer func() {
		cerr := out.Close()

		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return
	}

	out.Sync()

	return
}
