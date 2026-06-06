package engine

import (
	"math"
	"testing"
)

func TestParseValue(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{"integer", "42", 42, false},
		{"negative", "-7", -7, false},
		{"decimal", "3.14", 3.14, false},
		{"negative decimal", "-0.5", -0.5, false},
		{"scientific", "6.022e23", 6.022e23, false},
		{"scientific negative exp", "1e-10", 1e-10, false},
		{"zero", "0", 0, false},
		{"leading spaces", "  42  ", 42, false},
		{"empty", "", 0, true},
		{"letters", "abc", 0, true},
		{"mixed", "12abc", 0, true},
		{"double dot", "1.2.3", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := ParseValue(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseValue(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if v.Float64() != tt.want {
				t.Errorf("ParseValue(%q) = %v, want %v", tt.input, v.Float64(), tt.want)
			}
		})
	}
}

func TestRealValueString(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		want string
	}{
		{"integer", 42, "42"},
		{"zero", 0, "0"},
		{"negative", -7, "-7"},
		{"decimal", 3.14, "3.14"},
		{"large integer", 1000000, "1000000"},
		{"scientific", 6.022e23, "6.022e+23"},
		{"small", 1e-10, "1e-10"},
		{"infinity", math.Inf(1), "+Inf"},
		{"neg infinity", math.Inf(-1), "-Inf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewRealValue(tt.val)
			got := v.String()
			if got != tt.want {
				t.Errorf("RealValue(%v).String() = %q, want %q", tt.val, got, tt.want)
			}
		})
	}
}

func TestRealValueIsComplex(t *testing.T) {
	v := NewRealValue(42)
	if v.IsComplex() {
		t.Error("RealValue should not be complex")
	}
}

func TestRealValueComplex128(t *testing.T) {
	v := NewRealValue(3.14)
	c := v.Complex128()
	if real(c) != 3.14 || imag(c) != 0 {
		t.Errorf("RealValue(3.14).Complex128() = %v, want (3.14+0i)", c)
	}
}

func TestRealValueStringInBase(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		base BaseMode
		want string
	}{
		{"dec 255", 255, BaseDec, "255"},
		{"hex 255", 255, BaseHex, "#FF"},
		{"oct 255", 255, BaseOct, "377o"},
		{"bin 255", 255, BaseBin, "11111111b"},
		{"hex 0", 0, BaseHex, "#0"},
		{"hex 16", 16, BaseHex, "#10"},
		{"bin 10", 10, BaseBin, "1010b"},
		{"oct 8", 8, BaseOct, "10o"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewRealValue(tt.val)
			got := v.StringInBase(tt.base)
			if got != tt.want {
				t.Errorf("RealValue(%v).StringInBase(%v) = %q, want %q", tt.val, tt.base, got, tt.want)
			}
		})
	}
}

func TestBaseModeString(t *testing.T) {
	tests := []struct {
		mode BaseMode
		want string
	}{
		{BaseDec, "DEC"},
		{BaseHex, "HEX"},
		{BaseOct, "OCT"},
		{BaseBin, "BIN"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("BaseMode.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
