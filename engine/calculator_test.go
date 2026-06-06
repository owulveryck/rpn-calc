package engine

import (
	"encoding/json"
	"math"
	"testing"
)

func TestCalculatorSequence(t *testing.T) {
	tests := []struct {
		name     string
		sequence []string // "ENTER:x" or "OP:x"
		want     []string
	}{
		{
			"2 + 3 = 5",
			[]string{"ENTER:2", "ENTER:3", "OP:+"},
			[]string{"5"},
		},
		{
			"2 * (3 + 4) = 14",
			[]string{"ENTER:2", "ENTER:3", "ENTER:4", "OP:+", "OP:*"},
			[]string{"14"},
		},
		{
			"10 - 3 = 7",
			[]string{"ENTER:10", "ENTER:3", "OP:-"},
			[]string{"7"},
		},
		{
			"stack accumulation",
			[]string{"ENTER:1", "ENTER:2", "ENTER:3"},
			[]string{"1", "2", "3"},
		},
		{
			"swap then add",
			[]string{"ENTER:5", "ENTER:3", "OP:SWAP", "OP:-"},
			[]string{"-2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCalculator()
			var stateJSON string
			for _, cmd := range tt.sequence {
				if len(cmd) > 6 && cmd[:6] == "ENTER:" {
					stateJSON = c.Enter(cmd[6:])
				} else if len(cmd) > 3 && cmd[:3] == "OP:" {
					stateJSON = c.Execute(cmd[3:])
				}
			}
			var s State
			json.Unmarshal([]byte(stateJSON), &s)
			if s.Error != "" {
				t.Fatalf("unexpected error: %s", s.Error)
			}
			if len(s.Stack) != len(tt.want) {
				t.Fatalf("stack = %v, want %v", s.Stack, tt.want)
			}
			for i, w := range tt.want {
				if s.Stack[i] != w {
					t.Errorf("stack[%d] = %s, want %s", i, s.Stack[i], w)
				}
			}
		})
	}
}

func TestUndo(t *testing.T) {
	c := NewCalculator()
	c.Enter("10")
	c.Enter("5")
	c.Execute("+")

	s := parseState(t, c.Undo())
	if s.Depth != 2 || s.Stack[0] != "10" || s.Stack[1] != "5" {
		t.Errorf("after undo: stack = %v, want [10, 5]", s.Stack)
	}
}

func TestUndoMultiple(t *testing.T) {
	c := NewCalculator()
	c.Enter("1")
	c.Enter("2")
	c.Execute("+")
	c.Enter("3")
	c.Execute("*")

	c.Undo()
	c.Undo()
	s := parseState(t, c.Undo())
	if s.Depth != 2 || s.Stack[0] != "1" || s.Stack[1] != "2" {
		t.Errorf("after 3 undos: stack = %v, want [1, 2]", s.Stack)
	}
}

func TestUndoEmpty(t *testing.T) {
	c := NewCalculator()
	s := parseState(t, c.Undo())
	if s.Error == "" {
		t.Error("expected error for undo on empty undo stack")
	}
}

func TestLast(t *testing.T) {
	c := NewCalculator()
	c.Enter("2")
	c.Enter("3")
	c.Execute("+")

	s := parseState(t, c.Last())
	if s.Error != "" {
		t.Fatalf("unexpected error: %s", s.Error)
	}
	// stack should be: [5, 2, 3]
	if s.Depth != 3 {
		t.Fatalf("depth = %d, want 3", s.Depth)
	}
	if s.Stack[0] != "5" || s.Stack[1] != "2" || s.Stack[2] != "3" {
		t.Errorf("after LAST: stack = %v, want [5, 2, 3]", s.Stack)
	}
}

func TestLastEmpty(t *testing.T) {
	c := NewCalculator()
	s := parseState(t, c.Last())
	if s.Error == "" {
		t.Error("expected error for LAST with no previous operation")
	}
}

