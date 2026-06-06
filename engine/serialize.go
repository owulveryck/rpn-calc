package engine

import (
	"encoding/json"
	"fmt"
)

type serializedValue struct {
	Type string `json:"type"`
	Val  string `json:"value"`
}

type serializedState struct {
	Stack     []serializedValue `json:"stack"`
	AngleMode string            `json:"angleMode"`
	BaseMode  string            `json:"baseMode"`
	History   []HistoryEntry    `json:"history,omitempty"`
}

func (c *Calculator) Serialize() string {
	sv := make([]serializedValue, c.stack.Depth())
	for i, v := range c.stack.values {
		sv[i] = serializedValue{Type: "real", Val: v.BigFloat().Text('g', 80)}
	}
	s := serializedState{
		Stack:     sv,
		AngleMode: c.angleMode.String(),
		BaseMode:  c.baseMode.String(),
		History:   c.history,
	}
	b, _ := json.Marshal(s)
	return string(b)
}

func (c *Calculator) Restore(data string) error {
	var s serializedState
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return fmt.Errorf("invalid state: %w", err)
	}
	c.stack.Clear()
	c.undoStack = nil
	c.lastArgs = nil
	for _, sv := range s.Stack {
		v, err := ParseValue(sv.Val)
		if err != nil {
			return fmt.Errorf("invalid value %q: %w", sv.Val, err)
		}
		c.stack.Push(v)
	}
	c.SetAngleMode(s.AngleMode)
	c.SetBaseMode(s.BaseMode)
	c.history = s.History
	return nil
}
