package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"path/filepath"

	"m3u8-downloader/internal/errors"
)

// ============ 文件操作 ============

// PathExists 检查路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// EnsureDir 确保目录存在
func EnsureDir(path string) error {
	exists, err := PathExists(path)
	if err != nil {
		return errors.New(errors.DirCreate, "检查目录失败", err)
	}

	if !exists {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return errors.New(errors.DirCreate, "创建目录失败", err)
		}
	}

	return nil
}

// ReadFile 读取文件内容
func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New(errors.FileRead, fmt.Sprintf("读取文件失败: %s", path), err)
	}
	return data, nil
}

// WriteFile 写入文件内容
func WriteFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0666)
	if err != nil {
		return errors.New(errors.FileWrite, fmt.Sprintf("写入文件失败: %s", path), err)
	}
	return nil
}

// GetFileSize 获取文件大小（单位：MB）
func GetFileSize(path string) (float64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return float64(info.Size()) / (1024 * 1024), nil
}

// RemoveDir 删除目录
func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

// ============ 加解密操作 ============

// PKCS7Padding PKCS7 填充
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding PKCS7 去填充
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return origData
	}
	unpadding := int(origData[length-1])
	if unpadding > length {
		return origData
	}
	return origData[:(length - unpadding)]
}

// AesEncrypt AES 加密
func AesEncrypt(origData, key []byte, ivs ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	var iv []byte
	if len(ivs) == 0 {
		iv = key
	} else {
		iv = ivs[0]
	}
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AesDecrypt AES 解密
func AesDecrypt(crypted, key []byte, ivs ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	var iv []byte
	if len(ivs) == 0 {
		iv = key
	} else {
		iv = ivs[0]
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

// ============ TS 文件处理 ============

const (
	// SyncByte TS 同步字节
	SyncByte = uint8(0x47) // 71 in decimal
)

// RemoveTSPadding 移除 TS 文件前的填充字节
func RemoveTSPadding(data []byte) []byte {
	for i := 0; i < len(data); i++ {
		if data[i] == SyncByte {
			return data[i:]
		}
	}
	return data
}

// ============ URL 处理 ============

// ListTSFiles 列出目录中的所有 TS 文件
func ListTSFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New(errors.FileRead, "读取目录失败", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".ts" {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}
