package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"
)

func TestHMACTokenGeneration(t *testing.T) {
	secret := []byte("test-secret-key")
	claims := MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	// 生成令牌
	tokenString, err := GenerateHS256(secret, claims)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if tokenString == "" {
		t.Fatal("Token string should not be empty")
	}

	// 验证令牌格式
	if !IsTokenValid(tokenString) {
		t.Error("Generated token should have valid format")
	}
}

func TestHMACTokenParsing(t *testing.T) {
	secret := []byte("test-secret-key")
	claims := MapClaims{
		"sub":  "test-user",
		"name": "Test User",
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	// 生成令牌
	tokenString, err := GenerateHS256(secret, claims)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 解析令牌
	token, err := ParseHS256(tokenString, secret)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !token.Valid {
		t.Error("Token should be valid")
	}

	// 验证声明
	parsedClaims, ok := ExtractClaims(token)
	if !ok {
		t.Fatal("Failed to extract claims")
	}

	if sub, exists := GetClaimString(parsedClaims, "sub"); !exists || sub != "test-user" {
		t.Errorf("Expected sub='test-user', got '%s'", sub)
	}

	if name, exists := GetClaimString(parsedClaims, "name"); !exists || name != "Test User" {
		t.Errorf("Expected name='Test User', got '%s'", name)
	}
}

func TestTokenExpiration(t *testing.T) {
	secret := []byte("test-secret-key")
	claims := MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(-time.Hour).Unix(), // 已过期
	}

	tokenString, err := GenerateHS256(secret, claims)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 尝试解析过期令牌
	_, err = ParseHS256(tokenString, secret)
	if err != ErrTokenExpired {
		t.Errorf("Expected ErrTokenExpired, got %v", err)
	}
}

func TestTokenNotYetValid(t *testing.T) {
	secret := []byte("test-secret-key")
	claims := MapClaims{
		"sub": "test-user",
		"nbf": time.Now().Add(time.Hour).Unix(), // 1小时后生效
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	}

	tokenString, err := GenerateHS256(secret, claims)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 尝试解析还未生效的令牌
	_, err = ParseHS256(tokenString, secret)
	if err != ErrTokenNotYetValid {
		t.Errorf("Expected ErrTokenNotYetValid, got %v", err)
	}
}

func TestInvalidSignature(t *testing.T) {
	secret1 := []byte("secret-key-1")
	secret2 := []byte("secret-key-2")
	claims := MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	// 用 secret1 生成令牌
	tokenString, err := GenerateHS256(secret1, claims)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 用 secret2 验证令牌
	_, err = ParseHS256(tokenString, secret2)
	if err != ErrInvalidSignature {
		t.Errorf("Expected ErrInvalidSignature, got %v", err)
	}
}

func TestJWTBuilder(t *testing.T) {
	secret := []byte("test-secret-key")

	tokenString, err := NewBuilder(SigningMethodHS256, secret).
		SetIssuer("test-issuer").
		SetSubject("test-user").
		SetAudience("test-audience").
		SetExpirationFromNow(time.Hour).
		SetIssuedNow().
		SetClaim("role", "admin").
		Build()

	if err != nil {
		t.Fatalf("Failed to build token: %v", err)
	}

	// 解析并验证
	token, err := ParseHS256(tokenString, secret)
	if err != nil {
		t.Fatalf("Failed to parse built token: %v", err)
	}

	claims, ok := ExtractClaims(token)
	if !ok {
		t.Fatal("Failed to extract claims")
	}

	if iss, exists := GetClaimString(claims, "iss"); !exists || iss != "test-issuer" {
		t.Errorf("Expected iss='test-issuer', got '%s'", iss)
	}

	if role, exists := GetClaimString(claims, "role"); !exists || role != "admin" {
		t.Errorf("Expected role='admin', got '%s'", role)
	}
}

func TestRSATokenGeneration(t *testing.T) {
	// 生成 RSA 密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	claims := MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	// 生成令牌
	tokenString, err := GenerateRS256(privateKey, claims)
	if err != nil {
		t.Fatalf("Failed to generate RSA token: %v", err)
	}

	// 验证令牌
	token, err := ParseRS256(tokenString, &privateKey.PublicKey)
	if err != nil {
		t.Fatalf("Failed to parse RSA token: %v", err)
	}

	if !token.Valid {
		t.Error("RSA token should be valid")
	}
}

