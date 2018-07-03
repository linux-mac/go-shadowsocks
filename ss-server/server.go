package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	comm "github.com/go-shadowsocks/common"
)

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

var debug comm.DebugLog

func getRequest(conn *comm.Conn) (host string, err error) {
	debug.Println("get Request...")
	comm.SetReadTimeout(conn)

	buf := make([]byte, 269)
	// io.ReadFull(conn, buf[:7])
	// return "", err
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
		err = errors.New(`😓decode error, cause this happened maybe:
			1. client and server password is different
			2. your server cannot connect to the website you are aiming to visit
			3. if not the above reasons, please email to maintainer : kunnsh@gmail.com`)
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
	debug.Println("finish get Request...")
	return
}

func handleClient(conn *comm.Conn, port string) {
	var host string
	debug.Println("start handle client...")
	debug.Printf("new client %s->%s\n", conn.RemoteAddr().String(), conn.LocalAddr())
	closed := false
	defer func() {
		debug.Printf("closed pipe %s<->%s\n", conn.RemoteAddr().String(), host)
		if !closed {
			conn.Close()
		}
	}()

	host, err := getRequest(conn)
	if err != nil {
		debug.Printf("error get request: %s", err)
		closed = true
		return
	}

	if strings.ContainsRune(host, 0x00) {
		log.Println("invalid domain")
		closed = true
		return
	}

	debug.Println("connecting: ", host)
	//remote, err := net.DialTimeout("tcp", host, time.Second*12)
	remote, err := net.Dial("tcp", host)
	if err != nil {
		debug.Printf("dial error: %s", err)
		closed = true
		return
	}
	defer func() {
		if !closed {
			remote.Close()
		}
	}()
	debug.Printf("piping %s<->%s", conn.RemoteAddr().String(), host)
	go comm.PipeThenClose(conn, remote)
	comm.PipeThenClose(remote, conn)
	closed = true
	return
}

func run(srv comm.Server) {
	debug.Println("start listen port:", srv.Port)
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(srv.Port))
	if err != nil {
		debug.Printf("error listening port %v: %v\n", srv.Port, err)
		os.Exit(1)
	}
	var cipher *comm.Cipher
	for {
		conn, err := ln.Accept()
		if err != nil {
			debug.Printf("Accept error: %s", err)
			return
		}
		//if cipher == nil {
		cipher = comm.NewCipher(srv)
		debug.Println("create cipher for port: ", srv.Port)
		//}
		debug.Println("start accept...")
		go handleClient(comm.NewConn(conn, cipher), strconv.Itoa(srv.Port))
	}
}

func main() {
	var configPath string
	var version bool

	flag.BoolVar((*bool)(&debug), "d", false, "debug mode")
	flag.BoolVar((*bool)(&version), "v", false, "current version")
	flag.StringVar(&configPath, "c", os.Getenv("HOME")+"/.shadowsocks/config.json", "config path")
	flag.Parse()

	if version {
		comm.PrintVersion()
		os.Exit(0)
	}
	comm.SetDebug(debug)
	debug.Println("loading config file: ", configPath)
	config, err := comm.ParseConfig(configPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for index, srv := range config.Servers {
		if srv.Password == "" || srv.Port == 0 {
			log.Printf("index of [%d] config has errors，it won't be run on the server", index+1)
			continue
		}
		if srv.Method == "" {
			srv.Method = "chacha20-ietf-poly1305"
		}

		err := comm.CheckCipherMethod(srv.Method)
		if err != nil {
			log.Println(err)
			continue
		}
		go run(srv)
	}
	waitSignal()
}

func waitSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)
	for sig := range sigChan {
		if sig == syscall.SIGHUP {
			log.Println("Todo: update config")
		} else {
			log.Printf("caught signal %v, exit", sig)
			os.Exit(0)
		}
	}
}
