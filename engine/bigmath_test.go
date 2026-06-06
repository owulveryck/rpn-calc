package engine

import (
	"math"
	"math/big"
	"testing"
)

func TestBigLn(t *testing.T) {
	tests := []struct {
		name    string
		input   float64
		want    float64
		wantErr bool
	}{
		{"ln(1) = 0", 1, 0, false},
		{"ln(e) = 1", math.E, 1, false},
		{"ln(2)", 2, math.Ln2, false},
		{"ln(10)", 10, math.Ln10, false},
		{"ln(0.5)", 0.5, math.Log(0.5), false},
		{"ln(100)", 100, math.Log(100), false},
		{"ln(0.001)", 0.001, math.Log(0.001), false},
		{"ln(0) error", 0, 0, true},
		{"ln(-1) error", -1, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := new(big.Float).SetPrec(defaultPrec).SetFloat64(tt.input)
			got, err := bigLn(x)
			if (err != nil) != tt.wantErr {
				t.Fatalf("bigLn(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			gotF, _ := got.Float64()
			if math.Abs(gotF-tt.want) > 1e-14 {
				t.Errorf("bigLn(%v) = %v, want %v (diff %v)", tt.input, gotF, tt.want, gotF-tt.want)
			}
		})
	}
}

func TestBigExp(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{"exp(0) = 1", 0, 1},
		{"exp(1) = e", 1, math.E},
		{"exp(-1)", -1, math.Exp(-1)},
		{"exp(2)", 2, math.Exp(2)},
		{"exp(10)", 10, math.Exp(10)},
		{"exp(-10)", -10, math.Exp(-10)},
		{"exp(0.5)", 0.5, math.Exp(0.5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := new(big.Float).SetPrec(defaultPrec).SetFloat64(tt.input)
			got := bigExp(x)
			gotF, _ := got.Float64()
			if math.Abs(gotF-tt.want)/math.Abs(tt.want) > 1e-14 {
				t.Errorf("bigExp(%v) = %v, want %v (rel diff %v)", tt.input, gotF, tt.want, (gotF-tt.want)/tt.want)
			}
		})
	}
}

func TestBigLnExpRoundTrip(t *testing.T) {
	values := []float64{1, 2, 3, 42, 85, 100, 0.5, 0.001, 1e10, 1e-10}
	for _, v := range values {
		x := new(big.Float).SetPrec(defaultPrec).SetFloat64(v)
		lnx, err := bigLn(x)
		if err != nil {
			t.Fatalf("bigLn(%v) error: %v", v, err)
		}
		result := bigExp(lnx)
		resultF, _ := result.Float64()
		if math.Abs(resultF-v)/math.Abs(v) > 1e-14 {
			t.Errorf("exp(ln(%v)) = %v (rel diff %v)", v, resultF, (resultF-v)/v)
		}
	}
}

func TestBigTripleLnExpRoundTrip(t *testing.T) {
	x := new(big.Float).SetPrec(defaultPrec).SetFloat64(85)
	var err error
	for i := 0; i < 3; i++ {
		x, err = bigLn(x)
		if err != nil {
			t.Fatalf("bigLn failed at step %d: %v", i, err)
		}
	}
	for i := 0; i < 3; i++ {
		x = bigExp(x)
	}
	result, _ := x.Float64()
	if result != 85 {
		t.Errorf("85 -> ln -> ln -> ln -> exp -> exp -> exp = %v, want 85", result)
	}
}

func TestBigPow(t *testing.T) {
	tests := []struct {
		name    string
		y, x    float64
		want    float64
		wantErr bool
	}{
		{"2^3", 2, 3, 8, false},
		{"2^0", 2, 0, 1, false},
		{"0^5", 0, 5, 0, false},
		{"0^0", 0, 0, 1, false},
		{"10^2", 10, 2, 100, false},
		{"(-2)^3", -2, 3, -8, false},
		{"(-2)^2", -2, 2, 4, false},
		{"(-2)^0.5 error", -2, 0.5, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := new(big.Float).SetPrec(defaultPrec).SetFloat64(tt.y)
			x := new(big.Float).SetPrec(defaultPrec).SetFloat64(tt.x)
			got, err := bigPow(y, x)
			if (err != nil) != tt.wantErr {
				t.Fatalf("bigPow(%v, %v) error = %v, wantErr %v", tt.y, tt.x, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			gotF, _ := got.Float64()
			if math.Abs(gotF-tt.want) > 1e-10 {
				t.Errorf("bigPow(%v, %v) = %v, want %v", tt.y, tt.x, gotF, tt.want)
			}
		})
	}
}

func TestBigConstants(t *testing.T) {
	piF, _ := bigPi().Float64()
	if math.Abs(piF-math.Pi) > 1e-15 {
		t.Errorf("bigPi = %v, want %v", piF, math.Pi)
	}

	eF, _ := bigEConst().Float64()
	if math.Abs(eF-math.E) > 1e-15 {
		t.Errorf("bigE = %v, want %v", eF, math.E)
	}

	ln2F, _ := bigLn2().Float64()
	if math.Abs(ln2F-math.Ln2) > 1e-15 {
		t.Errorf("bigLn2 = %v, want %v", ln2F, math.Ln2)
	}

	ln10F, _ := bigLn10().Float64()
	if math.Abs(ln10F-math.Log(10)) > 1e-15 {
		t.Errorf("bigLn10 = %v, want %v", ln10F, math.Log(10))
	}
}
