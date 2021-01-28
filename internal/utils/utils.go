// The utils package provides a couple of helpers that make creating
// and testing Timeflakes easier
package utils

import (
	"math"
	"math/big"
	"strings"
)

// The 'strconv' package provides a Itoa function, but it can only deal with Int values.
// This converts a big.Int number to a string for a given alphabet.
func BigIntToASCII(value *big.Int, alphabet string, padding int) (string, error) {
	alphabetSlice := strings.Split(alphabet, "")
	v := new(big.Int)
	v.SetBytes(value.Bytes())

	z := big.NewInt(0)
	if v.Cmp(z) == 0 {
		return alphabetSlice[0], nil
	}

	var result string

	for v.Cmp(z) != 0 {
		rem := big.NewInt(0)
		base := big.NewInt(int64(len(alphabet)))
		v.DivMod(v, base, rem) // value.QuoRem64(uint64(base))
		result = alphabetSlice[rem.Int64()] + result
	}

	if padding != 0 {
		fill := math.Max(float64(padding-len(result)), 0)
		result = FillString(alphabetSlice[0], int(fill)) + result
	}

	return result, nil
}

// The 'strconv' package provides a Atoi function, but it can only deal with Int values.
// This converts a string to a big.Int value for a given alphabet.
func ASCIIToBigInt(value string, alphabet string) *big.Int {
	alphabetSlice := strings.Split(alphabet, "")
	if value == alphabetSlice[0] {
		return big.NewInt(0)
	}
	index := IndexAlphabet(alphabetSlice)
	result := big.NewInt(0)
	base := new(big.Int)
	base.SetInt64(int64(len(alphabet)))
	for _, char := range strings.Split(value, "") {
		m := new(big.Int)
		m.Mul(result, base)

		i := new(big.Int)
		i.SetInt64(int64(index[char]))
		result.Add(m, i)
	}
	return result
}

// Fills a string with a given character.
func FillString(char string, length int) string {
	o := char
	n := o
	for i := 1; i < length; i++ {
		n += o
	}
	return n
}

// Creates a map from a slice of strings.
func IndexAlphabet(alphabet []string) map[string]int {
	res := make(map[string]int)
	for k, v := range alphabet {
		res[v] = k
	}
	return res
}
