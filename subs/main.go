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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the full path to process: ")
	fullPath, _ := reader.ReadString('\n')

	if (fullPath == "\n") {
		log.Fatal("Error unspecified path to pricess")
	}

	fmt.Print("Enter the wanted SRT file: ")
	srtFilename, _ := reader.ReadString('\n')

	if (srtFilename == "\n") {
		log.Fatal("Error unspecified SRT file")
	}

	srtFilename = strings.TrimSpace(srtFilename)
	srtFilename = strings.TrimSuffix(srtFilename, ".srt") + ".srt"

	handleContents(strings.TrimSpace(fullPath), srtFilename)
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
			if (item.Name() != srtFilename) {
				log.Printf("Skipping: %s\n", currentPath)
				continue
			}

			var newPath = filepath.Join(fullPath, "..", filepath.Base(fullPath) + ".srt")

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
