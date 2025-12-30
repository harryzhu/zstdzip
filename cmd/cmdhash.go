package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// hashCmd represents the hash command
var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "--sum=blake3 | xxhash | md5 | sha1 | sha256 --source=path/of/file.txt",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if IsDebug {
			PrintArgs("source", "sum")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		Sum = strings.ToLower(Sum)

		fmt.Printf("%s: %s\n", Sum, hashFile(Sum))
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
	rootCmd.MarkFlagRequired("source")

	hashCmd.Flags().StringVar(&Sum, "sum", "sha256", "sum method: md5, sha1, sha256, blake3, xxhash")
}
