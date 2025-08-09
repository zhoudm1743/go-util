package types

import (
	"encoding/json"
	"fmt"
	"sync"

	"golang.org/x/exp/constraints"
)

// XEnum 泛型枚举类型
type XEnum[T constraints.Ordered] struct {
	value T
	name  string
	desc  string
}

// EnumRegistry 枚举注册表，用于管理枚举定义
type EnumRegistry[T constraints.Ordered] struct {
	values map[T]*XEnum[T]
	names  map[string]*XEnum[T]
	mu     sync.RWMutex
	all    []*XEnum[T]
}

// NewEnumRegistry 创建新的枚举注册表
func NewEnumRegistry[T constraints.Ordered]() *EnumRegistry[T] {
	return &EnumRegistry[T]{
		values: make(map[T]*XEnum[T]),
		names:  make(map[string]*XEnum[T]),
		all:    make([]*XEnum[T], 0),
	}
}

// Define 定义枚举值
func (r *EnumRegistry[T]) Define(value T, name, desc string) *XEnum[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	enum := &XEnum[T]{
		value: value,
		name:  name,
		desc:  desc,
	}

	r.values[value] = enum
	r.names[name] = enum
	r.all = append(r.all, enum)

	return enum
}

// DefineSimple 简单定义枚举值（名称和描述相同）
func (r *EnumRegistry[T]) DefineSimple(value T, name string) *XEnum[T] {
	return r.Define(value, name, name)
}

// FromValue 根据值获取枚举
func (r *EnumRegistry[T]) FromValue(value T) (*XEnum[T], bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	enum, exists := r.values[value]
	return enum, exists
}

// FromName 根据名称获取枚举
func (r *EnumRegistry[T]) FromName(name string) (*XEnum[T], bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	enum, exists := r.names[name]
	return enum, exists
}

// All 获取所有枚举值
func (r *EnumRegistry[T]) All() []*XEnum[T] {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*XEnum[T], len(r.all))
	copy(result, r.all)
	return result
}

// Values 获取所有枚举值的原始值
func (r *EnumRegistry[T]) Values() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()
	values := make([]T, 0, len(r.all))
	for _, enum := range r.all {
		values = append(values, enum.value)
	}
	return values
}

// Names 获取所有枚举名称
func (r *EnumRegistry[T]) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.all))
	for _, enum := range r.all {
		names = append(names, enum.name)
	}
	return names
}

// IsValid 检查值是否为有效的枚举值
func (r *EnumRegistry[T]) IsValid(value T) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.values[value]
	return exists
}

// IsValidName 检查名称是否为有效的枚举名称
func (r *EnumRegistry[T]) IsValidName(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.names[name]
	return exists
}

// Count 获取枚举项数量
func (r *EnumRegistry[T]) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.all)
}

// XEnum 方法

// Value 获取枚举值
func (e *XEnum[T]) Value() T {
	return e.value
}

// Name 获取枚举名称
func (e *XEnum[T]) Name() string {
	return e.name
}

// Desc 获取枚举描述
func (e *XEnum[T]) Desc() string {
	return e.desc
}

// String 实现 Stringer 接口
func (e *XEnum[T]) String() string {
	return e.name
}

// Equal 比较两个枚举是否相等
func (e *XEnum[T]) Equal(other *XEnum[T]) bool {
	if other == nil {
		return false
	}
	return e.value == other.value
}

// EqualValue 比较枚举值是否相等
func (e *XEnum[T]) EqualValue(value T) bool {
	return e.value == value
}

// EqualName 比较枚举名称是否相等
func (e *XEnum[T]) EqualName(name string) bool {
	return e.name == name
}

// MarshalJSON 实现 JSON 序列化
func (e *XEnum[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.name)
}

// UnmarshalJSON 实现 JSON 反序列化
func (e *XEnum[T]) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	e.name = name
	return nil
}

// GoString 实现 GoStringer 接口
func (e *XEnum[T]) GoString() string {
	return fmt.Sprintf("Enum{value: %v, name: %q, desc: %q}", e.value, e.name, e.desc)
}

// 预定义的常用枚举类型

