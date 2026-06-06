package engine

import "errors"

var ErrStackUnderflow = errors.New("stack underflow")

type Stack struct {
	values []Value
}

func (s *Stack) Push(v Value) {
	s.values = append(s.values, v)
}

func (s *Stack) Pop() (Value, error) {
	if len(s.values) == 0 {
		return nil, ErrStackUnderflow
	}
	v := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return v, nil
}

func (s *Stack) Peek(n int) (Value, error) {
	idx := len(s.values) - 1 - n
	if idx < 0 || idx >= len(s.values) {
		return nil, ErrStackUnderflow
	}
	return s.values[idx], nil
}

func (s *Stack) Depth() int {
	return len(s.values)
}

func (s *Stack) Clear() {
	s.values = s.values[:0]
}

func (s *Stack) Swap() error {
	n := len(s.values)
	if n < 2 {
		return ErrStackUnderflow
	}
	s.values[n-1], s.values[n-2] = s.values[n-2], s.values[n-1]
	return nil
}

func (s *Stack) Dup() error {
	if len(s.values) == 0 {
		return ErrStackUnderflow
	}
	s.values = append(s.values, s.values[len(s.values)-1])
	return nil
}

func (s *Stack) Drop() error {
	if len(s.values) == 0 {
		return ErrStackUnderflow
	}
	s.values = s.values[:len(s.values)-1]
	return nil
}

func (s *Stack) Over() error {
	n := len(s.values)
	if n < 2 {
		return ErrStackUnderflow
	}
	s.values = append(s.values, s.values[n-2])
	return nil
}

func (s *Stack) Rot() error {
	n := len(s.values)
	if n < 3 {
		return ErrStackUnderflow
	}
	s.values[n-3], s.values[n-2], s.values[n-1] = s.values[n-2], s.values[n-1], s.values[n-3]
	return nil
}

func (s *Stack) Pick(n int) error {
	idx := len(s.values) - n
	if n < 1 || idx < 0 {
		return ErrStackUnderflow
	}
	s.values = append(s.values, s.values[idx])
	return nil
}

func (s *Stack) Roll(n int) error {
	sz := len(s.values)
	if n < 1 || n > sz {
		return ErrStackUnderflow
	}
	idx := sz - n
	v := s.values[idx]
	copy(s.values[idx:], s.values[idx+1:])
	s.values[sz-1] = v
	return nil
}

func (s *Stack) Snapshot(base BaseMode) []string {
	result := make([]string, len(s.values))
	for i, v := range s.values {
		if base == BaseDec {
			result[i] = v.String()
		} else {
			result[i] = v.StringInBase(base)
		}
	}
	return result
}

func (s *Stack) Clone() []Value {
	if len(s.values) == 0 {
		return nil
	}
	clone := make([]Value, len(s.values))
	copy(clone, s.values)
	return clone
}

func (s *Stack) Restore(values []Value) {
	s.values = values
}
