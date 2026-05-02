package cmd

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/zeebo/blake3"
	"github.com/zeebo/xxh3"
)

func argsValidate() error {
	if IsDebug {
		fmt.Printf("--source = %v\n", Source)
		fmt.Printf("--target = %v\n", Target)
	}

	return nil
}

func positionalArgs(pArgs []string) {
	if Source == "" {
		if len(pArgs) == 1 || len(pArgs) == 2 {
			Source = pArgs[0]
		}
	}

	if Target == "" {
		if len(pArgs) == 2 {
			Target = pArgs[1]
		}
	}
}

func ToUnixSlash(s string) string {
	// for windows
	return strings.ReplaceAll(s, "\\", "/")
}

func GetNowUnix() int64 {
	return time.Now().UTC().Unix()
}

func FileExists(fpath string) bool {
	_, err := os.Stat(fpath)
	if err != nil {
		return false
	}
	return true
}

func Int2Str(n int) string {
	return fmt.Sprintf("%d", n)
}

func isFileMatched(fpath string, finfo os.FileInfo, fextreg *regexp.Regexp) bool {
	if fpath == "" || fpath == "." || fpath == ".." {
		return false
	}

	if RegExt != "" {
		if fextreg.MatchString(filepath.Ext(fpath)) == false {
			return false
		}
	}

	if MaxSizeMB != -1 && (finfo.Size() > (MaxSizeMB << 20)) {
		return false
	}

	if MinSizeMB != -1 && (finfo.Size() < (MinSizeMB << 20)) {
		return false
	}

	if MinAge != "" {
		if finfo.ModTime().Unix() < TimeStr2Unix(MinAge) {
			return false
		}
	}

	if MaxAge != "" {
		if finfo.ModTime().Unix() > TimeStr2Unix(MaxAge) {
			return false
		}
	}

	if IsIgnoreDotFile {
		if strings.HasPrefix(filepath.Base(fpath), ".") {
			return false
		}
	}

	return true
}

func hashFile(m string, fpath string) string {
	var hasher hash.Hash
	switch m {
	case "xxhash":
		hasher = xxh3.New()
	case "blake3":
		hasher = blake3.New()
	case "md5":
		hasher = md5.New()
	case "sha256":
		hasher = sha256.New()
	case "sha1":
		hasher = sha1.New()
	default:
		hasher = xxh3.New()
	}

	fh, err := os.Open(fpath)
	if err != nil {
		fh.Close()
		PrintError("HashFile", err)
	}

	r := bufio.NewReader(fh)

	var buf []byte = make([]byte, 8192)
	for {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			PrintError("HashFile", err)
		}
		hasher.Write(buf[:n])
	}

	fh.Close()
	return hex.EncodeToString(hasher.Sum(nil))
}

func TimeStr2Unix(s string) int64 {
	layout := "2006-01-02,15:04:05"
	var parsedTime time.Time
	var err error

	parsedTime, err = time.ParseInLocation(layout, s, time.Local)

	if err != nil {
		PrintError("TimeStr2Unix", NewError(err.Error()+". Correct Format: "+layout))
		return 0
	}

	return parsedTime.Unix()
}

func getCompressLevel(n int) zstd.EncoderLevel {
	cLevel := zstd.SpeedDefault
	switch n {
	case 0:
		cLevel = zstd.SpeedFastest
	case 1:
		cLevel = zstd.SpeedDefault
	case 2:
		cLevel = zstd.SpeedBetterCompression
	case 3:
		cLevel = zstd.SpeedBestCompression
	default:
		cLevel = zstd.SpeedDefault
	}

	return cLevel
}

func MakeDirs(dpath string) error {
	dpath = ToUnixSlash(dpath)
	_, err := os.Stat(dpath)
	if err != nil {
		err = os.MkdirAll(dpath, os.ModePerm)
		PrintError("MakeDirs:MkdirAll", err)
		return err
	}
	return nil
}
