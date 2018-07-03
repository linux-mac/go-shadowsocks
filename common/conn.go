package common

import (
	"io"
	"net"
)

//Conn struct
type Conn struct {
	net.Conn
	*Cipher
	ReadBuf  []byte
	WriteBuf []byte
}

//NewConn create
func NewConn(c net.Conn, cipher *Cipher) *Conn {
	return &Conn{
		Conn:     c,
		Cipher:   cipher,
		ReadBuf:  leakyBuf.Get(),
		WriteBuf: leakyBuf.Get(),
	}
}

//Close Conn
func (c *Conn) Close() error {
	leakyBuf.Put(c.ReadBuf)
	leakyBuf.Put(c.WriteBuf)
	return c.Conn.Close()
}

//DialWithRawAddr create remote connection
func DialWithRawAddr(rawaddr []byte, server string, cipher *Cipher) (c *Conn, err error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return
	}
	c = NewConn(conn, cipher.Copy())
	if _, err = c.Write(rawaddr); err != nil {
		c.Close()
		return nil, err
	}
	return
}

func (c *Conn) Read(b []byte) (n int, err error) {
	debug.Println("start read...")
	if c.dec == nil {
		debug.Println("initDecrypt")
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
	debug.Println("before decrypt:", cipherData)
	if n > 0 {
		c.decrypt(b[0:n], cipherData[0:n])
		debug.Println("after decrypt:", b)
	}
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	debug.Println("start write...")
	debug.Println("before encrypt:", b)
	var iv []byte
	if c.enc == nil {
		debug.Println("initEncrypt")
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
	debug.Println("after encrypt:", cipherData[len(iv):])
	n, err = c.Conn.Write(cipherData)
	return
}
