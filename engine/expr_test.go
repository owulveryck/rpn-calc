package engine

import (
	"math"
	"testing"
)

func TestSinPiExact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("RAD")
	c.Execute("PI")
	s := parseState(t, c.Execute("SIN"))
	if s.Stack[0] != "0" {
		t.Errorf("sin(π) = %s, want 0", s.Stack[0])
	}
}

func TestCosPiExact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("RAD")
	c.Execute("PI")
	s := parseState(t, c.Execute("COS"))
	if s.Stack[0] != "-1" {
		t.Errorf("cos(π) = %s, want -1", s.Stack[0])
	}
}

func TestSinPiOver2Exact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("RAD")
	c.Execute("PI")
	c.Enter("2")
	c.Execute("/")
	s := parseState(t, c.Execute("SIN"))
	if s.Stack[0] != "1" {
		t.Errorf("sin(π/2) = %s, want 1", s.Stack[0])
	}
}

func TestCosPiOver2Exact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("RAD")
	c.Execute("PI")
	c.Enter("2")
	c.Execute("/")
	s := parseState(t, c.Execute("COS"))
	if s.Stack[0] != "0" {
		t.Errorf("cos(π/2) = %s, want 0", s.Stack[0])
	}
}

func TestSin2PiExact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("RAD")
	c.Enter("2")
	c.Execute("PI")
	c.Execute("*")
	s := parseState(t, c.Execute("SIN"))
	if s.Stack[0] != "0" {
		t.Errorf("sin(2π) = %s, want 0", s.Stack[0])
	}
}

func TestCos2PiExact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("RAD")
	c.Enter("2")
	c.Execute("PI")
	c.Execute("*")
	s := parseState(t, c.Execute("COS"))
	if s.Stack[0] != "1" {
		t.Errorf("cos(2π) = %s, want 1", s.Stack[0])
	}
}

func TestTanPiExact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("RAD")
	c.Execute("PI")
	s := parseState(t, c.Execute("TAN"))
	if s.Stack[0] != "0" {
		t.Errorf("tan(π) = %s, want 0", s.Stack[0])
	}
}

func TestSinPiOver6Exact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("RAD")
	c.Execute("PI")
	c.Enter("6")
	c.Execute("/")
	s := parseState(t, c.Execute("SIN"))
	got, _ := ParseValue(s.Stack[0])
	if math.Abs(got.Float64()-0.5) > 1e-15 {
		t.Errorf("sin(π/6) = %s, want 0.5", s.Stack[0])
	}
}

func TestCosPiOver3Exact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("RAD")
	c.Execute("PI")
	c.Enter("3")
	c.Execute("/")
	s := parseState(t, c.Execute("COS"))
	got, _ := ParseValue(s.Stack[0])
	if math.Abs(got.Float64()-0.5) > 1e-15 {
		t.Errorf("cos(π/3) = %s, want 0.5", s.Stack[0])
	}
}

func TestSin180DegExact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("DEG")
	c.Enter("180")
	s := parseState(t, c.Execute("SIN"))
	if s.Stack[0] != "0" {
		t.Errorf("sin(180°) = %s, want 0", s.Stack[0])
	}
}

func TestCos180DegExact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("DEG")
	c.Enter("180")
	s := parseState(t, c.Execute("COS"))
	if s.Stack[0] != "-1" {
		t.Errorf("cos(180°) = %s, want -1", s.Stack[0])
	}
}

func TestSin90DegExact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("DEG")
	c.Enter("90")
	s := parseState(t, c.Execute("SIN"))
	if s.Stack[0] != "1" {
		t.Errorf("sin(90°) = %s, want 1", s.Stack[0])
	}
}

func TestCos90DegExact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("DEG")
	c.Enter("90")
	s := parseState(t, c.Execute("COS"))
	if s.Stack[0] != "0" {
		t.Errorf("cos(90°) = %s, want 0", s.Stack[0])
	}
}

func TestSin360DegExact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("DEG")
	c.Enter("360")
	s := parseState(t, c.Execute("SIN"))
	if s.Stack[0] != "0" {
		t.Errorf("sin(360°) = %s, want 0", s.Stack[0])
	}
}

func TestSin30DegExact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("DEG")
	c.Enter("30")
	s := parseState(t, c.Execute("SIN"))
	got, _ := ParseValue(s.Stack[0])
	if math.Abs(got.Float64()-0.5) > 1e-15 {
		t.Errorf("sin(30°) = %s, want 0.5", s.Stack[0])
	}
}

func TestSin100GradExact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("GRAD")
	c.Enter("100")
	s := parseState(t, c.Execute("SIN"))
	if s.Stack[0] != "1" {
		t.Errorf("sin(100 grad) = %s, want 1", s.Stack[0])
	}
}

