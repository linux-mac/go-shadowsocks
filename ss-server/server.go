package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	comm "github.com/go-shadowsocks/common"
)

var debug comm.DebugLog

func getRequest(conn net.Conn) (host string, err error) {
	const (
		idType  = 0 // address type index
		idIP0   = 1 // ip address start index
		idDmLen = 1 // domain address length index
		idDm0   = 2 // domain address start index

		typeIPv4 = 1 // type is ipv4 address
		typeDm   = 3 // type is domain address
		typeIPv6 = 4 // type is ipv6 address

		lenIPv4     = net.IPv4len + 2 // ipv4 + 2port
		lenIPv6     = net.IPv6len + 2 // ipv6 + 2port
		lenDmBase   = 2               // 1addrLen + 2port, plus addrLen
		lenHmacSha1 = 10
	)
	buf := make([]byte, 269)
	if _, err := io.ReadFull(conn, buf[:idType+1]); err != nil {
		return "", err
	}

	var reqStart, reqEnd int
	addrType := buf[idType]
	switch addrType {
	case typeIPv4:
		reqStart, reqEnd = idIP0, idIP0+lenIPv4
	case typeIPv6:
		reqStart, reqEnd = idIP0, idIP0+lenIPv6
	case typeDm:
		if _, err = io.ReadFull(conn, buf[idType+1:idDmLen+1]); err != nil {
			return
		}
		reqStart, reqEnd = idDm0, idDm0+int(buf[idDmLen])+lenDmBase
	default:
		err = fmt.Errorf("addr type %d not supported", addrType)
		return
	}
	if _, err = io.ReadFull(conn, buf[reqStart:reqEnd]); err != nil {
		return
	}

	switch addrType {
	case typeIPv4:
		host = net.IP(buf[idIP0 : idIP0+net.IPv4len]).String()
	case typeIPv6:
		host = net.IP(buf[idIP0 : idIP0+net.IPv6len]).String()
	case typeDm:
		host = string(buf[idDm0 : idDm0+int(buf[idDmLen])])
	}
	port := binary.BigEndian.Uint16(buf[reqEnd-2 : reqEnd])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))
	return
}

func handleClient(conn net.Conn) {
	debug.Println("start handle client...")
	conn.SetReadDeadline(time.Now().Add(comm.ReadTimeout))
	defer conn.Close()
	host, err := getRequest(conn)
	if err != nil {
		debug.Printf("error get request: %s", err)
		return
	}
	remote, err := net.Dial("tcp", host)
	if err != nil {
		debug.Printf("dial error: %s", err)
		return
	}
	go comm.PipeThenClose(conn, remote, true, false)
	comm.PipeThenClose(remote, conn, false, true)
}

func run(port string) {
	debug.Println("start listen port:", port)
	ln, err := net.Listen("tcp", port)
	if err != nil {
		debug.Printf("error listening port %v: %v\n", port, err)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			debug.Printf("Accept error: %s", err)
			return
		}
		debug.Println("start accept...")
		go handleClient(conn)
	}
}

func main() {
	var configPath string
	flag.BoolVar((*bool)(&debug), "d", false, "调试环境")
	flag.StringVar(&configPath, "c", "~/.shadowsocks/config.json", "配置路径")
	flag.Parse()
	debug.Println(configPath)
	comm.SetDebug(debug)
	config, err := comm.ParseConfig(configPath)
	if err != nil {
		log.Println(err)
		return
	}
	comm.InitCipher(config)
	run(":" + strconv.Itoa(config.Port))
}
