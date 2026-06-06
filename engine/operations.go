package engine

import (
	"errors"
	"math"
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

func binaryOp(c *Calculator, fn func(x, y float64) (float64, error)) error {
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
	result, err := fn(x.Float64(), y.Float64())
	if err != nil {
		c.stack.Push(y)
		c.stack.Push(x)
		return err
	}
	c.stack.Push(NewRealValue(result))
	return nil
}

func unaryOp(c *Calculator, fn func(x float64) (float64, error)) error {
	x, err := c.stack.Pop()
	if err != nil {
		return err
	}
	c.lastArgs = []Value{x}
	result, err := fn(x.Float64())
	if err != nil {
		c.stack.Push(x)
		return err
	}
	c.stack.Push(NewRealValue(result))
	return nil
}

func opAdd(c *Calculator) error {
	return binaryOp(c, func(x, y float64) (float64, error) { return y + x, nil })
}

func opSub(c *Calculator) error {
	return binaryOp(c, func(x, y float64) (float64, error) { return y - x, nil })
}

func opMul(c *Calculator) error {
	return binaryOp(c, func(x, y float64) (float64, error) { return y * x, nil })
}

func opDiv(c *Calculator) error {
	return binaryOp(c, func(x, y float64) (float64, error) {
		if x == 0 {
			return 0, ErrDivisionByZero
		}
		return y / x, nil
	})
}

func opNeg(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) { return -x, nil })
}

func opInv(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) {
		if x == 0 {
			return 0, ErrDivisionByZero
		}
		return 1 / x, nil
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
	return unaryOp(c, func(x float64) (float64, error) {
		return math.Sin(c.toRadians(x)), nil
	})
}

func opCos(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) {
		return math.Cos(c.toRadians(x)), nil
	})
}

func opTan(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) {
		return math.Tan(c.toRadians(x)), nil
	})
}

func opAsin(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) {
		return c.fromRadians(math.Asin(x)), nil
	})
}

func opAcos(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) {
		return c.fromRadians(math.Acos(x)), nil
	})
}

func opAtan(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) {
		return c.fromRadians(math.Atan(x)), nil
	})
}

func opLog(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) {
		if x <= 0 {
			return 0, errors.New("logarithm of non-positive number")
		}
		return math.Log10(x), nil
	})
}

func opLn(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) {
		if x <= 0 {
			return 0, errors.New("logarithm of non-positive number")
		}
		return math.Log(x), nil
	})
}

func opExp(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) { return math.Exp(x), nil })
}

func op10X(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) { return math.Pow(10, x), nil })
}

func opPow(c *Calculator) error {
	return binaryOp(c, func(x, y float64) (float64, error) { return math.Pow(y, x), nil })
}

func opSqrt(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) {
		if x < 0 {
			return 0, errors.New("square root of negative number")
		}
		return math.Sqrt(x), nil
	})
}

func opSq(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) { return x * x, nil })
}

func opFactorial(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) {
		if x < 0 {
			return 0, errors.New("factorial of negative number")
		}
		if x == math.Trunc(x) && x <= 170 {
			n := int(x)
			result := 1.0
			for i := 2; i <= n; i++ {
				result *= float64(i)
			}
			return result, nil
		}
		return math.Gamma(x + 1), nil
	})
}

func opAbs(c *Calculator) error {
	return unaryOp(c, func(x float64) (float64, error) { return math.Abs(x), nil })
}

func opPi(c *Calculator) error {
	c.stack.Push(NewRealValue(math.Pi))
	return nil
}

func opE(c *Calculator) error {
	c.stack.Push(NewRealValue(math.E))
	return nil
}

func opMin(c *Calculator) error {
	return binaryOp(c, func(x, y float64) (float64, error) { return math.Min(y, x), nil })
}

func opMax(c *Calculator) error {
	return binaryOp(c, func(x, y float64) (float64, error) { return math.Max(y, x), nil })
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
