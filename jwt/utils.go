package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"strings"
	"time"
)

// JWTBuilder JWT 构建器，提供链式调用
type JWTBuilder struct {
	method SigningMethod
	key    interface{}
	claims MapClaims
}

// NewBuilder 创建 JWT 构建器
func NewBuilder(method SigningMethod, key interface{}) *JWTBuilder {
	return &JWTBuilder{
		method: method,
		key:    key,
		claims: make(MapClaims),
	}
}

// SetIssuer 设置签发者
func (b *JWTBuilder) SetIssuer(issuer string) *JWTBuilder {
	b.claims["iss"] = issuer
	return b
}

// SetSubject 设置主题
func (b *JWTBuilder) SetSubject(subject string) *JWTBuilder {
	b.claims["sub"] = subject
	return b
}

// SetAudience 设置受众
func (b *JWTBuilder) SetAudience(audience string) *JWTBuilder {
	b.claims["aud"] = audience
	return b
}

// SetExpiration 设置过期时间
func (b *JWTBuilder) SetExpiration(exp time.Time) *JWTBuilder {
	b.claims["exp"] = exp.Unix()
	return b
}

// SetExpirationFromNow 设置从现在开始的过期时间
func (b *JWTBuilder) SetExpirationFromNow(duration time.Duration) *JWTBuilder {
	b.claims["exp"] = time.Now().Add(duration).Unix()
	return b
}

// SetNotBefore 设置生效时间
func (b *JWTBuilder) SetNotBefore(nbf time.Time) *JWTBuilder {
	b.claims["nbf"] = nbf.Unix()
	return b
}

// SetIssuedAt 设置签发时间
func (b *JWTBuilder) SetIssuedAt(iat time.Time) *JWTBuilder {
	b.claims["iat"] = iat.Unix()
	return b
}

// SetIssuedNow 设置签发时间为当前时间
func (b *JWTBuilder) SetIssuedNow() *JWTBuilder {
	b.claims["iat"] = time.Now().Unix()
	return b
}

// SetJWTID 设置 JWT ID
func (b *JWTBuilder) SetJWTID(jti string) *JWTBuilder {
	b.claims["jti"] = jti
	return b
}

// SetClaim 设置自定义声明
func (b *JWTBuilder) SetClaim(key string, value interface{}) *JWTBuilder {
	b.claims[key] = value
	return b
}

// SetClaims 设置多个声明
func (b *JWTBuilder) SetClaims(claims map[string]interface{}) *JWTBuilder {
	for k, v := range claims {
		b.claims[k] = v
	}
	return b
}

// Build 构建 JWT 令牌字符串
func (b *JWTBuilder) Build() (string, error) {
	jwt := New(b.method, b.key)
	return jwt.Generate(b.claims)
}

// BuildToken 构建 JWT 令牌对象
func (b *JWTBuilder) BuildToken() (*Token, error) {
	token := NewWithClaims(b.method, b.claims)
	return token, nil
}

// 通用便捷函数

// Generate 使用指定算法生成 JWT
func Generate(method SigningMethod, key interface{}, claims Claims) (string, error) {
	jwt := New(method, key)
	return jwt.Generate(claims)
}

// Parse 使用指定算法解析 JWT
func Parse(method SigningMethod, tokenString string, key interface{}) (*Token, error) {
	jwt := New(method, key)
	claims := make(MapClaims)
	return jwt.ParseWithClaims(tokenString, claims)
}

// ParseWithClaims 使用指定算法和声明类型解析 JWT
func ParseWithClaims(method SigningMethod, tokenString string, key interface{}, claims Claims) (*Token, error) {
	jwt := New(method, key)
	return jwt.ParseWithClaims(tokenString, claims)
}

// HMAC 算法便捷函数

// GenerateHS256 使用 HS256 算法生成 JWT
func GenerateHS256(secret []byte, claims Claims) (string, error) {
	return Generate(SigningMethodHS256, secret, claims)
}

// GenerateHS384 使用 HS384 算法生成 JWT
func GenerateHS384(secret []byte, claims Claims) (string, error) {
	return Generate(SigningMethodHS384, secret, claims)
}

// GenerateHS512 使用 HS512 算法生成 JWT
func GenerateHS512(secret []byte, claims Claims) (string, error) {
	return Generate(SigningMethodHS512, secret, claims)
}

// ParseHS256 使用 HS256 算法解析 JWT
func ParseHS256(tokenString string, secret []byte) (*Token, error) {
	return Parse(SigningMethodHS256, tokenString, secret)
}

// ParseHS384 使用 HS384 算法解析 JWT
func ParseHS384(tokenString string, secret []byte) (*Token, error) {
	return Parse(SigningMethodHS384, tokenString, secret)
}

// ParseHS512 使用 HS512 算法解析 JWT
func ParseHS512(tokenString string, secret []byte) (*Token, error) {
	return Parse(SigningMethodHS512, tokenString, secret)
}

// RSA 算法便捷函数

// GenerateRS256 使用 RS256 算法生成 JWT
func GenerateRS256(privateKey interface{}, claims Claims) (string, error) {
	return Generate(SigningMethodRS256, privateKey, claims)
}

// GenerateRS384 使用 RS384 算法生成 JWT
func GenerateRS384(privateKey interface{}, claims Claims) (string, error) {
	return Generate(SigningMethodRS384, privateKey, claims)
}

// GenerateRS512 使用 RS512 算法生成 JWT
func GenerateRS512(privateKey interface{}, claims Claims) (string, error) {
	return Generate(SigningMethodRS512, privateKey, claims)
}

