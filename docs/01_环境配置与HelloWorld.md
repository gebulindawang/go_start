# 01 - 环境配置与 Hello World

## 1. 安装 Go

### Windows
1. 访问 https://go.dev/dl/
2. 下载 Windows 安装包（`.msi`），例如 `go1.22.0.windows-amd64.msi`
3. 双击安装，默认路径 `C:\Program Files\Go`
4. 安装程序会自动配置环境变量

### 验证安装
打开终端（cmd 或 PowerShell），输入：

```bash
go version
# 输出: go version go1.22.0 windows/amd64
```

### 配置国内代理（推荐）

```bash
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

> 类比 Java：Go 的 module 系统相当于 Maven/Gradle，GOPROXY 相当于配置国内镜像源。

## 2. 关键环境变量

| 变量 | 说明 | Java 类比 |
|------|------|-----------|
| `GOROOT` | Go 安装目录 | `JAVA_HOME` |
| `GOPATH` | 工作区（旧模式用） | 类似 `~/.m2` |
| `GO111MODULE` | 模块模式开关 | 默认 on，无需关心 |
| `GOPROXY` | 模块下载代理 | Maven mirror |

> **新版本（1.16+）默认使用 Module 模式**，无需关心 GOPATH。

## 3. 开发工具推荐

- **VS Code** + Go 扩展（推荐，轻量免费）
- **GoLand**（JetBrains 出品，对学生免费，体验最好）

### VS Code 配置
1. 安装 Go 扩展
2. `Ctrl+Shift+P` 输入 `Go: Install/Update Tools`，全选安装

## 4. 第一个 Go 程序

### 创建项目目录

```bash
mkdir D:\code\go\go_start
cd D:\code\go\go_start
```

### 初始化模块

```bash
go mod init go_start
```

执行后会生成 `go.mod` 文件，内容类似：

```
module go_start

go 1.22
```

> 类比 Java：相当于 `pom.xml` 或 `build.gradle`，`module go_start` 是模块名。

### 创建 main.go

在项目根目录新建 `main.go`：

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

### 运行

```bash
go run main.go
# 输出: Hello, Go!
```

### 编译

```bash
go build main.go
# 生成 main.exe（Windows）或 main（Linux/Mac）
./main.exe
```

> 类比 Java：`go run` 相当于直接运行源码，`go build` 相当于 `javac` + 打包成可执行文件（不需要 JVM）。

## 5. Go 程序结构

```go
package main          // 包声明，main 表示可执行程序（类似 Java 的 main 类）

import "fmt"          // 导入标准库

func main() {         // 主函数，程序入口
    fmt.Println("Hi")
}
```

### 注意事项

- `package main` + `func main()` 是可执行程序的入口
- 左花括号 `{` 必须跟函数签名在同一行（强制风格）
- 语句末尾不需要分号
- 未使用的 import 和变量会**编译报错**（不是警告，是错误！）

## 6. 练习

1. 修改 `main.go`，输出你的名字
2. 编译并运行生成的可执行文件
3. 尝试故意留一个未使用的 import，观察报错信息

## 下一节

[02 - Go 基础语法（Java 视角）](02_Go基础语法_Java视角.md)
