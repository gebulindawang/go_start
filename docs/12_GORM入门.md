# 12 - GORM 入门

## 1. GORM 是什么

GORM 是 Go 语言最流行的 ORM 框架，类似 Java 的 MyBatis-Plus 或 Hibernate/JPA。

### 特点

- 全功能 ORM
- 支持关联（一对一、一对多、多对多）
- 支持 Hook（Before/After Create 等）
- 支持事务、迁移、预加载
- 支持多种数据库

### 对比 Java

| Java | GORM |
|------|------|
| `@Entity` | `gorm:"..."` tag |
| `@Table` | `TableName()` 方法 |
| `@Id @GeneratedValue` | `gorm:"primaryKey"` + `gorm:"autoIncrement"` |
| Repository / Mapper | `*gorm.DB` 方法 |
| `@Column` | `gorm:"column:xxx"` |
| JPQL/HQL | 链式 API |

## 2. 安装

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

## 3. 连接数据库

```go
package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
    dsn := "root:password@tcp(127.0.0.1:3306)/go_demo?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }

    DB = db
    return nil
}

func main() {
    if err := InitDB(); err != nil {
        panic(err)
    }
    // 使用 DB 操作数据库...
}
```

## 4. 模型定义

```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"size:50;not null"`
    Age       int
    Email     string         `gorm:"size:100;uniqueIndex"`
    Birthday  *time.Time
    CreatedAt time.Time      // GORM 自动管理
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`  // 软删除
}
```

### 常用 Tag

```go
type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"column:name;size:50;not null;uniqueIndex"`
    Age   int    `gorm:"default:18"`
    Role  string `gorm:"type:enum('admin','user');default:'user'"`
}
```

| Tag | 说明 |
|-----|------|
| `primaryKey` | 主键 |
| `column:xxx` | 列名 |
| `size:50` | 字段长度 |
| `not null` | 非空 |
| `unique` / `uniqueIndex` | 唯一 |
| `index` | 普通索引 |
| `default:18` | 默认值 |
| `type:xxx` | 指定类型 |
| `autoIncrement` | 自增 |

### 自动建表

```go
// 自动迁移（创建或更新表结构）
db.AutoMigrate(&User{})
```

> 类比 Java：相当于 JPA 的 `ddl-auto=update`。
> 生产环境一般不用，应该用手写 SQL 迁移脚本（如 golang-migrate）。

## 5. CRUD 操作

### 增（Create）

```go
// 创建单条
user := User{Name: "Tom", Age: 20, Email: "tom@example.com"}
result := db.Create(&user)
fmt.Println(user.ID)              // 自增 ID 已回填
fmt.Println(result.RowsAffected) // 影响行数
fmt.Println(result.Error)        // 错误

// 批量创建
users := []User{
    {Name: "A", Age: 20},
    {Name: "B", Age: 22},
}
db.Create(&users)

// 选择字段创建
db.Select("Name", "Age").Create(&User{Name: "C", Age: 25})
// INSERT INTO users (name, age) VALUES ('C', 25)
```

### 查（Read）

```go
var user User

// 根据主键查询
db.First(&user, 1)
// SELECT * FROM users WHERE id = 1

// 根据条件查询
db.First(&user, "name = ?", "Tom")
// SELECT * FROM users WHERE name = 'Tom' LIMIT 1

// 查询多条
var users []User
db.Find(&users)                       // 查所有
db.Where("age > ?", 18).Find(&users)  // 条件查询
db.Where("name LIKE ?", "%o%").Find(&users)

// 查询多条并计数
var count int64
db.Model(&User{}).Where("age > ?", 18).Count(&count)

// 排序、分页、限制
db.Order("age desc").Limit(10).Offset(0).Find(&users)

// 选择字段
db.Select("name, age").Find(&users)

// Distinct
db.Distinct("name").Find(&users)

// 取第一条（无数据返回 ErrRecordNotFound）
var first User
result := db.First(&first)
if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    fmt.Println("没有数据")
}
```

### 链式条件

```go
var users []User
db.Where("age >= ?", 18).
    Where("age <= ?", 30).
    Where("name LIKE ?", "T%").
    Order("age asc").
    Limit(10).
    Find(&users)
