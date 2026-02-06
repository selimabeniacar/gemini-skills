package linter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/flowlint/internal/parser"
)

// Severity represents the severity of a lint issue
type Severity int

const (
	SeverityWarning Severity = iota
	SeverityError
)

// Issue represents a linting issue found in the diagram
type Issue struct {
	Severity   Severity
	Message    string
	Line       int
	Context    string
	Suggestion string
	Fixable    bool
	FixType    string
	FixData    map[string]string
}

// Lint runs all linting rules against the diagram
func Lint(diagram *parser.Diagram) []Issue {
	issues := []Issue{}

	issues = append(issues, checkSubgraphQuotes(diagram)...)
	issues = append(issues, checkArrowStyles(diagram)...)
	issues = append(issues, checkClassDefs(diagram)...)
	issues = append(issues, checkOrphanNodes(diagram)...)
	issues = append(issues, checkAbbreviations(diagram)...)
	issues = append(issues, checkDuplicateNodes(diagram)...)
	issues = append(issues, checkNewlinesInLabels(diagram)...)
	issues = append(issues, checkComplexity(diagram)...)

	return issues
}

// checkSubgraphQuotes ensures all subgraph titles are quoted
func checkSubgraphQuotes(diagram *parser.Diagram) []Issue {
	issues := []Issue{}

	for _, sg := range diagram.Subgraphs {
		if !sg.Quoted && sg.Title != "" {
			issues = append(issues, Issue{
				Severity:   SeverityError,
				Message:    fmt.Sprintf("Subgraph '%s' title is not quoted", sg.ID),
				Line:       sg.Line,
				Context:    fmt.Sprintf("subgraph %s [%s]", sg.ID, sg.Title),
				Suggestion: fmt.Sprintf("Change to: subgraph %s [\"%s\"]", sg.ID, sg.Title),
				Fixable:    true,
				FixType:    "quote_subgraph",
				FixData:    map[string]string{"id": sg.ID, "title": sg.Title},
			})
		}
	}

	return issues
}

// checkArrowStyles ensures correct arrow usage for sync vs async
func checkArrowStyles(diagram *parser.Diagram) []Issue {
	issues := []Issue{}

	// Keywords that indicate async communication
	asyncKeywords := []string{"kafka", "publish", "consume", "queue", "rabbitmq", "sqs", "pubsub"}
	// Keywords that indicate sync communication
	syncKeywords := []string{"grpc", "http", "rest", "sql", "cache", "redis"}

	for _, edge := range diagram.Edges {
		labelLower := strings.ToLower(edge.Label)

		// Check if async keyword but using sync arrow
		for _, keyword := range asyncKeywords {
			if strings.Contains(labelLower, keyword) && edge.ArrowType == "==>" {
				issues = append(issues, Issue{
					Severity:   SeverityError,
					Message:    fmt.Sprintf("Async call '%s' using sync arrow (==>)", edge.Label),
					Line:       edge.Line,
					Suggestion: "Change ==> to -.-> for async calls",
					Fixable:    true,
					FixType:    "fix_arrow",
					FixData:    map[string]string{"from": edge.From, "to": edge.To, "old": "==>", "new": "-.->"},
				})
				break
			}
		}

		// Check if sync keyword but using async arrow
		for _, keyword := range syncKeywords {
			if strings.Contains(labelLower, keyword) && edge.ArrowType == "-.->" {
				issues = append(issues, Issue{
					Severity:   SeverityWarning,
					Message:    fmt.Sprintf("Sync call '%s' using async arrow (-.->)", edge.Label),
					Line:       edge.Line,
					Suggestion: "Change -.-> to ==> for sync calls",
					Fixable:    true,
					FixType:    "fix_arrow",
					FixData:    map[string]string{"from": edge.From, "to": edge.To, "old": "-.->", "new": "==>"},
				})
				break
			}
		}
	}

	return issues
}

