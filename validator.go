package helper

import (
	"fmt"
	"reflect"
	"strings"
)

type V map[string]interface{}

//Validator полномосштабная валидация
func Validator(isnull bool, data V) error {
	for k, v := range data {
		elem := reflect.ValueOf(v).Elem()
		if err := valid(elem, isnull, k); err != nil {
			return err
		}
	}
	return nil
}

func valid(elem reflect.Value, isnull bool, k string) error {
	switch kind := elem.Kind(); {
	case kind == reflect.String:
		if err := vString(elem, isnull); err != nil {
			return fmt.Errorf("%v: %v", k, err)
		}
	case kind == reflect.Int || kind == reflect.Int64:
		if err := vInt(elem, isnull); err != nil {
			return fmt.Errorf("%v: %v", k, err)
		}
	case kind == reflect.Float32 || kind == reflect.Float64:
		if err := vFloat(elem, isnull); err != nil {
			return fmt.Errorf("%v: %v", k, err)
		}
	case kind == reflect.Slice:
		if err := vArr(elem, isnull); err != nil {
			return fmt.Errorf("%v: %v", k, err)
		}
	case kind == reflect.Struct:
		if err := vStruct(elem, isnull); err != nil {
			return fmt.Errorf("%v: %v", k, err)
		}
	case kind == reflect.Ptr:
		if err := valid(elem.Elem(), isnull, k); err != nil {
			return fmt.Errorf("%v: %v", k, err)
		}
	}
	return nil
}

func vString(elem reflect.Value, isnull bool) error {
	ee := elem
	for ee.Kind() == reflect.Interface {
		ee = ee.Elem()
	}
	ss := strings.TrimSpace(ee.String())
	if isnull && ss == "" {
		return fmt.Errorf("is nil")
	}
	elem.Set(reflect.ValueOf(ss))
	return nil
}

func vInt(elem reflect.Value, isnull bool) error {
	x := elem.Int()
	if isnull {
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

func vFloat(elem reflect.Value, isnull bool) error {
	x := elem.Float()
	if isnull {
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

func vArr(elem reflect.Value, isnull bool) error {
	switch kind := elem.Type().Elem().Kind(); {
	case kind == reflect.String:
		for i := 0; i < elem.Len(); i++ {
			if err := vString(elem.Index(i), isnull); err != nil {
				return fmt.Errorf("[%v]: %v", i, err)
			}
		}
	case kind == reflect.Int || kind == reflect.Int64:
		for i := 0; i < elem.Len(); i++ {
			if err := vInt(elem.Index(i), isnull); err != nil {
				return fmt.Errorf("[%v]: %v", i, err)
			}
		}
	case kind == reflect.Float32 || kind == reflect.Float64:
		for i := 0; i < elem.Len(); i++ {
			if err := vFloat(elem.Index(i), isnull); err != nil {
				return fmt.Errorf("[%v]: %v", i, err)
			}
		}
	case kind == reflect.Interface:
		for i := 0; i < elem.Len(); i++ {
			e := elem.Index(i)
			switch kind := e.Elem().Kind(); {
			case kind == reflect.String:
				if err := vString(e, isnull); err != nil {
					return fmt.Errorf("[%v]: %v", i, err)
				}
			case kind == reflect.Int || kind == reflect.Int64:
				if err := vInt(elem.Index(i).Elem(), isnull); err != nil {
					return fmt.Errorf("[%v]: %v", i, err)
				}
			case kind == reflect.Float32 || kind == reflect.Float64:
				if err := vFloat(elem.Index(i).Elem(), isnull); err != nil {
					return fmt.Errorf("[%v]: %v", i, err)
				}
			case kind == reflect.Slice:
				if err := vArr(elem.Index(i).Elem(), isnull); err != nil {
					return fmt.Errorf("[%v]: %v", i, err)
				}
			}
		}
	case kind == reflect.Slice:
		for i := 0; i < elem.Len(); i++ {
			if err := vArr(elem.Index(i), isnull); err != nil {
				return fmt.Errorf("[%v]: %v", i, err)
			}
		}
	}
	return nil
}

func vStruct(elem reflect.Value, isnull bool) error {
	for i := 0; i < elem.NumField(); i++ {
		if k := strings.TrimSpace(elem.Type().Field(i).Tag.Get("valid")); k != "" {
			if err := valid(elem.Field(i), isnull, k); err != nil {
				return err
			}
		}
	}
	return nil
}
