package jsonx

import (
	"fmt"
	"strings"
)

// Builder JSON 构建器，支持链式调用
type Builder struct {
	json *JSON
}

// NewBuilder 创建新的构建器
func NewBuilder() *Builder {
	return &Builder{json: Object()}
}

// NewArrayBuilder 创建数组构建器
func NewArrayBuilder() *Builder {
	return &Builder{json: Array()}
}

// AddString 添加字符串字段
func (b *Builder) AddString(key, value string) *Builder {
	b.json.Set(key, value)
	return b
}

// AddInt 添加整数字段
func (b *Builder) AddInt(key string, value int) *Builder {
	b.json.Set(key, value)
	return b
}

// AddInt64 添加 int64 字段
func (b *Builder) AddInt64(key string, value int64) *Builder {
	b.json.Set(key, value)
	return b
}

// AddFloat 添加浮点数字段
func (b *Builder) AddFloat(key string, value float64) *Builder {
	b.json.Set(key, value)
	return b
}

// AddBool 添加布尔字段
func (b *Builder) AddBool(key string, value bool) *Builder {
	b.json.Set(key, value)
	return b
}

// AddObject 添加对象字段
func (b *Builder) AddObject(key string, value *JSON) *Builder {
	b.json.Set(key, value.data)
	return b
}

// AddArray 添加数组字段
func (b *Builder) AddArray(key string, value *JSON) *Builder {
	b.json.Set(key, value.data)
	return b
}

// AddRaw 添加原始值
func (b *Builder) AddRaw(key string, value interface{}) *Builder {
	b.json.Set(key, value)
	return b
}

// AddNull 添加 null 字段
func (b *Builder) AddNull(key string) *Builder {
	b.json.Set(key, nil)
	return b
}

// AddIf 条件添加字段
func (b *Builder) AddIf(condition bool, key string, value interface{}) *Builder {
	if condition {
		b.json.Set(key, value)
	}
	return b
}

// AddStringIf 条件添加字符串字段
func (b *Builder) AddStringIf(condition bool, key, value string) *Builder {
	if condition {
		b.json.Set(key, value)
	}
	return b
}

// AddMany 批量添加字段
func (b *Builder) AddMany(fields map[string]interface{}) *Builder {
	for k, v := range fields {
		b.json.Set(k, v)
	}
	return b
}

// 数组构建器方法

// AppendString 向数组添加字符串
func (b *Builder) AppendString(value string) *Builder {
	b.json.Append(value)
	return b
}

// AppendInt 向数组添加整数
func (b *Builder) AppendInt(value int) *Builder {
	b.json.Append(value)
	return b
}

// AppendFloat 向数组添加浮点数
func (b *Builder) AppendFloat(value float64) *Builder {
	b.json.Append(value)
	return b
}

// AppendBool 向数组添加布尔值
func (b *Builder) AppendBool(value bool) *Builder {
	b.json.Append(value)
	return b
}

// AppendObject 向数组添加对象
func (b *Builder) AppendObject(value *JSON) *Builder {
	b.json.Append(value.data)
	return b
}

// AppendArray 向数组添加数组
func (b *Builder) AppendArray(value *JSON) *Builder {
	b.json.Append(value.data)
	return b
}

// AppendRaw 向数组添加原始值
func (b *Builder) AppendRaw(value interface{}) *Builder {
	b.json.Append(value)
	return b
}

// AppendNull 向数组添加 null
func (b *Builder) AppendNull() *Builder {
	b.json.Append(nil)
	return b
}

// AppendMany 向数组批量添加值
func (b *Builder) AppendMany(values ...interface{}) *Builder {
	b.json.Append(values...)
	return b
}

// Build 构建最终的 JSON
func (b *Builder) Build() *JSON {
	return b.json
}

// BuildString 构建 JSON 字符串
func (b *Builder) BuildString() (string, error) {
	return b.json.ToJSON()
}

// BuildPrettyString 构建格式化的 JSON 字符串
func (b *Builder) BuildPrettyString() (string, error) {
	return b.json.ToPrettyJSON()
}

// 快速构建器函数

// QuickObject 快速创建对象
func QuickObject(fields map[string]interface{}) *JSON {
	return NewBuilder().AddMany(fields).Build()
}

// QuickArray 快速创建数组
func QuickArray(values ...interface{}) *JSON {
	return NewArrayBuilder().AppendMany(values...).Build()
}

// 模板构建器

// TemplateBuilder 模板构建器
type TemplateBuilder struct {
	template string
	values   map[string]interface{}
}

