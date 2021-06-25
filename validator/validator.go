package validator

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

//Validator виледирует данные
func Validator(isNull bool, data interface{}) error {
	if _, err := router(reflect.ValueOf(data), isNull, ""); err != nil {
		return err
	}
	return nil
}

func router(elem reflect.Value, isNull bool, key string) (out reflect.Value, err error) {
	switch elem.Kind() {
	case reflect.Ptr:
		return router(elem.Elem(), isNull, key)
	case reflect.String:
		return rString(elem, isNull, key)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rInt(elem, isNull, key)
	case reflect.Float32, reflect.Float64:
		return rFloat(elem, isNull, key)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rUint(elem, isNull, key)
	case reflect.Interface:
		return router(elem.Elem(), isNull, key)
	case reflect.Slice:
		return rArr(elem, isNull, key)
	case reflect.Struct:
		return rStruct(elem, isNull, key)
	case reflect.Map:
		return rMap(elem, isNull, key)
	case reflect.Invalid:
		err = fmt.Errorf("%v: is nil", key)
		return
	}
	if elem.IsValid() {
		if t, ok := elem.Interface().(time.Time); ok {
			return rTime(t, isNull, key)
		}
	}
	return
}

func rString(elem reflect.Value, isNull bool, key string) (out reflect.Value, err error) {
	ee := elem
	for ee.Kind() == reflect.Interface {
		ee = ee.Elem()
	}
	ss := strings.TrimSpace(ee.String())
	out = reflect.ValueOf(ss)
	if isNull && ss == "" {
		err = fmt.Errorf("%v: is nil", key)
		return
	}
	if elem.CanSet() {
		elem.Set(out)
	}
	return
}

func rInt(elem reflect.Value, isNull bool, key string) (out reflect.Value, err error) {
	out = elem
	x := elem.Int()
	if isNull {
		if x <= 0 {
			err = fmt.Errorf("%v: is nil", key)
			return
		}
	} else {
		if x < 0 {
			err = fmt.Errorf("%v: is negative", key)
			return
		}
	}
	return
}

func rUint(elem reflect.Value, isNull bool, key string) (out reflect.Value, err error) {
	out = elem
	x := elem.Uint()
	if isNull {
		if x == 0 {
			err = fmt.Errorf("%v: is nil", key)
			return
		}
	}
	return
}

func rFloat(elem reflect.Value, isNull bool, key string) (out reflect.Value, err error) {
	out = elem
	x := elem.Float()
	if isNull {
		if x <= 0 {
			err = fmt.Errorf("%v: is nil", key)
			return
		}
	} else {
		if x < 0 {
			err = fmt.Errorf("%v: is negative", key)
			return
		}
	}
	return
}

func rTime(elem time.Time, isNull bool, key string) (out reflect.Value, err error) {
	out = reflect.ValueOf(elem)
	if isNull {
		if elem.IsZero() {
			err = fmt.Errorf("%v: is nil", key)
			return
		}
	}
	return
}

func rArr(elem reflect.Value, isNull bool, key string) (out reflect.Value, err error) {
	out = elem
	for i := 0; i < elem.Len(); i++ {
		k := fmt.Sprintf("[%v]", i)
		if key != "" {
			k = fmt.Sprintf("%v[%v]", key, i)
		}
		value, err := router(elem.Index(i), isNull, k)
		if err != nil {
			return out, err
		}
		if !elem.Index(i).CanSet() || elem.Index(i).IsZero() {
			continue
		}
		if elem.Index(i).Kind() == reflect.Ptr {
			value = value.Addr()
		}
		if value.Type().ConvertibleTo(elem.Index(i).Type()) {
			value = value.Convert(elem.Index(i).Type())
		}
		elem.Index(i).Set(value)
	}
	return
}

func rStruct(elem reflect.Value, isNull bool, key string) (out reflect.Value, err error) {
	out = elem
	for i := 0; i < elem.NumField(); i++ {
		if k := strings.TrimSpace(elem.Type().Field(i).Tag.Get("valid")); k != "" {
			if key != "" {
				k = fmt.Sprintf("%v.%v", key, k)
			}
			value, err := router(elem.Field(i), isNull, k)
			if err != nil {
				return out, err
			}
			if !elem.Field(i).CanSet() || elem.Field(i).IsZero() {
				continue
			}
			if elem.Field(i).Kind() == reflect.Ptr {
				value = value.Addr()
			}
			if value.Type().ConvertibleTo(elem.Field(i).Type()) {
				value = value.Convert(elem.Field(i).Type())
			}
			elem.Field(i).Set(value)
		}
	}
	return
}

func rMap(elem reflect.Value, isNull bool, key string) (out reflect.Value, err error) {
	out = elem
	maps := elem.MapRange()
	for maps.Next() {
		k := maps.Key().String()
		if key != "" {
			k = fmt.Sprintf("%v.%v", key, maps.Key().String())
		}
		value, err := router(maps.Value(), isNull, k)
		if err != nil {
			return out, err
		}
		if maps.Value().Kind() == reflect.Ptr {
			value = value.Addr()
		}
		if value.Type().ConvertibleTo(maps.Value().Type()) {
			value = value.Convert(maps.Value().Type())
		}
		elem.SetMapIndex(maps.Key(), value)
	}
	return
}
