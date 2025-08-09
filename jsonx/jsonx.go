// Package jsonx 提供了一个高效、简单易用的 JSON 操作库
// 支持链式调用，不依赖第三方库，专注于性能和易用性
package jsonx

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

// JSON 主要结构体，支持链式调用
type JSON struct {
	data interface{}
	err  error
}

// New 创建一个新的 JSON 实例
func New(data interface{}) *JSON {
	return &JSON{data: data}
}

// Parse 解析 JSON 字符串
func Parse(jsonStr string) *JSON {
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	return &JSON{data: data, err: err}
}

// ParseBytes 解析 JSON 字节数组
func ParseBytes(jsonBytes []byte) *JSON {
	var data interface{}
	err := json.Unmarshal(jsonBytes, &data)
	return &JSON{data: data, err: err}
}

// Object 创建一个新的 JSON 对象
func Object() *JSON {
	return &JSON{data: make(map[string]interface{})}
}

// Array 创建一个新的 JSON 数组
func Array() *JSON {
	return &JSON{data: make([]interface{}, 0)}
}

// FromMap 从 map 创建 JSON
func FromMap(m map[string]interface{}) *JSON {
	return &JSON{data: m}
}

// FromSlice 从 slice 创建 JSON
func FromSlice(s []interface{}) *JSON {
	return &JSON{data: s}
}

// FromStruct 从结构体创建 JSON
func FromStruct(v interface{}) *JSON {
	data, err := structToMap(v)
	return &JSON{data: data, err: err}
}

// 核心访问方法

// Get 获取指定路径的值
func (j *JSON) Get(path string) *JSON {
	if j.err != nil {
		return j
	}

	value, err := j.getByPath(path)
	return &JSON{data: value, err: err}
}

// Set 设置指定路径的值
func (j *JSON) Set(path string, value interface{}) *JSON {
	if j.err != nil {
		return j
	}

	err := j.setByPath(path, value)
	return &JSON{data: j.data, err: err}
}

// Delete 删除指定路径的值
func (j *JSON) Delete(path string) *JSON {
	if j.err != nil {
		return j
	}

	err := j.deleteByPath(path)
	return &JSON{data: j.data, err: err}
}

// Has 检查指定路径是否存在
func (j *JSON) Has(path string) bool {
	if j.err != nil {
		return false
	}

	_, err := j.getByPath(path)
	return err == nil
}

// 类型检查方法

// IsObject 检查是否为对象
func (j *JSON) IsObject() bool {
	_, ok := j.data.(map[string]interface{})
	return ok
}

// IsArray 检查是否为数组
func (j *JSON) IsArray() bool {
	_, ok := j.data.([]interface{})
	return ok
}

// IsString 检查是否为字符串
func (j *JSON) IsString() bool {
	_, ok := j.data.(string)
	return ok
}

