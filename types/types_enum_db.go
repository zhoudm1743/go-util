package types

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	"golang.org/x/exp/constraints"
)

// 数据库接口实现，兼容 GORM 和其他 ORM

// DBValue 实现 driver.Valuer 接口，用于将枚举值写入数据库
func (e *XEnum[T]) DBValue() (driver.Value, error) {
	if e == nil {
		return nil, nil
	}
	return e.value, nil
}

// Scanner 接口实现，用于从数据库读取值到枚举
func (e *XEnum[T]) Scan(value interface{}) error {
	if value == nil {
		*e = XEnum[T]{}
		return nil
	}

	var v T
	switch any(v).(type) {
	case int:
		switch val := value.(type) {
		case int64:
			e.value = any(int(val)).(T)
		case int32:
			e.value = any(int(val)).(T)
		case int:
			e.value = any(val).(T)
		case []byte:
			if parsed := Str(string(val)).Int(); parsed != 0 || string(val) == "0" {
				e.value = any(parsed).(T)
			} else {
				return fmt.Errorf("cannot convert %v to int", value)
			}
		case string:
			if parsed := Str(val).Int(); parsed != 0 || val == "0" {
				e.value = any(parsed).(T)
			} else {
				return fmt.Errorf("cannot convert %s to int", val)
			}
		default:
			return fmt.Errorf("cannot scan %T into XEnum[int]", value)
		}
	case string:
		switch val := value.(type) {
		case string:
			e.value = any(val).(T)
		case []byte:
			e.value = any(string(val)).(T)
		default:
			return fmt.Errorf("cannot scan %T into XEnum[string]", value)
		}
	default:
		// 尝试直接类型转换
		if reflect.TypeOf(value).ConvertibleTo(reflect.TypeOf(v)) {
			converted := reflect.ValueOf(value).Convert(reflect.TypeOf(v))
			e.value = converted.Interface().(T)
		} else {
			return fmt.Errorf("cannot scan %T into XEnum[%T]", value, v)
		}
	}

	return nil
}

// GormDataType 返回 GORM 数据类型（可选实现）
func (e *XEnum[T]) GormDataType() string {
	var v T
	switch any(v).(type) {
	case int, int8, int16, int32, int64:
		return "integer"
	case uint, uint8, uint16, uint32, uint64:
		return "integer"
	case string:
		return "varchar(255)"
	case float32, float64:
		return "decimal"
	default:
		return "text"
	}
}

// EnumField GORM 枚举字段包装器
type EnumField[T constraints.Ordered] struct {
	*XEnum[T]
	registry *EnumRegistry[T]
}

// NewEnumField 创建 GORM 枚举字段
func NewEnumField[T constraints.Ordered](registry *EnumRegistry[T]) *EnumField[T] {
	return &EnumField[T]{
		XEnum:    &XEnum[T]{},
		registry: registry,
	}
}

// SetValue 设置枚举值
func (f *EnumField[T]) SetValue(value T) error {
	if enum, exists := f.registry.FromValue(value); exists {
		f.XEnum = enum
		return nil
	}
	return fmt.Errorf("invalid enum value: %v", value)
}

// SetName 设置枚举名称
func (f *EnumField[T]) SetName(name string) error {
	if enum, exists := f.registry.FromName(name); exists {
		f.XEnum = enum
		return nil
	}
	return fmt.Errorf("invalid enum name: %s", name)
}

// Value 实现 driver.Valuer 接口
func (f *EnumField[T]) Value() (driver.Value, error) {
	if f.XEnum == nil {
		return nil, nil
	}
	return f.XEnum.DBValue()
}

// Scan 实现 sql.Scanner 接口，自动验证枚举值
func (f *EnumField[T]) Scan(value interface{}) error {
	if value == nil {
		f.XEnum = nil
		return nil
	}

	// 临时枚举用于扫描
	temp := &XEnum[T]{}
	if err := temp.Scan(value); err != nil {
		return err
	}

	// 验证枚举值是否有效
	if enum, exists := f.registry.FromValue(temp.value); exists {
		f.XEnum = enum
		return nil
	}

	return fmt.Errorf("invalid enum value in database: %v", temp.value)
}

// GormDataType 实现 GORM 数据类型接口
func (f *EnumField[T]) GormDataType() string {
	if f.XEnum != nil {
		return f.XEnum.GormDataType()
	}
	var v T
	switch any(v).(type) {
	case int, int8, int16, int32, int64:
		return "integer"
	case uint, uint8, uint16, uint32, uint64:
		return "integer"
	case string:
		return "varchar(255)"
	case float32, float64:
		return "decimal"
	default:
		return "text"
	}
}

