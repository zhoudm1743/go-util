package six

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"strconv"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func ArrayInteger2Int64[T constraints.Integer](in XArray[T]) XArray[int64] {
	out := make(XArray[int64], in.Len())
	for i := 0; i < in.Len(); i++ {
		out[i] = int64(in[i])
	}
	return out
}

func ArrayFloat2Int64[T constraints.Float](in XArray[T]) XArray[int64] {
	out := make(XArray[int64], in.Len())
	for i := 0; i < in.Len(); i++ {
		v, _ := strconv.Atoi(fmt.Sprintf("%.0f", in[i]))
		out[i] = int64(v)
	}
	return out
}

func ArrayString2Int64(in XArray[string]) XArray[int64] {
	out := make(XArray[int64], in.Len())
	for i := 0; i < in.Len(); i++ {
		v, _ := strconv.Atoi(in[i])
		out[i] = int64(v)
	}
	return out
}

func ArrayInteger2Int[T constraints.Integer](in XArray[T]) XArray[int] {
	out := make(XArray[int], in.Len())
	for i := 0; i < in.Len(); i++ {
		out[i] = int(in[i])
	}
	return out
}

func ArrayFloat2Int[T constraints.Float](in XArray[T]) XArray[int] {
	out := make(XArray[int], in.Len())
	for i := 0; i < in.Len(); i++ {
		out[i], _ = strconv.Atoi(fmt.Sprintf("%.0f", in[i]))
	}
	return out
}

func ArrayString2Int(in XArray[string]) XArray[int] {
	out := make(XArray[int], in.Len())
	for i := 0; i < in.Len(); i++ {
		out[i], _ = strconv.Atoi(in[i])
	}
	return out
}

func ArrayNumber2String[T Number](in XArray[T]) XArray[string] {
	out := make(XArray[string], in.Len())
	for i := 0; i < in.Len(); i++ {
		out[i] = fmt.Sprint(in[i])
	}
	return out
}

func ArrayInteger2Float64[T constraints.Integer](in XArray[T]) XArray[float64] {
	out := make(XArray[float64], in.Len())
	for i := 0; i < in.Len(); i++ {
		out[i] = float64(in[i])
	}
	return out
}

func ArrayFloat2Float64[T constraints.Float](in XArray[T]) XArray[float64] {
	out := make(XArray[float64], in.Len())
	for i := 0; i < in.Len(); i++ {
		out[i] = float64(in[i])
	}
	return out
}

func ArrayString2Float64(in XArray[string]) XArray[float64] {
	out := make(XArray[float64], in.Len())
	for i := 0; i < in.Len(); i++ {
		out[i], _ = strconv.ParseFloat(in[i], 64)
	}
	return out
}
