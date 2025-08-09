# 🚀 JSONx - 高效的 Go JSON 操作库

JSONx 是一个高效、简单易用的 JSON 操作库，支持链式调用，专注于性能和易用性。

## ✨ 特性

- 🔥 **链式调用** - 流畅的 API 设计，支持方法链
- 🚀 **高性能** - 零拷贝字符串转换，高效的内存使用
- 💡 **简单易用** - 直观的 API，快速上手
- 🛠️ **功能丰富** - 支持路径操作、数组处理、对象合并等
- 🎯 **类型安全** - 内置类型检查和转换
- 📦 **零依赖** - 只使用 Go 标准库
- 🔍 **深度路径** - 支持 `obj.user.profile.name` 形式的路径访问
- 🏗️ **构建器模式** - 支持链式构建复杂 JSON 结构

## 📦 安装

```bash
go get github.com/zhoudm1743/go-util/jsonx
```

## 🚀 快速开始

### 基础操作

```go
package main

import (
    "fmt"
    "github.com/zhoudm1743/go-util/jsonx"
)

func main() {
    // 解析 JSON 字符串
    j := jsonx.Parse(`{"name": "张三", "age": 25, "active": true}`)
    
    // 获取值
    name := j.Get("name").String()        // "张三"
    age := j.Get("age").Int()            // 25
    active := j.Get("active").Bool()      // true
    
    // 链式修改
    j.Set("age", 26).
      Set("email", "zhangsan@example.com").
      Set("tags", []interface{}{"开发者", "Go"})
    
    // 输出 JSON
    result, _ := j.ToPrettyJSON()
    fmt.Println(result)
}
```

### 深度路径操作

```go
// 创建嵌套结构
j := jsonx.Object().
    Set("user.profile.name", "李四").
    Set("user.profile.age", 30).
    Set("user.settings.theme", "dark").
    Set("user.settings.lang", "zh-CN")

// 获取嵌套值
userName := j.Get("user.profile.name").String()  // "李四"
theme := j.Get("user.settings.theme").String()   // "dark"

// 检查路径是否存在
if j.Has("user.profile.email") {
    email := j.Get("user.profile.email").String()
}

// 删除嵌套路径
j.Delete("user.settings.theme")
```

### 构建器模式

```go
// 对象构建器
user := jsonx.NewBuilder().
    AddString("name", "王五").
    AddInt("age", 28).
    AddBool("verified", true).
    AddObject("address", jsonx.NewBuilder().
        AddString("city", "北京").
        AddString("street", "长安街").
        Build()).
    Build()

// 数组构建器
arr := jsonx.NewArrayBuilder().
    AppendString("Go").
    AppendString("Python").
    AppendString("JavaScript").
    Build()

// 快速构建
product := jsonx.QuickObject(map[string]interface{}{
    "id":    "P001",
    "name":  "智能手机",
    "price": 2999.99,
    "tags":  []interface{}{"电子产品", "手机"},
})
```

## 📚 API 文档

### 创建 JSON

```go
// 解析 JSON 字符串
j := jsonx.Parse(`{"key": "value"}`)

// 从字节数组解析
j := jsonx.ParseBytes(jsonBytes)

// 创建空对象
j := jsonx.Object()

// 创建空数组
j := jsonx.Array()

// 从 map 创建
j := jsonx.FromMap(map[string]interface{}{"key": "value"})

// 从 slice 创建
j := jsonx.FromSlice([]interface{}{1, 2, 3})

// 从结构体创建
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
j := jsonx.FromStruct(User{Name: "张三", Age: 25})
```

### 数据访问

```go
// 获取值（支持深度路径）
j.Get("name")                    // 简单路径
j.Get("user.profile.name")       // 深度路径
j.Get("items.0.name")           // 数组索引

// 设置值
j.Set("name", "新名字")
j.Set("user.age", 30)

// 检查路径是否存在
exists := j.Has("user.email")

// 删除路径
j.Delete("user.settings")
```

### 类型检查和转换

```go
// 类型检查
j.IsObject()    // 是否为对象
j.IsArray()     // 是否为数组
j.IsString()    // 是否为字符串
j.IsNumber()    // 是否为数字
j.IsBool()      // 是否为布尔值
j.IsNull()      // 是否为 null

// 类型转换
j.String()      // 转换为字符串
j.Int()         // 转换为整数
j.Int64()       // 转换为 int64
j.Float64()     // 转换为 float64
j.Bool()        // 转换为布尔值
```

### 数组操作

```go
arr := jsonx.Array()

// 添加元素
arr.Append("item1", "item2", "item3")
arr.Prepend("first")

// 获取元素
first := arr.Index(0)
length := arr.Length()

// 删除元素
arr.Remove(1)

// 迭代
arr.ForEach(func(key string, value *jsonx.JSON) bool {
    fmt.Printf("%s: %s\n", key, value.String())
    return true // 继续迭代
})

// 映射
doubled := arr.Map(func(key string, value *jsonx.JSON) interface{} {
    return value.Int() * 2
})

// 过滤
evens := arr.Filter(func(key string, value *jsonx.JSON) bool {
    return value.Int()%2 == 0
})
```

### 对象操作

