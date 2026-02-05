package linter

import (
	"fmt"
	"regexp"
	"strings"
)

// Default class definitions to inject if missing
var defaultClassDefs = map[string]string{
	"service":  "classDef service fill:#228be6,stroke:#1971c2,color:#fff",
	"entry":    "classDef entry fill:#40c057,stroke:#2f9e44,color:#fff",
	"kafka":    "classDef kafka fill:#12b886,stroke:#099268,color:#fff",
	"database": "classDef database fill:#fab005,stroke:#f59f00,color:#000",
	"cache":    "classDef cache fill:#be4bdb,stroke:#9c36b5,color:#fff",
	"external": "classDef external fill:#868e96,stroke:#495057,color:#fff",
}

// Fix applies automatic fixes to the mermaid code based on issues
func Fix(code string, issues []Issue) (string, int) {
	fixCount := 0
	lines := strings.Split(code, "\n")

	for _, issue := range issues {
		if !issue.Fixable {
			continue
		}

		switch issue.FixType {
		case "quote_subgraph":
			id := issue.FixData["id"]
			title := issue.FixData["title"]
			// Find and fix the subgraph line
			for i, line := range lines {
				if strings.Contains(line, "subgraph "+id) && !strings.Contains(line, `"`) {
					// Replace unquoted with quoted
					oldPattern := fmt.Sprintf(`subgraph\s+%s\s*\[([^\]"]+)\]`, regexp.QuoteMeta(id))
					re := regexp.MustCompile(oldPattern)
					if re.MatchString(line) {
						lines[i] = re.ReplaceAllString(line, fmt.Sprintf(`subgraph %s ["%s"]`, id, title))
						fixCount++
					} else {
						// Try without brackets
						oldPattern2 := fmt.Sprintf(`subgraph\s+%s\s+([^"\[\]]+)$`, regexp.QuoteMeta(id))
						re2 := regexp.MustCompile(oldPattern2)
						if re2.MatchString(line) {
							lines[i] = re2.ReplaceAllString(line, fmt.Sprintf(`subgraph %s ["%s"]`, id, title))
							fixCount++
						}
					}
					break
				}
			}

		case "fix_arrow":
			from := issue.FixData["from"]
			to := issue.FixData["to"]
			oldArrow := issue.FixData["old"]
			newArrow := issue.FixData["new"]

			// Find and fix the edge line
			for i, line := range lines {
				if strings.Contains(line, from) && strings.Contains(line, to) && strings.Contains(line, oldArrow) {
					lines[i] = strings.Replace(line, oldArrow, newArrow, 1)
					fixCount++
					break
				}
			}

		case "add_classdef":
			className := issue.FixData["class"]
			if def, ok := defaultClassDefs[className]; ok {
				// Find where to insert (after flowchart declaration or other classDefs)
				insertIdx := -1
				for i, line := range lines {
					if strings.HasPrefix(strings.TrimSpace(line), "flowchart") {
						insertIdx = i + 1
					}
					if strings.HasPrefix(strings.TrimSpace(line), "classDef") {
						insertIdx = i + 1
					}
				}
				if insertIdx > 0 && insertIdx < len(lines) {
					// Insert the classDef
					newLines := make([]string, 0, len(lines)+1)
					newLines = append(newLines, lines[:insertIdx]...)
					newLines = append(newLines, "    "+def)
					newLines = append(newLines, lines[insertIdx:]...)
					lines = newLines
					fixCount++
				}
			}
		}
	}

	return strings.Join(lines, "\n"), fixCount
}
