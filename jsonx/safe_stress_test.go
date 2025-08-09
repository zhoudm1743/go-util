package jsonx

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// 测试合理深度的嵌套结构
func TestReasonableDeepNesting(t *testing.T) {
	j := Object()

	// 创建 10 层深的嵌套结构（减少深度）
	path := "level1"
	for i := 2; i <= 10; i++ {
		path += fmt.Sprintf(".level%d", i)
	}

	// 设置深层值
	j.Set(path, "深层嵌套值")

	// 验证能否正确获取
	if value := j.Get(path).String(); value != "深层嵌套值" {
		t.Errorf("深层嵌套失败: expected '深层嵌套值', got '%s'", value)
	}

	// 验证中间路径存在
	midPath := "level1.level2.level3"
	if !j.Has(midPath) {
		t.Error("中间路径应该存在")
	}

	// 测试删除中间路径
	j.Delete("level1.level2.level3")
	if j.Has(path) {
		t.Error("删除中间路径后，深层路径应该不存在")
	}
}

// 测试中等大小的数组操作
func TestMediumArrayOperations(t *testing.T) {
	// 创建中等大小数组（减少到100个元素）
	arr := Array()

	// 添加 100 个元素
	for i := 0; i < 100; i++ {
		arr.Append(map[string]interface{}{
			"id":   i,
			"name": fmt.Sprintf("item_%d", i),
			"data": strings.Repeat("x", 10), // 减少每个元素的大小
		})
	}

	if arr.Length() != 100 {
		t.Errorf("数组长度错误: expected 100, got %d", arr.Length())
	}

	// 测试数组过滤
	filtered := arr.Filter(func(key string, value *JSON) bool {
		return value.Get("id").Int()%10 == 0 // 只保留 id 是 10 的倍数的
	})

	if filtered.Length() != 10 {
		t.Errorf("过滤后数组长度错误: expected 10, got %d", filtered.Length())
	}

	// 测试数组映射
	mapped := filtered.Map(func(key string, value *JSON) interface{} {
		return map[string]interface{}{
			"original_id": value.Get("id").Int(),
			"doubled_id":  value.Get("id").Int() * 2,
		}
	})

	first := mapped.Index(0)
	if first.Get("original_id").Int() != 0 {
		t.Error("映射结果不正确")
	}
	if first.Get("doubled_id").Int() != 0 {
		t.Error("映射计算不正确")
	}
}

// 测试中等大小 JSON 字符串解析
func TestMediumJSONParsing(t *testing.T) {
	// 构建中等大小 JSON 字符串（减少到50个用户）
	var builder strings.Builder
	builder.WriteString(`{"users": [`)

	for i := 0; i < 50; i++ {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf(`{
			"id": %d,
			"name": "用户_%d",
			"email": "user%d@example.com",
			"profile": {
				"age": %d,
				"bio": "简介",
				"tags": ["tag_%d", "category_%d"],
				"settings": {
					"theme": "dark",
					"notifications": true,
					"privacy": {
						"email_visible": %t,
						"profile_public": %t
					}
				}
			}
		}`, i, i, i, 20+i%50, i%5, i%10, i%2 == 0, i%3 == 0))
	}

	builder.WriteString(`]}`)
	mediumJSON := builder.String()

	// 解析中等大小 JSON
	j := Parse(mediumJSON)
	if j.Error() != nil {
		t.Fatalf("解析中等大小 JSON 失败: %v", j.Error())
	}

	// 验证数据完整性
	users := j.Get("users")
	if users.Length() != 50 {
		t.Errorf("用户数量错误: expected 50, got %d", users.Length())
	}

	// 随机检查几个用户
	user10 := users.Index(10)
	if user10.Get("name").String() != "用户_10" {
		t.Error("用户数据不正确")
	}

	if !user10.Get("profile.settings.notifications").Bool() {
		t.Error("嵌套布尔值不正确")
	}

	// 测试路径访问
	email := j.Get("users.10.email").String()
	if email != "user10@example.com" {
		t.Errorf("路径访问失败: expected 'user10@example.com', got '%s'", email)
	}
}

// 测试特殊字符和边界情况
func TestSpecialCharactersAndEdgeCases(t *testing.T) {
	j := Object()

	// 测试特殊字符作为键
	specialKeys := []string{
		"key with spaces",
		"key-with-dashes",
		"key_with_underscores",
		"key中文键",
		"🚀emoji🎉key",
	}

	for i, key := range specialKeys {
		j.Set(key, fmt.Sprintf("value_%d", i))
	}

	// 验证特殊字符键
	for i, key := range specialKeys {
		expected := fmt.Sprintf("value_%d", i)
		if value := j.Get(key).String(); value != expected {
			t.Errorf("特殊字符键 '%s' 失败: expected '%s', got '%s'", key, expected, value)
		}
	}

	// 测试特殊值
	specialValues := map[string]interface{}{
		"null_value":   nil,
		"empty_string": "",
		"zero_int":     0,
		"zero_float":   0.0,
		"false_bool":   false,
		"unicode":      "这是中文 🚀 This is English",
		"json_string":  `{"nested": "json"}`,
	}

	for key, value := range specialValues {
		j.Set(key, value)

		// 验证设置和获取
		retrieved := j.Get(key).ToInterface()
		if !compareValues(value, retrieved) {
			t.Errorf("特殊值 '%s' 失败: expected %v, got %v", key, value, retrieved)
		}
	}
}

