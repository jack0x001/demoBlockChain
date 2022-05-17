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
			ip := cmd.Flag(flagIP).Value.String()
			port := cmd.Flag(flagPort).Value.String()

			n := node.NewNode(dataDir, ip, port, *node.BootNode)

			err := n.Run()
			if err != nil {
				log.Fatal("启动节点失败: ", err)
			}
		},
	}

	addDefaultRequiredFlags(&cmd)
	return &cmd
}
