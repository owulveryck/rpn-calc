package engine

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var ErrInvalidNumber = errors.New("invalid number")

type Value interface {
	String() string
	Float64() float64
	IsComplex() bool
	Complex128() complex128
	StringInBase(base BaseMode) string
}

type BaseMode int

const (
	BaseDec BaseMode = iota
	BaseHex
	BaseOct
	BaseBin
)

func (b BaseMode) String() string {
	switch b {
	case BaseHex:
		return "HEX"
	case BaseOct:
		return "OCT"
	case BaseBin:
		return "BIN"
	default:
		return "DEC"
	}
}

type RealValue struct {
	val float64
}

func NewRealValue(v float64) RealValue {
	return RealValue{val: v}
}

func ParseValue(s string) (Value, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, ErrInvalidNumber
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidNumber, s)
	}
	return RealValue{val: f}, nil
}

func (v RealValue) String() string {
	if v.val == math.Trunc(v.val) && !math.IsInf(v.val, 0) && !math.IsNaN(v.val) && math.Abs(v.val) < 1e15 {
		return strconv.FormatFloat(v.val, 'f', -1, 64)
	}
	return strconv.FormatFloat(v.val, 'g', -1, 64)
}

func (v RealValue) Float64() float64 {
	return v.val
}

func (v RealValue) IsComplex() bool {
	return false
}

func (v RealValue) Complex128() complex128 {
	return complex(v.val, 0)
}

func (v RealValue) StringInBase(base BaseMode) string {
	n := int64(v.val)
	switch base {
	case BaseHex:
		return fmt.Sprintf("#%X", n)
	case BaseOct:
		return fmt.Sprintf("%oo", n)
	case BaseBin:
		return fmt.Sprintf("%bb", n)
	default:
		return v.String()
	}
}
