package main

import (
	"demoBlockChain/database"
	"fmt"
	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "status",
		Short: "Show the status of the chain",
		Long:  "Show the status of the chain",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, err := cmd.Flags().GetString(flagDataDir)
			if err != nil {
				return
			}
			state, err := database.NewStateFromDisk(dataDir)
			if err != nil {
				return
			}
			fmt.Printf("%+v\n", prettyPrint(state))
		}}

	addDefaultRequiredFlags(&cmd)
	return &cmd
}
