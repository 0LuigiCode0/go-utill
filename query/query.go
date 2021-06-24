package query

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"
)

//Генерирует строку query параметров
func QueryMarshal(in interface{}) (out string) {
	query := url.Values{}
	router(reflect.ValueOf(in), &query, "")
	out = query.Encode()
	return
}

func router(elem reflect.Value, query *url.Values, key string) {
	switch elem.Kind() {
	case reflect.Ptr:
		router(elem.Elem(), query, key)
	case reflect.String:
		rString(elem, query, key)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rInt(elem, query, key)
	case reflect.Float32, reflect.Float64:
		rFloat(elem, query, key)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rUint(elem, query, key)
	case reflect.Interface:
		router(elem.Elem(), query, key)
	case reflect.Slice:
		rArr(elem, query, key)
	case reflect.Struct:
		rStruct(elem, query, key)
	case reflect.Map:
		rMap(elem, query, key)
	}
	if elem.IsValid() {
		if t, ok := elem.Interface().(time.Time); ok {
			write(query, key, strings.TrimSpace(t.Format(time.RFC3339)))
		}
	}
}

func rString(elem reflect.Value, query *url.Values, key string) {
	write(query, key, strings.TrimSpace(elem.String()))
}
func rInt(elem reflect.Value, query *url.Values, key string) {
	write(query, key, strings.TrimSpace(fmt.Sprint(elem.Int())))
}
func rUint(elem reflect.Value, query *url.Values, key string) {
	write(query, key, strings.TrimSpace(fmt.Sprint(elem.Uint())))
}
func rFloat(elem reflect.Value, query *url.Values, key string) {
	write(query, key, strings.TrimSpace(fmt.Sprint(elem.Float())))
}

func rArr(elem reflect.Value, query *url.Values, key string) {
	for i := 0; i < elem.Len(); i++ {
		router(elem.Index(i), query, key)
	}
}

func rStruct(elem reflect.Value, query *url.Values, key string) {
	for i := 0; i < elem.NumField(); i++ {
		if k := strings.TrimSpace(elem.Type().Field(i).Tag.Get("query")); k != "" {
			if key != "" {
				k = strings.Join([]string{key, k}, ".")
			}
			router(elem.Field(i), query, k)
		}
	}
}

func rMap(elem reflect.Value, query *url.Values, key string) {
	maps := elem.MapRange()
	for maps.Next() {
		if k := fmt.Sprint(maps.Key().Interface()); k != "" {
			if key != "" {
				k = strings.Join([]string{key, k}, ".")
			}
			router(maps.Value(), query, k)
		}
	}
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
