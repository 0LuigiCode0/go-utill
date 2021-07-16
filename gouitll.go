package goutill

import (
	"bytes"
	"io"

	"github.com/0LuigiCode0/go-utill/form"
	"github.com/0LuigiCode0/go-utill/helper"
	"github.com/0LuigiCode0/go-utill/query"
	"github.com/0LuigiCode0/go-utill/validator"
)

//Генерирует mutlipart/form и content-type
func FormMarshal(in interface{}) (*bytes.Buffer, string, error) { return form.FormMarshal(in) }

//Генерирует строку query параметров
func QueryMarshal(in interface{}) (out string) { return query.QueryMarshal(in) }

//Validator виледирует данные
func Validator(isNull bool, data interface{}) error { return validator.Validator(isNull, data) }

//Парсит Json
func JsonParse(in io.Reader, out interface{}) (err error) { return helper.JsonParse(in, out) }
