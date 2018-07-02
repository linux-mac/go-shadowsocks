package common

import (
	"io"
	"log"
	"net"
	"time"
)

//SetReadTimeout conn read time out
func SetReadTimeout(c net.Conn) {
	if ReadTimeout != 0 {
		c.SetReadDeadline(time.Now().Add(ReadTimeout))
	}
}

//PipeThenClose data transfer
func PipeThenClose(src, dst net.Conn) {
	defer dst.Close()
	buf := leakyBuf.Get()
	defer leakyBuf.Put(buf)
	for {
		SetReadTimeout(src)
		n, err := src.Read(buf)
		if n > 0 {
			if _, err := dst.Write(buf[:n]); err != nil {
				log.Printf("write error: %s", err)
				break
			}
		}
		if err != nil {
			if err == io.EOF {
				log.Println("read EOF")
			}
			break
		}
	}
	return
}
