package common

import (
	"log"
	"os"
)

//DebugLog Log开关
type DebugLog bool

var debug DebugLog

var dLog = log.New(os.Stdin, "[Debug]", log.Ltime)

//Printf 格式化输出
func (d DebugLog) Printf(format string, args ...interface{}) {
	if d {
		dLog.Printf(format, args...)
	}
}

//Println 输出
func (d DebugLog) Println(args ...interface{}) {
	if d {
		dLog.Println(args...)
	}
}

//SetDebug 调试开关
func SetDebug(d DebugLog) {
	debug = d
}
