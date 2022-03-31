package utils

import (
	"reflect"
)

// IsNil checks whether `value` is nil.
func IsNil(value interface{}) bool {
	return value == nil
}

// IsEmpty checks whether `value` is empty.
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	switch value := value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return value == 0
	case float32, float64:
		return value == 0
	case bool:
		return value == false
	case string:
		return value == ""
	case []byte:
		return len(value) == 0
	case []rune:
		return len(value) == 0
	case []int:
		return len(value) == 0
	case []string:
		return len(value) == 0
	case []float32:
		return len(value) == 0
	case []float64:
		return len(value) == 0
	case map[string]interface{}:
		return len(value) == 0
	default:
		var rv reflect.Value
		if v, ok := value.(reflect.Value); ok {
			rv = v
		} else {
			rv = reflect.ValueOf(value)
		}

		switch rv.Kind() {
		case reflect.Bool:
			return !rv.Bool()
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			return rv.Int() == 0
		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Uintptr:
			return rv.Uint() == 0
		case reflect.Float32,
			reflect.Float64:
			return rv.Float() == 0
		case reflect.String:
			return rv.Len() == 0
		case reflect.Struct:
			for i := 0; i < rv.NumField(); i++ {
				if !IsEmpty(rv) {
					return false
				}
			}
			return true
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Array:
			return rv.Len() == 0
		case reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return true
			}
		}
	}
	return false
}

// IsInt checks whether `value` is type of int.
func IsInt(value interface{}) bool {
	switch value.(type) {
	case int, *int, int8, *int8, int16, *int16, int32, *int32, int64, *int64:
		return true
	}
	return false
}

// IsUint checks whether `value` is type of uint.
func IsUint(value interface{}) bool {
	switch value.(type) {
	case uint, *uint, uint8, *uint8, uint16, *uint16, uint32, *uint32, uint64, *uint64:
		return true
	}
	return false
}

// IsFloat checks whether `value` is type of float.
func IsFloat(value interface{}) bool {
	switch value.(type) {
	case float32, *float32, float64, *float64:
		return true
	}
	return false
}

// IsSlice checks whether `value` is type of slice.
func IsSlice(value interface{}) bool {
	var (
		reflectValue = reflect.ValueOf(value)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		return true
	}
	return false
}

// IsMap checks whether `value` is type of map.
func IsMap(value interface{}) bool {
	var (
		reflectValue = reflect.ValueOf(value)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	switch reflectKind {
	case reflect.Map:
		return true
	}
	return false
}

// IsStruct checks whether `value` is type of struct.
func IsStruct(value interface{}) bool {
	var (
		reflectValue = reflect.ValueOf(value)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Struct:
		return true
	}
	return false
}
