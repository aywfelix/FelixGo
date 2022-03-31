package utils

import (
	"reflect"
)

type ObjService struct {
	servers map[string]reflect.Method
	rcvr    reflect.Value
	typ     reflect.Type
}

func NewObjService(rep interface{}) *ObjService {
	service := new(ObjService)
	service.typ = reflect.TypeOf(rep)
	service.rcvr = reflect.ValueOf(rep)
	service.servers = map[string]reflect.Method{}
	for i := 0; i < service.typ.NumMethod(); i++ {
		method := service.typ.Method(i)
		service.servers[method.Name] = method
	}
	return service
}

func (s *ObjService) Call(methodName string, args ...interface{}) {
	for funcName, method := range s.servers {
		if funcName == methodName {
			values := []reflect.Value{}
			values = append(values, s.rcvr)
			for _, arg := range args {
				v := reflect.ValueOf(arg)
				values = append(values, v)
			}
			method.Func.Call(values)
		}
	}
}
