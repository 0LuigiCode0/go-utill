package helper

import (
	"fmt"
	"reflect"
	"strings"
)

type V map[string]interface{}

//Validator полномосштабная валидация
func Validator(isNull bool, data V) error {
	if err := vMap(reflect.ValueOf(data), isNull); err != nil {
		return err
	}
	return nil
}

func valid(elem reflect.Value, isNull bool, k string) error {
	switch kind := elem.Kind(); {
	case kind == reflect.Ptr:
		if err := valid(elem.Elem(), isNull, k); err != nil {
			return fmt.Errorf("%v: %v", k, err)
		}
	case kind == reflect.String:
		if err := vString(elem, isNull); err != nil {
			return fmt.Errorf("%v: %v", k, err)
		}
	case kind == reflect.Int || kind == reflect.Int64:
		if err := vInt(elem, isNull); err != nil {
			return fmt.Errorf("%v: %v", k, err)
		}
	case kind == reflect.Float32 || kind == reflect.Float64:
		if err := vFloat(elem, isNull); err != nil {
			return fmt.Errorf("%v: %v", k, err)
		}
	case kind == reflect.Interface:
		if err := valid(elem.Elem(), isNull, k); err != nil {
			return fmt.Errorf("%v: %v", k, err)
		}
	case kind == reflect.Slice:
		if err := vArr(elem, isNull); err != nil {
			return fmt.Errorf("%v: %v", k, err)
		}
	case kind == reflect.Struct:
		if err := vStruct(elem, isNull); err != nil {
			return fmt.Errorf("%v: %v", k, err)
		}
		// case kind == reflect.Map:
		// 	if err := vMap(elem, isNull); err != nil {
		// 		return fmt.Errorf("%v: %v", k, err)
		// 	}
	}
	return nil
}

func vString(elem reflect.Value, isNull bool) error {
	// fmt.Println(elem.Type().)
	if !elem.CanSet() {
		return fmt.Errorf("cannot set")
	}
	ee := elem
	for ee.Kind() == reflect.Interface {
		ee = ee.Elem()
	}
	ss := strings.TrimSpace(ee.String())
	if isNull && ss == "" {
		return fmt.Errorf("is nil")
	}
	elem.Set(reflect.ValueOf(ss))
	return nil
}

func vInt(elem reflect.Value, isNull bool) error {
	x := elem.Int()
	if isNull {
		if x <= 0 {
			return fmt.Errorf("is nil")
		}
	} else {
		if x < 0 {
			return fmt.Errorf("is negative")
		}
	}
	return nil
}

func vFloat(elem reflect.Value, isNull bool) error {
	x := elem.Float()
	if isNull {
		if x <= 0 {
			return fmt.Errorf("is nil")
		}
	} else {
		if x < 0 {
			return fmt.Errorf("is negative")
		}
	}
	return nil
}

func vArr(elem reflect.Value, isNull bool) error {
	for i := 0; i < elem.Len(); i++ {
		if err := valid(elem.Index(i), isNull, fmt.Sprint(i)); err != nil {
			return fmt.Errorf("[%v]: %v", i, err)
		}
	}
	return nil
}

func vStruct(elem reflect.Value, isNull bool) error {
	for i := 0; i < elem.NumField(); i++ {
		if k := strings.TrimSpace(elem.Type().Field(i).Tag.Get("valid")); k != "" {
			if err := valid(elem.Field(i), isNull, k); err != nil {
				return err
			}
		}
	}
	return nil
}

func vMap(elem reflect.Value, isNull bool) error {
	maps := elem.MapRange()
	for maps.Next() {
		if err := valid(maps.Value(), isNull, maps.Key().String()); err != nil {
			return fmt.Errorf("[%v]: %v", maps.Key().String(), err)
		}
	}
	return nil
}
