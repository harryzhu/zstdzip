package cmd

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"io"

	"encoding/hex"

	"os"
)

const AESCHUNKSIZE int64 = 16 << 20
const AESBLOCKSIZE int = 16

// PwdKey length can be 16

var (
	pwdKey []byte
	ivKey  []byte
)

type CryptFile struct {
	SrcPath  string
	DstPath  string
	Password string
}

func NewCryptFile(srcpath string, dstpath string, pswd string) *CryptFile {
	if pswd == "" {
		FatalError("NewCryptFile", NewError("password cannot be empty"))
	}

	cf := &CryptFile{
		SrcPath:  srcpath,
		DstPath:  dstpath,
		Password: pswd,
	}

	cf.setKeyPasswordIV()

	return cf
}

func (cf *CryptFile) AESEncode() {
	aesEncodeFile(cf.SrcPath, cf.DstPath)
}

func (cf *CryptFile) AESDecode() {
	aesDecodeFile(cf.SrcPath, cf.DstPath)
}

// ------------

func (cf *CryptFile) setKeyPasswordIV() *CryptFile {
	var salt string = SHA256("Cu5t0m-s@lt")

	if cf.Password == "" {
		FatalError("setKeyPasswordIV", NewError("you did not set any password"))
	}

	pk := SHA256(MD5(cf.Password) + ":" + salt)
	ivk := SHA256(MD5(pk) + ":" + salt)

	pwdKey = []byte(pk)[:AESBLOCKSIZE]
	ivKey = []byte(ivk)[:AESBLOCKSIZE]

	if pwdKey == nil || ivKey == nil {
		FatalError("setKeyPasswordIV", NewError("password and iv key cannot be empty"))
	}

	return cf
}

func aesEncodeFile(src string, dst string) {
	fsrc, err := os.Open(src)
	FatalError("aesEncodeFile", err)

	defer fsrc.Close()
	dst_temp := dst + ".temp"
	fdst, fhdst := NewBufWriter(dst_temp)

	iv := []byte(ivKey)

	block, err := aes.NewCipher(pwdKey)
	FatalError("aesEncodeFile", err)

	stream := cipher.NewCTR(block, iv)

	srcReader := bufio.NewReader(fsrc)
	buf := make([]byte, AESCHUNKSIZE)

	for {
		n, err := srcReader.Read(buf)
		if n == 0 {
			if err == io.EOF {
				break
			}

			if err != nil {
				PrintError("aesEncodeFile", err)
				break
			}
		}
		encByte := make([]byte, n)
		stream.XORKeyStream(encByte, buf[:n])

		_, err = fdst.Write(encByte)
		FatalError("aesEncodeFile", err)
	}
	fdst.Flush()
	fhdst.Close()

	os.Rename(dst_temp, dst)
}

func aesDecodeFile(src string, dst string) {
	fsrc, err := os.Open(src)
	FatalError("aesDecodeFile", err)
	defer fsrc.Close()

	dst_temp := dst + ".temp"
	fdst, fhdst := NewBufWriter(dst_temp)

	iv := []byte(ivKey)

	block, err := aes.NewCipher(pwdKey)
	FatalError("aesDecodeFile", err)

	stream := cipher.NewCTR(block, iv)

	srcReader := bufio.NewReader(fsrc)
	buf := make([]byte, AESCHUNKSIZE)

	for {
		n, err := srcReader.Read(buf)

		if n == 0 {
			if err == io.EOF {
				break
			}

			if err != nil {
				PrintError("aesDecodeFile", err)
				break
			}
		}
		decByte := make([]byte, n)
		stream.XORKeyStream(decByte, buf[:n])

		_, err = fdst.Write(decByte)
		FatalError("aesDecodeFile", err)
	}
	fdst.Flush()
	fhdst.Close()

	os.Rename(dst_temp, dst)
}

func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func NewBufWriter(f string) (*bufio.Writer, *os.File) {
	fh, err := os.Create(f)
	if err != nil {
		fh.Close()
		FatalError("NewBufWriter", err)
	}

	return bufio.NewWriter(fh), fh
}
