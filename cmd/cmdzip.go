package cmd

import (
	"fmt"
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
	--regext: regular pattern(Case Insensitive);
	--min-age: ignore files if file's last-modified-time is earlier than --min-age;
	--max-age: ignore files if file's last-modified-time is newer than --max-age;
	--min-size-mb: ignore files if file's size is less than --min-size-mb;
	--max-size-mb: ignore files if file's size is greater than --max-size-mb;`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		PrintArgs("source", "target", "threads", "serial", "level",
			"min-age", "max-age", "min-size-mb", "max-size-mb", "ext", "ignore-dot-file", "ignore-empty-dir")
		if strings.HasPrefix(Target, Source) || Source == "" || Target == "" {
			FatalError("zip", NewError("invalid --source= or --target="))
		}
		fmt.Println(" *** start:", timeBoot.Format("15:04:05"), "***")
	},
	Run: func(cmd *cobra.Command, args []string) {
		finfo, err := os.Stat(Source)
		if err != nil {
			FatalError("zip", err)
		}

		if finfo.IsDir() {
			compressDir()
			if IsDebug || IsDryRun {
				PrintSpinner(Int2Str(int(atomic.LoadInt32(&DeComTotalNum))))
			}
		} else {
			compressFile(finfo)
		}

	},
}

func init() {
	rootCmd.AddCommand(zipCmd)
	rootCmd.MarkFlagRequired("source")
	rootCmd.MarkFlagRequired("target")

	zipCmd.Flags().IntVar(&Level, "level", 1, "compress level: 0 | 1 | 2 | 3 ")
	//
	zipCmd.Flags().BoolVar(&IsIgnoreDotFile, "ignore-dot-file", false, "ignore files start with dot(.), i.e.: .DS_Store .Thumb")
	//
	zipCmd.Flags().StringVar(&RegExt, "ext", "", "regex pattern of file extension(Case Insensitive). i.e.: .(mp4|txt|png)")
	//
	zipCmd.Flags().StringVar(&MinAge, "min-age", "", "format: 2023-12-03,15:09:08, means 2023-12-03 15:09:08")
	zipCmd.Flags().StringVar(&MaxAge, "max-age", "", "format: 2023-12-25,23:59:59, means 2023-12-25 23:59:59")
	//
	zipCmd.Flags().Int64Var(&MinSizeMB, "min-size-mb", -1, "i.e.: 16 means 16MB, 16*1024*1024")
	zipCmd.Flags().Int64Var(&MaxSizeMB, "max-size-mb", -1, "i.e.: 32 means 32MB, 32*1024*1024")

}
