package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/plexusone/graphfs/pkg/format"
)

func formatCmd() *cobra.Command {
	var dryRun bool
	var quiet bool

	cmd := &cobra.Command{
		Use:   "format <path>",
		Short: "Format JSON files in a graph database",
		Long: `Format all JSON files in a graph database to canonical form.

Canonical formatting ensures:
  - Sorted keys (alphabetical)
  - Consistent 2-space indentation
  - No trailing newlines

This produces deterministic output for clean git diffs.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			// Check if path exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return fmt.Errorf("path does not exist: %s", path)
			}

			var formatted, unchanged, failed int

			// Process nodes and edges directories
			dirs := []string{
				filepath.Join(path, "nodes"),
				filepath.Join(path, "edges"),
			}

			for _, dir := range dirs {
				entries, err := os.ReadDir(dir)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
					return fmt.Errorf("failed to read directory %s: %w", dir, err)
				}

				for _, entry := range entries {
					if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
						continue
					}

					filePath := filepath.Join(dir, entry.Name())
					wasFormatted, err := formatFile(filePath, dryRun)
					if err != nil {
						if !quiet {
							fmt.Fprintf(os.Stderr, "✗ %s: %v\n", filePath, err)
						}
						failed++
						continue
					}

					if wasFormatted {
						if !quiet {
							if dryRun {
								fmt.Printf("would format: %s\n", filePath)
							} else {
								fmt.Printf("formatted: %s\n", filePath)
							}
						}
						formatted++
					} else {
						unchanged++
					}
				}
			}

			// Summary
			if !quiet {
				fmt.Println()
				if dryRun {
					fmt.Printf("Would format %d file(s), %d already canonical", formatted, unchanged)
				} else {
					fmt.Printf("Formatted %d file(s), %d already canonical", formatted, unchanged)
				}
				if failed > 0 {
					fmt.Printf(", %d failed", failed)
				}
				fmt.Println()
			}

			if failed > 0 {
				return fmt.Errorf("%d file(s) failed to format", failed)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be formatted without writing")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress output except errors")

	return cmd
}

// formatFile formats a JSON file to canonical form.
// Returns true if the file was modified, false if already canonical.
func formatFile(path string, dryRun bool) (bool, error) {
	// Read original content
	original, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("reading file: %w", err)
	}

	// Parse and re-marshal to canonical form
	var data any
	if err := format.UnmarshalCanonical(original, &data); err != nil {
		return false, fmt.Errorf("parsing JSON: %w", err)
	}

	canonical, err := format.MarshalCanonical(data)
	if err != nil {
		return false, fmt.Errorf("marshaling JSON: %w", err)
	}

	// Check if already canonical
	if string(original) == string(canonical) {
		return false, nil
	}

	// Write if not dry run
	if !dryRun {
		if err := os.WriteFile(path, canonical, 0600); err != nil {
			return false, fmt.Errorf("writing file: %w", err)
		}
	}

	return true, nil
}
