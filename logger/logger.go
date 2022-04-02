package logger

import (
	"fmt"
	"strings"

	. "github.com/aywfelix/felixgo/logger/handler"
)

const (
	LOG_LEVEL      = LOG_DEBUG
	OUTPUT_CONSOLE = true
	OUTPUT_FILE    = true
)

//日志头部信息标记位，采用bitmap方式，用户可以选择头部需要哪些标记位被打印
const (
	BitDateTime  = 1 << iota
	BitLongFile                                                       // 完整文件名称 /home/go/src/zinx/server.go
	BitShortFile                                                      // 最后文件名   server.go
	BitFunc                                                           // 打印所在函数名字
	BitFileLine                                                       // 文件行号
	BitLevel                                                          // 当前日志级别： 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
	BitStdFlag   = BitDateTime                                        // 标准头部日志格式
	BitDefault   = BitStdFlag | BitLevel | BitShortFile | BitFileLine // 默认日志头部格式
)

const (
	LOG_DEBUG = iota
	LOG_INFO
	LOG_WARN
	LOG_ERROR
	LOG_FATAL
)

var LogLevel = []string{
	"[D]",
	"[I]",
	"[W]",
	"[E]",
	"[F]",
}

type Logger struct {
	flag     int
	builder  strings.Builder
	logQueue chan *LogRecord
	handler  IHandler
}

var (
	DefaultLogger = newDefaultLogger("./gamelog/gamelog")
	LogDebug      = DefaultLogger.Debug
	LogInfo       = DefaultLogger.Info
	LogWarn       = DefaultLogger.Warn
	LogError      = DefaultLogger.Error
	LogFatal      = DefaultLogger.Fatal
)

func NewLogger(baseName string) *Logger {
	logger := new(Logger)
	logger.logQueue = make(chan *LogRecord, 1000)
	logger.handler = NewTimeRotatingFileHandler(baseName, WhenHour)
	if logger.handler == nil {
		panic("logger handler is nil")
	}
	logger.flag |= BitDateTime | BitLevel | BitShortFile | BitFileLine
	return logger
}

func newDefaultLogger(logPath string) *Logger {
	logger := new(Logger)
	logger.logQueue = make(chan *LogRecord, 1000)
	logger.handler = NewTimeRotatingFileHandler(logPath, WhenHour)
	if logger.handler == nil {
		panic("logger handler is nil")
	}
	logger.flag |= BitDateTime | BitLevel | BitShortFile | BitFileLine
	go logger.DefaultLogRun()
	return logger
}

func (l *Logger) SetFlag(flag int) {
	l.flag = flag
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.writeLog(LOG_DEBUG, fmt.Sprintf(format, args...))
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.writeLog(LOG_INFO, fmt.Sprintf(format, args...))
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.writeLog(LOG_WARN, fmt.Sprintf(format, args...))
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.writeLog(LOG_ERROR, fmt.Sprintf(format, args...))
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.writeLog(LOG_FATAL, fmt.Sprintf(format, args...))
}
