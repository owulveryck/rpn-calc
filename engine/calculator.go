// Package engine implémente un moteur de calculatrice RPN (Reverse Polish Notation)
// avec arithmétique en précision arbitraire, fonctions trigonométriques,
// logarithmiques et de manipulation de pile.
package engine

import (
	"encoding/json"
	"fmt"
)

// AngleMode représente l'unité angulaire utilisée pour les fonctions trigonométriques.
type AngleMode int

const (
	AngleRad  AngleMode = iota // Radians
	AngleDeg                   // Degrés
	AngleGrad                  // Gradians
)

// String retourne la représentation textuelle du mode angulaire ("RAD", "DEG" ou "GRAD").
func (m AngleMode) String() string {
	switch m {
	case AngleDeg:
		return "DEG"
	case AngleGrad:
		return "GRAD"
	default:
		return "RAD"
	}
}

// Calculator est le moteur principal de la calculatrice RPN.
// Il gère une pile de valeurs, un historique d'annulation et un journal des opérations.
type Calculator struct {
	stack     Stack
	undoStack [][]Value
	lastArgs  []Value
	angleMode AngleMode
	baseMode  BaseMode
	errorMsg  string
	history   []HistoryEntry
}

// HistoryEntry représente une entrée dans l'historique des opérations de la calculatrice.
type HistoryEntry struct {
	Expression string `json:"expr"`
}

// State représente l'état complet de la calculatrice, sérialisable en JSON.
type State struct {
	Stack     []string       `json:"stack"`
	Depth     int            `json:"depth"`
	Error     string         `json:"error"`
	AngleMode string         `json:"angleMode"`
	BaseMode  string         `json:"baseMode"`
	History   []HistoryEntry `json:"history"`
}

const maxUndo = 100
const maxHistory = 50

// NewCalculator crée une nouvelle calculatrice RPN avec le mode degrés et la base décimale par défaut.
func NewCalculator() *Calculator {
	return &Calculator{
		angleMode: AngleDeg,
		baseMode:  BaseDec,
	}
}

func (c *Calculator) saveUndo() {
	snapshot := c.stack.Clone()
	c.undoStack = append(c.undoStack, snapshot)
	if len(c.undoStack) > maxUndo {
		c.undoStack = c.undoStack[1:]
	}
}

func (c *Calculator) addHistory(expr string) {
	c.history = append(c.history, HistoryEntry{Expression: expr})
	if len(c.history) > maxHistory {
		c.history = c.history[len(c.history)-maxHistory:]
	}
}

// GetState retourne l'état actuel de la calculatrice sous forme de struct State.
func (c *Calculator) GetState() State {
	h := c.history
	if h == nil {
		h = []HistoryEntry{}
	}
	return State{
		Stack:     c.stack.Snapshot(c.baseMode),
		Depth:     c.stack.Depth(),
		Error:     c.errorMsg,
		AngleMode: c.angleMode.String(),
		BaseMode:  c.baseMode.String(),
		History:   h,
	}
}

// GetStateJSON retourne l'état actuel de la calculatrice sous forme de chaîne JSON.
func (c *Calculator) GetStateJSON() string {
	s := c.GetState()
	b, _ := json.Marshal(s)
	return string(b)
}

// Enter analyse la chaîne input comme une valeur numérique et l'empile.
// Retourne l'état JSON après l'opération.
func (c *Calculator) Enter(input string) string {
	c.errorMsg = ""
	v, err := ParseValueInBase(input, c.baseMode)
	if err != nil {
		c.errorMsg = err.Error()
		return c.GetStateJSON()
	}
	c.saveUndo()
	ev := NewExprValue(literalNode(v.BigFloat()))
	c.stack.Push(ev)
	c.addHistory("ENTER " + ev.String())
	return c.GetStateJSON()
}

