/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// zipCmd represents the zip command
var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "zip --input=/path/of/file(or folder) --output=/path/of/target.zip --speed=0|1|6|9",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			log.Fatal("--input= cannot be empty")
		} else {
			Input = AbsToSlash(Input)
		}

		if Output == "" {
			log.Fatal("--output= cannot be empty")
		} else {
			Output = AbsToSlash(Output)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		PrintArgs()

		entity := NewEntity(Input, Output).WithZipLevel(Speed).SetZipFileMap()
		entity.OpenZipTempFile().Compress().CloseZipTempFile()

		
	},
}

func init() {
	rootCmd.AddCommand(zipCmd)

	zipCmd.Flags().IntVar(&Speed, "speed", 1, "0/1/6/9, 0 is fastest without compression, 9 is slowest but most compression.")

}
