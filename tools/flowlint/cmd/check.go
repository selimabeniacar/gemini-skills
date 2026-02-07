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
- All caches appear
- All external systems appear
- Internal steps (if present) appear as nodes`,
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

	for _, svc := range deps.Services {
		fmt.Printf("Service: %s\n", svc.Name)
		fmt.Println(strings.Repeat("-", 40))

		// Check sync dependencies
		fmt.Println("  Sync Dependencies:")
		for _, dep := range svc.Dependencies.Sync {
			total++
			if diagram.HasNodeWithLabel(dep.Name) {
				fmt.Printf("    ✓ %s\n", dep.Name)
				found++
			} else {
				fmt.Printf("    ✗ %s (MISSING)\n", dep.Name)
				missing = append(missing, fmt.Sprintf("%s > %s (sync, from %s:%d)", svc.Name, dep.Name, dep.SourceFile, dep.SourceLine))
			}
		}

		// Check async dependencies (Kafka topics)
		fmt.Println("  Kafka Topics:")
		for _, dep := range svc.Dependencies.Async {
			total++
			if diagram.HasNodeWithLabel(dep.Name) {
				fmt.Printf("    ✓ %s (%s)\n", dep.Name, dep.Direction)
				found++
			} else {
				fmt.Printf("    ✗ %s (%s) (MISSING)\n", dep.Name, dep.Direction)
				missing = append(missing, fmt.Sprintf("%s > %s (kafka %s, from %s:%d)", svc.Name, dep.Name, dep.Direction, dep.SourceFile, dep.SourceLine))
			}
		}

		// Check databases
		fmt.Println("  Databases:")
		for _, db := range svc.Databases {
			total++
			if diagram.HasNodeWithLabel(db.Name) {
				fmt.Printf("    ✓ %s\n", db.Name)
				found++
			} else {
				fmt.Printf("    ✗ %s (MISSING)\n", db.Name)
				missing = append(missing, fmt.Sprintf("%s > %s (database, from %s:%d)", svc.Name, db.Name, db.SourceFile, db.SourceLine))
			}
		}

		// Check caches
		fmt.Println("  Caches:")
		for _, cache := range svc.Caches {
			total++
			if diagram.HasNodeWithLabel(cache.Name) {
				fmt.Printf("    ✓ %s\n", cache.Name)
				found++
			} else {
				fmt.Printf("    ✗ %s (MISSING)\n", cache.Name)
				missing = append(missing, fmt.Sprintf("%s > %s (cache, from %s:%d)", svc.Name, cache.Name, cache.SourceFile, cache.SourceLine))
			}
		}

		// Check external systems
		fmt.Println("  External Systems:")
		for _, ext := range svc.External {
			total++
			if diagram.HasNodeWithLabel(ext.Name) {
				fmt.Printf("    ✓ %s\n", ext.Name)
				found++
			} else {
				fmt.Printf("    ✗ %s (MISSING)\n", ext.Name)
				missing = append(missing, fmt.Sprintf("%s > %s (external, from %s:%d)", svc.Name, ext.Name, ext.SourceFile, ext.SourceLine))
			}
		}

		// Check internal steps (if present)
		if len(svc.InternalSteps) > 0 {
			fmt.Println("  Internal Steps:")
			for _, step := range svc.InternalSteps {
				total++
				if diagram.HasNodeWithLabel(step.Name) {
					fmt.Printf("    ✓ %s\n", step.Name)
					found++
				} else {
					fmt.Printf("    ✗ %s (MISSING)\n", step.Name)
					missing = append(missing, fmt.Sprintf("%s > %s (internal step)", svc.Name, step.Name))
				}
			}
		}

		fmt.Println()
	}

	// Summary
	fmt.Println(strings.Repeat("=", 50))
	if total > 0 {
		fmt.Printf("Coverage: %d/%d (%.0f%%)\n", found, total, float64(found)/float64(total)*100)
	} else {
		fmt.Println("Coverage: 0/0 (no dependencies)")
	}

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
