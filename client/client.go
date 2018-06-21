package main

import (
	"encoding/binary"
	"errors"
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

var (
	errAddrType      = errors.New("socks addr type not supported")
	errVer           = errors.New("socks version not supported")
	errMethod        = errors.New("socks only support 1 method now")
	errAuthExtraData = errors.New("socks authentication get extra data")
	errReqExtraData  = errors.New("socks request get extra data")
	errCmd           = errors.New("socks command not supported")
)

const (
	socksVer5       = 5
	socksCmdConnect = 1
)

var debug comm.DebugLog

//ServerCipher 服务端数据结构
type ServerCipher struct {
	server string
	cipher *comm.Cipher
}

var servers struct {
	srv *ServerCipher
}

//handshake: sockets 握手阶段
func handshake(conn net.Conn) (err error) {
	debug.Println("start handshake...")

	//Identifier: the client will connect to the server
	buf := make([]byte, 258)
	conn.SetReadDeadline(time.Now().Add(time.Minute * 2))
	var n int
	if n, err = io.ReadAtLeast(conn, buf, 2); err != nil {
		return err
	}
	ver := buf[0]
	if ver != socksVer5 {
		return errVer
	}
	nmethod := int(buf[1])
	msglen := nmethod + 2
	if n == msglen { //正常方式完成握手
	} else if n < msglen { //存在用户名和密码
		if _, err = io.ReadFull(conn, buf); err != nil { //TODO: 通过用户名密码握手
			return
		}
	} else {
		return errAuthExtraData
	}
	_, err = conn.Write([]byte{5, 0})
	if debug {
		log.Println("finished handshake...")
	}
	return
}

//getRequest: 建立连接
func getRequest(conn net.Conn) (rawaddr []byte, host string, err error) {
	const (
		idVer   = 0
		idCmd   = 1
		idType  = 3
		idIP0   = 4
		idDmLen = 4
		idDm0   = 5

		typeIPv4 = 1
		typeDm   = 3
		typeIPv6 = 4

		lenIPv4   = 3 + 1 + net.IPv4len + 2 // 3(ver+cmd+rsv) + 1addrType + ipv4 + 2port
		lenIPv6   = 3 + 1 + net.IPv6len + 2 // 3(ver+cmd+rsv) + 1addrType + ipv6 + 2port
		lenDmBase = 3 + 1 + 1 + 2           // 3 + 1addrType + 1addrLen + 2port, plus addrLen
	)

	debug.Println("start get request...")

	conn.SetReadDeadline(time.Now().Add(time.Minute * 2))
	buf := make([]byte, 263)
	var n int
	if n, err = io.ReadAtLeast(conn, buf, 5); err != nil { // VER+CMD+RSV+ATYP=4
		return
	}
	if buf[idVer] != socksVer5 {
		err = errVer
		return
	}
	if buf[idCmd] != socksCmdConnect {
		err = errCmd
		return
	}

	reqLen := -1
	switch buf[idType] {
	case typeIPv4:
		reqLen = lenIPv4
	case typeIPv6:
		reqLen = lenIPv6
	case typeDm:
		reqLen = int(buf[idDmLen]) + lenDmBase
	default:
		err = errAddrType
		return
	}

	if n == reqLen {
		//common case, do nothing
	} else if n < reqLen { // rare case
		if _, err = io.ReadFull(conn, buf[n:reqLen]); err != nil {
			return
		}
	} else {
		err = errReqExtraData
		return
	}
	rawaddr = buf[idType:reqLen]
	if debug {
		switch buf[idType] {
		case typeIPv4:
			host = net.IP(buf[idIP0 : idIP0+net.IPv4len]).String()
		case typeDm:
			host = net.IP(buf[idDm0 : idDm0+buf[idDmLen]]).String()
		case typeIPv6:
			host = net.IP(buf[idIP0 : idIP0+net.IPv6len]).String()
		}
		port := binary.BigEndian.Uint16(buf[reqLen-2 : reqLen])
		host = net.JoinHostPort(host, strconv.Itoa(int(port)))
		log.Println("host:", host)
		log.Println("finished get request...")
	}
	return
}

//createServerConn: 连接远程服务器
func createServerConn(rawaddr []byte, addr string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", servers.srv.server)
	if err != nil {
		return
	}
	if _, err = conn.Write(rawaddr); err != nil {
		return nil, err
	}
	return
}

func handleConnection(conn net.Conn) {
	debug.Printf("socks connect to the server %s\n", conn.RemoteAddr().String())

	if err := handshake(conn); err != nil {
		log.Printf("handshake: %s", err)
		return
	}
	rawaddr, addr, err := getRequest(conn)
	if err != nil {
		log.Printf("error get request: %s\n", err)
		return
	}
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	if err != nil {
		log.Println("send connection confirmation:", err)
		return
	}
	remote, err := createServerConn(rawaddr, addr)
	if err != nil {
		return
	}
	go comm.PipeThenClose(conn, remote)
	comm.PipeThenClose(remote, conn)
	debug.Println("closed connection to", addr)
}

func run(addr string) {
	l, err := net.Listen("tcp", addr)
	checkError("listening: ", err)

	debug.Printf("start listening socks5 at %v...\n", addr)

	for {
		conn, err := l.Accept()
		checkError("accept: ", err)
		go handleConnection(conn)
	}
}

func main() {
	flag.BoolVar((*bool)(&debug), "d", false, "print debug message")
	flag.Parse()

	comm.SetDebug(debug)
	config, err := comm.ParseConfig("../config.json")
	if err != nil {
		log.Println(err)
		return
	}
	remote := config.Server + ":" + strconv.Itoa(config.Port)
	servers.srv = &ServerCipher{remote, nil}
	run(config.LocalServer + ":" + strconv.Itoa(config.LocalPort))
}

func checkError(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
		os.Exit(1)
	}
}
