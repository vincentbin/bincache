package main

import (
	"fmt"
	"main/eliminate"
)

type String string

func (s String) Len() int {
	return len(s)
}

func main() {
	cache := eliminate.New(13, func(s string, v eliminate.Value) {
		fmt.Println(v)
	})

	cache.Add("yyb0", String("123"))
	cache.Add("yyb1", String("123"))
	cache.Add("yyb2", String("123"))
	cache.Add("yyb3", String("123"))
	cache.Add("yyb4", String("123"))

	_, ok := cache.Get("yyb4")
	fmt.Println(ok)

}
