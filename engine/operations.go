package engine

import (
	"errors"
	"math/big"
)

// ErrDivisionByZero est retourné quand une division ou un inverse par zéro est tenté.
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
	"FACT": {opFact, arityUnary, "!"},

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

// binaryGraphOp pops two values, builds a binary graph node, and pushes the result.
func binaryGraphOp(c *Calculator, op ExprOp) error {
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
	node := binaryNode(op, toExprNode(y), toExprNode(x))
	c.stack.Push(NewExprValue(node))
	return nil
}

// binaryGraphOpChecked pops two values, runs a pre-check, builds a graph node.
func binaryGraphOpChecked(c *Calculator, op ExprOp, check func(x, y *big.Float) error) error {
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
	if check != nil {
		if err := check(x.BigFloat(), y.BigFloat()); err != nil {
			c.stack.Push(y)
			c.stack.Push(x)
			return err
		}
	}
	node := binaryNode(op, toExprNode(y), toExprNode(x))
	c.stack.Push(NewExprValue(node))
	return nil
}

// unaryGraphOp pops one value, builds a unary graph node, and pushes the result.
func unaryGraphOp(c *Calculator, op ExprOp) error {
	x, err := c.stack.Pop()
	if err != nil {
		return err
	}
	c.lastArgs = []Value{x}
	node := unaryNode(op, toExprNode(x))
	c.stack.Push(NewExprValue(node))
	return nil
}

// unaryGraphOpChecked pops one value, runs a pre-check, builds a graph node.
func unaryGraphOpChecked(c *Calculator, op ExprOp, check func(x *big.Float) error) error {
	x, err := c.stack.Pop()
	if err != nil {
		return err
	}
	c.lastArgs = []Value{x}
	if check != nil {
		if err := check(x.BigFloat()); err != nil {
			c.stack.Push(x)
			return err
		}
	}
	node := unaryNode(op, toExprNode(x))
	c.stack.Push(NewExprValue(node))
	return nil
}

func opAdd(c *Calculator) error { return binaryGraphOp(c, OpAdd) }
func opSub(c *Calculator) error { return binaryGraphOp(c, OpSub) }
func opMul(c *Calculator) error { return binaryGraphOp(c, OpMul) }

func opDiv(c *Calculator) error {
	return binaryGraphOpChecked(c, OpDiv, func(x, _ *big.Float) error {
		if x.Sign() == 0 {
			return ErrDivisionByZero
		}
		return nil
	})
}

func opNeg(c *Calculator) error { return unaryGraphOp(c, OpNeg) }

func opInv(c *Calculator) error {
	return unaryGraphOpChecked(c, OpInv, func(x *big.Float) error {
		if x.Sign() == 0 {
			return ErrDivisionByZero
		}
		return nil
	})
}

// trigGraphOp builds a trig graph node with angle conversion embedded in the graph.
func trigGraphOp(c *Calculator, op ExprOp) error {
	x, err := c.stack.Pop()
	if err != nil {
		return err
	}
	c.lastArgs = []Value{x}
	inputNode := toExprNode(x)
	var radiansNode *ExprNode
	switch c.angleMode {
	case AngleDeg:
		radiansNode = degToRadNode(inputNode)
	case AngleGrad:
		radiansNode = gradToRadNode(inputNode)
	default:
		radiansNode = inputNode
	}
	node := unaryNode(op, radiansNode)
	c.stack.Push(NewExprValue(node))
	return nil
}

// inverseTrigGraphOp builds an inverse trig graph node with result conversion.
func inverseTrigGraphOp(c *Calculator, op ExprOp) error {
	x, err := c.stack.Pop()
	if err != nil {
		return err
	}
	c.lastArgs = []Value{x}
	node := unaryNode(op, toExprNode(x))
	var resultNode *ExprNode
	switch c.angleMode {
	case AngleDeg:
		resultNode = radToDegNode(node)
	case AngleGrad:
		resultNode = radToGradNode(node)
	default:
		resultNode = node
	}
	c.stack.Push(NewExprValue(resultNode))
	return nil
}

func opSin(c *Calculator) error  { return trigGraphOp(c, OpSin) }
func opCos(c *Calculator) error  { return trigGraphOp(c, OpCos) }
func opTan(c *Calculator) error  { return trigGraphOp(c, OpTan) }
func opAsin(c *Calculator) error { return inverseTrigGraphOp(c, OpAsin) }
func opAcos(c *Calculator) error { return inverseTrigGraphOp(c, OpAcos) }
func opAtan(c *Calculator) error { return inverseTrigGraphOp(c, OpAtan) }

func opLog(c *Calculator) error {
	return unaryGraphOpChecked(c, OpLog, func(x *big.Float) error {
		if x.Sign() <= 0 {
			return errors.New("logarithm of non-positive number")
		}
		return nil
	})
}

func opLn(c *Calculator) error {
	return unaryGraphOpChecked(c, OpLn, func(x *big.Float) error {
		if x.Sign() <= 0 {
			return errors.New("logarithm of non-positive number")
		}
		return nil
	})
}

func opExp(c *Calculator) error  { return unaryGraphOp(c, OpExp) }
func op10X(c *Calculator) error  { return unaryGraphOp(c, Op10X) }
func opSqrt(c *Calculator) error {
	return unaryGraphOpChecked(c, OpSqrt, func(x *big.Float) error {
		if x.Sign() < 0 {
			return errors.New("square root of negative number")
		}
		return nil
	})
}
func opSq(c *Calculator) error { return unaryGraphOp(c, OpSq) }

func opFact(c *Calculator) error {
	return unaryGraphOpChecked(c, OpFact, func(x *big.Float) error {
		f, _ := x.Float64()
		if f < 0 {
			return errors.New("factorial of negative number")
		}
		return nil
	})
}

func opAbs(c *Calculator) error { return unaryGraphOp(c, OpAbs) }

func opPow(c *Calculator) error {
	return binaryGraphOpChecked(c, OpPow, func(x, y *big.Float) error {
		if y.Sign() < 0 {
			xi, acc := x.Int(nil)
			if acc != big.Exact {
				return errors.New("negative base with non-integer exponent")
			}
			_ = xi
		}
		return nil
	})
}

func opPi(c *Calculator) error {
	c.stack.Push(NewExprValue(piNode()))
	return nil
}

func opE(c *Calculator) error {
	c.stack.Push(NewExprValue(eNode()))
	return nil
}

func opMin(c *Calculator) error { return binaryGraphOp(c, OpMin) }
func opMax(c *Calculator) error { return binaryGraphOp(c, OpMax) }

func opSwap(c *Calculator) error  { return c.stack.Swap() }
func opDup(c *Calculator) error   { return c.stack.Dup() }
func opDrop(c *Calculator) error  { return c.stack.Drop() }
func opOver(c *Calculator) error  { return c.stack.Over() }
func opRot(c *Calculator) error   { return c.stack.Rot() }
func opClear(c *Calculator) error {
	c.stack.Clear()
	c.history = nil
	c.undoStack = nil
	c.lastArgs = nil
	return nil
}

func opDepth(c *Calculator) error {
	c.stack.Push(NewExprValue(literalFloat64Node(float64(c.stack.Depth()))))
	return nil
}
