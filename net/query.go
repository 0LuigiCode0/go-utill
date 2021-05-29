package net

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func QueryMarshal(in interface{}) (out string) {
	query := url.Values{}
	valid(reflect.ValueOf(in), &query, "")
	out = query.Encode()
	return
}

func valid(elem reflect.Value, query *url.Values, key string) {
	switch elem.Kind() {
	case reflect.Ptr:
		valid(elem.Elem(), query, key)
	case reflect.String:
		vString(elem, query, key)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		vInt(elem, query, key)
	case reflect.Float32, reflect.Float64:
		vFloat(elem, query, key)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		vUint(elem, query, key)
	case reflect.Interface:
		valid(elem.Elem(), query, key)
	case reflect.Slice:
		vArr(elem, query, key)
	case reflect.Struct:
		vStruct(elem, query, key)
	case reflect.Map:
		vMap(elem, query, key)
	}
}

func vString(elem reflect.Value, query *url.Values, key string) {
	write(query, key, strings.TrimSpace(elem.String()))
}
func vInt(elem reflect.Value, query *url.Values, key string) {
	write(query, key, strings.TrimSpace(fmt.Sprint(elem.Int())))
}
func vUint(elem reflect.Value, query *url.Values, key string) {
	write(query, key, strings.TrimSpace(fmt.Sprint(elem.Uint())))
}
func vFloat(elem reflect.Value, query *url.Values, key string) {
	write(query, key, strings.TrimSpace(fmt.Sprint(elem.Float())))
}

func vArr(elem reflect.Value, query *url.Values, key string) {
	for i := 0; i < elem.Len(); i++ {
		valid(elem.Index(i), query, key)
	}
}

func vStruct(elem reflect.Value, query *url.Values, key string) {
	for i := 0; i < elem.NumField(); i++ {
		if k := strings.TrimSpace(elem.Type().Field(i).Tag.Get("query")); k != "" {
			if key != "" {
				k = strings.Join([]string{key, k}, "_")
			}
			valid(elem.Field(i), query, k)
		}
	}
}

func vMap(elem reflect.Value, query *url.Values, key string) error {
	maps := elem.MapRange()
	for maps.Next() {
		if k := fmt.Sprint(maps.Key().Interface()); k != "" {
			if key != "" {
				k = strings.Join([]string{key, k}, "_")
			}
			valid(maps.Value(), query, k)
		}
	}
	return nil
}

func write(query *url.Values, key, value string) {
	if value != "" {
		if v, ok := (*query)[key]; ok {
			(*query)[key] = append(v, value)
		} else {
			(*query)[key] = []string{value}
		}
	}
}
