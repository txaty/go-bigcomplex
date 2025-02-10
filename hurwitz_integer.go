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

// HurwitzInt represents a Hurwitz quaternion (or Hurwitz integer) of the form
// a + bi + cj + dk, where the original scalars are stored doubled. (That is, each
// component here is twice the value of the corresponding original scalar.)
// This doubling lets us represent both integers and half‑integers without mixing.
type HurwitzInt struct {
	dblR *big.Int // Doubled real part
	dblI *big.Int // Doubled i part
	dblJ *big.Int // Doubled j part
	dblK *big.Int // Doubled k part
}

// Init initializes a Hurwitz integer by allocating the internal big.Int values.
func (h *HurwitzInt) Init() *HurwitzInt {
	h.dblR = new(big.Int)
	h.dblI = new(big.Int)
	h.dblJ = new(big.Int)
	h.dblK = new(big.Int)
	return h
}

// String returns the string representation of the Hurwitz integer.
func (h *HurwitzInt) String() string {
	rSign := h.dblR.Sign()
	iSign := h.dblI.Sign()
	jSign := h.dblJ.Sign()
	kSign := h.dblK.Sign()

	// Obtain absolute values using the pool.
	rABS := iPool.Get().(*big.Int).Abs(h.dblR)
	defer iPool.Put(rABS)
	iABS := iPool.Get().(*big.Int).Abs(h.dblI)
	defer iPool.Put(iABS)
	jABS := iPool.Get().(*big.Int).Abs(h.dblJ)
	defer iPool.Put(jABS)
	kABS := iPool.Get().(*big.Int).Abs(h.dblK)
	defer iPool.Put(kABS)

	// If all components are zero, return "0".
	if rSign == 0 && iSign == 0 && jSign == 0 && kSign == 0 {
		return "0"
	}

	res := ""
	// Compose the real part.
	if rABS.Cmp(big2) == 0 {
		if rSign < 0 {
			res += "-"
		}
		res += "1"
	} else {
		res += hiComposeString(0, rSign, rABS, "")
	}
	// Compose the i, j, and k parts.
	res += hiComposeString(rSign, iSign, iABS, "i")
	res += hiComposeString(iSign, jSign, jABS, "j")
	res += hiComposeString(jSign, kSign, kABS, "k")
	return res
}

// hiComposeString is a helper function for composing a single component of the string.
// lastSign is the sign of the previous component; thisSign is the sign of the current component.
func hiComposeString(lastSign, thisSign int, abs *big.Int, suffix string) string {
	res := ""
	if lastSign != 0 && thisSign == 1 {
		res += "+"
	}
	if abs.Cmp(big1) == 0 {
		if thisSign == 1 {
			res += "0.5" + suffix
		} else {
			res += "-0.5" + suffix
		}
	} else if abs.Cmp(big2) == 0 {
		if thisSign == 1 {
			res += suffix
		} else {
			res += "-" + suffix
		}
	} else if abs.Sign() != 0 {
		opt := iPool.Get().(*big.Int)
		opt.Rsh(abs, 1)
		res += opt.String()
		if abs.Bit(0) == 1 {
			res += ".5"
		}
		res += suffix
		iPool.Put(opt)
	}
	return res
}

// NewHurwitzInt creates a new Hurwitz integer given the components.
// If doubled is true, the provided values are assumed to be already doubled.
func NewHurwitzInt(r, i, j, k *big.Int, doubled bool) *HurwitzInt {
	if doubled {
		return &HurwitzInt{
			dblR: new(big.Int).Set(r),
			dblI: new(big.Int).Set(i),
			dblJ: new(big.Int).Set(j),
			dblK: new(big.Int).Set(k),
		}
	}
	// Otherwise, shift left by 1 (i.e., multiply by 2).
	return &HurwitzInt{
		dblR: new(big.Int).Lsh(r, 1),
		dblI: new(big.Int).Lsh(i, 1),
		dblJ: new(big.Int).Lsh(j, 1),
		dblK: new(big.Int).Lsh(k, 1),
	}
}

// Set assigns the value of another Hurwitz integer to this one.
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

// Val returns the value of the Hurwitz integer as four big.Float values (dividing by 2).
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

// ValInt returns the value of the Hurwitz integer as four big.Int values by rounding.
func (h *HurwitzInt) ValInt() (r, i, j, k *big.Int) {
	rF, iF, jF, kF := h.Val()
	r = roundFloat(rF)
	i = roundFloat(iF)
	j = roundFloat(jF)
	k = roundFloat(kF)
	return
}

// Update sets the components of the Hurwitz integer. If doubled is false,
// the provided values are shifted left by 1.
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

// Zero sets the Hurwitz integer to zero.
func (h *HurwitzInt) Zero() *HurwitzInt {
	h.dblR = big.NewInt(0)
	h.dblI = big.NewInt(0)
	h.dblJ = big.NewInt(0)
	h.dblK = big.NewInt(0)
	return h
}

// Add computes the sum of two Hurwitz integers and stores the result in the receiver.
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

// Sub subtracts one Hurwitz integer from another and stores the result in the receiver.
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

// Conj computes the conjugate of the Hurwitz integer (negating the imaginary parts)
// and stores it in the receiver.
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

// Norm returns the norm of the Hurwitz integer.
// The computation is: (dblR^2 + dblI^2 + dblJ^2 + dblK^2) >> 2.
func (h *HurwitzInt) Norm() *big.Int {
	norm := new(big.Int).Mul(h.dblR, h.dblR)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	opt.Mul(h.dblI, h.dblI)
	norm.Add(norm, opt)
	opt.Mul(h.dblJ, h.dblJ)
	norm.Add(norm, opt)
	opt.Mul(h.dblK, h.dblK)
	norm.Add(norm, opt)
	norm.Rsh(norm, 2)
	return norm
}

