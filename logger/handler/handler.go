package handler

import (
	"io"
)

type IHandler interface {
	Write(p []byte) (int, error)
	Close() error
}

type StreamHandler struct {
	w io.Writer
}

func NewStreamHandler(w io.Writer) *StreamHandler {
	h := new(StreamHandler)
	h.w = w
	return h
}

func (h *StreamHandler) Write(b []byte) (int, error) {
	return h.w.Write(b)
}

func (h *StreamHandler) Close() error {
	return nil
}

type NullHandler struct {
}

func NewNullHandler() *NullHandler {
	return new(NullHandler)
}

func (h *NullHandler) Write(b []byte) (int, error) {
	return 0, nil
}

func (h *NullHandler) Close() error {
	return nil
}
