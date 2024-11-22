package main

// Credit: https://github.com/5rahim/seanime (based on)

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	inFile  = "CHANGELOG.md"
	outFile = "whats-new.md"
)

var versionRegex = regexp.MustCompile(`(?m)^## v\d+(\.\d+)*([a-zA-Z0-9-]+)?`)

func main() {
	inFilePath := filepath.Join(".", inFile)
	changelogContent, err := os.ReadFile(inFilePath)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	versionSections := versionRegex.Split(string(changelogContent), -1)

	var validSections []string
	for _, section := range versionSections {
		if !strings.Contains(section, "<!--") && !strings.Contains(section, "-->") {
			validSections = append(validSections, section)
		}
	}

	if len(validSections) == 0 {
		fmt.Println("No valid version sections found in changelog")
		return
	}

	lines := strings.Split(validSections[0], "\n")
	var filteredLines []string
	for _, line := range lines {
		if strings.HasPrefix(line, "## v") {
			continue
		}

		if line != "" || len(filteredLines) > 0 {
			filteredLines = append(filteredLines, line)
		}
	}

	currentChangelog := strings.TrimSpace(strings.Join(filteredLines, "\n"))
	outFilePath := filepath.Join(".", outFile)

	if err = os.WriteFile(outFilePath, []byte(currentChangelog), 0644); err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		return
	}

	fmt.Printf("Changelog successfully written to %s\n", outFilePath)
}
