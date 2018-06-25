package common

import (
	"log"
	"net"
	"time"
)

//PipeThenClose 端口数据转发
func PipeThenClose(src, dst net.Conn, readDec, writeEnc bool) {
	defer dst.Close()
	for {
		src.SetReadDeadline(time.Now().Add(ReadTimeout))
		buf := make([]byte, 4096)
		n, err := src.Read(buf)
		if err != nil {
			break
		}
		if readDec {
			buf, err = DecryptAESCFB(buf)
			if err != nil {
				break
			}
		}
		if writeEnc {
			buf, err = EncryptAESCFB(buf)
			if err != nil {
				break
			}
		}
		if n > 0 {
			if _, err := dst.Write(buf[:n]); err != nil {
				log.Printf("write error: %s", err)
				break
			}
		}

	}
	return
}
