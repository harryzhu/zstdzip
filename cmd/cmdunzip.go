package cmd

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
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
		positionalArgs(args)
		//
		if Source == "" || Source == Target {
			FatalError("unzip", NewError("--source= or --target= cannot be empty or same"))
		}

		if Target == "" {
			ext := filepath.Ext(Source)
			autoTarget := strings.TrimSuffix(Source, ext)
			Target = strings.TrimSuffix(autoTarget, ".zip")
			Target = strings.TrimSuffix(Target, ".zstd")
			if FileExists(Target) {
				Target = strings.Join([]string{Target, "(2)"}, " ")
			}
		}
		if Target != "" {
			MakeDirs(Target)
		}

		argsValidate()
		bootstrap()

		timeBoot = time.Now()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if FileExists(Source) {
			taskDecompressFile(Source)
		} else {
			FatalError("unzip", NewError("file does not exist: --source=", Source))
		}

		//

		wg := sync.WaitGroup{}

		for num := 1; num < 8; num++ {
			s := strings.Join([]string{Source, Int2Str(num)}, ".")
			if FileExists(s) {
				DebugInfo("unzip: Extracting", s)
				wg.Add(1)
				go func(s string) {
					defer wg.Done()
					taskDecompressFile(s)
				}(s)
				if IsSerial {
					wg.Wait()
				}
			}
		}
		wg.Wait()

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
