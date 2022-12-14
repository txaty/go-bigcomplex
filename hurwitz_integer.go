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

// HurwitzInt implements Hurwitz quaternion (or Hurwitz integer) a + bi + cj + dk
// The set of all Hurwitz quaternion is H = {a + bi + cj + dk | a, b, c, d are all integers or all half-integers}
// A mixture of integers and half-integers is excluded
// In the struct each scalar is twice the original scalar so that all the scalars can be stored using
// big integers for computation efficiency
type HurwitzInt struct {
	dblR *big.Int // r part doubled
	dblI *big.Int // i part doubled
	dblJ *big.Int // j part doubled
	dblK *big.Int // k part doubled
}

// Init initialize a Hurwitz integer
func (h *HurwitzInt) Init() *HurwitzInt {
	h.dblR = new(big.Int)
	h.dblI = new(big.Int)
	h.dblJ = new(big.Int)
	h.dblK = new(big.Int)
	return h
}

// String returns the string representation of the integral quaternion
func (h *HurwitzInt) String() string {
	rSign := h.dblR.Sign()
	iSign := h.dblI.Sign()
	jSign := h.dblJ.Sign()
	kSign := h.dblK.Sign()
	rABS := iPool.Get().(*big.Int).Abs(h.dblR)
	defer iPool.Put(rABS)
	iABS := iPool.Get().(*big.Int).Abs(h.dblI)
	defer iPool.Put(iABS)
	jABS := iPool.Get().(*big.Int).Abs(h.dblJ)
	defer iPool.Put(jABS)
	kABS := iPool.Get().(*big.Int).Abs(h.dblK)
	defer iPool.Put(kABS)
	if rSign == 0 && iSign == 0 && jSign == 0 && kSign == 0 {
		return "0"
	}
	res := ""
	if rABS.Cmp(big2) == 0 {
		if rSign < 0 {
			res += "-"
		}
		res += "1"
	} else {
		res += hiComposeString(0, iSign, rABS, "")
	}
	res += hiComposeString(rSign, iSign, iABS, "i")
	res += hiComposeString(iSign, jSign, jABS, "j")
	res += hiComposeString(jSign, kSign, kABS, "k")
	return res
}

func hiComposeString(lastSign, thisSign int, abs *big.Int, sign string) string {
	res := ""
	if lastSign != 0 && thisSign == 1 {
		res += "+"
	}
	if abs.Cmp(big1) == 0 {
		if thisSign == 1 {
			res += "0.5" + sign
		} else {
			res += "-0.5" + sign
		}
	} else if abs.Cmp(big2) == 0 {
		if thisSign == 1 {
			res += sign
		} else {
			res += "-" + sign
		}
	} else if abs.Sign() != 0 {
		opt := iPool.Get().(*big.Int)
		res += opt.Rsh(abs, 1).String()
		if abs.Bit(0) == 1 {
			res += ".5"
		}
		res += sign
	}
	return res
}

// NewHurwitzInt declares a new integral quaternion with the real, i, j, and k parts
// If isDouble is true, the arguments r, i, j, k are twice the original scalars
func NewHurwitzInt(r, i, j, k *big.Int, doubled bool) *HurwitzInt {
	if doubled {
		return &HurwitzInt{
			dblR: new(big.Int).Set(r),
			dblI: new(big.Int).Set(i),
			dblJ: new(big.Int).Set(j),
			dblK: new(big.Int).Set(k),
		}
	}
	return &HurwitzInt{
		dblR: new(big.Int).Lsh(r, 1),
		dblI: new(big.Int).Lsh(i, 1),
		dblJ: new(big.Int).Lsh(j, 1),
		dblK: new(big.Int).Lsh(k, 1),
	}
}

