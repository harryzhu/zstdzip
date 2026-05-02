package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

type CryptFile struct {
	SrcPath   string
	DstPath   string
	Password  string
	passKey   []byte
	blockSize int
}

func NewCryptFile(srcpath string, dstpath string, pswd string) *CryptFile {
	if pswd == "" {
		PrintError("NewCryptFile", NewError("password cannot be empty"))
	}

	cf := &CryptFile{
		SrcPath:   srcpath,
		DstPath:   dstpath,
		Password:  pswd,
		blockSize: 32,
	}

	cf.setPassKey()

	return cf
}

func (cf *CryptFile) AESEncode(method string) {
	DebugInfo("AESEncode", cf.blockSize, ":", string(cf.passKey))
	if method == "ctr" {
		ctrEncryptFile(cf.SrcPath, cf.DstPath, cf.passKey)
	}
	if method == "gcm" {
		gcmEncryptFile(cf.SrcPath, cf.DstPath, cf.passKey)
	}
}

func (cf *CryptFile) AESDecode(method string) {
	if method == "ctr" {
		ctrDecryptFile(cf.SrcPath, cf.DstPath, cf.passKey)
	}
	if method == "gcm" {
		gcmDecryptFile(cf.SrcPath, cf.DstPath, cf.passKey)
	}
}

func (cf *CryptFile) WithBlockSize(size int) *CryptFile {
	switch {
	case size == 32:
		cf.blockSize = 32
	case size == 24:
		cf.blockSize = 24
	default:
		cf.blockSize = 16
	}

	return cf
}

// ------------

func (cf *CryptFile) setPassKey() *CryptFile {
	var salt string = _sha256("Cu5t0m-s@lt")

	if cf.Password == "" {
		PrintError("setPassKey", NewError("you did not set any password"))
	}

	pk := _sha256(_md5(cf.Password) + ":" + salt)

	//AESBLOCKSIZE: 16(AES-128)/24(AES-192)/32(AES-256) 字节
	var aesBlockSize int = cf.blockSize
	cf.passKey = []byte(pk)[:aesBlockSize]

	if cf.passKey == nil {
		PrintError("setPassKey", NewError("passKey key cannot be empty"))
	}

	return cf
}

func _md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func _sha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// -----

