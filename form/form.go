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

type F map[string]interface{}

func FormMarshal(in interface{}) (*bytes.Buffer, string, error) {
	out := &bytes.Buffer{}
	form := multipart.NewWriter(out)
	if err := valid(reflect.ValueOf(in), form, ""); err != nil {
		return nil, "", err
	}
	err := form.Close()
	return out, form.FormDataContentType(), err
}

func valid(elem reflect.Value, form *multipart.Writer, key string) (err error) {
	switch elem.Kind() {
	case reflect.Ptr:
		valid(elem.Elem(), form, key)
	case reflect.String:
		vString(elem, form, key)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		vInt(elem, form, key)
	case reflect.Float32, reflect.Float64:
		vFloat(elem, form, key)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		vUint(elem, form, key)
	case reflect.Interface:
		valid(elem.Elem(), form, key)
	case reflect.Slice:
		vArr(elem, form, key)
	case reflect.Struct:
		vStruct(elem, form, key)
	case reflect.Map:
		vMap(elem, form, key)
	}
	if elem.IsValid() {
		if t, ok := elem.Interface().(time.Time); ok {
			return write(form, key, strings.TrimSpace(t.Format(time.RFC3339)))
		}
	}
	return
}

func vString(elem reflect.Value, form *multipart.Writer, key string) error {
	return write(form, key, strings.TrimSpace(elem.String()))
}
func vInt(elem reflect.Value, form *multipart.Writer, key string) error {
	return write(form, key, strings.TrimSpace(fmt.Sprint(elem.Int())))
}
func vUint(elem reflect.Value, form *multipart.Writer, key string) error {
	return write(form, key, strings.TrimSpace(fmt.Sprint(elem.Uint())))
}
func vFloat(elem reflect.Value, form *multipart.Writer, key string) error {
	return write(form, key, strings.TrimSpace(fmt.Sprint(elem.Float())))
}

func vArr(elem reflect.Value, form *multipart.Writer, key string) (err error) {
	for i := 0; i < elem.Len(); i++ {
		if err = valid(elem.Index(i), form, key); err != nil {
			return fmt.Errorf("%q [%v] write failed: %v", key, i, err)
		}
	}
	return
}

func vStruct(elem reflect.Value, form *multipart.Writer, key string) (err error) {
	for i := 0; i < elem.NumField(); i++ {
		if k := strings.TrimSpace(elem.Type().Field(i).Tag.Get("form")); k != "" {
			if key != "" {
				k = strings.Join([]string{key, k}, "_")
			}
			if err = valid(elem.Field(i), form, k); err != nil {
				return fmt.Errorf("%q write failed: %v", k, err)
			}
		}
	}
	return
}

func vMap(elem reflect.Value, form *multipart.Writer, key string) (err error) {
	maps := elem.MapRange()
	for maps.Next() {
		if k := fmt.Sprint(maps.Key().Interface()); k != "" {
			if key != "" {
				k = strings.Join([]string{key, k}, "_")
			}
			if err = valid(maps.Value(), form, k); err != nil {
				return fmt.Errorf("%q write failed: %v", k, err)
			}
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
		_, err = field.Write([]byte(value))
	}
	return
}
