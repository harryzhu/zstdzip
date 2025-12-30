package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var (
	timeBoot time.Time
	Source   string
	Target   string
	Threads  int
	Level    int
	//
	Password         string
	Sum              string
	TotalSize        int64
	Speed            int64
	IsDryRun         bool
	IsIgnoreDotFile  bool
	IsIgnoreEmptyDir bool
	IsSerial         bool
	IsDebug          bool
	//
	RegExt    string
	MinAge    string
	MaxAge    string
	MaxSizeMB int64
	MinSizeMB int64
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zstdzip [zip | unzip | hash] [options]",
	Short: "(de)compress file(s) in zip format with ZSTD algorithm, is able to keep last-modified-time & permissions",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		timeDuration := time.Since(timeBoot)
		fmt.Println("\n *** elapse:", timeDuration, "***")

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
	timeBoot = time.Now()
	numCPU := runtime.NumCPU()

	rootCmd.PersistentFlags().StringVar(&Source, "source", "", "source file or folder")
	rootCmd.PersistentFlags().StringVar(&Target, "target", "", "target file, saved.zstd.zip")
	rootCmd.PersistentFlags().IntVar(&Threads, "threads", numCPU, "threads")
	rootCmd.PersistentFlags().StringVar(&Password, "password", "", "set your password")
	rootCmd.PersistentFlags().BoolVar(&IsSerial, "serial", false, "optimization for hard disk, not for ssd")
	rootCmd.PersistentFlags().BoolVar(&IsDryRun, "dry-run", false, "just show result, will not write files")
	rootCmd.PersistentFlags().BoolVarP(&IsDebug, "debug", "", false, "print debug info")

}
