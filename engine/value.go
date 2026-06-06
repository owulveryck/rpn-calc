package engine

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

var ErrInvalidNumber = errors.New("invalid number")

type Value interface {
	String() string
	Float64() float64
	BigFloat() *big.Float
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
	val *big.Float
}

func NewRealValue(v float64) RealValue {
	bf := new(big.Float).SetPrec(defaultPrec)
	if math.IsNaN(v) {
		bf.SetFloat64(0)
	} else {
		bf.SetFloat64(v)
	}
	return RealValue{val: bf}
}

func ParseValue(s string) (Value, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, ErrInvalidNumber
	}
	f, _, err := new(big.Float).SetPrec(defaultPrec).Parse(s, 10)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidNumber, s)
	}
	return RealValue{val: f}, nil
}

func ParseValueInBase(s string, base BaseMode) (Value, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, ErrInvalidNumber
	}
	if base == BaseDec {
		return ParseValue(s)
	}
	var b int
	switch base {
	case BaseHex:
		b = 16
	case BaseOct:
		b = 8
	case BaseBin:
		b = 2
	}
	n, err := strconv.ParseInt(s, b, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidNumber, s)
	}
	return NewRealValue(float64(n)), nil
}

func formatFloat(bf *big.Float) string {
	f, _ := bf.Float64()
	if !math.IsInf(f, 0) {
		if f == math.Trunc(f) && math.Abs(f) < 1e15 {
			return strconv.FormatFloat(f, 'f', -1, 64)
		}
		return strconv.FormatFloat(f, 'g', -1, 64)
	}
	return bf.Text('g', 15)
}

func (v RealValue) String() string {
	return formatFloat(v.val)
}

func (v RealValue) Float64() float64 {
	f, _ := v.val.Float64()
	return f
}

func (v RealValue) BigFloat() *big.Float {
	return new(big.Float).SetPrec(defaultPrec).Copy(v.val)
}

func (v RealValue) IsComplex() bool {
	return false
}

func (v RealValue) Complex128() complex128 {
	return complex(v.Float64(), 0)
}

func (v RealValue) StringInBase(base BaseMode) string {
	n := int64(v.Float64())
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

type ExprValue struct {
	node *ExprNode
}

func NewExprValue(node *ExprNode) ExprValue {
	return ExprValue{node: node}
}

func (v ExprValue) String() string {
	return formatFloat(v.node.evaluate())
}

func (v ExprValue) Float64() float64 {
	f, _ := v.node.evaluate().Float64()
	return f
}

func (v ExprValue) BigFloat() *big.Float {
	return v.node.evaluate()
}

func (v ExprValue) IsComplex() bool {
	return false
}

func (v ExprValue) Complex128() complex128 {
	return complex(v.Float64(), 0)
}

func (v ExprValue) StringInBase(base BaseMode) string {
	n := int64(v.Float64())
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
