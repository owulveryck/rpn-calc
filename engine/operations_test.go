package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
)

func parseState(t *testing.T, jsonStr string) State {
	t.Helper()
	var s State
	if err := json.Unmarshal([]byte(jsonStr), &s); err != nil {
		t.Fatalf("failed to parse state: %v", err)
	}
	return s
}

func calcWithStack(vals ...float64) *Calculator {
	c := NewCalculator()
	for _, v := range vals {
		c.Enter(fmt.Sprintf("%g", v))
	}
	return c
}

func TestArithmetic(t *testing.T) {
	// y is pushed first (level 2), x is on top (level 1)
	// RPN: y op x
	tests := []struct {
		name string
		y, x float64
		op   string
		want float64
	}{
		{"add", 2, 3, "+", 5},
		{"sub", 10, 3, "-", 7},
		{"mul", 4, 5, "*", 20},
		{"div", 10, 4, "/", 2.5},
		{"add negatives", -2, -3, "+", -5},
		{"sub negative result", 3, 10, "-", -7},
		{"mul by zero", 5, 0, "*", 0},
		{"div large", 1000000, 3, "/", 1000000.0 / 3.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calcWithStack(tt.y, tt.x)
			s := parseState(t, c.Execute(tt.op))
			if s.Error != "" {
				t.Fatalf("unexpected error: %s", s.Error)
			}
			if len(s.Stack) != 1 {
				t.Fatalf("stack len = %d, want 1", len(s.Stack))
			}
			got, _ := ParseValue(s.Stack[0])
			if math.Abs(got.Float64()-tt.want) > 1e-10 {
				t.Errorf("%g %s %g = %v, want %g", tt.y, tt.op, tt.x, s.Stack[0], tt.want)
			}
		})
	}
}

func TestDivisionByZero(t *testing.T) {
	c := calcWithStack(10, 0)
	s := parseState(t, c.Execute("/"))
	if s.Error == "" {
		t.Error("expected error for division by zero")
	}
	if s.Depth != 2 {
		t.Errorf("stack should be unchanged, depth = %d", s.Depth)
	}
}

func TestNeg(t *testing.T) {
	c := calcWithStack(42)
	s := parseState(t, c.Execute("NEG"))
	if s.Stack[0] != "-42" {
		t.Errorf("NEG(42) = %s, want -42", s.Stack[0])
	}
	// double negation
	s = parseState(t, c.Execute("NEG"))
	if s.Stack[0] != "42" {
		t.Errorf("NEG(NEG(42)) = %s, want 42", s.Stack[0])
	}
}

func TestInv(t *testing.T) {
	c := calcWithStack(4)
	s := parseState(t, c.Execute("INV"))
	got, _ := ParseValue(s.Stack[0])
	if got.Float64() != 0.25 {
		t.Errorf("INV(4) = %v, want 0.25", s.Stack[0])
	}
}

func TestInvZero(t *testing.T) {
	c := calcWithStack(0)
	s := parseState(t, c.Execute("INV"))
	if s.Error == "" {
		t.Error("expected error for INV(0)")
	}
}

func TestTrigDeg(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("DEG")

	c.Enter("90")
	s := parseState(t, c.Execute("SIN"))
	got, _ := ParseValue(s.Stack[0])
	if math.Abs(got.Float64()-1) > 1e-10 {
		t.Errorf("sin(90 DEG) = %v, want 1", s.Stack[0])
	}

	c.Enter("0")
	s = parseState(t, c.Execute("COS"))
	got, _ = ParseValue(s.Stack[0])
	if math.Abs(got.Float64()-1) > 1e-10 {
		t.Errorf("cos(0 DEG) = %v, want 1", s.Stack[0])
	}

	c.Enter("45")
	s = parseState(t, c.Execute("TAN"))
	got, _ = ParseValue(s.Stack[0])
	if math.Abs(got.Float64()-1) > 1e-10 {
		t.Errorf("tan(45 DEG) = %v, want 1", s.Stack[0])
	}
}

