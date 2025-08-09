package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type XJSON struct {
	data interface{}
}

// JSON 创建 XJSON 实例
func JSON(data interface{}) XJSON {
	return XJSON{data: data}
}

// ParseJSON 解析 JSON 字符串
func ParseJSON(jsonStr string) (XJSON, error) {
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return XJSON{}, err
	}
	return XJSON{data: data}, nil
}

// NewJSONObject 创建空的 JSON 对象
func NewJSONObject() XJSON {
	return XJSON{data: make(map[string]interface{})}
}

// NewJSONArray 创建空的 JSON 数组
func NewJSONArray() XJSON {
	return XJSON{data: make([]interface{}, 0)}
}

// String 转换为 JSON 字符串
func (j XJSON) String() string {
	result, _ := j.StringWithIndent("")
	return result
}

// StringWithIndent 转换为带缩进的 JSON 字符串
func (j XJSON) StringWithIndent(indent string) (string, error) {
	if indent == "" {
		data, err := json.Marshal(j.data)
		return string(data), err
	}
	data, err := json.MarshalIndent(j.data, "", indent)
	return string(data), err
}

// Pretty 返回格式化的 JSON 字符串
func (j XJSON) Pretty() string {
	result, _ := j.StringWithIndent("  ")
	return result
}

// Data 获取原始数据
func (j XJSON) Data() interface{} {
	return j.data
}

// IsNull 判断是否为 null
func (j XJSON) IsNull() bool {
	return j.data == nil
}

// IsObject 判断是否为对象
func (j XJSON) IsObject() bool {
	_, ok := j.data.(map[string]interface{})
	return ok
}

// IsArray 判断是否为数组
func (j XJSON) IsArray() bool {
	_, ok := j.data.([]interface{})
	return ok
}

// IsString 判断是否为字符串
func (j XJSON) IsString() bool {
	_, ok := j.data.(string)
	return ok
}

// IsNumber 判断是否为数字
func (j XJSON) IsNumber() bool {
	switch j.data.(type) {
	case float64, int, int64, float32:
		return true
	default:
		return false
	}
}

// IsBool 判断是否为布尔值
func (j XJSON) IsBool() bool {
	_, ok := j.data.(bool)
	return ok
}

// Get 获取指定路径的值
func (j XJSON) Get(path string) XJSON {
	return j.GetWithSeparator(path, ".")
}

// GetWithSeparator 使用指定分隔符获取路径值
func (j XJSON) GetWithSeparator(path, separator string) XJSON {
	if path == "" {
		return j
	}

	keys := strings.Split(path, separator)
	current := j.data

	for _, key := range keys {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[key]
		case []interface{}:
			if index, err := strconv.Atoi(key); err == nil && index >= 0 && index < len(v) {
				current = v[index]
			} else {
				return XJSON{data: nil}
			}
		default:
			return XJSON{data: nil}
		}

		if current == nil {
			break
		}
	}

	return XJSON{data: current}
}

// Set 设置指定路径的值
func (j XJSON) Set(path string, value interface{}) XJSON {
	return j.SetWithSeparator(path, value, ".")
}

// SetWithSeparator 使用指定分隔符设置路径值
func (j XJSON) SetWithSeparator(path string, value interface{}, separator string) XJSON {
	if path == "" {
		return XJSON{data: value}
	}

	keys := strings.Split(path, separator)
	j.ensureObject()

	current := j.data.(map[string]interface{})

	for i, key := range keys[:len(keys)-1] {
		next, exists := current[key]
		if !exists || next == nil {
			// 判断下一个键是否为数字（数组索引）
			nextKey := keys[i+1]
			if _, err := strconv.Atoi(nextKey); err == nil {
				next = make([]interface{}, 0)
			} else {
				next = make(map[string]interface{})
			}
			current[key] = next
		}

		switch v := next.(type) {
		case map[string]interface{}:
			current = v
		case []interface{}:
			// 处理数组情况
			if index, err := strconv.Atoi(keys[i+1]); err == nil {
				// 扩展数组大小
				for len(v) <= index {
					v = append(v, nil)
				}
				current[key] = v

				if index < len(v) {
					if v[index] == nil {
						v[index] = make(map[string]interface{})
					}
					if nextMap, ok := v[index].(map[string]interface{}); ok {
						current = nextMap
					}
				}
			}
		default:
			// 创建新的对象
			newObj := make(map[string]interface{})
			current[key] = newObj
			current = newObj
		}
	}

	lastKey := keys[len(keys)-1]
	current[lastKey] = value

	return j
}

