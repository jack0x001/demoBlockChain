package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

const (
	flagDataDir = "datadir"
	flagIP      = "ip"
	flagPort    = "port"
)

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

//为命令添加必须的 默认参数
func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagDataDir, "", "Absolute path to the node data dir where the DB will/is stored")
	cmd.Flags().String(flagIP, "", "ip address of node")
	cmd.Flags().String(flagPort, "", "port number of node")
	_ = cmd.MarkFlagRequired(flagDataDir)
	_ = cmd.MarkFlagRequired(flagIP)
	_ = cmd.MarkFlagRequired(flagPort)

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
	rootCmd.AddCommand(runCmd())

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
