# JWT 包使用指南

一个轻量级、高性能的 JWT (JSON Web Token) 实现，支持多种签名算法。

## 🚀 特性

- **零依赖**: 仅使用 Go 标准库实现
- **多算法支持**: 支持 HS256, HS384, HS512, RS256, RS384, RS512
- **类型安全**: 完整的类型定义和错误处理
- **高性能**: 优化的编码/解码实现
- **易于使用**: 提供链式调用的构建器模式
- **标准兼容**: 完全符合 RFC 7519 标准

## 📦 安装

```go
import "github.com/zhoudm1743/go-util/jwt"
```

## 🛠️ 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/zhoudm1743/go-util/jwt"
)

func main() {
    // 🔥 推荐: 使用通用 API (可轻松切换算法)
    secret := []byte("your-secret-key")
    
    // 创建声明
    claims := jwt.MapClaims{
        "sub": "1234567890",
        "name": "John Doe",
        "iat": time.Now().Unix(),
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    }
    
    // 生成 JWT - 算法作为参数传入
    tokenString, err := jwt.Generate(jwt.SigningMethodHS256, secret, claims)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("生成的 JWT: %s\n", tokenString)
    
    // 解析 JWT - 算法作为参数传入
    token, err := jwt.Parse(jwt.SigningMethodHS256, tokenString, secret)
    if err != nil {
        log.Fatal(err)
    }
    
    // 提取声明
    if claims, ok := jwt.ExtractClaims(token); ok {
        if name, exists := jwt.GetClaimString(claims, "name"); exists {
            fmt.Printf("用户名: %s\n", name)
        }
    }
    
    // 传统方式 (仍然支持)
    tokenString2, err := jwt.GenerateHS256(secret, claims)
    token2, err := jwt.ParseHS256(tokenString2, secret)
}
```

### 使用构建器模式

```go
func main() {
    secret := []byte("your-secret-key")
    
    // 使用构建器创建 JWT
    tokenString, err := jwt.NewBuilder(jwt.SigningMethodHS256, secret).
        SetIssuer("your-app").
        SetSubject("user123").
        SetAudience("your-audience").
        SetExpirationFromNow(time.Hour * 24).
        SetIssuedNow().
        SetClaim("role", "admin").
        SetClaim("permissions", []string{"read", "write"}).
        Build()
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("JWT: %s\n", tokenString)
    
    // 解析构建器生成的令牌
    token, err := jwt.Parse(jwt.SigningMethodHS256, tokenString, secret)
    if err != nil {
        log.Fatal(err)
    }
    
    // 轻松切换算法 - 只需修改 SigningMethod
    hs512Token, err := jwt.NewBuilder(jwt.SigningMethodHS512, secret).
        SetSubject("user456").
        SetExpirationFromNow(time.Hour * 2).
        Build()
}
```

## 🔐 支持的签名算法

### HMAC 算法

```go
secret := []byte("your-secret-key")

// 🔥 推荐: 使用通用 API
algorithms := []jwt.SigningMethod{
    jwt.SigningMethodHS256,
    jwt.SigningMethodHS384,
    jwt.SigningMethodHS512,
}

for _, method := range algorithms {
    tokenString, err := jwt.Generate(method, secret, claims)
    token, err := jwt.Parse(method, tokenString, secret)
    fmt.Printf("算法 %s: %s\n", method.Alg(), tokenString[:20])
}

// 传统方式 (仍然支持)
// HS256
tokenString, err := jwt.GenerateHS256(secret, claims)
token, err := jwt.ParseHS256(tokenString, secret)

// HS384
tokenString, err := jwt.GenerateHS384(secret, claims)
token, err := jwt.ParseHS384(tokenString, secret)

// HS512
tokenString, err := jwt.GenerateHS512(secret, claims)
token, err := jwt.ParseHS512(tokenString, secret)
```

### RSA 算法

```go
// 生成 RSA 密钥对
privateKey, err := jwt.GenerateRSAKeyPair(2048)
if err != nil {
    log.Fatal(err)
}

