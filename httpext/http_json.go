package httpext

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"

	. "github.com/felix/felixgo/logger"
	. "github.com/felix/felixgo/container/string"
)

var (
	errJsonFormatInvalid = errors.New("error json format")
	errJsonTypeInvalid   = errors.New("error json object type")
	errDataToBytes       = errors.New("data convert to bytes failed")
)

type JsonObject interface{}
type JsonMap map[string]interface{}
type JsonArray []interface{}

const error_key = "error_keys"

//======================json map==============================================
func (this JsonMap) Get(key string) JsonObject {
	value, ok := this[key]
	if !ok {
		keys, ok := this[error_key]
		if ok {
			if !strings.Contains(keys.(string), key) {
				keys = fmt.Sprintf("%s,%s", keys, key)
				this[error_key] = keys
			}
		} else {
			this[error_key] = key
		}
		return nil
	}
	return value
}

func (this JsonMap) GetStr(key string) string {
	value := this.Get(key)
	if value == nil {
		return ""
	}
	return value.(string)
}

func (this JsonMap) GetInt(key string) int {
	value := this.Get(key)
	if value == nil {
		return 0
	}

	vf64, ok := value.(float64)
	if ok {
		return int(vf64)
	}

	vInt, ok := value.(int)
	if ok {
		return vInt
	}

	LogError("Json object map get int failed: %s", key)
	return 0
}

func (this JsonMap) GetBool(key string) bool {
	value := this.Get(key)
	if value == nil {
		return false
	}

	vb, ok := value.(bool)
	if !ok {
		LogError("Json object map get int failed: %s", key)
		return false
	}

	return vb
}

func (this JsonMap) GetUint(key string) uint {
	value := this.Get(key)
	if value == nil {
		return 0
	}

	if val, ok := value.(string); ok {
		num, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			LogError("Json object map get uint failed: %s, key=%s, value=%v", err.Error(), key, value)
			return 0
		}
		return uint(num)
	}

	if vf64, ok := value.(float64); ok {
		return uint(vf64)
	}

	if val, ok := value.(uint); ok {
		return val
	}

	LogError("Json object map get int failed: %s", key)
	return 0
}

func (this JsonMap) GetInt64(key string) int64 {
	value := this.Get(key)
	if value == nil {
		return 0
	}

	vf64, ok := value.(float64)
	if !ok {
		LogError("Json object map get int failed: %s", key)
		return 0
	}

	return int64(vf64)
}

func (this JsonMap) GetUint64(key string) uint64 {
	value := this.Get(key)
	if value == nil {
		return 0
	}

	vf64, ok := value.(float64)
	if !ok {
		LogError("Json object map get int failed: %s", key)
		return 0
	}

	return uint64(vf64)
}

func (this JsonMap) GetMap(key string) JsonMap {
	m, ok := this.Get(key).(JsonMap)
	if !ok {
		return nil
	}
	return m
}

func (this JsonMap) GetArray(key string) JsonArray {
	arr, ok := this.Get(key).(JsonArray)
	if !ok {
		return nil
	}
	return arr
}

func (this JsonMap) GetErrorKeys() string {
	keys, ok := this[error_key]
	if ok {
		return keys.(string)
	}
	return ""
}

func (this JsonMap) ToUrl(sb *StringBuilder) {
	first := true
	for key, value := range this {
		if key == error_key {
			continue
		}

		if !first {
			sb.Append("&")
		} else {
			first = false
		}

		sb.Append(key)
		sb.Append("=")

		text := fmt.Sprintf("%v", value)
		text = url.QueryEscape(text)
		sb.Append(text)
	}
}

func (this JsonMap) WriteToHttp(writer http.ResponseWriter) error {
	bytes, err := JsonHelper.Marshal(this)
	if err != nil {
		LogError("DataMap encode data to json error: %s, data=%s", err.Error(), this)
		return err
	}
	_, err = writer.Write(bytes)
	return err
}

//======================json arrays==============================================
func (this JsonArray) Get(i int) JsonObject {
	if i < len(this) {
		return this[i]
	}
	return nil
}

