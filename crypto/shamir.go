package crypto

import (
	"fmt"
	"math/big"
	"math/rand"
)

// evaluatePolynomial
// TODO: add ref to horners method & add proper docstring
func evaluatePolynomial(coefficients []*big.Int, K, x *big.Int) *big.Int {
	evaluation := big.NewInt(0)
	evaluation.Add(evaluation, coefficients[0])

	for i := 1; i < len(coefficients); i++ {
		evaluation.Mul(evaluation, x)
		evaluation.Mod(evaluation, P)
		evaluation.Add(evaluation, coefficients[i])
		evaluation.Mod(evaluation, P)
	}
	evaluation.Mul(evaluation, x)
	evaluation.Mod(evaluation, P)
	evaluation.Add(evaluation, K)
	evaluation.Mod(evaluation, P)

	return evaluation
}

// GenKeyShares generates n key shares (xi, P(xi)) where xi is [1...n] and P(xi) is a 256-bit integer
// t - 1 represents the degree of the polynomial P used to generate the n key shares
// which makes only t of them necessary to calculate K key used to encrypt the content
func GenKeyShares(secret [32]byte, t, n int) ([][32]byte, error) {
	if t < 2 {
		return nil, fmt.Errorf("unmet constraint t: %d >= 2", t)
	}
	if n < 3 {
		return nil, fmt.Errorf("unmet constraint n: %d >= 3", n)
	}
	if t > n {
		return nil, fmt.Errorf("unmet constraint t: %d <= n: %d", t, n)
	}

	coefficients := make([]*big.Int, t-1)
	for i := 0; i < t-1; i++ {
		randomCoefficient := make([]byte, 32)
		if _, err := rand.Read(randomCoefficient); err != nil {
			return nil, err
		}
		coefficient := big.NewInt(0)
		coefficient.SetBytes(randomCoefficient)
		coefficients[i] = coefficient
	}

	key := big.NewInt(0)
	key.SetBytes(secret[:])

	result := make([][32]byte, n)
	for i := 1; i <= n; i++ {
		evaluation := evaluatePolynomial(coefficients, key, big.NewInt(int64(i)))
		copy(result[i-1][:], evaluation.Bytes())
	}

	return result, nil
}

type Point struct {
	X  int
	Fx [32]byte
}

func GetKeyFromKeyShares(points []Point) ([32]byte, error) {
	if len(points) < 2 {
		return [32]byte{}, fmt.Errorf("got %d, wants at least 2", len(points))
	}

	lagrangeBasis := getLagrangeBasis(points)
	polynomialEvaluations := make([]*big.Int, len(points))
	for i, p := range points {
		bigInt := big.NewInt(0)
		bigInt.SetBytes(p.Fx[:])
		polynomialEvaluations[i] = bigInt
	}

	return findPolynomialRoot(lagrangeBasis, polynomialEvaluations)
}

func getLagrangeBasis(points []Point) []*big.Int {
	res := make([]*big.Int, len(points))
	for i := range points {
		pi0 := big.NewInt(0)
		numerator := big.NewInt(1)
		denominator := big.NewInt(1)

		// Calculate numerator
		for j := range points {
			currentFactor := big.NewInt(int64(-1 * points[j].X))
			numerator.Mul(numerator, currentFactor)
		}

		// Calculate denominator
		for j := range points {
			if i == j {
				continue
			}
			currentFactor := big.NewInt(0)
			xi := big.NewInt(int64(points[i].X))
			xj := big.NewInt(int64(points[j].X))
			currentFactor.Sub(xi, xj)
			currentFactor.Mod(currentFactor, P)
			denominator.Mul(denominator, currentFactor)
			denominator.Mod(denominator, P)
		}
		denominator.ModInverse(denominator, P)

		// Calculate division (numerator * denominator) mod p
		pi0.Mul(numerator, denominator)
		pi0.Mod(pi0, P)
		res[i] = pi0
	}

	return res
}

func findPolynomialRoot(lagrangeBasis, polynomialEvaluations []*big.Int) ([32]byte, error) {
	if len(polynomialEvaluations) != len(lagrangeBasis) {
		return [32]byte{}, fmt.Errorf("there must be as many lagrange basis as there are polynomial evaluations")
	}
	res := big.NewInt(0)
	for i := 0; i < len(polynomialEvaluations); i++ {
		res.Mul(lagrangeBasis[i], polynomialEvaluations[i])
		res.Mod(res, P)
	}
	var resBytes [32]byte
	copy(resBytes[:], res.Bytes())

	return resBytes, nil
}
