# 03 - 数据结构：切片与 Map

## 1. 数组（Array）

### Java 写法
```java
int[] arr = new int[5];
int[] arr2 = {1, 2, 3};
```

### Go 写法

```go
// 数组长度是类型的一部分！[5]int 和 [3]int 是不同类型
var arr [5]int               // 默认值 [0,0,0,0,0]
arr2 := [3]int{1, 2, 3}
arr3 := [...]int{1, 2, 3, 4} // 长度自动推断为 4

// 访问
fmt.Println(arr2[0])  // 1
arr2[0] = 10

// 遍历
for i, v := range arr2 {
    fmt.Println(i, v)
}
```

### 重要区别

- **Java 数组是对象，长度可变（通过 List）**
- **Go 数组是值类型，长度固定且属于类型的一部分**
- Go 数组赋值会**整体复制**，不是引用

```go
a := [3]int{1, 2, 3}
b := a       // 复制！b 和 a 是独立的
b[0] = 100
fmt.Println(a[0])  // 1（a 没变）
fmt.Println(b[0])  // 100
```

> 实际开发中**很少直接用数组**，几乎都用**切片**。

## 2. 切片（Slice）⭐ 重点

切片是 Go 中最常用的数据结构，相当于 Java 的 `ArrayList`。

### 创建切片

```go
// 方式 1：字面量
s := []int{1, 2, 3}

// 方式 2：make
s1 := make([]int, 5)       // 长度 5，容量 5
s2 := make([]int, 3, 10)   // 长度 3，容量 10

// 方式 3：从数组切
arr := [5]int{1, 2, 3, 4, 5}
s3 := arr[1:4]   // [2, 3, 4]，左闭右开
```

### 基本操作

```go
s := []int{1, 2, 3}

// 添加元素（返回新切片，原切片可能不变）
s = append(s, 4)
s = append(s, 5, 6)
s = append(s, []int{7, 8}...)  // 追加另一个切片

// 长度和容量
fmt.Println(len(s), cap(s))

// 切片操作
sub := s[1:3]      // 索引 1 到 2
first := s[:2]     // 前 2 个
last := s[len(s)-2:] // 后 2 个

// 删除索引 1 的元素（Go 没有内置 remove）
s = append(s[:1], s[2:]...)

// 复制
dst := make([]int, len(s))
copy(dst, s)
```

### 切片是引用类型

```go
a := []int{1, 2, 3}
b := a          // 引用同一底层数组
b[0] = 100
fmt.Println(a[0])  // 100（a 也变了！）
```

> 类比 Java：Go 切片赋值类似 `List<Integer>` 引用赋值，两个变量指向同一份数据。

### 遍历

```go
s := []string{"Go", "Java", "Python"}

// 带索引
for i, v := range s {
    fmt.Printf("%d: %s\n", i, v)
}

// 只要值
for _, v := range s {
    fmt.Println(v)
}

// 只要索引
for i := range s {
    fmt.Println(i)
}
```

## 3. Map

类似 Java 的 `HashMap`。

### 创建 Map

```go
// 方式 1：make
m := make(map[string]int)
m["Go"] = 1
m["Java"] = 2

// 方式 2：字面量
m2 := map[string]int{
    "Go":   1,
    "Java": 2,
}
```

### 基本操作

```go
m := map[string]int{"a": 1, "b": 2}

// 取值
v := m["a"]           // 1

// 判断 key 是否存在（重要！）
v, ok := m["c"]       // ok=false，v=0（零值）
if !ok {
    fmt.Println("key 不存在")
}

// 修改
m["a"] = 100

// 删除
delete(m, "a")

// 遍历（顺序随机！）
for k, v := range m {
    fmt.Println(k, v)
}

// 长度
fmt.Println(len(m))
```

### 注意事项

- Map 是**引用类型**，赋值是共享底层数据
- **遍历顺序是随机的**（和 Java HashMap 一样不保证顺序）
- 取不存在的 key 返回零值，不会 panic，所以要用 `v, ok` 模式判断

## 4. 嵌套结构

```go
// 切片的切片（二维）
matrix := [][]int{
    {1, 2, 3},
    {4, 5, 6},
}

// Map 的值是切片
m := map[string][]int{
    "even": {2, 4, 6},
    "odd":  {1, 3, 5},
}

// Map 的值是 Map
users := map[string]map[string]string{
    "tom": {"age": "18", "city": "北京"},
}
```

## 5. 练习

1. 创建一个切片，添加 5 个元素，遍历打印
2. 创建一个 `map[string]int` 存储水果价格，遍历输出
3. 写一个函数接收切片，返回最大值和最小值
4. 用切片实现一个简单的栈（push/pop）

```go
// 栈示例参考
stack := []int{}
stack = append(stack, 1)        // push
top := stack[len(stack)-1]      // peek
stack = stack[:len(stack)-1]    // pop
```

## 下一节

[04 - 函数与方法](04_函数与方法.md)
