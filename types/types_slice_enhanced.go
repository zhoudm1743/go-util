package types

import (
	"math"
	"sort"
	"strings"
)

// Sort 排序（升序）
func (a XArray[T]) Sort() XArray[T] {
	sorted := make(XArray[T], len(a))
	copy(sorted, a)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	return sorted
}

// SortDesc 排序（降序）
func (a XArray[T]) SortDesc() XArray[T] {
	sorted := make(XArray[T], len(a))
	copy(sorted, a)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] > sorted[j]
	})
	return sorted
}

// SortBy 自定义排序
func (a XArray[T]) SortBy(less func(T, T) bool) XArray[T] {
	sorted := make(XArray[T], len(a))
	copy(sorted, a)
	sort.Slice(sorted, func(i, j int) bool {
		return less(sorted[i], sorted[j])
	})
	return sorted
}

// Reverse 反转数组
func (a XArray[T]) Reverse() XArray[T] {
	reversed := make(XArray[T], len(a))
	for i, v := range a {
		reversed[len(a)-1-i] = v
	}
	return reversed
}

// Shuffle 随机打乱数组
func (a XArray[T]) Shuffle() XArray[T] {
	shuffled := make(XArray[T], len(a))
	copy(shuffled, a)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := int(math.Floor(math.Mod(float64(i+1), float64(len(shuffled)))))
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}

// Take 取前 n 个元素
func (a XArray[T]) Take(n int) XArray[T] {
	if n <= 0 {
		return XArray[T]{}
	}
	if n >= len(a) {
		return a
	}
	return a[:n]
}

// Skip 跳过前 n 个元素
func (a XArray[T]) Skip(n int) XArray[T] {
	if n <= 0 {
		return a
	}
	if n >= len(a) {
		return XArray[T]{}
	}
	return a[n:]
}

// TakeWhile 取满足条件的前部分元素
func (a XArray[T]) TakeWhile(predicate func(T) bool) XArray[T] {
	for i, v := range a {
		if !predicate(v) {
			return a[:i]
		}
	}
	return a
}

// SkipWhile 跳过满足条件的前部分元素
func (a XArray[T]) SkipWhile(predicate func(T) bool) XArray[T] {
	for i, v := range a {
		if !predicate(v) {
			return a[i:]
		}
	}
	return XArray[T]{}
}

