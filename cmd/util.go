package cmd

import (
	"archive/zip"
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/klauspost/compress/zstd"
	"github.com/zeebo/blake3"
	"github.com/zeebo/xxh3"
)

func PrintArgs(args ...string) error {
	// if IsDebug == false {
	// 	return nil
	// }
	fmt.Println(Green("--------------------"))
	if Contains(args, "source") {
		fmt.Println("--source=", Source)
	}

	if Contains(args, "target") {
		fmt.Println("--target=", Target)
	}

	if Contains(args, "serial") {
		fmt.Println("--serial=", IsSerial)
	}

	if Contains(args, "sum") {
		fmt.Println("--sum=", Sum)
	}

	if Contains(args, "threads") {
		fmt.Println("--threads=", Threads)
	}

	if Contains(args, "level") {
		fmt.Println("--level=", Level)
	}

	if Contains(args, "min-age") && MinAge != "" {
		fmt.Println("--min-age=", TimeStr2Unix(MinAge))
	}

	if Contains(args, "max-age") && MaxAge != "" {
		fmt.Println("--max-age=", TimeStr2Unix(MaxAge))
	}

	if Contains(args, "min-size-mb") && MinSizeMB != -1 {
		fmt.Println("--min-size-mb=", MinSizeMB)
	}

	if Contains(args, "max-size-mb") && MaxSizeMB != -1 {
		fmt.Println("--max-size-mb=", MaxSizeMB)
	}

	if Contains(args, "ignore-dot-file") {
		fmt.Println("--ignore-dot-file=", IsIgnoreDotFile)
	}

	if Contains(args, "ignore-empty-dir") {
		fmt.Println("--ignore-empty-dir=", IsIgnoreEmptyDir)
	}

	if Contains(args, "ext") {
		fmt.Println("--ext=", RegExt)
	}

	if IsDebug {
		fmt.Println("--debug=", IsDebug)
	}

	fmt.Println(Green("--------------------"))

	return nil
}

func GetNowUnix() int64 {
	return time.Now().UTC().Unix()
}

func ToUnixSlash(s string) string {
	// for windows
	return strings.ReplaceAll(s, "\\", "/")
}

func GetXxhashString(b []byte) string {
	return strconv.FormatUint(xxhash.Sum64(b), 10)
	//return strconv.FormatUint(xxh3.Hash(b), 10)
}

func GetMD5String(b []byte) string {
	hasher := md5.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}

func Int2Str(n int) string {
	return strconv.Itoa(n)
}

func Contains(arr []string, target string) bool {
	for _, val := range arr {
		if val == target {
			return true
		}
	}
	return false
}

