package logger

import (
	"context"
	"fmt"
	ct "github.com/daviddengcn/go-colortext"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LogRecord struct {
	level int
	text  string
}

var (
	logRecordPool = &sync.Pool{
		New: func() interface{} {
			return new(LogRecord)
		},
	}
	consoleColor = map[int]ct.Color{
		LOG_DEBUG: ct.White,
		LOG_INFO:  ct.Green,
		LOG_WARN:  ct.Yellow,
		LOG_ERROR: ct.Red,
		LOG_FATAL: ct.Red,
	}
	// runtime.Caller(2)
	CallDepth = 2
)

const limitLine = 100

func (l *Logger) formatHeader(file string, line int, level int, funcName string) {
	now := time.Now()
	if l.flag&(BitDateTime) != 0 {
		fmt.Fprintf(&l.builder, "[%s]", now.Format("2006-01-02 15:04:05"))
	}
	if l.flag&BitLevel != 0 {
		l.builder.WriteString(LogLevel[level])
	}
	if l.flag&(BitShortFile|BitLongFile|BitFileLine) != 0 {
		if l.flag&BitShortFile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}

		if l.flag&BitFileLine != 0 {
			fmt.Fprintf(&l.builder, "[%s:%d]", file, line)
		} else {
			fmt.Fprintf(&l.builder, "[%s]", file)
		}
	}
	if l.flag&BitFunc != 0 {
		fmt.Fprintf(&l.builder, "[%s]", funcName)
	}
}

func (l *Logger) writeLog(level int, text string) {
	if level < LOG_LEVEL {
		return
	}

	fileName := "unknow-file"
	funcName := "unknow"
	line := 0
	var ok bool
	var pc uintptr
	if l == nil {
		fmt.Println("logger is nil")
		return
	}
	if l.flag&(BitLongFile|BitShortFile) != 0 {
		pc, fileName, line, ok = runtime.Caller(CallDepth)
		if ok {
			fileName = path.Base(fileName)
			funcName = runtime.FuncForPC(pc).Name()
		}
	}

	l.formatHeader(fileName, line, level, funcName)

	l.builder.WriteString(text)
	if len(text) > 0 && text[len(text)-1] != '\n' {
		l.builder.WriteByte('\n')
	}

	text = l.builder.String()
	l.builder.Reset()
	// 从池中获取临时buffer
	record := logRecordPool.Get().(*LogRecord)
	record.level = level
	record.text = text
	l.logQueue <- record
}

func (l *Logger) saveLog(length int) {
	var builder strings.Builder
	builder.Grow(length * 256)

	for i := 0; i < length; i++ {
		logRecord := <-l.logQueue
		level := logRecord.level
		if OUTPUT_CONSOLE {
			// cmd 打印
			ct.ChangeColor(consoleColor[level], true, ct.Black, false)
			fmt.Printf(logRecord.text)
			ct.ResetColor()
		}
		builder.WriteString(logRecord.text)
		// 重新放回池中
		logRecordPool.Put(logRecord)
	}

	if OUTPUT_FILE {
		l.handler.Write([]byte(builder.String()))
	}
}

func (l *Logger) LogClose() {
	length := len(l.logQueue)
	if length >= 0 {
		l.saveLog(length)
	}
	l.handler.Close()
	close(l.logQueue)
}

func (l *Logger) LogStart(args ...interface{}) {
	ctx := args[0].(context.Context)
	wg := args[1].(*sync.WaitGroup)
	defer wg.Done()
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ctx.Done():
			l.LogClose()
			return
		case <-ticker.C:
			for length := len(l.logQueue); length > 0; {
				if length >= limitLine {
					l.saveLog(limitLine)
					length -= limitLine
				} else {
					l.saveLog(length)
					length = 0
				}
			}
		}
	}
}

func (l *Logger) DefaultLogRun() {
	ticker := time.NewTicker(time.Millisecond * 1)
	for {
		select {
		case <-ticker.C:
			for length := len(l.logQueue); length > 0; {
				if length >= limitLine {
					l.saveLog(limitLine)
					length -= limitLine
				} else {
					l.saveLog(length)
					length = 0
				}
			}
		}
	}
}