func TestTrigRad(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("RAD")

	c.Enter(fmt.Sprintf("%g", math.Pi/2))
	s := parseState(t, c.Execute("SIN"))
	got, _ := ParseValue(s.Stack[0])
	if math.Abs(got.Float64()-1) > 1e-10 {
		t.Errorf("sin(pi/2 RAD) = %v, want 1", s.Stack[0])
	}
}

func TestTrigGrad(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("GRAD")

	c.Enter("100")
	s := parseState(t, c.Execute("SIN"))
	got, _ := ParseValue(s.Stack[0])
	if math.Abs(got.Float64()-1) > 1e-10 {
		t.Errorf("sin(100 GRAD) = %v, want 1", s.Stack[0])
	}
}

func TestInverseTrigRoundtrip(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("DEG")

	c.Enter("30")
	c.Execute("SIN")
	s := parseState(t, c.Execute("ASIN"))
	got, _ := ParseValue(s.Stack[0])
	if math.Abs(got.Float64()-30) > 1e-10 {
		t.Errorf("asin(sin(30)) = %v, want 30", s.Stack[0])
	}
}

func TestLogLn(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		op   string
		want float64
	}{
		{"log(100)", 100, "LOG", 2},
		{"log(1)", 1, "LOG", 0},
		{"ln(e)", math.E, "LN", 1},
		{"ln(1)", 1, "LN", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calcWithStack(tt.val)
			s := parseState(t, c.Execute(tt.op))
			if s.Error != "" {
				t.Fatalf("unexpected error: %s", s.Error)
			}
			got, _ := ParseValue(s.Stack[0])
			if math.Abs(got.Float64()-tt.want) > 1e-10 {
				t.Errorf("%s(%g) = %v, want %g", tt.op, tt.val, s.Stack[0], tt.want)
			}
		})
	}
}

func TestLogNonPositive(t *testing.T) {
	for _, v := range []float64{0, -1} {
		c := calcWithStack(v)
		s := parseState(t, c.Execute("LOG"))
		if s.Error == "" {
			t.Errorf("expected error for LOG(%g)", v)
		}
	}
}

func TestExpAnd10X(t *testing.T) {
	c := calcWithStack(0)
	s := parseState(t, c.Execute("EXP"))
	got, _ := ParseValue(s.Stack[0])
	if math.Abs(got.Float64()-1) > 1e-10 {
		t.Errorf("exp(0) = %v, want 1", s.Stack[0])
	}

	c = calcWithStack(2)
	s = parseState(t, c.Execute("10^X"))
	got, _ = ParseValue(s.Stack[0])
	if math.Abs(got.Float64()-100) > 1e-10 {
		t.Errorf("10^2 = %v, want 100", s.Stack[0])
	}
}

func TestPow(t *testing.T) {
	c := calcWithStack(2, 10)
	s := parseState(t, c.Execute("Y^X"))
	got, _ := ParseValue(s.Stack[0])
	if got.Float64() != 1024 {
		t.Errorf("2^10 = %v, want 1024", s.Stack[0])
	}
}

func TestSqrt(t *testing.T) {
	c := calcWithStack(144)
	s := parseState(t, c.Execute("SQRT"))
	got, _ := ParseValue(s.Stack[0])
	if got.Float64() != 12 {
		t.Errorf("sqrt(144) = %v, want 12", s.Stack[0])
	}
}

func TestSqrtNegative(t *testing.T) {
	c := calcWithStack(-1)
	s := parseState(t, c.Execute("SQRT"))
	if s.Error == "" {
		t.Error("expected error for sqrt(-1)")
	}
}

