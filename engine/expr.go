package engine

import (
	"math"
	"math/big"
	"sync"
)

// ExprOp identifie une opération dans l'arbre d'expressions.
type ExprOp int

const (
	OpLiteral ExprOp = iota // Valeur littérale
	OpPi                    // Constante pi
	OpE                     // Constante e (nombre d'Euler)
	OpAdd                   // Addition
	OpSub                   // Soustraction
	OpMul                   // Multiplication
	OpDiv                   // Division
	OpNeg                   // Négation
	OpInv                   // Inverse (1/x)
	OpSin                   // Sinus
	OpCos                   // Cosinus
	OpTan                   // Tangente
	OpAsin                  // Arc sinus
	OpAcos                  // Arc cosinus
	OpAtan                  // Arc tangente
	OpLog                   // Logarithme décimal
	OpLn                    // Logarithme naturel
	OpExp                   // Exponentielle (e^x)
	Op10X                   // Puissance de 10 (10^x)
	OpPow                   // Puissance (y^x)
	OpSqrt                  // Racine carrée
	OpSq                    // Carré (x^2)
	OpFact                  // Factorielle
	OpAbs                   // Valeur absolue
	OpMin                   // Minimum
	OpMax                   // Maximum
)

// String retourne le nom textuel de l'opération.
func (op ExprOp) String() string {
	switch op {
	case OpLiteral:
		return "literal"
	case OpPi:
		return "pi"
	case OpE:
		return "e"
	case OpAdd:
		return "add"
	case OpSub:
		return "sub"
	case OpMul:
		return "mul"
	case OpDiv:
		return "div"
	case OpNeg:
		return "neg"
	case OpInv:
		return "inv"
	case OpSin:
		return "sin"
	case OpCos:
		return "cos"
	case OpTan:
		return "tan"
	case OpAsin:
		return "asin"
	case OpAcos:
		return "acos"
	case OpAtan:
		return "atan"
	case OpLog:
		return "log"
	case OpLn:
		return "ln"
	case OpExp:
		return "exp"
	case Op10X:
		return "10x"
	case OpPow:
		return "pow"
	case OpSqrt:
		return "sqrt"
	case OpSq:
		return "sq"
	case OpFact:
		return "fact"
	case OpAbs:
		return "abs"
	case OpMin:
		return "min"
	case OpMax:
		return "max"
	default:
		return "unknown"
	}
}

func parseExprOp(s string) (ExprOp, bool) {
	switch s {
	case "literal":
		return OpLiteral, true
	case "pi":
		return OpPi, true
	case "e":
		return OpE, true
	case "add":
		return OpAdd, true
	case "sub":
		return OpSub, true
	case "mul":
		return OpMul, true
	case "div":
		return OpDiv, true
	case "neg":
		return OpNeg, true
	case "inv":
		return OpInv, true
	case "sin":
		return OpSin, true
	case "cos":
		return OpCos, true
	case "tan":
		return OpTan, true
	case "asin":
		return OpAsin, true
	case "acos":
		return OpAcos, true
	case "atan":
		return OpAtan, true
	case "log":
		return OpLog, true
	case "ln":
		return OpLn, true
	case "exp":
		return OpExp, true
	case "10x":
		return Op10X, true
	case "pow":
		return OpPow, true
	case "sqrt":
		return OpSqrt, true
	case "sq":
		return OpSq, true
	case "fact":
		return OpFact, true
	case "abs":
		return OpAbs, true
	case "min":
		return OpMin, true
	case "max":
		return OpMax, true
	default:
		return 0, false
	}
}

// ExprNode représente un noeud dans l'arbre d'expressions de la calculatrice.
// Chaque noeud contient une opération, des enfants optionnels et un cache de résultat.
type ExprNode struct {
	op       ExprOp
	children []*ExprNode
	literal  *big.Float

	cached   *big.Float
	cacheMu  sync.Once
}

func literalNode(v *big.Float) *ExprNode {
	return &ExprNode{op: OpLiteral, literal: new(big.Float).SetPrec(defaultPrec).Copy(v)}
}

func literalFloat64Node(v float64) *ExprNode {
	bf := new(big.Float).SetPrec(defaultPrec).SetFloat64(v)
	return &ExprNode{op: OpLiteral, literal: bf}
}