func (this JsonArray) GetStr(i int) string {
	value := this.Get(i)
	if value == nil {
		return ""
	}
	return value.(string)
}

func (this JsonArray) GetInt(i int) int {
	value := this.Get(i)
	if value == nil {
		return -1
	}

	vi, ok := value.(int)
	if ok {
		return vi
	}

	vf, ok := value.(float64)
	if ok {
		return int(vf)
	}

	vs, ok := value.(string)
	if ok {
		num, err := strconv.ParseUint(vs, 10, 32)
		if err != nil {
			LogError("Json object array get int failed: %s, value=%v", err.Error(), value)
			return -1
		}
		return int(num)
	}
	return -1
}

func (this JsonArray) GetUint(i int) uint {
	value := this.GetInt(i)
	return uint(value)
}

func (this JsonArray) GetMap(i int) JsonMap {
	m, ok := this.Get(i).(JsonMap)
	if ok {
		return m
	}
	return nil
}

func (this JsonArray) GetArray(i int) JsonArray {
	arr, ok := this.Get(i).(JsonArray)
	if ok {
		return arr
	}
	return nil
}

type JsonType struct{}

var JsonHelper JsonType

func (this *JsonType) Marshal(obj JsonObject) ([]byte, error) {
	defer func() {
		err := recover()
		if err != nil {
			LogError("Marshal json %v error: %s", obj, err)
			LogError("stack: %s", debug.Stack())
		}
	}()
	return json.Marshal(obj)
}

func (this *JsonType) Unmarshal(bytes []byte) (JsonObject, error) {
	defer func() {
		err := recover()
		if err != nil {
			LogError("Unmarshal json %v error: %s", string(bytes), err)
			LogError("stack: %s", debug.Stack())
		}
	}()

	var flag byte
	for i := range bytes {
		v := bytes[i]
		if v != ' ' {
			flag = v
			break
		}
	}

	if flag == '{' {
		jMap := JsonMap{}
		err := json.Unmarshal(bytes, &jMap)
		return jMap, err
	}

	if flag == '[' {
		jArr := JsonArray{}
		err := json.Unmarshal(bytes, &jArr)
		return jArr, err
	}

	return nil, errJsonFormatInvalid
}

func (this *JsonType) UnmarshalMap(bytes []byte) (JsonMap, error) {
	obj, err := this.Unmarshal(bytes)
	if err != nil {
		return nil, err
	}

	vmap, ok := obj.(JsonMap)
	if ok {
		return vmap, nil
	}
	return nil, errJsonTypeInvalid
}

func (this *JsonType) UnmarshalArray(bytes []byte) (JsonArray, error) {
	obj, err := this.Unmarshal(bytes)
	if err != nil {
		return nil, err
	}

	array, ok := obj.(JsonArray)
	if ok {
		return array, nil
	}
	return nil, errJsonTypeInvalid
}

func (this *JsonType) DataToBytes(data interface{}) ([]byte, error) {
	var bytes []byte
	switch data.(type) {
	case string:
		bytes = []byte(data.(string))
	case []byte:
		bytes = data.([]byte)
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		bytes = []byte(fmt.Sprintf("%d", data))
	case nil:
		bytes = []byte("")
	case JsonMap:
		dataBytes, err := JsonHelper.Marshal(data)
		if err != nil {
			return nil, err
		}
		bytes = dataBytes
	case *JsonMap:
		dataBytes, err := JsonHelper.Marshal(data)
		if err != nil {
			return nil, err
		}
		bytes = dataBytes
	case JsonArray:
		dataBytes, err := JsonHelper.Marshal(data)
		if err != nil {
			return nil, err
		}
		bytes = dataBytes
	case *JsonArray:
		dataBytes, err := JsonHelper.Marshal(data)
		if err != nil {
			return nil, err
		}
		bytes = dataBytes
	default:
		str := fmt.Sprintf("%v", data)
		bytes = []byte(str)
	}
	if bytes == nil {
		return nil, errDataToBytes
	}
	return bytes, nil
}

//====================================================================================
