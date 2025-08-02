package mule

import (
	"reflect"
)

// XArrayAny 不限类型的数组/切片 可提供方法有限
type XArrayAny[T any] []T

func ArrayAny[T any](arr []T) XArrayAny[T] {
	return arr
}

func ArrayAnys[T any](arr ...T) XArrayAny[T] {
	return arr
}

func (a XArrayAny[T]) Len() int {
	return len(a)
}

func (a XArrayAny[T]) Index(x T) int {
	for i, y := range a {
		if a.EqualItem(x, y) {
			return i
		}
	}
	return -1
}

func (a XArrayAny[T]) Exist(x T) bool {
	for _, y := range a {
		if a.EqualItem(x, y) {
			return true
		}
	}
	return false
}

func (a XArrayAny[T]) Equal(y XArrayAny[T]) bool {
	return reflect.DeepEqual(a, y)
}

func (a XArrayAny[T]) EqualItemIndex(index int, y T) bool {
	return reflect.DeepEqual(a[index], y)
}

func (a XArrayAny[T]) EqualItem(x, y T) bool {
	return reflect.DeepEqual(x, y)
}

func (a XArrayAny[T]) Merge(x XArrayAny[T]) XArrayAny[T] {
	return append(a, x...)
}

// Unique 数组去重, 如果数组元素的类型是 integer|float|string这些元素 请使用 ArrayOrdered.Unique 或 ArrayOrdered.UniqueOrdered
func (a XArrayAny[T]) Unique() XArrayAny[T] {
	var y = make(XArrayAny[T], 0)
	for _, x := range a {
		if !y.Exist(x) {
			y = append(y, x)
		}
	}
	return y
}
