package crypto

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_EncryptDecrypt(t *testing.T) {
	tests := []struct {
		name             string
		encryptionKey    string
		content          []byte
		decryptionKey    []byte
		wantDecryptError bool
	}{
		{
			name:          "successful encrypt-decrypt round",
			encryptionKey: "a super secret key 🤫!!",
			content:       []byte(`this a TOP-SECRET, no one should know about this`),
		},
		{
			name:             "unsuccessful encrypt-decrypt round (using bytes rep of encryption key)",
			encryptionKey:    "a super secret key 🤫!!",
			decryptionKey:    []byte(`a super secret key 🤫!!`), // Since a 256-bit key is generated from input key
			wantDecryptError: true,
			content:          []byte(`this a TOP-SECRET, no one should know about this`),
		},
		{
			name:             "unsuccessful encrypt-decrypt round (empty key)",
			encryptionKey:    "a super secret key 🤫!!",
			decryptionKey:    []byte{},
			wantDecryptError: true,
			content:          []byte(`this a TOP-SECRET, no one should know about this`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			derivedKey, cipher, err := Encrypt(test.encryptionKey, test.content)
			if err != nil {
				t.Fatalf("Encrypt(%s, %s) returned unexpected error; %v", test.encryptionKey, test.content, err)
			}
			var decryptionKey []byte
			if test.decryptionKey != nil {
				decryptionKey = test.decryptionKey
			} else {
				decryptionKey = derivedKey
			}
			gotContent, err := Decrypt(decryptionKey, cipher)
			if test.wantDecryptError {
				if err == nil {
					t.Errorf("Decyrpt(%s, cipher) returned nil error, expected error", decryptionKey)
				}
				return
			}
			if err != nil {
				t.Fatalf("Decrypt(%s, cipher) returned unexepcted error; %v", decryptionKey, err)
			}
			if diff := cmp.Diff(gotContent, test.content); diff != "" {
				t.Errorf("Decrypt(%s, cipher) malformed original content", derivedKey)
			}
		})
	}
}
