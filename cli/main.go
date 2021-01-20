package main

import (
	"fmt"

	"github.com/pablotrinidad/shamir/crypto"
)

func main() {
	fmt.Println("Hola")
	fmt.Println(crypto.P)

	var key [32]byte
	keyShares, err := crypto.GenKeyShares(key, 2, 3)

	if err != nil {
		panic(err)
	}

	fmt.Println(keyShares)
}
