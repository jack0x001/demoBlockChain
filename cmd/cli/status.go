package main

import (
	"demoBlockChain/database"
	"fmt"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of the chain",
	Long:  "Show the status of the chain",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := database.NewStateFromDisk()
		if err != nil {
			return
		}
		fmt.Printf("%+v\n", prettyPrint(state))
	},
}
