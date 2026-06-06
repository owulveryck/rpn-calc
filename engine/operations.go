package engine

import (
	"errors"
	"math"
	"math/big"
)

var ErrDivisionByZero = errors.New("division by zero")

type opFunc func(c *Calculator) error

const (
	arityBinary = 2
	arityUnary  = 1
	arityConst  = 0
	arityStack  = -1
)

type opInfo struct {
	fn     opFunc
	arity  int
	symbol string
}

var operations = map[string]opInfo{
	"+":   {opAdd, arityBinary, "+"},
	"-":   {opSub, arityBinary, "-"},
	"*":   {opMul, arityBinary, "*"},
	"/":   {opDiv, arityBinary, "/"},
	"NEG": {opNeg, arityUnary, "neg"},
	"INV": {opInv, arityUnary, "1/"},

	"SIN":  {opSin, arityUnary, "sin"},
	"COS":  {opCos, arityUnary, "cos"},
	"TAN":  {opTan, arityUnary, "tan"},
	"ASIN": {opAsin, arityUnary, "asin"},
	"ACOS": {opAcos, arityUnary, "acos"},
	"ATAN": {opAtan, arityUnary, "atan"},

	"LOG":  {opLog, arityUnary, "log"},
	"LN":   {opLn, arityUnary, "ln"},
	"EXP":  {opExp, arityUnary, "exp"},
	"10^X": {op10X, arityUnary, "10^"},
	"Y^X":  {opPow, arityBinary, "^"},
	"SQRT": {opSqrt, arityUnary, "sqrt"},
	"SQ":   {opSq, arityUnary, "sq"},
	"FACT": {opFactorial, arityUnary, "!"},

	"ABS": {opAbs, arityUnary, "abs"},
	"PI":  {opPi, arityConst, "PI"},
	"E":   {opE, arityConst, "E"},
	"MIN": {opMin, arityBinary, "min"},
	"MAX": {opMax, arityBinary, "max"},

	"SWAP":  {opSwap, arityStack, "SWAP"},
	"DUP":   {opDup, arityStack, "DUP"},
	"DROP":  {opDrop, arityStack, "DROP"},
	"OVER":  {opOver, arityStack, "OVER"},
	"ROT":   {opRot, arityStack, "ROT"},
	"DEPTH": {opDepth, arityStack, "DEPTH"},
	"CLEAR": {opClear, arityStack, "CLEAR"},
}

func binaryOp(c *Calculator, fn func(x, y *big.Float) (*big.Float, error)) error {
	x, err := c.stack.Pop()
	if err != nil {
		return err
	}
	y, err := c.stack.Pop()
	if err != nil {
		c.stack.Push(x)
		return err
	}
	c.lastArgs = []Value{y, x}
	result, err := fn(x.BigFloat(), y.BigFloat())
	if err != nil {
		c.stack.Push(y)
		c.stack.Push(x)
		return err
	}
	c.stack.Push(newRealValueFromBig(result))
	return nil
}

func unaryOp(c *Calculator, fn func(x *big.Float) (*big.Float, error)) error {
	x, err := c.stack.Pop()
	if err != nil {
		return err
	}
	c.lastArgs = []Value{x}
	result, err := fn(x.BigFloat())
	if err != nil {
		c.stack.Push(x)
		return err
	}
	c.stack.Push(newRealValueFromBig(result))
	return nil
}

func opAdd(c *Calculator) error {
	return binaryOp(c, func(x, y *big.Float) (*big.Float, error) {
		return new(big.Float).SetPrec(defaultPrec).Add(y, x), nil
	})
}

func opSub(c *Calculator) error {
	return binaryOp(c, func(x, y *big.Float) (*big.Float, error) {
		return new(big.Float).SetPrec(defaultPrec).Sub(y, x), nil
	})
}

func opMul(c *Calculator) error {
	return binaryOp(c, func(x, y *big.Float) (*big.Float, error) {
		return new(big.Float).SetPrec(defaultPrec).Mul(y, x), nil
	})
}

func opDiv(c *Calculator) error {
	return binaryOp(c, func(x, y *big.Float) (*big.Float, error) {
		if x.Sign() == 0 {
			return nil, ErrDivisionByZero
		}
		return new(big.Float).SetPrec(defaultPrec).Quo(y, x), nil
	})
}

func opNeg(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		return new(big.Float).SetPrec(defaultPrec).Neg(x), nil
	})
}