// IsNumber 检查是否为数字
func (j *JSON) IsNumber() bool {
	switch j.data.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

// IsBool 检查是否为布尔值
func (j *JSON) IsBool() bool {
	_, ok := j.data.(bool)
	return ok
}

// IsNull 检查是否为 null
func (j *JSON) IsNull() bool {
	return j.data == nil
}

// 类型转换方法

// String 转换为字符串
func (j *JSON) String() string {
	if j.err != nil {
		return ""
	}

	switch v := j.data.(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Int 转换为整数
func (j *JSON) Int() int {
	if j.err != nil {
		return 0
	}

	switch v := j.data.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// Int64 转换为 int64
func (j *JSON) Int64() int64 {
	if j.err != nil {
		return 0
	}

	switch v := j.data.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

// Float64 转换为 float64
func (j *JSON) Float64() float64 {
	if j.err != nil {
		return 0
	}

	switch v := j.data.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0
}

// Bool 转换为布尔值
func (j *JSON) Bool() bool {
	if j.err != nil {
		return false
	}

	switch v := j.data.(type) {
	case bool:
		return v
	case string:
		return v == "true"
	case int:
		return v != 0
	case float64:
		return v != 0
	}
	return false
}

// 数组操作方法

// Length 获取数组或对象的长度
func (j *JSON) Length() int {
	if j.err != nil {
		return 0
	}

	switch v := j.data.(type) {
	case []interface{}:
		return len(v)
	case map[string]interface{}:
		return len(v)
	case string:
		return len(v)
	}
	return 0
}

// Index 获取数组指定索引的元素
func (j *JSON) Index(i int) *JSON {
	if j.err != nil {
		return j
	}

	arr, ok := j.data.([]interface{})
	if !ok {
		return &JSON{err: fmt.Errorf("not an array")}
	}

	if i < 0 || i >= len(arr) {
		return &JSON{err: fmt.Errorf("index out of range")}
	}

	return &JSON{data: arr[i]}
}

// Append 向数组添加元素
func (j *JSON) Append(values ...interface{}) *JSON {
	if j.err != nil {
		return j
	}

	arr, ok := j.data.([]interface{})
	if !ok {
		return &JSON{data: j.data, err: fmt.Errorf("not an array")}
	}

	j.data = append(arr, values...)
	return j
}

// Prepend 向数组开头添加元素
func (j *JSON) Prepend(values ...interface{}) *JSON {
	if j.err != nil {
		return j
	}

	arr, ok := j.data.([]interface{})
	if !ok {
		return &JSON{data: j.data, err: fmt.Errorf("not an array")}
	}

	j.data = append(values, arr...)
	return j
}

// Remove 删除数组指定索引的元素
func (j *JSON) Remove(index int) *JSON {
	if j.err != nil {
		return j
	}

	arr, ok := j.data.([]interface{})
	if !ok {
		return &JSON{data: j.data, err: fmt.Errorf("not an array")}
	}

	if index < 0 || index >= len(arr) {
		return &JSON{data: j.data, err: fmt.Errorf("index out of range")}
	}

	j.data = append(arr[:index], arr[index+1:]...)
	return j
}

// 对象操作方法

// Keys 获取对象的所有键
func (j *JSON) Keys() []string {
	if j.err != nil {
		return nil
	}

	obj, ok := j.data.(map[string]interface{})
	if !ok {
		return nil
	}

	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Values 获取对象的所有值
func (j *JSON) Values() []*JSON {
	if j.err != nil {
		return nil
	}

	obj, ok := j.data.(map[string]interface{})
	if !ok {
		return nil
	}

	values := make([]*JSON, 0, len(obj))
	for _, v := range obj {
		values = append(values, &JSON{data: v})
	}
	return values
}

// 迭代方法

// ForEach 遍历数组或对象
func (j *JSON) ForEach(fn func(key string, value *JSON) bool) *JSON {
	if j.err != nil {
		return j
	}

	switch v := j.data.(type) {
	case []interface{}:
		for i, item := range v {
			if !fn(strconv.Itoa(i), &JSON{data: item}) {
				break
			}
		}
	case map[string]interface{}:
		for k, item := range v {
			if !fn(k, &JSON{data: item}) {
				break
			}
		}
	}

	return j
}

// Map 映射数组或对象的值
func (j *JSON) Map(fn func(key string, value *JSON) interface{}) *JSON {
	if j.err != nil {
		return j
	}

	switch v := j.data.(type) {
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = fn(strconv.Itoa(i), &JSON{data: item})
		}
		return &JSON{data: result}

	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, item := range v {
			result[k] = fn(k, &JSON{data: item})
		}
		return &JSON{data: result}
	}

	return j
}

// Filter 过滤数组或对象
func (j *JSON) Filter(fn func(key string, value *JSON) bool) *JSON {
	if j.err != nil {
		return j
	}

	switch v := j.data.(type) {
	case []interface{}:
		result := make([]interface{}, 0)
		for i, item := range v {
			if fn(strconv.Itoa(i), &JSON{data: item}) {
				result = append(result, item)
			}
		}
		return &JSON{data: result}

	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, item := range v {
			if fn(k, &JSON{data: item}) {
				result[k] = item
			}
		}
		return &JSON{data: result}
	}

	return j
}

// 序列化方法

// ToJSON 转换为 JSON 字符串
func (j *JSON) ToJSON() (string, error) {
	if j.err != nil {
		return "", j.err
	}

	jsonBytes, err := json.Marshal(j.data)
	if err != nil {
		return "", err
	}

	return bytesToString(jsonBytes), nil
}

// ToPrettyJSON 转换为格式化的 JSON 字符串
func (j *JSON) ToPrettyJSON() (string, error) {
	if j.err != nil {
		return "", j.err
	}

	jsonBytes, err := json.MarshalIndent(j.data, "", "  ")
	if err != nil {
		return "", err
	}

	return bytesToString(jsonBytes), nil
}

// ToBytes 转换为 JSON 字节数组
func (j *JSON) ToBytes() ([]byte, error) {
	if j.err != nil {
		return nil, j.err
	}

	return json.Marshal(j.data)
}

// ToMap 转换为 map
func (j *JSON) ToMap() (map[string]interface{}, error) {
	if j.err != nil {
		return nil, j.err
	}

	if obj, ok := j.data.(map[string]interface{}); ok {
		return obj, nil
	}

	return nil, fmt.Errorf("not an object")
}

// ToSlice 转换为 slice
func (j *JSON) ToSlice() ([]interface{}, error) {
	if j.err != nil {
		return nil, j.err
	}

	if arr, ok := j.data.([]interface{}); ok {
		return arr, nil
	}

	return nil, fmt.Errorf("not an array")
}

// ToInterface 获取原始数据
func (j *JSON) ToInterface() interface{} {
	return j.data
}

// 克隆和合并

// Clone 深度克隆
func (j *JSON) Clone() *JSON {
	if j.err != nil {
		return &JSON{err: j.err}
	}

	cloned := deepClone(j.data)
	return &JSON{data: cloned}
}

// Merge 合并另一个 JSON 对象
func (j *JSON) Merge(other *JSON) *JSON {
	if j.err != nil {
		return j
	}
	if other.err != nil {
		return &JSON{data: j.data, err: other.err}
	}

	thisObj, thisOk := j.data.(map[string]interface{})
	otherObj, otherOk := other.data.(map[string]interface{})

	if !thisOk || !otherOk {
		return &JSON{data: j.data, err: fmt.Errorf("both values must be objects")}
	}

	result := make(map[string]interface{})

	// 复制当前对象的所有字段
	for k, v := range thisObj {
		result[k] = v
	}

	// 覆盖或添加其他对象的字段
	for k, v := range otherObj {
		result[k] = v
	}

	return &JSON{data: result}
}

// DeepMerge 深度合并另一个 JSON 对象
func (j *JSON) DeepMerge(other *JSON) *JSON {
	if j.err != nil {
		return j
	}
	if other.err != nil {
		return &JSON{data: j.data, err: other.err}
	}

	result := deepMerge(j.data, other.data)
	return &JSON{data: result}
}

// 错误处理

// Error 获取错误信息
func (j *JSON) Error() error {
	return j.err
}

// MustGet 获取值，如果出错则 panic
func (j *JSON) MustGet(path string) *JSON {
	result := j.Get(path)
	if result.err != nil {
		panic(result.err)
	}
	return result
}

// MustString 获取字符串，如果出错则 panic
func (j *JSON) MustString() string {
	if j.err != nil {
		panic(j.err)
	}
	return j.String()
}

// MustInt 获取整数，如果出错则 panic
func (j *JSON) MustInt() int {
	if j.err != nil {
		panic(j.err)
	}
	return j.Int()
}

// MustBool 获取布尔值，如果出错则 panic
func (j *JSON) MustBool() bool {
	if j.err != nil {
		panic(j.err)
	}
	return j.Bool()
}

// MustJSON 获取 JSON 字符串，如果出错则 panic
func (j *JSON) MustJSON() string {
	result, err := j.ToJSON()
	if err != nil {
		panic(err)
	}
	return result
}

// 内部方法

// getByPath 根据路径获取值
func (j *JSON) getByPath(path string) (interface{}, error) {
	if path == "" {
		return j.data, nil
	}

	parts := strings.Split(path, ".")
	current := j.data

	for _, part := range parts {
		// 检查是否为数组索引
		if idx, err := strconv.Atoi(part); err == nil {
			if arr, ok := current.([]interface{}); ok {
				if idx < 0 || idx >= len(arr) {
					return nil, fmt.Errorf("array index out of range: %d", idx)
				}
				current = arr[idx]
				continue
			}
		}

		// 处理对象属性
		if obj, ok := current.(map[string]interface{}); ok {
			if value, exists := obj[part]; exists {
				current = value
			} else {
				return nil, fmt.Errorf("path not found: %s", path)
			}
		} else {
			return nil, fmt.Errorf("cannot access property '%s' on non-object", part)
		}
	}

	return current, nil
}

// setByPath 根据路径设置值
func (j *JSON) setByPath(path string, value interface{}) error {
	if path == "" {
		j.data = value
		return nil
	}

	parts := strings.Split(path, ".")

	// 确保根数据结构存在
	if j.data == nil {
		if idx, err := strconv.Atoi(parts[0]); err == nil {
			j.data = make([]interface{}, idx+1)
		} else {
			j.data = make(map[string]interface{})
		}
	}

	return j.setByPathRecursive(j.data, parts, value, "")
}

// setByPathRecursive 递归设置路径值
func (j *JSON) setByPathRecursive(current interface{}, parts []string, value interface{}, currentPath string) error {
	if len(parts) == 0 {
		return fmt.Errorf("empty path parts")
	}

	part := parts[0]
	isLast := len(parts) == 1

	// 如果是最后一个部分，直接设置值
	if isLast {
		if idx, err := strconv.Atoi(part); err == nil {
			// 设置数组元素
			if arr, ok := current.([]interface{}); ok {
				if idx < 0 || idx >= 10000 {
					return fmt.Errorf("invalid array index: %d", idx)
				}
				// 扩展数组
				for len(arr) <= idx {
					arr = append(arr, nil)
				}
				arr[idx] = value
				j.updateDataReference(current, arr, currentPath)
				return nil
			}
			return fmt.Errorf("cannot set array index on non-array")
		} else {
			// 设置对象属性
			if obj, ok := current.(map[string]interface{}); ok {
				obj[part] = value
				return nil
			}
			return fmt.Errorf("cannot set property on non-object")
		}
	}

	// 不是最后一个部分，需要递归
	nextPart := parts[1]
	var nextCurrent interface{}

	if idx, err := strconv.Atoi(part); err == nil {
		// 当前部分是数组索引
		if arr, ok := current.([]interface{}); ok {
			if idx < 0 || idx >= 10000 {
				return fmt.Errorf("invalid array index: %d", idx)
			}
			// 扩展数组
			for len(arr) <= idx {
				arr = append(arr, nil)
			}

			// 如果该位置为 nil，需要创建新的容器
			if arr[idx] == nil {
				if nextIdx, nextErr := strconv.Atoi(nextPart); nextErr == nil {
					arr[idx] = make([]interface{}, nextIdx+1)
				} else {
					arr[idx] = make(map[string]interface{})
				}
			}

			nextCurrent = arr[idx]
			j.updateDataReference(current, arr, currentPath)
		} else {
			return fmt.Errorf("cannot access array index on non-array")
		}
	} else {
		// 当前部分是对象属性
		if obj, ok := current.(map[string]interface{}); ok {
			if _, exists := obj[part]; !exists {
				// 创建新的容器
				if _, nextErr := strconv.Atoi(nextPart); nextErr == nil {
					obj[part] = make([]interface{}, 0)
				} else {
					obj[part] = make(map[string]interface{})
				}
			}
			nextCurrent = obj[part]
		} else {
			return fmt.Errorf("cannot access property on non-object")
		}
	}

	// 递归处理剩余路径
	nextPath := currentPath
	if nextPath != "" {
		nextPath += "."
	}
	nextPath += part

	return j.setByPathRecursive(nextCurrent, parts[1:], value, nextPath)
}

// updateDataReference 更新数据引用
func (j *JSON) updateDataReference(oldRef, newRef interface{}, path string) {
	if path == "" {
		j.data = newRef
		return
	}

	// 更新嵌套引用的逻辑
	parts := strings.Split(path, ".")
	current := j.data

	for i, part := range parts {
		if i == len(parts)-1 {
			// 最后一个部分，更新引用
			if obj, ok := current.(map[string]interface{}); ok {
				obj[part] = newRef
			}
			break
		}

		// 继续深入
		if idx, err := strconv.Atoi(part); err == nil {
			if arr, ok := current.([]interface{}); ok && idx >= 0 && idx < len(arr) {
				current = arr[idx]
			}
		} else {
			if obj, ok := current.(map[string]interface{}); ok {
				current = obj[part]
			}
		}
	}
}

// deleteByPath 根据路径删除值
func (j *JSON) deleteByPath(path string) error {
	if path == "" {
		j.data = nil
		return nil
	}

	parts := strings.Split(path, ".")
	current := j.data

	// 遍历到最后一个部分之前
	for _, part := range parts[:len(parts)-1] {
		if idx, err := strconv.Atoi(part); err == nil {
			if arr, ok := current.([]interface{}); ok {
				if idx < 0 || idx >= len(arr) {
					return fmt.Errorf("array index out of range: %d", idx)
				}
				current = arr[idx]
				continue
			}
		}

		if obj, ok := current.(map[string]interface{}); ok {
			if value, exists := obj[part]; exists {
				current = value
			} else {
				return fmt.Errorf("path not found: %s", path)
			}
		} else {
			return fmt.Errorf("cannot access property '%s' on non-object", part)
		}
	}

	// 删除最后一个部分
	lastPart := parts[len(parts)-1]
	if obj, ok := current.(map[string]interface{}); ok {
		delete(obj, lastPart)
		return nil
	}

	return fmt.Errorf("cannot delete from non-object")
}

// updateArrayReference 更新数组引用
func (j *JSON) updateArrayReference(path []string, newArray []interface{}) {
	if len(path) == 0 {
		j.data = newArray
		return
	}

	current := j.data
	for _, part := range path[:len(path)-1] {
		if obj, ok := current.(map[string]interface{}); ok {
			current = obj[part]
		}
	}

	if obj, ok := current.(map[string]interface{}); ok {
		obj[path[len(path)-1]] = newArray
	}
}

// updateArrayAtPath 更新指定路径的数组
func (j *JSON) updateArrayAtPath(path []string, newArray []interface{}) {
	if len(path) == 0 {
		j.data = newArray
		return
	}

	current := j.data
	for _, part := range path[:len(path)-1] {
		if obj, ok := current.(map[string]interface{}); ok {
			current = obj[part]
		}
	}

	if obj, ok := current.(map[string]interface{}); ok {
		obj[path[len(path)-1]] = newArray
	}
}

// 工具函数

// structToMap 将结构体转换为 map
func structToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// deepClone 深度克隆
func deepClone(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, v := range val {
			result[k] = deepClone(v)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, v := range val {
			result[i] = deepClone(v)
		}
		return result
	default:
		return v
	}
}

// deepMerge 深度合并
func deepMerge(dst, src interface{}) interface{} {
	dstMap, dstOk := dst.(map[string]interface{})
	srcMap, srcOk := src.(map[string]interface{})

	if dstOk && srcOk {
		result := make(map[string]interface{})

		// 复制目标对象
		for k, v := range dstMap {
			result[k] = v
		}

		// 深度合并源对象
		for k, v := range srcMap {
			if dstVal, exists := result[k]; exists {
				result[k] = deepMerge(dstVal, v)
			} else {
				result[k] = deepClone(v)
			}
		}

		return result
	}

	return deepClone(src)
}

// bytesToString 高效的字节数组转字符串
func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// stringToBytes 高效的字符串转字节数组
func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
