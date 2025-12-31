package cmd

import (
	"archive/zip"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/spf13/cobra"
)

var (
	DeComTotalNum    int32
	DeComDirInfoList []*zip.File
	DeComLock        sync.Mutex
)

// unzipCmd represents the unzip command
var unzipCmd = &cobra.Command{
	Use:   "unzip",
	Short: "--source=/path/to/file.zstd.zip --target=/path/to/extract-dir ",
	Long: `optional args:
	--password= : use this password for decryption;
	--serial : decompress files one-by-one;
	--min-age= :, 
	--max-age= :, 
	--min-size-mb= :, 
	--max-size-mb= :, 
	--ext= : regular pattern(Case Insensitive)`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		PrintArgs("source", "target", "threads", "serial", "min-age", "max-age", "min-size-mb", "max-size-mb", "ext")
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

				decompressFile(decomFile)
			}(decomFile)

			if IsSerial {
				wg.Wait()
			}
		}
		wg.Wait()

		// sync dir's modTime and modePerm
		decompressDirMod()

		PrintSpinner(Int2Str(int(atomic.LoadInt32(&DeComTotalNum))))

	},
}

func init() {
	rootCmd.AddCommand(unzipCmd)
	rootCmd.MarkFlagRequired("source")
	rootCmd.MarkFlagRequired("target")

	unzipCmd.Flags().StringVar(&RegExt, "ext", "", "regex pattern of file extension(Case Insensitive). i.e.: .(mp4|txt|png)")
	//
	unzipCmd.Flags().StringVar(&MinAge, "min-age", "", "format: 2023-12-03,15:09:08, means 2023-12-03 15:09:08")
	unzipCmd.Flags().StringVar(&MaxAge, "max-age", "", "format: 2023-12-25,23:59:59, means 2023-12-25 23:59:59")
	//
	unzipCmd.Flags().Int64Var(&MinSizeMB, "min-size-mb", -1, "i.e.: 16 means 16MB, 16*1024*1024")
	unzipCmd.Flags().Int64Var(&MaxSizeMB, "max-size-mb", -1, "i.e.: 32 means 32MB, 32*1024*1024")
	//
	unzipCmd.Flags().BoolVar(&IsIgnoreEmptyDir, "ignore-empty-dir", true, "ignore empty folder")
	unzipCmd.Flags().BoolVar(&IsIgnoreDotFile, "ignore-dot-file", false, "ignore files start with dot(.), i.e.: .DS_Store .Thumb")
}
