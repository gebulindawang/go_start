# 08 - 包管理与 Go Modules

## 1. Go Modules 概述

Go Modules 是 Go 1.11+ 引入的依赖管理系统，相当于 Java 的 Maven/Gradle。

### 核心文件

- `go.mod`：模块定义和依赖声明（类似 `pom.xml`）
- `go.sum`：依赖的哈希校验（自动生成，不要手动改）

## 2. 初始化模块

```bash
mkdir myproject
cd myproject
go mod init myproject
```

生成的 `go.mod`：

```
module myproject

go 1.22
```

### 模块命名约定

- **个人项目**：`github.com/username/project`
- **本地学习**：直接用项目名即可，如 `go_start`

## 3. 项目结构

### 推荐结构

```
myproject/
├── go.mod
├── go.sum
├── main.go              # 程序入口
├── config/
│   └── config.go        # 子包：配置
├── model/
│   └── user.go          # 子包：用户模型
├── service/
│   └── user_service.go  # 子包：用户服务
└── utils/
    └── string_utils.go  # 子包：工具
```

## 4. 包的导入

### 创建子包

```go
// model/user.go
package model

type User struct {
    ID   int
    Name string
}

func NewUser(id int, name string) *User {
    return &User{ID: id, Name: name}
}
```

### 导入使用

```go
// main.go
package main

import (
    "fmt"
    "myproject/model"   // 导入本地包：模块名 + 路径
)

func main() {
    u := model.NewUser(1, "Tom")
    fmt.Println(u.Name)
}
```

### 导入别名

```go
import (
    m "myproject/model"     // 别名 m
    _ "myproject/init"      // 只执行 init，不导入名字
    . "fmt"                 // 直接导入到当前命名空间（不推荐）
)

func main() {
    u := m.NewUser(1, "Tom")
    Println(u.Name)  // 不需要 fmt.
}
```

## 5. 包的可见性

```go
// model/user.go
package model

type User struct {
    ID   int       // 大写：公开
    Name string    // 大写：公开
    age  int       // 小写：私有（包外不可见）
}

func NewUser(id int, name string, age int) *User {
    return &User{ID: id, Name: name, age: age}  // 包内可访问 age
}

func (u *User) GetAge() int {  // 大写：公开
    return u.age
}

func helper() {  // 小写：私有
    // ...
}
```

> **重要**：可见性是**包级别**的，不是类型级别。
> 同一个包内的所有文件可以互相访问私有成员。

## 6. 添加第三方依赖

```bash
# 安装包
go get github.com/gin-gonic/gin

# 安装指定版本
go get github.com/gin-gonic/gin@v1.9.1

# 升级
go get -u github.com/gin-gonic/gin

# 清理未使用的依赖
go mod tidy

# 下载依赖到本地缓存
go mod download
```

`go.mod` 会自动更新：

```
module myproject

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/mattn/go-sqlite3 v1.14.17
)
```

## 7. init 函数执行顺序

```go
// pkgA/a.go
package pkgA

import "fmt"

func init() {
    fmt.Println("pkgA init")
}
```

执行顺序：

1. 导入的包先初始化（按导入顺序，深度优先）
2. 包内变量初始化
3. 包内 `init()` 函数（按文件名顺序，可多个）
4. `main()` 函数

```
main 导入 pkgA -> pkgA 的全局变量 -> pkgA.init() -> main 的全局变量 -> main.init() -> main.main()
```

## 8. 常用命令

| 命令 | 说明 |
|------|------|
| `go mod init <name>` | 初始化模块 |
| `go mod tidy` | 添加缺失的依赖、删除多余的 |
| `go mod download` | 下载依赖到本地缓存 |
| `go mod graph` | 查看依赖图 |
| `go mod vendor` | 将依赖复制到 vendor 目录 |
| `go get <pkg>` | 添加/更新依赖 |
| `go get -u <pkg>` | 升级到最新版本 |
| `go get <pkg>@version` | 安装指定版本 |
| `go build` | 编译 |
| `go run` | 直接运行 |
| `go test` | 运行测试 |
| `go fmt` | 格式化代码 |
| `go vet` | 静态检查 |

## 9. go.mod 实例

```
module github.com/myname/myapp

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/jinzhu/gorm v1.9.16
    github.com/spf13/viper v1.18.2
)

require (
    github.com/bytedance/sonic v1.9.1 // indirect
    github.com/json-iterator/go v1.1.12 // indirect
)
```

- `require`：直接依赖
- `// indirect`：间接依赖（被你导入的包所依赖）

## 10. 实战练习

创建如下项目结构：

```
myapp/
├── go.mod
├── main.go
├── utils/
│   └── math.go
└── model/
    └── user.go
```

### utils/math.go

```go
package utils

func Add(a, b int) int {
    return a + b
}

func Multiply(a, b int) int {
    return a * b
}
```

### model/user.go

```go
package model

import "fmt"

type User struct {
    ID   int
    Name string
}

func (u User) String() string {
    return fmt.Sprintf("User{id=%d, name=%s}", u.ID, u.Name)
}
```

### main.go

```go
package main

import (
    "fmt"
    "myapp/model"
    "myapp/utils"
)

func main() {
    sum := utils.Add(2, 3)
    fmt.Println("2+3 =", sum)

    u := model.User{ID: 1, Name: "Tom"}
    fmt.Println(u)
}
```

### 运行

```bash
go mod init myapp
go run main.go
```

## 下一节

[09 - Go 操作 MySQL](09_Go操作MySQL.md)
