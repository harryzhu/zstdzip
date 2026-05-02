package cmd

import (
	"io/fs"
	"runtime"
)

var (
	chanFile    chan Item
	copyAllDone Item
	DirInfoList map[string]fs.FileInfo
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
}
