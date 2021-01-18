package main

import (
	"fmt"
	"syscall/js"

	"github.com/pablotrinidad/shamir/crypto"
)

func main() {
	fmt.Println("Go!")
	loadJSFuncs()
	<-make(chan bool)
}

func loadJSFuncs() {
	fns := map[string]func(this js.Value, args []js.Value) interface{}{
		"GoEncrypt": encrypt,
	}
	for k, v := range fns {
		js.Global().Set(k, js.FuncOf(v))
	}
}

func encrypt(_ js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return handleError(fmt.Errorf("got %d, want 2", len(args)))
	}

	gotKey := args[0].String()
	fileContent := make([]byte, args[1].Length())
	for i := 0; i < args[1].Length(); i++ {
		fileContent[i] = byte(args[1].Index(i).Int())
	}

	key, encrypted, err := crypto.Encrypt(gotKey, fileContent)
	if err != nil {
		return handleError(err)
	}
	return [2][]byte{key, encrypted}
}

func handleError(err error) string {
	return fmt.Sprintf("ERROR: %s", err)
}