// 转换为 PEM 格式
privatePEM := jwt.PrivateKeyToPEM(privateKey)
publicPEM, err := jwt.PublicKeyToPEM(&privateKey.PublicKey)
if err != nil {
    log.Fatal(err)
}

// 使用 RSA 私钥生成 JWT
tokenString, err := jwt.GenerateRS256(privatePEM, claims)
if err != nil {
    log.Fatal(err)
}

// 使用 RSA 公钥验证 JWT
token, err := jwt.ParseRS256(tokenString, publicPEM)
if err != nil {
    log.Fatal(err)
}
```

## 📋 声明管理

### 标准声明

```go
// 使用标准声明结构
claims := &jwt.StandardClaims{
    Issuer:    "your-app",
    Subject:   "user123",
    Audience:  "your-audience",
    ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
    IssuedAt:  time.Now().Unix(),
    NotBefore: time.Now().Unix(),
    ID:        "token-id-123",
}

tokenString, err := jwt.GenerateHS256(secret, claims)
```

### 自定义声明

```go
// 使用 MapClaims 添加自定义声明
claims := jwt.MapClaims{
    // 标准声明
    "iss": "your-app",
    "sub": "user123",
    "exp": time.Now().Add(time.Hour * 24).Unix(),
    
    // 自定义声明
    "role": "admin",
    "permissions": []string{"read", "write", "delete"},
    "department": "engineering",
    "level": 5,
}

tokenString, err := jwt.GenerateHS256(secret, claims)
```

## 🔍 令牌验证

### 基本验证

```go
// 验证令牌并提取声明
token, err := jwt.ParseHS256(tokenString, secret)
if err != nil {
    switch err {
    case jwt.ErrTokenExpired:
        fmt.Println("令牌已过期")
    case jwt.ErrTokenNotYetValid:
        fmt.Println("令牌还未生效")
    case jwt.ErrInvalidSignature:
        fmt.Println("无效的签名")
    default:
        fmt.Printf("验证失败: %v\n", err)
    }
    return
}

if token.Valid {
    fmt.Println("令牌验证成功")
}
```

### 高级验证

```go
// 创建 JWT 实例进行高级验证
j := jwt.New(jwt.SigningMethodHS256, secret)

// 使用自定义声明类型解析
claims := jwt.MapClaims{}
token, err := j.ParseWithClaims(tokenString, claims)
if err != nil {
    log.Fatal(err)
}

// 验证标准声明
err = jwt.ValidateStandardClaims(claims, "your-audience", "your-app", "user123")
if err != nil {
    log.Fatal(err)
}

// 验证自定义声明
if role, exists := jwt.GetClaimString(claims, "role"); exists {
    if role != "admin" {
        log.Fatal("权限不足")
    }
}
```

## 🛡️ 安全实践

### 密钥管理

```go
// 生成安全的 HMAC 密钥
secret, err := jwt.GenerateHMACSecret(32) // 256 位
if err != nil {
    log.Fatal(err)
}

// 生成 RSA 密钥对
privateKey, err := jwt.GenerateRSAKeyPair(2048)
if err != nil {
    log.Fatal(err)
}
```

### 令牌生命周期

```go
// 设置合理的过期时间
builder := jwt.NewBuilder(jwt.SigningMethodHS256, secret).
    SetIssuedNow().
    SetExpirationFromNow(time.Hour * 2). // 2小时后过期
    SetNotBefore(time.Now())             // 立即生效

// 对于敏感操作，使用更短的过期时间
sensitiveToken, err := jwt.NewBuilder(jwt.SigningMethodHS256, secret).
    SetSubject("user123").
    SetClaim("action", "password-reset").
    SetExpirationFromNow(time.Minute * 15). // 15分钟过期
    Build()
