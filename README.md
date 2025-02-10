# Big Complex

[![Go Reference](https://pkg.go.dev/badge/github.com/txaty/go-bigcomplex.svg)](https://pkg.go.dev/github.com/txaty/go-bigcomplex)
[![Go Report Card](https://goreportcard.com/badge/github.com/txaty/go-bigcomplex)](https://goreportcard.com/report/github.com/txaty/go-bigcomplex)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/3a9bb5ff5cb64dcf83903ca998a9144d)](https://app.codacy.com/gh/txaty/go-merkletree/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

Big complex number calculation library for Go (with [math/big](https://pkg.go.dev/math/big)).

Currently, the library supports:

1. Gaussian integer, complex numbers whose real and imaginary parts are both integers:

$$
Z[i] = \{ a + bi \;|\; a, b \in \mathbb{Z} \}, \quad \text{where } i^2 = -1.
$$

2. Hurwitz quaternion, quaternions whose components are either all integers or all half-integers (halves of odd
   integers; a mixture of integers and half-integers is excluded):

$$
H = \{ a + bi + cj + dk \in \mathbb{H} \;|\; a, b, c, d \in \mathbb{Z} \; \text{or} \; b, c, d \in \mathbb{Z} + \frac{1}{2}  \}.
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

Fan fact: Golang has native complex number types: ```complex64``` and ```complex128```.

```go
c1 := complex(10, 11) // constructor init
c2 := 10 + 11i        // complex number init syntax

realPart := real(c1)    // gets real part
imagPart := imag(c1)    // gets imaginary part
```

```complex64``` represents ```float64```  real and imaginary data, and ```complex128``` represents ```float128``` real
and imaginary data.
They are easy to use, but unfortunately they are incapable for handling very large complex numbers.

For instance, in finding the Lagrange Four Square Sum of a very large integer (1792 bits in size) for cryptographic
range proof,
we need to compute the Greatest Common Divisor (GCD) of Gaussian integers and the Greatest Common Right Divisor of
Hurwitz integers. And the built-in complex number types can not handle such large numbers.

So I came up with the idea of building a library for large complex number calculation with Golang ```math/big```
package.