func TestAngleModeChange(t *testing.T) {
	c := NewCalculator()
	c.Enter("42")

	s := parseState(t, c.SetAngleMode("RAD"))
	if s.AngleMode != "RAD" {
		t.Errorf("angle mode = %s, want RAD", s.AngleMode)
	}
	if s.Stack[0] != "42" {
		t.Error("changing angle mode should not modify the stack")
	}

	s = parseState(t, c.SetAngleMode("GRAD"))
	if s.AngleMode != "GRAD" {
		t.Errorf("angle mode = %s, want GRAD", s.AngleMode)
	}
}

func TestBaseModeChange(t *testing.T) {
	c := NewCalculator()
	c.Enter("255")

	s := parseState(t, c.SetBaseMode("HEX"))
	if s.BaseMode != "HEX" {
		t.Errorf("base mode = %s, want HEX", s.BaseMode)
	}
	if s.Stack[0] != "#FF" {
		t.Errorf("255 in HEX = %s, want #FF", s.Stack[0])
	}

	s = parseState(t, c.SetBaseMode("DEC"))
	if s.Stack[0] != "255" {
		t.Errorf("back to DEC = %s, want 255", s.Stack[0])
	}
}

func TestInvalidInput(t *testing.T) {
	c := NewCalculator()
	c.Enter("42")
	s := parseState(t, c.Enter("abc"))
	if s.Error == "" {
		t.Error("expected error for invalid input")
	}
	if s.Depth != 1 || s.Stack[0] != "42" {
		t.Error("stack should be unchanged after invalid input")
	}
}

func TestOperationNoSideEffectOnError(t *testing.T) {
	c := NewCalculator()
	c.Enter("42")
	s := parseState(t, c.Execute("/"))
	if s.Error == "" {
		t.Fatal("expected underflow error")
	}
	if s.Depth != 1 || s.Stack[0] != "42" {
		t.Errorf("stack should be unchanged, got %v", s.Stack)
	}
}

func TestQuadraticFormula(t *testing.T) {
	// Solve x^2 - 5x + 6 = 0 -> x = 2 or x = 3
	// Using (-b + sqrt(b^2 - 4ac)) / 2a
	// a=1, b=-5, c=6
	c := NewCalculator()
	c.SetAngleMode("RAD")

	// Compute discriminant: b^2 - 4ac
	c.Enter("-5")
	c.Execute("SQ") // 25
	c.Enter("4")
	c.Enter("1")
	c.Execute("*") // 4
	c.Enter("6")
	c.Execute("*") // 24
	c.Execute("-") // 25 - 24 = 1
	c.Execute("SQRT") // 1

	// (-b + sqrt) / 2a = (5 + 1) / 2
	c.Enter("-5")
	c.Execute("NEG") // 5
	c.Execute("SWAP")
	c.Execute("DUP") // stack: [5, 1, 1]
	c.Execute("ROT") // stack: [1, 1, 5]
	c.Execute("ROT") // stack: [1, 5, 1]
	c.Execute("+") // stack: [1, 6]
	c.Enter("2")
	c.Execute("/") // stack: [1, 3] -> x1 = 3

	s := parseState(t, c.GetStateJSON())
	got1, _ := ParseValue(s.Stack[len(s.Stack)-1])
	if math.Abs(got1.Float64()-3) > 1e-10 {
		t.Errorf("x1 = %v, want 3", s.Stack[len(s.Stack)-1])
	}

	// Now compute x2: (-b - sqrt) / 2a
	c.Execute("SWAP") // bring sqrt(discriminant) to top
	c.Enter("-5")
	c.Execute("NEG") // 5
	c.Execute("SWAP")
	c.Execute("-") // 5 - 1 = 4
	c.Enter("2")
	s = parseState(t, c.Execute("/")) // 4/2 = 2

	got2, _ := ParseValue(s.Stack[len(s.Stack)-1])
	if math.Abs(got2.Float64()-2) > 1e-10 {
		t.Errorf("x2 = %v, want 2", s.Stack[len(s.Stack)-1])
	}
}