func opInv(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		if x.Sign() == 0 {
			return nil, ErrDivisionByZero
		}
		one := new(big.Float).SetPrec(defaultPrec).SetInt64(1)
		return new(big.Float).SetPrec(defaultPrec).Quo(one, x), nil
	})
}

func (c *Calculator) toRadians(angle float64) float64 {
	switch c.angleMode {
	case AngleDeg:
		return angle * math.Pi / 180
	case AngleGrad:
		return angle * math.Pi / 200
	default:
		return angle
	}
}

func (c *Calculator) fromRadians(rad float64) float64 {
	switch c.angleMode {
	case AngleDeg:
		return rad * 180 / math.Pi
	case AngleGrad:
		return rad * 200 / math.Pi
	default:
		return rad
	}
}

func opSin(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		f, _ := x.Float64()
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(math.Sin(c.toRadians(f))), nil
	})
}

func opCos(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		f, _ := x.Float64()
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(math.Cos(c.toRadians(f))), nil
	})
}

func opTan(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		f, _ := x.Float64()
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(math.Tan(c.toRadians(f))), nil
	})
}

func opAsin(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		f, _ := x.Float64()
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(c.fromRadians(math.Asin(f))), nil
	})
}

func opAcos(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		f, _ := x.Float64()
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(c.fromRadians(math.Acos(f))), nil
	})
}

func opAtan(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		f, _ := x.Float64()
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(c.fromRadians(math.Atan(f))), nil
	})
}

func opLog(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		lnx, err := bigLn(x)
		if err != nil {
			return nil, err
		}
		return new(big.Float).SetPrec(defaultPrec).Quo(lnx, bigLn10()), nil
	})
}

func opLn(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		return bigLn(x)
	})
}

func opExp(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		return bigExp(x), nil
	})
}

func op10X(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		product := new(big.Float).SetPrec(defaultPrec).Mul(x, bigLn10())
		return bigExp(product), nil
	})
}

func opPow(c *Calculator) error {
	return binaryOp(c, func(x, y *big.Float) (*big.Float, error) {
		return bigPow(y, x)
	})
}

func opSqrt(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		if x.Sign() < 0 {
			return nil, errors.New("square root of negative number")
		}
		return new(big.Float).SetPrec(defaultPrec).Sqrt(x), nil
	})
}

func opSq(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		return new(big.Float).SetPrec(defaultPrec).Mul(x, x), nil
	})
}

func opFactorial(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		f, _ := x.Float64()
		if f < 0 {
			return nil, errors.New("factorial of negative number")
		}
		if x.IsInt() && f <= 170 {
			n := int(f)
			result := new(big.Float).SetPrec(defaultPrec).SetInt64(1)
			iBig := new(big.Float).SetPrec(defaultPrec)
			for i := 2; i <= n; i++ {
				iBig.SetInt64(int64(i))
				result.Mul(result, iBig)
			}
			return result, nil
		}
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(math.Gamma(f + 1)), nil
	})
}

func opAbs(c *Calculator) error {
	return unaryOp(c, func(x *big.Float) (*big.Float, error) {
		return new(big.Float).SetPrec(defaultPrec).Abs(x), nil
	})
}

func opPi(c *Calculator) error {
	c.stack.Push(newRealValueFromBig(bigPi()))
	return nil
}

func opE(c *Calculator) error {
	c.stack.Push(newRealValueFromBig(bigEConst()))
	return nil
}

func opMin(c *Calculator) error {
	return binaryOp(c, func(x, y *big.Float) (*big.Float, error) {
		if y.Cmp(x) <= 0 {
			return new(big.Float).SetPrec(defaultPrec).Copy(y), nil
		}
		return new(big.Float).SetPrec(defaultPrec).Copy(x), nil
	})
}

func opMax(c *Calculator) error {
	return binaryOp(c, func(x, y *big.Float) (*big.Float, error) {
		if y.Cmp(x) >= 0 {
			return new(big.Float).SetPrec(defaultPrec).Copy(y), nil
		}
		return new(big.Float).SetPrec(defaultPrec).Copy(x), nil
	})
}

func opSwap(c *Calculator) error  { return c.stack.Swap() }
func opDup(c *Calculator) error   { return c.stack.Dup() }
func opDrop(c *Calculator) error  { return c.stack.Drop() }
func opOver(c *Calculator) error  { return c.stack.Over() }
func opRot(c *Calculator) error   { return c.stack.Rot() }
func opClear(c *Calculator) error { c.stack.Clear(); return nil }

func opDepth(c *Calculator) error {
	c.stack.Push(NewRealValue(float64(c.stack.Depth())))
	return nil
}