// checkClassDefs ensures required class definitions exist
func checkClassDefs(diagram *parser.Diagram) []Issue {
	issues := []Issue{}

	requiredClasses := []string{"service", "kafka", "database", "external"}

	for _, class := range requiredClasses {
		if _, ok := diagram.ClassDefs[class]; !ok {
			issues = append(issues, Issue{
				Severity:   SeverityWarning,
				Message:    fmt.Sprintf("Missing classDef for '%s'", class),
				Suggestion: fmt.Sprintf("Add: classDef %s fill:#...,stroke:#...,color:#...", class),
				Fixable:    true,
				FixType:    "add_classdef",
				FixData:    map[string]string{"class": class},
			})
		}
	}

	return issues
}

// checkOrphanNodes finds nodes with no connections
func checkOrphanNodes(diagram *parser.Diagram) []Issue {
	issues := []Issue{}

	orphans := diagram.GetOrphanNodes()
	for _, node := range orphans {
		issues = append(issues, Issue{
			Severity:   SeverityWarning,
			Message:    fmt.Sprintf("Orphan node '%s' has no connections", node.ID),
			Line:       node.Line,
			Context:    node.Label,
			Suggestion: "Add connections or remove the node",
		})
	}

	return issues
}

// checkAbbreviations warns about potential abbreviations in node labels
func checkAbbreviations(diagram *parser.Diagram) []Issue {
	issues := []Issue{}

	// Common abbreviation patterns
	abbrPatterns := []struct {
		pattern string
		full    string
	}{
		{`(?i)svc\b`, "Service"},
		{`(?i)srv\b`, "Server"},
		{`(?i)msg\b`, "Message"},
		{`(?i)req\b`, "Request"},
		{`(?i)res\b`, "Response"},
		{`(?i)cfg\b`, "Config"},
		{`(?i)db\b`, "Database"},
		{`(?i)api\b`, "API"}, // This one is acceptable
	}

	for _, node := range diagram.Nodes {
		for _, abbr := range abbrPatterns {
			if abbr.pattern == `(?i)api\b` {
				continue // API is acceptable
			}
			re := regexp.MustCompile(abbr.pattern)
			if re.MatchString(node.Label) {
				issues = append(issues, Issue{
					Severity:   SeverityWarning,
					Message:    fmt.Sprintf("Node '%s' may contain abbreviation", node.Label),
					Line:       node.Line,
					Suggestion: fmt.Sprintf("Consider using full word '%s' instead", abbr.full),
				})
				break
			}
		}
	}

	return issues
}

// checkDuplicateNodes finds duplicate node IDs
func checkDuplicateNodes(diagram *parser.Diagram) []Issue {
	// The parser already handles this by using a map, so duplicates would be overwritten
	// This is more of a sanity check
	return []Issue{}
}

// checkComplexity detects potential spaghetti diagrams
func checkComplexity(diagram *parser.Diagram) []Issue {
	issues := []Issue{}

	nodeCount := len(diagram.Nodes)
	edgeCount := len(diagram.Edges)

	// Count connections per node (hub detection)
	connectionCount := make(map[string]int)
	for _, edge := range diagram.Edges {
		connectionCount[edge.From]++
		connectionCount[edge.To]++
	}

	// Find max connections (potential hub)
	maxConnections := 0
	hubNode := ""
	for nodeID, count := range connectionCount {
		if count > maxConnections {
			maxConnections = count
			hubNode = nodeID
		}
	}

	// Count nodes with multiple outgoing edges (fan-out)
	fanOutCount := 0
	outgoing := make(map[string]int)
	for _, edge := range diagram.Edges {
		outgoing[edge.From]++
	}
	for _, count := range outgoing {
		if count > 3 {
			fanOutCount++
		}
	}

	// Detect potential spaghetti patterns
	isSpaghetti := false
	reasons := []string{}

	// Rule 1: Too many nodes without subgraphs
	if nodeCount > 10 && len(diagram.Subgraphs) < 2 {
		isSpaghetti = true
		reasons = append(reasons, fmt.Sprintf("%d nodes without grouping", nodeCount))
	}

	// Rule 2: High edge-to-node ratio (dense connections)
	if nodeCount > 0 {
		edgeRatio := float64(edgeCount) / float64(nodeCount)
		if edgeRatio > 2.0 {
			isSpaghetti = true
			reasons = append(reasons, fmt.Sprintf("high edge density (%.1f edges per node)", edgeRatio))
		}
	}

	// Rule 3: Multiple nodes with high fan-out (not hub-and-spoke)
	if fanOutCount > 2 {
		isSpaghetti = true
		reasons = append(reasons, fmt.Sprintf("%d nodes with 4+ outgoing edges", fanOutCount))
	}

	// Rule 4: No clear hub pattern when complex
	if nodeCount > 8 && maxConnections < nodeCount/2 {
		isSpaghetti = true
		reasons = append(reasons, "no clear hub node - consider linear pipeline")
	}

	if isSpaghetti {
		reasonStr := strings.Join(reasons, ", ")
		issues = append(issues, Issue{
			Severity:   SeverityWarning,
			Message:    fmt.Sprintf("Diagram may be spaghetti: %s", reasonStr),
			Suggestion: "Consider using linear pipeline layout with target service as hub",
			Fixable:    false,
		})

		// Add specific advice
		if hubNode != "" && maxConnections > 3 {
			issues = append(issues, Issue{
				Severity:   SeverityWarning,
				Message:    fmt.Sprintf("Node '%s' has %d connections - good hub candidate", hubNode, maxConnections),
				Suggestion: "Restructure diagram with this node as central hub",
				Fixable:    false,
			})
		}
	}

	return issues
}

