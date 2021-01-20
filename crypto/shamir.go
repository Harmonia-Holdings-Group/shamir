package crypto

import (
	"fmt"
	"math/big"
	"math/rand"
)

// Evaluates polynomial of degree len(coefficients)
func evaluatePolynomial(coefficients []*big.Int, K *big.Int, x int) *big.Int {
	xBigInt := big.NewInt(int64(x))
	evaluation := coefficients[0]

	i := 1
	for ; i < len(coefficients); i++ {
		evaluation.Mul(evaluation, xBigInt)
		evaluation.Mod(evaluation, P)
		evaluation.Add(evaluation, coefficients[i])
		evaluation.Mod(evaluation, P)
	}
	evaluation.Mul(evaluation, xBigInt)
	evaluation.Mod(evaluation, P)
	evaluation.Add(evaluation, K)
	evaluation.Mod(evaluation, P)

	return evaluation
}

// Generates n key shares (xi, P(xi)) where xi is [1...n] and P(xi) is a 256-bit integer
// t - 1 represents the degree of the polynomial P used to generate the n key shares
// which makes only t of them necessary to calculate K key used to encrypt the content
func GenKeyShares(secret [32]byte, t, n int) ([][32]byte, error) {
	if t < 2 {
		return nil, fmt.Errorf("Unmet constraint t: %d >= 2", t)
	}

	if n < 3 {
		return nil, fmt.Errorf("Unmet constraint n: %d >= 3", n)
	}

	if t > n {
		return nil, fmt.Errorf("Unmet constraint t: %d <= n: %d", t, n)
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
		evaluation := evaluatePolynomial(coefficients, key, i)
		copy(result[i-1][:], evaluation.Bytes())
	}

	return result, nil
}

func getLagrangeBasis(points []Point) ([]*big.Int, error) {
	res := make([]*big.Int, len(points))
	for i, p := range points {
		Pi0 := big.NewInt(0)
		numerator := big.NewInt(1)
		denominator := big.NewInt(1)

		// Calculate numerator
		for j, _ := range points {
			currentFactor := big.NewInt(int64(-1 * points[j].X))
			numerator.Mul(numerator, currentFactor)
		}

		// Calculate denominator
		for j, _ := range points {
			if i == j {
				continue
			}
			currentFactor := big.NewInt(0)
			Xi := big.NewInt(int64(points[i].X))
			Xj := big.NewInt(int64(points[j].X))
			currentFactor.Sub(Xi, Xj)
			currentFactor.Mod(currentFactor, P)
			denominator.Mul(denominator, currentFactor)
			denominator.Mod(denominator, P)
		}
		denominator.ModInverse(denominator, P)

		// Calculate division (numerator * denominator) mod p
		Pi0.Mul(numerator, denominator)
		Pi0.Mod(Pi0, P)
	}

	return res, nil
}

func findPolynomialRoot(lagrangeBasis, polynomialEvaluations []*big.Int) ([32]byte, error) {
	if len(polynomialEvaluations) != len(lagrangeBasis) {
		return [32]byte{}, fmt.Errorf("There must be as many lagrange basis as there are polynomial evaluations")
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

type Point struct {
	X  int
	Fx [32]byte
}

func GetKeyFromKeyShares(points []Point) ([32]byte, error) {
	var key [32]byte
	if len(points) < 2 {
		return key, fmt.Errorf("got %d, wants at least 2", len(points))
	}

	lagrangeBasis, e := getLagrangeBasis(points)
	if e != nil {
		return [32]byte{}, fmt.Errorf("")
	}

	polynomialEvaluations := make([]*big.Int, len(points))
	for i, p := range points {
		bigInt := big.NewInt(0)
		bigInt.SetBytes(p.Fx[:])
		polynomialEvaluations[i] = bigInt
	}

	key, err := findPolynomialRoot(lagrangeBasis, polynomialEvaluations)
	if err != nil {
		return [32]byte{}, err
	}

	return key, nil
}
