package buffer

import (
	"fmt"
	"sync"
)

//===============================================================================
type Buffer struct {
	reader int
	writer int
	buffer []byte
	length int
	next   *Buffer
}

func newBuffer(len int) *Buffer {
	return &Buffer{
		reader: 0,
		writer: 0,
		length: len,
		buffer: make([]byte, len),
		next:   nil,
	}
}

func (b *Buffer) write(data []byte) int {
	dataLen := len(data)

	left := b.length - b.writer
	if left >= dataLen {
		copy(b.buffer[b.writer:], data)
		b.writer += dataLen
		b.toString()
		return dataLen
	}
	copy(b.buffer[b.writer:], data[:left])
	b.writer = b.length
	b.toString()
	return left
}

func (b *Buffer) readAll() []byte {
	data := b.buffer[b.reader:b.writer]
	b.reader = 0
	b.writer = 0
	b.toString()
	return data
}

func (b *Buffer) read(readlen int) ([]byte, int) {
	datalen := b.writer - b.reader
	retlen := 0
	var data []byte
	if datalen <= readlen {
		data = b.buffer[b.reader:b.writer]
		b.reader = b.writer
		retlen = datalen
	} else {
		data = b.buffer[b.reader : b.reader+readlen]
		b.reader += readlen
		retlen = readlen
	}
	b.toString()
	return data, retlen
}

func (b *Buffer) dataLen() int {
	return b.writer - b.reader
}

func (b *Buffer) isEmpty() bool {
	return b.reader == b.writer
}

func (b *Buffer) isFull() bool {
	return b.writer == b.length
}

func (b *Buffer) clear() {
	b.reader = 0
	b.writer = 0
}

func (b *Buffer) toString() {
	fmt.Println("buf wlen=", b.dataLen(), "reader=", b.reader, "writer=", b.writer, "isfull=", b.isFull(), "isempty=", b.isEmpty())
}

//==============================================================
const buffer_len = 10

var bufferPool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return newBuffer(buffer_len)
	},
}

type BufferList struct {
	head    *Buffer
	tail    *Buffer
	dataLen int
}

func NewBufferList() *BufferList {
	bl := &BufferList{
		head: nil,
		tail: nil,
	}

	return bl
}

func (bl *BufferList) Write(data []byte) {
	writeable := len(data)
	bl.dataLen += writeable

	var buffer *Buffer = bl.tail
	if buffer == nil {
		buffer = bufferPool.Get().(*Buffer)
		bl.head = buffer
		bl.tail = buffer
	}
	writeLen := 0
	for writeable > 0 {
		if buffer.isFull() {
			nextBuffer := bufferPool.Get().(*Buffer)
			buffer.next = nextBuffer
			bl.tail = nextBuffer
			buffer = nextBuffer
		}
		wlen := buffer.write(data[writeLen:])
		writeLen += wlen
		writeable -= wlen
	}
}

func (bl *BufferList) ReadAll() ([]byte, int) {
	readbuf := make([]byte, bl.dataLen)
	readLen := 0
	for bl.head != nil {
		buffer := bl.head
		rlen := buffer.dataLen()
		copy(readbuf[readLen:], buffer.readAll())
		readLen += rlen

		bufferPool.Put(buffer)

		bl.head = buffer.next
	}
	bl.dataLen = 0
	return readbuf, readLen
}

func (bl *BufferList) Read(len int) ([]byte, int) {
	var buffer *Buffer = bl.head
	if buffer == nil || len <= 0 {
		return nil, 0
	}
	var readbuf []byte
	readLen := 0
	readable := 0
	if bl.dataLen >= len {
		readable = len
	} else {
		readable = bl.dataLen
	}
	readbuf = make([]byte, readable)
	for bl.head != nil && readable > 0 {
		buffer = bl.head
		rbuf, rlen := buffer.read(readable)
		copy(readbuf[readLen:], rbuf)
		readLen += rlen
		readable -= rlen
		bl.dataLen -= rlen

		if buffer.isEmpty() {
			bl.head = bl.head.next
			buffer.clear()
			bufferPool.Put(buffer)
		}
	}

	return readbuf, readLen
}

func (bl *BufferList) DataLen() int {
	return bl.dataLen
}
