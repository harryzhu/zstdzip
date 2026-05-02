package cmd

import (
	"os"
)

type Item struct {
	SrcPath  string
	DstPath  string
	SrcFinfo os.FileInfo
	CopyMode int
	ChanFlag int
	// 4k align
	z struct{}
}

func NewItem(srcPath, dstPath string, srcFinfo os.FileInfo) Item {
	item := Item{}
	item.SrcPath = srcPath
	item.DstPath = dstPath
	item.SrcFinfo = srcFinfo
	item.CopyMode = 0
	item.ChanFlag = 0
	return item
}

func (m Item) WithChanFlag(num int) Item {
	m.ChanFlag = num
	return m
}

func (m Item) WithCopyMode(num int) Item {
	m.CopyMode = num
	return m
}
