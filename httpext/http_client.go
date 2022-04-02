package httpext

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	. "github.com/aywfelix/felixgo/container/string"
	. "github.com/aywfelix/felixgo/logger"
)

const (
	max_idle_count          = 100
	max_idle_conns_per_host = 100
	idle_conn_timeout       = 60

	request_retry   = 6
	request_timeout = 6 * time.Second
)

//============================================================================================
var httpClient *http.Client

//============================================================================================

type HttpRequest struct {
	host    string
	reqType string
}

func NewHttpRequest(host string, reqType string) *HttpRequest {
	return &HttpRequest{host: host, reqType: reqType}
}

func (r *HttpRequest) fixUri(router string) (*StringBuilder, string) {
	builder := new(StringBuilder)
	if !strings.HasPrefix(router, "http") {
		builder.Append(r.host)
		if !strings.HasSuffix(r.host, "/") {
			builder.Append("/")
		}
	}
	builder.Append(router)
	return builder, builder.String()
}

func (r *HttpRequest) Request(router string, params interface{}) ([]byte, error) {
	defer func() {
		err := recover()
		if err != nil {
			LogError("http client post request error, %s", err)
			LogError("error stack, %s", debug.Stack())
			return
		}
	}()

	var (
		request  *http.Request
		response *http.Response
		data     []byte
		err      error
	)

	builder, uri := r.fixUri(router)

	if r.reqType == http.MethodPost {
		if params != nil {
			data, err = JsonHelper.DataToBytes(params)
		}
		request, err = http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(data))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	} else {
		if params != nil {
			jsonMap, ok := params.(JsonMap)
			if ok && jsonMap != nil {
				builder.Append("?")
				jsonMap.ToUrl(builder)
			}
		}
		uri := builder.String()
		request, err = http.NewRequest(http.MethodGet, uri, nil)
		request.Header.Add("Content-Type", "application/octet-stream")
	}

	for i := 0; i < request_retry; i++ {
		response, err = httpClient.Do(request)
		if err != nil {
			LogError("http client request server failed, %s", err.Error())
			continue
		}
		break
	}
	if response == nil || err != nil {
		return nil, err
	}

	defer response.Body.Close()
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		LogError("http client read data failed, err: %s", err.Error())
		return nil, err
	}
	return data, nil
}

func (r *HttpRequest) RequestPost(router string, data []byte, header map[string]string) ([]byte, error) {
	defer func() {
		err := recover()
		if err != nil {
			LogError("http client post request failed, err: %s", err)
			LogError("error stack, %s", debug.Stack())
			return
		}
	}()

	var request *http.Request
	var response *http.Response
	var err error

	_, uri := r.fixUri(router)

	request, err = http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(data))
	if request == nil || err != nil {
		LogError("new http request failed, err: %s", err.Error())
		return nil, err
	}

	request.Header.Add("Content-Type", "application/octet-stream")
	if header != nil {
		for key, value := range header {
			request.Header.Set(key, value)
		}
	}
	for i := 0; i < request_retry; i++ {
		response, err = httpClient.Do(request)
		if err == nil {
			break
		}
		LogError("client request server failed, err: %s", err.Error())
	}
	if response == nil || err != nil {
		return nil, err
	}

	defer response.Body.Close()
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		LogError("http client read data failed, err: %s", err.Error())
		return nil, err
	}
	return data, nil
}

func (c *HttpRequest) RequestGet(router string, header map[string]string) (*http.Response, error) {
	defer func() {
		err := recover()
		if err != nil {
			LogError("http client post request error, %s", err)
			LogError("error stack, %s", debug.Stack())
			return
		}
	}()

	var request *http.Request
	var response *http.Response
	var err error

	_, uri := c.fixUri(router)

	request, err = http.NewRequest(http.MethodGet, uri, nil)
	if request == nil || err != nil {
		LogError("new http request failed, err: %s", err.Error())
		return nil, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	if header != nil {
		for key, value := range header {
			request.Header.Set(key, value)
		}
	}

	for i := 0; i < request_retry; i++ {
		response, err = httpClient.Do(request)
		if err == nil {
			break
		}
		LogError("client request server failed, err: %s", err.Error())
	}
	if response == nil || err != nil {
		return nil, err
	}
	return response, nil
}

func (r *HttpRequest) RequestJsonMap(router string, params interface{}) (JsonMap, error) {
	data, err := r.Request(router, params)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return JsonHelper.UnmarshalMap(data)
	}
	return nil, nil
}

func (r *HttpRequest) RequestJsonArray(router string, params interface{}) (JsonArray, error) {
	data, err := r.Request(router, params)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return JsonHelper.UnmarshalArray(data)
	}
	return nil, nil
}

var HttpPost *HttpRequest
var HttpGet *HttpRequest

func init() {
	httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			MaxIdleConns:        max_idle_count,
			MaxIdleConnsPerHost: max_idle_conns_per_host,
			IdleConnTimeout:     time.Duration(idle_conn_timeout) * time.Second,
		},
		Timeout: request_timeout,
	}
	HttpPost = NewHttpRequest("http://127.0.0.1:8080", "POST")
	HttpGet = NewHttpRequest("http://127.0.0.1:8080", "GET")
}