```

### 令牌刷新

```go
// 实现令牌刷新机制
func RefreshToken(oldTokenString string, secret []byte) (string, error) {
    // 解析旧令牌（可能已过期）
    token, err := jwt.ParseHS256(oldTokenString, secret)
    if err != nil && err != jwt.ErrTokenExpired {
        return "", err
    }
    
    // 提取声明
    claims, ok := jwt.ExtractClaims(token)
    if !ok {
        return "", jwt.ErrInvalidToken
    }
    
    // 创建新令牌
    newClaims := jwt.MapClaims{
        "sub": claims["sub"],
        "iat": time.Now().Unix(),
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    }
    
    return jwt.GenerateHS256(secret, newClaims)
}
```

## 🔧 实用工具

### 令牌信息提取

```go
// 不验证签名，仅解码头部和声明
header, err := jwt.DecodeHeader(tokenString)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("算法: %s\n", header.Algorithm)

claims, err := jwt.DecodeClaims(tokenString)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("声明: %+v\n", claims)
```

### 时间工具

```go
// 检查时间戳
exp := int64(1640995200) // 示例时间戳

if jwt.IsExpired(exp) {
    fmt.Println("已过期")
}

if jwt.IsNotYetValid(exp) {
    fmt.Println("还未生效")
}

// 时间转换
now := time.Now()
unix := jwt.TimeToUnix(now)
backToTime := jwt.UnixToTime(unix)
```

## 🌐 Web 框架集成

### HTTP 中间件示例

```go
import (
    "net/http"
    "strings"
)

func JWTMiddleware(secret []byte) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 从 Authorization 头部获取令牌
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
                return
            }
            
            // 提取 Bearer 令牌
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            if tokenString == authHeader {
                http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
                return
            }
            
            // 验证令牌
            token, err := jwt.ParseHS256(tokenString, secret)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }
            
            // 将声明添加到上下文
            claims, _ := jwt.ExtractClaims(token)
            ctx := context.WithValue(r.Context(), "claims", claims)
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// 使用中间件
func main() {
    secret := []byte("your-secret-key")
    
    mux := http.NewServeMux()
    mux.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
        claims := r.Context().Value("claims").(jwt.MapClaims)
        userID, _ := jwt.GetClaimString(claims, "sub")
        fmt.Fprintf(w, "Hello, user %s!", userID)
    })
    
    // 应用中间件
    handler := JWTMiddleware(secret)(mux)
    http.ListenAndServe(":8080", handler)
}
```

### 登录接口示例

```go
func LoginHandler(secret []byte) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 验证用户凭据（示例）
        username := r.FormValue("username")
        password := r.FormValue("password")
        
        if !validateCredentials(username, password) {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }
        
        // 生成 JWT
        tokenString, err := jwt.NewBuilder(jwt.SigningMethodHS256, secret).
            SetIssuer("your-app").
            SetSubject(username).
            SetExpirationFromNow(time.Hour * 24).
            SetIssuedNow().
            SetClaim("role", getUserRole(username)).
            Build()
        
        if err != nil {
            http.Error(w, "Token generation failed", http.StatusInternalServerError)
            return
        }
        
        // 返回令牌
        response := map[string]string{
            "token": tokenString,
            "type":  "Bearer",
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}
```

## ⚡ 性能优化

### 预编译签名方法

```go
// 在应用启动时预编译签名方法
var (
    hmacMethod = jwt.SigningMethodHS256
    rsaMethod  = jwt.SigningMethodRS256
)

// 重复使用 JWT 实例
var (
    hmacJWT = jwt.New(hmacMethod, secret)
    rsaJWT  = jwt.New(rsaMethod, publicKey)
)

