package main

import (
	"fmt"

	"github.com/pablotrinidad/shamir/crypto"
)

func main() {
	var key [32]byte
	key[13] = 1
	keyShares, err := crypto.GenKeyShares(key, 3, 5)
	points := make([]crypto.Point, 0, len(keyShares))
	fmt.Println("POINTS!!!!")
	for i := 0; i < len(keyShares); i++ {
		point := crypto.Point{X: i + 1, Fx: keyShares[i]}
		points = append(points, point)
		fmt.Printf("\t%v\n", keyShares[i])
	}

	if err != nil {
		panic(err)
	}

	derivedKey, err := crypto.GetKeyFromKeyShares(points)
	if err != nil {
		panic("")
	}

	fmt.Printf("ORIGINAL KEY: %v\n", key)
	fmt.Printf("DERIVED KEY: %v\n", derivedKey)
}