// Filter 过滤元素
func (a XArray[T]) Filter(predicate func(T) bool) XArray[T] {
	var result XArray[T]
	for _, v := range a {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map 映射转换
func (a XArray[T]) Map(mapper func(T) T) XArray[T] {
	result := make(XArray[T], len(a))
	for i, v := range a {
		result[i] = mapper(v)
	}
	return result
}

// Reduce 归约操作
func (a XArray[T]) Reduce(initial T, reducer func(T, T) T) T {
	result := initial
	for _, v := range a {
		result = reducer(result, v)
	}
	return result
}

// All 判断是否所有元素都满足条件
func (a XArray[T]) All(predicate func(T) bool) bool {
	for _, v := range a {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Any 判断是否有任一元素满足条件
func (a XArray[T]) Any(predicate func(T) bool) bool {
	for _, v := range a {
		if predicate(v) {
			return true
		}
	}
	return false
}

// None 判断是否没有元素满足条件
func (a XArray[T]) None(predicate func(T) bool) bool {
	return !a.Any(predicate)
}

// First 获取第一个元素
func (a XArray[T]) First() (T, bool) {
	if len(a) == 0 {
		var zero T
		return zero, false
	}
	return a[0], true
}

// Last 获取最后一个元素
func (a XArray[T]) Last() (T, bool) {
	if len(a) == 0 {
		var zero T
		return zero, false
	}
	return a[len(a)-1], true
}

// FirstOrDefault 获取第一个元素或默认值
func (a XArray[T]) FirstOrDefault(defaultValue T) T {
	if len(a) == 0 {
		return defaultValue
	}
	return a[0]
}

// LastOrDefault 获取最后一个元素或默认值
func (a XArray[T]) LastOrDefault(defaultValue T) T {
	if len(a) == 0 {
		return defaultValue
	}
	return a[len(a)-1]
}

// Find 查找第一个满足条件的元素
func (a XArray[T]) Find(predicate func(T) bool) (T, bool) {
	for _, v := range a {
		if predicate(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex 查找第一个满足条件的元素的索引
func (a XArray[T]) FindIndex(predicate func(T) bool) int {
	for i, v := range a {
		if predicate(v) {
			return i
		}
	}
	return -1
}

// FindLast 查找最后一个满足条件的元素
func (a XArray[T]) FindLast(predicate func(T) bool) (T, bool) {
	for i := len(a) - 1; i >= 0; i-- {
		if predicate(a[i]) {
			return a[i], true
		}
	}
	var zero T
	return zero, false
}

// FindLastIndex 查找最后一个满足条件的元素的索引
func (a XArray[T]) FindLastIndex(predicate func(T) bool) int {
	for i := len(a) - 1; i >= 0; i-- {
		if predicate(a[i]) {
			return i
		}
	}
	return -1
}

// IndexOf 查找元素的索引
func (a XArray[T]) IndexOf(item T) int {
	return a.Index(item)
}

// LastIndexOf 查找元素最后出现的索引
func (a XArray[T]) LastIndexOf(item T) int {
	for i := len(a) - 1; i >= 0; i-- {
		if a.EqualItem(a[i], item) {
			return i
		}
	}
	return -1
}

// Contains 判断是否包含指定元素
func (a XArray[T]) Contains(item T) bool {
	return a.Exist(item)
}

// Count 统计满足条件的元素数量
func (a XArray[T]) Count(predicate func(T) bool) int {
	count := 0
	for _, v := range a {
		if predicate(v) {
			count++
		}
	}
	return count
}

// Chunk 分块
func (a XArray[T]) Chunk(size int) []XArray[T] {
	if size <= 0 {
		return nil
	}

	var chunks []XArray[T]
	for i := 0; i < len(a); i += size {
		end := i + size
		if end > len(a) {
			end = len(a)
		}
		chunks = append(chunks, a[i:end])
	}
	return chunks
}

// GroupBy 分组
func (a XArray[T]) GroupBy(keySelector func(T) string) map[string]XArray[T] {
	groups := make(map[string]XArray[T])
	for _, item := range a {
		key := keySelector(item)
		groups[key] = append(groups[key], item)
	}
	return groups
}

// Distinct 去重
func (a XArray[T]) Distinct() XArray[T] {
	return a.UniqueOrdered()
}

// DistinctBy 按指定键去重
func (a XArray[T]) DistinctBy(keySelector func(T) string) XArray[T] {
	seen := make(map[string]bool)
	var result XArray[T]

	for _, item := range a {
		key := keySelector(item)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}
	return result
}

// Intersect 交集
func (a XArray[T]) Intersect(other XArray[T]) XArray[T] {
	var result XArray[T]
	for _, item := range a {
		if other.Contains(item) && !result.Contains(item) {
			result = append(result, item)
		}
	}
	return result
}

// Union 并集
func (a XArray[T]) Union(other XArray[T]) XArray[T] {
	result := a.Distinct()
	for _, item := range other {
		if !result.Contains(item) {
			result = append(result, item)
		}
	}
	return result
}

// Except 差集
func (a XArray[T]) Except(other XArray[T]) XArray[T] {
	var result XArray[T]
	for _, item := range a {
		if !other.Contains(item) {
			result = append(result, item)
		}
	}
	return result
}

// ZipPairs 压缩两个数组为元组切片
func (a XArray[T]) ZipPairs(other XArray[T]) [][2]T {
	minLen := len(a)
	if len(other) < minLen {
		minLen = len(other)
	}

	result := make([][2]T, minLen)
	for i := 0; i < minLen; i++ {
		result[i] = [2]T{a[i], other[i]}
	}
	return result
}

// Partition 分区
func (a XArray[T]) Partition(predicate func(T) bool) (XArray[T], XArray[T]) {
	var trueItems, falseItems XArray[T]
	for _, item := range a {
		if predicate(item) {
			trueItems = append(trueItems, item)
		} else {
			falseItems = append(falseItems, item)
		}
	}
	return trueItems, falseItems
}

// ForEach 遍历执行操作
func (a XArray[T]) ForEach(action func(T)) {
	for _, item := range a {
		action(item)
	}
}

// ForEachIndexed 带索引遍历执行操作
func (a XArray[T]) ForEachIndexed(action func(int, T)) {
	for i, item := range a {
		action(i, item)
	}
}

// 数字类型特有方法

// Sum 求和（仅适用于数字类型）
func (a XArray[T]) Sum() T {
	var sum T
	for _, v := range a {
		sum = sum + v
	}
	return sum
}

// Average 求平均值（仅适用于数字类型）
func (a XArray[T]) Average() float64 {
	if len(a) == 0 {
		return 0
	}

	var sum T
	for _, v := range a {
		sum = sum + v
	}

	// 类型转换为 float64
	switch any(sum).(type) {
	case int:
		return float64(any(sum).(int)) / float64(len(a))
	case int32:
		return float64(any(sum).(int32)) / float64(len(a))
	case int64:
		return float64(any(sum).(int64)) / float64(len(a))
	case float32:
		return float64(any(sum).(float32)) / float64(len(a))
	case float64:
		return any(sum).(float64) / float64(len(a))
	default:
		return 0
	}
}

// Min 求最小值
func (a XArray[T]) Min() (T, bool) {
	if len(a) == 0 {
		var zero T
		return zero, false
	}

	min := a[0]
	for _, v := range a[1:] {
		if v < min {
			min = v
		}
	}
	return min, true
}

// Max 求最大值
func (a XArray[T]) Max() (T, bool) {
	if len(a) == 0 {
		var zero T
		return zero, false
	}

	max := a[0]
	for _, v := range a[1:] {
		if v > max {
			max = v
		}
	}
	return max, true
}

// MinBy 按指定函数求最小值
func (a XArray[T]) MinBy(selector func(T) T) (T, bool) {
	if len(a) == 0 {
		var zero T
		return zero, false
	}

	minItem := a[0]
	minValue := selector(minItem)

	for _, item := range a[1:] {
		value := selector(item)
		if value < minValue {
			minItem = item
			minValue = value
		}
	}
	return minItem, true
}

// MaxBy 按指定函数求最大值
func (a XArray[T]) MaxBy(selector func(T) T) (T, bool) {
	if len(a) == 0 {
		var zero T
		return zero, false
	}

	maxItem := a[0]
	maxValue := selector(maxItem)

	for _, item := range a[1:] {
		value := selector(item)
		if value > maxValue {
			maxItem = item
			maxValue = value
		}
	}
	return maxItem, true
}

// ToSlice 转换为原生切片
func (a XArray[T]) ToSlice() []T {
	return []T(a)
}

// Clone 深度复制
func (a XArray[T]) Clone() XArray[T] {
	clone := make(XArray[T], len(a))
	copy(clone, a)
	return clone
}

// IsEmpty 判断是否为空
func (a XArray[T]) IsEmpty() bool {
	return len(a) == 0
}

// IsNotEmpty 判断是否不为空
func (a XArray[T]) IsNotEmpty() bool {
	return len(a) > 0
}

// 字符串数组特有方法

// JoinStrings 连接字符串数组（仅适用于字符串类型）
func JoinStrings[T ~string](a XArray[T], separator string) string {
	if len(a) == 0 {
		return ""
	}

	strs := make([]string, len(a))
	for i, v := range a {
		strs[i] = string(v)
	}
	return strings.Join(strs, separator)
}

// FilterEmpty 过滤空字符串（仅适用于字符串类型）
func FilterEmpty[T ~string](a XArray[T]) XArray[T] {
	return a.Filter(func(s T) bool {
		return string(s) != ""
	})
}

// TrimAll 去除所有字符串的空白字符（仅适用于字符串类型）
func TrimAll[T ~string](a XArray[T]) XArray[T] {
	return a.Map(func(s T) T {
		return T(strings.TrimSpace(string(s)))
	})
}

// ToLowerAll 转换所有字符串为小写（仅适用于字符串类型）
func ToLowerAll[T ~string](a XArray[T]) XArray[T] {
	return a.Map(func(s T) T {
		return T(strings.ToLower(string(s)))
	})
}

// ToUpperAll 转换所有字符串为大写（仅适用于字符串类型）
func ToUpperAll[T ~string](a XArray[T]) XArray[T] {
	return a.Map(func(s T) T {
		return T(strings.ToUpper(string(s)))
	})
}
