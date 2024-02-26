package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/klauspost/compress/zstd"
)

var (
	Speed  int
	numCPU int
)

type Entity struct {
	InputFullPath      string
	OutputFullPath     string
	ZipLevel           zstd.EncoderLevel
	ZipFileMap         map[string]string
	ZipTempFile        string
	ZipTempFileHandler *os.File
	ZipWriter          *zip.Writer
	UnzipReader        *zip.Reader
}

func NewEntity(inpath string, outpath string) *Entity {
	ett := &Entity{}
	ett.InputFullPath = inpath
	ett.OutputFullPath = outpath

	ett.ZipTempFile = strings.Join([]string{ett.OutputFullPath, "ing"}, ".")
	return ett
}

func (ett *Entity) OpenZipTempFile() *Entity {
	compr := zstd.ZipCompressor(
		zstd.WithWindowSize(1<<20),
		zstd.WithEncoderConcurrency(numCPU),
		zstd.WithEncoderLevel(ett.ZipLevel),
		zstd.WithEncoderCRC(false))

	var err error
	ett.ZipTempFileHandler, err = os.Create(ett.ZipTempFile)
	if err != nil {
		ett.ZipTempFileHandler.Close()
		FatalError(err)
	}

	ett.ZipWriter = zip.NewWriter(ett.ZipTempFileHandler)
	ett.ZipWriter.RegisterCompressor(zstd.ZipMethodWinZip, compr)
	return ett
}

func (ett *Entity) CloseZipTempFile() {
	ett.ZipWriter.Close()
	ett.ZipTempFileHandler.Close()

	err := os.Rename(ett.ZipTempFile, ett.OutputFullPath)
	FatalError(err)
}

func (ett *Entity) WithZipLevel(n int) *Entity {
	cLevel := zstd.SpeedDefault
	switch n {
	case 0:
		cLevel = zstd.SpeedFastest
	case 1:
		cLevel = zstd.SpeedDefault
	case 6:
		cLevel = zstd.SpeedBetterCompression
	case 9:
		cLevel = zstd.SpeedBestCompression
	default:
		cLevel = zstd.SpeedDefault
	}

	ett.ZipLevel = cLevel
	return ett
}

func (ett *Entity) SetZipFileMap() *Entity {
	finfo, err := os.Stat(ett.InputFullPath)
	FatalError(err)

	var zfp map[string]string
	if !finfo.IsDir() {
		zfp = make(map[string]string, 1)
		zfp[ett.InputFullPath] = filepath.Base(ett.InputFullPath)
	} else {
		fpathSkip := filepath.ToSlash(filepath.Dir(ett.InputFullPath))
		nameInZip := ""
		zfp = make(map[string]string, 4096)
		var walkFunc = func(p string, info os.FileInfo, err error) error {
			p = AbsToSlash(p)

			nameInZip = strings.Trim(strings.Replace(p, fpathSkip, "", 1), "/")
			nameInZip = filepath.ToSlash(nameInZip)
			if nameInZip != "" && nameInZip != "." && nameInZip != ".." {
				zfp[p] = nameInZip
			}
			return nil
		}
		err = filepath.Walk(ett.InputFullPath, walkFunc)
		FatalError(err)
	}

	ett.ZipFileMap = zfp
	return ett
}

func (ett *Entity) Compress() *Entity {
	var header *zip.FileHeader
	for abspath, zipname := range ett.ZipFileMap {

		fmt.Println(zipname)

		finfo, err := os.Stat(abspath)
		FatalError(err)

		fp, err := os.Open(abspath)
		FatalError(err)

		header, err = zip.FileInfoHeader(finfo)
		if err != nil {
			FatalError(err)
		} else {
			header.Name = zipname
		}

		header.Method = zstd.ZipMethodWinZip

		w, err := ett.ZipWriter.CreateHeader(header)
		FatalError(err)

		if !finfo.IsDir() {
			_, err = io.Copy(w, fp)
			FatalError(err)
		}

	}

	return ett
}

func (ett *Entity) Decompress() *Entity {
	fh, err := os.Open(ett.InputFullPath)
	FatalError(err)

	finfo, _ := fh.Stat()

	ett.UnzipReader, err = zip.NewReader(fh, finfo.Size())
	FatalError(err)

	decomp := zstd.ZipDecompressor(
		zstd.WithDecoderConcurrency(numCPU),
	)

	ett.UnzipReader.RegisterDecompressor(zstd.ZipMethodWinZip, decomp)

	var dstPath, dstDir string

	for _, fzip := range ett.UnzipReader.File {
		dstPath = filepath.Join(ett.OutputFullPath, fzip.Name)
		dstPath = filepath.ToSlash(dstPath)
		dstDir = filepath.Dir(dstPath)
		if _, err := os.Stat(dstDir); err != nil {
			os.MkdirAll(dstDir, os.ModePerm)
		}

		header := fzip.FileHeader
		if header.FileInfo().IsDir() {
			os.MkdirAll(dstPath, header.Mode())
			continue
		}

		dst, _ := os.Create(dstPath)
		funzip, err := fzip.Open()
		PrintlnError(err)

		if _, err := io.Copy(dst, funzip); err != nil {
			PrintlnError(err)
		} else {
			fmt.Println(fzip.Name)
		}
		if err := funzip.Close(); err != nil {
			PrintlnError(err)
		}
		dst.Close()

		os.Chmod(dstPath, header.FileInfo().Mode())
		os.Chtimes(dstPath, header.FileInfo().ModTime(), header.FileInfo().ModTime())

	}
	return ett
}

func (ett *Entity) DecompressAsync() *Entity {
	fh, err := os.Open(ett.InputFullPath)
	FatalError(err)

	finfo, _ := fh.Stat()

	ett.UnzipReader, err = zip.NewReader(fh, finfo.Size())
	FatalError(err)

	decomp := zstd.ZipDecompressor(
		zstd.WithDecoderConcurrency(numCPU),
	)

	ett.UnzipReader.RegisterDecompressor(zstd.ZipMethodWinZip, decomp)

	var dstPath, dstDir string
	var wg sync.WaitGroup

	for _, fzip := range ett.UnzipReader.File {
		dstPath = filepath.Join(ett.OutputFullPath, fzip.Name)
		dstPath = filepath.ToSlash(dstPath)
		dstDir = filepath.Dir(dstPath)
		if _, err := os.Stat(dstDir); err != nil {
			os.MkdirAll(dstDir, os.ModePerm)
		}

		header := fzip.FileHeader
		if header.FileInfo().IsDir() {
			os.MkdirAll(dstPath, header.Mode())
			continue
		}

		wg.Add(1)

		go func(dstPath string, fzip *zip.File, header zip.FileHeader) {
			defer wg.Done()

			dst, _ := os.Create(dstPath)
			funzip, err := fzip.Open()
			PrintlnError(err)

			if _, err := io.Copy(dst, funzip); err != nil {
				PrintlnError(err)
			} else {
				fmt.Println(fzip.Name)
			}
			if err := funzip.Close(); err != nil {
				PrintlnError(err)
			}
			dst.Close()

			os.Chmod(dstPath, header.FileInfo().Mode())
			os.Chtimes(dstPath, header.FileInfo().ModTime(), header.FileInfo().ModTime())

		}(dstPath, fzip, header)

		wg.Wait()
	}
	return ett
}