// ParseRS256 使用 RS256 算法解析 JWT
func ParseRS256(tokenString string, publicKey interface{}) (*Token, error) {
	return Parse(SigningMethodRS256, tokenString, publicKey)
}

// ParseRS384 使用 RS384 算法解析 JWT
func ParseRS384(tokenString string, publicKey interface{}) (*Token, error) {
	return Parse(SigningMethodRS384, tokenString, publicKey)
}

// ParseRS512 使用 RS512 算法解析 JWT
func ParseRS512(tokenString string, publicKey interface{}) (*Token, error) {
	return Parse(SigningMethodRS512, tokenString, publicKey)
}

// 密钥生成工具

// GenerateHMACSecret 生成 HMAC 密钥
func GenerateHMACSecret(length int) ([]byte, error) {
	secret := make([]byte, length)
	_, err := rand.Read(secret)
	return secret, err
}

// GenerateRSAKeyPair 生成 RSA 密钥对
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

// PrivateKeyToPEM 将 RSA 私钥转换为 PEM 格式
func PrivateKeyToPEM(key *rsa.PrivateKey) []byte {
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	})
	return keyPEM
}

// PublicKeyToPEM 将 RSA 公钥转换为 PEM 格式
func PublicKeyToPEM(key *rsa.PublicKey) ([]byte, error) {
	keyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	})
	return keyPEM, nil
}

// 时间工具

// TimeToUnix 将时间转换为 Unix 时间戳
func TimeToUnix(t time.Time) int64 {
	return t.Unix()
}

// UnixToTime 将 Unix 时间戳转换为时间
func UnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// IsExpired 检查 Unix 时间戳是否已过期
func IsExpired(exp int64) bool {
	return time.Now().Unix() > exp
}

// IsNotYetValid 检查 Unix 时间戳是否还未生效
func IsNotYetValid(nbf int64) bool {
	return time.Now().Unix() < nbf
}

// 声明工具

// NewStandardClaims 创建标准声明
func NewStandardClaims() *StandardClaims {
	return &StandardClaims{}
}

// NewMapClaims 创建映射声明
func NewMapClaims() MapClaims {
	return make(MapClaims)
}

// GetClaimString 从 MapClaims 中获取字符串声明
func GetClaimString(claims MapClaims, key string) (string, bool) {
	if value, exists := claims[key]; exists {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}

// GetClaimInt64 从 MapClaims 中获取 int64 声明
func GetClaimInt64(claims MapClaims, key string) (int64, bool) {
	if value, exists := claims[key]; exists {
		switch v := value.(type) {
		case int64:
			return v, true
		case float64:
			return int64(v), true
		case int:
			return int64(v), true
		}
	}
	return 0, false
}

// GetClaimFloat64 从 MapClaims 中获取 float64 声明
func GetClaimFloat64(claims MapClaims, key string) (float64, bool) {
	if value, exists := claims[key]; exists {
		switch v := value.(type) {
		case float64:
			return v, true
		case int64:
			return float64(v), true
		case int:
			return float64(v), true
		}
	}
	return 0, false
}

// GetClaimBool 从 MapClaims 中获取布尔声明
func GetClaimBool(claims MapClaims, key string) (bool, bool) {
	if value, exists := claims[key]; exists {
		if b, ok := value.(bool); ok {
			return b, true
		}
	}
	return false, false
}

// ExtractClaims 从令牌中提取声明为 MapClaims
func ExtractClaims(token *Token) (MapClaims, bool) {
	if claims, ok := token.Claims.(MapClaims); ok {
		return claims, true
	}
	return nil, false
}

// 验证工具

// ValidateStandardClaims 验证标准声明
func ValidateStandardClaims(claims MapClaims, audience, issuer, subject string) error {
	now := time.Now().Unix()

	// 验证过期时间
	if exp, exists := GetClaimInt64(claims, "exp"); exists && now > exp {
		return ErrTokenExpired
	}

	// 验证生效时间
	if nbf, exists := GetClaimInt64(claims, "nbf"); exists && now < nbf {
		return ErrTokenNotYetValid
	}

	// 验证受众
	if audience != "" {
		if aud, exists := GetClaimString(claims, "aud"); !exists || aud != audience {
			return ErrInvalidAudience
		}
	}

	// 验证签发者
	if issuer != "" {
		if iss, exists := GetClaimString(claims, "iss"); !exists || iss != issuer {
			return ErrInvalidIssuer
		}
	}

	// 验证主题
	if subject != "" {
		if sub, exists := GetClaimString(claims, "sub"); !exists || sub != subject {
			return ErrInvalidSubject
		}
	}

	return nil
}

// IsTokenValid 检查令牌是否有效（不验证签名）
func IsTokenValid(tokenString string) bool {
	parts := strings.Split(tokenString, ".")
	return len(parts) == 3
}

// DecodeHeader 仅解码 JWT 头部（不验证签名）
func DecodeHeader(tokenString string) (*Header, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	headerBytes, err := base64URLDecode(parts[0])
	if err != nil {
		return nil, err
	}

	header := &Header{}
	if err := json.Unmarshal(headerBytes, header); err != nil {
		return nil, err
	}

	return header, nil
}

// DecodeClaims 仅解码 JWT 声明（不验证签名）
func DecodeClaims(tokenString string) (MapClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	claimsBytes, err := base64URLDecode(parts[1])
	if err != nil {
		return nil, err
	}

	claims := make(MapClaims)
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, err
	}

	return claims, nil
}
