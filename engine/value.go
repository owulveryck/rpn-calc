package engine

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// ErrInvalidNumber est retourné quand une chaîne ne peut pas être interprétée comme un nombre valide.
var ErrInvalidNumber = errors.New("invalid number")

// Value est l'interface commune à toutes les valeurs manipulées par la calculatrice.
// Elle permet d'obtenir des représentations numériques en différentes précisions et bases.
type Value interface {
	String() string
	Float64() float64
	BigFloat() *big.Float
	IsComplex() bool
	Complex128() complex128
	StringInBase(base BaseMode) string
}

// BaseMode représente la base numérique utilisée pour l'affichage des valeurs.
type BaseMode int

const (
	BaseDec BaseMode = iota // Décimal
	BaseHex                 // Hexadécimal
	BaseOct                 // Octal
	BaseBin                 // Binaire
)

// String retourne la représentation textuelle de la base ("DEC", "HEX", "OCT" ou "BIN").
func (b BaseMode) String() string {
	switch b {
	case BaseHex:
		return "HEX"
	case BaseOct:
		return "OCT"
	case BaseBin:
		return "BIN"
	default:
		return "DEC"
	}
}

// RealValue représente une valeur réelle en précision arbitraire via big.Float.
type RealValue struct {
	val *big.Float
}

// NewRealValue crée une RealValue à partir d'un float64.
func NewRealValue(v float64) RealValue {
	bf := new(big.Float).SetPrec(defaultPrec)
	if math.IsNaN(v) {
		bf.SetFloat64(0)
	} else {
		bf.SetFloat64(v)
	}
	return RealValue{val: bf}
}

// ParseValue analyse une chaîne en base 10 et retourne la Value correspondante.
func ParseValue(s string) (Value, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, ErrInvalidNumber
	}
	f, _, err := new(big.Float).SetPrec(defaultPrec).Parse(s, 10)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidNumber, s)
	}
	return RealValue{val: f}, nil
}

// ParseValueInBase analyse une chaîne dans la base numérique spécifiée et retourne la Value correspondante.
func ParseValueInBase(s string, base BaseMode) (Value, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, ErrInvalidNumber
	}
	if base == BaseDec {
		return ParseValue(s)
	}
	var b int
	switch base {
	case BaseHex:
		b = 16
	case BaseOct:
		b = 8
	case BaseBin:
		b = 2
	}
	n, err := strconv.ParseInt(s, b, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidNumber, s)
	}
	return NewRealValue(float64(n)), nil
}

func formatFloat(bf *big.Float) string {
	f, _ := bf.Float64()
	if !math.IsInf(f, 0) {
		if f == math.Trunc(f) && math.Abs(f) < 1e15 {
			return strconv.FormatFloat(f, 'f', -1, 64)
		}
		return strconv.FormatFloat(f, 'g', -1, 64)
	}
	return bf.Text('g', 15)
}

// String retourne la représentation décimale de la valeur.
func (v RealValue) String() string {
	return formatFloat(v.val)
}

// Float64 retourne la valeur convertie en float64.
func (v RealValue) Float64() float64 {
	f, _ := v.val.Float64()
	return f
}

// BigFloat retourne une copie de la valeur interne en big.Float.
func (v RealValue) BigFloat() *big.Float {
	return new(big.Float).SetPrec(defaultPrec).Copy(v.val)
}

// IsComplex retourne false car RealValue ne représente pas un nombre complexe.
func (v RealValue) IsComplex() bool {
	return false
}

// Complex128 retourne la valeur comme nombre complexe avec partie imaginaire nulle.
func (v RealValue) Complex128() complex128 {
	return complex(v.Float64(), 0)
}

// StringInBase retourne la représentation de la valeur dans la base spécifiée.
func (v RealValue) StringInBase(base BaseMode) string {
	n := int64(v.Float64())
	switch base {
	case BaseHex:
		return fmt.Sprintf("#%X", n)
	case BaseOct:
		return fmt.Sprintf("%oo", n)
	case BaseBin:
		return fmt.Sprintf("%bb", n)
	default:
		return v.String()
	}
}

// ExprValue représente une valeur calculée définie par un arbre d'expressions.
// La valeur numérique est évaluée paresseusement à partir du graphe d'opérations.
type ExprValue struct {
	node *ExprNode
}

// NewExprValue crée une ExprValue à partir d'un noeud d'expression.
func NewExprValue(node *ExprNode) ExprValue {
	return ExprValue{node: node}
}

// String évalue l'expression et retourne sa représentation décimale.
func (v ExprValue) String() string {
	return formatFloat(v.node.evaluate())
}

// Float64 évalue l'expression et retourne le résultat en float64.
func (v ExprValue) Float64() float64 {
	f, _ := v.node.evaluate().Float64()
	return f
}

// BigFloat évalue l'expression et retourne le résultat en big.Float.
func (v ExprValue) BigFloat() *big.Float {
	return v.node.evaluate()
}

// IsComplex retourne false car ExprValue ne représente pas un nombre complexe.
func (v ExprValue) IsComplex() bool {
	return false
}

// Complex128 évalue l'expression et retourne le résultat comme nombre complexe.
func (v ExprValue) Complex128() complex128 {
	return complex(v.Float64(), 0)
}

// StringInBase évalue l'expression et retourne le résultat dans la base spécifiée.
func (v ExprValue) StringInBase(base BaseMode) string {
	n := int64(v.Float64())
	switch base {
	case BaseHex:
		return fmt.Sprintf("#%X", n)
	case BaseOct:
		return fmt.Sprintf("%oo", n)
	case BaseBin:
		return fmt.Sprintf("%bb", n)
	default:
		return v.String()
	}
}
