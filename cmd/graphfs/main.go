// Package main provides the graphfs CLI.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "graphfs",
		Short:   "Git-friendly filesystem graph database",
		Version: version,
	}

	rootCmd.AddCommand(validateCmd())
	rootCmd.AddCommand(formatCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
