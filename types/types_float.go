package types

import (
	"fmt"
	"math"
	"strconv"

	"golang.org/x/exp/constraints"
)

type XFloat float64

func Float[T constraints.Float](f T) XFloat {
	return XFloat(f)
}

func (x XFloat) Equal(f float64) bool {
	return float64(x) == f
}

func (x XFloat) EqualX(f XFloat) bool {
	return x == f
}

func (x XFloat) Float() float64 {
	return float64(x)
}

func (x XFloat) Decimal() XDecimal {
	return Decimal(x.Float())
}

func (x XFloat) Int() int {
	i, _ := strconv.Atoi(x.Format(0))
	return i
}

func (x XFloat) Int64() int64 {
	return int64(x.Int())
}

func (x XFloat) String() string {
	return fmt.Sprintf("%f", x)
}

func (x XFloat) Format(decimals int) string {
	format := fmt.Sprintf("%%.%df", decimals)
	return fmt.Sprintf(format, x)
}

func (x XFloat) Ceil(decimals float64) XFloat {
	pow := math.Pow(10, decimals)
	return XFloat(math.Ceil(float64(x)*pow) / pow)
}

func (x XFloat) Floor(decimals float64) XFloat {
	pow := math.Pow(10, decimals)
	return XFloat(math.Floor(float64(x)*pow) / pow)
}

func (x XFloat) Round(decimals float64) XFloat {
	pow := math.Pow(10, decimals)
	return XFloat(math.Round(float64(x)*pow) / pow)
}

func (x XFloat) IsZero() bool {
	return x == 0
}