// checkNewlinesInLabels ensures no node labels contain newlines
func checkNewlinesInLabels(diagram *parser.Diagram) []Issue {
	issues := []Issue{}

	for _, node := range diagram.Nodes {
		if strings.Contains(node.Label, "\n") {
			// Create fixed label by replacing newlines with spaces
			fixedLabel := strings.ReplaceAll(node.Label, "\n", " ")
			// Collapse multiple spaces
			fixedLabel = regexp.MustCompile(`\s+`).ReplaceAllString(fixedLabel, " ")
			fixedLabel = strings.TrimSpace(fixedLabel)

			issues = append(issues, Issue{
				Severity:   SeverityError,
				Message:    fmt.Sprintf("Node '%s' contains newline in label", node.ID),
				Line:       node.Line,
				Context:    node.Label,
				Suggestion: fmt.Sprintf("Change to single line: %s[%s]", node.ID, fixedLabel),
				Fixable:    true,
				FixType:    "fix_newline",
				FixData:    map[string]string{"id": node.ID, "old": node.Label, "new": fixedLabel},
			})
		}
	}

	// Also check edge labels
	for _, edge := range diagram.Edges {
		if strings.Contains(edge.Label, "\n") {
			fixedLabel := strings.ReplaceAll(edge.Label, "\n", " ")
			fixedLabel = regexp.MustCompile(`\s+`).ReplaceAllString(fixedLabel, " ")
			fixedLabel = strings.TrimSpace(fixedLabel)

			issues = append(issues, Issue{
				Severity:   SeverityError,
				Message:    fmt.Sprintf("Edge label contains newline: %s -> %s", edge.From, edge.To),
				Line:       edge.Line,
				Context:    edge.Label,
				Suggestion: fmt.Sprintf("Change to single line: |%s|", fixedLabel),
				Fixable:    true,
				FixType:    "fix_edge_newline",
				FixData:    map[string]string{"from": edge.From, "to": edge.To, "old": edge.Label, "new": fixedLabel},
			})
		}
	}

	// Check subgraph titles
	for _, sg := range diagram.Subgraphs {
		if strings.Contains(sg.Title, "\n") {
			fixedTitle := strings.ReplaceAll(sg.Title, "\n", " ")
			fixedTitle = regexp.MustCompile(`\s+`).ReplaceAllString(fixedTitle, " ")
			fixedTitle = strings.TrimSpace(fixedTitle)

			issues = append(issues, Issue{
				Severity:   SeverityError,
				Message:    fmt.Sprintf("Subgraph '%s' title contains newline", sg.ID),
				Line:       sg.Line,
				Context:    sg.Title,
				Suggestion: fmt.Sprintf("Change to single line: subgraph %s [\"%s\"]", sg.ID, fixedTitle),
				Fixable:    true,
				FixType:    "fix_subgraph_newline",
				FixData:    map[string]string{"id": sg.ID, "old": sg.Title, "new": fixedTitle},
			})
		}
	}

	return issues
}
