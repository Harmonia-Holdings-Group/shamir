package main

import (
	"fmt"
	"math/big"

	"github.com/pablotrinidad/shamir/crypto"
)

func main() {
	var key [32]byte
	key[30] += 2
	keyShares, err := crypto.GenKeyShares(key, 3, 5)
	points := make([]crypto.Point, 0, len(keyShares))
	for i := 0; i < len(keyShares); i++ {
		point := crypto.Point{X: i + 1, Fx: keyShares[i]}
		points = append(points, point)
	}

	if err != nil {
		panic(err)
	}

	derivedKey, err := crypto.GetKeyFromKeyShares(points)
	if err != nil {
		panic("")
	}
	OMG := big.NewInt(0)
	OMG.SetBytes(key[:])
	derivedKeyOMG := big.NewInt(0)
	derivedKeyOMG.SetBytes(derivedKey[:])

	fmt.Println()
	fmt.Println(OMG)
	fmt.Println(derivedKeyOMG)
}
