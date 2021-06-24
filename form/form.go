package form

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
	"strings"
	"time"
)

//Генерирует mutlipart/form и content-type
func FormMarshal(in interface{}) (*bytes.Buffer, string, error) {
	out := &bytes.Buffer{}
	form := multipart.NewWriter(out)
	if err := router(reflect.ValueOf(in), form, ""); err != nil {
		return nil, "", err
	}
	err := form.Close()
	return out, form.FormDataContentType(), err
}

func router(elem reflect.Value, form *multipart.Writer, key string) (err error) {
	switch elem.Kind() {
	case reflect.Ptr:
		return router(elem.Elem(), form, key)
	case reflect.String:
		return rString(elem, form, key)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rInt(elem, form, key)
	case reflect.Float32, reflect.Float64:
		return rFloat(elem, form, key)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rUint(elem, form, key)
	case reflect.Interface:
		return router(elem.Elem(), form, key)
	case reflect.Slice:
		return rArr(elem, form, key)
	case reflect.Struct:
		return rStruct(elem, form, key)
	case reflect.Map:
		return rMap(elem, form, key)
	}
	if elem.IsValid() {
		if t, ok := elem.Interface().(time.Time); ok {
			return write(form, key, strings.TrimSpace(t.Format(time.RFC3339)))
		}
	}
	return
}

func rString(elem reflect.Value, form *multipart.Writer, key string) error {
	return write(form, key, strings.TrimSpace(elem.String()))
}
func rInt(elem reflect.Value, form *multipart.Writer, key string) error {
	return write(form, key, strings.TrimSpace(fmt.Sprint(elem.Int())))
}
func rUint(elem reflect.Value, form *multipart.Writer, key string) error {
	return write(form, key, strings.TrimSpace(fmt.Sprint(elem.Uint())))
}
func rFloat(elem reflect.Value, form *multipart.Writer, key string) error {
	return write(form, key, strings.TrimSpace(fmt.Sprint(elem.Float())))
}

func rArr(elem reflect.Value, form *multipart.Writer, key string) (err error) {
	for i := 0; i < elem.Len(); i++ {
		if err = router(elem.Index(i), form, key); err != nil {
			return
		}
	}
	return
}

func rStruct(elem reflect.Value, form *multipart.Writer, key string) (err error) {
	for i := 0; i < elem.NumField(); i++ {
		if k := strings.TrimSpace(elem.Type().Field(i).Tag.Get("form")); k != "" {
			if key != "" {
				k = strings.Join([]string{key, k}, ".")
			}
			if err = router(elem.Field(i), form, k); err != nil {
				return
			}
		}
	}
	return
}

func rMap(elem reflect.Value, form *multipart.Writer, key string) (err error) {
	maps := elem.MapRange()
	for maps.Next() {
		k := maps.Key().String()
		if key != "" {
			k = strings.Join([]string{key, k}, ".")
		}
		if err = router(maps.Value(), form, k); err != nil {
			return
		}
	}
	return
}

func write(form *multipart.Writer, key, value string) (err error) {
	if value != "" {
		var field io.Writer
		field, err = form.CreateFormField(key)
		if err != nil {
			return
		}
		if _, err = field.Write([]byte(value)); err != nil {
			err = fmt.Errorf("%v write error: %v", key, err)
		}
	}
	return
}