// Copy returns a deep copy of the Hurwitz integer.
func (h *HurwitzInt) Copy() *HurwitzInt {
	return NewHurwitzInt(h.dblR, h.dblI, h.dblJ, h.dblK, true)
}

// Prod computes the Hamilton product (quaternion multiplication) of two Hurwitz integers,
// storing the result in the receiver.
func (h *HurwitzInt) Prod(a, b *HurwitzInt) *HurwitzInt {
	// Temporary variables for each component.
	r := new(big.Int)
	i := new(big.Int)
	j := new(big.Int)
	k := new(big.Int)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)

	// Compute the real component:
	// r = (a.dblR*b.dblR - a.dblI*b.dblI - a.dblJ*b.dblJ - a.dblK*b.dblK) >> 1
	r.Mul(a.dblR, b.dblR)
	r.Sub(r, opt.Mul(a.dblI, b.dblI))
	r.Sub(r, opt.Mul(a.dblJ, b.dblJ))
	r.Sub(r, opt.Mul(a.dblK, b.dblK))
	r.Rsh(r, 1)

	// Compute the i component:
	// i = (a.dblR*b.dblI + a.dblI*b.dblR + a.dblJ*b.dblK - a.dblK*b.dblJ) >> 1
	i.Mul(a.dblR, b.dblI)
	i.Add(i, opt.Mul(a.dblI, b.dblR))
	i.Add(i, opt.Mul(a.dblJ, b.dblK))
	i.Sub(i, opt.Mul(a.dblK, b.dblJ))
	i.Rsh(i, 1)

	// Compute the j component:
	// j = (a.dblR*b.dblJ - a.dblI*b.dblK + a.dblJ*b.dblR + a.dblK*b.dblI) >> 1
	j.Mul(a.dblR, b.dblJ)
	j.Sub(j, opt.Mul(a.dblI, b.dblK))
	j.Add(j, opt.Mul(a.dblJ, b.dblR))
	j.Add(j, opt.Mul(a.dblK, b.dblI))
	j.Rsh(j, 1)

	// Compute the k component:
	// k = (a.dblR*b.dblK + a.dblI*b.dblJ - a.dblJ*b.dblI + a.dblK*b.dblR) >> 1
	k.Mul(a.dblR, b.dblK)
	k.Add(k, opt.Mul(a.dblI, b.dblJ))
	k.Sub(k, opt.Mul(a.dblJ, b.dblI))
	k.Add(k, opt.Mul(a.dblK, b.dblR))
	k.Rsh(k, 1)

	h.dblR, h.dblI, h.dblJ, h.dblK = r, i, j, k
	return h
}

// Div performs Euclidean division of two Hurwitz integers (a / b).
// The remainder is stored in the receiver and the quotient is returned as a new Hurwitz integer.
func (h *HurwitzInt) Div(a, b *HurwitzInt) *HurwitzInt {
	// Make copies of the operands.
	ac := a.Copy()
	bc := b.Copy()

	// Compute the conjugate of bc.
	bConj := new(HurwitzInt).Conj(bc)
	// Numerator = a * conjugate(b)
	numerator := new(HurwitzInt).Prod(ac, bConj)
	// Denominator = b * conjugate(b)
	denominator := new(HurwitzInt).Prod(bc, bConj)

	// Use the real part of the denominator for the division.
	deFloat := fPool.Get().(*big.Float).SetInt(denominator.dblR)
	defer fPool.Put(deFloat)

	// Compute each component of the quotient as a float.
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

	// Round the computed float values to the nearest integers.
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

	// Create the quotient. Note: since the Hurwitz integer stores doubled values,
	// we pass false for 'doubled' to let NewHurwitzInt adjust appropriately.
	quotient := NewHurwitzInt(rsInt, isInt, jsInt, ksInt, false)

	// Compute the remainder: remainder = a - (quotient * b)
	opt := new(HurwitzInt).Prod(quotient, bc)
	h.Sub(ac, opt)
	return quotient
}

// GCRD computes the greatest common right-divisor (GCRD) of two Hurwitz integers
// using the Euclidean algorithm. (The result is unique up to multiplication by a unit.)
// The result is stored in the receiver and also returned as a new Hurwitz integer.
func (h *HurwitzInt) GCRD(a, b *HurwitzInt) *HurwitzInt {
	ac := new(HurwitzInt).Set(a)
	bc := new(HurwitzInt).Set(b)

	if ac.CmpNorm(bc) < 0 {
		ac, bc = bc, ac
	}

	remainder := new(HurwitzInt)
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

// Equals returns true if the Hurwitz integer is equal to the provided Hurwitz integer.
func (h *HurwitzInt) Equals(a *HurwitzInt) bool {
	return h.dblR.Cmp(a.dblR) == 0 &&
		h.dblI.Cmp(a.dblI) == 0 &&
		h.dblJ.Cmp(a.dblJ) == 0 &&
		h.dblK.Cmp(a.dblK) == 0
}

// IsZero returns true if the Hurwitz integer is zero.
func (h *HurwitzInt) IsZero() bool {
	return h.dblR.Sign() == 0 &&
		h.dblI.Sign() == 0 &&
		h.dblJ.Sign() == 0 &&
		h.dblK.Sign() == 0
}

// CmpNorm compares the norms of two Hurwitz integers.
func (h *HurwitzInt) CmpNorm(a *HurwitzInt) int {
	return h.Norm().Cmp(a.Norm())
}
