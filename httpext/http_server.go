package httpext

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	. "github.com/aywfelix/felixgo/logger"
	. "github.com/aywfelix/felixgo/thread"
)

type RequestHandler func(form *HttpForm) interface{}

type HttpServer struct {
	server *http.Server
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		server: nil,
	}
}

func (h *HttpServer) StartServer(uri string) {
	GoRun(func() {
		h.server = &http.Server{
			Addr:           uri,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		err := h.server.ListenAndServe()
		if err != nil {
			LogError("http server listen error, %s", err.Error())
			return
		}
		LogInfo("start http server listen...")
	})
}

func (h *HttpServer) RegisterRouter(router string, handler RequestHandler) {
	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		var res interface{}
		defer func() {
			err := recover()
			if err != nil {
				log := fmt.Sprintf("http server handler error (%s): %s\nstack: %s", router, err, debug.Stack())
				res = log
			}
			LogInfo("http server handle %s request", router)
			resp, ok := res.(HttpResponse)
			if ok {
				if resp.Header != nil {
					for k, v := range resp.Header {
						w.Header().Set(k, v)
					}
				}
				h.writeResponse(w, resp.Body)
			} else {
				if res != nil {
					h.writeResponse(w, res)
				}
			}
		}()
		form := &HttpForm{Request: r, Response: w}
		res = handler(form)
	}
	http.HandleFunc(router, httpHandler)
}

func (h *HttpServer) StopServer() {
	h.server.Close()
	h.server = nil
	LogInfo("http server stop...")
}

func (h *HttpServer) writeResponse(w http.ResponseWriter, data interface{}) {
	bytes, err := JsonHelper.DataToBytes(data)
	if err != nil {
		return
	}

	count := len(bytes)
	for offset := 0; offset < count; {
		wlen, err := w.Write(bytes[offset:])
		if err != nil {
			LogError("http server response error: %s", err.Error())
			return
		}
		offset += wlen
	}
}
