package common

import (
	"fmt"
	"io"
	"net"
)

//Conn 自定义Conn结构
type Conn struct {
	net.Conn
	*Cipher
	ReadBuf  []byte
	WriteBuf []byte
}

//NewConn 创建Conn
func NewConn(c net.Conn, cipher *Cipher) *Conn {
	return &Conn{
		Conn:     c,
		Cipher:   cipher,
		ReadBuf:  leakyBuf.Get(),
		WriteBuf: leakyBuf.Get(),
	}
}

//Close 关闭Conn
func (c *Conn) Close() error {
	leakyBuf.Put(c.ReadBuf)
	leakyBuf.Put(c.WriteBuf)
	return c.Conn.Close()
}

//DialWithRawAddr 封装后的Conn发起远程请求
func DialWithRawAddr(rawaddr []byte, server string, cipher *Cipher) (c *Conn, err error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return
	}
	c = NewConn(conn, cipher)
	if _, err = c.Write(rawaddr); err != nil {
		c.Close()
		return nil, err
	}
	return
}

func (c *Conn) Read(b []byte) (n int, err error) {
	debug.Println("start read...")
	if c.dec == nil {
		iv := make([]byte, c.info.ivLen)
		if _, err = io.ReadFull(c.Conn, iv); err != nil {
			return
		}
		if err = c.initDecrypt(iv); err != nil {
			return
		}
		if len(c.iv) == 0 {
			c.iv = iv
		}
	}

	cipherData := c.ReadBuf
	if len(b) > len(cipherData) {
		cipherData = make([]byte, len(b))
	} else {
		cipherData = cipherData[:len(b)]
	}
	n, err = c.Conn.Read(cipherData)
	if n > 0 {
		c.decrypt(b[0:n], cipherData[0:n])
		fmt.Println(b)
	}
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	debug.Println("start write...")
	fmt.Println(b)
	var iv []byte
	if c.enc == nil {
		iv, err = c.initEncrypt()
		if err != nil {
			return
		}
	}

	cipherData := c.WriteBuf
	dataSize := len(b) + len(iv)
	if dataSize > len(cipherData) {
		cipherData = make([]byte, dataSize)
	} else {
		cipherData = cipherData[:dataSize]
	}

	if iv != nil {
		copy(cipherData, iv)
	}

	c.encrypt(cipherData[len(iv):], b)
	// cipherData := make([]byte, len(b))
	// c.encrypt(cipherData, b)
	n, err = c.Conn.Write(cipherData)
	return
}
