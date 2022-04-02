package httpext

import (
	"io/ioutil"
	"net/http"
	"strconv"

	. "github.com/aywfelix/felixgo/logger"
)

//======================================================================================
type HttpForm struct {
	Request  *http.Request
	Response http.ResponseWriter
	reqData  JsonObject
	err      error
}

func (h *HttpForm) GetError() error {
	return h.err
}

func (h *HttpForm) HeaderValue(key string) string {
	return h.Request.Header.Get(key)
}

func (h *HttpForm) readReqData() JsonObject {
	if h.Request.Method == http.MethodPost {
		defer h.Request.Body.Close()
		data, err := ioutil.ReadAll(h.Request.Body)
		if err != nil {
			LogError("read data from http form failed, err: %s", err.Error())
			return nil
		}
		h.reqData, err = JsonHelper.Unmarshal(data)
		if err != nil {
			LogError("json unmarshal failed, err: %s", err.Error())
			return nil
		}
		return h.reqData
	}
	LogError("read nil from http form")
	return nil
}

func (h *HttpForm) getReqData() interface{} {
	if h.reqData == nil {
		return h.readReqData()
	}
	return h.reqData
}

func (h *HttpForm) reqDataToJsonMap() JsonMap {
	jsonMap, ok := h.getReqData().(JsonMap)
	if !ok {
		LogError("form data is nil")
		return nil
	}
	return jsonMap
}

func (h *HttpForm) reqDataToJsonArray() JsonArray {
	jsonArr, ok := h.getReqData().(JsonArray)
	if !ok {
		LogError("form data is nil")
		return nil
	}
	return jsonArr
}

func (h *HttpForm) GetString(key string) string {
	jsonMap := h.reqDataToJsonMap()
	if jsonMap != nil {
		return jsonMap.GetStr(key)
	}
	return h.Request.FormValue(key)
}

func (h *HttpForm) GetInt(key string) int {
	jsonMap := h.reqDataToJsonMap()
	if jsonMap != nil {
		return jsonMap.GetInt(key)
	} else {
		str := h.Request.FormValue(key)
		v, err := strconv.Atoi(str)
		if err != nil {
			h.err = err
			return 0
		}
		return v
	}
}

func (h *HttpForm) GetInt64(key string) int64 {
	jsonMap := h.reqDataToJsonMap()
	if jsonMap != nil {
		return jsonMap.GetInt64(key)
	} else {
		str := h.Request.FormValue(key)
		v, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			h.err = err
			return 0
		}
		return v
	}
}

func (h *HttpForm) GetBool(key string) bool {
	jsonMap := h.reqDataToJsonMap()
	if jsonMap != nil {
		return jsonMap.GetBool(key)
	} else {
		str := h.Request.FormValue(key)
		b, err := strconv.ParseBool(str)
		if err != nil {
			h.err = err
			return false
		}
		return b
	}
}

func (h *HttpForm) GetUint(key string) uint {
	jsonMap := h.reqDataToJsonMap()
	if jsonMap != nil {
		return jsonMap.GetUint(key)
	} else {
		str := h.Request.FormValue(key)
		v, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			h.err = err
			return 0
		}
		return uint(v)
	}
}

func (h *HttpForm) GetUint64(key string) uint64 {
	jsonMap := h.reqDataToJsonMap()
	if jsonMap != nil {
		return jsonMap.GetUint64(key)
	} else {
		str := h.Request.FormValue(key)
		v, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			h.err = err
			return 0
		}
		return v
	}
}

func (h *HttpForm) GetMap(key string) JsonMap {
	jsonMap := h.reqDataToJsonMap()
	if jsonMap != nil {
		return jsonMap.GetMap(key)
	} else {
		str := h.Request.FormValue(key)
		jsonMap, err := JsonHelper.UnmarshalMap([]byte(str))
		if err != nil {
			return nil
		}
		return jsonMap
	}
}

func (h *HttpForm) GetArray(key string) JsonArray {
	jsonMap := h.reqDataToJsonMap()
	if jsonMap != nil {
		return jsonMap.GetArray(key)
	} else {
		str := h.Request.FormValue(key)
		jsonArr, err := JsonHelper.UnmarshalArray([]byte(str))
		if err != nil {
			return nil
		}
		return jsonArr
	}
}

//======================================================================================

type HttpResponse struct {
	Header map[string]string
	Body   interface{}
}

func (h *HttpResponse) SetHeader(key string, value string) {
	if h.Header == nil {
		h.Header = make(map[string]string)
	}
	h.Header[key] = value
}
