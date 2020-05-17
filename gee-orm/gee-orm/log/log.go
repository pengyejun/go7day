package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// log.Lshortfile 显示中文和行号
var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers = []*log.Logger{errorLog, infoLog}
	mu sync.Mutex
)

var (
	Error = errorLog.Println
	Errorf = errorLog.Printf
	Info = infoLog.Println
	Infof = infoLog.Printf
)

const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}

	if InfoLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}
}