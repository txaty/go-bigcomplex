package complex

import "math/big"

// roundFloat rounds the given big float to the nearest big integer
func roundFloat(f *big.Float, res *big.Int) *big.Int {
	if f.Sign() < 0 {
		f.Sub(f, rDelta)
	} else {
		f.Add(f, rDelta)
	}
	if res == nil {
		res = new(big.Int)
	}
	f.Int(res)
	return res
}