func literalIntNode(v int64) *ExprNode {
	bf := new(big.Float).SetPrec(defaultPrec).SetInt64(v)
	return &ExprNode{op: OpLiteral, literal: bf}
}

func piNode() *ExprNode {
	return &ExprNode{op: OpPi}
}

func eNode() *ExprNode {
	return &ExprNode{op: OpE}
}

func unaryNode(op ExprOp, child *ExprNode) *ExprNode {
	return &ExprNode{op: op, children: []*ExprNode{child}}
}

func binaryNode(op ExprOp, left, right *ExprNode) *ExprNode {
	return &ExprNode{op: op, children: []*ExprNode{left, right}}
}

func toExprNode(v Value) *ExprNode {
	if ev, ok := v.(ExprValue); ok {
		return ev.node
	}
	return literalNode(v.BigFloat())
}

// piRational tries to express the node as (p/q)*π.
// Returns (p, q, true) if successful, with q > 0 and gcd(|p|,q) == 1.
func piRational(n *ExprNode) (p, q int64, ok bool) {
	switch n.op {
	case OpPi:
		return 1, 1, true

	case OpLiteral:
		if n.literal.Sign() == 0 {
			return 0, 1, true
		}
		return 0, 0, false

	case OpNeg:
		if len(n.children) == 1 {
			pp, qq, ok := piRational(n.children[0])
			if ok {
				return -pp, qq, true
			}
		}

	case OpMul:
		if len(n.children) == 2 {
			// Lit(k) * piExpr or piExpr * Lit(k)
			if k, ok := literalInt(n.children[0]); ok {
				if pp, qq, ok := piRational(n.children[1]); ok {
					return reduceFrac(k*pp, qq)
				}
			}
			if k, ok := literalInt(n.children[1]); ok {
				if pp, qq, ok := piRational(n.children[0]); ok {
					return reduceFrac(k*pp, qq)
				}
			}
		}

	case OpDiv:
		if len(n.children) == 2 {
			// piExpr / Lit(k)
			if k, ok := literalInt(n.children[1]); ok && k != 0 {
				if pp, qq, ok := piRational(n.children[0]); ok {
					return reduceFrac(pp, qq*k)
				}
			}
		}

	case OpAdd:
		if len(n.children) == 2 {
			p1, q1, ok1 := piRational(n.children[0])
			p2, q2, ok2 := piRational(n.children[1])
			if ok1 && ok2 {
				// p1/q1 + p2/q2 = (p1*q2 + p2*q1) / (q1*q2)
				return reduceFrac(p1*q2+p2*q1, q1*q2)
			}
		}

	case OpSub:
		if len(n.children) == 2 {
			p1, q1, ok1 := piRational(n.children[0])
			p2, q2, ok2 := piRational(n.children[1])
			if ok1 && ok2 {
				return reduceFrac(p1*q2-p2*q1, q1*q2)
			}
		}
	}
	return 0, 0, false
}

func literalInt(n *ExprNode) (int64, bool) {
	if n.op != OpLiteral || n.literal == nil {
		return 0, false
	}
	if !n.literal.IsInt() {
		return 0, false
	}
	v, acc := n.literal.Int64()
	if acc != big.Exact {
		return 0, false
	}
	return v, true
}

func reduceFrac(p, q int64) (int64, int64, bool) {
	if q == 0 {
		return 0, 0, false
	}
	if q < 0 {
		p, q = -p, -q
	}
	if p == 0 {
		return 0, 1, true
	}
	g := gcd64(abs64(p), q)
	return p / g, q / g, true
}

