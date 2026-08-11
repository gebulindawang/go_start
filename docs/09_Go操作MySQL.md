# 09 - Go 操作 MySQL

## 1. 准备工作

### 安装 MySQL 驱动

```bash
go get github.com/go-sql-driver/mysql
```

### 准备数据库

```sql
CREATE DATABASE IF NOT EXISTS go_demo DEFAULT CHARSET utf8mb4;

USE go_demo;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    age INT,
    email VARCHAR(100)
);

INSERT INTO users (name, age, email) VALUES
('Tom', 20, 'tom@example.com'),
('Jerry', 22, 'jerry@example.com');
```

## 2. 连接数据库

### Java 写法
```java
Connection conn = DriverManager.getConnection(
    "jdbc:mysql://localhost:3306/go_demo?useSSL=false",
    "root", "password"
);
```

### Go 写法

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"  // 匿名导入，执行 init 注册驱动
)

func main() {
    // DSN: 用户名:密码@tcp(主机:端口)/数据库名?参数
    dsn := "root:password@tcp(127.0.0.1:3306)/go_demo?charset=utf8mb4&parseTime=True&loc=Local"

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 测试连接
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("数据库连接成功！")
}
```

### DSN 格式说明

```
root:password@tcp(127.0.0.1:3306)/go_demo?charset=utf8mb4&parseTime=True&loc=Local
```

- `root:password`：用户名:密码
- `tcp(127.0.0.1:3306)`：协议://地址:端口
- `/go_demo`：数据库名
- 参数：
  - `charset=utf8mb4`：字符集
  - `parseTime=True`：自动解析时间类型
  - `loc=Local`：使用本地时区

## 3. 连接池配置

```go
db, _ := sql.Open("mysql", dsn)

// 配置连接池
db.SetMaxOpenConns(20)   // 最大连接数
db.SetMaxIdleConns(10)   // 最大空闲连接数
db.SetConnMaxLifetime(time.Hour)  // 连接最大存活时间
```

## 4. 增删改查（CRUD）

### 增 - Insert

```go
func createUser(db *sql.DB, name string, age int, email string) (int64, error) {
    result, err := db.Exec(
        "INSERT INTO users (name, age, email) VALUES (?, ?, ?)",
        name, age, email,
    )
    if err != nil {
        return 0, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }
    return id, nil
}

// 调用
id, err := createUser(db, "Alice", 25, "alice@example.com")
fmt.Println("新用户 ID:", id)
```

### 查 - Select 单行

```go
type User struct {
    ID    int
    Name  string
    Age   int
    Email string
}

func getUserByID(db *sql.DB, id int) (*User, error) {
    var u User
    err := db.QueryRow(
        "SELECT id, name, age, email FROM users WHERE id = ?",
        id,
    ).Scan(&u.ID, &u.Name, &u.Age, &u.Email)

    if err == sql.ErrNoRows {
        return nil, nil  // 没找到
    }
    if err != nil {
        return nil, err
    }
    return &u, nil
}
```

### 查 - Select 多行

```go
func getAllUsers(db *sql.DB) ([]User, error) {
    rows, err := db.Query("SELECT id, name, age, email FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()  // 重要！释放连接

    var users []User
    for rows.Next() {
        var u User
        err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.Email)
        if err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    // 检查遍历过程中是否有错误
    if err = rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}
```

### 改 - Update

```go
func updateUserAge(db *sql.DB, id int, age int) error {
    result, err := db.Exec(
        "UPDATE users SET age = ? WHERE id = ?",
        age, id,
    )
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("没有用户被更新")
    }
    return nil
}
```

### 删 - Delete

```go
func deleteUser(db *sql.DB, id int) error {
    result, err := db.Exec("DELETE FROM users WHERE id = ?", id)
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("用户不存在")
    }
    return nil
}
```

## 5. 预处理（Prepared Statement）

```go
func batchInsert(db *sql.DB, users []User) error {
    stmt, err := db.Prepare("INSERT INTO users (name, age, email) VALUES (?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, u := range users {
        _, err := stmt.Exec(u.Name, u.Age, u.Email)
        if err != nil {
            return err
        }
    }
    return nil
}
```

## 6. 事务（Transaction）

```go
func transferMoney(db *sql.DB, fromID, toID int, amount float64) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    // 注意：defer 中处理回滚
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    _, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
    if err != nil {
        return err
    }

    _, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)
    if err != nil {
        return err
    }

    err = tx.Commit()
    return err
}
```

## 7. 完整示例

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    ID    int
    Name  string
    Age   int
    Email string
}

func main() {
    dsn := "root:password@tcp(127.0.0.1:3306)/go_demo?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    // 增
    result, _ := db.Exec(
        "INSERT INTO users (name, age, email) VALUES (?, ?, ?)",
        "Bob", 30, "bob@example.com",
    )
    id, _ := result.LastInsertId()
    fmt.Println("新增用户 ID:", id)

    // 查
    rows, _ := db.Query("SELECT id, name, age, email FROM users")
    defer rows.Close()

    fmt.Println("用户列表:")
    for rows.Next() {
        var u User
        rows.Scan(&u.ID, &u.Name, &u.Age, &u.Email)
        fmt.Printf("  %d: %s, %d, %s\n", u.ID, u.Name, u.Age, u.Email)
    }

    // 改
    db.Exec("UPDATE users SET age = ? WHERE name = ?", 31, "Bob")
    fmt.Println("已更新 Bob 的年龄")

    // 删
    db.Exec("DELETE FROM users WHERE name = ?", "Bob")
    fmt.Println("已删除 Bob")
}
```

## 8. 注意事项

### 重要点

1. **永远使用 `?` 占位符**，不要拼接 SQL 字符串（防 SQL 注入）
2. **`rows.Close()` 必须调用**，否则会泄漏连接
3. **`sql.ErrNoRows` 单独处理**，不是真正的错误
4. **`defer db.Close()`** 在程序结束时关闭

### 错误处理示例

```go
err := db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
switch {
case err == sql.ErrNoRows:
    fmt.Println("用户不存在")
case err != nil:
    log.Fatal(err)
default:
    fmt.Println("用户名:", name)
}
```

## 9. 练习

1. 创建一个 `users` 表，实现完整的 CRUD
2. 写一个函数根据 name 模糊查询用户
3. 实现一个简单的转账事务（两个账户扣款加款）
4. 用预处理批量插入 10 个用户

> **提示**：实际开发中很少直接用 `database/sql`，会用 ORM 框架如 **GORM**，后面会学到。

## 下一节

[10 - Gin 框架入门](10_Gin框架入门.md)
