package crypto

import (
	"fmt"
	"math/big"
	"math/rand"
)

// GenKeyShares returns n keys from which only t of them are required to reconstruct the original
// secret. The keys are generated using Shamir Secret Sharing Scheme so each element of the returned
// array corresponds to an evaluation from x=1 up to x=n of the randomly generated polynomial f(x)
// of degree t-1, i.e: data[i] = (i+1, f(i+1)).
// Note: Polynomial evaluations and coefficients operate under an arithmetic finite field of size P.
// P is defined under primes.go and it is known to use exactly 257 for its representation.
func GenKeyShares(secret [32]byte, t, n int) ([][]byte, error) {
	if t < 2 {
		return nil, fmt.Errorf("unmet constraint t: %d >= 2", t)
	}
	if n < 3 {
		return nil, fmt.Errorf("unmet constraint n: %d >= 3", n)
	}
	if t > n {
		return nil, fmt.Errorf("unmet constraint t: %d <= n: %d", t, n)
	}

	coefficients, err := genRandomCoefficients(t - 1)
	if err != nil {
		return nil, err
	}

	key := big.NewInt(0)
	key.SetBytes(secret[:])
	key = key.Abs(key)
	key = key.Mod(key, P)

	coefficients = append(coefficients, key)
	fxs := evaluatePolynomial(coefficients, n)

	keys := make([][]byte, n)
	for i := range keys {
		fx := fxs[i]
		keys[i] = fx.Bytes()
	}

	return keys, nil
}

// genRandomCoefficients returns n numbers of 32 random bytes under Zp.
func genRandomCoefficients(n int) ([]*big.Int, error) {
	nums := make([]*big.Int, n)
	for i := 0; i < n; i++ {
		bytes := make([]byte, 32)
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}
		c := big.NewInt(0)
		c.SetBytes(bytes)
		c = c.Abs(c)
		c = c.Mod(c, P)
		nums[i] = c
	}
	return nums, nil
}

// evaluatePolynomial takes the coefficients of a polynomial of degree n=len(coefficients)-1
// in descending order, i.e: coefficients[0] = a_n and coefficients[n-1] = a_0, and computes the
// result of evaluating said polynomial for points [1, n]. Computation is done in linear time
// using Horner's method https://en.wikipedia.org/wiki/Horner%27s_method.
func evaluatePolynomial(coefficients []*big.Int, n int) []*big.Int {
	fxs := make([]*big.Int, n)
	for x := 1; x <= n; x++ {
		y := big.NewInt(0)
		y.Add(y, coefficients[0])

		for i := 1; i < len(coefficients); i++ {
			y.Mul(y, big.NewInt(int64(x)))
			y.Mod(y, P)
			y.Add(y, coefficients[i])
			y.Mod(y, P)
		}

		fxs[x-1] = y
	}
	return fxs
}

// Point is contains the evaluation of a polynomial at point X, i.e: f(X).
type Point struct {
	X  int
	Fx []byte
}

func GetKeyFromKeyShares(points []Point) ([]byte, error) {
	if len(points) < 2 {
		return []byte{}, fmt.Errorf("got %d, wants at least 2", len(points))
	}

	basis := lagrangeBasis(points)
	evs := make([]*big.Int, len(points))
	for i := range points {
		bytes := points[i].Fx
		fx := big.NewInt(0)
		fx.SetBytes(bytes)
		evs[i] = fx
	}
	root, err := findPolynomialRoot(basis, evs)
	if err != nil {
		return nil, err
	}
	return root.Bytes(), nil
}

// lagrangeBasis computes the basis required for the polynomial interpolation using Lagrange the
// polynomials. https://en.wikipedia.org/wiki/Lagrange_polynomial
func lagrangeBasis(points []Point) []*big.Int {
	basis := make([]*big.Int, len(points))
	for i := range points {
		pi0 := big.NewInt(0)
		numerator := big.NewInt(1)
		denominator := big.NewInt(1)

		// Calculate numerator
		for j := range points {
			if i == j {
				continue
			}
			currentFactor := big.NewInt(0)
			currentFactor.Sub(P, big.NewInt(int64(points[j].X)))
			numerator.Mul(numerator, currentFactor)
			numerator.Mod(numerator, P)
		}

		// Calculate denominator
		for j := range points {
			if i == j {
				continue
			}
			currentFactor := big.NewInt(0)
			xi := big.NewInt(int64(points[i].X))
			xj := big.NewInt(int64(points[j].X))
			inverseXj := big.NewInt(0)
			inverseXj.Sub(P, xj)
			currentFactor.Add(xi, inverseXj)
			currentFactor.Mod(currentFactor, P)
			denominator.Mul(denominator, currentFactor)
			denominator.Mod(denominator, P)
		}
		denominator.ModInverse(denominator, P)

		// Calculate division (numerator * denominator) mod p
		pi0.Mul(numerator, denominator)
		pi0.Mod(pi0, P)
		basis[i] = pi0
	}

	return basis
}

// findPolynomialRoot returns the polynomial interpolation at x=0 using the lagrange polynomial and
// basis.
func findPolynomialRoot(lagrangeBasis, polynomialEvaluations []*big.Int) (*big.Int, error) {
	if len(polynomialEvaluations) != len(lagrangeBasis) {
		return nil, fmt.Errorf("there must be as many lagrange basis as there are polynomial evaluations")
	}

	root := big.NewInt(0)
	for i := 0; i < len(polynomialEvaluations); i++ {
		currentAddend := big.NewInt(0)
		currentAddend.Mul(lagrangeBasis[i], polynomialEvaluations[i])
		currentAddend.Mod(currentAddend, P)
		root.Add(root, currentAddend)
		root.Mod(root, P)
	}
	return root, nil
}
