package jsonx

import (
	"testing"
)

func TestBasicOperations(t *testing.T) {
	// 测试解析 JSON
	jsonStr := `{"name": "test", "age": 25, "active": true, "score": 98.5}`
	j := Parse(jsonStr)

	if j.Error() != nil {
		t.Fatalf("Failed to parse JSON: %v", j.Error())
	}

	// 测试获取不同类型的值
	if name := j.Get("name").String(); name != "test" {
		t.Errorf("Expected name='test', got '%s'", name)
	}

	if age := j.Get("age").Int(); age != 25 {
		t.Errorf("Expected age=25, got %d", age)
	}

	if active := j.Get("active").Bool(); !active {
		t.Errorf("Expected active=true, got %t", active)
	}

	if score := j.Get("score").Float64(); score != 98.5 {
		t.Errorf("Expected score=98.5, got %f", score)
	}
}

func TestPathOperations(t *testing.T) {
	j := Object()

	// 测试设置嵌套路径
	j.Set("user.profile.name", "张三")
	j.Set("user.profile.age", 30)
	j.Set("user.settings.theme", "dark")

	// 测试获取嵌套路径
	if name := j.Get("user.profile.name").String(); name != "张三" {
		t.Errorf("Expected name='张三', got '%s'", name)
	}

	if age := j.Get("user.profile.age").Int(); age != 30 {
		t.Errorf("Expected age=30, got %d", age)
	}

	// 测试路径是否存在
	if !j.Has("user.profile.name") {
		t.Error("Path 'user.profile.name' should exist")
	}

	if j.Has("user.profile.email") {
		t.Error("Path 'user.profile.email' should not exist")
	}

	// 测试删除路径
	j.Delete("user.settings")
	if j.Has("user.settings.theme") {
		t.Error("Path 'user.settings.theme' should be deleted")
	}
}

func TestArrayOperations(t *testing.T) {
	arr := Array()

	// 测试添加元素
	arr.Append("first", 42, true, nil)

	if length := arr.Length(); length != 4 {
		t.Errorf("Expected length=4, got %d", length)
	}

	// 测试获取元素
	if first := arr.Index(0).String(); first != "first" {
		t.Errorf("Expected first element='first', got '%s'", first)
	}

	if second := arr.Index(1).Int(); second != 42 {
		t.Errorf("Expected second element=42, got %d", second)
	}

	// 测试前置添加
	arr.Prepend("zero")
	if zero := arr.Index(0).String(); zero != "zero" {
		t.Errorf("Expected first element='zero', got '%s'", zero)
	}

	// 测试删除元素
	arr.Remove(0)
	if first := arr.Index(0).String(); first != "first" {
		t.Errorf("After remove, expected first element='first', got '%s'", first)
	}
}

func TestTypeChecking(t *testing.T) {
	tests := []struct {
		value    interface{}
		isObject bool
		isArray  bool
		isString bool
		isNumber bool
		isBool   bool
		isNull   bool
	}{
		{map[string]interface{}{"key": "value"}, true, false, false, false, false, false},
		{[]interface{}{1, 2, 3}, false, true, false, false, false, false},
		{"hello", false, false, true, false, false, false},
		{42, false, false, false, true, false, false},
		{3.14, false, false, false, true, false, false},
		{true, false, false, false, false, true, false},
		{nil, false, false, false, false, false, true},
	}

	for i, test := range tests {
		j := New(test.value)

		if j.IsObject() != test.isObject {
			t.Errorf("Test %d: Expected IsObject()=%t, got %t", i, test.isObject, j.IsObject())
		}

		if j.IsArray() != test.isArray {
			t.Errorf("Test %d: Expected IsArray()=%t, got %t", i, test.isArray, j.IsArray())
		}

		if j.IsString() != test.isString {
			t.Errorf("Test %d: Expected IsString()=%t, got %t", i, test.isString, j.IsString())
		}

		if j.IsNumber() != test.isNumber {
			t.Errorf("Test %d: Expected IsNumber()=%t, got %t", i, test.isNumber, j.IsNumber())
		}

		if j.IsBool() != test.isBool {
			t.Errorf("Test %d: Expected IsBool()=%t, got %t", i, test.isBool, j.IsBool())
		}

		if j.IsNull() != test.isNull {
			t.Errorf("Test %d: Expected IsNull()=%t, got %t", i, test.isNull, j.IsNull())
		}
	}
}

func TestBuilder(t *testing.T) {
	// 测试对象构建器
	obj := NewBuilder().
		AddString("name", "test").
		AddInt("age", 25).
		AddBool("active", true).
		AddFloat("score", 98.5).
		Build()

	if name := obj.Get("name").String(); name != "test" {
		t.Errorf("Expected name='test', got '%s'", name)
	}

	if age := obj.Get("age").Int(); age != 25 {
		t.Errorf("Expected age=25, got %d", age)
	}

	// 测试数组构建器
	arr := NewArrayBuilder().
		AppendString("first").
		AppendInt(42).
		AppendBool(true).
		Build()

	if length := arr.Length(); length != 3 {
		t.Errorf("Expected length=3, got %d", length)
	}

	if first := arr.Index(0).String(); first != "first" {
		t.Errorf("Expected first='first', got '%s'", first)
	}
}