func TestSin200GradExact(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("GRAD")
	c.Enter("200")
	s := parseState(t, c.Execute("SIN"))
	if s.Stack[0] != "0" {
		t.Errorf("sin(200 grad) = %s, want 0", s.Stack[0])
	}
}

func TestLnEExact(t *testing.T) {
	c := NewCalculator()
	c.Execute("E")
	s := parseState(t, c.Execute("LN"))
	if s.Stack[0] != "1" {
		t.Errorf("ln(e) = %s, want 1", s.Stack[0])
	}
}

func TestLn1Exact(t *testing.T) {
	c := NewCalculator()
	c.Enter("1")
	s := parseState(t, c.Execute("LN"))
	if s.Stack[0] != "0" {
		t.Errorf("ln(1) = %s, want 0", s.Stack[0])
	}
}

func TestExp0Exact(t *testing.T) {
	c := NewCalculator()
	c.Enter("0")
	s := parseState(t, c.Execute("EXP"))
	if s.Stack[0] != "1" {
		t.Errorf("exp(0) = %s, want 1", s.Stack[0])
	}
}

func TestExpLnRoundtripExact(t *testing.T) {
	c := NewCalculator()
	c.Enter("42")
	c.Execute("LN")
	s := parseState(t, c.Execute("EXP"))
	if s.Stack[0] != "42" {
		t.Errorf("exp(ln(42)) = %s, want 42", s.Stack[0])
	}
}

func TestLnExpRoundtripExact(t *testing.T) {
	c := NewCalculator()
	c.Enter("7")
	c.Execute("EXP")
	s := parseState(t, c.Execute("LN"))
	if s.Stack[0] != "7" {
		t.Errorf("ln(exp(7)) = %s, want 7", s.Stack[0])
	}
}

func TestGraphSerializeRoundTrip(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("RAD")
	c.Execute("PI")

	data := c.Serialize()
	c2 := NewCalculator()
	if err := c2.Restore(data); err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	// The restored value should still be recognized as PI
	s := parseState(t, c2.Execute("SIN"))
	if s.Stack[0] != "0" {
		t.Errorf("sin(restored PI) = %s, want 0 — graph not preserved", s.Stack[0])
	}
}

func TestGraphSerializeComplexExpr(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("RAD")
	c.Execute("PI")
	c.Enter("2")
	c.Execute("/")

	data := c.Serialize()
	c2 := NewCalculator()
	if err := c2.Restore(data); err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	// sin(π/2) should still be exact after restore
	s := parseState(t, c2.Execute("SIN"))
	if s.Stack[0] != "1" {
		t.Errorf("sin(restored π/2) = %s, want 1 — graph not preserved", s.Stack[0])
	}
}

func TestBackwardCompatRestore(t *testing.T) {
	// Simulate old "real" format
	oldData := `{"stack":[{"type":"real","value":"3.1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679"}],"angleMode":"RAD","baseMode":"DEC"}`
	c := NewCalculator()
	if err := c.Restore(oldData); err != nil {
		t.Fatalf("backward compat restore failed: %v", err)
	}
	s := c.GetState()
	if s.Depth != 1 {
		t.Fatalf("expected depth 1, got %d", s.Depth)
	}
	got, _ := ParseValue(s.Stack[0])
	if math.Abs(got.Float64()-math.Pi) > 1e-15 {
		t.Errorf("restored value = %s, want π", s.Stack[0])
	}
}

func TestPiRationalDetection(t *testing.T) {
	tests := []struct {
		name     string
		node     *ExprNode
		wantP    int64
		wantQ    int64
		wantOk   bool
	}{
		{"pi", piNode(), 1, 1, true},
		{"pi/2", binaryNode(OpDiv, piNode(), literalIntNode(2)), 1, 2, true},
		{"pi/6", binaryNode(OpDiv, piNode(), literalIntNode(6)), 1, 6, true},
		{"2*pi", binaryNode(OpMul, literalIntNode(2), piNode()), 2, 1, true},
		{"neg pi", unaryNode(OpNeg, piNode()), -1, 1, true},
		{"literal 0", literalIntNode(0), 0, 1, true},
		{"literal 5", literalIntNode(5), 0, 0, false},
		{"3*pi/4", binaryNode(OpDiv, binaryNode(OpMul, literalIntNode(3), piNode()), literalIntNode(4)), 3, 4, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, q, ok := piRational(tt.node)
			if ok != tt.wantOk {
				t.Errorf("piRational ok = %v, want %v", ok, tt.wantOk)
				return
			}
			if ok && (p != tt.wantP || q != tt.wantQ) {
				t.Errorf("piRational = (%d/%d), want (%d/%d)", p, q, tt.wantP, tt.wantQ)
			}
		})
	}
}
