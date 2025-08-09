# XEnum 完整使用指南

XEnum 是一个简易、高可用、高性能的泛型枚举实现，支持类型安全的枚举操作，完全兼容 GORM 数据库操作。

## 📖 目录

1. [特性概览](#-特性概览)
2. [基础用法](#-基础用法)
3. [查找和验证](#-查找和验证)
4. [枚举操作](#-枚举操作)
5. [高级功能](#-高级功能)
6. [GORM 数据库集成](#-gorm-数据库集成)
7. [企业级项目集成](#️-企业级项目集成)
8. [高级查询模式](#-高级查询模式)
9. [性能优化策略](#-性能优化策略)
10. [生产环境最佳实践](#️-生产环境最佳实践)
11. [实际应用场景](#-实际应用场景)
12. [最佳实践总结](#-最佳实践总结)

## 🚀 特性概览

- **类型安全** - 基于 Go 泛型，编译时类型检查
- **高性能** - 内置快速查找和批量验证机制
- **并发安全** - 内置读写锁，支持并发访问
- **JSON 支持** - 内置 JSON 序列化/反序列化
- **GORM 兼容** - 完全兼容 GORM，支持直接在数据库模型中使用
- **简易易用** - 链式 API 设计，简洁的定义方式
- **功能丰富** - 支持枚举集合、范围查询、批量操作等

## 📖 基础用法

### 1. 创建简单枚举

```go
package main

import (
    "fmt"
    util "github.com/zhoudm1743/go-util"
)

func main() {
    // 创建状态枚举
    Status := util.NewIntEnum()
    
    // 定义枚举值
    inactive := Status.DefineSimple(0, "INACTIVE")
    active := Status.DefineSimple(1, "ACTIVE")
    pending := Status.DefineSimple(2, "PENDING")
    
    fmt.Printf("状态: %s (值: %d)\n", active.Name(), active.Value())
    // 输出: 状态: ACTIVE (值: 1)
}
```

### 2. 使用构建器模式

```go
// 链式定义枚举
UserRole := util.NewEnumBuilder[int]().
    Add(1, "ADMIN", "管理员").
    Add(2, "USER", "普通用户").
    Add(3, "GUEST", "访客").
    Build()

// 获取枚举
admin, _ := UserRole.FromValue(1)
fmt.Printf("角色: %s - %s\n", admin.Name(), admin.Desc())
// 输出: 角色: ADMIN - 管理员
```

### 3. 使用便捷定义方法

```go
// 快速定义整数枚举
Priority := util.DefineEnum(map[int]string{
    1: "LOW",
    2: "MEDIUM", 
    3: "HIGH",
    4: "URGENT",
})

// 快速定义带描述的枚举
Level := util.DefineEnumWithDesc(map[string][2]string{
    "DEBUG": {"DEBUG", "调试级别"},
    "INFO":  {"INFO", "信息级别"},
    "WARN":  {"WARN", "警告级别"},
    "ERROR": {"ERROR", "错误级别"},
})
```

## 🔍 查找和验证

### 基础查找

```go
// 根据值查找
if status, exists := Status.FromValue(1); exists {
    fmt.Printf("找到状态: %s\n", status.Name())
}

// 根据名称查找
if status, exists := Status.FromName("ACTIVE"); exists {
    fmt.Printf("状态值: %d\n", status.Value())
}

// 验证值的有效性
isValid := Status.IsValid(1)        // true
isValidName := Status.IsValidName("ACTIVE") // true
```

### 高性能查找

```go
// 创建快速查找器（适用于频繁查找的场景）
lookup := Status.NewFastLookup()

// 快速查找（无锁操作）
if status, exists := lookup.GetByValue(1); exists {
    fmt.Printf("快速查找: %s\n", status.Name())
}
```

### 批量验证

```go
// 创建批量验证器
validator := Status.NewBatchValidator()

// 批量验证
values := []int{1, 2, 99, 3}
results := validator.ValidateAll(values)
fmt.Printf("验证结果: %v\n", results) // [true true false true]

// 过滤有效值
validValues := validator.FilterValid(values)
fmt.Printf("有效值: %v\n", validValues) // [1 2 3]
```

## 📊 枚举操作

### 获取枚举信息

```go
// 获取所有枚举值
all := Status.All()
for _, enum := range all {
    fmt.Printf("%s: %d\n", enum.Name(), enum.Value())
}

// 获取所有值和名称
values := Status.Values()   // [0, 1, 2]
names := Status.Names()     // ["INACTIVE", "ACTIVE", "PENDING"]
count := Status.Count()     // 3
```

### 枚举比较

```go
active, _ := Status.FromName("ACTIVE")
pending, _ := Status.FromName("PENDING")

// 枚举比较
fmt.Printf("相等: %t\n", active.Equal(pending))           // false
fmt.Printf("值相等: %t\n", active.EqualValue(1))         // true
fmt.Printf("名称相等: %t\n", active.EqualName("ACTIVE"))  // true
```

### 范围查询

```go
// 获取指定范围内的枚举
rangeEnums := util.EnumRange(Priority, 2, 4)
for _, enum := range rangeEnums {
    fmt.Printf("优先级: %s\n", enum.Name())
}
// 输出: MEDIUM, HIGH, URGENT
```

## 🔧 高级功能

### 枚举集合

```go
// 创建枚举集合
activeStatuses := util.NewEnumSet[int]()

// 添加多个状态
active, _ := Status.FromName("ACTIVE")
pending, _ := Status.FromName("PENDING")

activeStatuses.Add(active).Add(pending)

// 检查集合
fmt.Printf("集合大小: %d\n", activeStatuses.Size())
fmt.Printf("包含ACTIVE: %t\n", activeStatuses.Contains(active))
fmt.Printf("包含值1: %t\n", activeStatuses.ContainsValue(1))
```

### JSON 序列化

```go
type Config struct {
    LogLevel *util.XEnum[string] `json:"log_level"`
    Status   *util.XEnum[int]    `json:"status"`
}

config := Config{
    LogLevel: util.MustGetEnum(LogLevel, "INFO"),
    Status:   util.MustGetEnum(Status, 1),
}

// 序列化为 JSON
jsonData, _ := json.Marshal(config)
fmt.Printf("JSON: %s\n", string(jsonData))
// 输出: {"log_level":"INFO","status":"ACTIVE"}
```

### 字符串解析

```go
// 根据字符串解析枚举（支持名称和值）
func parseLogLevel(levelStr string) (*util.XEnum[string], error) {
    return util.ParseEnum(LogLevel, levelStr)
}

// 使用示例
if level, err := parseLogLevel("INFO"); err == nil {
    fmt.Printf("日志级别: %s - %s\n", level.Name(), level.Desc())
}
```

## 🗄️ GORM 数据库集成

XEnum 完全兼容 GORM，支持直接在数据库模型中使用枚举类型。

### 基本使用

```go
import (
    util "github.com/zhoudm1743/go-util"
    "gorm.io/gorm"
)

// 定义用户状态枚举
var UserStatus = util.DefineGormEnum("user_status", map[int]string{
    0: "INACTIVE",
    1: "ACTIVE", 
    2: "PENDING",
    3: "SUSPENDED",
})

// 用户模型
type User struct {
    ID     uint              `gorm:"primaryKey"`
    Name   string            `gorm:"size:100"`
    Email  string            `gorm:"uniqueIndex"`
    Status *util.XEnum[int]  `gorm:"type:int;default:0;index"` // 直接使用枚举
}

func main() {
    // 自动迁移
    db.AutoMigrate(&User{})
    
    // 创建用户
    activeStatus, _ := UserStatus.FromValue(1) // ACTIVE
    user := User{
        Name:   "张三",
        Email:  "zhangsan@example.com", 
        Status: activeStatus,
    }
    
    db.Create(&user)  // 枚举值自动存储为数据库原始值
    fmt.Printf("创建用户: %s, 状态: %s\n", user.Name, user.Status.Name())
}
```

### 数据库查询

```go
// 查询指定状态的用户
query, value := util.CreateEnumQuery("status", activeStatus)
db.Where(query, value).Find(&users)

// IN 查询
query, values := util.CreateEnumInQuery("status", []*util.XEnum[int]{active, pending})
db.Where(query, values).Find(&users)

// 复杂查询
var users []User
active, _ := UserStatus.FromName("ACTIVE")
suspended, _ := UserStatus.FromName("SUSPENDED")

db.Where("email LIKE ?", "%@company.com").
   Where("status IN ?", []int{active.Value(), suspended.Value()}).
   Find(&users)
```

### 使用 EnumField 包装器

```go
// 增强的用户模型，支持自动验证
type UserWithValidation struct {
    ID     uint                         `gorm:"primaryKey"`
    Name   string                       `gorm:"size:100"`
    Status *util.EnumField[int]         `gorm:"type:int"`
}

// 创建时自动设置枚举字段
func CreateUserWithValidation(db *gorm.DB, name string, statusValue int) error {
    user := &UserWithValidation{
        Name:   name,
        Status: util.NewEnumField(UserStatus),
    }
    
    // 设置状态值，自动验证
    if err := user.Status.SetValue(statusValue); err != nil {
        return fmt.Errorf("invalid status: %v", err)
    }
    
    return db.Create(user).Error
}
```

## 🏗️ 企业级项目集成

### 项目结构建议

```
project/
├── enums/
│   ├── user.go          // 用户相关枚举
│   ├── order.go         // 订单相关枚举
│   └── common.go        // 通用枚举
├── models/
│   ├── user.go          // 用户模型
│   └── order.go         // 订单模型
├── repositories/
│   ├── user_repository.go
│   └── order_repository.go
├── services/
│   ├── user_service.go
│   └── order_service.go
└── main.go
```

### 枚举定义最佳实践

```go
// enums/user.go
package enums

import util "github.com/zhoudm1743/go-util"

// 用户状态枚举 - 使用常量确保类型安全
const (
    UserStatusInactive  = 0
    UserStatusActive    = 1
    UserStatusPending   = 2
    UserStatusSuspended = 3
    UserStatusDeleted   = 4
)

var UserStatus = util.DefineGormEnumWithDesc("user_status", map[int][2]string{
    UserStatusInactive:  {"INACTIVE", "非活跃"},
    UserStatusActive:    {"ACTIVE", "活跃"},
    UserStatusPending:   {"PENDING", "待处理"},
    UserStatusSuspended: {"SUSPENDED", "暂停"},
    UserStatusDeleted:   {"DELETED", "已删除"},
})

// 提供便捷的构造函数
func NewUserStatus(value int) (*util.XEnum[int], error) {
    return util.ConvertToEnum(UserStatus, value)
}

func MustUserStatus(value int) *util.XEnum[int] {
    return util.MustGetEnum(UserStatus, value)
}
```

### 模型定义

```go
// models/user.go
package models

import (
    "time"
    "your-project/enums"
    util "github.com/zhoudm1743/go-util"
    "gorm.io/gorm"
)

type User struct {
    ID        uint              `gorm:"primaryKey"`
    Username  string            `gorm:"uniqueIndex;size:50;not null"`
    Email     string            `gorm:"uniqueIndex;size:100;not null"`
    Password  string            `gorm:"size:255;not null"`
    Status    *util.XEnum[int]  `gorm:"type:int;not null;default:0;index;comment:用户状态"`
    Role      *util.XEnum[int]  `gorm:"type:int;not null;default:4;index;comment:用户角色"`
    LastLogin *time.Time        `gorm:"comment:最后登录时间"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt    `gorm:"index"`
}

// GORM 钩子：创建前验证
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // 设置默认状态
    if u.Status == nil {
        u.Status = enums.MustUserStatus(enums.UserStatusInactive)
    }
    
    // 验证枚举值
    if err := util.ValidateEnumField(enums.UserStatus, u.Status.Value()); err != nil {
        return err
    }
    
    return nil
}
```

## 🔍 高级查询模式

### 仓储模式 (Repository Pattern)

```go
// repositories/user_repository.go
package repositories

import (
    "your-project/enums"
    "your-project/models"
    util "github.com/zhoudm1743/go-util"
    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

// 按状态查询用户
func (r *UserRepository) FindByStatus(status *util.XEnum[int]) ([]models.User, error) {
    var users []models.User
    query, value := util.CreateEnumQuery("status", status)
    err := r.db.Where(query, value).Find(&users).Error
    return users, err
}

// 状态统计
func (r *UserRepository) GetStatusStatistics() (map[string]int64, error) {
    stats := make(map[string]int64)
    
    for _, status := range enums.UserStatus.All() {
        var count int64
        if err := r.db.Model(&models.User{}).
            Where("status = ?", status.Value()).Count(&count).Error; err != nil {
            return nil, err
        }
        stats[status.Name()] = count
    }
    
    return stats, nil
}

// 批量更新状态
func (r *UserRepository) BatchUpdateStatus(userIDs []uint, newStatus *util.XEnum[int]) error {
    return r.db.Model(&models.User{}).
        Where("id IN ?", userIDs).
        Update("status", newStatus.Value()).Error
}
```

### 查询构建器模式

```go
type UserQueryBuilder struct {
    db *gorm.DB
}

func (r *UserRepository) Query() *UserQueryBuilder {
    return &UserQueryBuilder{db: r.db}
}

func (q *UserQueryBuilder) WithStatus(status *util.XEnum[int]) *UserQueryBuilder {
    if status != nil {
        q.db = q.db.Where("status = ?", status.Value())
    }
    return q
}

func (q *UserQueryBuilder) ActiveOnly() *UserQueryBuilder {
    return q.WithStatus(enums.MustUserStatus(enums.UserStatusActive))
}

func (q *UserQueryBuilder) Find() ([]models.User, error) {
    var users []models.User
    err := q.db.Find(&users).Error
    return users, err
}

// 使用示例
func ExampleUsage(repo *UserRepository) {
    // 查找活跃用户
    activeUsers, _ := repo.Query().ActiveOnly().Find()
    
    // 统计活跃用户数量
    count, _ := repo.Query().ActiveOnly().Count()
}
```

## 🚀 性能优化策略

### 1. 数据库层面优化

```sql
-- 创建合适的索引
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status_role ON users(status, role);  -- 复合索引

-- 为频繁查询的组合创建覆盖索引
CREATE INDEX idx_users_lookup ON users(status, role, id, username);
```

### 2. 应用层面优化

```go
// services/user_service.go
package services

import (
    "sync"
    "your-project/enums"
    util "github.com/zhoudm1743/go-util"
)

// 全局预加载的查找表
var (
    userStatusLookup *util.FastLookup[int]
    lookupOnce       sync.Once
)

// 初始化查找表
func initLookups() {
    lookupOnce.Do(func() {
        userStatusLookup = enums.UserStatus.NewFastLookup()
    })
}

type UserService struct {
    repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
    initLookups()
    return &UserService{repo: repo}
}

// 高性能状态名称获取
func (s *UserService) GetStatusName(statusValue int) string {
    if status, exists := userStatusLookup.GetByValue(statusValue); exists {
        return status.Name()
    }
    return "UNKNOWN"
}
```

### 3. 连接池和事务优化

```go
// config/database.go
func SetupDatabase() (*gorm.DB, error) {
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        PrepareStmt: true,  // 预编译语句，提高性能
    })
    if err != nil {
        return nil, err
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }

    // 连接池配置
    sqlDB.SetMaxOpenConns(100)           // 最大打开连接数
    sqlDB.SetMaxIdleConns(10)            // 最大空闲连接数
    sqlDB.SetConnMaxLifetime(time.Hour)  // 连接最大生存时间

    return db, nil
}

// 事务中的枚举操作
func (s *UserService) UpdateUserStatusInTransaction(userID uint, newStatus int) error {
    return s.repo.db.Transaction(func(tx *gorm.DB) error {
        // 验证状态
        if err := util.ValidateEnumField(enums.UserStatus, newStatus); err != nil {
            return err
        }

        // 状态转换验证
        // ... 业务逻辑

        // 更新状态
        return tx.Model(&models.User{}).Where("id = ?", userID).
            Update("status", newStatus).Error
    })
}
```

## 🛡️ 生产环境最佳实践

### 1. 错误处理和日志

```go
// utils/enum_logger.go
package typess

import (
    "log/slog"
    "context"
)

type EnumLogger struct {
    logger *slog.Logger
}

func NewEnumLogger() *EnumLogger {
    return &EnumLogger{
        logger: slog.Default(),
    }
}

func (l *EnumLogger) LogEnumValidationError(ctx context.Context, field string, value interface{}, err error) {
    l.logger.ErrorContext(ctx, "Enum validation failed",
        slog.String("field", field),
        slog.Any("value", value),
        slog.String("error", err.Error()),
    )
}

func (l *EnumLogger) LogStatusTransition(ctx context.Context, userID uint, from, to string) {
    l.logger.InfoContext(ctx, "User status transition",
        slog.Uint64("user_id", uint64(userID)),
        slog.String("from_status", from),
        slog.String("to_status", to),
    )
}
```

### 2. 监控和指标

```go
// monitoring/enum_metrics.go
package monitoring

import (
    "time"
    "your-project/enums"
)

type EnumMetrics struct {
    statusCounts    map[int]int64
    lastUpdate      time.Time
    updateInterval  time.Duration
}

func NewEnumMetrics() *EnumMetrics {
    return &EnumMetrics{
        statusCounts:   make(map[int]int64),
        updateInterval: 1 * time.Minute,
    }
}

// 定期更新指标
func (m *EnumMetrics) UpdateMetrics(repo *repositories.UserRepository) error {
    if time.Since(m.lastUpdate) < m.updateInterval {
        return nil
    }

    stats, err := repo.GetStatusStatistics()
    if err != nil {
        return err
    }

    for _, status := range enums.UserStatus.All() {
        count := stats[status.Name()]
        m.statusCounts[status.Value()] = count
        
        // 导出到监控系统（如 Prometheus）
        // prometheus.UserStatusGauge.WithLabelValues(status.Name()).Set(float64(count))
    }

    m.lastUpdate = time.Now()
    return nil
}
```

### 3. 配置管理

```go
// config/enum_config.go
package config

type EnumConfig struct {
    EnableValidation     bool
    EnableStatusLogging  bool
    CacheTTL            int // seconds
    MaxBatchSize        int
}

func LoadEnumConfig() *EnumConfig {
    return &EnumConfig{
        EnableValidation:    getEnvBool("ENUM_ENABLE_VALIDATION", true),
        EnableStatusLogging: getEnvBool("ENUM_ENABLE_LOGGING", true),
        CacheTTL:           getEnvInt("ENUM_CACHE_TTL", 300),
        MaxBatchSize:       getEnvInt("ENUM_MAX_BATCH_SIZE", 1000),
    }
}
```

## 🎯 实际应用场景

### 1. 用户状态管理

```go
// 定义用户状态枚举
UserStatus := util.CreateStatusEnum() // 内置的状态枚举

type User struct {
    ID     int                `json:"id"`
    Name   string             `json:"name"`
    Status *util.XEnum[int]   `json:"status"`
}

// 创建用户
user := User{
    ID:     1,
    Name:   "张三",
    Status: util.MustGetEnum(UserStatus, 1), // ACTIVE
}

// 状态判断
if user.Status.EqualName("ACTIVE") {
    fmt.Println("用户已激活")
}
```

### 2. 电商订单系统

```go
// 订单状态枚举
var OrderStatus = util.DefineGormEnumWithDesc("order_status", map[int][2]string{
    1: {"PENDING", "待支付"},
    2: {"PAID", "已支付"},
    3: {"SHIPPED", "已发货"},
    4: {"DELIVERED", "已送达"},
    5: {"CANCELLED", "已取消"},
    6: {"REFUNDED", "已退款"},
})

// 订单模型
type Order struct {
    ID          uint              `gorm:"primaryKey"`
    OrderNumber string            `gorm:"uniqueIndex;size:50"`
    UserID      uint              `gorm:"index"`
    Amount      float64           `gorm:"type:decimal(10,2)"`
    Status      *util.XEnum[int]  `gorm:"type:int;index"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// 订单状态流转
func UpdateOrderStatus(db *gorm.DB, orderID uint, newStatusValue int) error {
    // 验证新状态是否有效
    if err := util.ValidateEnumField(OrderStatus, newStatusValue); err != nil {
        return err
    }
    
    // 获取枚举实例
    newStatus, _ := OrderStatus.FromValue(newStatusValue)
    
    // 更新订单状态
    return db.Model(&Order{}).
        Where("id = ?", orderID).
        Update("status", newStatus.Value()).Error
}
```

### 3. 权限管理系统

```go
// 用户角色枚举
var UserRole = util.DefineGormEnumWithDesc("user_role", map[int][2]string{
    1: {"SUPER_ADMIN", "超级管理员"},
    2: {"ADMIN", "管理员"},
    3: {"MODERATOR", "版主"},
    4: {"USER", "普通用户"},
    5: {"GUEST", "访客"},
})

// 权限检查
func HasPermission(db *gorm.DB, userID uint, requiredLevel int) (bool, error) {
    var user models.User
    err := db.First(&user, userID).Error
    if err != nil {
        return false, err
    }
    
    // 根据角色检查权限
    switch user.Role.Name() {
    case "SUPER_ADMIN":
        return true, nil
    case "ADMIN":
        return requiredLevel <= 4, nil
    case "MODERATOR":
        return requiredLevel <= 3, nil
    default:
        return requiredLevel <= 1, nil
    }
}
```

## 🏆 最佳实践总结

### ✅ 核心优势

- **类型安全**: 编译时检查，避免运行时枚举错误
- **高性能**: O(1) 查找，预编译查找表，优化的数据库操作
- **易维护**: 清晰的项目结构，标准化的枚举定义
- **生产就绪**: 完整的错误处理、日志记录和监控支持

### 🎯 设计原则

1. **项目结构**: 分离枚举定义、模型和仓储层
2. **性能优化**: 使用索引、连接池、预加载查找表
3. **错误处理**: 完善的验证和日志记录机制
4. **监控**: 实时跟踪枚举状态分布和转换

### 📋 开发检查清单

- [ ] 使用常量定义枚举值，确保类型安全
- [ ] 为频繁查询的枚举字段创建数据库索引
- [ ] 实现 GORM 钩子进行枚举值验证
- [ ] 使用预加载查找表优化频繁访问
- [ ] 添加完善的错误处理和日志记录
- [ ] 实现枚举状态监控和指标收集
- [ ] 编写单元测试覆盖枚举操作
- [ ] 文档化枚举定义和状态转换规则

### 🔄 持续改进

- 定期审查枚举定义的合理性
- 监控数据库查询性能
- 根据业务需求调整缓存策略
- 收集枚举使用指标，持续优化

---

通过遵循这些实践，您可以构建出高性能、类型安全、易维护的枚举系统，让 Go 代码更加简洁优雅！

**让枚举更简单，让开发更高效！** 🚀 