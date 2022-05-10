package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "cli",
		Short: "The DEMO Blockchain  CLI",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("incorrect usage")
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(balancesCmd())
	rootCmd.AddCommand(txCmd())

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