// Delete 删除指定路径的值
func (j XJSON) Delete(path string) XJSON {
	return j.DeleteWithSeparator(path, ".")
}

// DeleteWithSeparator 使用指定分隔符删除路径值
func (j XJSON) DeleteWithSeparator(path, separator string) XJSON {
	if path == "" {
		return XJSON{data: nil}
	}

	keys := strings.Split(path, separator)
	if len(keys) == 0 {
		return j
	}

	current := j.data
	parents := make([]interface{}, len(keys)-1)

	// 遍历到父级
	for i, key := range keys[:len(keys)-1] {
		parents[i] = current
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[key]
		case []interface{}:
			if index, err := strconv.Atoi(key); err == nil && index >= 0 && index < len(v) {
				current = v[index]
			} else {
				return j
			}
		default:
			return j
		}
	}

	// 删除最后一个键
	lastKey := keys[len(keys)-1]
	switch v := current.(type) {
	case map[string]interface{}:
		delete(v, lastKey)
	case []interface{}:
		if index, err := strconv.Atoi(lastKey); err == nil && index >= 0 && index < len(v) {
			// 从数组中删除元素
			copy(v[index:], v[index+1:])
			v = v[:len(v)-1]

			// 更新父级的引用
			if len(parents) > 0 {
				parentKey := keys[len(keys)-2]
				switch parent := parents[len(parents)-1].(type) {
				case map[string]interface{}:
					parent[parentKey] = v
				}
			}
		}
	}

	return j
}

// Has 判断是否存在指定路径
func (j XJSON) Has(path string) bool {
	return !j.Get(path).IsNull()
}