func TestRSAPEMKeys(t *testing.T) {
	// 生成密钥对
	privateKey, err := GenerateRSAKeyPair(2048)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// 转换为 PEM
	privatePEM := PrivateKeyToPEM(privateKey)
	publicPEM, err := PublicKeyToPEM(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("Failed to convert public key to PEM: %v", err)
	}

	claims := MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	// 使用 PEM 密钥生成令牌
	tokenString, err := GenerateRS256(privatePEM, claims)
	if err != nil {
		t.Fatalf("Failed to generate token with PEM key: %v", err)
	}

	// 使用 PEM 公钥验证令牌
	token, err := ParseRS256(tokenString, publicPEM)
	if err != nil {
		t.Fatalf("Failed to parse token with PEM key: %v", err)
	}

	if !token.Valid {
		t.Error("Token should be valid")
	}
}

func TestStandardClaims(t *testing.T) {
	secret := []byte("test-secret-key")
	now := time.Now()

	claims := &StandardClaims{
		Issuer:    "test-issuer",
		Subject:   "test-user",
		Audience:  "test-audience",
		ExpiresAt: now.Add(time.Hour).Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ID:        "test-token-id",
	}

	tokenString, err := GenerateHS256(secret, claims)
	if err != nil {
		t.Fatalf("Failed to generate token with StandardClaims: %v", err)
	}

	// 解析令牌
	jwt := New(SigningMethodHS256, secret)
	parsedClaims := &StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, parsedClaims)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !token.Valid {
		t.Error("Token should be valid")
	}

	if parsedClaims.Subject != "test-user" {
		t.Errorf("Expected subject='test-user', got '%s'", parsedClaims.Subject)
	}
}

func TestDecodeHeaderAndClaims(t *testing.T) {
	secret := []byte("test-secret-key")
	claims := MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	tokenString, err := GenerateHS256(secret, claims)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 解码头部
	header, err := DecodeHeader(tokenString)
	if err != nil {
		t.Fatalf("Failed to decode header: %v", err)
	}

	if header.Algorithm != "HS256" {
		t.Errorf("Expected alg='HS256', got '%s'", header.Algorithm)
	}

	if header.Type != "JWT" {
		t.Errorf("Expected typ='JWT', got '%s'", header.Type)
	}

	// 解码声明
	decodedClaims, err := DecodeClaims(tokenString)
	if err != nil {
		t.Fatalf("Failed to decode claims: %v", err)
	}

	if sub, exists := GetClaimString(decodedClaims, "sub"); !exists || sub != "test-user" {
		t.Errorf("Expected sub='test-user', got '%s'", sub)
	}
}

func TestValidateStandardClaims(t *testing.T) {
	now := time.Now()
	claims := MapClaims{
		"iss": "test-issuer",
		"sub": "test-user",
		"aud": "test-audience",
		"exp": now.Add(time.Hour).Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
	}

	// 验证有效声明
	err := ValidateStandardClaims(claims, "test-audience", "test-issuer", "test-user")
	if err != nil {
		t.Errorf("Valid claims should pass validation: %v", err)
	}

	// 验证无效受众
	err = ValidateStandardClaims(claims, "wrong-audience", "test-issuer", "test-user")
	if err != ErrInvalidAudience {
		t.Errorf("Expected ErrInvalidAudience, got %v", err)
	}

	// 验证无效签发者
	err = ValidateStandardClaims(claims, "test-audience", "wrong-issuer", "test-user")
	if err != ErrInvalidIssuer {
		t.Errorf("Expected ErrInvalidIssuer, got %v", err)
	}

	// 验证无效主题
	err = ValidateStandardClaims(claims, "test-audience", "test-issuer", "wrong-user")
	if err != ErrInvalidSubject {
		t.Errorf("Expected ErrInvalidSubject, got %v", err)
	}
}

func TestHMACSecretGeneration(t *testing.T) {
	secret, err := GenerateHMACSecret(32)
	if err != nil {
		t.Fatalf("Failed to generate HMAC secret: %v", err)
	}

	if len(secret) != 32 {
		t.Errorf("Expected secret length 32, got %d", len(secret))
	}

	// 生成的密钥应该是随机的
	secret2, err := GenerateHMACSecret(32)
	if err != nil {
		t.Fatalf("Failed to generate second HMAC secret: %v", err)
	}

	// 两个密钥不应该相同（概率极低）
	if string(secret) == string(secret2) {
		t.Error("Generated secrets should be different")
	}
}

func TestTimeUtilities(t *testing.T) {
	now := time.Now()
	unix := TimeToUnix(now)

	if unix != now.Unix() {
		t.Errorf("TimeToUnix conversion failed")
	}

	backToTime := UnixToTime(unix)
	if backToTime.Unix() != now.Unix() {
		t.Errorf("UnixToTime conversion failed")
	}

	// 测试过期检查
	expiredTime := time.Now().Add(-time.Hour).Unix()
	if !IsExpired(expiredTime) {
		t.Error("Should detect expired timestamp")
	}

	futureTime := time.Now().Add(time.Hour).Unix()
	if IsExpired(futureTime) {
		t.Error("Should not detect future timestamp as expired")
	}

	// 测试未生效检查
	if !IsNotYetValid(futureTime) {
		t.Error("Should detect future timestamp as not yet valid")
	}

	pastTime := time.Now().Add(-time.Hour).Unix()
	if IsNotYetValid(pastTime) {
		t.Error("Should not detect past timestamp as not yet valid")
	}
}

