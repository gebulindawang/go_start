package main

import "fmt"

func main() {
	m := map[int]string{1:"java",2:"go"}
	v := m[1]
	fmt.Println(v)
	v,ok := m[3]
	if !ok{
		fmt.Println(v)
		fmt.Println("key不存在")
	}
}
