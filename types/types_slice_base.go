package types

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/exp/constraints"
)

type XArray[T constraints.Ordered] []T

func Array[T constraints.Ordered](a []T) XArray[T] {
	return a
}
func Arrays[T constraints.Ordered](a ...T) XArray[T] {
	return a
}

func (a XArray[T]) Len() int {
	return len(a)
}

func (a XArray[T]) Index(x T) int {
	for i, y := range a {
		if a.EqualItem(x, y) {
			return i
		}
	}
	return -1
}

func (a XArray[T]) Exist(x T) bool {
	for _, y := range a {
		if a.EqualItem(x, y) {
			return true
		}
	}
	return false
}

func (a XArray[T]) Equal(y XArray[T]) bool {
	l := a.Len()
	if l != y.Len() {
		return false
	}

	for i := 0; i < l; i++ {
		if !a.EqualItem(a[i], y[i]) {
			return false
		}
	}
	return true
}

func (a XArray[T]) EqualItemIndex(index int, y T) bool {
	return a[index] == y
}

func (a XArray[T]) EqualItem(x, y T) bool {
	return x == y
}

func (a XArray[T]) Merge(x XArray[T]) XArray[T] {
	return append(a, x...)
}

// Unique 不保留原数组顺序
func (a XArray[T]) Unique() XArray[T] {
	m := make(map[T]bool)
	for _, x := range a {
		if _, ok := m[x]; !ok {
			m[x] = true
		}
	}
	var y XArray[T]
	for k := range m {
		y = append(y, k)
	}
	return y
}

// UniqueOrdered 保留原数组顺序
func (a XArray[T]) UniqueOrdered() XArray[T] {
	var y = make(XArray[T], 0)
	for _, x := range a {
		if !y.Exist(x) {
			y = append(y, x)
		}
	}
	return y
}

func (a XArray[T]) Remove(index int) XArray[T] {
	return slices.Delete(a, index, index+1)
}

func (a XArray[T]) RemoveValue(value T) XArray[T] {
	return slices.DeleteFunc(a, func(v T) bool {
		return v == value
	})
}

func (a XArray[T]) Push(v T) XArray[T] {
	return append(a, v)
}

func (a XArray[T]) Replace(oldVal, newVal T) XArray[T] {
	if a.Len() > 0 {
		for i, x := range a {
			if x == oldVal {
				a[i] = newVal
			}
		}
	}
	return a
}

func (a XArray[T]) Join(sep string) string {
	if a.Len() == 0 {
		return ""
	}
	return strings.Join(a.ToString(), sep)
}

func (a XArray[T]) ToString() XArray[string] {
	var y = make([]string, a.Len())
	for i, v := range a {
		y[i] = fmt.Sprint(v)
	}
	return y
}

func (a XArray[T]) ToInt() XArray[int] {
	var x any = a
	switch x.(type) {
	case XArray[int]:
		return x.(XArray[int])
	case XArray[int8]:
		return ArrayInteger2Int(x.(XArray[int8]))
	case XArray[int16]:
		return ArrayInteger2Int(x.(XArray[int16]))
	case XArray[int32]:
		return ArrayInteger2Int(x.(XArray[int32]))
	case XArray[int64]:
		return ArrayInteger2Int(x.(XArray[int64]))
	case XArray[uint]:
		return ArrayInteger2Int(x.(XArray[uint]))
	case XArray[uint8]:
		return ArrayInteger2Int(x.(XArray[uint8]))
	case XArray[uint16]:
		return ArrayInteger2Int(x.(XArray[uint16]))
	case XArray[uint32]:
		return ArrayInteger2Int(x.(XArray[uint32]))
	case XArray[uint64]:
		return ArrayInteger2Int(x.(XArray[uint64]))
	case XArray[float32]:
		return ArrayFloat2Int(x.(XArray[float32]))
	case XArray[float64]:
		return ArrayFloat2Int(x.(XArray[float64]))
	case XArray[string]:
		return ArrayString2Int(x.(XArray[string]))
	default:
		return nil
	}
}

func (a XArray[T]) ToInt64() XArray[int64] {
	var x any = a
	switch x.(type) {
	case XArray[int]:
		return ArrayInteger2Int64(x.(XArray[int]))
	case XArray[int8]:
		return ArrayInteger2Int64(x.(XArray[int8]))
	case XArray[int16]:
		return ArrayInteger2Int64(x.(XArray[int16]))
	case XArray[int32]:
		return ArrayInteger2Int64(x.(XArray[int32]))
	case XArray[int64]:
		return x.(XArray[int64])
	case XArray[uint]:
		return ArrayInteger2Int64(x.(XArray[uint]))
	case XArray[uint8]:
		return ArrayInteger2Int64(x.(XArray[uint8]))
	case XArray[uint16]:
		return ArrayInteger2Int64(x.(XArray[uint16]))
	case XArray[uint32]:
		return ArrayInteger2Int64(x.(XArray[uint32]))
	case XArray[uint64]:
		return ArrayInteger2Int64(x.(XArray[uint64]))
	case XArray[float32]:
		return ArrayFloat2Int64(x.(XArray[float32]))
	case XArray[float64]:
		return ArrayFloat2Int64(x.(XArray[float64]))
	case XArray[string]:
		return ArrayString2Int64(x.(XArray[string]))
	default:
		return nil
	}
}

func (a XArray[T]) ToFloat64() XArray[float64] {
	var x any = a
	switch x.(type) {
	case XArray[int]:
		return ArrayInteger2Float64(x.(XArray[int]))
	case XArray[int8]:
		return ArrayInteger2Float64(x.(XArray[int8]))
	case XArray[int16]:
		return ArrayInteger2Float64(x.(XArray[int16]))
	case XArray[int32]:
		return ArrayInteger2Float64(x.(XArray[int32]))
	case XArray[int64]:
		return ArrayInteger2Float64(x.(XArray[int64]))
	case XArray[uint]:
		return ArrayInteger2Float64(x.(XArray[uint]))
	case XArray[uint8]:
		return ArrayInteger2Float64(x.(XArray[uint8]))
	case XArray[uint16]:
		return ArrayInteger2Float64(x.(XArray[uint16]))
	case XArray[uint32]:
		return ArrayInteger2Float64(x.(XArray[uint32]))
	case XArray[uint64]:
		return ArrayInteger2Float64(x.(XArray[uint64]))
	case XArray[float32]:
		return ArrayFloat2Float64(x.(XArray[float32]))
	case XArray[float64]:
		return x.(XArray[float64])
	case XArray[string]:
		return ArrayString2Float64(x.(XArray[string]))
	default:
		return nil
	}
}