// 全局枚举注册表管理器
var globalEnumRegistries = make(map[string]interface{})

// RegisterGlobalEnum 注册全局枚举（用于 GORM 标签）
func RegisterGlobalEnum[T constraints.Ordered](name string, registry *EnumRegistry[T]) {
	globalEnumRegistries[name] = registry
}

// GetGlobalEnum 获取全局枚举
func GetGlobalEnum[T constraints.Ordered](name string) (*EnumRegistry[T], bool) {
	if registry, exists := globalEnumRegistries[name]; exists {
		if typedRegistry, ok := registry.(*EnumRegistry[T]); ok {
			return typedRegistry, true
		}
	}
	return nil, false
}

// 便捷的 GORM 集成函数

// DefineGormEnum 定义并注册 GORM 枚举
func DefineGormEnum[T constraints.Ordered](name string, definitions map[T]string) *EnumRegistry[T] {
	registry := DefineEnum(definitions)
	RegisterGlobalEnum(name, registry)
	return registry
}

// DefineGormEnumWithDesc 定义并注册带描述的 GORM 枚举
func DefineGormEnumWithDesc[T constraints.Ordered](name string, definitions map[T][2]string) *EnumRegistry[T] {
	registry := DefineEnumWithDesc(definitions)
	RegisterGlobalEnum(name, registry)
	return registry
}

// NewGormEnumField 创建 GORM 枚举字段（通过注册表名称）
func NewGormEnumField[T constraints.Ordered](registryName string) *EnumField[T] {
	if registry, exists := GetGlobalEnum[T](registryName); exists {
		return NewEnumField(registry)
	}
	panic(fmt.Sprintf("enum registry not found: %s", registryName))
}

// 示例预定义枚举（带 GORM 注册）

// 用户状态枚举（注册到 GORM）
var UserStatusRegistry = DefineGormEnum("user_status", map[int]string{
	0: "INACTIVE",
	1: "ACTIVE",
	2: "PENDING",
	3: "SUSPENDED",
	4: "DELETED",
})

// 日志级别枚举（注册到 GORM）
var LogLevelRegistry = DefineGormEnum("log_level", map[string]string{
	"DEBUG": "DEBUG",
	"INFO":  "INFO",
	"WARN":  "WARN",
	"ERROR": "ERROR",
	"FATAL": "FATAL",
})

// 用户角色枚举（注册到 GORM）
var UserRoleRegistry = DefineGormEnumWithDesc("user_role", map[int][2]string{
	1: {"SUPER_ADMIN", "超级管理员"},
	2: {"ADMIN", "管理员"},
	3: {"MODERATOR", "版主"},
	4: {"USER", "普通用户"},
	5: {"GUEST", "访客"},
})

// 辅助函数，用于 GORM 查询

// CreateEnumQuery 创建枚举查询条件辅助函数
func CreateEnumQuery[T constraints.Ordered](column string, enum *XEnum[T]) (string, interface{}) {
	if enum == nil {
		return column + " IS NULL", nil
	}
	return column + " = ?", enum.value
}

// CreateEnumInQuery 创建枚举 IN 查询条件
func CreateEnumInQuery[T constraints.Ordered](column string, enums []*XEnum[T]) (string, []interface{}) {
	if len(enums) == 0 {
		return "1 = 0", nil // 返回空结果
	}

	values := make([]interface{}, 0, len(enums))
	for _, enum := range enums {
		if enum != nil {
			values = append(values, enum.value)
		}
	}

	return column + " IN ?", values
}

// ValidateEnumField 验证枚举字段值
func ValidateEnumField[T constraints.Ordered](registry *EnumRegistry[T], value T) error {
	if !registry.IsValid(value) {
		return fmt.Errorf("invalid enum value: %v", value)
	}
	return nil
}

// ConvertToEnum 将数据库值转换为枚举
func ConvertToEnum[T constraints.Ordered](registry *EnumRegistry[T], value T) (*XEnum[T], error) {
	if enum, exists := registry.FromValue(value); exists {
		return enum, nil
	}
	return nil, fmt.Errorf("invalid enum value: %v", value)
}

// GetEnumValues 获取枚举的所有有效值（用于数据库约束）
func GetEnumValues[T constraints.Ordered](registry *EnumRegistry[T]) []T {
	return registry.Values()
}

// GetEnumNames 获取枚举的所有有效名称（用于数据库约束）
func GetEnumNames[T constraints.Ordered](registry *EnumRegistry[T]) []string {
	return registry.Names()
}
