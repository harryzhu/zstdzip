/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var(
	Async bool
)

// unzipCmd represents the unzip command
var unzipCmd = &cobra.Command{
	Use:   "unzip",
	Short: "unzip --input=/path/of/file.zip --output=/path/of/extract/folder",
	Long:  `or you can use https://github.com/mcmilk/7-Zip-zstd to unzip`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			log.Fatal("--input= cannot be empty")
		} else {
			Input = AbsToSlash(Input)
		}
		if Output == "" {
			Output = "./" + strings.Replace(filepath.Base(Input), filepath.Ext(Input), "", 1)
			fmt.Println("you can use --output= to specify the extract folder, default is current folder")
		}
		Output = AbsToSlash(Output)

	},
	Run: func(cmd *cobra.Command, args []string) {
		PrintArgs()

		entity := NewEntity(Input, Output)
		if Async == true{
			entity.DecompressAsync()
		}else{
			entity.Decompress()
		}
		

	},
}

func init() {
	rootCmd.AddCommand(unzipCmd)
	unzipCmd.Flags().BoolVarP(&Async, "async", "",false, "unzip files in async mode")

}
