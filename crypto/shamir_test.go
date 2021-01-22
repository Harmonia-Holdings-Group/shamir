package crypto

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Test_KeyGeneration tests GenKeyShares and GetKeyFromKeyShares methods.
func ATest_KeyGeneration(t *testing.T) {
	tests := []struct {
		name                    string
		inputKey                string // 256-bit base64 encoded key
		keyShares               int
		keyThreshold            int
		wantKeyMatch            bool
		useKeys                 int
		wantCreationError       bool
		wantReconstructionError bool
	}{
		{
			name:         "successful key generation (t: 5, n: 20) with exact number of keys",
			inputKey:     "7ezI6kdQ7fX4ekK8dRRSmOriR9SNKhV3OZ0k2CYnVTE=",
			keyShares:    20,
			keyThreshold: 5,
			useKeys:      5,
			wantKeyMatch: true,
		},
		{
			name:         "successful key generation (t: 5, n: 20) with extra keys",
			inputKey:     "NplIPa8T6fhgq+Mr6aGPekoaDeyjZymLLLSTzz1Me8I=",
			keyShares:    20,
			keyThreshold: 5,
			useKeys:      11,
			wantKeyMatch: true,
		},
		{
			name:         "successful key generation (t: 5, n: 20) with all keys",
			inputKey:     "+jbb2q+dZgHQ4czlTpUmOVc7cQ1XS3BFf4ENVs5lF9c=",
			keyShares:    20,
			keyThreshold: 10,
			useKeys:      20,
			wantKeyMatch: true,
		},
		{
			name:         "unsuccessful key generation (t: 5, n: 10) missing one key",
			inputKey:     "JBP7NwmwWTnwTPLpL30Il/wllvmtC4qeqFXHv+uq6JI=",
			keyShares:    10,
			keyThreshold: 5,
			useKeys:      4,
			wantKeyMatch: false,
		},
		{
			name:                    "unsuccessful key generation (t: 5, n: 10) without keys",
			inputKey:                "pgKBrF7WTdfmantSLeGpV4IyIFpIUAHshHk8DBAuypU=",
			keyShares:               10,
			keyThreshold:            5,
			useKeys:                 0,
			wantReconstructionError: true,
			wantKeyMatch:            false,
		},
		{
			name:              "unsuccessful key generation (t: 2, n: 2) invalid number of key shares",
			inputKey:          "zM5q1g8atTJ7Egd0zuQjYtr288p0xo6aQ6c1z641ROQ=",
			keyShares:         2,
			keyThreshold:      2,
			useKeys:           2,
			wantCreationError: true,
			wantKeyMatch:      false,
		},
		{
			name:              "unsuccessful key generation (t: 1, n: 3) invalid key threshold number",
			inputKey:          "Qlco35ne5pbmhu5eBbFOp2yCLMdNFD/IJ9sq+ZcdoRY=",
			keyShares:         3,
			keyThreshold:      1,
			useKeys:           3,
			wantCreationError: true,
			wantKeyMatch:      false,
		},
		{
			name:              "unsuccessful key generation (t: 4, n: 3) larger key threshold number",
			inputKey:          "Vzf3CgLH9b655mSoZ2HrMbr9VoVv7YrOIjflbe4IcKs=",
			keyShares:         3,
			keyThreshold:      4,
			useKeys:           3,
			wantCreationError: true,
			wantKeyMatch:      false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			keyBytes, err := base64.StdEncoding.DecodeString(test.inputKey)
			if err != nil {
				t.Fatalf("Unexpected error occurred while decoding key; %v", err)
			}
			var key [32]byte
			copy(key[:], keyBytes)

			keys, err := GenKeyShares(key, test.keyThreshold, test.keyShares)
			if err != nil && !test.wantCreationError {
				t.Fatalf("GenKeyShares(%s, t:%d, n:%d) returned unexpected error; %v", test.inputKey, test.keyThreshold, test.keyShares, err)
			}
			if err == nil && test.wantCreationError {
				t.Fatalf("GenKeyShares(%s, t:%d, n:%d) returned nil error, want errror", test.inputKey, test.keyThreshold, test.keyShares)
			}
			if test.wantCreationError {
				return
			}

			pointsMap := make(map[int]bool, test.useKeys)
			points := make([]Point, 0, test.useKeys)

			// Computed just for readable test output
			var encodedPoints strings.Builder
			encodedPoints.WriteString("\n")

			for len(points) < test.useKeys {
				x := rand.Intn(test.keyShares) + 1
				if ok := pointsMap[x]; ok {
					continue
				}
				pointsMap[x] = true
				points = append(points, Point{X: x, Fx: keys[x-1]})
				encodedPoints.WriteString(fmt.Sprintf("\t'%d-%s'\n", x, base64.StdEncoding.EncodeToString(keys[x-1][:])))
			}

			derivedKey, err := GetKeyFromKeyShares(points)
			if err != nil && !test.wantReconstructionError {
				t.Fatalf("GetKeyFromKeyShares(%s) returned unexpected error; %v", encodedPoints.String(), err)
			}
			if err == nil && test.wantReconstructionError {
				t.Fatalf("GetKeyFromKeyShares(%s) returned nil error, want error", encodedPoints.String())
			}
			if test.wantReconstructionError {
				return
			}

			decodedKey := base64.StdEncoding.EncodeToString(derivedKey[:])
			diff := cmp.Diff(test.inputKey, decodedKey)
			if diff != "" && test.wantKeyMatch {
				t.Errorf("GetKeyFromKeyShares(%s): %q, want %q, diff want->got: %s", encodedPoints.String(), decodedKey, test.inputKey, diff)
			}
			if diff == "" && !test.wantKeyMatch {
				t.Errorf("GetKeyFromKeyShares(%s) returned matching key %q when unexpected", encodedPoints.String(), test.inputKey)
			}
		})
	}
}