// Set sets the Hurwitz integer to the given Hurwitz integer
func (h *HurwitzInt) Set(a *HurwitzInt) *HurwitzInt {
	if h.dblR == nil {
		h.dblR = new(big.Int)
	}
	h.dblR.Set(a.dblR)
	if h.dblI == nil {
		h.dblI = new(big.Int)
	}
	h.dblI.Set(a.dblI)
	if h.dblJ == nil {
		h.dblJ = new(big.Int)
	}
	h.dblJ.Set(a.dblJ)
	if h.dblK == nil {
		h.dblK = new(big.Int)
	}
	h.dblK.Set(a.dblK)
	return h
}

// SetFloat set scalars of a Hurwitz integer by big float variables
//func (h *HurwitzInt) SetFloat(r, i, j, k *big.Float) *HurwitzInt {
//	panic("not implemented")
//}

// Val reveals value of a Hurwitz integer
func (h *HurwitzInt) Val() (r, i, j, k *big.Float) {
	r = new(big.Float).SetInt(h.dblR)
	r.Quo(r, big2f)
	i = new(big.Float).SetInt(h.dblI)
	i.Quo(i, big2f)
	j = new(big.Float).SetInt(h.dblJ)
	j.Quo(j, big2f)
	k = new(big.Float).SetInt(h.dblK)
	k.Quo(k, big2f)
	return
}

// ValInt reveals value of a Hurwitz integer in integer
func (h *HurwitzInt) ValInt() (r, i, j, k *big.Int) {
	rF, iF, jF, kF := h.Val()
	r = roundFloat(rF)
	i = roundFloat(iF)
	j = roundFloat(jF)
	k = roundFloat(kF)
	return
}

// Update updates the integral quaternion with the given real, i, j, and k parts
func (h *HurwitzInt) Update(r, i, j, k *big.Int, doubled bool) *HurwitzInt {
	if doubled {
		h.dblR = r
		h.dblI = i
		h.dblJ = j
		h.dblK = k
	} else {
		if h.dblR == nil {
			h.dblR = new(big.Int)
		}
		h.dblR.Lsh(r, 1)
		if h.dblI == nil {
			h.dblI = new(big.Int)
		}
		h.dblI.Lsh(i, 1)
		if h.dblJ == nil {
			h.dblJ = new(big.Int)
		}
		h.dblJ.Lsh(j, 1)
		if h.dblK == nil {
			h.dblK = new(big.Int)
		}
		h.dblK.Lsh(k, 1)
	}
	return h
}

// Zero sets the Hurwitz integer to zero
func (h *HurwitzInt) Zero() *HurwitzInt {
	h.dblR = big.NewInt(0)
	h.dblI = big.NewInt(0)
	h.dblJ = big.NewInt(0)
	h.dblK = big.NewInt(0)
	return h
}

// Add adds two integral quaternions
func (h *HurwitzInt) Add(a, b *HurwitzInt) *HurwitzInt {
	if h.dblR == nil {
		h.dblR = new(big.Int)
	}
	h.dblR.Add(a.dblR, b.dblR)
	if h.dblI == nil {
		h.dblI = new(big.Int)
	}
	h.dblI.Add(a.dblI, b.dblI)
	if h.dblJ == nil {
		h.dblJ = new(big.Int)
	}
	h.dblJ.Add(a.dblJ, b.dblJ)
	if h.dblK == nil {
		h.dblK = new(big.Int)
	}
	h.dblK.Add(a.dblK, b.dblK)
	return h
}

// Sub subtracts two integral quaternions
func (h *HurwitzInt) Sub(a, b *HurwitzInt) *HurwitzInt {
	if h.dblR == nil {
		h.dblR = new(big.Int)
	}
	h.dblR.Sub(a.dblR, b.dblR)
	if h.dblI == nil {
		h.dblI = new(big.Int)
	}
	h.dblI.Sub(a.dblI, b.dblI)
	if h.dblJ == nil {
		h.dblJ = new(big.Int)
	}
	h.dblJ.Sub(a.dblJ, b.dblJ)
	if h.dblK == nil {
		h.dblK = new(big.Int)
	}
	h.dblK.Sub(a.dblK, b.dblK)
	return h
}

