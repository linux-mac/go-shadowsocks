package common

import (
	"log"
	"os"
)

//DebugLog debug log type
type DebugLog bool

var debug DebugLog

var dLog = log.New(os.Stdin, "[Debug]", log.Ltime)

//Printf log output format
func (d DebugLog) Printf(format string, args ...interface{}) {
	if d {
		dLog.Printf(format, args...)
	}
}

//Println log output with newline
func (d DebugLog) Println(args ...interface{}) {
	if d {
		dLog.Println(args...)
	}
}

//SetDebug debug switch
func SetDebug(d DebugLog) {
	debug = d
}
