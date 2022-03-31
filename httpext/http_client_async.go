package httpext

import (
	"net/http"
	
	. "github.com/felix/felixgo/logger"
	. "github.com/felix/felixgo/thread"
)

type RespHandler func(reqResult *RequestAsyncRet)
type RespErrorHandler func(err error)

type HttpRequestAsync struct {
	httpRequest    *HttpRequest
	taskPool       *TaskPool
	reqResultQueue chan *RequestAsyncRet
}

type RequestAsyncRet struct {
	retData     []byte
	retResp     *http.Response
	retMap      JsonMap
	retArray    JsonArray
	retErr      error
	respHandler RespHandler
	errHandler  RespErrorHandler
}

func NewHttpRequestAsync(host string, reqType string) *HttpRequestAsync {
	return &HttpRequestAsync{
		httpRequest:    NewHttpRequest(host, reqType),
		taskPool:       NewTaskPool(500, 1000),
		reqResultQueue: make(chan *RequestAsyncRet, 1000),
	}
}

func (r *HttpRequestAsync) Request(router string, params interface{}, respHandler RespHandler, errHandler RespErrorHandler) {
	r.taskPool.Submit(NewTask(func(args ...interface{}) error {
		data, err := r.httpRequest.Request(router, params)
		if err != nil {
			LogError("HttpRequestAsync request error, %s", err.Error())
			return err
		}
		reqAsyncRet := &RequestAsyncRet{
			retData:     data,
			retResp:     nil,
			retMap:      nil,
			retArray:    nil,
			retErr:      err,
			respHandler: respHandler,
			errHandler:  errHandler,
		}
		r.reqResultQueue <- reqAsyncRet
		return nil
	}))
}

func (r *HttpRequestAsync) RequestPost(router string, data []byte, header map[string]string, respHandler RespHandler, errHandler RespErrorHandler) {
	r.taskPool.Submit(NewTask(func(args ...interface{}) error {
		data, err := r.httpRequest.RequestPost(router, data, header)
		if err != nil {
			LogError("HttpRequestAsync request error, %s", err.Error())
			return err
		}
		reqAsyncRet := &RequestAsyncRet{
			retData:     data,
			retResp:     nil,
			retMap:      nil,
			retArray:    nil,
			retErr:      err,
			respHandler: respHandler,
			errHandler:  errHandler,
		}
		r.reqResultQueue <- reqAsyncRet
		return nil
	}))
}

func (r *HttpRequestAsync) RequestGet(router string, header map[string]string, respHandler RespHandler, errHandler RespErrorHandler) {
	r.taskPool.Submit(NewTask(func(args ...interface{}) error {
		resp, err := r.httpRequest.RequestGet(router, header)
		if err != nil {
			LogError("HttpRequestAsync, RequestGet, request error, %s", err.Error())
			return err
		}
		reqAsyncRet := &RequestAsyncRet{
			retData:     nil,
			retResp:     resp,
			retMap:      nil,
			retArray:    nil,
			retErr:      err,
			respHandler: respHandler,
			errHandler:  errHandler,
		}
		r.reqResultQueue <- reqAsyncRet
		return nil
	}))
}

func (r *HttpRequestAsync) RequestJsonMap(router string, params interface{}, respHandler RespHandler, errHandler RespErrorHandler) {
	r.taskPool.Submit(NewTask(func(args ...interface{}) error {
		retMap, err := r.httpRequest.RequestJsonMap(router, params)
		if err != nil {
			LogError("HttpRequestAsync, RequestJsonMap, request error, %s", err.Error())
			return err
		}
		reqAsyncRet := &RequestAsyncRet{
			retData:     nil,
			retResp:     nil,
			retMap:      retMap,
			retArray:    nil,
			retErr:      err,
			respHandler: respHandler,
			errHandler:  errHandler,
		}
		r.reqResultQueue <- reqAsyncRet
		return nil
	}))
}

func (r *HttpRequestAsync) RequestJsonArray(router string, params interface{}, respHandler RespHandler, errHandler RespErrorHandler) {
	r.taskPool.Submit(NewTask(func(args ...interface{}) error {
		retArray, err := r.httpRequest.RequestJsonArray(router, params)
		if err != nil {
			LogError("HttpRequestAsync, RequestJsonArray, request error, %s", err.Error())
			return err
		}
		reqAsyncRet := &RequestAsyncRet{
			retData:     nil,
			retResp:     nil,
			retMap:      nil,
			retArray:    retArray,
			retErr:      err,
			respHandler: respHandler,
			errHandler:  errHandler,
		}
		r.reqResultQueue <- reqAsyncRet
		return nil
	}))
}

func (r *HttpRequestAsync) HandleHttpResponse() {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				LogError("HttpRequestAsync handle response error, %s", err)
			}
		}()

		for resp := range r.reqResultQueue {
			resp.respHandler(resp)
			if resp.errHandler != nil && resp.retErr != nil {
				resp.errHandler(resp.retErr)
			}
		}
	}()
}

func (r *HttpRequestAsync) Release() {
	close(r.reqResultQueue)
}

var HttpPostAsync *HttpRequestAsync
var HttpGetAsync *HttpRequestAsync

func StartHttpPostAsync() {
	HttpPostAsync = NewHttpRequestAsync("http://127.0.0.1:8080", "POST")
	HttpPostAsync.taskPool.Start()
}

func StartHttpGetAsync() {
	HttpGetAsync = NewHttpRequestAsync("http://127.0.0.1:8080", "GET")
	HttpGetAsync.taskPool.Start()
}
