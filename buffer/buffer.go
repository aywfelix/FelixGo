package buffer

import (
	"fmt"
)

const singlebuffer_len = 10

type SingleBuffer struct {
	reader int
	writer int
	buffer []byte
	length int
}

func NewSingleBuffer() *SingleBuffer {
	return &SingleBuffer{
		reader: 0,
		writer: 0,
		length: singlebuffer_len,
		buffer: make([]byte, singlebuffer_len),
	}
}

// 一定可以写入所有数据
func (sb *SingleBuffer) Write(data []byte) {
	dataLen := len(data)

	left := sb.length - sb.writer
	if left >= dataLen {
		copy(sb.buffer[sb.writer:], data)
		sb.writer += dataLen
		fmt.Println("write data directly")
		return
	}
	if sb.reader >= dataLen {
		copy(sb.buffer, sb.buffer[sb.reader:sb.writer])
		sb.writer -= sb.reader
		sb.reader = 0
		copy(sb.buffer[sb.writer:], data)
		fmt.Println("move reader to front, write data directly")
		return
	}

	needLen := (sb.writer - sb.reader) + dataLen
	for sb.length < needLen {
		sb.length *= 2
	}
	buffer := make([]byte, sb.length)
	copy(buffer, sb.buffer[sb.reader:sb.writer])
	sb.reader = 0
	sb.writer -= sb.reader
	sb.buffer = buffer
	copy(buffer, data)
	sb.writer += dataLen
	fmt.Println("expand buffer, write data")
}

func (sb *SingleBuffer) ReadAll() ([]byte, int) {
	data := sb.buffer[sb.reader:sb.writer]
	datalen := sb.writer - sb.reader
	sb.reader = 0
	sb.writer = 0
	return data, datalen
}

func (sb *SingleBuffer) Read(readlen int) ([]byte, int) {
	datalen := sb.writer - sb.reader
	retlen := 0
	var data []byte
	if datalen <= readlen {
		data = sb.buffer[sb.reader:sb.writer]
		sb.reader = sb.writer
		retlen = datalen
	} else {
		data = sb.buffer[sb.reader : sb.reader+readlen]
		sb.reader += readlen
		retlen = readlen
	}
	if sb.reader == sb.writer {
		sb.reader = 0
		sb.writer = 0
	}
	return data, retlen
}

func (sb *SingleBuffer) Reader() []byte {
	return sb.buffer[sb.reader:]
}

func (sb *SingleBuffer) Writer() []byte {
	return sb.buffer[sb.writer:]
}

func (sb *SingleBuffer) ReaderOffset() int {
	return sb.reader
}

func (sb *SingleBuffer) WriterOffset() int {
	return sb.writer
}

func (sb *SingleBuffer) Len() int {
	return sb.length
}

func (sb *SingleBuffer) DataLen() int {
	return sb.writer - sb.reader
}

func (sb *SingleBuffer) IsEmpty() bool {
	return sb.reader == sb.writer
}

func (sb *SingleBuffer) IsFull() bool {
	return sb.writer == sb.length
}

func (sb *SingleBuffer) Clear() {
	sb.reader = 0
	sb.writer = 0
}
