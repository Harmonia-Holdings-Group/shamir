package crypto

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Test_LargePrimeGeneration verifies the prime factorization produces the
// desired number.
func Test_LargePrimeGeneration(t *testing.T) {
	wantPrime := "208351617316091241234326746312124448251235562226470491514186331217050270460481"
	gotPrime := fmt.Sprintf("%s", P)
	if diff := cmp.Diff(wantPrime, gotPrime); diff != "" {
		t.Errorf("Prime generation failed, diff want->got:\n%s", diff)
	}
}
