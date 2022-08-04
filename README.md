# Big Complex

[![Go Reference](https://pkg.go.dev/badge/github.com/tommytim0515/go-bigcomplex.svg)](https://pkg.go.dev/github.com/tommytim0515/go-bigcomplex)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommytim0515/go-bigcomplex)](https://goreportcard.com/report/github.com/tommytim0515/go-bigcomplex)
![Coverage](https://img.shields.io/badge/Coverage-19.2%25-red)

Big complex number calculation library for Go (with [math/big](https://pkg.go.dev/math/big)).

Currently, the library supports:

1. Gaussian
   integer <img src="http://www.sciweavers.org/tex2img.php?eq=Z%5Bi%5D%20%3D%20%5C%7Ba%20%2B%20bi%5C%20%7C%5C%20a%2C%20b%20%20%5Cin%20Z%5C%7D%2C%20where%5C%20i%5E2%20%3D%20-1&bc=White&fc=Black&im=jpg&fs=12&ff=arev&edit=0" align="center" border="0" alt="Z[i] = \{a + bi\ |\ a, b  \in Z\}, where\ i^2 = -1" width="333" height="21" />:
   Complex numbers whose real and imaginary parts are both integers.
2. Hurwitz
   quaternion <img src="http://www.sciweavers.org/tex2img.php?eq=H%20%3D%20%5C%7Ba%20%2B%20bi%20%2B%20cj%20%2B%20dk%20%20%5Cin%20%5Cmathbb%7BH%7D%5C%20%7C%5C%20a%2Cb%2Cc%2Cd%20%20%5Cin%20%5Cmathbb%7BZ%7D%5C%20or%5C%20%5Ca%2Cb%2Cc%2Cd%20%20%5Cin%20%5Cmathbb%7BZ%7D%20%2B%20%5Cfrac%7B1%7D%7B2%7D%5C%7D&bc=White&fc=Black&im=jpg&fs=12&ff=arev&edit=0" align="center" border="0" alt="H = \{a + bi + cj + dk  \in \mathbb{H}\ |\ a,b,c,d  \in \mathbb{Z}\ or\ \a,b,c,d  \in \mathbb{Z} + \frac{1}{2}\}" width="475" height="43" />:
   Quaternions whose components are either all integers or all half-integers (halves of odd integers; a mixture of
   integers and half-integers is excluded).

## Installation

```bash
go get github.com/tommytim0515/go-bigcomplex
```

## Examples

The usage is quite similar to Golang ```math/big``` package.

```go
package main

import (
   "fmt"
   "math/big"

   bc "github.com/tommytim0515/go-bigcomplex"
)

func main() {
   // Gaussian integer calculation
   g1 := bc.NewGaussianInt(big.NewInt(1), big.NewInt(2)) // 1 + 2i
   g2 := bc.NewGaussianInt(big.NewInt(5), big.NewInt(6)) // 5 + 6i
   gcd := new(bc.GaussianInt).GCD(g1, g2)
   fmt.Println(gcd)

   // Hurwitz integer calculation
   h1 := bc.NewHurwitzInt(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), false) // 1 + i + j + k
   h2 := bc.NewHurwitzInt(big.NewInt(3), big.NewInt(2), big.NewInt(2), big.NewInt(3), true) // 3/2 + i + j + 3k/2
   prod := new(bc.HurwitzInt).Pord(h1, h2)
   fmt.Println(prod)
}
````

## Why this Library?

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