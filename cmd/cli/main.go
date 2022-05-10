package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "cli",
		Short: "The DEMO Blockchain  CLI",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("incorrect usage")
		},
		Run: func(cmd *cobra.Command, args []string) {
			//do nothing
		},
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(balancesCmd())
	rootCmd.AddCommand(txCmd())

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