// IntEnum 整数枚举
type IntEnum = *EnumRegistry[int]

// StringEnum 字符串枚举
type StringEnum = *EnumRegistry[string]

// 便捷函数

// NewIntEnum 创建整数枚举
func NewIntEnum() IntEnum {
	return NewEnumRegistry[int]()
}

// NewStringEnum 创建字符串枚举
func NewStringEnum() StringEnum {
	return NewEnumRegistry[string]()
}

// 枚举构建器 - 更流畅的API

// EnumBuilder 枚举构建器
type EnumBuilder[T constraints.Ordered] struct {
	registry *EnumRegistry[T]
}

// NewEnumBuilder 创建枚举构建器
func NewEnumBuilder[T constraints.Ordered]() *EnumBuilder[T] {
	return &EnumBuilder[T]{
		registry: NewEnumRegistry[T](),
	}
}

// Add 添加枚举值
func (b *EnumBuilder[T]) Add(value T, name, desc string) *EnumBuilder[T] {
	b.registry.Define(value, name, desc)
	return b
}

// AddSimple 简单添加枚举值
func (b *EnumBuilder[T]) AddSimple(value T, name string) *EnumBuilder[T] {
	b.registry.DefineSimple(value, name)
	return b
}

// Build 构建枚举注册表
func (b *EnumBuilder[T]) Build() *EnumRegistry[T] {
	return b.registry
}

// 常用枚举示例和工厂函数

// CreateStatusEnum 创建状态枚举（示例）
func CreateStatusEnum() IntEnum {
	return NewEnumBuilder[int]().
		Add(0, "INACTIVE", "非活跃").
		Add(1, "ACTIVE", "活跃").
		Add(2, "PENDING", "待处理").
		Add(3, "SUSPENDED", "暂停").
		Add(4, "DELETED", "已删除").
		Build()
}

// CreateLevelEnum 创建级别枚举（示例）
func CreateLevelEnum() StringEnum {
	return NewEnumBuilder[string]().
		Add("DEBUG", "DEBUG", "调试").
		Add("INFO", "INFO", "信息").
		Add("WARN", "WARN", "警告").
		Add("ERROR", "ERROR", "错误").
		Add("FATAL", "FATAL", "致命").
		Build()
}

// 高性能查找相关

// FastLookup 快速查找结构体，预编译查找映射
type FastLookup[T constraints.Ordered] struct {
	valueMap map[T]*XEnum[T]
	nameMap  map[string]*XEnum[T]
}

// NewFastLookup 创建快速查找结构
func (r *EnumRegistry[T]) NewFastLookup() *FastLookup[T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return &FastLookup[T]{
		valueMap: r.values,
		nameMap:  r.names,
	}
}

// GetByValue 通过值快速查找
func (f *FastLookup[T]) GetByValue(value T) (*XEnum[T], bool) {
	enum, exists := f.valueMap[value]
	return enum, exists
}

// GetByName 通过名称快速查找
func (f *FastLookup[T]) GetByName(name string) (*XEnum[T], bool) {
	enum, exists := f.nameMap[name]
	return enum, exists
}

// 批量操作

// BatchValidator 批量验证器
type BatchValidator[T constraints.Ordered] struct {
	registry *EnumRegistry[T]
	valid    map[T]bool
	mu       sync.RWMutex
}

// NewBatchValidator 创建批量验证器
func (r *EnumRegistry[T]) NewBatchValidator() *BatchValidator[T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	valid := make(map[T]bool, len(r.values))
	for value := range r.values {
		valid[value] = true
	}

	return &BatchValidator[T]{
		registry: r,
		valid:    valid,
	}
}

// IsValid 快速验证单个值
func (v *BatchValidator[T]) IsValid(value T) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.valid[value]
}

// ValidateAll 批量验证
func (v *BatchValidator[T]) ValidateAll(values []T) []bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	results := make([]bool, len(values))
	for i, value := range values {
		results[i] = v.valid[value]
	}
	return results
}