func TestSq(t *testing.T) {
	c := calcWithStack(7)
	s := parseState(t, c.Execute("SQ"))
	got, _ := ParseValue(s.Stack[0])
	if got.Float64() != 49 {
		t.Errorf("7^2 = %v, want 49", s.Stack[0])
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		want float64
	}{
		{"0!", 0, 1},
		{"1!", 1, 1},
		{"5!", 5, 120},
		{"10!", 10, 3628800},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calcWithStack(tt.val)
			s := parseState(t, c.Execute("FACT"))
			got, _ := ParseValue(s.Stack[0])
			if math.Abs(got.Float64()-tt.want) > 1e-6 {
				t.Errorf("%g! = %v, want %g", tt.val, s.Stack[0], tt.want)
			}
		})
	}
}

func TestFactorialNonInteger(t *testing.T) {
	c := calcWithStack(5.5)
	s := parseState(t, c.Execute("FACT"))
	if s.Error != "" {
		t.Fatalf("unexpected error: %s", s.Error)
	}
	got, _ := ParseValue(s.Stack[0])
	expected := math.Gamma(6.5)
	if math.Abs(got.Float64()-expected) > 1e-6 {
		t.Errorf("5.5! = %v, want %g", s.Stack[0], expected)
	}
}

func TestFactorialNegative(t *testing.T) {
	c := calcWithStack(-1)
	s := parseState(t, c.Execute("FACT"))
	if s.Error == "" {
		t.Error("expected error for factorial(-1)")
	}
}

func TestConstants(t *testing.T) {
	c := NewCalculator()
	s := parseState(t, c.Execute("PI"))
	got, _ := ParseValue(s.Stack[0])
	if math.Abs(got.Float64()-math.Pi) > 1e-15 {
		t.Errorf("PI = %v, want %g", s.Stack[0], math.Pi)
	}

	s = parseState(t, c.Execute("E"))
	got, _ = ParseValue(s.Stack[len(s.Stack)-1])
	if math.Abs(got.Float64()-math.E) > 1e-15 {
		t.Errorf("E = %v, want %g", s.Stack[len(s.Stack)-1], math.E)
	}
}

func TestAbs(t *testing.T) {
	c := calcWithStack(-42)
	s := parseState(t, c.Execute("ABS"))
	if s.Stack[0] != "42" {
		t.Errorf("ABS(-42) = %s, want 42", s.Stack[0])
	}
}

func TestMinMax(t *testing.T) {
	c := calcWithStack(3, 7)
	s := parseState(t, c.Execute("MIN"))
	if s.Stack[0] != "3" {
		t.Errorf("MIN(3,7) = %s, want 3", s.Stack[0])
	}

	c = calcWithStack(3, 7)
	s = parseState(t, c.Execute("MAX"))
	if s.Stack[0] != "7" {
		t.Errorf("MAX(3,7) = %s, want 7", s.Stack[0])
	}
}

func TestDepthOp(t *testing.T) {
	c := calcWithStack(1, 2, 3)
	s := parseState(t, c.Execute("DEPTH"))
	if s.Stack[len(s.Stack)-1] != "3" {
		t.Errorf("DEPTH of 3-element stack = %s, want 3", s.Stack[len(s.Stack)-1])
	}
	if s.Depth != 4 {
		t.Errorf("after DEPTH, stack depth = %d, want 4", s.Depth)
	}
}

func TestUnknownOp(t *testing.T) {
	c := NewCalculator()
	s := parseState(t, c.Execute("FOOBAR"))
	if s.Error == "" {
		t.Error("expected error for unknown operation")
	}
}

func TestOperationUnderflow(t *testing.T) {
	ops := []string{"+", "-", "*", "/", "Y^X", "MIN", "MAX"}
	for _, op := range ops {
		t.Run(op+"_empty", func(t *testing.T) {
			c := NewCalculator()
			s := parseState(t, c.Execute(op))
			if s.Error == "" {
				t.Errorf("expected underflow for %s on empty stack", op)
			}
		})
		t.Run(op+"_one", func(t *testing.T) {
			c := calcWithStack(42)
			s := parseState(t, c.Execute(op))
			if s.Error == "" {
				t.Errorf("expected underflow for %s with 1 element", op)
			}
			if s.Depth != 1 {
				t.Errorf("stack should be unchanged, depth = %d", s.Depth)
			}
		})
	}
}
