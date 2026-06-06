package engine

import (
	"errors"
	"math/big"
	"sync"
)

const defaultPrec = 256

var (
	cachedLn2  *big.Float
	cachedLn10 *big.Float
	cachedPi   *big.Float
	cachedE    *big.Float
	initOnce   sync.Once
)

func initConstants() {
	initOnce.Do(func() {
		cachedLn2, _, _ = new(big.Float).SetPrec(defaultPrec).Parse(
			"0.6931471805599453094172321214581765680755001343602552541206800094933936219696947156058633269964186875", 10)
		cachedLn10, _, _ = new(big.Float).SetPrec(defaultPrec).Parse(
			"2.3025850929940456840179914546843642076011014886287729760333279009675726096773524802359972050895982983", 10)
		cachedPi, _, _ = new(big.Float).SetPrec(defaultPrec).Parse(
			"3.1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679", 10)
		cachedE, _, _ = new(big.Float).SetPrec(defaultPrec).Parse(
			"2.7182818284590452353602874713526624977572470936999595749669676277240766303535475945713821785251664274", 10)
	})
}

func bigLn2() *big.Float {
	initConstants()
	return new(big.Float).SetPrec(defaultPrec).Copy(cachedLn2)
}

func bigLn10() *big.Float {
	initConstants()
	return new(big.Float).SetPrec(defaultPrec).Copy(cachedLn10)
}

func bigPi() *big.Float {
	initConstants()
	return new(big.Float).SetPrec(defaultPrec).Copy(cachedPi)
}

func bigEConst() *big.Float {
	initConstants()
	return new(big.Float).SetPrec(defaultPrec).Copy(cachedE)
}

// bigLn computes ln(x) for x > 0 using argument reduction and the
// arctanh series: ln(m) = 2*(t + t³/3 + t⁵/5 + ...) where t = (m-1)/(m+1).
// With MantExp giving m in [0.5, 1), |t| ≤ 1/3, so convergence is fast.
func bigLn(x *big.Float) (*big.Float, error) {
	if x.Sign() <= 0 {
		return nil, errors.New("logarithm of non-positive number")
	}

	prec := uint(defaultPrec)

	mant := new(big.Float).SetPrec(prec)
	exp := x.MantExp(mant) // mant in [0.5, 1), x = mant * 2^exp

	one := new(big.Float).SetPrec(prec).SetInt64(1)
	two := new(big.Float).SetPrec(prec).SetInt64(2)

	// t = (mant - 1) / (mant + 1)
	num := new(big.Float).SetPrec(prec).Sub(mant, one)
	den := new(big.Float).SetPrec(prec).Add(mant, one)
	t := new(big.Float).SetPrec(prec).Quo(num, den)

	t2 := new(big.Float).SetPrec(prec).Mul(t, t)
	term := new(big.Float).SetPrec(prec).Copy(t)
	sum := new(big.Float).SetPrec(prec).Copy(t)

	eps := new(big.Float).SetPrec(prec).SetMantExp(one, -int(prec+10))

	contrib := new(big.Float).SetPrec(prec)
	divisor := new(big.Float).SetPrec(prec)
	absTerm := new(big.Float).SetPrec(prec)

	for k := int64(3); ; k += 2 {
		term.Mul(term, t2)
		divisor.SetInt64(k)
		contrib.Quo(term, divisor)
		sum.Add(sum, contrib)
		absTerm.Abs(contrib)
		if absTerm.Cmp(eps) < 0 {
			break
		}
	}

	lnMant := new(big.Float).SetPrec(prec).Mul(two, sum)

	if exp != 0 {
		expBig := new(big.Float).SetPrec(prec).SetInt64(int64(exp))
		correction := new(big.Float).SetPrec(prec).Mul(expBig, bigLn2())
		lnMant.Add(lnMant, correction)
	}

	return lnMant, nil
}

// bigExp computes exp(x) using argument reduction and Taylor series.
// Reduction: x = k*ln(2) + r, then r = r/2^s, compute Taylor, square s times, multiply by 2^k.
func bigExp(x *big.Float) *big.Float {
	prec := uint(defaultPrec)

	if x.Sign() == 0 {
		return new(big.Float).SetPrec(prec).SetInt64(1)
	}

	ln2 := bigLn2()

	// k = round(x / ln2)
	kFloat := new(big.Float).SetPrec(prec).Quo(x, ln2)
	kInt, _ := kFloat.Int(nil)

	// Check rounding: if frac >= 0.5, adjust k
	kf := new(big.Float).SetPrec(prec).SetInt(kInt)
	frac := new(big.Float).SetPrec(prec).Sub(kFloat, kf)
	half := new(big.Float).SetPrec(prec).SetFloat64(0.5)
	negHalf := new(big.Float).SetPrec(prec).SetFloat64(-0.5)
	if frac.Cmp(half) >= 0 {
		kInt.Add(kInt, big.NewInt(1))
	} else if frac.Cmp(negHalf) < 0 {
		kInt.Sub(kInt, big.NewInt(1))
	}
	k := kInt.Int64()

	// r = x - k * ln2
	kBig := new(big.Float).SetPrec(prec).SetInt64(k)
	r := new(big.Float).SetPrec(prec).Mul(kBig, ln2)
	r.Sub(x, r)

	// Further reduce: r = r / 2^s
	const s = 10
	scaleFactor := new(big.Float).SetPrec(prec).SetMantExp(
		new(big.Float).SetPrec(prec).SetInt64(1), s)
	r.Quo(r, scaleFactor)

	// Taylor series: exp(r) = 1 + r + r²/2! + r³/3! + ...
	sum := new(big.Float).SetPrec(prec).SetInt64(1)
	term := new(big.Float).SetPrec(prec).SetInt64(1)
	eps := new(big.Float).SetPrec(prec).SetMantExp(
		new(big.Float).SetPrec(prec).SetInt64(1), -int(prec+10))
	absVal := new(big.Float).SetPrec(prec)
	nf := new(big.Float).SetPrec(prec)

	for n := int64(1); ; n++ {
		term.Mul(term, r)
		nf.SetInt64(n)
		term.Quo(term, nf)
		sum.Add(sum, term)
		absVal.Abs(term)
		if absVal.Cmp(eps) < 0 {
			break
		}
	}

	// Square s times
	for i := 0; i < s; i++ {
		sum.Mul(sum, sum)
	}

	// Multiply by 2^k
	sum.SetMantExp(sum, int(k))

	return sum
}

// bigPow computes y^x = exp(x * ln(y)).
func bigPow(y, x *big.Float) (*big.Float, error) {
	prec := uint(defaultPrec)

	if y.Sign() == 0 {
		if x.Sign() == 0 {
			return new(big.Float).SetPrec(prec).SetInt64(1), nil
		}
		return new(big.Float).SetPrec(prec).SetInt64(0), nil
	}

	if y.Sign() < 0 {
		xi, acc := x.Int(nil)
		if acc != big.Exact {
			return nil, errors.New("negative base with non-integer exponent")
		}
		absY := new(big.Float).SetPrec(prec).Abs(y)
		lnAbsY, err := bigLn(absY)
		if err != nil {
			return nil, err
		}
		product := new(big.Float).SetPrec(prec).Mul(x, lnAbsY)
		result := bigExp(product)
		if xi.Bit(0) == 1 {
			result.Neg(result)
		}
		return result, nil
	}

	lnY, err := bigLn(y)
	if err != nil {
		return nil, err
	}
	product := new(big.Float).SetPrec(prec).Mul(x, lnY)
	return bigExp(product), nil
}
