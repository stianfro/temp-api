package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: commit-msg <path to commit message file>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening commit message file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	commitMessage := scanner.Text()

	matched, err := regexp.MatchString(`^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert|bump)(\(\w+\))?:\s.+$`, commitMessage)
	if err != nil {
		fmt.Printf("Error in regex matching: %s\n", err)
		os.Exit(1)
	}

	if !matched {
		fmt.Println("ERROR: Commit message does not follow Conventional Commits format.")
		fmt.Println("Example of a valid message: 'feat(database): add new indexing feature'")
		os.Exit(1)
	}
}
