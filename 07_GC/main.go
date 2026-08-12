package main

import (
	"fmt"
	"time"
)

func sayHello(name string) {
	for i := range 3 {
		fmt.Printf("%s:%d\n", name, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	ch := make(chan string)

	go func() {
		ch <- "hello" // 发送方会阻塞，直到有人接收
		fmt.Println("发送完成")
	}()

	msg := <-ch // 接收前，发送方一直阻塞
	fmt.Println("收到:", msg)

	time.Sleep(time.Second)
}
