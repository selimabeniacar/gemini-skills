package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/flowlint/internal/linter"
	"github.com/user/flowlint/internal/parser"
)

var (
	lintFix    bool
	lintOutput string
)

var lintCmd = &cobra.Command{
	Use:   "lint <diagram.md>",
	Short: "Check diagram against style guide",
	Long: `Checks the Mermaid diagram against the style guide and reports
any violations.

Checks performed:
- Sync calls use ==> arrows
- Async calls use -.-> arrows
- All subgraph titles are quoted
- classDef styles are defined and applied
- No abbreviations in node labels
- No orphan nodes (unconnected)
- No duplicate node IDs

Use --fix to automatically fix issues where possible.`,
	Args: cobra.ExactArgs(1),
	RunE: runLint,
}

func init() {
	lintCmd.Flags().BoolVar(&lintFix, "fix", false, "Automatically fix issues")
	lintCmd.Flags().StringVarP(&lintOutput, "output", "o", "", "Output file for fixed diagram")
}

func runLint(cmd *cobra.Command, args []string) error {
	diagramPath := args[0]

	// Read the markdown file
	content, err := os.ReadFile(diagramPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Extract mermaid code block
	mermaidCode, err := parser.ExtractMermaid(string(content))
	if err != nil {
		return fmt.Errorf("failed to extract mermaid: %w", err)
	}

	// Parse the mermaid diagram
	diagram, err := parser.ParseMermaid(mermaidCode)
	if err != nil {
		return fmt.Errorf("failed to parse mermaid: %w", err)
	}

	// Run linting rules
	issues := linter.Lint(diagram)

	// Print issues
	hasErrors := false
	for _, issue := range issues {
		switch issue.Severity {
		case linter.SeverityError:
			fmt.Printf("❌ ERROR: %s\n", issue.Message)
			if issue.Line > 0 {
				fmt.Printf("   Line %d: %s\n", issue.Line, issue.Context)
			}
			if issue.Suggestion != "" {
				fmt.Printf("   Suggestion: %s\n", issue.Suggestion)
			}
			hasErrors = true
		case linter.SeverityWarning:
			fmt.Printf("⚠️  WARNING: %s\n", issue.Message)
			if issue.Suggestion != "" {
				fmt.Printf("   Suggestion: %s\n", issue.Suggestion)
			}
		}
		fmt.Println()
	}

	if len(issues) == 0 {
		fmt.Println("✓ No style issues found")
		return nil
	}

	// Apply fixes if requested
	if lintFix {
		fixedCode, fixCount := linter.Fix(mermaidCode, issues)
		if fixCount > 0 {
			fmt.Printf("\n✓ Applied %d automatic fixes\n", fixCount)

			// Replace mermaid block in original content
			fixedContent := parser.ReplaceMermaid(string(content), fixedCode)

			// Write output
			outputPath := lintOutput
			if outputPath == "" {
				outputPath = diagramPath
			}

			if err := os.WriteFile(outputPath, []byte(fixedContent), 0644); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			fmt.Printf("✓ Written to %s\n", outputPath)
		} else {
			fmt.Println("\nNo automatic fixes available for remaining issues")
		}
	}

	if hasErrors && !lintFix {
		return fmt.Errorf("found %d issues (use --fix to auto-fix)", len(issues))
	}

	return nil
}