// NewTemplate 创建模板构建器
func NewTemplate(template string) *TemplateBuilder {
	return &TemplateBuilder{
		template: template,
		values:   make(map[string]interface{}),
	}
}

// Set 设置模板变量
func (t *TemplateBuilder) Set(key string, value interface{}) *TemplateBuilder {
	t.values[key] = value
	return t
}

// SetMany 批量设置模板变量
func (t *TemplateBuilder) SetMany(values map[string]interface{}) *TemplateBuilder {
	for k, v := range values {
		t.values[k] = v
	}
	return t
}

// Build 构建 JSON（简单的字符串替换）
func (t *TemplateBuilder) Build() *JSON {
	result := t.template
	for k, v := range t.values {
		placeholder := fmt.Sprintf("{{%s}}", k)
		var replacement string

		switch val := v.(type) {
		case string:
			replacement = fmt.Sprintf(`"%s"`, escapeJSONString(val))
		case nil:
			replacement = "null"
		case bool:
			if val {
				replacement = "true"
			} else {
				replacement = "false"
			}
		default:
			replacement = fmt.Sprintf("%v", val)
		}

		result = strings.ReplaceAll(result, placeholder, replacement)
	}

	return Parse(result)
}

// 实用工具函数

// Merge 合并多个 JSON 对象
func Merge(jsons ...*JSON) *JSON {
	if len(jsons) == 0 {
		return Object()
	}

	result := jsons[0].Clone()
	for i := 1; i < len(jsons); i++ {
		result = result.Merge(jsons[i])
		if result.err != nil {
			return result
		}
	}

	return result
}

// DeepMergeAll 深度合并多个 JSON 对象
func DeepMergeAll(jsons ...*JSON) *JSON {
	if len(jsons) == 0 {
		return Object()
	}

	result := jsons[0].Clone()
	for i := 1; i < len(jsons); i++ {
		result = result.DeepMerge(jsons[i])
		if result.err != nil {
			return result
		}
	}

	return result
}

// Compare 比较两个 JSON 是否相等
func Compare(j1, j2 *JSON) bool {
	if j1.err != nil || j2.err != nil {
		return false
	}

	json1, err1 := j1.ToJSON()
	json2, err2 := j2.ToJSON()

	if err1 != nil || err2 != nil {
		return false
	}

	return json1 == json2
}

// Flatten 扁平化 JSON 对象
func Flatten(j *JSON) map[string]interface{} {
	result := make(map[string]interface{})
	flattenRecursive(j.data, "", result)
	return result
}

// flattenRecursive 递归扁平化
func flattenRecursive(data interface{}, prefix string, result map[string]interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for k, val := range v {
			key := k
			if prefix != "" {
				key = prefix + "." + k
			}
			flattenRecursive(val, key, result)
		}
	case []interface{}:
		for i, val := range v {
			key := fmt.Sprintf("%d", i)
			if prefix != "" {
				key = prefix + "." + key
			}
			flattenRecursive(val, key, result)
		}
	default:
		result[prefix] = v
	}
}

// Unflatten 反扁平化为 JSON 对象
func Unflatten(flat map[string]interface{}) *JSON {
	result := Object()

	for path, value := range flat {
		result.Set(path, value)
	}

	return result
}

// Transform 转换 JSON 结构
func Transform(j *JSON, transformer func(key string, value *JSON) interface{}) *JSON {
	return j.Map(transformer)
}

// Pick 选择指定字段
func Pick(j *JSON, fields ...string) *JSON {
	if !j.IsObject() {
		return &JSON{err: fmt.Errorf("not an object")}
	}

	result := Object()
	for _, field := range fields {
		if j.Has(field) {
			value := j.Get(field)
			result.Set(field, value.data)
		}
	}

	return result
}

// Omit 排除指定字段
func Omit(j *JSON, fields ...string) *JSON {
	if !j.IsObject() {
		return &JSON{err: fmt.Errorf("not an object")}
	}

	result := j.Clone()
	for _, field := range fields {
		result.Delete(field)
	}

	return result
}

// Schema 简单的 JSON Schema 验证
type Schema struct {
	Type       string             `json:"type"`
	Properties map[string]*Schema `json:"properties,omitempty"`
	Items      *Schema            `json:"items,omitempty"`
	Required   []string           `json:"required,omitempty"`
	MinLength  *int               `json:"minLength,omitempty"`
	MaxLength  *int               `json:"maxLength,omitempty"`
	Minimum    *float64           `json:"minimum,omitempty"`
	Maximum    *float64           `json:"maximum,omitempty"`
}

