package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

const (
	flagDataDir = "datadir"
)

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

//为命令添加必须的 默认参数
func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagDataDir, "", "Absolute path to the node data dir where the DB will/is stored")
	_ = cmd.MarkFlagRequired(flagDataDir)

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

	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(statusCmd())
	rootCmd.AddCommand(balancesCmd())
	rootCmd.AddCommand(runCmd())

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
