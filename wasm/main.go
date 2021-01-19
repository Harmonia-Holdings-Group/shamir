package main

import (
	"crypto/rand"
	"encoding/base64"
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
		"GoGenKeys": genKeys,
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
	return []interface{}{
		base64.StdEncoding.EncodeToString(key),
		base64.StdEncoding.EncodeToString(encrypted),
	}
}

func genKeys(_ js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return handleError(fmt.Errorf("got %d args, want 2", len(args)))
	}

	k := args[0].Int()
	n := args[1].Int()
	fmt.Printf("k: %s, n:%s\n", k, n)

	keys := make([]interface{}, n)
	for i := 0; i < n; i++ {
		randomData := make([]byte, 256/8)
		if _, err := rand.Read(randomData); err != nil {
			return handleError(fmt.Errorf("unexpected err; %v", err))
		}
		keys[i] = base64.StdEncoding.EncodeToString(randomData)
	}

	return keys
}

func handleError(err error) string {
	return fmt.Sprintf("ERROR: %s", err)
}