func GetCompressLevel(n int) zstd.EncoderLevel {
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

func OpenZipTempFile(zipTempFile string) (zipTempFileHandler *os.File, zipTempWriter *zip.Writer) {
	compr := zstd.ZipCompressor(
		zstd.WithWindowSize(1<<20),
		zstd.WithEncoderConcurrency(Threads),
		zstd.WithEncoderLevel(GetCompressLevel(Level)),
		zstd.WithEncoderCRC(false))

	var err error
	zipTempFileHandler, err = os.Create(zipTempFile)
	if err != nil {
		zipTempFileHandler.Close()
		FatalError("OpenZipTempFile", err)
	}

	zipTempWriter = zip.NewWriter(zipTempFileHandler)
	zipTempWriter.RegisterCompressor(zstd.ZipMethodWinZip, compr)
	return zipTempFileHandler, zipTempWriter
}

func CloseZipTempFile(zipTempFile string, zipTempFileHandler *os.File, zipTempWriter *zip.Writer) {
	zipTempWriter.Close()
	zipTempFileHandler.Close()

	if Password != "" {
		NewCryptFile(zipTempFile, zipTempFile+".encrypted", Password).AESEncode()
		os.Rename(zipTempFile+".encrypted", zipTempFile)
	}

	TargetFile := strings.Replace(zipTempFile, ".ing", "", 1)
	err := os.Rename(zipTempFile, TargetFile)
	FatalError("OpenZipTempFile", err)
}

func CopyFile(src, dst string) error {
	src = ToUnixSlash(src)
	dst = ToUnixSlash(dst)
	srcFileHandler, err := os.Open(src)
	if err != nil {
		PrintError("CopyFile: os.Open", err)
		return err
	}
	defer srcFileHandler.Close()

	dstTemp := dst + ".ing"

	MakeDirs(filepath.Dir(dstTemp))

	dstFileHandler, err := os.Create(dstTemp)
	if err != nil {
		PrintError("CopyFile: os.Create", err)
		return err
	}
	defer dstFileHandler.Close()

	srcReader := bufio.NewReader(srcFileHandler)
	dstWriter := bufio.NewWriter(dstFileHandler)
	_, err = io.Copy(dstWriter, srcReader)
	if err != nil {
		PrintError("CopyFile: io.Copy", err)
		return err
	}

	dstWriter.Flush()

	finfo, err := srcFileHandler.Stat()
	PrintError("CopyFile", err)

	err = os.Chtimes(dstTemp, finfo.ModTime(), finfo.ModTime())
	PrintError("CopyFile: os.Chtimes", err)

	srcFileHandler.Close()
	dstFileHandler.Close()

	err = os.Rename(dstTemp, dst)
	PrintError("CopyFile: os.Rename", err)

	err = os.Chmod(dst, finfo.Mode())
	PrintError("CopyFile: os.Chmod", err)

	return nil
}

func MakeDirs(dpath string) error {
	dpath = ToUnixSlash(dpath)
	_, err := os.Stat(dpath)
	if err != nil {
		DebugInfo("MakeDirs", dpath)
		err = os.MkdirAll(dpath, os.ModePerm)
		PrintError("MakeDirs:MkdirAll", err)
		return err
	}
	return nil
}

func SendFileToChanFile(srcPath string, dstPath string) (ele map[string]string, err error) {
	srcPath = ToUnixSlash(srcPath)
	dstPath = ToUnixSlash(dstPath)

	ele = make(map[string]string)

	ele["srcPath"] = srcPath
	ele["dstPath"] = dstPath

	return ele, nil
}

func GetChanFileToDisk(chanFileNum chan map[string]string, tw *zip.Writer) error {
	for {
		cf := <-chanFileNum
		if val, ok := cf["_COPYSTATUS"]; ok {
			DebugInfo("_COPYSTATUS:", val)
			break
		}

		if srcPath, ok := cf["srcPath"]; ok {
			dstPath := cf["dstPath"]
			if IsDebug {
				fmt.Printf("%s <== %s\n", dstPath, srcPath)
			}

			finfo, err := os.Stat(srcPath)
			if err != nil {
				PrintError("os.Stat:"+srcPath, err)
				return err
			}

			header, err := zip.FileInfoHeader(finfo)
			if err != nil {
				PrintError("zip.FileInfoHeader:"+srcPath, err)
				return err
			} else {
				header.Name = dstPath
			}

			header.Method = zstd.ZipMethodWinZip

			w, err := tw.CreateHeader(header)
			if err != nil {
				PrintError("tw.CreateHeader:"+srcPath, err)
				return err
			}

			if !finfo.IsDir() {
				fp, err := os.Open(srcPath)
				defer fp.Close()

				if err != nil {
					PrintError("os.Open:"+srcPath, err)
					return err
				}
				_, err = io.Copy(w, fp)

				if err != nil {
					PrintError("io.Copy:"+srcPath, err)
					return err
				}
				fp.Close()
			}

		}

	}

	return nil
}

func CompressDir() error {
	var qcap int = 10
	var chanFile chan map[string]string = make(chan map[string]string, qcap)
	var chanFile1 chan map[string]string = make(chan map[string]string, qcap)
	var chanFile2 chan map[string]string = make(chan map[string]string, qcap)
	var chanFile3 chan map[string]string = make(chan map[string]string, qcap)
	var chanFile4 chan map[string]string = make(chan map[string]string, qcap)
	var chanFile5 chan map[string]string = make(chan map[string]string, qcap)
	var chanFile6 chan map[string]string = make(chan map[string]string, qcap)
	var chanFile7 chan map[string]string = make(chan map[string]string, qcap)

	wg := sync.WaitGroup{}

	wg.Add(3)

	go func() error {
		defer wg.Done()

		t0 := Target + ".ing"
		t0FileHandler, t0Writer := OpenZipTempFile(t0)

		GetChanFileToDisk(chanFile, t0Writer)

		CloseZipTempFile(t0, t0FileHandler, t0Writer)

		return nil
	}()

	go func() error {
		defer wg.Done()

		if !IsSerial {
			wgCompress := sync.WaitGroup{}
			var t1, t2, t3, t4, t5, t6, t7 string
			var t1FileHandler, t2FileHandler, t3FileHandler, t4FileHandler, t5FileHandler, t6FileHandler, t7FileHandler *os.File
			var t1Writer, t2Writer, t3Writer, t4Writer, t5Writer, t6Writer, t7Writer *zip.Writer

			t1 = Target + ".1.ing"
			t1FileHandler, t1Writer = OpenZipTempFile(t1)

			t2 = Target + ".2.ing"
			t2FileHandler, t2Writer = OpenZipTempFile(t2)

			t3 = Target + ".3.ing"
			t3FileHandler, t3Writer = OpenZipTempFile(t3)

			t4 = Target + ".4.ing"
			t4FileHandler, t4Writer = OpenZipTempFile(t4)

			t5 = Target + ".5.ing"
			t5FileHandler, t5Writer = OpenZipTempFile(t5)

			t6 = Target + ".6.ing"
			t6FileHandler, t6Writer = OpenZipTempFile(t6)

			t7 = Target + ".7.ing"
			t7FileHandler, t7Writer = OpenZipTempFile(t7)

			wgCompress.Add(7)
			go func() {
				defer wgCompress.Done()
				GetChanFileToDisk(chanFile1, t1Writer)
			}()
			go func() {
				defer wgCompress.Done()
				GetChanFileToDisk(chanFile2, t2Writer)
			}()
			go func() {
				defer wgCompress.Done()
				GetChanFileToDisk(chanFile3, t3Writer)
			}()
			go func() {
				defer wgCompress.Done()
				GetChanFileToDisk(chanFile4, t4Writer)
			}()
			go func() {
				defer wgCompress.Done()
				GetChanFileToDisk(chanFile5, t5Writer)
			}()
			go func() {
				defer wgCompress.Done()
				GetChanFileToDisk(chanFile6, t6Writer)
			}()
			go func() {
				defer wgCompress.Done()
				GetChanFileToDisk(chanFile7, t7Writer)
			}()

			wgCompress.Wait()

			CloseZipTempFile(t1, t1FileHandler, t1Writer)
			CloseZipTempFile(t2, t2FileHandler, t2Writer)
			CloseZipTempFile(t3, t3FileHandler, t3Writer)
			CloseZipTempFile(t4, t4FileHandler, t4Writer)
			CloseZipTempFile(t5, t5FileHandler, t5Writer)
			CloseZipTempFile(t6, t6FileHandler, t6Writer)
			CloseZipTempFile(t7, t7FileHandler, t7Writer)
		}
		return nil
	}()

	go func() error {
		defer wg.Done()

		var nameInZip string

		regExt := regexp.MustCompile("(?i)" + RegExt)

		num := 0
		filepath.Walk(Source, func(fpath string, finfo os.FileInfo, err error) error {
			if err != nil {
				PrintError("CompressFiles: walkdir", err)
				return err
			}

			if IsIgnoreEmptyDir {
				if finfo.IsDir() {
					return nil
				}
			}

			if IsIgnoreDotFile {
				if strings.HasPrefix(filepath.Base(fpath), ".") {
					return nil
				}
			}

			if RegExt != "" {
				if regExt.MatchString(filepath.Ext(fpath)) == false {
					return nil
				}
			}

			if MaxSizeMB != -1 && (finfo.Size() > (MaxSizeMB << 20)) {
				return nil
			}

			if MinSizeMB != -1 && (finfo.Size() < (MinSizeMB << 20)) {
				return nil
			}

			if MinAge != "" {
				if finfo.ModTime().Unix() < TimeStr2Unix(MinAge) {
					return nil
				}
			}

			if MaxAge != "" {
				if finfo.ModTime().Unix() > TimeStr2Unix(MaxAge) {
					return nil
				}
			}

			fpath = ToUnixSlash(fpath)
			nameInZip = strings.Trim(fpath[len(Source):], "/")

			if nameInZip == "" || nameInZip == "." || nameInZip == ".." {
				return nil
			}

			ele, err := SendFileToChanFile(fpath, nameInZip)
			if err != nil {
				return err
			}

			if !IsSerial {
				twhash := strings.ToLower(GetMD5String([]byte(fpath)))
				switch twhash[0:1] {
				case "0":
					chanFile <- ele
				case "1":
					chanFile <- ele
				case "2":
					chanFile1 <- ele
				case "3":
					chanFile1 <- ele
				case "4":
					chanFile2 <- ele
				case "5":
					chanFile2 <- ele
				case "6":
					chanFile3 <- ele
				case "7":
					chanFile3 <- ele
				case "8":
					chanFile4 <- ele
				case "9":
					chanFile4 <- ele
				case "a":
					chanFile5 <- ele
				case "b":
					chanFile5 <- ele
				case "c":
					chanFile6 <- ele
				case "d":
					chanFile6 <- ele
				case "e":
					chanFile7 <- ele
				case "f":
					chanFile7 <- ele
				default:
					chanFile <- ele
					DebugInfo(nameInZip, " ==== ", twhash)
				}
			} else {
				chanFile <- ele
			}

			num++
			if !IsDebug {
				if num < 100 || num%10 == 0 {
					PrintSpinner(Int2Str(num))
				}
			}

			return nil
		})

		PrintSpinner(Int2Str(num))
		atomic.StoreInt32(&DeComTotalNum, int32(num))

		copyDone := make(map[string]string)
		copyDone["_COPYSTATUS"] = "DONE"
		chanFile <- copyDone
		chanFile1 <- copyDone
		chanFile2 <- copyDone
		chanFile3 <- copyDone
		chanFile4 <- copyDone
		chanFile5 <- copyDone
		chanFile6 <- copyDone
		chanFile7 <- copyDone
		return nil
	}()

	wg.Wait()

	return nil
}

func CompressFile(finfo os.FileInfo) error {
	header, err := zip.FileInfoHeader(finfo)
	if err != nil {
		FatalError(Source, err)
		return err
	}
	fpath := ToUnixSlash(Source)
	nameInPath := filepath.Base(fpath)

	header.Name = nameInPath

	header.Method = zstd.ZipMethodWinZip

	t0 := Target + ".ing"
	t0FileHandler, t0Writer := OpenZipTempFile(t0)

	fp, err := os.Open(Source)
	if err != nil {
		FatalError(Source, err)
		return err
	}

	w, err := t0Writer.CreateHeader(header)
	if err != nil {
		FatalError(Source, err)
		return err
	}

	if !finfo.IsDir() {
		_, err = io.Copy(w, fp)
		if err != nil {
			PrintError(Source, err)
			return err
		}

	}

	fp.Close()

	CloseZipTempFile(t0, t0FileHandler, t0Writer)

	return nil
}

func DecompressFile(fpath string) error {
	var fh *os.File
	var err error

	if Password != "" {
		NewCryptFile(fpath, fpath+".decrypted", Password).AESDecode()
		fh, err = os.Open(fpath + ".decrypted")
	} else {
		fh, err = os.Open(fpath)
	}

	FatalError("DecompressFile", err)

	finfo, _ := fh.Stat()

	unzipReader, err := zip.NewReader(fh, finfo.Size())
	FatalError("DecompressFile", err)

	decomp := zstd.ZipDecompressor(
		zstd.WithDecoderConcurrency(Threads),
	)

	unzipReader.RegisterDecompressor(zstd.ZipMethodWinZip, decomp)

	var dstPath, dstDir string

	num := 0
	regExt := regexp.MustCompile("(?i)" + RegExt)

	for _, fzip := range unzipReader.File {
		header := fzip.FileHeader
		finfo := header.FileInfo()

		dstPath = filepath.ToSlash(filepath.Join(Target, fzip.Name))
		dstDir = filepath.Dir(dstPath)

		if finfo.IsDir() {
			DeComDirInfoList = append(DeComDirInfoList, fzip)
			continue
		}

		if RegExt != "" {
			if regExt.MatchString(filepath.Ext(fzip.Name)) == false {
				continue
			}
		}

		if MaxSizeMB != -1 && (finfo.Size() > (MaxSizeMB << 20)) {
			continue
		}

		if MinSizeMB != -1 && (finfo.Size() < (MinSizeMB << 20)) {
			continue
		}

		if MinAge != "" {
			if header.FileInfo().ModTime().Unix() < TimeStr2Unix(MinAge) {
				continue
			}
		}

		if MaxAge != "" {
			if finfo.ModTime().Unix() > TimeStr2Unix(MaxAge) {
				continue
			}
		}

		if _, err := os.Stat(dstDir); err != nil {
			MakeDirs(dstDir)
		}

		atomic.AddInt32(&DeComTotalNum, 1)

		num++
		if !IsDebug {
			if num < 100 || num%10 == 0 {
				PrintSpinner(Int2Str(int(atomic.LoadInt32(&DeComTotalNum))))
			}
		} else {
			fmt.Printf("%s ==> %s\n", fzip.Name, dstPath)
		}

		dst, _ := os.Create(dstPath)
		funzip, err := fzip.Open()
		PrintError("DecompressFile:fzip.Open", err)

		if _, err := io.Copy(dst, funzip); err != nil {
			PrintError("DecompressFile:io.Copy", err)
		}

		if err := funzip.Close(); err != nil {
			PrintError("DecompressFile:funzip.Close", err)
		}
		dst.Close()

		err = os.Chtimes(dstPath, finfo.ModTime(), finfo.ModTime())
		PrintError("DecompressFile:os.Chtimes", err)

		err = os.Chmod(dstPath, finfo.Mode())
		PrintError("DecompressFile:os.Chmod", err)

	}

	if IsDebug {
		PrintSpinner(Int2Str(num))
	}

	fh.Close()

	if Password != "" {
		_, err = os.Stat(fpath + ".decrypted")
		if err == nil {
			os.Remove(fpath + ".decrypted")
		}
		PrintError("DecompressFile", err)
	}

	return nil
}

func DecompressDirMod() error {
	if len(DeComDirInfoList) == 0 {
		return nil
	}

	var modTimeList map[string]time.Time = make(map[string]time.Time, 1)
	var modeList map[string]fs.FileMode = make(map[string]fs.FileMode, 1)

	for _, m := range DeComDirInfoList {
		info := m.FileInfo()
		modTimeList[info.Name()] = m.ModTime()
		modeList[info.Name()] = m.Mode()
	}

	var dstPath string

	filepath.WalkDir(Target, func(fpath string, finfo fs.DirEntry, err error) error {
		if err != nil {
			PrintError("DecompressDirMod", err)
		}

		if finfo.IsDir() {
			dstPath = ToUnixSlash(fpath)
			if modtime, ok := modTimeList[finfo.Name()]; ok {
				DebugInfo("DecompressDirMod", modtime, dstPath)
				err = os.Chtimes(dstPath, modtime, modtime)
				PrintError("DecompressDirMod:os.Chtimes", err)
			}

			if modelist, ok := modeList[finfo.Name()]; ok {
				DebugInfo("DecompressDirMod", modelist, dstPath)
				err = os.Chmod(dstPath, modelist)
				PrintError("DecompressDirMod:os.Chmod", err)
			}

		}

		return nil
	})

	return nil
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

	fh, err := os.Open(Source)
	if err != nil {
		fh.Close()
		FatalError("HashFile", err)
	}

	r := bufio.NewReader(fh)

	var buf []byte = make([]byte, 8192)
	for {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			FatalError("HashFile", err)
		}
		hasher.Write(buf[:n])
	}

	fh.Close()
	return hex.EncodeToString(hasher.Sum(nil))
}

func TimeStr2Unix(s string) int64 {
	layout := "20060102150405"
	var parsedTime time.Time
	var err error

	parsedTime, err = time.ParseInLocation(layout, s, time.Local)

	if err != nil {
		PrintError("TimeStr2Unix", err)
		os.Exit(0)
	}

	return parsedTime.Unix()
}
