package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/flowlint/internal/parser"
)

var validateCmd = &cobra.Command{
	Use:   "validate <diagram.md>",
	Short: "Validate Mermaid syntax using mmdc",
	Long: `Extracts the Mermaid code block from a markdown file and
validates it using the Mermaid CLI (mmdc).

Requires mmdc to be installed:
  npm install -g @mermaid-js/mermaid-cli

Returns exit code 0 if valid, 1 if invalid.`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
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

	// Create temp file for mermaid code
	tmpDir, err := os.MkdirTemp("", "flowlint-")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	mermaidFile := filepath.Join(tmpDir, "diagram.mmd")
	if err := os.WriteFile(mermaidFile, []byte(mermaidCode), 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	outputFile := filepath.Join(tmpDir, "diagram.svg")

	// Check if mmdc is available
	mmdc, err := exec.LookPath("mmdc")
	if err != nil {
		// Try npx
		npx, npxErr := exec.LookPath("npx")
		if npxErr != nil {
			return fmt.Errorf("mmdc not found. Install with: npm install -g @mermaid-js/mermaid-cli")
		}
		mmdc = npx
		args = []string{"mmdc", "-i", mermaidFile, "-o", outputFile, "-q"}
	} else {
		args = []string{"-i", mermaidFile, "-o", outputFile, "-q"}
	}

	// Run mmdc
	execCmd := exec.Command(mmdc, args...)
	output, err := execCmd.CombinedOutput()
	if err != nil {
		// Parse error output for helpful messages
		errMsg := string(output)
		if strings.Contains(errMsg, "Parse error") {
			fmt.Println("❌ Mermaid syntax error:")
			fmt.Println(errMsg)
			return fmt.Errorf("diagram has syntax errors")
		}
		fmt.Println("❌ Validation failed:")
		fmt.Println(errMsg)
		return fmt.Errorf("mmdc validation failed")
	}

	fmt.Println("✓ Mermaid syntax is valid")
	return nil
}
