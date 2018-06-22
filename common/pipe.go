package common

import (
	"log"
	"net"
	"time"
)

//PipeThenClose 端口数据转发
func PipeThenClose(src, dst net.Conn) {
	defer dst.Close()
	for {
		src.SetReadDeadline(time.Now().Add(ReadTimeout))
		buf := make([]byte, 4096)
		n, err := src.Read(buf)
		if n > 0 {
			if _, err := dst.Write(buf[:n]); err != nil {
				log.Printf("write error: %s", err)
				break
			}
		}
		if err != nil {
			break
		}
	}
	return
}
