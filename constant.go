package complex

import "math/big"

const (
	roundingDelta = 0.49
)

var (
	// big integer
	big1    = big.NewInt(1)
	bigNeg1 = big.NewInt(-1)

	// big float
	big2f = big.NewFloat(2)
	// delta for rounding big float to big int
	rDelta = big.NewFloat(roundingDelta)
)
