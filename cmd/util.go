package cmd

import (
	"bufio"
	"crypto/sha1"
	"hash"

	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zeebo/blake3"
	"github.com/zeebo/xxh3"
)

var (
	GlobalStatus string = "ok"
)

func PrintArgs() {
	fmt.Println("--input=", Input)
	fmt.Println("--output=", Output)
	fmt.Println("--speed=", Speed)
	fmt.Println("--threads=", numCPU)
	fmt.Println("--async=", Async)
	fmt.Println("--debug=", IsDebug)
	fmt.Println("--sum=", Sum)
	fmt.Println(" ")
}

func DivideFloat64(a, b float64) float64 {
	if b == 0.0 {
		panic("a cannot be divided by ZERO")
	}

	return a / b
}

func PrintSpeed(fsize float64, tsec float64) {
	MB := float64(1 << 20)
	fspeed := DivideFloat64(fsize, tsec)
	fspeed_mb := DivideFloat64(fspeed, MB)
	fmt.Printf("size: %.2f Bytes( %.2f MB ), seconds: %.7f\n", fsize, DivideFloat64(fsize, MB), tsec)
	fmt.Printf("Speed: %.2f MB/s\n", fspeed_mb)

}

func GetTimeNowUnix() int64 {
	return time.Now().Unix()
}

func GetTimeNow() time.Time {
	return time.Now()
}

func FatalError(err any) {
	if err != nil {
		GlobalStatus = "error"
		log.Fatal(err)
	}
}

func PrintlnError(err any) {
	if err != nil {
		GlobalStatus = "error"
		log.Println(err)
	}
}

func PrintlnDebug(s string) {
	if IsDebug {
		fmt.Println(s)
	}
}

func PrintSpinner(s string) {
	if IsDebug {
		fmt.Printf("... %5.30s\r", s)
	}
}

func AbsToSlash(s string) string {
	s, err := filepath.Abs(s)
	FatalError(err)
	return strings.TrimRight(filepath.ToSlash(s), "/")
}

func SaveJson(p string, m map[string]string) {
	fp, err := os.Create(p)
	FatalError(err)

	j, err := json.Marshal(m)
	FatalError(err)

	_, err = fp.Write(j)
	defer fp.Close()

	FatalError(err)
}

func HashFile(m string) string {
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

	fh, err := os.Open(Input)
	if err != nil {
		fh.Close()
		FatalError(err)
	}

	if IsDebug {
		fhinfo, _ := fh.Stat()
		TotalSize += fhinfo.Size()
	}

	r := bufio.NewReader(fh)

	var buf []byte = make([]byte, 8192)
	for {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			FatalError(err)
		}
		hasher.Write(buf[:n])
	}

	fh.Close()
	return hex.EncodeToString(hasher.Sum(nil))
}
