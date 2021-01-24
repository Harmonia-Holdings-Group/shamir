package crypto

import (
	"math/big"
)

// P is a 257 bits prime number with decimal representation:
// 208351617316091241234326746312124448251235562226470491514186331217050270460481.
// https://primes.utm.edu/curios/page.php?number_id=3746
var P *big.Int

func init() {
	P = big.NewInt(257)
	// pFactorization is the fundamental factorization of P:
	// 312^2 + 312^3 + 312^5 + 312^7 + 312^11 + 312^13 + 312^17 + 312^19 + 312^23 + 312^29 + 312^31 + 1
	//pFactorization := [12][2]float64{
	//	{312, 2},
	//	{312, 3},
	//	{312, 5},
	//	{312, 7},
	//	{312, 11},
	//	{312, 13},
	//	{312, 17},
	//	{312, 19},
	//	{312, 23},
	//	{312, 29},
	//	{312, 31},
	//	{1, 1},
	//}
	//for _, f := range pFactorization {
	//	currentFactor := big.NewInt(0)
	//	base := big.NewInt(int64(f[0]))
	//	exponent := big.NewInt(int64(f[1]))
	//	currentFactor.Exp(base, exponent, nil)
	//	P.Add(P, currentFactor)
	//}
}