```go
obj := jsonx.Object()

// 获取键和值
keys := obj.Keys()           // []string
values := obj.Values()       // []*jsonx.JSON

// 迭代对象
obj.ForEach(func(key string, value *jsonx.JSON) bool {
    fmt.Printf("%s: %v\n", key, value.ToInterface())
    return true
})
```

### 序列化

```go
// 转换为 JSON 字符串
jsonStr, err := j.ToJSON()

// 转换为格式化的 JSON
prettyJSON, err := j.ToPrettyJSON()

// 转换为字节数组
jsonBytes, err := j.ToBytes()

// 转换为原始类型
rawData := j.ToInterface()
mapData, err := j.ToMap()           // map[string]interface{}
sliceData, err := j.ToSlice()       // []interface{}
```

### 克隆和合并

```go
// 深度克隆
cloned := j.Clone()

// 浅合并（覆盖相同键）
merged := j1.Merge(j2)

// 深度合并（递归合并对象）
deepMerged := j1.DeepMerge(j2)

// 合并多个对象
result := jsonx.Merge(j1, j2, j3)
result := jsonx.DeepMergeAll(j1, j2, j3)
```

## 🔧 高级功能

### 模板构建器

```go
template := `{
    "user": "{{username}}",
    "message": "{{message}}",
    "timestamp": {{timestamp}},
    "active": {{active}}
}`

j := jsonx.NewTemplate(template).
    Set("username", "张三").
    Set("message", "Hello World").
    Set("timestamp", 1640995200).
    Set("active", true).
    Build()
```

### 扁平化和反扁平化

```go
nested := jsonx.QuickObject(map[string]interface{}{
    "user": map[string]interface{}{
        "profile": map[string]interface{}{
            "name": "张三",
            "age":  30,
        },
    },
})

// 扁平化
flattened := jsonx.Flatten(nested)
// 结果: {"user.profile.name": "张三", "user.profile.age": 30}

// 反扁平化
unflattened := jsonx.Unflatten(flattened)
// 恢复原始嵌套结构
```

### 字段选择和排除

```go
user := jsonx.QuickObject(map[string]interface{}{
    "id":       1,
    "name":     "张三",
    "email":    "zhang@example.com",
    "password": "secret",
    "internal": "data",
})

// 只选择公开字段
public := jsonx.Pick(user, "id", "name", "email")

// 排除敏感字段
safe := jsonx.Omit(user, "password", "internal")
```

### Schema 验证

```go
schema := &jsonx.Schema{
    Type: "object",
    Properties: map[string]*jsonx.Schema{
        "name": {Type: "string", MinLength: intPtr(1)},
        "age":  {Type: "number", Minimum: float64Ptr(0)},
    },
    Required: []string{"name", "age"},
}

user := jsonx.QuickObject(map[string]interface{}{
    "name": "张三",
    "age":  25,
})

if err := schema.Validate(user); err != nil {
    log.Printf("验证失败: %v", err)
}
```

## 🛠️ 实用工具

```go
// JSON 字符串验证
valid := jsonx.IsValid(`{"name": "test"}`)  // true

// 格式化 JSON
pretty, _ := jsonx.Pretty(`{"name":"test"}`)

// 压缩 JSON
compact, _ := jsonx.Minify(prettyJSON)

// 比较两个 JSON
equal := jsonx.Compare(j1, j2)

// 获取 JSON 信息
size := jsonx.Size(j)        // 字节大小
depth := jsonx.Depth(j)      // 嵌套深度
jsonType := jsonx.GetType(j) // 类型名称
```

## 🔥 链式调用示例

```go
// 复杂的链式操作
result := jsonx.Object().
    Set("user.name", "张三").
    Set("user.age", 25).
    Set("user.tags", []interface{}{"开发者", "Go"}).
    Get("user").
    Set("verified", true).
    Set("last_login", time.Now().Unix()).
    Clone().
    Merge(jsonx.QuickObject(map[string]interface{}{
        "preferences": map[string]interface{}{
            "theme": "dark",
            "lang":  "zh-CN",
        },
    }))

// 数组链式操作
numbers := jsonx.Array().
    Append(1, 2, 3, 4, 5).
    Filter(func(key string, value *jsonx.JSON) bool {
        return value.Int()%2 == 0  // 过滤偶数
    }).
    Map(func(key string, value *jsonx.JSON) interface{} {
        return value.Int() * 10    // 乘以 10
    })
```

## ⚡ 性能优化

JSONx 包含多项性能优化：

- **零拷贝字符串转换** - 使用 `unsafe` 包进行高效转换
- **内存复用** - 减少不必要的内存分配
- **快速路径访问** - 优化的路径解析算法
- **类型断言缓存** - 减少重复的类型检查

## 🧪 测试

```bash
# 运行所有测试
go test -v

# 运行基准测试
go test -bench=.

# 生成测试覆盖率报告
go test -cover
```

## 📝 示例

查看 `example/main.go` 文件了解完整的使用示例。

```bash
cd example
go run main.go
```

## 🤝 贡献

欢迎提交 Pull Request 和 Issue！

## 📄 许可证

本项目采用 MIT 许可证。详细信息请查看 [LICENSE](LICENSE) 文件。

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

---

**让 JSON 操作更简单，让 Go 开发更高效！** 🚀 