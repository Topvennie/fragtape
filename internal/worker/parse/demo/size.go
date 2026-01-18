package demo

import (
	"encoding/json"
	"reflect"
	"unsafe"
)

func getDeepSize(val reflect.Value) int64 {
	switch val.Kind() {
	case reflect.Pointer:
		if val.IsNil() {
			return 0
		}
		return int64(unsafe.Sizeof(val.Pointer())) + getDeepSize(val.Elem())

	case reflect.Slice:
		if val.IsNil() {
			return 0
		}
		size := int64(val.Type().Size())
		for i := 0; i < val.Len(); i++ {
			size += getDeepSize(val.Index(i))
		}
		return size

	case reflect.String:
		return int64(unsafe.Sizeof(val.String())) + int64(val.Len())

	case reflect.Struct:
		var size int64
		size = int64(val.Type().Size())

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			switch field.Kind() {
			case reflect.Slice, reflect.String, reflect.Pointer, reflect.Map:
				size += getDeepSize(field) - int64(field.Type().Size())
			}
		}
		return size

	case reflect.Map:
		if val.IsNil() {
			return 0
		}
		size := int64(unsafe.Sizeof(val.Pointer()))
		iter := val.MapRange()
		for iter.Next() {
			size += getDeepSize(iter.Key())
			size += getDeepSize(iter.Value())
			size += 16
		}
		return size

	default:
		return int64(val.Type().Size())
	}
}

func getJSONSize(v any) int64 {
	bytes, err := json.Marshal(v)
	if err != nil {
		return 0
	}

	return int64(len(bytes))
}
