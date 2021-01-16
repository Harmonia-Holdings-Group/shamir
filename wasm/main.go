package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"syscall/js"
)

func main() {
	fmt.Println("Go!")
	loadJSFuncs()
	<-make(chan bool)
}

func loadJSFuncs() {
	funcs := map[string]func(this js.Value, args []js.Value) interface{}{
		"genKey":         genKey,
		"encryptWithKey": encryptWithKey,
	}
	for k, v := range funcs {
		js.Global().Set(k, js.FuncOf(v))
	}
}

func genKey(_ js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("got %d, want 1", len(args))
	}
	passwd := args[0].String()
	h := sha256.New()
	h.Write([]byte(passwd))
	key := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(key)
}

func encryptWithKey(_ js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("got %d, want 2", len(args))
	}
	fileContent := make([]byte, args[0].Length())
	for i := 0; i < args[0].Length(); i++ {
		fileContent[i] = byte(args[0].Index(i).Int())
	}
	key, err := base64.StdEncoding.DecodeString(args[1].String())
	if err != nil {
		return err.Error()
	}
	fmt.Println(1)
	block, err := aes.NewCipher([]byte(key))
	fmt.Println(2)
	if err != nil {
		return err.Error()
	}
	fmt.Println(3)
	aesGCM, err := cipher.NewGCM(block)
	fmt.Println(4)
	nonce := make([]byte, aesGCM.NonceSize())
	fmt.Println(5)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	fmt.Println(6)
	content := aesGCM.Seal(nonce, nonce, []byte(fileContent), nil)
	fmt.Println(7)
	return base64.StdEncoding.EncodeToString(content)
}
