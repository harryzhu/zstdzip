package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	fmt.Println(" ")
}

func GetTimeNowUnix() int64 {
	return time.Now().Unix()
}

func GetTimeNow() time.Time {
	return time.Now()
}

func FatalError(err error) {
	if err != nil {
		GlobalStatus = "error"
		log.Fatal(err)
	}
}

func PrintlnError(err error) {
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
		fmt.Printf("... %-32s \r", s[0:31])
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
