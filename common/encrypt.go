package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
)

type cipherInfo struct {
	key []byte
	iv  []byte
}

var c cipherInfo

//InitCipher 初始化
func InitCipher(config *Config) {
	pwd := sha256.Sum256([]byte(config.Password))
	c.iv = pwd[:aes.BlockSize]
	switch config.Method {
	case "aes-128-cfb":
		c.key = pwd[:16]
		break
	case "aes-192-cfb":
		c.key = pwd[:24]
		break
	case "aes-256-cfb":
		c.key = pwd[:32]
		break
	}
}

//EncryptAESCFB 加密
func EncryptAESCFB(src []byte) (dst []byte, err error) {
	dst = make([]byte, len(src))
	blocker, err := aes.NewCipher(c.key)
	if err != nil {
		return
	}
	aesEncryter := cipher.NewCFBEncrypter(blocker, c.iv)
	aesEncryter.XORKeyStream(dst, src)
	return
}

//DecryptAESCFB 解密
func DecryptAESCFB(src []byte) (dst []byte, err error) {
	dst = make([]byte, len(src))
	blocker, err := aes.NewCipher(c.key)
	if err != nil {
		return
	}
	aesDecrypter := cipher.NewCFBDecrypter(blocker, c.iv)
	aesDecrypter.XORKeyStream(dst, src)
	return
}