// Execute exécute l'opération nommée op sur la pile.
// Retourne l'état JSON après l'opération.
func (c *Calculator) Execute(op string) string {
	c.errorMsg = ""
	info, ok := operations[op]
	if !ok {
		c.errorMsg = fmt.Sprintf("unknown operation: %s", op)
		return c.GetStateJSON()
	}

	// capture operands before execution for history
	var argX, argY string
	if info.arity == arityBinary && c.stack.Depth() >= 2 {
		x, _ := c.stack.Peek(0)
		y, _ := c.stack.Peek(1)
		argX = x.String()
		argY = y.String()
	} else if info.arity == arityUnary && c.stack.Depth() >= 1 {
		x, _ := c.stack.Peek(0)
		argX = x.String()
	}

	c.saveUndo()
	if err := info.fn(c); err != nil {
		c.errorMsg = err.Error()
		if len(c.undoStack) > 0 {
			c.stack.Restore(c.undoStack[len(c.undoStack)-1])
			c.undoStack = c.undoStack[:len(c.undoStack)-1]
		}
		return c.GetStateJSON()
	}

	c.addHistory(c.formatHistory(info, argX, argY))
	return c.GetStateJSON()
}

func (c *Calculator) formatHistory(info opInfo, argX, argY string) string {
	switch info.arity {
	case arityBinary:
		result, _ := c.stack.Peek(0)
		return fmt.Sprintf("%s %s %s = %s", argY, info.symbol, argX, result.String())
	case arityUnary:
		result, _ := c.stack.Peek(0)
		if info.symbol == "!" {
			return fmt.Sprintf("%s! = %s", argX, result.String())
		}
		if info.symbol == "neg" {
			return fmt.Sprintf("neg(%s) = %s", argX, result.String())
		}
		if info.symbol == "sq" {
			return fmt.Sprintf("%s² = %s", argX, result.String())
		}
		if info.symbol == "1/" {
			return fmt.Sprintf("1/%s = %s", argX, result.String())
		}
		if info.symbol == "10^" {
			return fmt.Sprintf("10^%s = %s", argX, result.String())
		}
		return fmt.Sprintf("%s(%s) = %s", info.symbol, argX, result.String())
	case arityConst:
		result, _ := c.stack.Peek(0)
		return fmt.Sprintf("%s = %s", info.symbol, result.String())
	default:
		return info.symbol
	}
}

// Undo annule la dernière opération en restaurant l'état précédent de la pile.
// Retourne l'état JSON après l'annulation.
func (c *Calculator) Undo() string {
	c.errorMsg = ""
	if len(c.undoStack) == 0 {
		c.errorMsg = "nothing to undo"
		return c.GetStateJSON()
	}
	snapshot := c.undoStack[len(c.undoStack)-1]
	c.undoStack = c.undoStack[:len(c.undoStack)-1]
	c.stack.Restore(snapshot)
	return c.GetStateJSON()
}

// Last réempile les derniers arguments consommés par une opération.
// Retourne l'état JSON après l'opération.
func (c *Calculator) Last() string {
	c.errorMsg = ""
	if len(c.lastArgs) == 0 {
		c.errorMsg = "no last arguments"
		return c.GetStateJSON()
	}
	c.saveUndo()
	for _, v := range c.lastArgs {
		c.stack.Push(v)
	}
	return c.GetStateJSON()
}

// SetAngleMode change le mode angulaire de la calculatrice ("DEG", "RAD" ou "GRAD").
// Retourne l'état JSON après le changement.
func (c *Calculator) SetAngleMode(mode string) string {
	switch mode {
	case "DEG":
		c.angleMode = AngleDeg
	case "RAD":
		c.angleMode = AngleRad
	case "GRAD":
		c.angleMode = AngleGrad
	}
	return c.GetStateJSON()
}

// SetBaseMode change la base numérique d'affichage ("DEC", "HEX", "OCT" ou "BIN").
// Retourne l'état JSON après le changement.
func (c *Calculator) SetBaseMode(mode string) string {
	switch mode {
	case "HEX":
		c.baseMode = BaseHex
	case "OCT":
		c.baseMode = BaseOct
	case "BIN":
		c.baseMode = BaseBin
	default:
		c.baseMode = BaseDec
	}
	return c.GetStateJSON()
}