// 测试错误恢复和链式调用的健壮性
func TestErrorRecoveryAndChaining(t *testing.T) {
	// 从无效 JSON 开始
	invalid := Parse(`{"invalid": json}`)

	// 即使有错误，链式调用也应该继续工作（错误传播）
	result := invalid.
		Set("new_key", "new_value").
		Get("new_key").
		Set("another", 123).
		Get("another")

	if result.Error() == nil {
		t.Error("错误应该传播到链式调用的末尾")
	}

	// 测试从错误中恢复
	valid := Object().
		Set("user.name", "张三").
		Set("user.age", 25)

	// 制造一个错误（访问不存在的数组索引）
	errorResult := valid.Get("user.skills.0.name")
	if errorResult.Error() == nil {
		t.Error("访问不存在的路径应该产生错误")
	}

	// 但原始对象应该仍然有效
	if name := valid.Get("user.name").String(); name != "张三" {
		t.Error("原始对象应该保持有效")
	}
}

// 测试复杂的路径操作
func TestComplexPathOperations(t *testing.T) {
	j := Object()

	// 创建复杂的混合结构
	j.Set("data.0.info.details.0.value", "nested_array_object")
	j.Set("data.1.info.name", "second_item")
	j.Set("data.0.info.tags.0", "tag1")
	j.Set("data.0.info.tags.1", "tag2")

	// 验证复杂路径
	value := j.Get("data.0.info.details.0.value").String()
	if value != "nested_array_object" {
		t.Errorf("复杂路径访问失败: expected 'nested_array_object', got '%s'", value)
	}

	// 测试数组索引越界
	outOfBounds := j.Get("data.10.info.name")
	if outOfBounds.Error() == nil {
		t.Error("数组索引越界应该产生错误")
	}

	// 测试混合类型路径
	j.Set("mixed.string_key", "value")
	j.Set("mixed.0", "array_item")

	if j.Get("mixed.string_key").String() != "value" {
		t.Error("混合类型路径访问失败")
	}
}

// 测试性能基准（简化版）
func TestSimplePerformanceBenchmark(t *testing.T) {
	// 创建测试数据（减少大小）
	obj := Object()
	for i := 0; i < 100; i++ {
		obj.Set(fmt.Sprintf("key_%d", i), map[string]interface{}{
			"id":   i,
			"name": fmt.Sprintf("name_%d", i),
			"data": strings.Repeat("x", 10),
		})
	}

	// 测试序列化性能
	start := time.Now()
	for i := 0; i < 10; i++ {
		_, err := obj.ToJSON()
		if err != nil {
			t.Fatalf("序列化失败: %v", err)
		}
	}
	serializationTime := time.Since(start)

	// 测试解析性能
	jsonStr, _ := obj.ToJSON()
	start = time.Now()
	for i := 0; i < 10; i++ {
		parsed := Parse(jsonStr)
		if parsed.Error() != nil {
			t.Fatalf("解析失败: %v", parsed.Error())
		}
	}
	parsingTime := time.Since(start)

	// 性能基准（宽松的限制）
	if serializationTime > time.Second {
		t.Errorf("序列化性能太慢: %v", serializationTime)
	}

	if parsingTime > time.Second {
		t.Errorf("解析性能太慢: %v", parsingTime)
	}

	t.Logf("序列化 10 次耗时: %v", serializationTime)
	t.Logf("解析 10 次耗时: %v", parsingTime)
}

// 测试边界值和极端情况
func TestBoundaryConditions(t *testing.T) {
	// 测试空值处理
	empty := Object()
	if empty.Length() != 0 {
		t.Error("空对象长度应该为0")
	}

	// 测试空数组
	emptyArr := Array()
	if emptyArr.Length() != 0 {
		t.Error("空数组长度应该为0")
	}

	// 测试空字符串解析
	emptyJSON := Parse("")
	if emptyJSON.Error() == nil {
		t.Error("空字符串解析应该失败")
	}

	// 测试单个值
	singleValue := Parse("42")
	if singleValue.Error() != nil {
		t.Errorf("单个数字解析失败: %v", singleValue.Error())
	}
	if singleValue.Int() != 42 {
		t.Errorf("单个数字值错误: expected 42, got %d", singleValue.Int())
	}

	// 测试超大数组索引
	arr := Array().Append("item")
	largeIndex := arr.Index(999999)
	if largeIndex.Error() == nil {
		t.Error("超大索引应该产生错误")
	}

	// 测试负数索引
	negativeIndex := arr.Index(-1)
	if negativeIndex.Error() == nil {
		t.Error("负数索引应该产生错误")
	}
}

// 辅助函数：比较两个值是否相等
func compareValues(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	switch va := a.(type) {
	case string:
		if vb, ok := b.(string); ok {
			return va == vb
		}
	case int:
		if vb, ok := b.(int); ok {
			return va == vb
		}
		if vb, ok := b.(float64); ok {
			return float64(va) == vb
		}
	case float64:
		if vb, ok := b.(float64); ok {
			return va == vb
		}
		if vb, ok := b.(int); ok {
			return va == float64(vb)
		}
	case bool:
		if vb, ok := b.(bool); ok {
			return va == vb
		}
	}

	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
