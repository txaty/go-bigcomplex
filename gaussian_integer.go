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

// GaussianInt represents a Gaussian integer, that is, a complex number
// whose real and imaginary parts are both integers.
type GaussianInt struct {
	R *big.Int // Real part
	I *big.Int // Imaginary part
}

// String returns the string representation of the Gaussian integer.
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

// NewGaussianInt creates a new Gaussian integer with the specified real and imaginary parts.
func NewGaussianInt(r *big.Int, i *big.Int) *GaussianInt {
	return &GaussianInt{
		R: new(big.Int).Set(r),
		I: new(big.Int).Set(i),
	}
}

// Set assigns the value of another Gaussian integer to this one.
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

// Update sets the real and imaginary parts of the Gaussian integer.
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

// Add computes the sum of two Gaussian integers and stores the result in the receiver.
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

// Sub subtracts one Gaussian integer from another and stores the result in the receiver.
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

// Prod calculates the product of two Gaussian integers and stores the result in the receiver.
func (g *GaussianInt) Prod(a, b *GaussianInt) *GaussianInt {
	// Compute: (a.R + a.I*i) * (b.R + b.I*i)
	// = (a.R*b.R - a.I*b.I) + (a.R*b.I + a.I*b.R)*i
	r := new(big.Int).Mul(a.R, b.R)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	r.Sub(r, opt.Mul(a.I, b.I))
	i := new(big.Int).Mul(a.R, b.I)
	i.Add(i, opt.Mul(a.I, b.R))
	g.R, g.I = r, i
	return g
}

// Conj computes the conjugate of the given Gaussian integer and stores it in the receiver.
func (g *GaussianInt) Conj(origin *GaussianInt) *GaussianInt {
	// The conjugate of (R + I*i) is (R - I*i).
	img := new(big.Int).Neg(origin.I)
	g.Update(origin.R, img)
	return g
}

// Norm returns the norm of the Gaussian integer (R^2 + I^2).
func (g *GaussianInt) Norm() *big.Int {
	norm := new(big.Int).Mul(g.R, g.R)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	opt.Mul(g.I, g.I)
	norm.Add(norm, opt)
	return norm
}

// Copy creates a deep copy of the Gaussian integer.
func (g *GaussianInt) Copy() *GaussianInt {
	return NewGaussianInt(
		new(big.Int).Set(g.R),
		new(big.Int).Set(g.I),
	)
}

// Div performs Euclidean division of two Gaussian integers (a / b).
// The remainder is stored in the receiver, and the quotient is returned as a new Gaussian integer.
func (g *GaussianInt) Div(a, b *GaussianInt) *GaussianInt {
	// Compute the conjugate of b.
	bConj := new(GaussianInt).Conj(b)
	// Numerator = a * conjugate(b)
	numerator := new(GaussianInt).Prod(a, bConj)
	// Denominator = b * conjugate(b)
	denominator := new(GaussianInt).Prod(b, bConj)
	// Use the real part of the denominator for the division.
	deFloat := fPool.Get().(*big.Float).SetInt(denominator.R)
	defer fPool.Put(deFloat)

	// Compute the real quotient component.
	realScalar := fPool.Get().(*big.Float).SetInt(numerator.R)
	defer fPool.Put(realScalar)
	realScalar.Quo(realScalar, deFloat)
	// Compute the imaginary quotient component.
	imgScalar := fPool.Get().(*big.Float).SetInt(numerator.I)
	defer fPool.Put(imgScalar)
	imgScalar.Quo(imgScalar, deFloat)

	// Round the computed float values to the nearest integers.
	rsInt := iPool.Get().(*big.Int)
	defer iPool.Put(rsInt)
	rsInt = roundFloat(realScalar)
	isInt := iPool.Get().(*big.Int)
	defer iPool.Put(isInt)
	isInt = roundFloat(imgScalar)
	quotient := NewGaussianInt(rsInt, isInt)

	// Compute the remainder: remainder = a - (quotient * b)
	opt := new(GaussianInt).Prod(quotient, b)
	g.Sub(a, opt)
	return quotient
}

// Equals returns true if the Gaussian integer is equal to the given Gaussian integer.
func (g *GaussianInt) Equals(a *GaussianInt) bool {
	return g.R.Cmp(a.R) == 0 && g.I.Cmp(a.I) == 0
}

// IsZero returns true if the Gaussian integer is zero.
func (g *GaussianInt) IsZero() bool {
	return g.R.Sign() == 0 && g.I.Sign() == 0
}

// IsOne returns true if the Gaussian integer equals one.
func (g *GaussianInt) IsOne() bool {
	return g.R.Sign() == 1 && g.I.Sign() == 0
}

// CmpNorm compares the norms of two Gaussian integers.
func (g *GaussianInt) CmpNorm(a *GaussianInt) int {
	return g.Norm().Cmp(a.Norm())
}

// GCD calculates the greatest common divisor of two Gaussian integers using the Euclidean algorithm.
// The result is stored in the receiver and also returned as a new Gaussian integer.
func (g *GaussianInt) GCD(a, b *GaussianInt) *GaussianInt {
	ac := new(GaussianInt).Set(a)
	bc := new(GaussianInt).Set(b)

	if ac.CmpNorm(bc) < 0 {
		ac, bc = bc, ac
	}
	remainder := new(GaussianInt)
	for {
		remainder.Div(ac, bc)
		if remainder.IsZero() {
			g.Set(bc)
			return NewGaussianInt(bc.R, bc.I)
		}
		ac.Set(bc)
		bc.Set(remainder)
	}
}
