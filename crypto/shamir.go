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
	//
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

type Point struct {
	X  int
	Fx [32]byte
}

func GetKeyFromKeyShares(points []Point) ([32]byte, error) {
	var key [32]byte
	if len(points) < 2 {
		return key, fmt.Errorf("got %d, wants at least 2", len(points))
	}
	return points[0].Fx, nil
}
