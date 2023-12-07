package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"
)

func PrintArgs() {
	fmt.Println(" *** start:", timeBoot.Format("15:04:05"))
	fmt.Println("--input=", Input)
	fmt.Println("--output=", Output)
	fmt.Println("--speed=", Speed)
	fmt.Println("threads:", numCPU)
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
		log.Fatal(err)
	}
}

func PrintlnError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func AbsToSlash(s string) string {
	s, err := filepath.Abs(s)
	FatalError(err)
	return strings.TrimRight(filepath.ToSlash(s), "/")
}