// FilterValid 过滤有效值
func (v *BatchValidator[T]) FilterValid(values []T) []T {
	v.mu.RLock()
	defer v.mu.RUnlock()

	var valid []T
	for _, value := range values {
		if v.valid[value] {
			valid = append(valid, value)
		}
	}
	return valid
}

// 枚举集合操作

// EnumSet 枚举集合
type EnumSet[T constraints.Ordered] struct {
	items map[T]*XEnum[T]
	mu    sync.RWMutex
}

// NewEnumSet 创建枚举集合
func NewEnumSet[T constraints.Ordered]() *EnumSet[T] {
	return &EnumSet[T]{
		items: make(map[T]*XEnum[T]),
	}
}

// Add 添加枚举到集合
func (s *EnumSet[T]) Add(enum *XEnum[T]) *EnumSet[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[enum.value] = enum
	return s
}

// Remove 从集合移除枚举
func (s *EnumSet[T]) Remove(enum *XEnum[T]) *EnumSet[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, enum.value)
	return s
}

// Contains 检查集合是否包含枚举
func (s *EnumSet[T]) Contains(enum *XEnum[T]) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.items[enum.value]
	return exists
}

// ContainsValue 检查集合是否包含值
func (s *EnumSet[T]) ContainsValue(value T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.items[value]
	return exists
}

// Size 获取集合大小
func (s *EnumSet[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// ToSlice 转换为切片
func (s *EnumSet[T]) ToSlice() []*XEnum[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*XEnum[T], 0, len(s.items))
	for _, enum := range s.items {
		result = append(result, enum)
	}
	return result
}

// Clear 清空集合
func (s *EnumSet[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make(map[T]*XEnum[T])
}

// 实用工具函数

// EnumToMap 将枚举转换为映射
func EnumToMap[T constraints.Ordered](registry *EnumRegistry[T]) map[T]string {
	all := registry.All()
	result := make(map[T]string, len(all))
	for _, enum := range all {
		result[enum.value] = enum.name
	}
	return result
}

// EnumToSlice 将枚举转换为切片
func EnumToSlice[T constraints.Ordered](registry *EnumRegistry[T]) []T {
	return registry.Values()
}

// ParseEnum 解析字符串为枚举（支持名称和值）
func ParseEnum[T constraints.Ordered](registry *EnumRegistry[T], str string) (*XEnum[T], error) {
	// 先尝试按名称查找
	if enum, exists := registry.FromName(str); exists {
		return enum, nil
	}

	// 尝试解析为值
	var value T
	switch any(value).(type) {
	case int:
		parsed := Str(str).Int()
		if enum, exists := registry.FromValue(any(parsed).(T)); exists {
			return enum, nil
		}
	case string:
		if enum, exists := registry.FromValue(any(str).(T)); exists {
			return enum, nil
		}
	}

	return nil, fmt.Errorf("invalid enum value: %s", str)
}

// EnumRange 创建枚举范围
func EnumRange[T constraints.Ordered](registry *EnumRegistry[T], start, end T) []*XEnum[T] {
	var result []*XEnum[T]
	for _, enum := range registry.All() {
		if enum.value >= start && enum.value <= end {
			result = append(result, enum)
		}
	}
	return result
}

// MustGetEnum 必须获取枚举，不存在则panic
func MustGetEnum[T constraints.Ordered](registry *EnumRegistry[T], value T) *XEnum[T] {
	if enum, exists := registry.FromValue(value); exists {
		return enum
	}
	panic(fmt.Sprintf("enum value not found: %v", value))
}

// 类型安全的枚举定义宏

// DefineEnum 定义枚举的便捷宏
func DefineEnum[T constraints.Ordered](definitions map[T]string) *EnumRegistry[T] {
	registry := NewEnumRegistry[T]()
	for value, name := range definitions {
		registry.DefineSimple(value, name)
	}
	return registry
}

// DefineEnumWithDesc 定义带描述的枚举
func DefineEnumWithDesc[T constraints.Ordered](definitions map[T][2]string) *EnumRegistry[T] {
	registry := NewEnumRegistry[T]()
	for value, nameDesc := range definitions {
		registry.Define(value, nameDesc[0], nameDesc[1])
	}
	return registry
}
