package httpext

import (
	"fmt"
	. "github.com/felix/felixgo/logger"
	"net/http"
	"runtime/debug"
	. "github.com/felix/felixgo/thread"
)

type HttpAsyncHandler func(form *HttpFormAsync) interface{}

type AsyncHttpServer struct {
	httpServer *HttpServer
	taskPool   *TaskPool
}

type HttpFormAsync struct {
	router string
	Resp   chan interface{}
	HttpForm
}

func NewAsyncHttpServer() *AsyncHttpServer {
	return &AsyncHttpServer{
		httpServer: NewHttpServer(),
		taskPool:   NewTaskPool(500, 1000),
	}
}

func (h *AsyncHttpServer) StartServer(uri string) {
	h.httpServer.StartServer(uri)
	h.taskPool.Start()
}

func (h *AsyncHttpServer) RegisterRouter(router string, handler HttpAsyncHandler) {
	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		var res interface{}
		defer func() {
			err := recover()
			if err != nil {
				log := fmt.Sprintf("http server handler(%s) error: %s\nstack: %s", router, err, debug.Stack())
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

		form := &HttpFormAsync{
			HttpForm: HttpForm{Request: r, Response: w},
			Resp:     make(chan interface{}),
			router:   router,
		}

		h.taskPool.Submit(NewTask(func(args ...interface{}) error {
			handler(form)
			return nil
		}))

		res = <-form.Resp
	}
	http.HandleFunc(router, httpHandler)
}

func (h *AsyncHttpServer) StopServer() {
	LogInfo("http server stop...")
	h.httpServer.server.Close()
	h.taskPool.Stop()
}

func (h *AsyncHttpServer) writeResponse(w http.ResponseWriter, data interface{}) {
	bytes, err := JsonHelper.DataToBytes(data)
	if err != nil {
		return
	}

	count := len(bytes)
	for offset := 0; offset < count; {
		wlen, err := w.Write(bytes[offset:])
		if err != nil {
			LogError("http server response failed, err: %s", err.Error())
			return
		}
		offset += wlen
	}
}
