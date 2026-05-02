package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// hashCmd represents the hash command
var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "--sum=blake3 | xxhash | md5 | sha1 | sha256 --source=path/of/file.txt",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		positionalArgs(args)
		//
		argsValidate()
		bootstrap()

		timeBoot = time.Now()
	},
	Run: func(cmd *cobra.Command, args []string) {
		var hashResult []map[string]string = make([]map[string]string, 0)
		var result []byte
		Sum = strings.ToLower(Sum)

		finfo, err := os.Stat(Source)
		if err != nil {
			FatalError("hash", err)
		}
		if finfo.IsDir() {
			filepath.WalkDir(Source, func(fpath string, dirInfo fs.DirEntry, err error) error {
				if err != nil {
					PrintError("hash: filepath.WalkDir", err)
					return err
				}
				if dirInfo.IsDir() {
					return nil
				}

				fpath = ToUnixSlash(fpath)
				var m map[string]string = make(map[string]string, 0)
				fkey := strings.TrimPrefix(strings.TrimPrefix(fpath, Source), "/")
				m[fkey] = hashFile(Sum, fpath)
				hashResult = append(hashResult, m)

				return nil
			})
			result, err = json.Marshal(hashResult)
			if len(hashResult) > 10 {
				result4print, err := json.Marshal(hashResult[0:10])
				PrintError("hash: json.Marshal", err)
				fmt.Printf("%s ...\n", result4print)
			} else {
				fmt.Printf("%s\n", result)
			}

		} else {
			var m map[string]string = make(map[string]string, 1)
			m[Source] = hashFile(Sum, Source)

			hashResult = hashResult[0:0]
			hashResult = append(hashResult, m)

			result, err = json.Marshal(hashResult)
			FatalError("hash: json.Marshal", err)
			fmt.Printf("%s: %s\n", Sum, m[Source])

		}

		if Target != "" {
			MakeDirs((filepath.Dir(Target)))
			fjson := strings.Join([]string{Target, Sum, "json"}, ".")
			fp, err := os.Create(fjson)
			if err != nil {
				PrintError("hash: os.Create", err)
			}
			_, err = fp.Write(result)
			if err != nil {
				PrintError("hash: WriteString", err)
			} else {
				PrintlnInfo("json saved", fjson)
			}

			fp.Close()
			//
			ftxt := strings.Join([]string{Target, Sum, "txt"}, ".")
			fp, err = os.Create(ftxt)
			if err != nil {
				PrintError("hash: os.Create", err)
			}
			for _, kv := range hashResult {
				for k, v := range kv {
					_, err = fp.WriteString(fmt.Sprintf("%s: %v\n", v, k))
				}

			}
			fp.Close()
			PrintlnInfo("txt saved", ftxt)
		}
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
	rootCmd.MarkFlagRequired("source")

	hashCmd.Flags().StringVar(&Sum, "sum", "sha256", "sum method: md5, sha1, sha256, blake3, xxhash")
}
