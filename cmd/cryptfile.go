package cmd

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"io"
	"strings"

	"encoding/hex"

	"os"
)

type CryptFile struct {
	SrcPath  string
	DstPath  string
	Password string
}

func NewCryptFile(srcpath string, dstpath string, pswd string) *CryptFile {
	if pswd == "" {
		FatalError("password cannot be empty")
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

const AESCHUNKSIZE int64 = 16 << 20
const AESBLOCKSIZE int = 16

// PwdKey length can be 16

var (
	pwdKey []byte
	ivKey  []byte
)

// ------------

func (cf *CryptFile) setKeyPasswordIV() *CryptFile {
	var salt string = SHA256("Cu5t0m-s@lt")

	if cf.Password == "" {
		FatalError("you did not set any password")
	}

	pk := SHA256(MD5(cf.Password) + ":" + salt)
	ivk := SHA256(MD5(pk) + ":" + salt)

	pwdKey = []byte(pk)[:AESBLOCKSIZE]
	ivKey = []byte(ivk)[:AESBLOCKSIZE]

	if pwdKey == nil || ivKey == nil {
		FatalError("password and iv key cannot be empty")
	}

	return cf
}

func aesEncodeFile(src string, dst string) {
	fsrc, err := os.Open(src)
	FatalError(err)

	defer fsrc.Close()
	dst_temp := dst + ".temp"
	fdst, fhdst := NewBufWriter(dst_temp)

	iv := []byte(ivKey)

	block, err := aes.NewCipher(pwdKey)
	FatalError(err)

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
				PrintlnError(err)
				break
			}
		}
		encByte := make([]byte, n)
		stream.XORKeyStream(encByte, buf[:n])

		_, err = fdst.Write(encByte)
		FatalError(err)
	}
	fdst.Flush()
	fhdst.Close()

	os.Rename(dst_temp, dst)
}

func aesDecodeFile(src string, dst string) {
	fsrc, err := os.Open(src)
	FatalError(err)
	defer fsrc.Close()

	dst_temp := dst + ".temp"
	fdst, fhdst := NewBufWriter(dst_temp)

	iv := []byte(ivKey)

	block, err := aes.NewCipher(pwdKey)
	FatalError(err)

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
				PrintlnError(err)
				break
			}
		}
		decByte := make([]byte, n)
		stream.XORKeyStream(decByte, buf[:n])

		_, err = fdst.Write(decByte)
		FatalError(err)
	}
	fdst.Flush()
	fhdst.Close()

	os.Rename(dst_temp, dst)
}

func GetEnv(s string, vDefault string) string {
	v := os.Getenv(s)
	if v == "" {
		return vDefault
	}
	return strings.Trim(v, " ")
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
		FatalError(err)
	}

	return bufio.NewWriter(fh), fh
}
