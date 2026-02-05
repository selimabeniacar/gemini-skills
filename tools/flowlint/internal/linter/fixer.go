package linter

import (
	"fmt"
	"regexp"
	"strings"
)

// Default class definitions to inject if missing (muted professional colors)
var defaultClassDefs = map[string]string{
	"service":  "classDef service fill:#a5d8ff,stroke:#339af0,color:#1864ab",
	"entry":    "classDef entry fill:#b2f2bb,stroke:#51cf66,color:#2b8a3e",
	"kafka":    "classDef kafka fill:#96f2d7,stroke:#38d9a9,color:#087f5b",
	"database": "classDef database fill:#ffec99,stroke:#fcc419,color:#e67700",
	"cache":    "classDef cache fill:#d0bfff,stroke:#9775fa,color:#6741d9",
	"external": "classDef external fill:#dee2e6,stroke:#adb5bd,color:#495057",
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

		case "fix_newline":
			// Fix newlines in node labels
			nodeID := issue.FixData["id"]
			newLabel := issue.FixData["new"]
			code = fixNodeLabel(code, nodeID, newLabel)
			lines = strings.Split(code, "\n")
			fixCount++

		case "fix_edge_newline":
			// Fix newlines in edge labels
			oldLabel := issue.FixData["old"]
			newLabel := issue.FixData["new"]
			// Replace the multi-line label with single-line
			escapedOld := regexp.QuoteMeta(oldLabel)
			// Handle the label possibly spanning multiple lines
			re := regexp.MustCompile(`\|` + escapedOld + `\|`)
			code = re.ReplaceAllString(code, "|"+newLabel+"|")
			lines = strings.Split(code, "\n")
			fixCount++

		case "fix_subgraph_newline":
			// Fix newlines in subgraph titles
			sgID := issue.FixData["id"]
			newTitle := issue.FixData["new"]
			// Find and fix the subgraph line
			for i, line := range lines {
				if strings.Contains(line, "subgraph "+sgID) {
					// Replace the title
					re := regexp.MustCompile(`subgraph\s+` + regexp.QuoteMeta(sgID) + `\s*\[.*\]`)
					lines[i] = re.ReplaceAllString(line, fmt.Sprintf(`subgraph %s ["%s"]`, sgID, newTitle))
					fixCount++
					break
				}
			}
		}
	}

	return strings.Join(lines, "\n"), fixCount
}

// fixNodeLabel replaces a node's label, handling multi-line cases
func fixNodeLabel(code string, nodeID string, newLabel string) string {
	// Match node definitions like: A[Label] or A[Multi\nLine]
	// This regex handles the node ID followed by various bracket types
	patterns := []string{
		`(\s` + regexp.QuoteMeta(nodeID) + `)\[([^\]]*)\]`,         // [Label]
		`(\s` + regexp.QuoteMeta(nodeID) + `)\[\(([^\)]*)\)\]`,     // [(Label)]
		`(\s` + regexp.QuoteMeta(nodeID) + `)\(\(([^\)]*)\)\)`,     // ((Label))
		`(\s` + regexp.QuoteMeta(nodeID) + `)\(([^\)]*)\)`,         // (Label)
		`(\s` + regexp.QuoteMeta(nodeID) + `)\[\[([^\]]*)\]\]`,     // [[Label]]
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(code) {
			// Determine the bracket type from the match
			match := re.FindStringSubmatch(code)
			if len(match) > 0 {
				// Replace preserving the bracket type
				if strings.Contains(pattern, `\[\(`) {
					code = re.ReplaceAllString(code, fmt.Sprintf(`$1[(%s)]`, newLabel))
				} else if strings.Contains(pattern, `\(\(`) {
					code = re.ReplaceAllString(code, fmt.Sprintf(`$1((%s))`, newLabel))
				} else if strings.Contains(pattern, `\[\[`) {
					code = re.ReplaceAllString(code, fmt.Sprintf(`$1[[%s]]`, newLabel))
				} else if strings.Contains(pattern, `\(`) {
					code = re.ReplaceAllString(code, fmt.Sprintf(`$1(%s)`, newLabel))
				} else {
					code = re.ReplaceAllString(code, fmt.Sprintf(`$1[%s]`, newLabel))
				}
				break
			}
		}
	}

	return code
}
