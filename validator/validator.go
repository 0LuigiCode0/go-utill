package validator

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type V map[string]interface{}

//Validator полномосштабная валидация
func Validator(isNull bool, data V) error {
	if err := vMap(reflect.ValueOf(data), isNull); err != nil {
		return err
	}
	return nil
}

func valid(elem reflect.Value, isNull bool) (out reflect.Value, err error) {
	switch elem.Kind() {
	case reflect.Ptr:
		out, err = valid(elem.Elem(), isNull)
		if err != nil {
			return out, err
		}
	case reflect.String:
		out, err = vString(elem, isNull)
		if err != nil {
			return out, err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if err := vInt(elem, isNull); err != nil {
			return out, err
		}
	case reflect.Float32, reflect.Float64:
		if err := vFloat(elem, isNull); err != nil {
			return out, err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if err := vUint(elem, isNull); err != nil {
			return out, err
		}
	case reflect.Interface:
		out, err = valid(elem.Elem(), isNull)
		if err != nil {
			return out, err
		}
	case reflect.Slice:
		if err := vArr(elem, isNull); err != nil {
			return out, err
		}
	case reflect.Struct:
		if err := vStruct(elem, isNull); err != nil {
			return out, err
		}
	case reflect.Map:
		if err := vMap(elem, isNull); err != nil {
			return out, err
		}
	}
	if elem.IsValid() {
		if t, ok := elem.Interface().(time.Time); ok {
			if err := vTime(t, isNull); err != nil {
				return out, err
			}
		}
	}
	return out, nil
}

func vString(elem reflect.Value, isNull bool) (out reflect.Value, err error) {
	ee := elem
	for ee.Kind() == reflect.Interface {
		ee = ee.Elem()
	}
	ss := strings.TrimSpace(ee.String())
	if isNull && ss == "" {
		return out, fmt.Errorf("is nil")
	}
	out = reflect.ValueOf(ss)
	if !elem.CanSet() {
		return out, nil
	}
	elem.Set(out)
	return out, nil
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

func vUint(elem reflect.Value, isNull bool) error {
	x := elem.Uint()
	if isNull {
		if x == 0 {
			return fmt.Errorf("is nil")
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

func vTime(elem time.Time, isNull bool) error {
	if isNull {
		if elem.IsZero() {
			return fmt.Errorf("is nil")
		}
	}
	return nil
}

func vArr(elem reflect.Value, isNull bool) error {
	for i := 0; i < elem.Len(); i++ {
		if _, err := valid(elem.Index(i), isNull); err != nil {
			return fmt.Errorf("index [%v]: %v", i, err)
		}
	}
	return nil
}

func vStruct(elem reflect.Value, isNull bool) error {
	for i := 0; i < elem.NumField(); i++ {
		if k := strings.TrimSpace(elem.Type().Field(i).Tag.Get("valid")); k != "" {
			if _, err := valid(elem.Field(i), isNull); err != nil {
				return fmt.Errorf("tag %q: %v", k, err)
			}
		}
	}
	return nil
}

func vMap(elem reflect.Value, isNull bool) error {
	maps := elem.MapRange()
	for maps.Next() {
		ee, err := valid(maps.Value(), isNull)
		if err != nil {
			return fmt.Errorf("key %q: %v", maps.Key().String(), err)
		}
		elem.SetMapIndex(maps.Key(), ee)
	}
	return nil
}
