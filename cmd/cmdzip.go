/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"
	"sync/atomic"

	"github.com/spf13/cobra"
)

// zipCmd represents the zip command
var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "--source=/path/to/document-dir --target=/path/to/file.zstd.zip",
	Long: `optional args: 
	--password= : if password is not empty, will enable encrypt;
	--threads=8 : 
	--level=0|1|2|3 : 0 is fastest and lowest compression, 3 is slowest but highest compression;
	--serial: compress files one-by-one into a single archive, will disable parallel feature, 
	          better performance for HDD, but slower in SSD;
	files FILTER:
	--ignore-dot-file: ;
	--ignore-empty-dir: ;
	--regext: ;
	--min-age: ignore files if file's last-modified-time is earlier than --min-age;
	--max-age: ignore files if file's last-modified-time is newer than --max-age;
	--min-size-mb: ignore files if file's size is less than --min-size-mb;
	--max-size-mb: ignore files if file's size is greater than --max-size-mb;`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		PrintArgs("source", "target", "threads", "serial", "level")
		if strings.HasPrefix(Target, Source) || Source == "" || Target == "" {
			FatalError("zip", NewError("invalid --source= or --target="))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		finfo, err := os.Stat(Source)
		if err != nil {
			FatalError("zip", err)
		}

		if finfo.IsDir() {
			CompressDir()
			PrintSpinner(Int2Str(int(atomic.LoadInt32(&DeComTotalNum))))
		} else {
			CompressFile(finfo)
		}

	},
}

func init() {
	rootCmd.AddCommand(zipCmd)

	zipCmd.Flags().IntVar(&Level, "level", 1, "compress level: 0 | 1 | 2 | 3 ")
	//
	zipCmd.Flags().BoolVar(&IsIgnoreDotFile, "ignore-dot-file", false, "ignore files start with dot(.), i.e.: .DS_Store .Thumb")
	zipCmd.Flags().BoolVar(&IsIgnoreEmptyDir, "ignore-empty-dir", false, "ignore files start with dot(.), i.e.: .DS_Store .Thumb")
	//
	zipCmd.Flags().StringVar(&RegExt, "regext", "", "regex pattern of file extension(Case Insensitive). i.e.: .(mp4|txt|png)")
	//
	zipCmd.Flags().StringVar(&MinAge, "min-age", "", "format: 20231203150908, means 2023-12-03 15:09:08")
	zipCmd.Flags().StringVar(&MaxAge, "max-age", "", "format: 20231225235959, means 2023-12-25 23:59:59")
	//
	zipCmd.Flags().Int64Var(&MinSizeMB, "min-size-mb", -1, "i.e.: 16 means 16MB, 16*1024*1024")
	zipCmd.Flags().Int64Var(&MaxSizeMB, "max-size-mb", -1, "i.e.: 32 means 32MB, 32*1024*1024")

}
