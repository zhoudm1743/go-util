package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zhoudm1743/go-util/jwt"
)

func main() {
	fmt.Println("🔐 JWT 包使用示例")
	fmt.Println("==================")

	// 1. 通用 API 示例（推荐）
	fmt.Println("\n1. 通用 API 示例（推荐）")
	genericAPIExample()

	// 2. HMAC 签名示例（传统方式）
	fmt.Println("\n2. HMAC 签名示例（传统方式）")
	hmacExample()

	// 3. RSA 签名示例
	fmt.Println("\n3. RSA 签名示例")
	rsaExample()

	// 4. 构建器模式示例
	fmt.Println("\n4. 构建器模式示例")
	builderExample()

	// 5. 令牌验证示例
	fmt.Println("\n5. 令牌验证示例")
	validationExample()

	// 6. 时间工具示例
	fmt.Println("\n6. 时间工具示例")
	timeUtilsExample()
}

func genericAPIExample() {
	secret := []byte("my-secret-key")

	// 创建声明
	claims := jwt.MapClaims{
		"sub":  "user123",
		"name": "李四",
		"role": "manager",
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"iat":  time.Now().Unix(),
	}

	// 🔥 使用通用 API 生成 JWT（推荐方式）
	// 可以轻松切换算法，只需要改变 SigningMethod
	algorithms := []struct {
		name   string
		method jwt.SigningMethod
		key    interface{}
	}{
		{"HS256", jwt.SigningMethodHS256, secret},
		{"HS384", jwt.SigningMethodHS384, secret},
		{"HS512", jwt.SigningMethodHS512, secret},
	}

	for _, alg := range algorithms {
		fmt.Printf("\n使用 %s 算法:\n", alg.name)

		// 生成 JWT
		tokenString, err := jwt.Generate(alg.method, alg.key, claims)
		if err != nil {
			fmt.Printf("生成令牌失败: %v\n", err)
			continue
		}

		fmt.Printf("生成的 JWT: %s...\n", tokenString[:50])

		// 解析 JWT
		token, err := jwt.Parse(alg.method, tokenString, alg.key)
		if err != nil {
			fmt.Printf("解析令牌失败: %v\n", err)
			continue
		}

		// 提取声明
		if parsedClaims, ok := jwt.ExtractClaims(token); ok {
			if name, exists := jwt.GetClaimString(parsedClaims, "name"); exists {
				fmt.Printf("用户名: %s\n", name)
			}
			if role, exists := jwt.GetClaimString(parsedClaims, "role"); exists {
				fmt.Printf("角色: %s\n", role)
			}
		}
	}

	// 演示使用自定义声明类型
	fmt.Printf("\n使用自定义声明类型:\n")
	standardClaims := &jwt.StandardClaims{
		Subject:   "user456",
		Issuer:    "go-util-jwt",
		Audience:  "web-app",
		ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	tokenString, err := jwt.Generate(jwt.SigningMethodHS256, secret, standardClaims)
	if err == nil {
		fmt.Printf("标准声明 JWT: %s...\n", tokenString[:50])

		// 使用相同的声明类型解析
		parsedClaims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(jwt.SigningMethodHS256, tokenString, secret, parsedClaims)
		if err == nil && token.Valid {
			fmt.Printf("解析成功 - 主题: %s, 签发者: %s\n", parsedClaims.Subject, parsedClaims.Issuer)
		}
	}
}

func hmacExample() {
	secret := []byte("my-secret-key")

	// 创建声明
	claims := jwt.MapClaims{
		"sub":  "user123",
		"name": "张三",
		"role": "admin",
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"iat":  time.Now().Unix(),
	}

	// 生成 JWT
	tokenString, err := jwt.GenerateHS256(secret, claims)
	if err != nil {
		log.Fatalf("生成令牌失败: %v", err)
	}

	fmt.Printf("生成的 JWT: %s\n", tokenString)

	// 解析 JWT
	token, err := jwt.ParseHS256(tokenString, secret)
	if err != nil {
		log.Fatalf("解析令牌失败: %v", err)
	}

	// 提取声明
	if parsedClaims, ok := jwt.ExtractClaims(token); ok {
		if name, exists := jwt.GetClaimString(parsedClaims, "name"); exists {
			fmt.Printf("用户名: %s\n", name)
		}
		if role, exists := jwt.GetClaimString(parsedClaims, "role"); exists {
			fmt.Printf("角色: %s\n", role)
		}
	}
}

func rsaExample() {
	// 生成 RSA 密钥对
	privateKey, err := jwt.GenerateRSAKeyPair(2048)
	if err != nil {
		log.Fatalf("生成 RSA 密钥失败: %v", err)
	}

	// 转换为 PEM 格式
	privatePEM := jwt.PrivateKeyToPEM(privateKey)
	publicPEM, err := jwt.PublicKeyToPEM(&privateKey.PublicKey)
	if err != nil {
		log.Fatalf("转换公钥失败: %v", err)
	}

	// 创建声明
	claims := jwt.MapClaims{
		"sub": "user456",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	// 生成 JWT
	tokenString, err := jwt.GenerateRS256(privatePEM, claims)
	if err != nil {
		log.Fatalf("生成 RSA 令牌失败: %v", err)
	}

	fmt.Printf("RSA JWT: %s...\n", tokenString[:50])

	// 验证 JWT
	token, err := jwt.ParseRS256(tokenString, publicPEM)
	if err != nil {
		log.Fatalf("验证 RSA 令牌失败: %v", err)
	}

	fmt.Printf("RSA 令牌验证成功: %t\n", token.Valid)
}

func builderExample() {
	secret := []byte("builder-secret")

	// 使用构建器创建复杂的 JWT
	tokenString, err := jwt.NewBuilder(jwt.SigningMethodHS256, secret).
		SetIssuer("go-util-jwt").
		SetSubject("user789").
		SetAudience("web-app").
		SetExpirationFromNow(time.Hour*2).
		SetIssuedNow().
		SetJWTID("unique-token-id").
		SetClaim("role", "user").
		SetClaim("permissions", []string{"read", "write"}).
		SetClaim("department", "engineering").
		Build()

	if err != nil {
		log.Fatalf("构建令牌失败: %v", err)
	}

	fmt.Printf("构建器生成的 JWT: %s...\n", tokenString[:50])

	// 解析并显示声明
	token, err := jwt.ParseHS256(tokenString, secret)
	if err != nil {
		log.Fatalf("解析令牌失败: %v", err)
	}

	if claims, ok := jwt.ExtractClaims(token); ok {
		fmt.Printf("签发者: %s\n", claims["iss"])
		fmt.Printf("主题: %s\n", claims["sub"])
		fmt.Printf("受众: %s\n", claims["aud"])
		fmt.Printf("角色: %s\n", claims["role"])

		if permissions, ok := claims["permissions"].([]interface{}); ok {
			fmt.Printf("权限: %v\n", permissions)
		}
	}
}

func validationExample() {
	secret := []byte("validation-secret")

	// 创建一个即将过期的令牌
	claims := jwt.MapClaims{
		"sub": "user999",
		"exp": time.Now().Add(time.Second * 5).Unix(), // 5秒后过期
		"iat": time.Now().Unix(),
	}

	tokenString, err := jwt.GenerateHS256(secret, claims)
	if err != nil {
		log.Fatalf("生成令牌失败: %v", err)
	}

	// 立即验证（应该成功）
	token, err := jwt.ParseHS256(tokenString, secret)
	if err != nil {
		fmt.Printf("验证失败: %v\n", err)
	} else {
		fmt.Printf("令牌验证成功: %t\n", token.Valid)
	}

	// 等待令牌过期
	fmt.Println("等待令牌过期...")
	time.Sleep(time.Second * 6)

	// 再次验证（应该失败）
	_, err = jwt.ParseHS256(tokenString, secret)
	if err != nil {
		if err == jwt.ErrTokenExpired {
			fmt.Println("令牌已过期（预期结果）")
		} else {
			fmt.Printf("验证失败: %v\n", err)
		}
	}

	// 演示错误的密钥
	wrongSecret := []byte("wrong-secret")
	_, err = jwt.ParseHS256(tokenString, wrongSecret)
	if err != nil {
		if err == jwt.ErrInvalidSignature {
			fmt.Println("签名验证失败（预期结果）")
		} else {
			fmt.Printf("验证失败: %v\n", err)
		}
	}
}

func timeUtilsExample() {
	now := time.Now()

	// 时间转换
	unix := jwt.TimeToUnix(now)
	backToTime := jwt.UnixToTime(unix)
	fmt.Printf("时间转换: %s -> %d -> %s\n",
		now.Format("15:04:05"), unix, backToTime.Format("15:04:05"))

	// 过期检查
	expiredTime := time.Now().Add(-time.Hour).Unix()
	futureTime := time.Now().Add(time.Hour).Unix()

	fmt.Printf("过期检查: 过去时间已过期 = %t, 未来时间已过期 = %t\n",
		jwt.IsExpired(expiredTime), jwt.IsExpired(futureTime))

	// 生效检查
	fmt.Printf("生效检查: 过去时间未生效 = %t, 未来时间未生效 = %t\n",
		jwt.IsNotYetValid(expiredTime), jwt.IsNotYetValid(futureTime))

	// 令牌信息提取（不验证签名）
	secret := []byte("decode-secret")
	claims := jwt.MapClaims{
		"sub": "decode-user",
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	tokenString, _ := jwt.GenerateHS256(secret, claims)

	// 仅解码头部
	header, err := jwt.DecodeHeader(tokenString)
	if err == nil {
		fmt.Printf("令牌算法: %s, 类型: %s\n", header.Algorithm, header.Type)
	}

	// 仅解码声明
	decodedClaims, err := jwt.DecodeClaims(tokenString)
	if err == nil {
		if sub, exists := jwt.GetClaimString(decodedClaims, "sub"); exists {
			fmt.Printf("令牌主题: %s\n", sub)
		}
	}
}
