package jwt

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"strings"
	"time"
)

// JWT 错误定义
var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenNotYetValid = errors.New("token not yet valid")
	ErrInvalidAudience  = errors.New("invalid audience")
	ErrInvalidIssuer    = errors.New("invalid issuer")
	ErrInvalidSubject   = errors.New("invalid subject")
	ErrInvalidKeyType   = errors.New("invalid key type")
	ErrKeyMustBePEM     = errors.New("key must be PEM encoded")
)

// SigningMethod 签名方法接口
type SigningMethod interface {
	Sign(signingString string, key interface{}) (string, error)
	Verify(signingString, signature string, key interface{}) error
	Alg() string
}

// Header JWT 头部
type Header struct {
	Type      string `json:"typ"`
	Algorithm string `json:"alg"`
	KeyID     string `json:"kid,omitempty"`
}

// StandardClaims 标准声明
type StandardClaims struct {
	Audience  string `json:"aud,omitempty"` // 受众
	ExpiresAt int64  `json:"exp,omitempty"` // 过期时间
	ID        string `json:"jti,omitempty"` // JWT ID
	IssuedAt  int64  `json:"iat,omitempty"` // 签发时间
	Issuer    string `json:"iss,omitempty"` // 签发者
	NotBefore int64  `json:"nbf,omitempty"` // 生效时间
	Subject   string `json:"sub,omitempty"` // 主题
}

// Valid 验证标准声明
func (c StandardClaims) Valid() error {
	now := time.Now().Unix()

	// 检查过期时间
	if c.ExpiresAt != 0 && now > c.ExpiresAt {
		return ErrTokenExpired
	}

	// 检查生效时间
	if c.NotBefore != 0 && now < c.NotBefore {
		return ErrTokenNotYetValid
	}

	return nil
}

// Claims 声明接口
type Claims interface {
	Valid() error
}

// MapClaims 映射类型的声明
type MapClaims map[string]interface{}

// Valid 验证映射声明
func (m MapClaims) Valid() error {
	now := time.Now().Unix()

	// 检查过期时间
	if exp, ok := m["exp"]; ok {
		switch exp := exp.(type) {
		case float64:
			if now > int64(exp) {
				return ErrTokenExpired
			}
		case int64:
			if now > exp {
				return ErrTokenExpired
			}
		}
	}

	// 检查生效时间
	if nbf, ok := m["nbf"]; ok {
		switch nbf := nbf.(type) {
		case float64:
			if now < int64(nbf) {
				return ErrTokenNotYetValid
			}
		case int64:
			if now < nbf {
				return ErrTokenNotYetValid
			}
		}
	}

	return nil
}

// Token JWT 令牌
type Token struct {
	Raw       string        // 原始令牌字符串
	Method    SigningMethod // 签名方法
	Header    *Header       // 头部
	Claims    Claims        // 声明
	Signature string        // 签名
	Valid     bool          // 是否有效
}

// SignedString 生成签名后的 JWT 字符串
func (t *Token) SignedString(key interface{}) (string, error) {
	// 编码头部
	headerBytes, err := json.Marshal(t.Header)
	if err != nil {
		return "", err
	}
	header := base64URLEncode(headerBytes)

	// 编码声明
	claimsBytes, err := json.Marshal(t.Claims)
	if err != nil {
		return "", err
	}
	claims := base64URLEncode(claimsBytes)

	// 生成签名字符串
	signingString := header + "." + claims

	// 签名
	signature, err := t.Method.Sign(signingString, key)
	if err != nil {
		return "", err
	}

	return signingString + "." + signature, nil
}

// JWT 主要结构体
type JWT struct {
	signingMethod SigningMethod
	key           interface{}
}

// New 创建新的 JWT 实例
func New(method SigningMethod, key interface{}) *JWT {
	return &JWT{
		signingMethod: method,
		key:           key,
	}
}

// NewWithClaims 创建带声明的新令牌
func NewWithClaims(method SigningMethod, claims Claims) *Token {
	return &Token{
		Header: &Header{
			Type:      "JWT",
			Algorithm: method.Alg(),
		},
		Claims: claims,
		Method: method,
	}
}

// Parse 解析 JWT 令牌
func (j *JWT) Parse(tokenString string) (*Token, error) {
	claims := make(MapClaims)
	return j.ParseWithClaims(tokenString, claims)
}

// ParseWithClaims 解析带指定声明类型的 JWT 令牌
func (j *JWT) ParseWithClaims(tokenString string, claims Claims) (*Token, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	token := &Token{
		Raw: tokenString,
	}

	// 解析头部
	headerBytes, err := base64URLDecode(parts[0])
	if err != nil {
		return nil, err
	}

	header := &Header{}
	if err := json.Unmarshal(headerBytes, header); err != nil {
		return nil, err
	}
	token.Header = header

	// 验证签名方法
	if header.Algorithm != j.signingMethod.Alg() {
		return nil, fmt.Errorf("unexpected signing method: %s", header.Algorithm)
	}
	token.Method = j.signingMethod

	// 解析声明
	claimsBytes, err := base64URLDecode(parts[1])
	if err != nil {
		return nil, err
	}

	// 根据 claims 类型进行不同的处理
	switch c := claims.(type) {
	case MapClaims:
		var tempClaims map[string]interface{}
		if err := json.Unmarshal(claimsBytes, &tempClaims); err != nil {
			return nil, err
		}
		for k, v := range tempClaims {
			c[k] = v
		}
	default:
		if err := json.Unmarshal(claimsBytes, claims); err != nil {
			return nil, err
		}
	}
	token.Claims = claims

	// 验证签名
	signingString := parts[0] + "." + parts[1]
	token.Signature = parts[2]

	if err := j.signingMethod.Verify(signingString, parts[2], j.key); err != nil {
		return nil, err
	}

	// 验证声明
	if err := claims.Valid(); err != nil {
		return nil, err
	}

	token.Valid = true
	return token, nil
}

