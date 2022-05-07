package main

import (
	"demoBlockChain/database"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var balancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List balances",
	Long:  "List balances",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := database.NewStateFromDisk()
		if err != nil {
			log.Fatal(err)
		}
		defer state.CloseDB()

		for account, balance := range state.Balances {
			fmt.Printf("%s: %d\n", account, balance)
		}
	},
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
	//添加子命令
	balanceCmd.AddCommand(balancesListCmd)
	return balanceCmd
}
