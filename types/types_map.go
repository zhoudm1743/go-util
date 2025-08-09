package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type XMap[K comparable, V any] map[K]V

func Map[K comparable, V any](m map[K]V) XMap[K, V] {
	return XMap[K, V](m)
}

func NewMap[K comparable, V any]() XMap[K, V] {
	return make(XMap[K, V])
}

// Len 返回 Map 的长度
func (m XMap[K, V]) Len() int {
	return len(m)
}

// IsEmpty 判断 Map 是否为空
func (m XMap[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

// Set 设置键值对
func (m XMap[K, V]) Set(key K, value V) XMap[K, V] {
	m[key] = value
	return m
}

// Get 获取值，如果不存在返回零值和 false
func (m XMap[K, V]) Get(key K) (V, bool) {
	value, exists := m[key]
	return value, exists
}

// GetOrDefault 获取值，如果不存在返回默认值
func (m XMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}

// Has 判断是否包含指定键
func (m XMap[K, V]) Has(key K) bool {
	_, exists := m[key]
	return exists
}

// Delete 删除指定键
func (m XMap[K, V]) Delete(key K) XMap[K, V] {
	delete(m, key)
	return m
}

// Clear 清空 Map
func (m XMap[K, V]) Clear() XMap[K, V] {
	for k := range m {
		delete(m, k)
	}
	return m
}

// Keys 获取所有键
func (m XMap[K, V]) Keys() []K {
	keys := make([]K, 0, m.Len())
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values 获取所有值
func (m XMap[K, V]) Values() []V {
	values := make([]V, 0, m.Len())
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Merge 合并另一个 Map
func (m XMap[K, V]) Merge(other XMap[K, V]) XMap[K, V] {
	for k, v := range other {
		m[k] = v
	}
	return m
}

// Copy 创建副本
func (m XMap[K, V]) Copy() XMap[K, V] {
	newMap := make(XMap[K, V], m.Len())
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

// Filter 过滤元素
func (m XMap[K, V]) Filter(predicate func(K, V) bool) XMap[K, V] {
	result := NewMap[K, V]()
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// ForEach 遍历执行函数
func (m XMap[K, V]) ForEach(fn func(K, V)) {
	for k, v := range m {
		fn(k, v)
	}
}

// Equal 判断两个 Map 是否相等
func (m XMap[K, V]) Equal(other XMap[K, V]) bool {
	if m.Len() != other.Len() {
		return false
	}
	for k, v := range m {
		if otherV, exists := other[k]; !exists || !reflect.DeepEqual(v, otherV) {
			return false
		}
	}
	return true
}

// ToJSON 转换为 JSON 字符串
func (m XMap[K, V]) ToJSON() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON 从 JSON 字符串创建 Map
func (m XMap[K, V]) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), &m)
}

// String 实现 Stringer 接口
func (m XMap[K, V]) String() string {
	if m.IsEmpty() {
		return "map[]"
	}

	pairs := make([]string, 0, m.Len())
	for k, v := range m {
		pairs = append(pairs, fmt.Sprintf("%v:%v", k, v))
	}
	return fmt.Sprintf("map[%s]", strings.Join(pairs, " "))
}

// SortedKeys 返回排序后的键（仅适用于可排序的键类型）
func (m XMap[K, V]) SortedKeys() []K {
	keys := m.Keys()

	// 使用类型断言来处理排序
	switch any(keys).(type) {
	case []string:
		stringKeys := any(keys).([]string)
		sort.Strings(stringKeys)
		return any(stringKeys).([]K)
	case []int:
		intKeys := any(keys).([]int)
		sort.Ints(intKeys)
		return any(intKeys).([]K)
	case []float64:
		floatKeys := any(keys).([]float64)
		sort.Float64s(floatKeys)
		return any(floatKeys).([]K)
	default:
		// 对于其他类型，尝试使用泛型约束
		if isOrdered(keys) {
			sortOrdered(keys)
		}
		return keys
	}
}

// 辅助函数：检查类型是否可排序
func isOrdered[K any](keys []K) bool {
	if len(keys) == 0 {
		return false
	}
	// 检查是否为基本可排序类型
	switch any(keys[0]).(type) {
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

// 辅助函数：对可排序类型进行排序
func sortOrdered[K any](keys []K) {
	if len(keys) <= 1 {
		return
	}

	// 简单的冒泡排序实现
	for i := 0; i < len(keys)-1; i++ {
		for j := 0; j < len(keys)-i-1; j++ {
			if shouldSwap(keys[j], keys[j+1]) {
				keys[j], keys[j+1] = keys[j+1], keys[j]
			}
		}
	}
}

// 辅助函数：比较两个值是否需要交换
func shouldSwap[K any](a, b K) bool {
	switch va := any(a).(type) {
	case string:
		return va > any(b).(string)
	case int:
		return va > any(b).(int)
	case int8:
		return va > any(b).(int8)
	case int16:
		return va > any(b).(int16)
	case int32:
		return va > any(b).(int32)
	case int64:
		return va > any(b).(int64)
	case uint:
		return va > any(b).(uint)
	case uint8:
		return va > any(b).(uint8)
	case uint16:
		return va > any(b).(uint16)
	case uint32:
		return va > any(b).(uint32)
	case uint64:
		return va > any(b).(uint64)
	case float32:
		return va > any(b).(float32)
	case float64:
		return va > any(b).(float64)
	default:
		return false
	}
}
