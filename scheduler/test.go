package main

import "fmt"

type mess struct {
	name int
}
func main() {
	var smap = make(map[string]mess)

	smap["hello"] = mess {name: 12}
	if _, ok := smap["ddd"]; ok {
		fmt.Println(smap["ddd"])
	}
	if _, ok := smap["hello"]; ok {
		fmt.Println(smap["hello"])
	}
}
