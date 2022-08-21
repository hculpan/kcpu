/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/hculpan/kcpu/cpu/executor"
	"github.com/spf13/cobra"
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Executes a KCPU binary program",
	Long:  `Executes a KCPU binary program`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Invalid number of args: Must pass KCPU executable file to process")
			return
		}
		if err := executor.Execute(args[0]); err != nil {
			if err.Error() == "halt" {
				fmt.Println("execution halted by program")
			} else {
				fmt.Println(fmt.Errorf("program execution terminated due to error: %s", err))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)

	// executeCmd.Flags().StringVar(&inputFile, "input", "", "config file (default is $HOME/.cobra.yaml)")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// executeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// executeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
