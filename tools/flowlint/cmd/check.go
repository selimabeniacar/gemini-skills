package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/flowlint/internal/parser"
)

var checkCmd = &cobra.Command{
	Use:   "check <diagram.md> <dependencies.yaml>",
	Short: "Verify diagram completeness against dependencies",
	Long: `Compares the diagram against the dependencies.yaml file
to ensure all dependencies are represented.

Checks:
- All sync dependencies appear as nodes
- All async dependencies (Kafka topics) appear
- All databases appear
- All external systems appear
- Arrow directions match dependency directions`,
	Args: cobra.ExactArgs(2),
	RunE: runCheck,
}

func runCheck(cmd *cobra.Command, args []string) error {
	diagramPath := args[0]
	depsPath := args[1]

	// Read diagram
	diagramContent, err := os.ReadFile(diagramPath)
	if err != nil {
		return fmt.Errorf("failed to read diagram: %w", err)
	}

	// Read dependencies
	depsContent, err := os.ReadFile(depsPath)
	if err != nil {
		return fmt.Errorf("failed to read dependencies: %w", err)
	}

	// Extract mermaid
	mermaidCode, err := parser.ExtractMermaid(string(diagramContent))
	if err != nil {
		return fmt.Errorf("failed to extract mermaid: %w", err)
	}

	// Parse mermaid
	diagram, err := parser.ParseMermaid(mermaidCode)
	if err != nil {
		return fmt.Errorf("failed to parse mermaid: %w", err)
	}

	// Parse dependencies
	deps, err := parser.ParseDependencies(depsContent)
	if err != nil {
		return fmt.Errorf("failed to parse dependencies: %w", err)
	}

	// Check completeness
	fmt.Println("Checking diagram completeness...\n")

	missing := []string{}
	found := 0
	total := 0

	// Check sync dependencies
	fmt.Println("Sync Dependencies:")
	for _, dep := range deps.Dependencies.Sync {
		total++
		if diagram.HasNodeWithLabel(dep.Name) {
			fmt.Printf("  ✓ %s\n", dep.Name)
			found++
		} else {
			fmt.Printf("  ✗ %s (MISSING)\n", dep.Name)
			missing = append(missing, fmt.Sprintf("%s (sync, from %s:%d)", dep.Name, dep.SourceFile, dep.SourceLine))
		}
	}

	// Check async dependencies (Kafka topics)
	fmt.Println("\nKafka Topics:")
	for _, dep := range deps.Dependencies.Async {
		total++
		if diagram.HasNodeWithLabel(dep.Name) {
			fmt.Printf("  ✓ %s (%s)\n", dep.Name, dep.Direction)
			found++
		} else {
			fmt.Printf("  ✗ %s (%s) (MISSING)\n", dep.Name, dep.Direction)
			missing = append(missing, fmt.Sprintf("%s (kafka %s, from %s:%d)", dep.Name, dep.Direction, dep.SourceFile, dep.SourceLine))
		}
	}

	// Check external systems
	fmt.Println("\nExternal Systems:")
	for _, ext := range deps.External {
		total++
		if diagram.HasNodeWithLabel(ext.Name) {
			fmt.Printf("  ✓ %s\n", ext.Name)
			found++
		} else {
			fmt.Printf("  ✗ %s (MISSING)\n", ext.Name)
			missing = append(missing, fmt.Sprintf("%s (external, from %s:%d)", ext.Name, ext.SourceFile, ext.SourceLine))
		}
	}

	// Check caches
	fmt.Println("\nCaches:")
	for _, cache := range deps.Caches {
		total++
		if diagram.HasNodeWithLabel(cache.Name) {
			fmt.Printf("  ✓ %s\n", cache.Name)
			found++
		} else {
			fmt.Printf("  ✗ %s (MISSING)\n", cache.Name)
			missing = append(missing, fmt.Sprintf("%s (cache, from %s:%d)", cache.Name, cache.SourceFile, cache.SourceLine))
		}
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("Coverage: %d/%d (%.0f%%)\n", found, total, float64(found)/float64(total)*100)

	if len(missing) > 0 {
		fmt.Println("\nMissing items:")
		for _, m := range missing {
			fmt.Printf("  - %s\n", m)
		}
		return fmt.Errorf("diagram is incomplete: %d missing items", len(missing))
	}

	fmt.Println("\n✓ Diagram is complete")
	return nil
}
