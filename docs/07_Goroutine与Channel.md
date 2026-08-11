# 07 - Goroutine 与 Channel

## 1. 并发 vs 并行

- **并发（Concurrency）**：多个任务交替执行（单核也能并发）
- **并行（Parallelism）**：多个任务同时执行（需要多核）

Go 的并发通过 **goroutine** 实现，比 Java 的 Thread 轻量得多。

## 2. Goroutine

### Java 创建线程
```java
new Thread(() -> System.out.println("hello")).start();
```

### Go 创建 goroutine

```go
// 只需在函数前加 go 关键字
go func() {
    fmt.Println("hello from goroutine")
}()
```

### 简单示例

```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("%s: %d\n", name, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    go sayHello("A")  // 启动 goroutine
    go sayHello("B")
    sayHello("C")     // 主 goroutine

    // 主 goroutine 结束 = 程序结束
    // 所以需要让 main 等一下
}
```

### 主 goroutine 退出会强制结束所有 goroutine

```go
func main() {
    go func() {
        for i := 0; i < 5; i++ {
            fmt.Println(i)
            time.Sleep(time.Second)
        }
    }()

    time.Sleep(2 * time.Second)  // 等 2 秒，主 goroutine 退出后程序结束
    fmt.Println("main 退出")
}
```

> 实际开发中不要用 `time.Sleep` 等待，应该用 `sync.WaitGroup` 或 channel。

## 3. sync.WaitGroup

```go
package main

import (
    "fmt"
    "sync"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()  // 完成时调用 Done

    fmt.Printf("worker %d 开始工作\n", id)
    // 模拟工作
    fmt.Printf("worker %d 完成\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)         // 计数器 +1
        go worker(i, &wg)
    }

    wg.Wait()  // 等待所有 goroutine 完成
    fmt.Println("所有 worker 完成")
}
```

> 类比 Java：相当于 `CountDownLatch`，但用法更简单。

## 4. Channel（通道）⭐

Channel 是 goroutine 之间通信的管道，是 Go 并发的核心。

### 核心理念

> **不要通过共享内存来通信，而要通过通信来共享内存。**

### 创建 channel

```go
// 无缓冲 channel
ch := make(chan int)

// 有缓冲 channel
ch := make(chan int, 10)
```

### 基本操作

```go
ch := make(chan string, 1)

// 发送数据
ch <- "hello"

// 接收数据
msg := <-ch

// 关闭 channel
close(ch)
```

### 无缓冲 channel（同步）

```go
func main() {
    ch := make(chan string)

    go func() {
        ch <- "hello"  // 发送方会阻塞，直到有人接收
        fmt.Println("发送完成")
    }()

    msg := <-ch  // 接收前，发送方一直阻塞
    fmt.Println("收到:", msg)

    time.Sleep(time.Second)
}
```

### 有缓冲 channel（异步）

```go
ch := make(chan int, 3)  // 缓冲区大小 3

ch <- 1  // 不阻塞
ch <- 2  // 不阻塞
ch <- 3  // 不阻塞
// ch <- 4  // 阻塞！缓冲区满了

fmt.Println(<-ch)  // 1
fmt.Println(<-ch)  // 2
```

### 遍历 channel

```go
ch := make(chan int, 5)

go func() {
    for i := 1; i <= 5; i++ {
        ch <- i
    }
    close(ch)  // 发送完必须关闭，否则 for-range 会死锁
}()

for v := range ch {
    fmt.Println(v)
}
```

### 检测 channel 是否关闭

```go
v, ok := <-ch
if !ok {
    fmt.Println("channel 已关闭")
}
```

## 5. select 语句

`select` 用于同时监听多个 channel，类似 Java 的 `select` 或 `switch`。

```go
func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "来自 ch1"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "来自 ch2"
    }()

    // 等先到的
    select {
    case msg := <-ch1:
        fmt.Println(msg)
    case msg := <-ch2:
        fmt.Println(msg)
    }
}
```

### select 的特点

- 多个 case 同时就绪，**随机选一个**
- 没有 case 就绪会阻塞
- 加 `default` 时不阻塞

### 超时控制

```go
select {
case msg := <-ch:
    fmt.Println("收到:", msg)
case <-time.After(3 * time.Second):
    fmt.Println("超时！")
}
```

### 非阻塞接收

```go
select {
case msg := <-ch:
    fmt.Println("收到:", msg)
default:
    fmt.Println("没有数据，跳过")
}
```

## 6. 实战示例：生产者消费者

```go
package main

import (
    "fmt"
    "time"
)

func producer(ch chan<- int) {  // chan<- 表示只发送
    for i := 1; i <= 5; i++ {
        fmt.Println("生产:", i)
        ch <- i
        time.Sleep(500 * time.Millisecond)
    }
    close(ch)
}

func consumer(ch <-chan int) {  // <-chan 表示只接收
    for v := range ch {
        fmt.Println("消费:", v)
    }
}

func main() {
    ch := make(chan int, 3)

    go producer(ch)
    go consumer(ch)

    time.Sleep(3 * time.Second)
    fmt.Println("主程序结束")
}
```

## 7. 并发安全：sync.Mutex

多个 goroutine 修改同一变量时需要加锁。

```go
package main

import (
    "fmt"
    "sync"
)

type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}

func main() {
    counter := &Counter{}
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }

    wg.Wait()
    fmt.Println("最终计数:", counter.Value())  // 1000
}
```

> 类比 Java：`sync.Mutex` 相当于 `synchronized` 或 `ReentrantLock`。

## 8. 常见陷阱

### 1. 循环变量捕获

```go
// 错误写法（Go 1.22 之前）
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // 可能打印 3 个 2！
    }()
}

// 正确写法 1：传参
for i := 0; i < 3; i++ {
    go func(i int) {
        fmt.Println(i)
    }(i)
}

// 正确写法 2：新建变量
for i := 0; i < 3; i++ {
    i := i  // 创建新变量
    go func() {
        fmt.Println(i)
    }()
}
```

> **Go 1.22+ 修复了这个问题**，循环变量每次迭代都是新变量，可以放心用。

### 2. 死锁

```go
func main() {
    ch := make(chan int)
    ch <- 1  // 死锁！无缓冲 channel 没人接收会一直阻塞
}
```

## 9. 练习

1. 启动 5 个 goroutine 同时打印 1~100，用 WaitGroup 等待
2. 用 channel 实现生产者消费者模型
3. 用 `select` + `time.After` 实现一个 3 秒超时
4. 用 `sync.Mutex` 实现一个线程安全的计数器

## 下一节

[08 - 包管理与 Go Modules](08_包管理GoModules.md)
