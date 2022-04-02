package event

import (
	"reflect"

	. "github.com/aywfelix/felixgo/logger"
)

type OnEventHandler func(et EventType, args ...interface{})

type handlerSlice []OnEventHandler

type IEvent interface {
	RegEvent(et EventType, handler OnEventHandler)
	DelEvent(et EventType)
	DispatchEvent(et EventType, params []interface{})
}

type Event struct {
	events map[EventType]handlerSlice
}

func NewEvent() *Event {
	return &Event{
		events: make(map[EventType]handlerSlice),
	}
}

func (e *Event) RegEvent(et EventType, handler OnEventHandler) {
	handlerSlice, ok := e.events[et]
	if !ok {
		handlerSlice = make([]OnEventHandler, 0)
	} else {
		for _, f := range handlerSlice {
			vf := reflect.ValueOf(f)
			vhandler := reflect.ValueOf(handler)
			fname := reflect.TypeOf(handler).String()
			if vf.Pointer() == vhandler.Pointer() {
				LogError("duplicate register event handler, f=", fname)
				return
			}
		}
	}
	handlerSlice = append(handlerSlice, handler)
	e.events[et] = handlerSlice
}

func (e *Event) DelEvent(et EventType) {
	if _, ok := e.events[et]; ok {
		delete(e.events, et)
	}
}

func (e *Event) DispatchEvent(et EventType, args ...interface{}) {
	handlerSlice, ok := e.events[et]
	if !ok {
		LogError("dipatch event error, event type=", et)
		return
	}
	for _, f := range handlerSlice {
		f(et, args...)
	}
	LogDebug("dispatch event, event type=", et)
}