```

### 改（Update）

```go
// 保存整个对象
user.Age = 25
db.Save(&user)  // UPDATE users SET ... WHERE id = ?

// 更新单个字段
db.Model(&User{}).Where("id = ?", 1).Update("age", 30)

// 更新多个字段
db.Model(&user).Updates(map[string]interface{}{
    "name": "Tom2",
    "age":  30,
})

// 用 struct 更新（零值字段不会更新）
db.Model(&user).Updates(User{Name: "Tom2", Age: 30})

// 更新选定字段（包括零值）
db.Model(&user).Select("age").Updates(User{Age: 0})  // age 会被设为 0
```

### 删（Delete）

```go
// 删除（软删除，因为有 DeletedAt 字段）
db.Delete(&user)
db.Delete(&User{}, 1)            // 按主键删
db.Where("age < ?", 18).Delete(&User{})

// 永久删除（绕过软删除）
db.Unscoped().Delete(&user)

// 查询包含软删除的记录
db.Unscoped().Find(&users)
```

> 软删除：模型有 `DeletedAt` 字段时，删除只是设置 `deleted_at` 时间，不真正删除。

## 6. 事务

### 自动事务

```go
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&User{Name: "A"}).Error; err != nil {
        return err  // 返回错误会自动回滚
    }
    if err := tx.Create(&User{Name: "B"}).Error; err != nil {
        return err
    }
    return nil  // 返回 nil 自动提交
})
```

### 手动事务

```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&User{Name: "A"}).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&User{Name: "B"}).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

## 7. 关联关系

### 一对多

```go
type User struct {
    ID      uint
    Name    string
    // 一对多：一个用户有多篇文章
    Articles []Article
}

type Article struct {
    ID     uint
    Title  string
    UserID uint  // 外键
}

// 创建时自动关联
db.Create(&User{
    Name: "Tom",
    Articles: []Article{
        {Title: "Go 入门"},
        {Title: "Gin 实战"},
    },
})

// 预加载（带出关联数据）
var user User
db.Preload("Articles").First(&user, 1)
fmt.Println(user.Articles)
```

### 多对多

```go
type Student struct {
    ID    uint
    Name  string
    // 多对多：学生 ↔ 课程
    Courses []Course `gorm:"many2many:student_courses;"`
}

type Course struct {
    ID    uint
    Title string
}
```

## 8. Hook（钩子）

类似 Java JPA 的 `@PrePersist`。

```go
type User struct {
    ID   uint
    Name string
    UUID string
}

// 创建前自动生成 UUID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.UUID = uuid.New().String()
    return
}

// 创建后的钩子
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
    fmt.Println("用户已创建:", u.Name)
    return
}
```

可用钩子：
- `BeforeSave` / `AfterSave`
- `BeforeCreate` / `AfterCreate`
- `BeforeUpdate` / `AfterUpdate`
- `BeforeDelete` / `AfterDelete`
- `AfterFind`

## 9. 完整示例

```go
package main

import (
    "fmt"
    "log"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Name      string         `gorm:"size:50;not null" json:"name"`
    Age       int            `json:"age"`
    Email     string         `gorm:"size:100;uniqueIndex" json:"email"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

var DB *gorm.DB

func initDB() {
    dsn := "root:password@tcp(127.0.0.1:3306)/go_demo?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    DB = db
    DB.AutoMigrate(&User{})
}

func main() {
    initDB()

    // 增
    user := User{Name: "Tom", Age: 20, Email: "tom@example.com"}
    DB.Create(&user)
    fmt.Println("创建:", user)

    // 查
    var u User
    DB.First(&u, user.ID)
    fmt.Println("查询:", u)

    // 改
    DB.Model(&u).Update("age", 21)
    fmt.Println("更新后:", u)

    // 删
    DB.Delete(&u)
    fmt.Println("已删除")
}
```

## 10. 练习

1. 定义 `Article` 模型，包含 `Title`、`Content`、`AuthorID`，实现 CRUD
2. 实现用户与文章的一对多关系，用 `Preload` 查询用户及其所有文章
3. 用事务实现一个简单的转账逻辑
4. 实现分页查询：`GET /users?page=1&size=10`

## 下一节

[13 - Gin 实战：RESTful API](13_Gin实战_RESTful_API.md)
