package main

import (
	"fmt"

	"github.com/pablotrinidad/shamir/crypto"
)

func main() {
	var key [32]byte
	keyShares, err := crypto.GenKeyShares(key, 3, 5)

	if err != nil {
		panic(err)
	}

	fmt.Println(keyShares)
}
