package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/plexusone/graphfs/pkg/schema"
	"github.com/plexusone/graphfs/pkg/store"
)

func validateCmd() *cobra.Command {
	var requireLabel bool
	var allowedNodeTypes []string
	var allowedEdgeTypes []string

	cmd := &cobra.Command{
		Use:   "validate <path>",
		Short: "Validate a graph database",
		Long: `Validate a graph database at the specified path.

Checks for:
  - Valid node structure (ID, type, optional label)
  - Valid edge structure (from, to, type, confidence)
  - Referential integrity (edge endpoints exist)
  - Filesystem-safe node IDs`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			// Check if path exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return fmt.Errorf("path does not exist: %s", path)
			}

			// Load the graph
			fs, err := store.NewFSStore(path)
			if err != nil {
				return fmt.Errorf("failed to open store: %w", err)
			}

			g, err := fs.LoadGraph()
			if err != nil {
				return fmt.Errorf("failed to load graph: %w", err)
			}

			// Configure validator
			v := schema.NewValidator()
			v.RequireNodeLabel = requireLabel
			if len(allowedNodeTypes) > 0 {
				v.AllowedNodeTypes = allowedNodeTypes
			}
			if len(allowedEdgeTypes) > 0 {
				v.AllowedEdgeTypes = allowedEdgeTypes
			}

			// Validate
			errs := v.ValidateGraph(g)

			if len(errs) == 0 {
				fmt.Printf("✓ Valid graph: %d nodes, %d edges\n", g.NodeCount(), g.EdgeCount())
				return nil
			}

			fmt.Printf("✗ Found %d validation error(s):\n\n", len(errs))
			for i, err := range errs {
				fmt.Printf("  %d. %v\n", i+1, err)
			}
			fmt.Println()

			return fmt.Errorf("validation failed with %d error(s)", len(errs))
		},
	}

	cmd.Flags().BoolVar(&requireLabel, "require-label", false, "Require all nodes to have a label")
	cmd.Flags().StringSliceVar(&allowedNodeTypes, "node-types", nil, "Allowed node types (comma-separated)")
	cmd.Flags().StringSliceVar(&allowedEdgeTypes, "edge-types", nil, "Allowed edge types (comma-separated)")

	return cmd
}
