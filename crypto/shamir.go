package crypto

import (
	"fmt"
	"math/big"
)

func GenKeys(secret *big.Int, t, n int) []*big.Int {
	// 1. Random selection of coeficcients
	// 2. Evaluate poly [1, n]
	return nil
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
