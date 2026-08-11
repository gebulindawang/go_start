# 11 - Gin 路由、参数与中间件

## 1. 路由进阶

### 路由参数

```go
r := gin.Default()

// 路径参数 :id
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"id": id})
})

// 通配符 *filepath（匹配任意路径）
r.GET("/static/*filepath", func(c *gin.Context) {
    path := c.Param("filepath")
    c.File("./static" + path)
})
```

### 查询参数

```go
r.GET("/search", func(c *gin.Context) {
    q := c.Query("q")              // 获取查询参数
    page := c.DefaultQuery("page", "1")  // 带默认值
    c.JSON(200, gin.H{"q": q, "page": page})
})
```

### Header 参数

```go
r.GET("/header", func(c *gin.Context) {
    ua := c.GetHeader("User-Agent")
    token := c.GetHeader("Authorization")
    c.JSON(200, gin.H{"ua": ua, "token": token})
})
```

### Cookie

```go
r.GET("/set_cookie", func(c *gin.Context) {
    c.SetCookie("token", "abc123", 3600, "/", "localhost", false, true)
    c.JSON(200, gin.H{"msg": "cookie 已设置"})
})

r.GET("/get_cookie", func(c *gin.Context) {
    token, err := c.Cookie("token")
    if err != nil {
        c.JSON(400, gin.H{"error": "no cookie"})
        return
    }
    c.JSON(200, gin.H{"token": token})
})
```

## 2. 参数绑定（Binding）

### ShouldBind vs MustBind

| 方法 | 行为 |
|------|------|
| `ShouldBindJSON` | 绑定失败返回 error，**自己处理** |
| `MustBind` / `Bind` | 绑定失败自动返回 400（不灵活，**不推荐**） |

**推荐统一用 `ShouldBind*` 系列**。

### 不同来源的绑定

```go
type User struct {
    Name  string `json:"name" form:"name" uri:"name" binding:"required"`
    Age   int    `json:"age" form:"age" uri:"age" binding:"gte=0,lte=150"`
    Email string `json:"email" form:"email" binding:"email"`
}

// JSON body
r.POST("/users", func(c *gin.Context) {
    var u User
    if err := c.ShouldBindJSON(&u); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, u)
})

// Form
r.POST("/form", func(c *gin.Context) {
    var u User
    if err := c.ShouldBind(&u); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, u)
})

// URI
r.GET("/users/:name/:age", func(c *gin.Context) {
    var u User
    if err := c.ShouldBindUri(&u); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, u)
})

// Query
r.GET("/query", func(c *gin.Context) {
    var u User
    if err := c.ShouldBindQuery(&u); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, u)
})
```

### 常用校验规则

```go
type RegisterReq struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Password string `json:"password" binding:"required,min=6"`
    Email    string `json:"email" binding:"required,email"`
    Age      int    `json:"age" binding:"required,gte=1,lte=150"`
    Phone    string `json:"phone" binding:"required,e164"`
    Gender   string `json:"gender" binding:"required,oneof=male female"`
    URL      string `json:"url" binding:"url"`
}
```

| 规则 | 说明 |
|------|------|
| `required` | 必填 |
| `min=N,max=M` | 字符串长度/数字范围 |
| `gte=N,lte=M` | 大于等于/小于等于 |
| `oneof=a b c` | 枚举值 |
| `email` | 邮箱格式 |
| `url` | URL 格式 |
| `len=N` | 长度等于 N |
| `dive` | 切片/Map 深入校验 |

## 3. 中间件（Middleware）⭐

中间件类似 Java 的 Filter/Interceptor，在请求前后执行。

### 内置中间件

```go
r := gin.Default()
// Default() 等于:
// r := gin.New()
// r.Use(gin.Logger(), gin.Recovery())
```

### 自定义中间件

