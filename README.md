# Shamir Secret Sharing Scheme

**A Golang-WebAssembly implementation.**


[Shamir's Secret Sharing Scheme](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing) (**_SSSS_**)
is a cryptographic algorithm used to share a secret into multiple parts. To reconstruct the original
secret, a minimum number of parts is required.

This project exposes a web interface for encrypting and decrypting files using
[AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) and deriving a set of keys using
**_SSSS_**.

## How it works:

1. Users must enter a human readable password `p` and select a file to encrypt.
2. A key `k` of 256 bits is generated from `p` using the [SHA-256](https://en.wikipedia.org/wiki/SHA-2)
   hash function, i.e: `k = SHA256(p)`.
3. The file is then encrypted using `k` and [**AES**](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard).
4. A number `n` of key shares is selected as well as a number `t` of minimum number or required keys.
5. `n` keys are generated using **_SSSS_**, during this process `k` is considered the secret
   (as is the value required to successfully decrypt the file).
1. Given at least `t` keys, the secret `k` is recovered an then used for decryption.

## Details
### Encryption
1. Based on an user-provided master password, a key **K** is generated using **SHA-256**.
2. Content is encrypted using **AES-256** with **K** as key.
3. A **t - 1** degree polynomial is randomly generated that will later be used to generate **n** key shares by randomly taking **xi** and its evaluation against the polynomial **P(xi)**. Key shares are pairs **(xi, P(xi))** that can be used to recover **K** by evaluating **P(0)** using _Horner's Method_.

### Decryption
With **t** of the **n** key shares generated, the **K** key used to encrypt the content, can be recovered. 

1. Using _Horner's Method_ of polynomial evaluation, we evaluate **P(0)**, thus recovering **K**.
2. Using **K** we can decrypt the content.

## The code

### Go

A package (**`crypto`**) and executable (**wasm**) are included in this repository.

* The package exposes an API for encrypting, decrypting, generating keys, and reconstructing the
  secret using a set of shared keys.
* The executable is the bridge between the browser's runtime and the compiled WebAssembly binary.

To run the package tests use `go test .` under the `crypto/` folder (**93.5% coverage**).

To build the WebAssembly binary use: `GOOS=js GOARCH=wasm go build -o ../ui/public/wasm/main.wasm`
under the `wasm/` folder. If the `ui/public/wasm/main.wasm` file is not present, you can copy it from
`GOROOT/misc/wasm/wasm_exec.js`. For more info about WebAssembly and Golang, visit
[**this wiki**](https://github.com/golang/go/wiki/WebAssembly).

### UI (ReactJS)

A basic UI was created for interacting with the input fields and communicating with the WebAssembly
binary.

Use `npm start` to run the development server and `npm build` to create the production bundle.
Notice that this process does not compile the Go source into the WebAssembly binary, nor does it
include the other required WASM runtime files.