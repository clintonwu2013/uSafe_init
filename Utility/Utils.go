package utility

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/robfig/cron"
)

var (
	// DebugLogRoot : root of debug log
	DebugLogRoot = "./Log/"
	// LogDebug : use print immediatilly
	LogDebug *log.Logger
	logFile  *os.File

	fileOpenSuccess = false
	mEnable         = false

	schedule *cron.Cron
)

func getDate() string {
	return time.Now().Format("2006-01-02")
}
func alterLog() {
	var err error
	tFile := logFile
	if mEnable {
		if _, err := os.Stat(DebugLogRoot); os.IsNotExist(err) {
			os.Mkdir(DebugLogRoot, os.ModePerm)
		}

		fileName := DebugLogRoot + "debugLog_" + getDate() + ".log"
		logFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("===== debug log file open failed, %s =====", fileName)
			fmt.Printf("===== error: %s =====", err.Error())
			fmt.Printf("===== can not output debug log to file =====")
			fileOpenSuccess = false
		} else {
			fileOpenSuccess = true
		}
	}
	tFile.Close()

	if mEnable {
		if fileOpenSuccess {
			LogDebug = log.New(io.MultiWriter(os.Stdout, logFile), "[Debug]", log.Ldate|log.Ltime)
		} else {
			LogDebug = log.New(os.Stdout, "[Debug]", log.Ldate|log.Ltime)
		}
	} else {
		LogDebug = log.New(ioutil.Discard, "[Debug]", log.Ldate|log.Ltime)
	}
}
func init() {
	var err error

	EnableDebug(true) // default enable debug

	schedule = cron.New()
	err = schedule.AddFunc(`0 0 0 * * *`, alterLog)
	if err != nil {
		return
	}
	schedule.Start()
}

// EnableDebug : enable debug and write to file(also print in console)
func EnableDebug(enable bool) {
	mEnable = enable
	alterLog()
}

// IfEnableDebug : get if debug enable
func IfEnableDebug() bool {
	return mEnable
}

// DebuggerPrintf : print log with goroutine id
func DebuggerPrintf(format string, v ...interface{}) {
	nv := make([]interface{}, 0, 1+len(v))
	nv = append(append(nv, runtime.GetGoroutineId()), v...)
	go LogDebug.Printf("[%d] "+format, nv...)
}
