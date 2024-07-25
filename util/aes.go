package util

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
)

// PKCS7Padding pads the plaintext to a multiple of the block size
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding unpads the plaintext
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AESEncrypt encrypts plaintext using AES algorithm in ECB mode
func AESEncrypt(plainText, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plainTextBytes := []byte(plainText)
	plainTextBytes = PKCS7Padding(plainTextBytes, block.BlockSize())

	ciphertext := make([]byte, len(plainTextBytes))
	// ECB mode does not use IV, so we simply encrypt blocks directly
	for bs, be := 0, block.BlockSize(); bs < len(plainTextBytes); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Encrypt(ciphertext[bs:be], plainTextBytes[bs:be])
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func GenerateKey(t, p string) string {
	key := ""

	// 从 T 中提取字符
	for _, e := range []int{2, 11, 22, 23, 29, 30, 33, 36} {
		key += charAt(t, e-1)
	}

	// 从 P 中提取字符
	for _, e := range []int{1, 7, 8, 12, 15, 18, 19, 28} {
		key += charAt(p, e-1)
	}

	return key
}

func charAt(s string, index int) string {
	if index < 0 || index >= len(s) {
		return ""
	}
	return string(s[index])
}
