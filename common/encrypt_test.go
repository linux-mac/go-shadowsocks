package common

import (
	"reflect"
	"testing"
)

func TestEvpBytesToKey(t *testing.T) {
	pwd := "foobar"
	key := evpBytesToKey(pwd, 32)
	keyTarget := []byte{0x38, 0x58, 0xf6, 0x22, 0x30, 0xac, 0x3c, 0x91, 0x5f, 0x30, 0x0c, 0x66, 0x43, 0x12, 0xc6, 0x3f, 0x56, 0x83, 0x78, 0x52, 0x96, 0x14, 0xd2, 0x2d, 0xdb, 0x49, 0x23, 0x7d, 0x2f, 0x60, 0xbf, 0xdf}
	if !reflect.DeepEqual(key, keyTarget) {
		t.Errorf("key not correct\n\texpect: %v\n\tgot:   %v\n", keyTarget, key)
	}
}

const text = "1"

func testCipher(t *testing.T, c *Cipher, msg string) {
	n := len(text)
	cipherBuf := make([]byte, n)
	originTxt := make([]byte, n)

	c.encrypt(cipherBuf, []byte(text))
	c.decrypt(originTxt, cipherBuf)

	if string(originTxt) != text {
		t.Error(msg, "encrypt then decrytp does not get original text")
	}
}

func testBlockCipher(t *testing.T, method string) {
	var cipher *Cipher
	var err error

	err = CheckCipherMethod(method)
	if err != nil {
		t.Error("cipher method not support")
	}
	srv := Server{}
	srv.Method = method
	srv.Password = "foobar"
	cipher = NewCipher(srv)
	iv, err := cipher.initEncrypt()
	if err != nil {
		t.Error(method, "initEncrypt:", err)
	}
	if err := cipher.initDecrypt(iv); err != nil {
		t.Error(method, "initDecrypt:", err)
	}
	testCipher(t, cipher, method)

	cipherCopy := cipher.Copy()
	iv, err = cipherCopy.initEncrypt()
	if err != nil {
		t.Error(method, "copy initEncrypt:", err)
	}
	if err = cipherCopy.initDecrypt(iv); err != nil {
		t.Error(method, "copy initDecrypt:", err)
	}
	testCipher(t, cipherCopy, method+" copy")
}

func TestAES256CFB(t *testing.T) {
	testBlockCipher(t, "aes-256-cfb")
}
