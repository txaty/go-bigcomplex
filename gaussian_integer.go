package complex

import (
	"math/big"
)

// GaussianInt implements Gaussian integer
// In number theory, a Gaussian integer is a complex number whose real and imaginary parts are both integers
type GaussianInt struct {
	R *big.Int // real part
	I *big.Int // imaginary part
}

// String returns the string representation of the Gaussian integer
func (g *GaussianInt) String() string {
	res := ""
	if g.R.Sign() != 0 {
		res += g.R.String()
	}
	gISign := g.I.Sign()
	if gISign == 0 {
		if res == "" {
			return "0"
		}
		return res
	}
	if gISign == 1 {
		res += "+"
	}
	if g.I.Cmp(bigNeg1) == 0 {
		res += "-"
	} else if g.I.Cmp(big1) != 0 {
		res += g.I.String()
	}
	res += "i"
	return res
}

// NewGaussianInt declares a new Gaussian integer with the real part and imaginary part
func NewGaussianInt(r *big.Int, i *big.Int) *GaussianInt {
	return &GaussianInt{
		R: new(big.Int).Set(r),
		I: new(big.Int).Set(i),
	}
}

// Set sets the Gaussian integer to the given Gaussian integer
func (g *GaussianInt) Set(a *GaussianInt) *GaussianInt {
	if g.R == nil {
		g.R = new(big.Int)
	}
	g.R.Set(a.R)
	if g.I == nil {
		g.I = new(big.Int)
	}
	g.I.Set(a.I)
	return g
}

// Update updates the Gaussian integer with the given real and imaginary parts
func (g *GaussianInt) Update(r, i *big.Int) *GaussianInt {
	if g.R == nil {
		g.R = new(big.Int)
	}
	g.R.Set(r)
	if g.I == nil {
		g.I = new(big.Int)
	}
	g.I.Set(i)
	return g
}

// Add adds two Gaussian integers
func (g *GaussianInt) Add(a, b *GaussianInt) *GaussianInt {
	if g.R == nil {
		g.R = new(big.Int)
	}
	g.R.Add(a.R, b.R)
	if g.I == nil {
		g.I = new(big.Int)
	}
	g.I.Add(a.I, b.I)
	return g
}

// Sub subtracts two Gaussian integers
func (g *GaussianInt) Sub(a, b *GaussianInt) *GaussianInt {
	if g.R == nil {
		g.R = new(big.Int)
	}
	g.R.Sub(a.R, b.R)
	if g.I == nil {
		g.I = new(big.Int)
	}
	g.I.Sub(a.I, b.I)
	return g
}

// Prod returns the products of two Gaussian integers
func (g *GaussianInt) Prod(a, b *GaussianInt) *GaussianInt {
	r := new(big.Int).Mul(a.R, b.R)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	r.Sub(r, opt.Mul(a.I, b.I))
	i := new(big.Int).Mul(a.R, b.I)
	i.Add(i, opt.Mul(a.I, b.R))
	g.R, g.I = r, i
	return g
}

// Conj obtains the conjugate of the original Gaussian integer
func (g *GaussianInt) Conj(origin *GaussianInt) *GaussianInt {
	img := new(big.Int).Neg(origin.I)
	g.Update(origin.R, img)
	return g
}

// Norm obtains the norm of the Gaussian integer
func (g *GaussianInt) Norm() *big.Int {
	norm := new(big.Int).Mul(g.R, g.R)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	norm.Add(norm, opt.Mul(g.I, g.I))
	return norm
}

// Copy copies the Gaussian integer
func (g *GaussianInt) Copy() *GaussianInt {
	return NewGaussianInt(
		new(big.Int).Set(g.R),
		new(big.Int).Set(g.I),
	)
}

// Div performs Euclidean division of two Gaussian integers, i.e. a/b
// the remainder is stored in the Gaussian integer that calls the method
// the quotient is returned as a new Gaussian integer
func (g *GaussianInt) Div(a, b *GaussianInt) *GaussianInt {
	bConj := giPool.Get().(*GaussianInt).Conj(b)
	defer giPool.Put(bConj)
	numerator := giPool.Get().(*GaussianInt).Prod(a, bConj)
	defer giPool.Put(numerator)
	denominator := giPool.Get().(*GaussianInt).Prod(b, bConj)
	defer giPool.Put(denominator)
	deFloat := fPool.Get().(*big.Float).SetInt(denominator.R)
	defer fPool.Put(deFloat)

	realScalar := fPool.Get().(*big.Float).SetInt(numerator.R)
	defer fPool.Put(realScalar)
	realScalar.Quo(realScalar, deFloat)
	imagScalar := fPool.Get().(*big.Float).SetInt(numerator.I)
	defer fPool.Put(imagScalar)
	imagScalar.Quo(imagScalar, deFloat)

	rsInt := iPool.Get().(*big.Int)
	defer iPool.Put(rsInt)
	roundFloat(realScalar, rsInt)
	isInt := iPool.Get().(*big.Int)
	defer iPool.Put(isInt)
	roundFloat(imagScalar, isInt)
	quotient := NewGaussianInt(rsInt, isInt)
	opt := giPool.Get().(*GaussianInt)
	defer giPool.Put(opt)
	g.Sub(a, opt.Prod(quotient, b))
	return quotient
}

// IsZero returns true if the Gaussian integer is equal to zero
func (g *GaussianInt) IsZero() bool {
	return g.R.Sign() == 0 && g.I.Sign() == 0
}

// IsOne returns true if the Gaussian integer is equal to one
func (g *GaussianInt) IsOne() bool {
	return g.R.Sign() == 1 && g.I.Sign() == 0
}

// CmpNorm compares the norm of two Gaussian integers
func (g *GaussianInt) CmpNorm(a *GaussianInt) int {
	return g.Norm().Cmp(a.Norm())
}

// GCD calculates the greatest common divisor of two Gaussian integers using Euclidean algorithm
// the result is stored in the Gaussian integer that calls the method and returned
func (g *GaussianInt) GCD(a, b *GaussianInt) *GaussianInt {
	ac := giPool.Get().(*GaussianInt).Set(a)
	defer giPool.Put(ac)
	bc := giPool.Get().(*GaussianInt).Set(b)
	defer giPool.Put(bc)

	if ac.CmpNorm(bc) < 0 {
		ac, bc = bc, ac
	}
	remainder := giPool.Get().(*GaussianInt)
	defer giPool.Put(remainder)
	for {
		remainder.Div(ac, bc)
		if remainder.IsZero() {
			g.Set(bc)
			return new(GaussianInt).Set(bc)
		}
		ac.Set(bc)
		bc.Set(remainder)
	}
}