func gcd64(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

// degToRadNode converts a degree-valued node to radians in the graph.
// For integer degree values, it tries to create an exact pi-multiple node.
func degToRadNode(n *ExprNode) *ExprNode {
	if deg, ok := literalInt(n); ok {
		g := gcd64(abs64(deg), 180)
		p := deg / g
		q := int64(180) / g
		return ratPiNode(p, q)
	}
	return binaryNode(OpMul, n, binaryNode(OpDiv, piNode(), literalIntNode(180)))
}

// gradToRadNode converts a gradian-valued node to radians in the graph.
func gradToRadNode(n *ExprNode) *ExprNode {
	if grad, ok := literalInt(n); ok {
		g := gcd64(abs64(grad), 200)
		p := grad / g
		q := int64(200) / g
		return ratPiNode(p, q)
	}
	return binaryNode(OpMul, n, binaryNode(OpDiv, piNode(), literalIntNode(200)))
}

// ratPiNode creates a graph node representing (p/q)*π.
func ratPiNode(p, q int64) *ExprNode {
	if p == 0 {
		return literalIntNode(0)
	}
	pi := piNode()
	if q == 1 {
		if p == 1 {
			return pi
		}
		if p == -1 {
			return unaryNode(OpNeg, pi)
		}
		return binaryNode(OpMul, literalIntNode(p), pi)
	}
	if p == 1 {
		return binaryNode(OpDiv, pi, literalIntNode(q))
	}
	if p == -1 {
		return unaryNode(OpNeg, binaryNode(OpDiv, pi, literalIntNode(q)))
	}
	return binaryNode(OpDiv, binaryNode(OpMul, literalIntNode(p), pi), literalIntNode(q))
}

// radToDegNode converts a radian-valued node to degrees.
func radToDegNode(n *ExprNode) *ExprNode {
	return binaryNode(OpMul, n, binaryNode(OpDiv, literalIntNode(180), piNode()))
}

// radToGradNode converts a radian-valued node to gradians.
func radToGradNode(n *ExprNode) *ExprNode {
	return binaryNode(OpMul, n, binaryNode(OpDiv, literalIntNode(200), piNode()))
}

// evaluate computes the numerical value of the expression tree.
// It applies symbolic simplification rules for exact results where possible.
func (n *ExprNode) evaluate() *big.Float {
	n.cacheMu.Do(func() {
		n.cached = n.evalImpl()
	})
	return new(big.Float).SetPrec(defaultPrec).Copy(n.cached)
}

func (n *ExprNode) evalImpl() *big.Float {
	switch n.op {
	case OpLiteral:
		return new(big.Float).SetPrec(defaultPrec).Copy(n.literal)
	case OpPi:
		return bigPi()
	case OpE:
		return bigEConst()

	case OpAdd:
		left := n.children[0].evaluate()
		right := n.children[1].evaluate()
		return new(big.Float).SetPrec(defaultPrec).Add(left, right)
	case OpSub:
		left := n.children[0].evaluate()
		right := n.children[1].evaluate()
		return new(big.Float).SetPrec(defaultPrec).Sub(left, right)
	case OpMul:
		left := n.children[0].evaluate()
		right := n.children[1].evaluate()
		return new(big.Float).SetPrec(defaultPrec).Mul(left, right)
	case OpDiv:
		left := n.children[0].evaluate()
		right := n.children[1].evaluate()
		return new(big.Float).SetPrec(defaultPrec).Quo(left, right)

	case OpNeg:
		child := n.children[0].evaluate()
		return new(big.Float).SetPrec(defaultPrec).Neg(child)
	case OpInv:
		child := n.children[0].evaluate()
		one := new(big.Float).SetPrec(defaultPrec).SetInt64(1)
		return new(big.Float).SetPrec(defaultPrec).Quo(one, child)

	case OpSin:
		if v, ok := evalSinExact(n.children[0]); ok {
			return v
		}
		arg := n.children[0].evaluate()
		f, _ := arg.Float64()
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(math.Sin(f))

	case OpCos:
		if v, ok := evalCosExact(n.children[0]); ok {
			return v
		}
		arg := n.children[0].evaluate()
		f, _ := arg.Float64()
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(math.Cos(f))

	case OpTan:
		if v, ok := evalTanExact(n.children[0]); ok {
			return v
		}
		arg := n.children[0].evaluate()
		f, _ := arg.Float64()
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(math.Tan(f))

	case OpAsin:
		arg := n.children[0].evaluate()
		f, _ := arg.Float64()
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(math.Asin(f))
	case OpAcos:
		arg := n.children[0].evaluate()
		f, _ := arg.Float64()
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(math.Acos(f))
	case OpAtan:
		arg := n.children[0].evaluate()
		f, _ := arg.Float64()
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(math.Atan(f))

	case OpLog:
		child := n.children[0].evaluate()
		lnx, err := bigLn(child)
		if err != nil {
			return new(big.Float).SetPrec(defaultPrec)
		}
		return new(big.Float).SetPrec(defaultPrec).Quo(lnx, bigLn10())
	case OpLn:
		if v, ok := evalLnExact(n.children[0]); ok {
			return v
		}
		child := n.children[0].evaluate()
		lnx, err := bigLn(child)
		if err != nil {
			return new(big.Float).SetPrec(defaultPrec)
		}
		return lnx
	case OpExp:
		if v, ok := evalExpExact(n.children[0]); ok {
			return v
		}
		child := n.children[0].evaluate()
		return bigExp(child)
	case Op10X:
		child := n.children[0].evaluate()
		product := new(big.Float).SetPrec(defaultPrec).Mul(child, bigLn10())
		return bigExp(product)

	case OpPow:
		base := n.children[0].evaluate()
		exp := n.children[1].evaluate()
		result, err := bigPow(base, exp)
		if err != nil {
			return new(big.Float).SetPrec(defaultPrec)
		}
		return result

	case OpSqrt:
		child := n.children[0].evaluate()
		return new(big.Float).SetPrec(defaultPrec).Sqrt(child)
	case OpSq:
		child := n.children[0].evaluate()
		return new(big.Float).SetPrec(defaultPrec).Mul(child, child)

	case OpFact:
		child := n.children[0].evaluate()
		f, _ := child.Float64()
		if child.IsInt() && f <= 170 && f >= 0 {
			n := int(f)
			result := new(big.Float).SetPrec(defaultPrec).SetInt64(1)
			iBig := new(big.Float).SetPrec(defaultPrec)
			for i := 2; i <= n; i++ {
				iBig.SetInt64(int64(i))
				result.Mul(result, iBig)
			}
			return result
		}
		return new(big.Float).SetPrec(defaultPrec).SetFloat64(math.Gamma(f + 1))

	case OpAbs:
		child := n.children[0].evaluate()
		return new(big.Float).SetPrec(defaultPrec).Abs(child)

	case OpMin:
		left := n.children[0].evaluate()
		right := n.children[1].evaluate()
		if left.Cmp(right) <= 0 {
			return new(big.Float).SetPrec(defaultPrec).Copy(left)
		}
		return new(big.Float).SetPrec(defaultPrec).Copy(right)
	case OpMax:
		left := n.children[0].evaluate()
		right := n.children[1].evaluate()
		if left.Cmp(right) >= 0 {
			return new(big.Float).SetPrec(defaultPrec).Copy(left)
		}
		return new(big.Float).SetPrec(defaultPrec).Copy(right)
	}

	return new(big.Float).SetPrec(defaultPrec)
}

func evalSinExact(argNode *ExprNode) (*big.Float, bool) {
	p, q, ok := piRational(argNode)
	if !ok {
		return nil, false
	}
	val, ok := sinPiRational(p, q)
	if !ok {
		return nil, false
	}
	return new(big.Float).SetPrec(defaultPrec).SetFloat64(val), true
}

func evalCosExact(argNode *ExprNode) (*big.Float, bool) {
	p, q, ok := piRational(argNode)
	if !ok {
		return nil, false
	}
	val, ok := cosPiRational(p, q)
	if !ok {
		return nil, false
	}
	return new(big.Float).SetPrec(defaultPrec).SetFloat64(val), true
}

func evalTanExact(argNode *ExprNode) (*big.Float, bool) {
	p, q, ok := piRational(argNode)
	if !ok {
		return nil, false
	}
	s, sok := sinPiRational(p, q)
	c, cok := cosPiRational(p, q)
	if !sok || !cok || c == 0 {
		return nil, false
	}
	if s == 0 {
		return new(big.Float).SetPrec(defaultPrec), true
	}
	return new(big.Float).SetPrec(defaultPrec).SetFloat64(s / c), true
}

// sinPiRational returns the exact value of sin((p/q)*π) for known fractions.
// p/q should be in reduced form with q > 0, p in [0, 2q).
func sinPiRational(p, q int64) (float64, bool) {
	// Normalize to [0, 2) -- p/q represents a fraction of π
	// We need p*π / q, normalized modulo 2π
	// So p/q mod 2 (in units of π)
	p = p % (2 * q)
	if p < 0 {
		p += 2 * q
	}
	g := gcd64(abs64(p), q)
	p = p / g
	q = q / g

	// Use symmetry: sin((p/q)π) with p/q > 1 → sin(π + x) = -sin(x)
	neg := false
	if 2*p > 2*q { // p/q > 1
		p = 2*q - p // reflect: sin((2-p/q)π) = -sin((p/q-1)π)... actually sin(2π - x) = -sin(x)
		neg = true
	}
	// Now p/q in [0, 1]
	// sin((p/q)π) with p/q in (1/2, 1] → sin(π-x) = sin(x)
	if 2*p > q { // p/q > 1/2
		p = q - p // sin((1 - p/q)π) = sin((p/q)π) -- no, sin(π-x) = sin(x)
	}
	// Now p/q in [0, 1/2]

	g = gcd64(abs64(p), q)
	p = p / g
	q = q / g

	var val float64
	found := true
	switch {
	case p == 0:
		val = 0
	case p == 1 && q == 6:
		val = 0.5
	case p == 1 && q == 4:
		val = math.Sqrt2 / 2
	case p == 1 && q == 3:
		val = math.Sqrt(3) / 2
	case p == 1 && q == 2:
		val = 1
	default:
		found = false
	}
	if !found {
		return 0, false
	}
	if neg {
		val = -val
	}
	return val, true
}

// cosPiRational returns the exact value of cos((p/q)*π) for known fractions.
func cosPiRational(p, q int64) (float64, bool) {
	// cos(x) = sin(π/2 + x) = sin((1/2 + p/q)π)
	// But simpler: cos((p/q)π) = sin((1/2 - p/q)π + π/2) ... let's just enumerate directly
	// Normalize to [0, 2)
	p = p % (2 * q)
	if p < 0 {
		p += 2 * q
	}
	g := gcd64(abs64(p), q)
	p = p / g
	q = q / g

	// cos(x) = cos(-x), cos(2π-x) = cos(x)
	neg := false
	if 2*p > 2*q { // p/q > 1
		p = 2*q - p
		// cos(2π - x) = cos(x)
	}
	g = gcd64(abs64(p), q)
	p = p / g
	q = q / g
	// Now p/q in [0, 1]
	// cos(π - x) = -cos(x) for x in (0, π)
	if 2*p > q { // p/q > 1/2
		p = q - p
		neg = true
	}
	g = gcd64(abs64(p), q)
	p = p / g
	q = q / g
	// Now p/q in [0, 1/2]

	var val float64
	found := true
	switch {
	case p == 0:
		val = 1
	case p == 1 && q == 6:
		val = math.Sqrt(3) / 2
	case p == 1 && q == 4:
		val = math.Sqrt2 / 2
	case p == 1 && q == 3:
		val = 0.5
	case p == 1 && q == 2:
		val = 0
	default:
		found = false
	}
	if !found {
		return 0, false
	}
	if neg {
		val = -val
	}
	return val, true
}

// evalLnExact tries to compute exact values for ln.
func evalLnExact(argNode *ExprNode) (*big.Float, bool) {
	// ln(e) = 1
	if argNode.op == OpE {
		return new(big.Float).SetPrec(defaultPrec).SetInt64(1), true
	}
	// ln(1) = 0
	if argNode.op == OpLiteral && argNode.literal != nil {
		if argNode.literal.Cmp(new(big.Float).SetInt64(1)) == 0 {
			return new(big.Float).SetPrec(defaultPrec).SetInt64(0), true
		}
	}
	// ln(exp(x)) = x
	if argNode.op == OpExp && len(argNode.children) == 1 {
		return argNode.children[0].evaluate(), true
	}
	return nil, false
}

// evalExpExact tries to compute exact values for exp.
func evalExpExact(argNode *ExprNode) (*big.Float, bool) {
	// exp(0) = 1
	if argNode.op == OpLiteral && argNode.literal != nil && argNode.literal.Sign() == 0 {
		return new(big.Float).SetPrec(defaultPrec).SetInt64(1), true
	}
	// exp(1) = e
	if argNode.op == OpLiteral && argNode.literal != nil {
		if argNode.literal.Cmp(new(big.Float).SetInt64(1)) == 0 {
			return bigEConst(), true
		}
	}
	// exp(ln(x)) = x
	if argNode.op == OpLn && len(argNode.children) == 1 {
		return argNode.children[0].evaluate(), true
	}
	return nil, false
}

