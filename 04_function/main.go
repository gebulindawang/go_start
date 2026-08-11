package main

import "fmt"
type Person struct{
	Name string
	age int
}
func (p Person) sayHello(){
	fmt.Printf("%s今年%d岁了",p.Name,p.age)
}

func main(){
	p := Person{
		Name: "马嘉祺",
		age: 18,
	}
	p.sayHello()
}