func ctrEncryptFile(src string, dst string, pwdKey []byte) error {
	fsrc, err := os.Open(src)
	if err != nil {
		PrintError("ctrEncryptFile", err)
		return err
	}

	block, err := aes.NewCipher(pwdKey)
	if err != nil {
		PrintError("ctrEncryptFile", err)
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		PrintError("ctrEncryptFile", err)
		return err
	}

	dstTemp := dst + ".ing"
	fdst, err := os.Create(dstTemp)
	if err != nil {
		PrintError("ctrEncryptFile", err)
		return err
	}

	if _, err := fdst.Write(iv); err != nil {
		PrintError("ctrEncryptFile", err)
		return err
	}

	stream := cipher.NewCTR(block, iv)

	AesChunkSize := 128 << 10
	buf := make([]byte, AesChunkSize)
	dstWriter := &cipher.StreamWriter{S: stream, W: fdst}

	for {
		n, err := fsrc.Read(buf)
		if n > 0 {
			if _, err := dstWriter.Write(buf[:n]); err != nil {
				PrintError("ctrEncryptFile", err)
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	//
	fdst.Close()
	fsrc.Close()

	err = os.Rename(dstTemp, dst)
	return err
}

func ctrDecryptFile(src string, dst string, pwdKey []byte) error {
	fsrc, err := os.Open(src)
	if err != nil {
		PrintError("ctrDecryptFile", err)
		return err
	}
	defer fsrc.Close()

	dstTemp := dst + ".ing"
	fdst, err := os.Create(dstTemp)
	if err != nil {
		PrintError("ctrDecryptFile", err)
		return err
	}

	block, err := aes.NewCipher(pwdKey)
	if err != nil {
		PrintError("ctrDecryptFile", err)
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(fsrc, iv); err != nil {
		PrintError("ctrDecryptFile", err)
		return err
	}

	stream := cipher.NewCTR(block, iv)

	AesChunkSize := 128 << 10
	buf := make([]byte, AesChunkSize)
	srcReader := &cipher.StreamReader{S: stream, R: fsrc}

	for {
		n, err := srcReader.Read(buf)
		if n > 0 {
			if _, err := fdst.Write(buf[:n]); err != nil {
				PrintError("ctrDecryptFile", err)
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	fdst.Close()
	fsrc.Close()

	err = os.Rename(dstTemp, dst)
	return err
}

// EncryptFile AES-GCM 加密文件（高性能流式）
// srcPath: 源文件路径
// dstPath: 加密后文件路径
// key: 32 字节(AES-256)/24 字节(AES-192)/16 字节(AES-128)
func gcmEncryptFile(srcPath, dstPath string, pwdKey []byte) error {
	// 打开源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		PrintError("gcmEncryptFile", err)
		return err
	}

	// 创建目标文件
	dstTemp := dstPath + ".ing"
	dstFile, err := os.Create(dstTemp)
	if err != nil {
		PrintError("gcmEncryptFile", err)
		return err
	}

	// 初始化 AES 密码块
	block, err := aes.NewCipher(pwdKey)
	if err != nil {
		PrintError("gcmEncryptFile", err)
		return err
	}

	// 初始化 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		PrintError("gcmEncryptFile", err)
		return err
	}

	// 生成随机 nonce（GCM 标准：12 字节，最高性能）
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		PrintError("gcmEncryptFile", err)
		return err
	}

	// 写入 nonce（解密时必须先读取）
	if _, err := dstFile.Write(nonce); err != nil {
		PrintError("gcmEncryptFile", err)
		return err
	}

	// 预分配缓冲区（零分配核心）
	AesChunkSize := 128 << 10
	buf := make([]byte, AesChunkSize)
	cipherBuf := make([]byte, AesChunkSize+gcm.Overhead())

	// 流式加密
	for {
		// 读取文件块
		n, err := srcFile.Read(buf)
		if n > 0 {
			// GCM 加密（无内存拷贝）
			cipherBuf = gcm.Seal(cipherBuf[:0], nonce, buf[:n], nil)
			// 写入加密数据
			if _, err := dstFile.Write(cipherBuf); err != nil {
				PrintError("gcmEncryptFile", err)
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	dstFile.Close()
	srcFile.Close()

	err = os.Rename(dstTemp, dstPath)

	return err
}

// DecryptFile AES-GCM 解密文件（高性能流式）
// srcPath: 加密文件路径
// dstPath: 解密后文件路径
// key: 与加密时相同的密钥
func gcmDecryptFile(srcPath, dstPath string, pwdKey []byte) error {
	// 打开加密文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		PrintError("gcmDecryptFile", err)
		return err
	}

	// 创建解密文件
	dstTemp := dstPath + ".ing"
	dstFile, err := os.Create(dstTemp)
	if err != nil {
		PrintError("gcmDecryptFile", err)
		return err
	}

	// 初始化 AES 密码块
	block, err := aes.NewCipher(pwdKey)
	if err != nil {
		PrintError("gcmDecryptFile", err)
		return err
	}

	// 初始化 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		PrintError("gcmDecryptFile", err)
		return err
	}

	// 读取 nonce（加密时写入的前 12 字节）
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(srcFile, nonce); err != nil {
		PrintError("gcmDecryptFile", err)
		return err
	}

	// 预分配缓冲区
	AesChunkSize := 128 << 10
	buf := make([]byte, AesChunkSize+gcm.Overhead())
	plainBuf := make([]byte, AesChunkSize)

	// 流式解密
	for {
		// 读取加密块
		n, err := srcFile.Read(buf)
		if n > 0 {
			// GCM 解密（无内存拷贝，自带完整性校验）
			plainBuf, err := gcm.Open(plainBuf[:0], nonce, buf[:n], nil)
			if err != nil {
				PrintError("gcmDecryptFile", err)
				return err
			}
			// 写入解密数据
			if _, err := dstFile.Write(plainBuf); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	dstFile.Close()
	srcFile.Close()

	err = os.Rename(dstTemp, dstPath)
	return nil
}