```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        // 请求前
        path := c.Request.URL.Path
        method := c.Request.Method

        // 让请求继续
        c.Next()

        // 请求后
        duration := time.Since(start)
        status := c.Writer.Status()
        log.Printf("[%s] %s %d %v", method, path, status, duration)
    }
}

func main() {
    r := gin.New()
    r.Use(Logger())
    r.Use(gin.Recovery())

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"msg": "pong"})
    })

    r.Run(":8080")
}
```

### 中间件的执行顺序

```
请求 -> 中间件1 前置 -> 中间件2 前置 -> Handler -> 中间件2 后置 -> 中间件1 后置 -> 响应
```

`c.Next()` 之前的代码是**前置**，之后的代码是**后置**。

### 中止请求

```go
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "未授权"})
            c.Abort()  // 中止后续处理
            return
        }
        // 验证 token...
        c.Set("userID", 1001)  // 在上下文中存数据
        c.Next()
    }
}

func main() {
    r := gin.Default()

    // 只对 /admin 路径生效
    admin := r.Group("/admin", Auth())
    admin.GET("/dashboard", func(c *gin.Context) {
        // 取出中间件设置的值
        uid := c.MustGet("userID").(int)
        c.JSON(200, gin.H{"uid": uid, "msg": "欢迎进入后台"})
    })

    r.Run(":8080")
}
```

### 上下文传值（c.Set / c.Get）

```go
func UserMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("user", &User{ID: 1, Name: "Tom"})
        c.Next()
    }
}

r.GET("/profile", UserMiddleware(), func(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(500, gin.H{"error": "用户不存在"})
        return
    }
    c.JSON(200, user)
})
```

### 常用中间件

#### CORS 跨域

```bash
go get github.com/gin-contrib/cors
```

```go
import "github.com/gin-contrib/cors"

r.Use(cors.Default())
// 或自定义配置
r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))
```

#### 日志中间件

```go
func LoggerToFile() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        log.Printf("%s %s %d %v\n",
            c.Request.Method,
            c.Request.URL.Path,
            c.Writer.Status(),
            time.Since(start),
        )
    }
}
```

## 4. 文件上传

### 单文件上传

```go
r.POST("/upload", func(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 保存文件
    dst := "./uploads/" + file.Filename
    if err := c.SaveUploadedFile(file, dst); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "filename": file.Filename,
        "size":     file.Size,
        "msg":      "上传成功",
    })
})
```

### 多文件上传

```go
r.POST("/uploads", func(c *gin.Context) {
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    files := form.File["files"]
    for _, file := range files {
        dst := "./uploads/" + file.Filename
        c.SaveUploadedFile(file, dst)
    }

    c.JSON(200, gin.H{"msg": fmt.Sprintf("上传 %d 个文件成功", len(files))})
})
```

## 5. 静态文件服务

```go
// 提供 /static 路径下的静态文件
r.Static("/static", "./static")

// 直接作为根路径资源
r.StaticFile("/favicon.ico", "./favicon.ico")

// 也可以这样
r.StaticFS("/files", gin.Dir("./uploads", false))
```

## 6. 模板渲染（可选）

```go
r.LoadHTMLGlob("templates/*")

r.GET("/index", func(c *gin.Context) {
    c.HTML(200, "index.html", gin.H{
        "title": "Go + Gin",
        "users": []string{"Tom", "Jerry"},
    })
})
```

templates/index.html:

```html
<!DOCTYPE html>
<html>
<head><title>{{.title}}</title></head>
<body>
    <h1>{{.title}}</h1>
    <ul>
    {{range .users}}
        <li>{{.}}</li>
    {{end}}
    </ul>
</body>
</html>
```

## 7. 练习

1. 写一个鉴权中间件，检查请求头中的 `Authorization` 字段
2. 实现一个简单的日志中间件，记录所有请求的方法、路径、状态码、耗时
3. 实现一个文件上传接口，限制文件大小为 10MB
4. 用分组路由 + 中间件实现 `/api/public` 和 `/api/private` 两组接口

## 下一节

[12 - GORM 入门](12_GORM入门.md)
