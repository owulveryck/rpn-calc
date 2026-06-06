package engine

import (
	"errors"
	"testing"
)

func newStack(vals ...float64) *Stack {
	s := &Stack{}
	for _, v := range vals {
		s.Push(NewRealValue(v))
	}
	return s
}

func stackEquals(s *Stack, want ...float64) bool {
	if s.Depth() != len(want) {
		return false
	}
	for i, w := range want {
		v, _ := s.Peek(len(want) - 1 - i)
		if v.Float64() != w {
			return false
		}
	}
	return true
}

func TestPushPop(t *testing.T) {
	s := &Stack{}
	s.Push(NewRealValue(42))
	s.Push(NewRealValue(7))

	v, err := s.Pop()
	if err != nil {
		t.Fatal(err)
	}
	if v.Float64() != 7 {
		t.Errorf("got %v, want 7", v.Float64())
	}

	v, err = s.Pop()
	if err != nil {
		t.Fatal(err)
	}
	if v.Float64() != 42 {
		t.Errorf("got %v, want 42", v.Float64())
	}
}

func TestPopEmpty(t *testing.T) {
	s := &Stack{}
	_, err := s.Pop()
	if !errors.Is(err, ErrStackUnderflow) {
		t.Errorf("expected ErrStackUnderflow, got %v", err)
	}
}

func TestPeek(t *testing.T) {
	s := newStack(1, 2, 3)

	v, err := s.Peek(0)
	if err != nil || v.Float64() != 3 {
		t.Errorf("Peek(0) = %v, err=%v, want 3", v, err)
	}

	v, err = s.Peek(1)
	if err != nil || v.Float64() != 2 {
		t.Errorf("Peek(1) = %v, err=%v, want 2", v, err)
	}

	v, err = s.Peek(2)
	if err != nil || v.Float64() != 1 {
		t.Errorf("Peek(2) = %v, err=%v, want 1", v, err)
	}

	_, err = s.Peek(3)
	if !errors.Is(err, ErrStackUnderflow) {
		t.Errorf("Peek(3) expected underflow, got %v", err)
	}

	_, err = s.Peek(-1)
	if !errors.Is(err, ErrStackUnderflow) {
		t.Errorf("Peek(-1) expected underflow, got %v", err)
	}
}

func TestDepth(t *testing.T) {
	s := &Stack{}
	if s.Depth() != 0 {
		t.Errorf("empty stack depth = %d, want 0", s.Depth())
	}
	s.Push(NewRealValue(1))
	s.Push(NewRealValue(2))
	if s.Depth() != 2 {
		t.Errorf("depth = %d, want 2", s.Depth())
	}
	s.Pop()
	if s.Depth() != 1 {
		t.Errorf("depth after pop = %d, want 1", s.Depth())
	}
}

func TestClear(t *testing.T) {
	s := newStack(1, 2, 3)
	s.Clear()
	if s.Depth() != 0 {
		t.Errorf("depth after clear = %d, want 0", s.Depth())
	}
}

func TestSwap(t *testing.T) {
	s := newStack(1, 2)
	if err := s.Swap(); err != nil {
		t.Fatal(err)
	}
	if !stackEquals(s, 2, 1) {
		t.Errorf("after swap: got %v", s.Snapshot(BaseDec))
	}
}

func TestSwapUnderflow(t *testing.T) {
	tests := []struct {
		name string
		s    *Stack
	}{
		{"empty", newStack()},
		{"one element", newStack(1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Swap(); !errors.Is(err, ErrStackUnderflow) {
				t.Errorf("expected underflow, got %v", err)
			}
		})
	}
}

func TestDup(t *testing.T) {
	s := newStack(5)
	if err := s.Dup(); err != nil {
		t.Fatal(err)
	}
	if !stackEquals(s, 5, 5) {
		t.Errorf("after dup: got %v", s.Snapshot(BaseDec))
	}
}

func TestDupEmpty(t *testing.T) {
	s := newStack()
	if err := s.Dup(); !errors.Is(err, ErrStackUnderflow) {
		t.Errorf("expected underflow, got %v", err)
	}
}

func TestDrop(t *testing.T) {
	s := newStack(1, 2, 3)
	if err := s.Drop(); err != nil {
		t.Fatal(err)
	}
	if !stackEquals(s, 1, 2) {
		t.Errorf("after drop: got %v", s.Snapshot(BaseDec))
	}
}

func TestDropEmpty(t *testing.T) {
	s := newStack()
	if err := s.Drop(); !errors.Is(err, ErrStackUnderflow) {
		t.Errorf("expected underflow, got %v", err)
	}
}