// Conj obtains the conjugate of the original integral quaternion
func (h *HurwitzInt) Conj(origin *HurwitzInt) *HurwitzInt {
	if h.dblR == nil {
		h.dblR = new(big.Int)
	}
	h.dblR.Set(origin.dblR)
	if h.dblI == nil {
		h.dblI = new(big.Int)
	}
	h.dblI.Neg(origin.dblI)
	if h.dblJ == nil {
		h.dblJ = new(big.Int)
	}
	h.dblJ.Neg(origin.dblJ)
	if h.dblK == nil {
		h.dblK = new(big.Int)
	}
	h.dblK.Neg(origin.dblK)
	return h
}

// Norm obtains the norm of the integral quaternion
func (h *HurwitzInt) Norm() *big.Int {
	norm := new(big.Int).Mul(h.dblR, h.dblR)
	opt := iPool.Get().(*big.Int).Mul(h.dblI, h.dblI)
	defer iPool.Put(opt)
	norm.Add(norm, opt)
	opt.Mul(h.dblJ, h.dblJ)
	norm.Add(norm, opt)
	opt.Mul(h.dblK, h.dblK)
	norm.Add(norm, opt)
	norm.Rsh(norm, 2)
	return norm
}

// Copy copies the integral quaternion
func (h *HurwitzInt) Copy() *HurwitzInt {
	return NewHurwitzInt(h.dblR, h.dblI, h.dblJ, h.dblK, true)
}

// Prod returns the Hamilton product of two integral quaternions
// the product (a1 + b1j + c1k + d1)(a2 + b2j + c2k + d2) is determined by the products of the
// basis elements and the distributive law
func (h *HurwitzInt) Prod(a, b *HurwitzInt) *HurwitzInt {
	r, i, j, k := new(big.Int), new(big.Int), new(big.Int), new(big.Int)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	// 1 part
	r.Mul(a.dblR, b.dblR)
	r.Sub(r, opt.Mul(a.dblI, b.dblI))
	r.Sub(r, opt.Mul(a.dblJ, b.dblJ))
	r.Sub(r, opt.Mul(a.dblK, b.dblK))
	r.Rsh(r, 1)

	// i part
	i.Mul(a.dblR, b.dblI)
	i.Add(i, opt.Mul(a.dblI, b.dblR))
	i.Add(i, opt.Mul(a.dblJ, b.dblK))
	i.Sub(i, opt.Mul(a.dblK, b.dblJ))
	i.Rsh(i, 1)

	// j part
	j.Mul(a.dblR, b.dblJ)
	j.Sub(j, opt.Mul(a.dblI, b.dblK))
	j.Add(j, opt.Mul(a.dblJ, b.dblR))
	j.Add(j, opt.Mul(a.dblK, b.dblI))
	j.Rsh(j, 1)

	// k part
	k.Mul(a.dblR, b.dblK)
	k.Add(k, opt.Mul(a.dblI, b.dblJ))
	k.Sub(k, opt.Mul(a.dblJ, b.dblI))
	k.Add(k, opt.Mul(a.dblK, b.dblR))
	k.Rsh(k, 1)

	h.dblR, h.dblI, h.dblJ, h.dblK = r, i, j, k

	return h
}

