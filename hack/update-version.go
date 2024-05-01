package main

import (
	"log"
	"os"
	"regexp"
)

const (
	filePermission = 0o644
	filePath       = "main.go"
	versionConst   = "const version ="
)

func main() {
	newVersion := os.Getenv("CZ_PRE_NEW_VERSION")
	if newVersion == "" {
		log.Fatalf("CZ_PRE_NEW_VERSION is not set")
	}
	newVersion = `"` + newVersion + `"`

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	r := regexp.MustCompile(`(?m)^const version =.*$`)
	newContent := r.ReplaceAllString(string(content), versionConst+" "+newVersion)

	if err = os.WriteFile(filePath, []byte(newContent), os.FileMode(filePermission)); err != nil {
		log.Fatalf("failed to write to file: %s", err)
	}
}
