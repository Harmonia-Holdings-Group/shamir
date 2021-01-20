package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
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
		"GoEncrypt":             encrypt,
		"GoGenKeys":             genKeys,
		"GoGetKeyFromKeyShares": getKeyFromKeyShares,
		"GoDecrypt":             decrypt,
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
	fmt.Printf("Go: [Content to encrypt] %v\n", fileContent)
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
	fmt.Printf("genKeys(k: %d, n:%d)\n", k, n)

	keys := make([]interface{}, n)
	for i := 0; i < n; i++ {
		randomData := make([]byte, 256/8)
		if _, err := rand.Read(randomData); err != nil {
			return handleError(fmt.Errorf("unexpected err; %v", err))
		}
		keys[i] = fmt.Sprintf("%d-%s", i+1, base64.StdEncoding.EncodeToString(randomData))
	}

	return keys
}

func getKeyFromKeyShares(_ js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return handleError(fmt.Errorf("got %d args, want 1", len(args)))
	}
	if args[0].Length() < 2 {
		return handleError(fmt.Errorf("cannot derive keys from less that 2 key shars"))
	}
	keys := make([]crypto.Point, args[0].Length())
	for i := 0; i < args[0].Length(); i++ {
		keyStr := args[0].Index(i).String()
		fmt.Printf("Go: %s\n", keyStr)
		key := strings.SplitN(keyStr, "-", 2)
		if len(key) != 2 {
			return handleError(fmt.Errorf("got invalid key format, wants x-f(x) for key %q", keyStr))
		}
		x, err := strconv.Atoi(key[0])
		if err != nil {
			return handleError(fmt.Errorf("%s is not an int in key %q", key[0], keyStr))
		}
		decoded, err := base64.StdEncoding.DecodeString(key[1])
		if err != nil {
			return handleError(fmt.Errorf("failed decoding key; %v", err))
		}
		var fx [32]byte
		copy(fx[:], decoded)
		keys[i] = crypto.Point{X: x, Fx: fx}
	}
	key, err := crypto.GetKeyFromKeyShares(keys)
	if err != nil {
		return handleError(fmt.Errorf("failed obtaining master key; %v", err))
	}
	return base64.StdEncoding.EncodeToString(key[:])
}

func decrypt(_ js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return handleError(fmt.Errorf("got %d, want 2", len(args)))
	}
	key, err := base64.StdEncoding.DecodeString(args[0].String())
	if err != nil {
		return handleError(fmt.Errorf("failed decoding key; %v", err))
	}

	fileContent := make([]byte, args[1].Length())
	for i := 0; i < args[1].Length(); i++ {
		fileContent[i] = byte(args[1].Index(i).Int())
	}
	fmt.Printf("Go [Content to decrypt] %v\n", fileContent)

	content, err := crypto.Decrypt(key, fileContent)
	if err != nil {
		return handleError(fmt.Errorf("failed decrypting file; %v", err))
	}
	fmt.Printf("Go: [Decrypted content] %v\n", content)
	return base64.StdEncoding.EncodeToString(content)
}

func handleError(err error) string {
	return fmt.Sprintf("ERROR: %s", err)
}
