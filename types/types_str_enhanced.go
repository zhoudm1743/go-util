package types

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// Contains 判断字符串是否包含子串
func (s XStr) Contains(substr string) bool {
	return strings.Contains(string(s), substr)
}

// ContainsAny 判断字符串是否包含任一字符
func (s XStr) ContainsAny(chars string) bool {
	return strings.ContainsAny(string(s), chars)
}

// HasPrefix 判断字符串是否以指定前缀开头
func (s XStr) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(s), prefix)
}

// HasSuffix 判断字符串是否以指定后缀结尾
func (s XStr) HasSuffix(suffix string) bool {
	return strings.HasSuffix(string(s), suffix)
}

// Index 返回子串第一次出现的位置，未找到返回 -1
func (s XStr) Index(substr string) int {
	return strings.Index(string(s), substr)
}

// LastIndex 返回子串最后一次出现的位置，未找到返回 -1
func (s XStr) LastIndex(substr string) int {
	return strings.LastIndex(string(s), substr)
}

// Count 统计子串出现的次数
func (s XStr) Count(substr string) int {
	return strings.Count(string(s), substr)
}

// Repeat 重复字符串 n 次
func (s XStr) Repeat(count int) XStr {
	return XStr(strings.Repeat(string(s), count))
}

// Reverse 反转字符串
func (s XStr) Reverse() XStr {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return XStr(runes)
}

// PadLeft 左填充到指定长度
func (s XStr) PadLeft(length int, pad string) XStr {
	if len(pad) == 0 {
		pad = " "
	}
	str := string(s)
	for len(str) < length {
		str = pad + str
	}
	if len(str) > length {
		str = str[len(str)-length:]
	}
	return XStr(str)
}

// PadRight 右填充到指定长度
func (s XStr) PadRight(length int, pad string) XStr {
	if len(pad) == 0 {
		pad = " "
	}
	str := string(s)
	for len(str) < length {
		str = str + pad
	}
	if len(str) > length {
		str = str[:length]
	}
	return XStr(str)
}

// PadCenter 居中填充到指定长度
func (s XStr) PadCenter(length int, pad string) XStr {
	if len(pad) == 0 {
		pad = " "
	}
	str := string(s)
	if len(str) >= length {
		return s
	}

	totalPad := length - len(str)
	leftPad := totalPad / 2
	rightPad := totalPad - leftPad

	return XStr(strings.Repeat(pad, leftPad) + str + strings.Repeat(pad, rightPad))
}

// Substring 截取子字符串
func (s XStr) Substring(start int) XStr {
	runes := []rune(s)
	if start < 0 || start >= len(runes) {
		return ""
	}
	return XStr(runes[start:])
}

// SubstringWithEnd 截取子字符串（指定起始和结束位置）
func (s XStr) SubstringWithEnd(start, end int) XStr {
	runes := []rune(s)
	if start < 0 {
		start = 0
	}
	if end > len(runes) {
		end = len(runes)
	}
	if start >= end {
		return ""
	}
	return XStr(runes[start:end])
}

// Left 从左边截取指定长度
func (s XStr) Left(length int) XStr {
	runes := []rune(s)
	if length >= len(runes) {
		return s
	}
	return XStr(runes[:length])
}

// Right 从右边截取指定长度
func (s XStr) Right(length int) XStr {
	runes := []rune(s)
	if length >= len(runes) {
		return s
	}
	return XStr(runes[len(runes)-length:])
}

// Mid 从中间截取指定长度
func (s XStr) Mid(start, length int) XStr {
	return s.SubstringWithEnd(start, start+length)
}

// IsNumeric 判断是否为纯数字
func (s XStr) IsNumeric() bool {
	str := string(s)
	if str == "" {
		return false
	}
	for _, char := range str {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// IsAlpha 判断是否为纯字母
func (s XStr) IsAlpha() bool {
	str := string(s)
	if str == "" {
		return false
	}
	for _, char := range str {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')) {
			return false
		}
	}
	return true
}

// IsAlphaNumeric 判断是否为字母和数字组合
func (s XStr) IsAlphaNumeric() bool {
	str := string(s)
	if str == "" {
		return false
	}
	for _, char := range str {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return false
		}
	}
	return true
}

// IsEmail 简单的邮箱格式验证
func (s XStr) IsEmail() bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(string(s))
}

// IsURL 简单的 URL 格式验证
func (s XStr) IsURL() bool {
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return urlRegex.MatchString(string(s))
}

// IsIPv4 判断是否为有效的 IPv4 地址
func (s XStr) IsIPv4() bool {
	ipRegex := regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
	return ipRegex.MatchString(string(s))
}

// MD5 计算 MD5 哈希
func (s XStr) MD5() XStr {
	hash := md5.Sum([]byte(s))
	return XStr(hex.EncodeToString(hash[:]))
}

// SHA1 计算 SHA1 哈希
func (s XStr) SHA1() XStr {
	hash := sha1.Sum([]byte(s))
	return XStr(hex.EncodeToString(hash[:]))
}

// SHA256 计算 SHA256 哈希
func (s XStr) SHA256() XStr {
	hash := sha256.Sum256([]byte(s))
	return XStr(hex.EncodeToString(hash[:]))
}

