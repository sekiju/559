package main

import (
	"flag"
	"fmt"
	"log"
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
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	version := flag.String("version", "latest", "The version to search for in the changelog (e.g., v1.2.3 or latest)")
	flag.Parse()

	inFilePath := filepath.Join(".", inFile)
	content, err := os.ReadFile(inFilePath)
	if err != nil {
		return fmt.Errorf("error reading input file: %v", err)
	}

	section, err := findVersionSection(string(content), *version)
	if err != nil {
		log.Fatalf("Error finding version section: %v", err)
	}

	lines := strings.Split(section, "\n")
	var filteredLines []string
	for _, line := range lines {
		if strings.HasPrefix(line, "## v") || strings.HasPrefix(line, "Date:") {
			continue
		}

		if line != "" || len(filteredLines) > 0 {
			filteredLines = append(filteredLines, line)
		}
	}

	currentChangelog := strings.TrimSpace(strings.Join(filteredLines, "\n"))
	outFilePath := filepath.Join(".", outFile)

	if err = os.WriteFile(outFilePath, []byte(currentChangelog), 0644); err != nil {
		return fmt.Errorf("error writing to output file: %v", err)
	}

	fmt.Printf("Changelog successfully written to %s\n", outFilePath)

	return nil
}

func findVersionSection(content, version string) (string, error) {
	versionMatches := versionRegex.FindAllStringIndex(content, -1)
	if len(versionMatches) < 2 {
		return "", fmt.Errorf("insufficient version headers found")
	}

	versionMatches = versionMatches[1:]

	if version == "latest" {
		if len(versionMatches) > 1 {
			return content[versionMatches[0][0]:versionMatches[1][0]], nil
		} else {
			return content[versionMatches[0][0]:], nil
		}
	}

	for i := 0; i < len(versionMatches); i++ {
		versionHeader := content[versionMatches[i][0]:versionMatches[i][1]]

		if strings.HasPrefix(versionHeader, "## v"+version) {
			nextStart := len(content)

			if i+1 < len(versionMatches) {
				nextStart = versionMatches[i+1][0]
			}

			return content[versionMatches[i][0]:nextStart], nil
		}
	}

	return "", fmt.Errorf("version %s not found", version)
}
