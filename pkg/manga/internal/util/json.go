package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ExtractFromHTML[T any](html, prefix, suffix string) (*T, error) {
	html = strings.ReplaceAll(html, "&quot;", "\"")

	startIdx := strings.Index(html, prefix)
	if startIdx == -1 {
		return nil, fmt.Errorf("prefix not found: %s", prefix)
	}

	startIdx += len(prefix)

	endIdx := strings.Index(html[startIdx:], suffix)
	if endIdx == -1 {
		return nil, fmt.Errorf("suffix not found: %s", suffix)
	}

	var result T
	if err := json.Unmarshal([]byte(html[startIdx:startIdx+endIdx]), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