// Div performs Euclidean division of two Hurwitz integers, i.e. a/b
// the remainder is stored in the Hurwitz integer that calls the method
// the quotient is returned as a new Hurwitz integer
func (h *HurwitzInt) Div(a, b *HurwitzInt) *HurwitzInt {
	ac := hiPool.Get().(*HurwitzInt)
	defer hiPool.Put(ac)
	ac = a.Copy()
	bc := hiPool.Get().(*HurwitzInt)
	defer hiPool.Put(bc)
	bc = b.Copy()

	bConj := hiPool.Get().(*HurwitzInt).Conj(bc)
	defer hiPool.Put(bConj)
	numerator := hiPool.Get().(*HurwitzInt).Prod(ac, bConj)
	defer hiPool.Put(numerator)
	denominator := hiPool.Get().(*HurwitzInt).Prod(bc, bConj)
	defer hiPool.Put(denominator)
	deFloat := fPool.Get().(*big.Float).SetInt(denominator.dblR)
	defer fPool.Put(deFloat)

	rScalar := fPool.Get().(*big.Float).SetInt(numerator.dblR)
	defer fPool.Put(rScalar)
	rScalar.Quo(rScalar, deFloat)
	iScalar := fPool.Get().(*big.Float).SetInt(numerator.dblI)
	defer fPool.Put(iScalar)
	iScalar.Quo(iScalar, deFloat)
	jScalar := fPool.Get().(*big.Float).SetInt(numerator.dblJ)
	defer fPool.Put(jScalar)
	jScalar.Quo(jScalar, deFloat)
	kScalar := fPool.Get().(*big.Float).SetInt(numerator.dblK)
	defer fPool.Put(kScalar)
	kScalar.Quo(kScalar, deFloat)

	rsInt := iPool.Get().(*big.Int)
	defer iPool.Put(rsInt)
	rsInt = roundFloat(rScalar)
	isInt := iPool.Get().(*big.Int)
	defer iPool.Put(isInt)
	isInt = roundFloat(iScalar)
	jsInt := iPool.Get().(*big.Int)
	defer iPool.Put(jsInt)
	jsInt = roundFloat(jScalar)
	ksInt := iPool.Get().(*big.Int)
	defer iPool.Put(ksInt)
	ksInt = roundFloat(kScalar)

	quotient := NewHurwitzInt(rsInt, isInt, jsInt, ksInt, false)
	opt := hiPool.Get().(*HurwitzInt)
	defer hiPool.Put(opt)
	h.Sub(ac, opt.Prod(quotient, bc))
	return quotient
}

// GCRD calculates the greatest common right-divisor of two Hurwitz integers using Euclidean algorithm
// The GCD is unique only up to multiplication by a unit (multiplication on the left in the case
// of a GCRD, and on the right in the case of a GCLD)
// the result is stored in the Hurwitz integer that calls the method and returned
func (h *HurwitzInt) GCRD(a, b *HurwitzInt) *HurwitzInt {
	ac := hiPool.Get().(*HurwitzInt).Set(a)
	defer hiPool.Put(ac)
	bc := hiPool.Get().(*HurwitzInt).Set(b)
	defer hiPool.Put(bc)

	if ac.CmpNorm(bc) < 0 {
		ac, bc = bc, ac
	}
	remainder := hiPool.Get().(*HurwitzInt)
	defer hiPool.Put(remainder)
	for {
		remainder.Div(ac, bc)
		if remainder.IsZero() {
			h.Set(bc)
			return new(HurwitzInt).Set(bc)
		}
		ac.Set(bc)
		bc.Set(remainder)
	}
}

// Equals checks if the two Hurwitz integers are equal
func (h *HurwitzInt) Equals(a *HurwitzInt) bool {
	return h.dblR.Cmp(a.dblR) == 0 &&
		h.dblI.Cmp(a.dblI) == 0 &&
		h.dblJ.Cmp(a.dblJ) == 0 &&
		h.dblK.Cmp(a.dblK) == 0
}

// IsZero returns true if the Hurwitz integer is zero
func (h *HurwitzInt) IsZero() bool {
	return h.dblR.Sign() == 0 &&
		h.dblI.Sign() == 0 &&
		h.dblJ.Sign() == 0 &&
		h.dblK.Sign() == 0
}

// CmpNorm compares the norm of two Hurwitz integers
func (h *HurwitzInt) CmpNorm(a *HurwitzInt) int {
	return h.Norm().Cmp(a.Norm())
}
