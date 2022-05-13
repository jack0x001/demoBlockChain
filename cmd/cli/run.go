package main

import (
	"demoBlockChain/node"
	"github.com/spf13/cobra"
	"log"
)

func runCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "run",
		Short: "Run the blockchain node",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir := cmd.Flag(flagDataDir).Value.String()
			err := node.Run(dataDir)
			if err != nil {
				log.Fatal("启动节点失败: ", err)
			}
		},
	}

	addDefaultRequiredFlags(&cmd)
	return &cmd
}
