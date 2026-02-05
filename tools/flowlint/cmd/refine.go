package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/flowlint/internal/linter"
	"github.com/user/flowlint/internal/parser"
)

var (
	refineOutput       string
	refineSkipValidate bool
)

var refineCmd = &cobra.Command{
	Use:   "refine <diagram.md> <dependencies.yaml>",
	Short: "Run full refinement pipeline",
	Long: `Runs the complete refinement pipeline:

1. Validate - Check Mermaid syntax with mmdc
2. Lint - Check style guide and auto-fix
3. Check - Verify completeness against dependencies

Use --skip-validate if mmdc is not installed.
Use --output to specify output file (defaults to overwriting input).`,
	Args: cobra.ExactArgs(2),
	RunE: runRefine,
}

func init() {
	refineCmd.Flags().StringVarP(&refineOutput, "output", "o", "", "Output file for refined diagram")
	refineCmd.Flags().BoolVar(&refineSkipValidate, "skip-validate", false, "Skip mmdc validation")
}

func runRefine(cmd *cobra.Command, args []string) error {
	diagramPath := args[0]
	depsPath := args[1]

	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║          flowlint refinement pipeline            ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")
	fmt.Println()

	// Read files
	diagramContent, err := os.ReadFile(diagramPath)
	if err != nil {
		return fmt.Errorf("failed to read diagram: %w", err)
	}

	depsContent, err := os.ReadFile(depsPath)
	if err != nil {
		return fmt.Errorf("failed to read dependencies: %w", err)
	}

	// Step 1: Validate (optional)
	fmt.Println("Step 1: Syntax Validation")
	fmt.Println("─────────────────────────")
	if refineSkipValidate {
		fmt.Println("⏭  Skipped (--skip-validate)")
	} else {
		// Run validation
		if err := runValidate(cmd, []string{diagramPath}); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}
	fmt.Println()

	// Step 2: Lint and fix
	fmt.Println("Step 2: Style Linting")
	fmt.Println("─────────────────────")

	mermaidCode, err := parser.ExtractMermaid(string(diagramContent))
	if err != nil {
		return fmt.Errorf("failed to extract mermaid: %w", err)
	}

	diagram, err := parser.ParseMermaid(mermaidCode)
	if err != nil {
		return fmt.Errorf("failed to parse mermaid: %w", err)
	}

	issues := linter.Lint(diagram)

	errorCount := 0
	warningCount := 0
	for _, issue := range issues {
		switch issue.Severity {
		case linter.SeverityError:
			fmt.Printf("  ❌ %s\n", issue.Message)
			errorCount++
		case linter.SeverityWarning:
			fmt.Printf("  ⚠️  %s\n", issue.Message)
			warningCount++
		}
	}

	if len(issues) == 0 {
		fmt.Println("  ✓ No style issues found")
	} else {
		fmt.Printf("\n  Found %d errors, %d warnings\n", errorCount, warningCount)

		// Apply fixes
		fixedCode, fixCount := linter.Fix(mermaidCode, issues)
		if fixCount > 0 {
			fmt.Printf("  ✓ Applied %d automatic fixes\n", fixCount)
			diagramContent = []byte(parser.ReplaceMermaid(string(diagramContent), fixedCode))
		}
	}
	fmt.Println()

	// Step 3: Completeness check
	fmt.Println("Step 3: Completeness Check")
	fmt.Println("──────────────────────────")

	deps, err := parser.ParseDependencies(depsContent)
	if err != nil {
		return fmt.Errorf("failed to parse dependencies: %w", err)
	}

	// Re-parse diagram with fixes applied
	mermaidCode, _ = parser.ExtractMermaid(string(diagramContent))
	diagram, _ = parser.ParseMermaid(mermaidCode)

	missing := checkCompleteness(diagram, deps)
	if len(missing) > 0 {
		fmt.Printf("  ⚠️  Missing %d items:\n", len(missing))
		for _, m := range missing {
			fmt.Printf("    - %s\n", m)
		}
	} else {
		fmt.Println("  ✓ All dependencies represented")
	}
	fmt.Println()

	// Write output
	outputPath := refineOutput
	if outputPath == "" {
		outputPath = diagramPath
	}

	if err := os.WriteFile(outputPath, diagramContent, 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// Summary
	fmt.Println("════════════════════════════════════════════════════")
	if errorCount == 0 && len(missing) == 0 {
		fmt.Println("✓ Refinement complete - diagram is ready")
	} else {
		fmt.Println("⚠️  Refinement complete with issues")
		if errorCount > 0 {
			fmt.Printf("   %d style errors remain (manual fix required)\n", errorCount)
		}
		if len(missing) > 0 {
			fmt.Printf("   %d missing items (regenerate diagram)\n", len(missing))
		}
	}
	fmt.Printf("\nOutput: %s\n", outputPath)

	return nil
}

func checkCompleteness(diagram *parser.Diagram, deps *parser.Dependencies) []string {
	missing := []string{}

	// Check sync dependencies
	for _, dep := range deps.Dependencies.Sync {
		if !diagram.HasNodeWithLabel(dep.Name) {
			missing = append(missing, fmt.Sprintf("%s (sync)", dep.Name))
		}
	}

	// Check async dependencies
	for _, dep := range deps.Dependencies.Async {
		if !diagram.HasNodeWithLabel(dep.Name) {
			missing = append(missing, fmt.Sprintf("%s (kafka)", dep.Name))
		}
	}

	// Check external systems
	for _, ext := range deps.External {
		if !diagram.HasNodeWithLabel(ext.Name) {
			missing = append(missing, fmt.Sprintf("%s (external)", ext.Name))
		}
	}

	// Check caches
	for _, cache := range deps.Caches {
		if !diagram.HasNodeWithLabel(cache.Name) {
			missing = append(missing, fmt.Sprintf("%s (cache)", cache.Name))
		}
	}

	return missing
}