func TestQuickFunctions(t *testing.T) {
	// 测试快速对象创建
	obj := QuickObject(map[string]interface{}{
		"name": "test",
		"age":  30,
	})

	if name := obj.Get("name").String(); name != "test" {
		t.Errorf("Expected name='test', got '%s'", name)
	}

	// 测试快速数组创建
	arr := QuickArray("a", "b", "c")

	if length := arr.Length(); length != 3 {
		t.Errorf("Expected length=3, got %d", length)
	}

	if first := arr.Index(0).String(); first != "a" {
		t.Errorf("Expected first='a', got '%s'", first)
	}
}

func TestSerialization(t *testing.T) {
	original := QuickObject(map[string]interface{}{
		"name":   "test",
		"age":    25,
		"active": true,
		"tags":   []interface{}{"go", "json"},
	})

	// 测试转换为 JSON 字符串
	jsonStr, err := original.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize to JSON: %v", err)
	}

	// 测试解析回来
	parsed := Parse(jsonStr)
	if parsed.Error() != nil {
		t.Fatalf("Failed to parse JSON: %v", parsed.Error())
	}

	// 验证数据完整性
	if name := parsed.Get("name").String(); name != "test" {
		t.Errorf("Expected name='test', got '%s'", name)
	}

	if age := parsed.Get("age").Int(); age != 25 {
		t.Errorf("Expected age=25, got %d", age)
	}

	// 测试美化 JSON
	prettyJSON, err := original.ToPrettyJSON()
	if err != nil {
		t.Fatalf("Failed to serialize to pretty JSON: %v", err)
	}

	if len(prettyJSON) <= len(jsonStr) {
		t.Error("Pretty JSON should be longer than compact JSON")
	}
}

func TestIterators(t *testing.T) {
	// 测试对象迭代
	obj := QuickObject(map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	})

	count := 0
	obj.ForEach(func(key string, value *JSON) bool {
		count++
		if value.Int() == 0 {
			t.Errorf("Unexpected zero value for key '%s'", key)
		}
		return true
	})

	if count != 3 {
		t.Errorf("Expected to iterate 3 times, got %d", count)
	}

	// 测试数组迭代
	arr := QuickArray(10, 20, 30)

	sum := 0
	arr.ForEach(func(key string, value *JSON) bool {
		sum += value.Int()
		return true
	})

	if sum != 60 {
		t.Errorf("Expected sum=60, got %d", sum)
	}
}

func TestMapAndFilter(t *testing.T) {
	// 测试数组映射
	numbers := QuickArray(1, 2, 3, 4, 5)

	doubled := numbers.Map(func(key string, value *JSON) interface{} {
		return value.Int() * 2
	})

	if doubled.Index(0).Int() != 2 {
		t.Errorf("Expected first doubled value=2, got %d", doubled.Index(0).Int())
	}

	if doubled.Index(4).Int() != 10 {
		t.Errorf("Expected last doubled value=10, got %d", doubled.Index(4).Int())
	}

	// 测试数组过滤
	evens := numbers.Filter(func(key string, value *JSON) bool {
		return value.Int()%2 == 0
	})

	if evens.Length() != 2 {
		t.Errorf("Expected 2 even numbers, got %d", evens.Length())
	}

	if evens.Index(0).Int() != 2 {
		t.Errorf("Expected first even=2, got %d", evens.Index(0).Int())
	}

	if evens.Index(1).Int() != 4 {
		t.Errorf("Expected second even=4, got %d", evens.Index(1).Int())
	}
}

func TestCloneAndMerge(t *testing.T) {
	original := QuickObject(map[string]interface{}{
		"name": "original",
		"age":  25,
		"settings": map[string]interface{}{
			"theme": "light",
		},
	})

	// 测试克隆
	cloned := original.Clone()
	cloned.Set("name", "cloned")

	if original.Get("name").String() == "cloned" {
		t.Error("Original should not be affected by cloned modification")
	}

	// 测试浅合并
	other := QuickObject(map[string]interface{}{
		"age":  30,
		"city": "北京",
		"settings": map[string]interface{}{
			"lang": "zh-CN",
		},
	})

	merged := original.Clone().Merge(other)

	if merged.Get("age").Int() != 30 {
		t.Errorf("Expected merged age=30, got %d", merged.Get("age").Int())
	}

	if merged.Get("city").String() != "北京" {
		t.Errorf("Expected merged city='北京', got '%s'", merged.Get("city").String())
	}

	// 浅合并应该覆盖 settings
	if merged.Has("settings.theme") {
		t.Error("Shallow merge should replace settings completely")
	}

	// 测试深度合并
	deepMerged := original.Clone().DeepMerge(other)

	if !deepMerged.Has("settings.theme") {
		t.Error("Deep merge should preserve original settings.theme")
	}

	if !deepMerged.Has("settings.lang") {
		t.Error("Deep merge should add new settings.lang")
	}
}

