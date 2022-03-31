package global

import (
	"context"
)

type (
	Ctx = context.Context
)

type (
	Map        = map[string]interface{}
	MapAnyAny  = map[interface{}]interface{}
	MapAnyStr  = map[interface{}]string
	MapAnyInt  = map[interface{}]int
	MapStrAny  = map[string]interface{}
	MapStrStr  = map[string]string
	MapStrInt  = map[string]int
	MapIntAny  = map[int]interface{}
	MapIntStr  = map[int]string
	MapIntInt  = map[int]int
	MapAnyBool = map[interface{}]bool
	MapStrBool = map[string]bool
	MapIntBool = map[int]bool
)

type (
	List        = []Map
	ListAnyAny  = []MapAnyAny
	ListAnyStr  = []MapAnyStr
	ListAnyInt  = []MapAnyInt
	ListStrAny  = []MapStrAny
	ListStrStr  = []MapStrStr
	ListStrInt  = []MapStrInt
	ListIntAny  = []MapIntAny
	ListIntStr  = []MapIntStr
	ListIntInt  = []MapIntInt
	ListAnyBool = []MapAnyBool
	ListStrBool = []MapStrBool
	ListIntBool = []MapIntBool
)

type (
	Array    = []interface{}
	ArrayAny = []interface{}
	ArrayStr = []string
	ArrayInt = []int
)
