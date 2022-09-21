package complex

import "math/big"

// roundFloat rounds the given big float to the nearest big integer
func roundFloat(f *big.Float) *big.Int {
	if f.Sign() < 0 {
		f.Sub(f, rDelta)
	} else {
		f.Add(f, rDelta)
	}
	res, _ := f.Int(nil)
	return res
}
