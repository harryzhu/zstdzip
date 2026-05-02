package cmd

import (
	"io/fs"
	"regexp"
	"runtime"
)

var (
	chanFile    chan Item
	copyAllDone Item
	DirInfoList map[string]fs.FileInfo
	fextReg     *regexp.Regexp
)

func bootstrap() {
	numCPU = runtime.NumCPU()
	Source = ToUnixSlash(Source)
	Target = ToUnixSlash(Target)
	//
	copyAllDone = NewItem("", "", nil).WithChanFlag(0)
	copyAllDone.CopyMode = -1
	//
	chanFile = make(chan Item, 256)
	//
	DirInfoList = make(map[string]fs.FileInfo, 2048)

	fextReg = regexp.MustCompile("(?i)" + RegExt)
}
