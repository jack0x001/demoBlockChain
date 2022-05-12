package main

import (
	"demoBlockChain/database"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func balancesListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "List balances",
		Long:  "List balances",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, err := cmd.Flags().GetString("datadir")
			if err != nil {
				log.Fatal(err)
			}
			state, err := database.NewStateFromDisk(dataDir)
			if err != nil {
				log.Fatal(err)
			}
			defer state.CloseDB()

			for account, balance := range state.Balances {
				fmt.Printf("%s: %d\n", account, balance)
			}
		},
	}

	addDefaultRequiredFlags(&cmd)
	return &cmd
}

func balancesCmd() *cobra.Command {
	balanceCmd := &cobra.Command{
		Use:   "balances",
		Short: "actions about balances",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("incorrect usage")
		},
		Run: func(cmd *cobra.Command, args []string) {
			// do nothing
			// 上面的PreRunE 会纠正拼写错误
		},
	}

	//子命令
	balanceCmd.AddCommand(balancesListCmd())
	return balanceCmd
}
