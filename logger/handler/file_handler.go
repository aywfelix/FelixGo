package handler

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	. "github.com/aywfelix/felixgo/utils"
)

type FileHandler struct {
	fd *os.File
}

func NewFileHandler(fileName string, flag int) (*FileHandler, error) {
	dir := path.Dir(fileName)
	os.Mkdir(dir, 0777)

	f, err := os.OpenFile(fileName, flag, 0)
	if err != nil {
		return nil, err
	}

	handler := &FileHandler{fd: f}
	return handler, nil
}

func (h *FileHandler) Write(b []byte) (int, error) {
	return h.fd.Write(b)
}

func (h *FileHandler) Close() error {
	return h.fd.Close()
}

type RotatingFileHandler struct {
	fd *os.File

	fileName    string
	maxBytes    int64
	curBytes    int64
	backupCount int
}

func NewRotatingFileHandler(fileName string, maxBytes int64, backupCount int) (*RotatingFileHandler, error) {
	dir := path.Dir(fileName)
	if !File.Exists(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}

	if maxBytes <= 0 {
		return nil, errors.New("invalid max bytes")
	}

	h := &RotatingFileHandler{
		fileName:    fileName,
		maxBytes:    maxBytes,
		curBytes:    0,
		backupCount: backupCount,
	}

	var err error
	h.fd, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	f, err := h.fd.Stat()
	if err != nil {
		return nil, err
	}
	h.curBytes = f.Size()
	return h, nil
}

func (h *RotatingFileHandler) Write(b []byte) (int, error) {
	h.doRollover()
	n, err := h.fd.Write(b)
	h.curBytes += int64(n)
	return n, err
}

func (h *RotatingFileHandler) Close() error {
	return h.fd.Close()
}

func (h *RotatingFileHandler) doRollover() {
	if h.curBytes < h.maxBytes {
		return
	}

	f, err := h.fd.Stat()
	if err != nil {
		return
	}
	if f.Size() < h.maxBytes {
		h.curBytes = f.Size()
		return
	}

	if h.backupCount > 0 {
		h.fd.Close()

		for i := h.backupCount - 1; i > 0; i-- {
			sfn := fmt.Sprintf("%s.%d", h.fileName, i)
			dfn := fmt.Sprintf("%s.%d", h.fileName, i+1)

			os.Rename(sfn, dfn)
		}
		dfn := fmt.Sprintf("%s.1", h.fileName)
		os.Rename(h.fileName, dfn)

		h.fd, _ = os.OpenFile(h.fileName, os.O_CREATE|os.O_APPEND, 0666)
		h.curBytes = 0
		f, err := h.fd.Stat()
		if err != nil {
			return
		}
		h.curBytes = f.Size()
	}
}

type TimeRotatingFileHandler struct {
	fd *os.File

	baseName   string
	interval   int64
	suffix     string
	rolloverAt int64
}

const (
	WhenSecond = iota
	WhenMinute
	WhenHour
	WhenDay
)

func NewTimeRotatingFileHandler(fileName string, when int8) *TimeRotatingFileHandler {
	dir := path.Dir(fileName)
	if !File.Exists(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}

	h := new(TimeRotatingFileHandler)
	h.baseName = fileName

	switch when {
	case WhenSecond:
		h.interval = 1
		h.suffix = "2006-01-02_15-04-05"
	case WhenMinute:
		h.interval = 60
		h.suffix = "2006-01-02_15-04"
	case WhenHour:
		h.interval = 3600
		h.suffix = "2006-01-02_15"
	case WhenDay:
		h.interval = 3600 * 24
		h.suffix = "2006-01-02"
	default:
		h.interval = 3600
		h.suffix = "2006-01-02_15"
	}

	var err error
	now := time.Now()
	fileName = h.baseName + now.Format(h.suffix)
	h.fd, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}
	fInfo, _ := h.fd.Stat()
	h.rolloverAt = fInfo.ModTime().Unix() + h.interval
	return h
}

func (h *TimeRotatingFileHandler) doRollover() {
	now := time.Now()

	if h.rolloverAt < now.Unix() {
		fName := h.baseName + now.Format(h.suffix)
		h.fd.Close()
		err := os.Rename(h.baseName, fName)
		if err != nil {
			return
		}
		h.fd, _ = os.OpenFile(h.baseName, os.O_CREATE|os.O_APPEND, 0666)
		h.rolloverAt = now.Unix() + h.interval
	}
}

func (h *TimeRotatingFileHandler) Write(b []byte) (int, error) {
	h.doRollover()
	return h.fd.Write(b)
}

func (h *TimeRotatingFileHandler) Close() error {
	return h.fd.Close()
}
