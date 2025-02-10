# Big Complex

[![Go Reference](https://pkg.go.dev/badge/github.com/txaty/go-bigcomplex.svg)](https://pkg.go.dev/github.com/txaty/go-bigcomplex)
[![Go Report Card](https://goreportcard.com/badge/github.com/txaty/go-bigcomplex)](https://goreportcard.com/report/github.com/txaty/go-bigcomplex)
[![codecov](https://codecov.io/gh/txaty/go-bigcomplex/graph/badge.svg?token=LPW23PAEH8)](https://codecov.io/gh/txaty/go-bigcomplex)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/f149e51e0475464d843477adba68b577)](https://app.codacy.com/gh/txaty/go-bigcomplex/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

Big complex number calculation library in Go (with [math/big](https://pkg.go.dev/math/big)).

Currently, the library supports:

1. **Gaussian Integers**  
   Complex numbers whose real and imaginary parts are both integers:
   $$
   Z[i] = \{ a + bi \ |\ a, b \in \mathbb{Z} \}, \quad \text{where } i^2 = -1.
   $$

2. **Hurwitz Quaternions**  
   Quaternions whose components are either all integers or all half‑integers (half‑integers being halves of odd integers; mixing integers and half‑integers is not allowed):
   $$
   H = \{ a + bi + cj + dk \in \mathbb{H} \ |\ a, b, c, d \in \mathbb{Z} \ \text{or} \ b, c, d \in \mathbb{Z} + \frac{1}{2}  \}.
   $$

## Installation

```bash
go get -u github.com/txaty/go-bigcomplex
```

## Examples

The usage is quite similar to Golang ```math/big``` package.

```go
package main

import (
   "fmt"
   "math/big"

   bc "github.com/txaty/go-bigcomplex"
)

func main() {
   // Gaussian integer calculation
   g1 := bc.NewGaussianInt(big.NewInt(5), big.NewInt(6)) // 5 + 6i
   g2 := bc.NewGaussianInt(big.NewInt(1), big.NewInt(2)) // 1 + 2i
   div := new(bc.GaussianInt).Div(g2, g1)
   fmt.Println(div) // 3 - i
   gcd := new(bc.GaussianInt).GCD(g1, g2)
   fmt.Println(gcd) // i

   // Hurwitz integer calculation
   // 1 + i + j + k
   h1 := bc.NewHurwitzInt(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), false)
   // 3/2 + i + j + 3k/2
   h2 := bc.NewHurwitzInt(big.NewInt(3), big.NewInt(2), big.NewInt(2), big.NewInt(3), true)
   prod := new(bc.HurwitzInt).Prod(h1, h2)
   fmt.Println(prod) // 2 + 3i + 2j + 3k
}
````

## Why This Library?

Fan fact: Golang has native complex number types: `complex64` and `complex128`.

```go
c1 := complex(10, 11) // Using the complex constructor
c2 := 10 + 11i        // Using literal syntax

realPart := real(c1)  // Retrieves the real part
imagPart := imag(c1)  // Retrieves the imaginary part
```

However, `complex64` (composed of two `float32` values) and `complex128` (composed of two `float64` values) are limited to fixed‑precision arithmetic and cannot handle very large numbers.
For example, in finding the Lagrange Four Square Sum of a very large integer (1792 bits in size) for cryptographic range proof, we need to compute the Greatest Common Divisor (GCD) of Gaussian integers and the Greatest Common Right Divisor of Hurwitz integers. And the built-in complex number types cannot handle such large numbers.

This motivated the development of Big Complex: a library for large‑scale complex number calculations using Go’s math/big package.

## License

This project is licensed under the MIT License.