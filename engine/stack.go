package engine

import "errors"

// ErrStackUnderflow est retourné quand une opération tente de lire ou dépiler
// plus d'éléments que la pile n'en contient.
var ErrStackUnderflow = errors.New("stack underflow")

// Stack est une pile LIFO de Value utilisée comme mémoire principale de la calculatrice RPN.
type Stack struct {
	values []Value
}

// Push empile une valeur au sommet de la pile.
func (s *Stack) Push(v Value) {
	s.values = append(s.values, v)
}

// Pop retire et retourne la valeur au sommet de la pile.
// Retourne ErrStackUnderflow si la pile est vide.
func (s *Stack) Pop() (Value, error) {
	if len(s.values) == 0 {
		return nil, ErrStackUnderflow
	}
	v := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return v, nil
}

// Peek retourne la valeur à la position n depuis le sommet (0 = sommet) sans la retirer.
// Retourne ErrStackUnderflow si la position est invalide.
func (s *Stack) Peek(n int) (Value, error) {
	idx := len(s.values) - 1 - n
	if idx < 0 || idx >= len(s.values) {
		return nil, ErrStackUnderflow
	}
	return s.values[idx], nil
}

// Depth retourne le nombre d'éléments dans la pile.
func (s *Stack) Depth() int {
	return len(s.values)
}

// Clear vide la pile de tous ses éléments.
func (s *Stack) Clear() {
	s.values = s.values[:0]
}

// Swap échange les deux éléments au sommet de la pile.
func (s *Stack) Swap() error {
	n := len(s.values)
	if n < 2 {
		return ErrStackUnderflow
	}
	s.values[n-1], s.values[n-2] = s.values[n-2], s.values[n-1]
	return nil
}

// Dup duplique l'élément au sommet de la pile.
func (s *Stack) Dup() error {
	if len(s.values) == 0 {
		return ErrStackUnderflow
	}
	s.values = append(s.values, s.values[len(s.values)-1])
	return nil
}

// Drop retire l'élément au sommet de la pile sans le retourner.
func (s *Stack) Drop() error {
	if len(s.values) == 0 {
		return ErrStackUnderflow
	}
	s.values = s.values[:len(s.values)-1]
	return nil
}

// Over copie le second élément de la pile et l'empile au sommet.
func (s *Stack) Over() error {
	n := len(s.values)
	if n < 2 {
		return ErrStackUnderflow
	}
	s.values = append(s.values, s.values[n-2])
	return nil
}

// Rot effectue une rotation des trois éléments au sommet de la pile :
// le troisième élément est déplacé au sommet.
func (s *Stack) Rot() error {
	n := len(s.values)
	if n < 3 {
		return ErrStackUnderflow
	}
	s.values[n-3], s.values[n-2], s.values[n-1] = s.values[n-2], s.values[n-1], s.values[n-3]
	return nil
}

// Pick copie le n-ième élément (1 = sommet) et l'empile au sommet.
func (s *Stack) Pick(n int) error {
	idx := len(s.values) - n
	if n < 1 || idx < 0 {
		return ErrStackUnderflow
	}
	s.values = append(s.values, s.values[idx])
	return nil
}

// Roll déplace le n-ième élément au sommet, décalant les éléments intermédiaires vers le bas.
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

// Snapshot retourne une représentation textuelle de tous les éléments de la pile
// dans la base numérique spécifiée.
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

// Clone retourne une copie superficielle du contenu de la pile pour la sauvegarde d'état.
func (s *Stack) Clone() []Value {
	if len(s.values) == 0 {
		return nil
	}
	clone := make([]Value, len(s.values))
	copy(clone, s.values)
	return clone
}

// Restore remplace le contenu de la pile par les valeurs fournies.
func (s *Stack) Restore(values []Value) {
	s.values = values
}
