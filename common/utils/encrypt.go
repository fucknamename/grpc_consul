package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// MD5加密
func Md5Encrypt(str string) string {
	md := md5.New()
	md.Write([]byte(str))                // 需要加密的字符串为 str
	cipherStr := md.Sum(nil)             // 不需要拼接额外的数据，如果约定了额外加密数据，可以在这里传递
	return hex.EncodeToString(cipherStr) // 输出加密结果
}

// HMACSHA256String 计算 HMAC-SHA256 签名
func HMACSHA256String(message, key string) string {
	signature := HMACSHA256(message, key)
	// 转换为16进制字符串
	return hex.EncodeToString(signature)
}

// HMACSHA256 计算 HMAC-SHA256 签名
func HMACSHA256(message, key string) []byte {
	// 将 key 和 message 转换为字节切片
	keyBytes := []byte(key)
	messageBytes := []byte(message)

	// 创建一个新的 HMAC 使用 SHA256
	h := hmac.New(sha256.New, keyBytes)
	h.Write(messageBytes)

	// 获取最终的哈希结果
	signature := h.Sum(nil)
	return signature
}

// AES-128-CBC 加密
func AesEncryptCBC2(plainText, key []byte) (string, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", nil, err
	}

	blockSize := block.BlockSize()
	plainText = pkcs7Padding(plainText, blockSize)

	// 生成随机 IV
	// iv := make([]byte, blockSize)
	// if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	// 	return "", nil, err
	// }
	iv := key[:16]

	mode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, plainText)

	// 返回 base64 编码的密文和 IV
	// return base64.RawURLEncoding.EncodeToString(cipherText), iv, nil
	return base64.StdEncoding.EncodeToString(cipherText), iv, nil
}

// AES-128-CBC 解密
func AesDecryptCBC2(cipherTextBase64 string, key, iv []byte) (string, error) {
	// cipherText, err := base64.RawURLEncoding.DecodeString(cipherTextBase64)
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	mode.CryptBlocks(plainText, cipherText)

	plainText = pkcs7Unpadding(plainText)
	return string(plainText), nil
}

// PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS7 去填充
func pkcs7Unpadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
