# 02 - Go 基础语法（Java 视角）

## 1. 变量声明

### Java 写法
```java
int a = 10;
String name = "Tom";
final int MAX = 100;
```

### Go 写法

```go
// 方式 1：var 关键字
var a int = 10
var name string = "Tom"

// 方式 2：类型推断（推荐）
var a = 10
var name = "Tom"

// 方式 3：短变量声明（只能在函数内用，最常用）
b := 20
age := 18

// 多变量声明
var x, y int = 1, 2
var m, n = "hi", 3.14
p, q := 1, 2
```

### 区别要点

| 特性 | Java | Go |
|------|------|-----|
| 声明位置 | 类内/方法内 | 包级/函数内 |
| 未初始化默认值 | null / 0 | 0 / "" / false / nil |
| 短声明 `:=` | 不支持 | 函数内可用 |
| 多返回值 | 不支持 | 支持 |

## 2. 基本数据类型

```go
// 整数
var i int = 10           // 平台相关（64位机器上是 int64）
var i8 int8 = 127        // -128 ~ 127
var ui uint = 10         // 无符号

// 浮点
var f float64 = 3.14

// 布尔
var ok bool = true

// 字符串
var s string = "hello"

// 字符（rune，本质是 int32）
var c rune = '中'
```

### 类型转换

Go **没有隐式类型转换**，必须显式：

```go
var a int = 10
var b float64 = float64(a)   // 必须显式
var c int = int(b)
```

> Java 中 `int -> long -> double` 会自动转换，Go 中**任何**类型转换都必须写出来。

### 字符串

```go
s := "hello"
// 字符串拼接
s2 := s + " world"
// 多行字符串（反引号，类似 Java 的 textBlock）
s3 := `第一行
第二行
第三行`
// 字符串长度（字节数）
n := len(s)
// 取字符（按字节）
ch := s[0]   // 'h' 的 byte 值
```

## 3. 常量

```go
// 单个常量
const Pi = 3.14159
const MaxSize int = 100

// 多常量一起声明
const (
    StatusOK = 200
    StatusNotFound = 404
    StatusError = 500
)

// iota：自增常量（类似枚举）
const (
    Sunday = iota    // 0
    Monday           // 1
    Tuesday          // 2
    Wednesday        // 3
)
```

> Java 没有原生的 `iota`，Go 用它模拟枚举。

## 4. 控制流

### if-else

```go
age := 18

// 写法 1：常规
if age >= 18 {
    fmt.Println("成年")
} else if age >= 12 {
    fmt.Println("少年")
} else {
    fmt.Println("儿童")
}

// 写法 2：带初始化语句（Go 特色）
if n := len("hello"); n > 3 {
    fmt.Println("长度大于 3:", n)
}
// n 的作用域仅限于 if-else 块内
```

> **注意**：条件表达式**不需要括号**，但花括号是必须的，且 `{` 必须在同一行。

### for 循环（Go 只有 for，没有 while）

```go
// 写法 1：经典 for
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

// 写法 2：类似 while
n := 0
for n < 5 {
    fmt.Println(n)
    n++
}

// 写法 3：死循环
for {
    fmt.Println("running...")
    break
}

// 写法 4：遍历（类似 Java 的 for-each）
nums := []int{1, 2, 3}
for index, value := range nums {
    fmt.Println(index, value)
}

// 忽略 index
for _, value := range nums {
    fmt.Println(value)
}
```

### switch

```go
day := "Monday"

// 经典 switch（Go 默认不穿透，不需要 break！）
switch day {
case "Monday", "Tuesday":
    fmt.Println("工作日")
case "Saturday", "Sunday":
    fmt.Println("周末")
default:
    fmt.Println("其他")
}

// 带初始化的 switch
switch n := 10; {
case n < 0:
    fmt.Println("负数")
case n == 0:
    fmt.Println("零")
default:
    fmt.Println("正数")
}

// 使用 fallthrough 强制穿透（少用）
switch x := 1; x {
case 1:
    fmt.Println("一")
    fallthrough
case 2:
    fmt.Println("二")
}
```

> **重要区别**：Java switch 需要 break 防止穿透，**Go 默认不穿透**，需要穿透才用 `fallthrough`。

## 5. 运算符

大部分和 Java 一样，但有几个差异：

| 运算符 | Java | Go | 说明 |
|--------|------|-----|------|
| 自增 | `i++` (可用在表达式中) | `i++` (只能作为语句) | Go 中 `j = i++` 会报错 |
| 逻辑非 | `!` | `!` | 相同 |
| 位运算 | `& | ^` | `& | ^` | 相同 |
| 取地址 | 不常用 | `&` | Go 中常用（指针） |
| 解引用 | 不常用 | `*` | Go 中常用（指针） |

```go
i := 5
i++         // OK
// j := i++ // 报错！Go 中 ++ 是语句不是表达式
```

## 6. 指针（Java 没有的概念）

Go 有指针，但比 C 简单：**没有指针运算**。

```go
a := 10
p := &a          // p 是指向 a 的指针，类型 *int
fmt.Println(p)   // 输出地址，如 0xc0000b2000
fmt.Println(*p)  // 解引用，输出 10

*p = 20          // 通过指针修改 a 的值
fmt.Println(a)   // 20
```

> Java 中所有对象引用其实就是指针，只不过 Java 把它隐藏了。Go 把指针显式化，但用法受限（更安全）。
>
> 日常使用中，**指针主要用于函数传参时避免拷贝或允许修改**。

## 7. 输入输出

```go
package main

import "fmt"

func main() {
    // 输出
    fmt.Print("不带换行")
    fmt.Println("带换行")
    fmt.Printf("格式化: %s 今年 %d 岁\n", "Tom", 18)

    // 输入
    var name string
    fmt.Print("请输入姓名: ")
    fmt.Scan(&name)
    fmt.Println("你好,", name)
}
```

### 常用格式化占位符

| 占位符 | 说明 | 示例 |
|--------|------|------|
| `%d` | 整数 | `fmt.Printf("%d", 10)` |
| `%s` | 字符串 | `fmt.Printf("%s", "hi")` |
| `%f` | 浮点 | `fmt.Printf("%.2f", 3.14)` |
| `%t` | 布尔 | `fmt.Printf("%t", true)` |
| `%v` | 通用（任意类型） | `fmt.Printf("%v", anything)` |
| `%+v` | 带字段名 | 输出 struct 时显示字段名 |
| `%T` | 类型 | `fmt.Printf("%T", x)` |

## 8. 练习

1. 声明一个字符串变量，用 `fmt.Printf` 输出其类型和值
2. 写一个 for 循环打印 1~10
3. 用 switch 实现一个简单的成绩等级判断（A/B/C/D）
4. 尝试用指针修改一个变量的值

## 下一节

[03 - 数据结构：切片与 Map](03_数据结构_切片与Map.md)
