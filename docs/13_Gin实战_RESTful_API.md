# 13 - Gin 实战：RESTful API

## 1. 项目结构

一个标准的 Gin + GORM 项目结构：

```
myapp/
├── go.mod
├── go.sum
├── main.go                  # 程序入口
├── config.yaml              # 配置文件
├── config/
│   └── config.go            # 配置加载
├── database/
│   └── db.go                # 数据库初始化
├── model/
│   └── user.go              # 数据模型
├── dto/
│   └── user.go              # 请求/响应结构
├── repository/
│   └── user.go              # 数据访问层
├── service/
│   └── user.go              # 业务逻辑层
├── controller/
│   └── user.go              # 控制器（HTTP 处理）
├── middleware/
│   ├── auth.go              # 鉴权
│   └── logger.go            # 日志
└── router/
    └── router.go            # 路由注册
```

> 类比 Java Spring Boot：
> - controller = @RestController
> - service = @Service
> - repository = @Repository / Mapper
> - model = @Entity
> - dto = DTO/VO
> - middleware = Filter/Interceptor

## 2. 统一响应格式

### 响应结构

```go
// dto/response.go
package dto

type Response struct {
    Code    int         `json:"code"`    // 业务状态码：0 成功，非 0 失败
    Message string      `json:"message"` // 提示信息
    Data    interface{} `json:"data"`    // 数据
}

func Success(data interface{}) Response {
    return Response{Code: 0, Message: "success", Data: data}
}

func Error(code int, msg string) Response {
    return Response{Code: code, Message: msg, Data: nil}
}
```

### 使用

```go
// 成功
c.JSON(200, dto.Success(user))

// 失败
c.JSON(200, dto.Error(1001, "用户不存在"))
```

## 3. 模型与 DTO

### model/user.go

```go
package model

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
    Password  string         `gorm:"size:100;not null" json:"-"`  // json:"-" 不返回
    Email     string         `gorm:"size:100;uniqueIndex" json:"email"`
    Age       int            `json:"age"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### dto/user.go

```go
package dto

type CreateUserReq struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Password string `json:"password" binding:"required,min=6"`
    Email    string `json:"email" binding:"required,email"`
    Age      int    `json:"age" binding:"gte=1,lte=150"`
}

type UpdateUserReq struct {
    Email string `json:"email" binding:"email"`
    Age   int    `json:"age" binding:"gte=1,lte=150"`
}

type LoginReq struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResp struct {
    Token string `json:"token"`
    User  User  `json:"user"`
}

type User struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
}

type PageReq struct {
    Page int `form:"page,default=1" binding:"gte=1"`
    Size int `form:"size,default=10" binding:"gte=1,lte=100"`
}

type PageResp struct {
    Total int64       `json:"total"`
    List  interface{} `json:"list"`
}
```

## 4. Repository 层

```go
// repository/user.go
package repository

import (
    "myapp/model"
    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
    var user model.User
    err := r.db.First(&user, id).Error
    return &user, err
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
    var user model.User
    err := r.db.Where("username = ?", username).First(&user).Error
    return &user, err
}

func (r *UserRepository) FindAll(page, size int) ([]model.User, int64, error) {
    var users []model.User
    var total int64

    r.db.Model(&model.User{}).Count(&total)
    err := r.db.Offset((page - 1) * size).Limit(size).Find(&users).Error
    return users, total, err
}

func (r *UserRepository) Update(user *model.User) error {
    return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
    return r.db.Delete(&model.User{}, id).Error
}
```

## 5. Service 层

```go
// service/user.go
package service

import (
    "errors"
    "myapp/dto"
    "myapp/model"
    "myapp/repository"

    "golang.org/x/crypto/bcrypt"
)

type UserService struct {
    repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) Create(req *dto.CreateUserReq) (*model.User, error) {
    // 检查用户名是否已存在
    existing, err := s.repo.FindByUsername(req.Username)
    if err == nil && existing.ID != 0 {
        return nil, errors.New("用户名已存在")
    }

    // 加密密码
    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &model.User{
        Username: req.Username,
        Password: string(hashed),
        Email:    req.Email,
        Age:      req.Age,
    }

    if err := s.repo.Create(user); err != nil {
        return nil, err
    }
    return user, nil
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
    return s.repo.FindByID(id)
}

func (s *UserService) List(page, size int) ([]model.User, int64, error) {
    return s.repo.FindAll(page, size)
}

func (s *UserService) Update(id uint, req *dto.UpdateUserReq) error {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return err
    }
    if req.Email != "" {
        user.Email = req.Email
    }
    if req.Age > 0 {
        user.Age = req.Age
    }
    return s.repo.Update(user)
}

func (s *UserService) Delete(id uint) error {
    return s.repo.Delete(id)
}

func (s *UserService) Login(req *dto.LoginReq) (*model.User, error) {
    user, err := s.repo.FindByUsername(req.Username)
    if err != nil {
        return nil, errors.New("用户不存在")
    }

    if err := bcrypt.CompareHashAndPassword(
        []byte(user.Password), []byte(req.Password),
    ); err != nil {
        return nil, errors.New("密码错误")
    }
    return user, nil
}
```

