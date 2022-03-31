package handler

import (
	"encoding/binary"
	"net"
	"time"
)

type SocketHandler struct {
	c        net.Conn
	protocol string
	addr     string
}

func NewSocketHandler(protocol string, addr string) *SocketHandler {
	s := new(SocketHandler)
	s.protocol = protocol
	s.addr = addr
	return s
}

func (h *SocketHandler) Write(b []byte) (int, error) {
	if err := h.connect(); err != nil {
		return 0, err
	}

	buf := make([]byte, len(b)+4)
	binary.LittleEndian.PutUint32(buf, uint32(len(b)))
	copy(buf[4:], b)
	n, err := h.c.Write(buf)
	if err != nil {
		h.c.Close()
		h.c = nil
		return n, err
	}
	return n, nil
}

func (h *SocketHandler) Close() error {
	if h.c != nil {
		return h.c.Close()
	}
	return nil
}

func (h *SocketHandler) connect() error {
	if h.c != nil {
		return nil
	}

	var err error
	h.c, err = net.DialTimeout(h.protocol, h.addr, 20*time.Second)
	if err != nil {
		return err
	}
	return nil
}
