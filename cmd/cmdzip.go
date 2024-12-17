/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

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

		if fh, err := os.Stat(Output); err == nil {
			if TotalSize > 0 && fh.Size() > 0 {
				fmt.Printf("*** Compression Rate: %.4f ***\n",
					DivideFloat64(float64(fh.Size()), float64(TotalSize)))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(zipCmd)

	zipCmd.Flags().IntVar(&Speed, "speed", 1, "0/1/6/9, 0 is fastest without compression, 9 is slowest but most compression.")

}
