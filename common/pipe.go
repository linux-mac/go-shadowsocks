package common

import (
	"log"
	"net"
	"time"
)

//PipeThenClose 端口数据转发
func PipeThenClose(src, dst net.Conn) {
	defer dst.Close()
	log.Println("start pipe...")
	for {
		src.SetReadDeadline(time.Now().Add(600 * time.Second))
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
		log.Println("data: ", string(buf[:n-1]))
	}
	log.Println("finished")
	return
}