func TestOver(t *testing.T) {
	s := newStack(1, 2)
	if err := s.Over(); err != nil {
		t.Fatal(err)
	}
	if !stackEquals(s, 1, 2, 1) {
		t.Errorf("after over: got %v", s.Snapshot(BaseDec))
	}
}

func TestOverUnderflow(t *testing.T) {
	s := newStack(1)
	if err := s.Over(); !errors.Is(err, ErrStackUnderflow) {
		t.Errorf("expected underflow, got %v", err)
	}
}

func TestRot(t *testing.T) {
	s := newStack(1, 2, 3)
	if err := s.Rot(); err != nil {
		t.Fatal(err)
	}
	if !stackEquals(s, 2, 3, 1) {
		t.Errorf("after rot: got %v, want [2,3,1]", s.Snapshot(BaseDec))
	}
}

func TestRotUnderflow(t *testing.T) {
	s := newStack(1, 2)
	if err := s.Rot(); !errors.Is(err, ErrStackUnderflow) {
		t.Errorf("expected underflow, got %v", err)
	}
}

func TestPick(t *testing.T) {
	s := newStack(10, 20, 30)
	if err := s.Pick(2); err != nil {
		t.Fatal(err)
	}
	if !stackEquals(s, 10, 20, 30, 20) {
		t.Errorf("after pick(2): got %v", s.Snapshot(BaseDec))
	}
}

func TestPickBoundary(t *testing.T) {
	s := newStack(10, 20, 30)
	if err := s.Pick(3); err != nil {
		t.Fatal(err)
	}
	if !stackEquals(s, 10, 20, 30, 10) {
		t.Errorf("after pick(3): got %v", s.Snapshot(BaseDec))
	}
}

func TestPickUnderflow(t *testing.T) {
	s := newStack(10, 20)
	if err := s.Pick(3); !errors.Is(err, ErrStackUnderflow) {
		t.Errorf("expected underflow, got %v", err)
	}
	if err := s.Pick(0); !errors.Is(err, ErrStackUnderflow) {
		t.Errorf("Pick(0) expected underflow, got %v", err)
	}
}

func TestRoll(t *testing.T) {
	s := newStack(1, 2, 3, 4)
	if err := s.Roll(3); err != nil {
		t.Fatal(err)
	}
	if !stackEquals(s, 1, 3, 4, 2) {
		t.Errorf("after roll(3): got %v, want [1,3,4,2]", s.Snapshot(BaseDec))
	}
}

func TestRollUnderflow(t *testing.T) {
	s := newStack(1, 2)
	if err := s.Roll(3); !errors.Is(err, ErrStackUnderflow) {
		t.Errorf("expected underflow, got %v", err)
	}
	if err := s.Roll(0); !errors.Is(err, ErrStackUnderflow) {
		t.Errorf("Roll(0) expected underflow, got %v", err)
	}
}

func TestSnapshot(t *testing.T) {
	s := newStack(1, 2.5, 3)
	snap := s.Snapshot(BaseDec)
	want := []string{"1", "2.5", "3"}
	if len(snap) != len(want) {
		t.Fatalf("snapshot len = %d, want %d", len(snap), len(want))
	}
	for i, w := range want {
		if snap[i] != w {
			t.Errorf("snapshot[%d] = %q, want %q", i, snap[i], w)
		}
	}
}

func TestSnapshotEmpty(t *testing.T) {
	s := newStack()
	snap := s.Snapshot(BaseDec)
	if len(snap) != 0 {
		t.Errorf("empty snapshot len = %d, want 0", len(snap))
	}
}

func TestSnapshotHex(t *testing.T) {
	s := newStack(255, 16)
	snap := s.Snapshot(BaseHex)
	if snap[0] != "#FF" || snap[1] != "#10" {
		t.Errorf("hex snapshot = %v, want [#FF, #10]", snap)
	}
}

func TestCloneRestore(t *testing.T) {
	s := newStack(1, 2, 3)
	clone := s.Clone()
	s.Clear()
	if s.Depth() != 0 {
		t.Fatal("clear didn't work")
	}
	s.Restore(clone)
	if !stackEquals(s, 1, 2, 3) {
		t.Errorf("after restore: got %v", s.Snapshot(BaseDec))
	}
}

func TestCloneEmpty(t *testing.T) {
	s := newStack()
	clone := s.Clone()
	if clone != nil {
		t.Errorf("clone of empty stack should be nil, got %v", clone)
	}
}