// Generate 生成 JWT 令牌
func (j *JWT) Generate(claims Claims) (string, error) {
	token := NewWithClaims(j.signingMethod, claims)
	return token.SignedString(j.key)
}

// Validate 验证 JWT 令牌
func (j *JWT) Validate(tokenString string) (*Token, error) {
	return j.Parse(tokenString)
}

// base64URL 编码/解码工具函数
func base64URLEncode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

func base64URLDecode(s string) ([]byte, error) {
	// 添加必要的填充
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

// HMAC 签名方法实现
type SigningMethodHMAC struct {
	Name string
	Hash crypto.Hash
}

var (
	SigningMethodHS256 = &SigningMethodHMAC{"HS256", crypto.SHA256}
	SigningMethodHS384 = &SigningMethodHMAC{"HS384", crypto.SHA384}
	SigningMethodHS512 = &SigningMethodHMAC{"HS512", crypto.SHA512}
)

func (m *SigningMethodHMAC) Alg() string {
	return m.Name
}

func (m *SigningMethodHMAC) Sign(signingString string, key interface{}) (string, error) {
	keyBytes, ok := key.([]byte)
	if !ok {
		return "", ErrInvalidKeyType
	}

	hasher := hmac.New(m.Hash.New, keyBytes)
	hasher.Write([]byte(signingString))

	return base64URLEncode(hasher.Sum(nil)), nil
}

func (m *SigningMethodHMAC) Verify(signingString, signature string, key interface{}) error {
	sig, err := m.Sign(signingString, key)
	if err != nil {
		return err
	}

	if !hmac.Equal([]byte(sig), []byte(signature)) {
		return ErrInvalidSignature
	}

	return nil
}

// RSA 签名方法实现
type SigningMethodRSA struct {
	Name string
	Hash crypto.Hash
}

var (
	SigningMethodRS256 = &SigningMethodRSA{"RS256", crypto.SHA256}
	SigningMethodRS384 = &SigningMethodRSA{"RS384", crypto.SHA384}
	SigningMethodRS512 = &SigningMethodRSA{"RS512", crypto.SHA512}
)

func (m *SigningMethodRSA) Alg() string {
	return m.Name
}

func (m *SigningMethodRSA) Sign(signingString string, key interface{}) (string, error) {
	var rsaKey *rsa.PrivateKey

	switch k := key.(type) {
	case *rsa.PrivateKey:
		rsaKey = k
	case []byte:
		var err error
		rsaKey, err = parseRSAPrivateKeyFromPEM(k)
		if err != nil {
			return "", err
		}
	default:
		return "", ErrInvalidKeyType
	}

	var hasher hash.Hash
	switch m.Hash {
	case crypto.SHA256:
		hasher = sha256.New()
	case crypto.SHA384:
		h := sha512.New384()
		hasher = h
	case crypto.SHA512:
		hasher = sha512.New()
	default:
		return "", errors.New("unsupported hash algorithm")
	}

	hasher.Write([]byte(signingString))
	hashed := hasher.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, m.Hash, hashed)
	if err != nil {
		return "", err
	}

	return base64URLEncode(signature), nil
}

func (m *SigningMethodRSA) Verify(signingString, signature string, key interface{}) error {
	var rsaKey *rsa.PublicKey

	switch k := key.(type) {
	case *rsa.PublicKey:
		rsaKey = k
	case *rsa.PrivateKey:
		rsaKey = &k.PublicKey
	case []byte:
		var err error
		rsaKey, err = parseRSAPublicKeyFromPEM(k)
		if err != nil {
			return err
		}
	default:
		return ErrInvalidKeyType
	}

	sig, err := base64URLDecode(signature)
	if err != nil {
		return err
	}

	var hasher hash.Hash
	switch m.Hash {
	case crypto.SHA256:
		hasher = sha256.New()
	case crypto.SHA384:
		h := sha512.New384()
		hasher = h
	case crypto.SHA512:
		hasher = sha512.New()
	default:
		return errors.New("unsupported hash algorithm")
	}

	hasher.Write([]byte(signingString))
	hashed := hasher.Sum(nil)

	return rsa.VerifyPKCS1v15(rsaKey, m.Hash, hashed, sig)
}

// PEM 密钥解析工具函数
func parseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, ErrKeyMustBePEM
	}

	if block.Type != "RSA PRIVATE KEY" && block.Type != "PRIVATE KEY" {
		return nil, errors.New("not an RSA private key")
	}

	if block.Type == "PRIVATE KEY" {
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		if rsaKey, ok := parsedKey.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, errors.New("not an RSA private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func parseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, ErrKeyMustBePEM
	}

	if block.Type != "RSA PUBLIC KEY" && block.Type != "PUBLIC KEY" {
		return nil, errors.New("not an RSA public key")
	}

	if block.Type == "PUBLIC KEY" {
		parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		if rsaKey, ok := parsedKey.(*rsa.PublicKey); ok {
			return rsaKey, nil
		}
		return nil, errors.New("not an RSA public key")
	}

	return x509.ParsePKCS1PublicKey(block.Bytes)
}
