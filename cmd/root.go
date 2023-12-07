/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var (
	timeBoot time.Time
	Input    string
	Output   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zstdzip [zip | unzip] [options]",
	Short: "(de)compress file(s) in zip format with ZSTD method",
	Long: `Compress: zstdzip zip --input=/path/of/file or folder  --output=/path/of/abc.zip  --speed=0|1|6|9
	Decompress: zstdzip unzip --input=abc.zip  --output=/path/of/target/folder
	or you can use https://github.com/mcmilk/7-Zip-zstd to unzip`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n *** elapse:", time.Since(timeBoot), "***")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	timeBoot = GetTimeNow()

	numCPU = runtime.NumCPU()
	if numCPU > 32 {
		numCPU = 32
	}

	rootCmd.PersistentFlags().StringVar(&Input, "input", "", "source file or folder")
	rootCmd.PersistentFlags().StringVar(&Output, "output", "", "target file")

}
