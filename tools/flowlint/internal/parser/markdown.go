package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// ExtractMermaid extracts the mermaid code block from markdown content
func ExtractMermaid(content string) (string, error) {
	// Match ```mermaid ... ```
	re := regexp.MustCompile("(?s)```mermaid\\s*\\n(.+?)\\n```")
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return "", fmt.Errorf("no mermaid code block found")
	}
	return strings.TrimSpace(matches[1]), nil
}

// ReplaceMermaid replaces the mermaid code block in markdown with new code
func ReplaceMermaid(content, newMermaid string) string {
	re := regexp.MustCompile("(?s)```mermaid\\s*\\n.+?\\n```")
	return re.ReplaceAllString(content, "```mermaid\n"+newMermaid+"\n```")
}
