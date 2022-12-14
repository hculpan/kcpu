/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/hculpan/kcpu/kasm/assembler"
	"github.com/spf13/cobra"
)

var outputFilename string = ""

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Assembles a KASM assembly language program",
	Long:  `Assembles a KASM assembly language program`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Must give a KASM program file name")
			return
		}

		if !fileExists(args[0]) {
			fmt.Printf("File '%s' not found\n", args[0])
			return
		}

		buildConfig := assembler.BuildConfig{
			InputFilename:       args[0],
			OutputAssembledFile: true,
		}
		if len(outputFilename) > 0 {
			buildConfig.OutputFilename = outputFilename
		} else {
			buildConfig.SetDefaultOutputFilename()
		}

		assembler.BuildFile(buildConfig)
	},
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//buildCmd.Flags().StringVar(&outputDir, "dir", "d", "Directory where files will be output")
	buildCmd.Flags().StringVar(&outputFilename, "output", "", "Name of generated file")
}