func TestHistoryBasicArithmetic(t *testing.T) {
	c := NewCalculator()
	c.Enter("2")
	c.Enter("3")
	s := parseState(t, c.Execute("+"))
	if len(s.History) < 3 {
		t.Fatalf("history len = %d, want >= 3", len(s.History))
	}
	last := s.History[len(s.History)-1]
	if last.Expression != "2 + 3 = 5" {
		t.Errorf("history = %q, want %q", last.Expression, "2 + 3 = 5")
	}
}

func TestHistoryUnary(t *testing.T) {
	c := NewCalculator()
	c.SetAngleMode("DEG")
	c.Enter("90")
	s := parseState(t, c.Execute("SIN"))
	last := s.History[len(s.History)-1]
	if last.Expression != "sin(90) = 1" {
		t.Errorf("history = %q, want %q", last.Expression, "sin(90) = 1")
	}
}

func TestHistoryFactorial(t *testing.T) {
	c := NewCalculator()
	c.Enter("5")
	s := parseState(t, c.Execute("FACT"))
	last := s.History[len(s.History)-1]
	if last.Expression != "5! = 120" {
		t.Errorf("history = %q, want %q", last.Expression, "5! = 120")
	}
}

func TestHistorySq(t *testing.T) {
	c := NewCalculator()
	c.Enter("7")
	s := parseState(t, c.Execute("SQ"))
	last := s.History[len(s.History)-1]
	if last.Expression != "7² = 49" {
		t.Errorf("history = %q, want %q", last.Expression, "7² = 49")
	}
}

func TestHistoryInv(t *testing.T) {
	c := NewCalculator()
	c.Enter("4")
	s := parseState(t, c.Execute("INV"))
	last := s.History[len(s.History)-1]
	if last.Expression != "1/4 = 0.25" {
		t.Errorf("history = %q, want %q", last.Expression, "1/4 = 0.25")
	}
}

func TestHistoryConstants(t *testing.T) {
	c := NewCalculator()
	s := parseState(t, c.Execute("PI"))
	if len(s.History) != 1 {
		t.Fatalf("history len = %d, want 1", len(s.History))
	}
	if s.History[0].Expression[:4] != "PI =" {
		t.Errorf("history = %q, want prefix 'PI ='", s.History[0].Expression)
	}
}

func TestHistoryStackOps(t *testing.T) {
	c := NewCalculator()
	c.Enter("1")
	c.Enter("2")
	s := parseState(t, c.Execute("SWAP"))
	last := s.History[len(s.History)-1]
	if last.Expression != "SWAP" {
		t.Errorf("history = %q, want %q", last.Expression, "SWAP")
	}
}

func TestHistoryEnter(t *testing.T) {
	c := NewCalculator()
	s := parseState(t, c.Enter("42"))
	if len(s.History) != 1 {
		t.Fatalf("history len = %d, want 1", len(s.History))
	}
	if s.History[0].Expression != "ENTER 42" {
		t.Errorf("history = %q, want %q", s.History[0].Expression, "ENTER 42")
	}
}

func TestHistoryMaxSize(t *testing.T) {
	c := NewCalculator()
	for i := 0; i < 60; i++ {
		c.Enter("1")
	}
	s := c.GetState()
	if len(s.History) != 50 {
		t.Errorf("history len = %d, want 50 (maxHistory)", len(s.History))
	}
}

func TestHistoryInState(t *testing.T) {
	c := NewCalculator()
	c.Enter("10")
	c.Enter("5")
	c.Execute("+")
	s := c.GetState()
	if len(s.History) == 0 {
		t.Fatal("history should not be empty")
	}
	if s.History == nil {
		t.Error("history should never be nil in state")
	}
}

func TestHistoryNotAddedOnError(t *testing.T) {
	c := NewCalculator()
	c.Enter("42")
	before := len(c.GetState().History)
	c.Execute("/") // underflow
	after := len(c.GetState().History)
	if after != before {
		t.Errorf("history should not grow on error: before=%d, after=%d", before, after)
	}
}