// Validate 验证 JSON 是否符合 Schema
func (s *Schema) Validate(j *JSON) error {
	return s.validateValue(j, "")
}

// validateValue 验证值
func (s *Schema) validateValue(j *JSON, path string) error {
	switch s.Type {
	case "object":
		if !j.IsObject() {
			return fmt.Errorf("expected object at %s", path)
		}

		// 验证必需字段
		for _, required := range s.Required {
			if !j.Has(required) {
				return fmt.Errorf("missing required field '%s' at %s", required, path)
			}
		}

		// 验证属性
		if s.Properties != nil {
			for prop, schema := range s.Properties {
				if j.Has(prop) {
					propPath := path + "." + prop
					if path == "" {
						propPath = prop
					}
					if err := schema.validateValue(j.Get(prop), propPath); err != nil {
						return err
					}
				}
			}
		}

	case "array":
		if !j.IsArray() {
			return fmt.Errorf("expected array at %s", path)
		}

		if s.Items != nil {
			length := j.Length()
			for i := 0; i < length; i++ {
				itemPath := fmt.Sprintf("%s[%d]", path, i)
				if err := s.Items.validateValue(j.Index(i), itemPath); err != nil {
					return err
				}
			}
		}

	case "string":
		if !j.IsString() {
			return fmt.Errorf("expected string at %s", path)
		}

		str := j.String()
		if s.MinLength != nil && len(str) < *s.MinLength {
			return fmt.Errorf("string too short at %s", path)
		}
		if s.MaxLength != nil && len(str) > *s.MaxLength {
			return fmt.Errorf("string too long at %s", path)
		}

	case "number":
		if !j.IsNumber() {
			return fmt.Errorf("expected number at %s", path)
		}

		num := j.Float64()
		if s.Minimum != nil && num < *s.Minimum {
			return fmt.Errorf("number too small at %s", path)
		}
		if s.Maximum != nil && num > *s.Maximum {
			return fmt.Errorf("number too large at %s", path)
		}

	case "boolean":
		if !j.IsBool() {
			return fmt.Errorf("expected boolean at %s", path)
		}

	case "null":
		if !j.IsNull() {
			return fmt.Errorf("expected null at %s", path)
		}
	}

	return nil
}

// 辅助函数

// escapeJSONString 转义 JSON 字符串
func escapeJSONString(s string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		"\n", `\n`,
		"\r", `\r`,
		"\t", `\t`,
		"\b", `\b`,
		"\f", `\f`,
	)
	return replacer.Replace(s)
}

// Pretty 格式化 JSON 字符串
func Pretty(jsonStr string) (string, error) {
	j := Parse(jsonStr)
	if j.err != nil {
		return "", j.err
	}
	return j.ToPrettyJSON()
}

// Minify 压缩 JSON 字符串
func Minify(jsonStr string) (string, error) {
	j := Parse(jsonStr)
	if j.err != nil {
		return "", j.err
	}
	return j.ToJSON()
}

// IsValid 检查 JSON 字符串是否有效
func IsValid(jsonStr string) bool {
	j := Parse(jsonStr)
	return j.err == nil
}

// GetType 获取 JSON 值的类型
func GetType(j *JSON) string {
	if j.IsObject() {
		return "object"
	}
	if j.IsArray() {
		return "array"
	}
	if j.IsString() {
		return "string"
	}
	if j.IsNumber() {
		return "number"
	}
	if j.IsBool() {
		return "boolean"
	}
	if j.IsNull() {
		return "null"
	}
	return "unknown"
}

// Size 计算 JSON 的大小（字节）
func Size(j *JSON) int {
	jsonStr, err := j.ToJSON()
	if err != nil {
		return 0
	}
	return len(jsonStr)
}

// Depth 计算 JSON 的深度
func Depth(j *JSON) int {
	return calculateDepth(j.data, 0)
}

// calculateDepth 递归计算深度
func calculateDepth(data interface{}, currentDepth int) int {
	switch v := data.(type) {
	case map[string]interface{}:
		maxDepth := currentDepth
		for _, val := range v {
			depth := calculateDepth(val, currentDepth+1)
			if depth > maxDepth {
				maxDepth = depth
			}
		}
		return maxDepth
	case []interface{}:
		maxDepth := currentDepth
		for _, val := range v {
			depth := calculateDepth(val, currentDepth+1)
			if depth > maxDepth {
				maxDepth = depth
			}
		}
		return maxDepth
	default:
		return currentDepth
	}
}
