// MIT License
//
// Copyright (c) 2022 Tommy TIAN
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
	rSign := g.R.Sign()
	iSign := g.I.Sign()
	res := ""
	if rSign != 0 {
		res += g.R.String()
	}
	if iSign == 0 {
		if res == "" {
			return "0"
		}
		return res
	}
	if iSign == 1 && rSign != 0 {
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
	opt := iPool.Get().(*big.Int).Mul(g.I, g.I)
	defer iPool.Put(opt)
	norm.Add(norm, opt)
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
	rsInt = roundFloat(realScalar)
	isInt := iPool.Get().(*big.Int)
	defer iPool.Put(isInt)
	isInt = roundFloat(imagScalar)
	quotient := NewGaussianInt(rsInt, isInt)
	opt := giPool.Get().(*GaussianInt)
	defer giPool.Put(opt)
	g.Sub(a, opt.Prod(quotient, b))
	return quotient
}

// Equals checks if two Gaussian integers are equal
func (g *GaussianInt) Equals(a *GaussianInt) bool {
	return g.R.Cmp(a.R) == 0 && g.I.Cmp(a.I) == 0
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