// Keys 获取对象的所有键
func (j XJSON) Keys() []string {
	if obj, ok := j.data.(map[string]interface{}); ok {
		keys := make([]string, 0, len(obj))
		for k := range obj {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

// Values 获取对象的所有值
func (j XJSON) Values() []XJSON {
	if obj, ok := j.data.(map[string]interface{}); ok {
		values := make([]XJSON, 0, len(obj))
		for _, v := range obj {
			values = append(values, XJSON{data: v})
		}
		return values
	}
	return nil
}

// Len 获取数组长度或对象键数量
func (j XJSON) Len() int {
	switch v := j.data.(type) {
	case []interface{}:
		return len(v)
	case map[string]interface{}:
		return len(v)
	default:
		return 0
	}
}

// Index 获取数组指定索引的元素
func (j XJSON) Index(index int) XJSON {
	if arr, ok := j.data.([]interface{}); ok && index >= 0 && index < len(arr) {
		return XJSON{data: arr[index]}
	}
	return XJSON{data: nil}
}

// Append 向数组追加元素
func (j XJSON) Append(values ...interface{}) XJSON {
	j.ensureArray()
	if arr, ok := j.data.([]interface{}); ok {
		arr = append(arr, values...)
		j.data = arr
	}
	return j
}

// Prepend 向数组前面插入元素
func (j XJSON) Prepend(values ...interface{}) XJSON {
	j.ensureArray()
	if arr, ok := j.data.([]interface{}); ok {
		newArr := make([]interface{}, len(values)+len(arr))
		copy(newArr, values)
		copy(newArr[len(values):], arr)
		j.data = newArr
	}
	return j
}

// Insert 在指定位置插入元素
func (j XJSON) Insert(index int, values ...interface{}) XJSON {
	j.ensureArray()
	if arr, ok := j.data.([]interface{}); ok {
		if index < 0 {
			index = 0
		}
		if index > len(arr) {
			index = len(arr)
		}

		newArr := make([]interface{}, len(arr)+len(values))
		copy(newArr, arr[:index])
		copy(newArr[index:], values)
		copy(newArr[index+len(values):], arr[index:])
		j.data = newArr
	}
	return j
}

// Remove 删除数组指定索引的元素
func (j XJSON) Remove(index int) XJSON {
	if arr, ok := j.data.([]interface{}); ok && index >= 0 && index < len(arr) {
		newArr := make([]interface{}, len(arr)-1)
		copy(newArr, arr[:index])
		copy(newArr[index:], arr[index+1:])
		j.data = newArr
	}
	return j
}

// ForEach 遍历数组或对象
func (j XJSON) ForEach(fn func(key interface{}, value XJSON)) {
	switch v := j.data.(type) {
	case []interface{}:
		for i, item := range v {
			fn(i, XJSON{data: item})
		}
	case map[string]interface{}:
		for key, value := range v {
			fn(key, XJSON{data: value})
		}
	}
}

// Map 映射数组元素
func (j XJSON) Map(fn func(XJSON) XJSON) XJSON {
	if arr, ok := j.data.([]interface{}); ok {
		newArr := make([]interface{}, len(arr))
		for i, item := range arr {
			newItem := fn(XJSON{data: item})
			newArr[i] = newItem.data
		}
		return XJSON{data: newArr}
	}
	return j
}

// Filter 过滤数组元素
func (j XJSON) Filter(fn func(XJSON) bool) XJSON {
	if arr, ok := j.data.([]interface{}); ok {
		newArr := make([]interface{}, 0)
		for _, item := range arr {
			if fn(XJSON{data: item}) {
				newArr = append(newArr, item)
			}
		}
		return XJSON{data: newArr}
	}
	return j
}

// AsString 转换为字符串
func (j XJSON) AsString() string {
	switch v := j.data.(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// AsInt 转换为整数
func (j XJSON) AsInt() int {
	switch v := j.data.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// AsInt64 转换为 int64
func (j XJSON) AsInt64() int64 {
	switch v := j.data.(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

// AsFloat 转换为浮点数
func (j XJSON) AsFloat() float64 {
	switch v := j.data.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0.0
}

// AsBool 转换为布尔值
func (j XJSON) AsBool() bool {
	switch v := j.data.(type) {
	case bool:
		return v
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
		return v != ""
	case float64:
		return v != 0
	case int:
		return v != 0
	}
	return false
}

// AsArray 转换为 XJSON 数组
func (j XJSON) AsArray() []XJSON {
	if arr, ok := j.data.([]interface{}); ok {
		result := make([]XJSON, len(arr))
		for i, item := range arr {
			result[i] = XJSON{data: item}
		}
		return result
	}
	return nil
}

// AsMap 转换为 map[string]XJSON
func (j XJSON) AsMap() map[string]XJSON {
	if obj, ok := j.data.(map[string]interface{}); ok {
		result := make(map[string]XJSON)
		for k, v := range obj {
			result[k] = XJSON{data: v}
		}
		return result
	}
	return nil
}

// Clone 深度复制
func (j XJSON) Clone() XJSON {
	return XJSON{data: j.deepCopy(j.data)}
}

// Merge 合并对象
func (j XJSON) Merge(other XJSON) XJSON {
	j.ensureObject()
	other.ensureObject()

	if obj1, ok1 := j.data.(map[string]interface{}); ok1 {
		if obj2, ok2 := other.data.(map[string]interface{}); ok2 {
			for k, v := range obj2 {
				obj1[k] = v
			}
		}
	}
	return j
}

// DeepMerge 深度合并对象
func (j XJSON) DeepMerge(other XJSON) XJSON {
	j.ensureObject()
	other.ensureObject()

	if obj1, ok1 := j.data.(map[string]interface{}); ok1 {
		if obj2, ok2 := other.data.(map[string]interface{}); ok2 {
			j.data = j.deepMerge(obj1, obj2)
		}
	}
	return j
}

// 辅助方法：确保数据是对象
func (j *XJSON) ensureObject() {
	if !j.IsObject() {
		j.data = make(map[string]interface{})
	}
}

// 辅助方法：确保数据是数组
func (j *XJSON) ensureArray() {
	if !j.IsArray() {
		j.data = make([]interface{}, 0)
	}
}

// 辅助方法：深度复制
func (j XJSON) deepCopy(src interface{}) interface{} {
	switch v := src.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for k, val := range v {
			newMap[k] = j.deepCopy(val)
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, len(v))
		for i, val := range v {
			newSlice[i] = j.deepCopy(val)
		}
		return newSlice
	default:
		return v
	}
}

// 辅助方法：深度合并
func (j XJSON) deepMerge(obj1, obj2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// 复制第一个对象
	for k, v := range obj1 {
		result[k] = j.deepCopy(v)
	}

	// 合并第二个对象
	for k, v := range obj2 {
		if existing, exists := result[k]; exists {
			if map1, ok1 := existing.(map[string]interface{}); ok1 {
				if map2, ok2 := v.(map[string]interface{}); ok2 {
					result[k] = j.deepMerge(map1, map2)
					continue
				}
			}
		}
		result[k] = j.deepCopy(v)
	}

	return result
}
