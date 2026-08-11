# 10 - Gin 框架入门

## 1. Gin 是什么

Gin 是 Go 语言最流行的 Web 框架，类似 Java 的 Spring Boot。

### 特点

- **高性能**：基于 httprouter，速度极快
- **简洁**：API 设计优雅
- **中间件支持**：类似过滤器/拦截器
- **JSON 处理方便**
- **错误管理**：内置机制

### 对比 Java

| Java | Go |
|------|-----|
| Spring Boot | Gin |
| @RestController | gin.Engine + handler 函数 |
| @RequestMapping | r.GET/POST |
| @RequestParam | c.Query / c.Param |
| @RequestBody | c.ShouldBindJSON |
| Filter / Interceptor | Middleware |

## 2. 安装

```bash
# 初始化模块
go mod init myginapp

# 安装 Gin
go get -u github.com/gin-gonic/gin
```

## 3. 第一个 Gin 程序

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    // 创建路由引擎
    r := gin.Default()

    // 定义路由
    r.GET("/hello", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello, Gin!",
        })
    })

    // 启动服务，默认端口 8080
    r.Run(":8080")
}
```

### 运行

```bash
go run main.go
```

### 测试

浏览器访问 `http://localhost:8080/hello`，或用 curl：

```bash
curl http://localhost:8080/hello
# {"message":"Hello, Gin!"}
```

## 4. 路由（Routes）

### 基本路由

```go
r.GET("/users", getUsers)        // GET 请求
r.POST("/users", createUser)     // POST 请求
r.PUT("/users/:id", updateUser)  // PUT 请求
r.DELETE("/users/:id", deleteUser)  // DELETE 请求
r.PATCH("/users/:id", patchUser)
r.HEAD("/users", headUsers)
r.OPTIONS("/users", optionsUsers)

// 匹配所有方法
r.Any("/test", anyHandler)

// 无匹配路由（404）
r.NoRoute(func(c *gin.Context) {
    c.JSON(404, gin.H{"message": "页面不存在"})
})
```

### 路由参数

```go
// 路径参数
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"id": id})
})
// 访问 /users/123  ->  {"id":"123"}

// 通配符
r.GET("/files/*filepath", func(c *gin.Context) {
    path := c.Param("filepath")
    c.JSON(200, gin.H{"path": path})
})
// 访问 /files/a/b.txt  ->  {"path":"/a/b.txt"}
```

### 查询参数

```go
r.GET("/search", func(c *gin.Context) {
    keyword := c.Query("keyword")           // 必须的查询参数
    page := c.DefaultQuery("page", "1")     // 带默认值
    size := c.Query("size")

    c.JSON(200, gin.H{
        "keyword": keyword,
        "page":    page,
        "size":    size,
    })
})
// 访问 /search?keyword=go&page=2
```

### 表单参数

```go
r.POST("/login", func(c *gin.Context) {
    username := c.PostForm("username")
    password := c.PostForm("password")

    c.JSON(200, gin.H{
        "username": username,
        "password": password,
    })
})
```

## 5. 返回数据

### 返回 JSON

```go
// 方式 1：gin.H（本质是 map[string]interface{}）
r.GET("/json1", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "ok", "code": 0})
})

// 方式 2：结构体
r.GET("/json2", func(c *gin.Context) {
    user := User{ID: 1, Name: "Tom"}
    c.JSON(200, user)
})

// 方式 3：更快捷的 JSON
r.GET("/json3", func(c *gin.Context) {
    c.IndentedJSON(200, gin.H{"msg": "美化输出"})
})
```

### 返回其他格式

```go
// 返回字符串
r.GET("/text", func(c *gin.Context) {
    c.String(200, "Hello %s", "World")
})

// 返回 XML
r.GET("/xml", func(c *gin.Context) {
    c.XML(200, gin.H{"message": "ok"})
})

// 返回文件
r.GET("/file", func(c *gin.Context) {
    c.File("./main.go")
})

// 重定向
r.GET("/old", func(c *gin.Context) {
    c.Redirect(302, "/new")
})
```

## 6. 请求绑定（参数解析）⭐

### JSON 绑定

```go
type LoginReq struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

r.POST("/login", func(c *gin.Context) {
    var req LoginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"msg": "登录成功", "user": req.Username})
})

// 请求：curl -X POST http://localhost:8080/login \
//   -H "Content-Type: application/json" \
//   -d '{"username":"admin","password":"123456"}'
```

### 表单绑定

```go
type RegisterReq struct {
    Username string `form:"username" binding:"required"`
    Email    string `form:"email" binding:"required,email"`
}

r.POST("/register", func(c *gin.Context) {
    var req RegisterReq
    if err := c.ShouldBind(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"msg": "注册成功"})
})
```

### URI 绑定

```go
type UserReq struct {
    ID int `uri:"id" binding:"required"`
}

r.GET("/users/:id", func(c *gin.Context) {
    var req UserReq
    if err := c.ShouldBindUri(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"id": req.ID})
})
```

### Binding 标签

| 标签 | 说明 |
|------|------|
| `json:"name"` | JSON 字段名 |
| `form:"name"` | 表单字段名 |
| `uri:"name"` | URI 参数名 |
| `binding:"required"` | 必填 |
| `binding:"email"` | 邮箱格式 |
| `binding:"min=6,max=20"` | 长度限制 |
| `binding:"one=a b c"` | 枚举值 |

## 7. 分组路由

```go
func main() {
    r := gin.Default()

    // API v1 分组
    v1 := r.Group("/api/v1")
    {
        v1.GET("/users", listUsers)
        v1.GET("/users/:id", getUser)
        v1.POST("/users", createUser)
        v1.PUT("/users/:id", updateUser)
        v1.DELETE("/users/:id", deleteUser)
    }

    // API v2 分组
    v2 := r.Group("/api/v2")
    {
        v2.GET("/users", listUsersV2)
    }

    r.Run(":8080")
}
```

> 类比 Java：相当于 `@RequestMapping("/api/v1")` 类级别的路径前缀。

## 8. 完整示例

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

var users = []User{
    {ID: 1, Name: "Tom", Age: 20},
    {ID: 2, Name: "Jerry", Age: 22},
}

func main() {
    r := gin.Default()

    // 分组路由
    api := r.Group("/api")
    {
        api.GET("/users", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "code": 0,
                "data": users,
            })
        })

        api.GET("/users/:id", func(c *gin.Context) {
            id := c.Param("id")
            for _, u := range users {
                if fmt.Sprintf("%d", u.ID) == id {
                    c.JSON(http.StatusOK, gin.H{"data": u})
                    return
                }
            }
            c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
        })

        api.POST("/users", func(c *gin.Context) {
            var u User
            if err := c.ShouldBindJSON(&u); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
            }
            u.ID = len(users) + 1
            users = append(users, u)
            c.JSON(http.StatusCreated, gin.H{"data": u})
        })
    }

    r.Run(":8080")
}
```

## 9. 热重载（可选）

开发时每次改代码都要重启，可以安装 `air`：

```bash
go install github.com/cosmtrek/air@latest

# 在项目目录运行
air
```

修改代码后会自动重新编译运行。

## 10. 练习

1. 实现一个完整的用户 CRUD API（GET/POST/PUT/DELETE）
2. 用分组路由组织 `/api/v1` 和 `/api/v2`
3. 给 `POST /users` 添加参数校验（name 必填，age >= 0）
4. 实现一个简单的查询接口，支持 `?keyword=xxx&page=1&size=10`

## 下一节

[11 - Gin 路由、参数与中间件](11_Gin路由_参数_中间件.md)
