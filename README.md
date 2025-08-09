# 🚀 Go-Util - 现代化 Go 工具库

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Release](https://img.shields.io/github/v/release/zhoudm1743/go-util?style=for-the-badge)
![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen?style=for-the-badge)

**让 Go 代码更简洁，让开发更高效！**

一个功能强大、类型安全、链式调用的现代化 Go 工具库

[快速开始](#-快速开始) •
[核心功能](#-核心功能) •
[文档](#-文档) •
[示例](#-示例)

</div>

---

## ✨ 为什么选择 Go-Util？

### 🎯 **链式调用，优雅编程**
```go
// 传统写法 - 繁琐且容易出错
emails := []string{}
for _, user := range users {
    if user.IsActive && isValidEmail(user.Email) {
        emails = append(emails, user.Email)
    }
}
sort.Strings(emails)
uniqueEmails := removeDuplicates(emails)

// Go-Util - 链式优雅
validEmails := util.Array(users).
    Filter(func(u User) bool { return u.IsActive }).
    Map(func(u User) string { return u.Email }).
    Filter(func(email string) bool { return util.Str(email).IsEmail() }).
    Sort().
    Distinct().
    ToSlice()
```

### 🛡️ **类型安全，泛型支持**
```go
// 编译时类型检查，无运行时类型错误
numbers := util.Arrays(1, 2, 3, 4, 5).     // XArray[int]
    Filter(func(n int) bool { return n > 2 }). // 类型安全的过滤
    Map(func(n int) int { return n * 2 }).     // 类型安全的映射
    ToSlice()                                   // []int
```

### ⚡ **高性能，零依赖核心**
- **极速 JSON 操作**：200层嵌套，5000个对象，0.01秒处理 🚀
- **内存优化**：零拷贝字符串转换，智能内存管理
- **并发安全**：线程安全的枚举系统，支持高并发场景

---

## 🚀 快速开始

### 安装
```bash
go get github.com/zhoudm1743/go-util
```

### 5分钟上手
```go
package main

import (
    "fmt"
    util "github.com/zhoudm1743/go-util"
)

func main() {
    // 🔤 字符串操作 - 链式调用
    result := util.Str("hello_world_example").
        Snake2BigCamel().        // 转大驼峰: "HelloWorldExample"
        ReplaceAll("World", "Go").// 替换: "HelloGoExample"
        Lower().                 // 转小写: "hellogoexample"
        FirstUpper().            // 首字母大写: "Hellogoexample"
        String()
    fmt.Println(result) // "Hellogoexample"
    
    // 📊 数组操作 - 函数式编程
    evenDoubled := util.Arrays(1, 2, 3, 4, 5, 6).
        Filter(func(n int) bool { return n%2 == 0 }). // [2, 4, 6]
        Map(func(n int) int { return n * 2 }).        // [4, 8, 12]
        ToSlice()
    fmt.Println(evenDoubled) // [4, 8, 12]
    
    // ⏰ 时间操作 - 直观易用
    nextWeek := util.Now().
        AddDays(7).
        StartOfDay().
        FormatChinese()
    fmt.Println(nextWeek) // "2024年1月22日 00时00分00秒"
    
    // 🌐 HTTP 请求 - 现代化
    resp, _ := util.Http().
        BaseURL("https://jsonplaceholder.typicode.com").
        Header("User-Agent", "Go-Util").
        Timeout(10).
        Get("/posts/1")
    
    if resp.IsOK() {
        fmt.Println("请求成功:", resp.String()[:50] + "...")
    }
}
```

---

## 🛠️ 核心功能

<table>
<tr>
<td>

### 🔤 字符串处理
```go
str := util.Str("hello@example.com")

// 验证和判断
str.IsEmail()        // true
str.IsURL()          // false
str.Contains("@")    // true
str.HasPrefix("hello") // true

// 格式转换
str.Upper()          // "HELLO@EXAMPLE.COM"
str.Camel2Snake()    // "hello_example_com"
str.Slugify()        // "hello-example-com"

// 安全操作
str.MD5()            // "5d41402abc4b2a76b9719d911017c592"
str.Base64Encode()   // "aGVsbG9AZXhhbXBsZS5jb20="

// 智能处理
str.WordCount()      // 1
str.Similarity(other) // 0.85
```

</td>
<td>

### 📊 数组处理
```go
arr := util.Arrays(5, 2, 8, 1, 9)

// 排序和筛选
arr.Sort()           // [1, 2, 5, 8, 9]
arr.Filter(func(n int) bool { 
    return n > 3 
})                   // [5, 8, 9]

// 函数式操作
arr.Map(func(n int) string {
    return fmt.Sprintf("num_%d", n)
})                   // ["num_5", "num_8", "num_9"]

// 聚合计算
arr.Sum()            // 22
arr.Average()        // 7.33
arr.Max()            // 9

// 集合操作
arr.Distinct()       // 去重
arr.Chunk(2)         // 分块: [[5,8], [9]]
```

</td>
</tr>
<tr>
<td>

### ⏰ 时间处理
```go
now := util.Now()
birthday := util.Date(1990, 5, 15, 14, 30, 0)

// 智能格式化
now.FormatRelative()    // "2小时前"
now.FormatChinese()     // "2024年1月15日"
now.FormatRFC3339()     // "2024-01-15T14:30:45Z"

// 时间计算
birthday.Age()          // 33
now.DaysTo(future)      // 15
now.Between(start, end) // true

// 时间范围
now.StartOfWeek()       // 本周开始
now.EndOfMonth()        // 本月结束
now.Quarter()           // 1 (第一季度)
```

</td>
<td>

### 🗺️ 映射操作
```go
m := util.NewMap[string, int]()

// 安全操作
m.Set("apple", 10)
m.GetOrDefault("grape", 0) // 0
m.Has("apple")             // true

// 批量处理
m.Keys()                   // ["apple"]
m.Values()                 // [10]
m.Filter(func(k string, v int) bool {
    return v > 5
})                         // 只保留值>5的项

// 便捷转换
m.ToJSON()                 // {"apple":10}
m.Equal(other)             // 深度比较
```

</td>
</tr>
</table>

---

## 🎯 高级功能

### 🚀 JSONx - 极致 JSON 操作
```go
import "github.com/zhoudm1743/go-util/jsonx"

// 🔥 深度路径操作 - 200层嵌套无压力
user := jsonx.Object().
    Set("profile.personal.name", "张三").
    Set("profile.contact.emails.0", "zhang@example.com").
    Set("settings.theme.color", "dark").
    Set("permissions.admin.access", true)

name := user.Get("profile.personal.name").String() // "张三"
email := user.Get("profile.contact.emails.0").String() // "zhang@example.com"

// 🔥 链式构建复杂结构
api := jsonx.NewBuilder().
    AddString("version", "1.0").
    AddObject("user", jsonx.NewBuilder().
        AddString("name", "李四").
        AddArray("roles", jsonx.NewArrayBuilder().
            AppendString("admin").
            AppendString("user").
            Build()).
        Build()).
    Build()

// 🔥 函数式数组处理
numbers := jsonx.QuickArray(1, 2, 3, 4, 5).
    Filter(func(key string, value *jsonx.JSON) bool {
        return value.Int()%2 == 0  // 过滤偶数
    }).
    Map(func(key string, value *jsonx.JSON) interface{} {
        return value.Int() * 10    // 乘以10
    })
// 结果: [20, 40]

// 🔥 高级工具
flattened := jsonx.Flatten(nested)  // 扁平化嵌套结构
merged := jsonx.DeepMerge(obj1, obj2) // 深度合并
picked := jsonx.Pick(user, "name", "email") // 选择字段
```

**极限性能测试通过**：
- ✅ **200层嵌套深度**
- ✅ **5000个对象处理** (0.01秒)
- ✅ **1000个数组元素**
- ✅ **复杂对象数组混合嵌套**

### 🔐 JWT - 企业级认证
```go
import "github.com/zhoudm1743/go-util/jwt"

// 🔥 通用 API - 算法参数化
secret := []byte("your-secret-key")
claims := jwt.MapClaims{
    "sub": "user123",
    "name": "张三",
    "role": "admin",
    "exp": time.Now().Add(24 * time.Hour).Unix(),
}

// 支持所有算法的统一接口
algorithms := []jwt.SigningMethod{
    jwt.SigningMethodHS256,  // HMAC SHA-256
    jwt.SigningMethodHS384,  // HMAC SHA-384  
    jwt.SigningMethodRS256,  // RSA SHA-256
}

for _, method := range algorithms {
    token, _ := jwt.Generate(method, secret, claims)
    parsed, _ := jwt.Parse(method, token, secret)
}

// 🔥 构建器模式
token, err := jwt.NewBuilder(jwt.SigningMethodHS512, secret).
    SetIssuer("go-util-app").
    SetSubject("user123").
    SetExpirationFromNow(24 * time.Hour).
    SetClaim("role", "admin").
    SetClaim("permissions", []string{"read", "write"}).
    Build()

// 🔥 RSA 密钥对生成
privateKey, publicKey := jwt.GenerateRSAKeyPair(2048)
rsaToken, _ := jwt.Generate(jwt.SigningMethodRS256, privateKey, claims)
```

### 🔢 XEnum - 类型安全枚举
```go
// 🔥 高性能枚举系统
UserStatus := util.NewEnumBuilder[int]().
    Add(0, "INACTIVE", "未激活").
    Add(1, "ACTIVE", "已激活").
    Add(2, "SUSPENDED", "已暂停").
    Build()

// O(1) 快速查找
lookup := UserStatus.NewFastLookup()
if status, exists := lookup.GetByValue(1); exists {
    fmt.Printf("状态: %s (%s)", status.Name(), status.Desc())
}

// 批量验证
validator := UserStatus.NewBatchValidator()
results := validator.ValidateAll([]int{0, 1, 99, 2})
// [true, true, false, true]

// 🔥 GORM 数据库集成
type User struct {
    ID     uint                    `gorm:"primaryKey"`
    Name   string                  `gorm:"size:100"`
    Status *util.XEnum[int]        `gorm:"type:int;index"`
}

// 直接存储和查询
db.Create(&User{Name: "张三", Status: active})
db.Where("status = ?", UserStatus.ACTIVE.Value()).Find(&users)
```

---

## 📖 文档

### 📚 完整指南
- 🚀 **[JSONx 完整使用指南](jsonx/README.md)** - 深度路径、构建器、性能优化
- 🔐 **[JWT 包使用指南](jwt/README.md)** - 认证、安全实践、Web框架集成
- 🔢 **[XEnum 完整使用指南](ENUM_COMPLETE_GUIDE.md)** - 枚举系统、GORM集成、企业实践

### 🎯 使用场景

<details>
<summary><b>📊 数据处理管道</b></summary>

```go
// 复杂的用户数据处理流程
result := util.Array(users).
    Filter(func(u User) bool { 
        return u.IsActive && u.LastLoginAt.After(lastWeek) 
    }).
    Map(func(u User) UserDTO {
        return UserDTO{
            ID:    u.ID,
            Name:  util.Str(u.Name).FirstUpper().String(),
            Email: util.Str(u.Email).Lower().String(),
            Role:  u.Role.Name(),
        }
    }).
    SortBy(func(a, b UserDTO) bool { return a.Name < b.Name }).
    Take(100).
    ToSlice()
```
</details>

<details>
<summary><b>🌐 API 数据处理</b></summary>

```go
// RESTful API 客户端
client := util.Http().
    BaseURL("https://api.github.com").
    Header("Authorization", "token "+githubToken).
    Header("Accept", "application/vnd.github.v3+json").
    Timeout(30 * time.Second).
    WithRetry(util.RetryConfig{
        MaxRetries: 3,
        Delay:      time.Second,
    })

var repos []Repository
if resp, err := client.Get("/user/repos"); err == nil && resp.IsOK() {
    resp.JSON(&repos)
    
    popularRepos := util.Array(repos).
        Filter(func(r Repository) bool { return r.Stars > 100 }).
        SortBy(func(a, b Repository) bool { return a.Stars > b.Stars }).
        Take(10).
        ToSlice()
}
```
</details>

<details>
<summary><b>📁 文件批处理</b></summary>

```go
// 批量图片处理
processed := 0
util.File("./photos").Walk(func(f util.XFile) error {
    if !f.IsFile() || !util.Array([]string{".jpg", ".png", ".gif"}).Contains(f.Ext()) {
        return nil
    }
    
    // 生成新文件名
    timestamp := util.Now().Format("20060102_150405")
    newName := fmt.Sprintf("%s_%03d%s", timestamp, processed, f.Ext())
    
    // 移动并重命名
    if err := f.Move(f.Dir() + "/processed/" + newName); err == nil {
        processed++
        fmt.Printf("处理完成: %s -> %s\n", f.Name(), newName)
    }
    return nil
})
```
</details>

---

## ⚡ 性能基准

### 🚀 JSONx 性能测试
```
深度嵌套测试 (200层):     ✅ 毫秒级完成
大数组处理 (5000对象):    ✅ 0.01秒完成  
序列化 (13万字符):       ✅ < 1毫秒
复杂路径解析:           ✅ O(1) 查找
内存使用:              ✅ 零拷贝优化
```

### 🔐 JWT 性能测试
```
HMAC 签名/验证:         ✅ 10万次/秒
RSA 签名/验证:          ✅ 1万次/秒
大载荷处理 (10KB):      ✅ 毫秒级完成
批量验证:              ✅ 并发安全
```

### 🔢 枚举性能测试
```
O(1) 快速查找:          ✅ 100万次/秒
批量验证 (1000项):      ✅ 微秒级完成
并发访问:              ✅ 无锁设计
内存占用:              ✅ 最小化
```

---

## 🔧 安装与配置

### 环境要求
- **Go 1.23+** (支持泛型和最新特性)
- **现代化开发环境** (VSCode, GoLand 等)

### 依赖管理
```bash
# 核心库 (零第三方依赖)
go get github.com/zhoudm1743/go-util

# JSONx 独立包
go get github.com/zhoudm1743/go-util/jsonx

# JWT 独立包  
go get github.com/zhoudm1743/go-util/jwt

# 可选: 高精度计算支持
go get github.com/shopspring/decimal
go get golang.org/x/exp
```

### IDE 配置
```json
// VSCode settings.json
{
    "go.toolsManagement.autoUpdate": true,
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.formatTool": "goimports"
}
```

---

## 🤝 社区与支持

### 🐛 问题反馈
- **Bug 报告**: [GitHub Issues](https://github.com/zhoudm1743/go-util/issues)
- **功能请求**: [Feature Requests](https://github.com/zhoudm1743/go-util/discussions)
- **安全漏洞**: security@go-util.com

### 📈 贡献指南
```bash
# 1. Fork 项目
git clone https://github.com/your-username/go-util.git

# 2. 创建功能分支
git checkout -b feature/amazing-feature

# 3. 运行测试
go test ./...

# 4. 提交更改
git commit -m "feat: add amazing feature"

# 5. 推送并创建 PR
git push origin feature/amazing-feature
```

### 🎖️ 贡献者
感谢所有为 Go-Util 做出贡献的开发者！

---

## 📄 许可证

本项目采用 **MIT 许可证** - 查看 [LICENSE](LICENSE) 文件了解详情。

```
MIT License - 自由使用、修改、分发
✅ 商业使用   ✅ 修改   ✅ 分发   ✅ 私有使用
```

---

<div align="center">

### 🌟 如果这个项目对你有帮助，请给个 Star！

**让 Go 代码更简洁，让开发更高效！** 🚀

[⭐ Star 项目](https://github.com/zhoudm1743/go-util) • 
[📖 查看文档](https://github.com/zhoudm1743/go-util/wiki) • 
[💬 加入讨论](https://github.com/zhoudm1743/go-util/discussions)

---

*Made with ❤️ by Go developers, for Go developers*

</div> 