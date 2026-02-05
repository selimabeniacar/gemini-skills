package parser

import (
	"regexp"
	"strings"
)

// Node represents a node in the mermaid diagram
type Node struct {
	ID       string
	Label    string
	Shape    string // rectangle, cylinder, stadium, etc.
	Line     int
	Classes  []string
	Subgraph string
}

// Edge represents a connection between nodes
type Edge struct {
	From      string
	To        string
	Label     string
	ArrowType string // -->, ==>, -.->
	Line      int
}

// Subgraph represents a subgraph grouping
type Subgraph struct {
	ID     string
	Title  string
	Quoted bool
	Line   int
	Nodes  []string
}

// Diagram represents a parsed mermaid diagram
type Diagram struct {
	Direction string // LR, TD, etc.
	Nodes     map[string]*Node
	Edges     []*Edge
	Subgraphs []*Subgraph
	ClassDefs map[string]string
	Classes   map[string][]string // node -> classes
	RawLines  []string
}

// ParseMermaid parses mermaid flowchart code into a structured diagram
func ParseMermaid(code string) (*Diagram, error) {
	diagram := &Diagram{
		Nodes:     make(map[string]*Node),
		Edges:     []*Edge{},
		Subgraphs: []*Subgraph{},
		ClassDefs: make(map[string]string),
		Classes:   make(map[string][]string),
		RawLines:  []string{},
	}

	lines := strings.Split(code, "\n")
	diagram.RawLines = lines

	var currentSubgraph *Subgraph

	// Regex patterns
	directionRe := regexp.MustCompile(`^flowchart\s+(LR|TD|TB|RL|BT)`)
	// Node patterns: [], [(...)], ([...]), [[...]], {...}, ((...)), (...)
	nodeRe := regexp.MustCompile(`^\s*([A-Za-z0-9_]+)(\[|\[\(|\(\[|\[\[|\{|\(\(|\()(.+?)(\]|\)\]|\]\)|\]\]|\}|\)\)|\))`)
	edgeRe := regexp.MustCompile(`^\s*([A-Za-z0-9_]+)\s*(-->|==>|-.->|-.-|--)\s*(\|[^|]+\|)?\s*([A-Za-z0-9_]+)`)
	subgraphStartRe := regexp.MustCompile(`^\s*subgraph\s+([A-Za-z0-9_-]+)\s*\[?"?([^"\]]*)"?\]?`)
	subgraphEndRe := regexp.MustCompile(`^\s*end\s*$`)
	classDefRe := regexp.MustCompile(`^\s*classDef\s+([A-Za-z0-9_]+)\s+(.+)`)
	classRe := regexp.MustCompile(`^\s*class\s+([A-Za-z0-9_,\s]+)\s+([A-Za-z0-9_]+)`)

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if strings.HasPrefix(line, "%%") || line == "" {
			continue
		}

		// Check direction
		if matches := directionRe.FindStringSubmatch(line); matches != nil {
			diagram.Direction = matches[1]
			continue
		}

		// Check subgraph start
		if matches := subgraphStartRe.FindStringSubmatch(line); matches != nil {
			quoted := strings.Contains(line, `"`)
			currentSubgraph = &Subgraph{
				ID:     matches[1],
				Title:  strings.Trim(matches[2], `"`),
				Quoted: quoted,
				Line:   lineNum + 1,
				Nodes:  []string{},
			}
			diagram.Subgraphs = append(diagram.Subgraphs, currentSubgraph)
			continue
		}

		// Check subgraph end
		if subgraphEndRe.MatchString(line) {
			currentSubgraph = nil
			continue
		}

		// Check classDef
		if matches := classDefRe.FindStringSubmatch(line); matches != nil {
			diagram.ClassDefs[matches[1]] = matches[2]
			continue
		}

		// Check class application
		if matches := classRe.FindStringSubmatch(line); matches != nil {
			nodes := strings.Split(matches[1], ",")
			className := matches[2]
			for _, node := range nodes {
				node = strings.TrimSpace(node)
				diagram.Classes[node] = append(diagram.Classes[node], className)
			}
			continue
		}

		// Check edges (must check before nodes since edges contain node references)
		if matches := edgeRe.FindStringSubmatch(line); matches != nil {
			label := ""
			if matches[3] != "" {
				label = strings.Trim(matches[3], "|")
			}
			edge := &Edge{
				From:      matches[1],
				To:        matches[4],
				ArrowType: matches[2],
				Label:     label,
				Line:      lineNum + 1,
			}
			diagram.Edges = append(diagram.Edges, edge)
			continue
		}

		// Check nodes
		if matches := nodeRe.FindStringSubmatch(line); matches != nil {
			shape := "rectangle"
			switch matches[2] {
			case "[(":
				shape = "cylinder"
			case "([":
				shape = "stadium"
			case "[[":
				shape = "double_rectangle"
			case "{":
				shape = "diamond"
			case "((":
				shape = "circle"
			case "(":
				shape = "rounded"
			}

			node := &Node{
				ID:    matches[1],
				Label: matches[3],
				Shape: shape,
				Line:  lineNum + 1,
			}

			if currentSubgraph != nil {
				node.Subgraph = currentSubgraph.ID
				currentSubgraph.Nodes = append(currentSubgraph.Nodes, node.ID)
			}

			diagram.Nodes[node.ID] = node
		}
	}

	// Apply classes to nodes
	for nodeID, classes := range diagram.Classes {
		if node, ok := diagram.Nodes[nodeID]; ok {
			node.Classes = classes
		}
	}

	return diagram, nil
}

// HasNodeWithLabel checks if a node with the given label exists
func (d *Diagram) HasNodeWithLabel(label string) bool {
	label = strings.ToLower(label)
	for _, node := range d.Nodes {
		if strings.ToLower(node.Label) == label {
			return true
		}
		// Also check if label is contained (for partial matches)
		if strings.Contains(strings.ToLower(node.Label), label) {
			return true
		}
	}
	return false
}

// GetOrphanNodes returns nodes with no incoming or outgoing edges
func (d *Diagram) GetOrphanNodes() []*Node {
	connected := make(map[string]bool)
	for _, edge := range d.Edges {
		connected[edge.From] = true
		connected[edge.To] = true
	}

	orphans := []*Node{}
	for id, node := range d.Nodes {
		if !connected[id] {
			orphans = append(orphans, node)
		}
	}
	return orphans
}
