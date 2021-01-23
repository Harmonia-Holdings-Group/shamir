package crypto

import (
	"math/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Test_KeyGeneration tests GenKeyShares and GetKeyFromKeyShares methods.
func Test_KeyGeneration(t *testing.T) {
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

func Test_EvaluatePoly(t *testing.T) {
	co := []*big.Int{
		big.NewInt(2),
		big.NewInt(3),
	}
	k := big.NewInt(6)
	n := 3
	points := make([]Point, n)

	// PASO 1
	// 		GENERA PUNTOS
	//		COPIA RAW BYTES into POINTS
	fmt.Println("EVALUATIONS:")
	for x := 1; x <= n; x++ {
		r := evaluatePolynomial(co, k, big.NewInt(int64(x)))
		fmt.Printf("f(%d) = %d\t (FROM POINTS: ", x, r)
		var y [32]byte
		copy(y[:], r.Bytes())
		points[x-1] = Point{X:  x, Fx: y}
		fmt.Printf("f(%d) = %d)\n", points[x-1].X, points[x-1].Fx)
	}

	fmt.Println(" ------------------------------------ ")

	fmt.Println("BASIS:")
	basis := getLagrangeBasis(points)
	for _, b := range basis {
		fmt.Println(b)
	}

	// PASO 3
	// 		LAS EVALUACIONES SE CREAN USANDO LOS BYTES VOLTEADOS
	//		SE MANDA A LLAMAD A findPolynomialRoot con esos puntos
	fmt.Println("ROOT:")
	evs := make([]*big.Int, len(points))
	for i := range points {
		fx := points[i].Fx
		y := big.NewInt(0)
		fx1 := make([]byte, 32)
		for i := range fx {
			fx1[31-i] = fx[i]
		}
		y.SetBytes(fx1[:])
		evs[i] = y
		fmt.Printf("y_%d: %d\n\tUsing bytes %v\n\n", i+1, y, fx1[:])
	}
	r, _ := findPolynomialRoot(basis, evs)
	fmt.Println(r)

}

// f(x) = 2x^2 + 3x + 6

// f(1) = 2 |  = 2 + 3 + 6 = (2+3)%11 + 6 = (5+6) % 11 = 11%11 = 0
// f(2) = 9 |  = 8 + 6 + 6 = (8 + 6) % 11 + 6 = (3 + 6) % 11 = 9
// f(3) = 0  |  = 9*2 + 3*3 + 6 = 18%11 + 9 + 6 = 7 + 9 + 6 = 16%11 + 6 = (5+6) % 11 = 0


// BASIS
//
// p1 = (-2 * -3) / [(1 + -2) (1 + -3)] = (9 * 8 % 11) / (1 + 9 % 11) (9) = (6) / [10 * 9] = 6 * inv(2) = 6 * 6 = 36%11 = 3
// p2 = (-1 * -3) / [(2 + -1) (2 + -3)] = (10 * 8 % 11) / (2+10)(2+8) = 3 / (1)(10) = 3 / 10 = 3 * inv(10) = 3*10 = 30%11 = 8
// p3 = (-1 * -2) / [(3 + -1) (3 + -2)] = (10 * 9 % 11) / (3+10)(2+9) = 2 / (2)(0) = 2 / 0 = 2 * inv(0) = 2 *


// p1 = (0-2/1-2) * (0-3/1-3) = (-2/-1) * (-3/-2) = (9 * inv10) (8 * inv9) = (9*10) (8*5) = 2 * 7 = 3
// p2 = (0-1/2-1) * (0-3/2-3) = (-1/1) * (-3/-1) = (10 * inv1) (8 * inv10) = (10*1) (8*10) = 10 * 3 = 8
// p3 = (0-1/3-1) * (0-2/3-2) = (-1/2) * (-2/1) = (10 * inv2) (9 * inv1) = (10*6) (9*1) = 5 * 9 = 1

// Lagrange form

//f(x) = y1p1 + y2p2 + y3p3
//f(0) = 0*3 + 9*8 + 0*1
//f(0) = 0 + 9*8 + 0
//f(0) = 6
//
