package engine

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type serializedNode struct {
	Op       string            `json:"op"`
	Value    string            `json:"value,omitempty"`
	Children []*serializedNode `json:"children,omitempty"`
}

type serializedValue struct {
	Type string          `json:"type"`
	Val  string          `json:"value,omitempty"`
	Expr *serializedNode `json:"expr,omitempty"`
}

type serializedState struct {
	Stack     []serializedValue `json:"stack"`
	AngleMode string            `json:"angleMode"`
	BaseMode  string            `json:"baseMode"`
	History   []HistoryEntry    `json:"history,omitempty"`
}

func serializeNode(n *ExprNode) *serializedNode {
	if n == nil {
		return nil
	}
	sn := &serializedNode{Op: n.op.String()}
	if n.op == OpLiteral && n.literal != nil {
		sn.Value = n.literal.Text('g', 80)
	}
	if len(n.children) > 0 {
		sn.Children = make([]*serializedNode, len(n.children))
		for i, child := range n.children {
			sn.Children[i] = serializeNode(child)
		}
	}
	return sn
}

func deserializeNode(sn *serializedNode) (*ExprNode, error) {
	if sn == nil {
		return nil, fmt.Errorf("nil node")
	}
	op, ok := parseExprOp(sn.Op)
	if !ok {
		return nil, fmt.Errorf("unknown op: %s", sn.Op)
	}
	node := &ExprNode{op: op}
	if op == OpLiteral {
		if sn.Value == "" {
			node.literal = new(big.Float).SetPrec(defaultPrec)
		} else {
			f, _, err := new(big.Float).SetPrec(defaultPrec).Parse(sn.Value, 10)
			if err != nil {
				return nil, fmt.Errorf("invalid literal %q: %w", sn.Value, err)
			}
			node.literal = f
		}
	}
	if len(sn.Children) > 0 {
		node.children = make([]*ExprNode, len(sn.Children))
		for i, child := range sn.Children {
			cn, err := deserializeNode(child)
			if err != nil {
				return nil, err
			}
			node.children[i] = cn
		}
	}
	return node, nil
}

// Serialize retourne l'état complet de la calculatrice sous forme de chaîne JSON,
// incluant les arbres d'expressions pour permettre une restauration fidèle.
func (c *Calculator) Serialize() string {
	sv := make([]serializedValue, c.stack.Depth())
	for i, v := range c.stack.values {
		if ev, ok := v.(ExprValue); ok {
			sv[i] = serializedValue{
				Type: "expr",
				Expr: serializeNode(ev.node),
			}
		} else {
			sv[i] = serializedValue{Type: "real", Val: v.BigFloat().Text('g', 80)}
		}
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

// Restore reconstruit l'état de la calculatrice à partir d'une chaîne JSON
// produite par Serialize. L'historique d'annulation est réinitialisé.
func (c *Calculator) Restore(data string) error {
	var s serializedState
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return fmt.Errorf("invalid state: %w", err)
	}
	c.stack.Clear()
	c.undoStack = nil
	c.lastArgs = nil
	for _, sv := range s.Stack {
		switch sv.Type {
		case "expr":
			if sv.Expr == nil {
				return fmt.Errorf("expr value with nil expression")
			}
			node, err := deserializeNode(sv.Expr)
			if err != nil {
				return fmt.Errorf("invalid expression: %w", err)
			}
			c.stack.Push(NewExprValue(node))
		default:
			// Backward compat: "real" type or unknown → literal node
			v, err := ParseValue(sv.Val)
			if err != nil {
				return fmt.Errorf("invalid value %q: %w", sv.Val, err)
			}
			c.stack.Push(NewExprValue(literalNode(v.BigFloat())))
		}
	}
	c.SetAngleMode(s.AngleMode)
	c.SetBaseMode(s.BaseMode)
	c.history = s.History
	return nil
}
