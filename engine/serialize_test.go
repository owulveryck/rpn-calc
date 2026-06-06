package engine

import (
	"math"
	"testing"
)

func TestSerializeRoundTrip(t *testing.T) {
	c := NewCalculator()
	c.Enter("3.14")
	c.Enter("42")
	c.Enter("-7.5")
	c.SetAngleMode("RAD")
	c.SetBaseMode("HEX")

	data := c.Serialize()

	c2 := NewCalculator()
	if err := c2.Restore(data); err != nil {
		t.Fatalf("Restore error: %v", err)
	}

	s1 := c.GetState()
	s2 := c2.GetState()

	if s1.Depth != s2.Depth {
		t.Fatalf("depth mismatch: %d vs %d", s1.Depth, s2.Depth)
	}
	if s1.AngleMode != s2.AngleMode {
		t.Errorf("angle mode: %s vs %s", s1.AngleMode, s2.AngleMode)
	}
	if s1.BaseMode != s2.BaseMode {
		t.Errorf("base mode: %s vs %s", s1.BaseMode, s2.BaseMode)
	}

	for i := range s1.Stack {
		if s1.Stack[i] != s2.Stack[i] {
			t.Errorf("stack[%d]: %s vs %s", i, s1.Stack[i], s2.Stack[i])
		}
	}
}

func TestSerializeEmpty(t *testing.T) {
	c := NewCalculator()
	data := c.Serialize()

	c2 := NewCalculator()
	c2.Enter("999")
	if err := c2.Restore(data); err != nil {
		t.Fatal(err)
	}
	if c2.stack.Depth() != 0 {
		t.Errorf("restored stack should be empty, depth = %d", c2.stack.Depth())
	}
}

func TestRestoreInvalidJSON(t *testing.T) {
	c := NewCalculator()
	if err := c.Restore("not json"); err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestRestoreInvalidValue(t *testing.T) {
	c := NewCalculator()
	if err := c.Restore(`{"stack":[{"type":"real","value":"abc"}],"angleMode":"DEG","baseMode":"DEC"}`); err == nil {
		t.Error("expected error for invalid value in stack")
	}
}

func TestSerializePreservesValues(t *testing.T) {
	c := NewCalculator()
	c.Enter("1e-10")
	c.Enter("6.022e23")
	c.Execute("PI")

	data := c.Serialize()
	c2 := NewCalculator()
	c2.Restore(data)

	v0, _ := c2.stack.Peek(2)
	if math.Abs(v0.Float64()-1e-10) > 1e-25 {
		t.Errorf("value 0: %g, want 1e-10", v0.Float64())
	}

	v1, _ := c2.stack.Peek(1)
	if math.Abs(v1.Float64()-6.022e23) > 1e18 {
		t.Errorf("value 1: %g, want 6.022e23", v1.Float64())
	}

	v2, _ := c2.stack.Peek(0)
	if math.Abs(v2.Float64()-math.Pi) > 1e-15 {
		t.Errorf("value 2: %g, want pi", v2.Float64())
	}
}

func TestSerializeHistoryRoundTrip(t *testing.T) {
	c := NewCalculator()
	c.Enter("2")
	c.Enter("3")
	c.Execute("+")
	c.Execute("PI")

	data := c.Serialize()
	c2 := NewCalculator()
	if err := c2.Restore(data); err != nil {
		t.Fatal(err)
	}

	s1 := c.GetState()
	s2 := c2.GetState()
	if len(s1.History) != len(s2.History) {
		t.Fatalf("history len: %d vs %d", len(s1.History), len(s2.History))
	}
	for i := range s1.History {
		if s1.History[i].Expression != s2.History[i].Expression {
			t.Errorf("history[%d]: %q vs %q", i, s1.History[i].Expression, s2.History[i].Expression)
		}
	}
}

func TestSerializeNoHistoryBackwardCompat(t *testing.T) {
	c := NewCalculator()
	err := c.Restore(`{"stack":[{"type":"real","value":"42"}],"angleMode":"DEG","baseMode":"DEC"}`)
	if err != nil {
		t.Fatal(err)
	}
	s := c.GetState()
	if s.Depth != 1 || s.Stack[0] != "42" {
		t.Errorf("stack = %v, want [42]", s.Stack)
	}
	if len(s.History) != 0 {
		t.Errorf("history should be empty for old format, got %d", len(s.History))
	}
}

func TestUndoNotSerialized(t *testing.T) {
	c := NewCalculator()
	c.Enter("1")
	c.Enter("2")
	c.Execute("+")

	data := c.Serialize()
	c2 := NewCalculator()
	c2.Restore(data)

	s := parseState(t, c2.Undo())
	if s.Error == "" {
		t.Error("undo should fail after restore (undo history not persisted)")
	}
}