// 在处理请求时重复使用
func ParseToken(tokenString string) (*jwt.Token, error) {
    return hmacJWT.Parse(tokenString)
}
```

### 批量操作

```go
// 批量生成令牌
func GenerateTokensForUsers(userIDs []string, secret []byte) (map[string]string, error) {
    tokens := make(map[string]string, len(userIDs))
    
    for _, userID := range userIDs {
        token, err := jwt.NewBuilder(jwt.SigningMethodHS256, secret).
            SetSubject(userID).
            SetExpirationFromNow(time.Hour * 24).
            SetIssuedNow().
            Build()
        
        if err != nil {
            return nil, err
        }
        
        tokens[userID] = token
    }
    
    return tokens, nil
}
```

## 🧪 测试

### 单元测试示例

```go
package jwt_test

import (
    "testing"
    "time"
    "github.com/zhoudm1743/go-util/jwt"
)

func TestJWTGeneration(t *testing.T) {
    secret := []byte("test-secret")
    claims := jwt.MapClaims{
        "sub": "test-user",
        "exp": time.Now().Add(time.Hour).Unix(),
    }
    
    // 生成令牌
    tokenString, err := jwt.GenerateHS256(secret, claims)
    if err != nil {
        t.Fatalf("Failed to generate token: %v", err)
    }
    
    // 验证令牌
    token, err := jwt.ParseHS256(tokenString, secret)
    if err != nil {
        t.Fatalf("Failed to parse token: %v", err)
    }
    
    if !token.Valid {
        t.Error("Token should be valid")
    }
    
    // 验证声明
    parsedClaims, ok := jwt.ExtractClaims(token)
    if !ok {
        t.Error("Failed to extract claims")
    }
    
    if sub, exists := jwt.GetClaimString(parsedClaims, "sub"); !exists || sub != "test-user" {
        t.Errorf("Expected sub='test-user', got '%s'", sub)
    }
}

func TestTokenExpiration(t *testing.T) {
    secret := []byte("test-secret")
    claims := jwt.MapClaims{
        "sub": "test-user",
        "exp": time.Now().Add(-time.Hour).Unix(), // 已过期
    }
    
    tokenString, err := jwt.GenerateHS256(secret, claims)
    if err != nil {
        t.Fatalf("Failed to generate token: %v", err)
    }
    
    // 尝试解析过期令牌
    _, err = jwt.ParseHS256(tokenString, secret)
    if err != jwt.ErrTokenExpired {
        t.Errorf("Expected ErrTokenExpired, got %v", err)
    }
}
```

## 🔍 错误处理

```go
func HandleJWTError(err error) {
    switch err {
    case jwt.ErrInvalidToken:
        fmt.Println("令牌格式无效")
    case jwt.ErrInvalidSignature:
        fmt.Println("签名验证失败")
    case jwt.ErrTokenExpired:
        fmt.Println("令牌已过期")
    case jwt.ErrTokenNotYetValid:
        fmt.Println("令牌还未生效")
    case jwt.ErrInvalidAudience:
        fmt.Println("无效的受众")
    case jwt.ErrInvalidIssuer:
        fmt.Println("无效的签发者")
    case jwt.ErrInvalidSubject:
        fmt.Println("无效的主题")
    case jwt.ErrInvalidKeyType:
        fmt.Println("无效的密钥类型")
    case jwt.ErrKeyMustBePEM:
        fmt.Println("密钥必须是 PEM 格式")
    default:
        fmt.Printf("未知错误: %v\n", err)
    }
}
```

## 📚 最佳实践

1. **安全的密钥管理**
   - 使用足够长的随机密钥
   - 定期轮换密钥
   - 安全存储私钥

2. **合理的过期时间**
   - 访问令牌：15分钟 - 1小时
   - 刷新令牌：7天 - 30天
   - 敏感操作：5分钟 - 15分钟

3. **声明最小化**
   - 只包含必要的信息
   - 避免敏感数据（如密码）
   - 使用简短的键名

4. **验证所有声明**
   - 检查过期时间
   - 验证签发者和受众
   - 验证自定义声明

5. **错误处理**
   - 区分不同类型的错误
   - 提供有意义的错误信息
   - 记录安全相关的错误

---

**让 JWT 更安全，让认证更简单！** 🔐 