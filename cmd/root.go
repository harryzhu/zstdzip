/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var (
	timeBoot  time.Time
	Input     string
	Output    string
	TotalSize int64
	LogStatus string
	YYYYMMDD  string
	Threads   int
	IsDebug   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zstdzip [zip | unzip] [options]",
	Short: "(de)compress file(s) in zip format with ZSTD method",
	Long: `Compress: zstdzip zip --input=/path/of/file or folder  --output=/path/of/abc.zip  --speed=0|1|6|9
	Decompress: zstdzip unzip --input=abc.zip  --output=/path/of/target/folder
	or you can use https://github.com/mcmilk/7-Zip-zstd to unzip`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println(" *** start:", timeBoot.Format("15:04:05"), "***")

		if Threads > numCPU || Threads < 1 {
			numCPU = runtime.NumCPU()
		} else {
			numCPU = Threads
		}

		runtime.GOMAXPROCS(numCPU)
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		timeDuration := time.Since(timeBoot)
		fmt.Println("\n *** elapse:", timeDuration, "***")
		if IsDebug {
			timeSeconds := timeDuration.Seconds()
			if timeSeconds <= 0 {
				timeSeconds = 1
			}
			if TotalSize > 0 {
				PrintSpeed(float64(TotalSize), float64(timeSeconds))
			}
		}

		if LogStatus != "" {
			fmt.Println(" *** global status:", GlobalStatus)
			result := make(map[string]string, 4)
			result["status"] = GlobalStatus
			result["start"] = timeBoot.Format("15:04:05")
			result["elapse"] = timeDuration.String()

			SaveJson(LogStatus+".zstdzip-"+YYYYMMDD+".log", result)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	timeBoot = GetTimeNow()
	YYYYMMDD = timeBoot.Format("20060102")
	numCPU = runtime.NumCPU()

	rootCmd.PersistentFlags().StringVar(&Input, "input", "", "source file or folder")
	rootCmd.PersistentFlags().StringVar(&Output, "output", "", "target file")
	rootCmd.PersistentFlags().StringVar(&LogStatus, "log", "", "log Global Status into this file")
	rootCmd.PersistentFlags().IntVar(&Threads, "threads", numCPU, "threads")
	rootCmd.PersistentFlags().BoolVarP(&IsDebug, "debug", "", false, "print debug info")

}
