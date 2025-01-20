/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var Sum string

// hashCmd represents the hash command
var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "hash --sum=blake3 | xxhash | md5 | sha1 | sha256 --input=path/of/file.txt",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		Sum = strings.ToLower(Sum)

		fmt.Println(Sum + ": ")
		fmt.Println(HashFile(Sum))
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
	rootCmd.MarkFlagRequired("input")

	hashCmd.Flags().StringVar(&Sum, "sum", "xxhash", "sum method: md5, sha1, sha256, blake3, xxhash")

}
