package six

import (
	"golang.org/x/exp/constraints"
	"strconv"
)

type XInt int64

func Int[T constraints.Integer](i T) XInt {
	return XInt(i)
}

func (x XInt) Int() int {
	return int(x)
}

func (x XInt) Int64() int64 {
	return int64(x)
}

func (x XInt) String() string {
	return strconv.Itoa(x.Int())
}

func (x XInt) Float() float64 {
	return float64(x)
}

func (x XInt) Bool() bool {
	return x != 0
}

func (x XInt) IsZero() bool {
	return x == 0
}