// Base64Encode Base64 编码
func (s XStr) Base64Encode() XStr {
	return XStr(base64.StdEncoding.EncodeToString([]byte(s)))
}

// Base64Decode Base64 解码
func (s XStr) Base64Decode() (XStr, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(s))
	if err != nil {
		return "", err
	}
	return XStr(decoded), nil
}

// URLEncode URL 编码
func (s XStr) URLEncode() XStr {
	return XStr(strings.ReplaceAll(strings.ReplaceAll(string(s), " ", "%20"), "&", "%26"))
}

// Slugify 转换为 URL 友好的字符串
func (s XStr) Slugify() XStr {
	str := strings.ToLower(string(s))
	// 替换空格为连字符
	str = regexp.MustCompile(`\s+`).ReplaceAllString(str, "-")
	// 移除非字母数字和连字符的字符
	str = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(str, "")
	// 移除多余的连字符
	str = regexp.MustCompile(`-+`).ReplaceAllString(str, "-")
	// 移除首尾连字符
	str = strings.Trim(str, "-")
	return XStr(str)
}

// WordCount 统计单词数量
func (s XStr) WordCount() int {
	str := strings.TrimSpace(string(s))
	if str == "" {
		return 0
	}
	words := strings.Fields(str)
	return len(words)
}

// Words 分割为单词数组
func (s XStr) Words() []string {
	return strings.Fields(string(s))
}

// Lines 分割为行数组
func (s XStr) Lines() []string {
	return strings.Split(string(s), "\n")
}

// Wrap 按指定宽度换行
func (s XStr) Wrap(width int) XStr {
	if width <= 0 {
		return s
	}

	words := s.Words()
	if len(words) == 0 {
		return s
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len() == 0 {
			currentLine.WriteString(word)
		} else if currentLine.Len()+1+len(word) <= width {
			currentLine.WriteString(" ")
			currentLine.WriteString(word)
		} else {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return XStr(strings.Join(lines, "\n"))
}

// Indent 缩进每一行
func (s XStr) Indent(indent string) XStr {
	lines := s.Lines()
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}
	return XStr(strings.Join(lines, "\n"))
}

// Escape 转义 HTML 特殊字符
func (s XStr) EscapeHTML() XStr {
	str := string(s)
	str = strings.ReplaceAll(str, "&", "&amp;")
	str = strings.ReplaceAll(str, "<", "&lt;")
	str = strings.ReplaceAll(str, ">", "&gt;")
	str = strings.ReplaceAll(str, "\"", "&quot;")
	str = strings.ReplaceAll(str, "'", "&#39;")
	return XStr(str)
}

// UnescapeHTML 反转义 HTML 特殊字符
func (s XStr) UnescapeHTML() XStr {
	str := string(s)
	str = strings.ReplaceAll(str, "&lt;", "<")
	str = strings.ReplaceAll(str, "&gt;", ">")
	str = strings.ReplaceAll(str, "&quot;", "\"")
	str = strings.ReplaceAll(str, "&#39;", "'")
	str = strings.ReplaceAll(str, "&amp;", "&")
	return XStr(str)
}

// Matches 正则表达式匹配
func (s XStr) Matches(pattern string) bool {
	matched, _ := regexp.MatchString(pattern, string(s))
	return matched
}

// FindAll 查找所有匹配的子串
func (s XStr) FindAll(pattern string) []string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	return re.FindAllString(string(s), -1)
}

// ReplaceRegex 使用正则表达式替换
func (s XStr) ReplaceRegex(pattern, replacement string) XStr {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return s
	}
	return XStr(re.ReplaceAllString(string(s), replacement))
}

// Random 生成指定长度的随机字符串
func RandomString(length int) XStr {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return XStr(b)
}

// RandomNumeric 生成指定长度的随机数字字符串
func RandomNumeric(length int) XStr {
	const charset = "0123456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return XStr(b)
}

// RandomAlpha 生成指定长度的随机字母字符串
func RandomAlpha(length int) XStr {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return XStr(b)
}

// UUID 生成简单的 UUID（非标准实现）
func GenerateUUID() XStr {
	return XStr(fmt.Sprintf("%s-%s-%s-%s-%s",
		RandomString(8),
		RandomString(4),
		RandomString(4),
		RandomString(4),
		RandomString(12)))
}

// Template 简单的模板替换
func (s XStr) Template(data map[string]interface{}) XStr {
	str := string(s)
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		str = strings.ReplaceAll(str, placeholder, fmt.Sprintf("%v", value))
	}
	return XStr(str)
}

// Similarity 计算与另一个字符串的相似度（简单实现）
func (s XStr) Similarity(other XStr) float64 {
	str1, str2 := string(s), string(other)
	if str1 == str2 {
		return 1.0
	}

	maxLen := len(str1)
	if len(str2) > maxLen {
		maxLen = len(str2)
	}

	if maxLen == 0 {
		return 1.0
	}

	// 简单的编辑距离计算
	distance := levenshteinDistance(str1, str2)
	return 1.0 - float64(distance)/float64(maxLen)
}

// levenshteinDistance 计算编辑距离
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}

	for j := 1; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// min 返回三个数中的最小值
func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}
