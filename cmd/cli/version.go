package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of cli",
	Long:  "Print the version number of cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cli version %s\n", version)
	},
}
