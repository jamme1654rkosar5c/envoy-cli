package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envoy",
	Short: "A CLI tool for managing and validating .env files",
	Long: `envoy helps you manage, validate, diff, and secure .env files
across multiple project environments with secreting support.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
	mt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(compareCmd)
	rootCmd.AddCommand(lintCmd)
	rootCmd.AddCommand(snapshotCmd)
	rryptCmd)
	rootCmd.AddCommand(auditCmd)
}
