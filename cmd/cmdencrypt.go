/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt --password=Your-password --input=/path/of/single_file --output=/path/of/target_file",
	Long: `you can set env variable ZSTDZIPPASSWORD as your default password for security.
	[export ZSTDZIPPASSWORD=Your-password]`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			FatalError("--input= cannot be empty")
		} else {
			Input = AbsToSlash(Input)
		}

		if Output == "" {
			FatalError("--output= cannot be empty")
		} else {
			Output = AbsToSlash(Output)
		}

		if Input == Output {
			FatalError("--input and --output cannot be same")
		}

		if Password == "" {
			Password = GetEnv("ZSTDZIPPASSWORD", "")
			PrintlnDebug("--password= is empty, will use the env variable ZSTDZIPPASSWORD as default.")
		}

		if Password == "" {
			FatalError("--password= cannot be empty")
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		PrintArgs()

		if _, err := os.Stat(Input); err != nil {
			FatalError(err)
		}

		NewCryptFile(Input, Output, Password).AESEncode()

	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.PersistentFlags().StringVar(&Password, "password", "", "--password= or set env variable: export ZSTDZIPPASSWORD=")
}