// 测试新的通用 API
func TestGenericAPI(t *testing.T) {
	secret := []byte("test-secret")
	claims := MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	// 测试通用生成函数
	tokenString, err := Generate(SigningMethodHS256, secret, claims)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if tokenString == "" {
		t.Fatal("Token string should not be empty")
	}

	// 测试通用解析函数
	token, err := Parse(SigningMethodHS256, tokenString, secret)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if !token.Valid {
		t.Error("Token should be valid")
	}

	// 验证声明
	parsedClaims, ok := ExtractClaims(token)
	if !ok {
		t.Fatal("Failed to extract claims")
	}

	if sub, exists := GetClaimString(parsedClaims, "sub"); !exists || sub != "test-user" {
		t.Errorf("Expected sub='test-user', got '%s'", sub)
	}
}

func TestGenericAPIWithDifferentAlgorithms(t *testing.T) {
	secret := []byte("test-secret")
	claims := MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	// 测试不同的 HMAC 算法
	algorithms := []SigningMethod{
		SigningMethodHS256,
		SigningMethodHS384,
		SigningMethodHS512,
	}

	for _, method := range algorithms {
		t.Run(method.Alg(), func(t *testing.T) {
			// 生成令牌
			tokenString, err := Generate(method, secret, claims)
			if err != nil {
				t.Fatalf("Generate with %s failed: %v", method.Alg(), err)
			}

			// 解析令牌
			token, err := Parse(method, tokenString, secret)
			if err != nil {
				t.Fatalf("Parse with %s failed: %v", method.Alg(), err)
			}

			if !token.Valid {
				t.Errorf("Token should be valid for %s", method.Alg())
			}

			// 验证算法
			if token.Header.Algorithm != method.Alg() {
				t.Errorf("Expected algorithm %s, got %s", method.Alg(), token.Header.Algorithm)
			}
		})
	}
}

func TestGenericAPIWithCustomClaims(t *testing.T) {
	secret := []byte("test-secret")
	now := time.Now()

	// 使用标准声明
	standardClaims := &StandardClaims{
		Issuer:    "test-issuer",
		Subject:   "test-user",
		Audience:  "test-audience",
		ExpiresAt: now.Add(time.Hour).Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ID:        "test-token-id",
	}

	// 生成令牌
	tokenString, err := Generate(SigningMethodHS256, secret, standardClaims)
	if err != nil {
		t.Fatalf("Generate with StandardClaims failed: %v", err)
	}

	// 解析令牌
	parsedClaims := &StandardClaims{}
	token, err := ParseWithClaims(SigningMethodHS256, tokenString, secret, parsedClaims)
	if err != nil {
		t.Fatalf("ParseWithClaims failed: %v", err)
	}

	if !token.Valid {
		t.Error("Token should be valid")
	}

	// 验证声明
	if parsedClaims.Subject != "test-user" {
		t.Errorf("Expected subject='test-user', got '%s'", parsedClaims.Subject)
	}

	if parsedClaims.Issuer != "test-issuer" {
		t.Errorf("Expected issuer='test-issuer', got '%s'", parsedClaims.Issuer)
	}

	if parsedClaims.Audience != "test-audience" {
		t.Errorf("Expected audience='test-audience', got '%s'", parsedClaims.Audience)
	}
}

func TestGenericAPIRSA(t *testing.T) {
	// 生成 RSA 密钥对
	privateKey, err := GenerateRSAKeyPair(2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	claims := MapClaims{
		"sub": "rsa-test-user",
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	// 测试 RSA 算法
	rsaAlgorithms := []SigningMethod{
		SigningMethodRS256,
		SigningMethodRS384,
		SigningMethodRS512,
	}

	for _, method := range rsaAlgorithms {
		t.Run(method.Alg(), func(t *testing.T) {
			// 生成令牌
			tokenString, err := Generate(method, privateKey, claims)
			if err != nil {
				t.Fatalf("Generate with %s failed: %v", method.Alg(), err)
			}

			// 解析令牌
			token, err := Parse(method, tokenString, &privateKey.PublicKey)
			if err != nil {
				t.Fatalf("Parse with %s failed: %v", method.Alg(), err)
			}

			if !token.Valid {
				t.Errorf("Token should be valid for %s", method.Alg())
			}
		})
	}
}
