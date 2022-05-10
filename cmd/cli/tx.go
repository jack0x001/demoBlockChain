package main

import (
	"demoBlockChain/database"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

const flagFrom = "from"
const flagTo = "to"
const flagValue = "value"
const flagData = "data"

func txCmd() *cobra.Command {
	tx := &cobra.Command{
		Use:   "tx",
		Short: "Transaction commands",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("incorrect usage")
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	tx.AddCommand(txAddCmd())

	return tx
}

//tbb tx add --from=andrej --to=andrej --value=3
func txAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a transaction",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			value, _ := cmd.Flags().GetUint(flagValue)
			data, _ := cmd.Flags().GetString(flagData)

			fromAccount := database.NewAccount(from)
			toAccount := database.NewAccount(to)
			tx := database.NewTx(fromAccount, toAccount, value, data)
			state, err := database.NewStateFromDisk()
			if err != nil {
				log.Fatal(err)
			}

			defer state.CloseDB()
			
			err = state.AddTx(tx)
			if err != nil {
				log.Fatal(err)
			}
			err = state.Persist()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("tx added, from:", from, "to:", to, "value:", value, "data:", data)
		},
	}

	cmd.Flags().String(flagFrom, "", "From what account to send tokens")
	_ = cmd.MarkFlagRequired(flagFrom)

	cmd.Flags().String(flagTo, "", "To what account to send tokens")
	_ = cmd.MarkFlagRequired(flagTo)

	cmd.Flags().Uint(flagValue, 0, "How many tokens to send")
	_ = cmd.MarkFlagRequired(flagValue)

	cmd.Flags().String(flagData, "", "Possible values: 'reward'")
	return cmd
}
