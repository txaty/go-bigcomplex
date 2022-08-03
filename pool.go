package complex

import (
	"math/big"
	"sync"
)

var (
	iPool = sync.Pool{
		New: func() interface{} { return new(big.Int) },
	}
	fPool = sync.Pool{
		New: func() interface{} { return new(big.Float) },
	}
	giPool = sync.Pool{
		New: func() interface{} { return new(GaussianInt) },
	}
	hiPool = sync.Pool{
		New: func() interface{} { return new(HurwitzInt) },
	}
)
