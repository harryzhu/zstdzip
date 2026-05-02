package cmd

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
)

func taskSendFileToChan() error {
	var nameInZip string
	SourceInfo, err := os.Stat(Source)
	if err != nil {
		PrintError("taskSendFileToChan: os.Stat", err)
		return err
	}
	if SourceInfo.IsDir() == false {
		fpath := ToUnixSlash(Source)
		nameInZip := filepath.Base(fpath)
		ele := NewItem(fpath, nameInZip, SourceInfo)
		chanFile <- ele
	} else {
		filepath.Walk(Source, func(fpath string, finfo fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fpath = ToUnixSlash(fpath)
			nameInZip = strings.TrimPrefix(fpath, Source)
			nameInZip = strings.TrimPrefix(nameInZip, "/")

			if finfo.IsDir() {
				return nil
			}

			if nameInZip == "" || nameInZip == "." || nameInZip == ".." {
				return nil
			}

			if isFileMatched(fpath, finfo, fextReg) == false {
				return nil
			}

			ele := NewItem(fpath, nameInZip, finfo)
			chanFile <- ele

			return nil
		})
	}

	chanFile <- NewItem("", "", nil).WithChanFlag(0).WithCopyMode(-1)
	if !IsSerial {
		chanFile <- NewItem("", "", nil).WithChanFlag(1).WithCopyMode(-1)
		chanFile <- NewItem("", "", nil).WithChanFlag(2).WithCopyMode(-1)
		chanFile <- NewItem("", "", nil).WithChanFlag(3).WithCopyMode(-1)
		chanFile <- NewItem("", "", nil).WithChanFlag(4).WithCopyMode(-1)
		chanFile <- NewItem("", "", nil).WithChanFlag(5).WithCopyMode(-1)
		chanFile <- NewItem("", "", nil).WithChanFlag(6).WithCopyMode(-1)
		chanFile <- NewItem("", "", nil).WithChanFlag(7).WithCopyMode(-1)
	}

	return nil
}

func taskGetChanFileToOneArchive(num int) error {
	compressor := zstd.ZipCompressor(
		zstd.WithWindowSize(1<<20),
		zstd.WithEncoderConcurrency(numCPU),
		zstd.WithEncoderLevel(getCompressLevel(Level)),
		zstd.WithEncoderCRC(false))

	targetTemp := strings.Join([]string{Target, Int2Str(num), "ing"}, ".")

	DebugInfo("taskGetChanFileToOneArchive", targetTemp)
	targetTempFileHandler, err := os.Create(targetTemp)
	if err != nil {
		PrintError("taskGetChanFileToDisk", err)
		return err
	}

	targetTempWriter := zip.NewWriter(targetTempFileHandler)
	targetTempWriter.RegisterCompressor(zstd.ZipMethodWinZip, compressor)

	//-----------
	for {
		item := <-chanFile
		if item.CopyMode == -1 {
			break
		}
		//DebugInfo("taskGetChanFileToDisk", item.DstPath)

		srcPath := item.SrcPath
		nameInZip := item.DstPath
		srcFinfo := item.SrcFinfo

		if srcPath == "" || nameInZip == "" || srcFinfo == nil {
			continue
		}
		//
		header, err := zip.FileInfoHeader(srcFinfo)
		if err != nil {
			PrintError("taskGetChanFileToDisk: zip.FileInfoHeader", err)
			continue
		}

		header.Name = nameInZip
		header.Method = zstd.ZipMethodWinZip

		if !srcFinfo.IsDir() {
			fp, err := os.Open(srcPath)
			if err != nil {
				PrintError("taskGetChanFileToDisk: os.Open", err)
				continue
			}

			w, err := targetTempWriter.CreateHeader(header)
			if err != nil {
				PrintError("taskGetChanFileToDisk: CreateHeader", err)
				continue
			}

			_, err = io.Copy(w, fp)
			if err != nil {
				PrintError("taskGetChanFileToDisk: io.Copy", err)
				continue
			}

			fp.Close()
		}

	}
	//
	targetTempWriter.Close()
	targetTempFileHandler.Close()

	if Password != "" {
		NewCryptFile(targetTemp, targetTemp+".encrypted", Password).AESEncode(AESMethod)
		err = os.Rename(targetTemp+".encrypted", targetTemp)
		PrintError("taskGetChanFileToDisk:NewCryptFile", err)
	}
	DebugInfo("taskGetChanFileToDisk", targetTemp)
	err = os.Rename(targetTemp, strings.TrimSuffix(strings.TrimSuffix(targetTemp, ".ing"), ".0"))
	PrintError("taskGetChanFileToDisk:os.Rename", err)

	return nil
}

func taskDecompressFile(sourceFile string) error {
	srcPath := ToUnixSlash(sourceFile)
	srcFinfo, err := os.Stat(srcPath)
	if err != nil {
		PrintError("taskDecompressFile:os.Stat", err)
		return err
	}

	if srcFinfo.IsDir() {
		PrintError("taskDecompressFile:os.Stat", NewError("--source= should be a file, not a folder"))
		return nil
	}

	var srcHandler *os.File
	if Password != "" {
		NewCryptFile(srcPath, srcPath+".decrypted", Password).AESDecode(AESMethod)
		srcHandler, err = os.Open(srcPath + ".decrypted")
	} else {
		srcHandler, err = os.Open(srcPath)
	}

	unzipReader, err := zip.NewReader(srcHandler, srcFinfo.Size())
	if err != nil {
		PrintError("taskDecompressFile:zip.NewReader", err)
		return err
	}

	decompressor := zstd.ZipDecompressor(
		zstd.WithDecoderConcurrency(numCPU),
	)

	unzipReader.RegisterDecompressor(zstd.ZipMethodWinZip, decompressor)

	var dstPath, dstDir string
	for _, fzip := range unzipReader.File {
		finfo := fzip.FileHeader.FileInfo()

		dstPath = ToUnixSlash(filepath.Join(Target, fzip.Name))
		dstDir = ToUnixSlash(filepath.Dir(dstPath))

		if finfo.IsDir() {
			MakeDirs(dstPath)
			continue
		} else {
			MakeDirs(dstDir)
		}

		if isFileMatched(dstPath, finfo, fextReg) == false {
			continue
		}

		dstFile, err := os.Create(dstPath)
		PrintError("taskDecompressFile:os.Create", err)

		funzip, err := fzip.Open()
		PrintError("taskDecompressFile:fzip.Open", err)

		if _, err := io.Copy(dstFile, funzip); err != nil {
			PrintError("taskDecompressFile:io.Copy", err)
		}

		if err := funzip.Close(); err != nil {
			PrintError("taskDecompressFile:funzip.Close", err)
		}
		dstFile.Close()

		err = os.Chtimes(dstPath, finfo.ModTime(), finfo.ModTime())
		PrintError("taskDecompressFile:os.Chtimes", err)

		err = os.Chmod(dstPath, finfo.Mode())
		PrintError("taskDecompressFile:os.Chmod", err)

	}

	if FileExists(srcPath + ".decrypted") {
		err = os.Remove(srcPath + ".decrypted")
		PrintError("taskDecompressFile:os.Remove", err)
	}

	return nil
}