func TestUtilityFunctions(t *testing.T) {
	// 测试 JSON 验证
	validJSON := `{"name": "test", "age": 25}`
	invalidJSON := `{"name": "test", "age":}`

	if !IsValid(validJSON) {
		t.Error("Valid JSON should be recognized as valid")
	}

	if IsValid(invalidJSON) {
		t.Error("Invalid JSON should be recognized as invalid")
	}

	// 测试格式化
	compactJSON := `{"name":"test","age":25}`
	prettyJSON, err := Pretty(compactJSON)
	if err != nil {
		t.Fatalf("Failed to prettify JSON: %v", err)
	}

	if len(prettyJSON) <= len(compactJSON) {
		t.Error("Pretty JSON should be longer than compact JSON")
	}

	// 测试压缩
	minifiedJSON, err := Minify(prettyJSON)
	if err != nil {
		t.Fatalf("Failed to minify JSON: %v", err)
	}

	if len(minifiedJSON) >= len(prettyJSON) {
		t.Error("Minified JSON should be shorter than pretty JSON")
	}
}

func TestFlatten(t *testing.T) {
	nested := QuickObject(map[string]interface{}{
		"user": map[string]interface{}{
			"name": "test",
			"profile": map[string]interface{}{
				"age": 25,
			},
		},
		"items": []interface{}{"a", "b"},
	})

	flattened := Flatten(nested)

	expectedKeys := []string{
		"user.name",
		"user.profile.age",
		"items.0",
		"items.1",
	}

	for _, key := range expectedKeys {
		if _, exists := flattened[key]; !exists {
			t.Errorf("Expected flattened key '%s' not found", key)
		}
	}

	if flattened["user.name"] != "test" {
		t.Errorf("Expected user.name='test', got '%v'", flattened["user.name"])
	}

	// 测试反扁平化
	unflattened := Unflatten(flattened)

	if unflattened.Get("user.name").String() != "test" {
		t.Error("Unflatten should restore original structure")
	}

	if unflattened.Get("user.profile.age").Int() != 25 {
		t.Error("Unflatten should restore nested values")
	}
}

func TestPickAndOmit(t *testing.T) {
	original := QuickObject(map[string]interface{}{
		"id":       1,
		"name":     "test",
		"email":    "test@example.com",
		"password": "secret",
		"internal": "data",
	})

	// 测试选择字段
	picked := Pick(original, "id", "name", "email")

	if picked.Length() != 3 {
		t.Errorf("Expected picked object to have 3 fields, got %d", picked.Length())
	}

	if !picked.Has("id") || !picked.Has("name") || !picked.Has("email") {
		t.Error("Picked object should contain specified fields")
	}

	if picked.Has("password") || picked.Has("internal") {
		t.Error("Picked object should not contain unspecified fields")
	}

	// 测试排除字段
	omitted := Omit(original, "password", "internal")

	if omitted.Length() != 3 {
		t.Errorf("Expected omitted object to have 3 fields, got %d", omitted.Length())
	}

	if omitted.Has("password") || omitted.Has("internal") {
		t.Error("Omitted object should not contain excluded fields")
	}

	if !omitted.Has("id") || !omitted.Has("name") || !omitted.Has("email") {
		t.Error("Omitted object should contain non-excluded fields")
	}
}

func TestStructConversion(t *testing.T) {
	type User struct {
		Name   string   `json:"name"`
		Age    int      `json:"age"`
		Active bool     `json:"active"`
		Tags   []string `json:"tags"`
	}

	user := User{
		Name:   "张三",
		Age:    30,
		Active: true,
		Tags:   []string{"developer", "go"},
	}

	j := FromStruct(user)

	if j.Error() != nil {
		t.Fatalf("Failed to convert struct: %v", j.Error())
	}

	if j.Get("name").String() != "张三" {
		t.Errorf("Expected name='张三', got '%s'", j.Get("name").String())
	}

	if j.Get("age").Int() != 30 {
		t.Errorf("Expected age=30, got %d", j.Get("age").Int())
	}

	if !j.Get("active").Bool() {
		t.Error("Expected active=true")
	}

	if j.Get("tags").Length() != 2 {
		t.Errorf("Expected tags length=2, got %d", j.Get("tags").Length())
	}
}

func TestErrorHandling(t *testing.T) {
	// 测试无效 JSON 解析
	invalid := Parse(`{"invalid": json}`)
	if invalid.Error() == nil {
		t.Error("Invalid JSON should produce an error")
	}

	// 测试链式调用中的错误传播
	result := invalid.Set("key", "value").Get("key")
	if result.Error() == nil {
		t.Error("Error should propagate through chain")
	}

	// 测试数组索引越界
	arr := QuickArray(1, 2, 3)
	outOfBounds := arr.Index(10)
	if outOfBounds.Error() == nil {
		t.Error("Array index out of bounds should produce an error")
	}

	// 测试非数组调用数组方法
	obj := QuickObject(map[string]interface{}{"key": "value"})
	appendResult := obj.Append("value")
	if appendResult.Error() == nil {
		t.Error("Append on non-array should produce an error")
	}
}