## 6. Controller 层

```go
// controller/user.go
package controller

import (
    "net/http"
    "strconv"
    "myapp/dto"
    "myapp/service"
    "github.com/gin-gonic/gin"
)

type UserController struct {
    svc *service.UserService
}

func NewUserController(svc *service.UserService) *UserController {
    return &UserController{svc: svc}
}

func (ctrl *UserController) Create(c *gin.Context) {
    var req dto.CreateUserReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
        return
    }

    user, err := ctrl.svc.Create(&req)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.Error(1001, err.Error()))
        return
    }
    c.JSON(http.StatusCreated, dto.Success(user))
}

func (ctrl *UserController) Get(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    user, err := ctrl.svc.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, dto.Error(1002, "用户不存在"))
        return
    }
    c.JSON(http.StatusOK, dto.Success(user))
}

func (ctrl *UserController) List(c *gin.Context) {
    var req dto.PageReq
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
        return
    }

    users, total, err := ctrl.svc.List(req.Page, req.Size)
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
        return
    }
    c.JSON(http.StatusOK, dto.Success(dto.PageResp{
        Total: total,
        List:  users,
    }))
}

func (ctrl *UserController) Update(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    var req dto.UpdateUserReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
        return
    }
    if err := ctrl.svc.Update(uint(id), &req); err != nil {
        c.JSON(http.StatusBadRequest, dto.Error(1003, err.Error()))
        return
    }
    c.JSON(http.StatusOK, dto.Success(nil))
}

func (ctrl *UserController) Delete(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    if err := ctrl.svc.Delete(uint(id)); err != nil {
        c.JSON(http.StatusBadRequest, dto.Error(1004, err.Error()))
        return
    }
    c.JSON(http.StatusOK, dto.Success(nil))
}
```

## 7. 路由注册

```go
// router/router.go
package router

import (
    "myapp/controller"
    "myapp/middleware"
    "github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, uc *controller.UserController) {
    r.Use(middleware.Logger(), gin.Recovery())

    api := r.Group("/api/v1")
    {
        // 公开接口
        api.POST("/users", uc.Create)
        api.GET("/users", uc.List)
        api.GET("/users/:id", uc.Get)

        // 需要鉴权的接口
        auth := api.Group("/", middleware.Auth())
        {
            auth.PUT("/users/:id", uc.Update)
            auth.DELETE("/users/:id", uc.Delete)
        }
    }
}
```

## 8. 主程序

```go
// main.go
package main

import (
    "log"
    "myapp/config"
    "myapp/controller"
    "myapp/database"
    "myapp/repository"
    "myapp/router"
    "myapp/service"

    "github.com/gin-gonic/gin"
)

func main() {
    // 1. 加载配置
    cfg := config.Load()

    // 2. 初始化数据库
    db := database.Init(cfg)

    // 3. 依赖注入（手动组装）
    userRepo := repository.NewUserRepository(db)
    userSvc := service.NewUserService(userRepo)
    userCtrl := controller.NewUserController(userSvc)

    // 4. 启动 Gin
    r := gin.Default()
    router.Setup(r, userCtrl)

    log.Println("服务启动: http://localhost:" + cfg.Server.Port)
    if err := r.Run(":" + cfg.Server.Port); err != nil {
        log.Fatal(err)
    }
}
```

## 9. 数据库初始化

```go
// database/db.go
package database

import (
    "myapp/config"
    "myapp/model"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func Init(cfg *config.Config) *gorm.DB {
    dsn := cfg.GetDSN()
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    // 自动建表（开发期用）
    db.AutoMigrate(&model.User{})

    return db
}
```

## 10. API 测试

启动服务后，可以用 curl 或 Postman 测试：

```bash
# 创建用户
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"tom","password":"123456","email":"tom@example.com","age":20}'

# 查询列表
curl http://localhost:8080/api/v1/users?page=1&size=10

# 查询单个
curl http://localhost:8080/api/v1/users/1

# 更新（需要 token）
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer mytoken" \
  -d '{"email":"newtom@example.com","age":25}'

# 删除
curl -X DELETE http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer mytoken"
```

## 11. 练习

1. 完成上面的项目骨架，跑通用户 CRUD
2. 加一个文章（Article）模块，实现文章 CRUD
3. 给文章加权限：只有作者本人能修改/删除
4. 实现登录接口，返回 JWT token

## 下一节

[14 - 项目配置与部署](14_项目配置与部署.md)
