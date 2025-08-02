package six

import "strconv"

type XBool bool

func (b XBool) Int() int {
	if b {
		return 1
	}
	return 0
}

func (b XBool) Int64() int64 {
	return int64(b.Int())
}

func (b XBool) Float() float64 {
	return float64(b.Int())
}

func (b XBool) String() string {
	return strconv.FormatBool(bool(b))
}

func (b XBool) Bool() bool {
	return bool(b)
}
