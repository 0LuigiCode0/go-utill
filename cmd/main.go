package main

import (
	"fmt"

	"github.com/0LuigiCode0/go-utill/net"
)

type Test struct {
	In     *Test               `query:"in"`
	String string              `query:"string"`
	Int    int8                `query:"int"`
	Flaot  float32             `query:"float"`
	Uint   uint32              `query:"uint"`
	Map    map[int]interface{} `query:"map"`
	Arr    []int               `query:"arr"`
}

func main() {
	t := Test{
		In: &Test{
			String: "sdsd",
			Int:    4,
			Map:    map[int]interface{}{1: "sdsd", 2: "sdsd"},
		},
		Flaot: 5.5,
		Uint:  45,
		Arr:   []int{4, 8},
	}

	s := net.QueryMarshal(t)
	fmt.Println(s)

	s = net.QueryMarshal(map[string][]string{"1": {"sdfdsf", "sfssdfd"}, "4": {"tthyrty", "xcvxc"}})
	fmt.Println(s)
}
