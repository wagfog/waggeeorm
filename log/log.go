package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	//创建两个日志记录器（logger）：errorLog 和 infoLog，并使用标准输出（os.Stdout）作为日志的目标输出
	//我知道你肯定会很奇怪，这个奇怪的字符串干嘛的 ，这个叫做ANSI转义码，我们可以用这个来修改error的颜色
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m", log.LstdFlags|log.Lshortfile) //使用红色 "[error]" 标识，并在日志消息前面添加时间戳和文件名
	infoLog  = log.New(os.Stdout, "\033[34m[info]\033[0m", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}

	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
	//如果level高于ErrorLevel，则将错误日志输出到标准输出；如果level高于InfoLevel，则将信息日志输出到标准输出。
	//对于低于对应级别的日志，将其输出重定向到丢弃文件。这样可以动态地控制日志输出级别，方便在不同情况下进行日志记录
}
