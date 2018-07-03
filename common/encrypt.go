package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"io"

	"github.com/Yawning/chacha20"
)

type cipherInfo struct {
	keyLen    int
	ivLen     int
	newStream func(key, iv []byte, doe DecOrEnc) (cipher.Stream, error)
}

//Cipher 加密结构
type Cipher struct {
	enc  cipher.Stream
	dec  cipher.Stream
	key  []byte
	info *cipherInfo
	iv   []byte
}

var c Cipher

var cipherMethod = map[string]*cipherInfo{
	"aes-256-cfb":            {32, 16, newAESCFBStream},
	"chacha20-ietf-poly1305": {32, 12, newChaCha20IETFStream},
}

//DecOrEnc enum type
type DecOrEnc int

//enum
const (
	Decrypt DecOrEnc = iota
	Encrypt
)

func newStream(block cipher.Block, err error, key, iv []byte,
	doe DecOrEnc) (cipher.Stream, error) {
	if err != nil {
		return nil, err
	}
	if doe == Encrypt {
		return cipher.NewCFBEncrypter(block, iv), nil
	}
	return cipher.NewCFBDecrypter(block, iv), nil
}

func newAESCFBStream(key, iv []byte, doe DecOrEnc) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	return newStream(block, err, key, iv, doe)
}

func newChaCha20IETFStream(key, iv []byte, _ DecOrEnc) (cipher.Stream, error) {
	return chacha20.NewCipher(key, iv)
}

//CheckCipherMethod 检查加密方法是否系统支持
func CheckCipherMethod(method string) error {
	_, ok := cipherMethod[method]
	if !ok {
		return errors.New("不支持的加密算法: " + method)
	}
	return nil
}

//NewCipher 初始化
func NewCipher(srv Server) (c *Cipher) {
	m := cipherMethod[srv.Method]
	key := evpBytesToKey(srv.Password, m.keyLen)
	c = &Cipher{key: key, info: m}
	return c
}

func md5sum(d []byte) []byte {
	h := md5.New()
	h.Write(d)
	return h.Sum(nil)
}

func evpBytesToKey(password string, keyLen int) (key []byte) {
	const md5Len = 16
	cnt := (keyLen-1)/md5Len + 1
	m := make([]byte, cnt*md5Len)
	copy(m, md5sum([]byte(password)))
	d := make([]byte, md5Len+len(password))
	start := 0
	for i := 1; i < cnt; i++ {
		start += md5Len
		copy(d, m[start-md5Len:start])
		copy(d[md5Len:], password)
		copy(m[start:], md5sum(d))
	}
	return m[:keyLen]
}

func (c *Cipher) initEncrypt() (iv []byte, err error) {
	if c.iv == nil {
		iv = make([]byte, c.info.ivLen)
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil, err
		}
		c.iv = iv
	} else {
		iv = c.iv
	}
	c.enc, err = c.info.newStream(c.key, iv, Encrypt)
	return
}

func (c *Cipher) initDecrypt(iv []byte) (err error) {
	c.dec, err = c.info.newStream(c.key, iv, Decrypt)
	return
}

func (c *Cipher) encrypt(dst, src []byte) {
	c.enc.XORKeyStream(dst, src)
}

func (c *Cipher) decrypt(dst, src []byte) {
	c.dec.XORKeyStream(dst, src)
}
