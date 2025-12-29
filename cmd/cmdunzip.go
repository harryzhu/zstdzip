/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/spf13/cobra"
)

var (
	DeComTotalNum int32
)

// unzipCmd represents the unzip command
var unzipCmd = &cobra.Command{
	Use:   "unzip",
	Short: "--source=/path/to/file.zstd.zip --target=/path/to/extract-dir ",
	Long: `optional args:
	--password= : use this password for decryption;
	--serial : decompress files one-by-one`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		PrintArgs("source", "target", "threads", "serial")
		if Source == Target || Source == "" || Target == "" {
			FatalError("unzip", NewError("invalid --source= or --target="))
		}
		fmt.Println(" *** start:", timeBoot.Format("15:04:05"), "***")

	},
	Run: func(cmd *cobra.Command, args []string) {

		var decomFile string
		wg := sync.WaitGroup{}
		for idx := range 8 {
			if idx == 0 {
				decomFile = Source
			} else {
				decomFile = Source + "." + Int2Str(idx)
			}

			_, err := os.Stat(decomFile)

			if err != nil {
				continue
			}
			DebugInfo("unzip", decomFile)
			wg.Add(1)
			go func(decomFile string) {
				defer wg.Done()

				DecompressFile(decomFile)
			}(decomFile)

			if IsSerial {
				wg.Wait()
			}
		}
		wg.Wait()

		PrintSpinner(Int2Str(int(atomic.LoadInt32(&DeComTotalNum))))

	},
}

func init() {
	rootCmd.AddCommand(unzipCmd)
	unzipCmd.Flags().StringVar(&RegExt, "regext", "", "regex pattern of file extension(Case Insensitive). i.e.: .(mp4|txt|png)")

}
