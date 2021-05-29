package helper

import (
	"encoding/json"
	"io"
)

//Парсит Json
func JsonParse(in io.Reader, out interface{}) (err error) {
	buf, err := io.ReadAll(in)
	if err != nil {
		return
	}
	return json.Unmarshal(buf, out)
}
