package utils

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"strings"
)

func Itoa(value *big.Int, alphabet string, padding int) (string, error) {
	alphabetSlice := strings.Split(alphabet, "")
	z := big.NewInt(0)
	if value.Cmp(z) == 0 {
		return alphabetSlice[0], nil
	}
	result := ""

	for value.Cmp(z) != 0 {
		rem := big.NewInt(0)
		base := big.NewInt(int64(len(alphabet)))
		value.DivMod(value, base, rem) // value.QuoRem64(uint64(base))
		result = alphabetSlice[rem.Int64()] + result
	}

	if padding != 0 {
		fill := Max(padding-len(result), 0)
		result = FillString(alphabetSlice[0], fill) + result
	}

	return result, nil
}

func Atoi(value string, alphabet string) *big.Int {
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

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func FillString(char string, amount int) string {
	o := char
	n := o
	for i := 1; i < amount; i++ {
		n += o
	}
	return n
}

func IndexAlphabet(alphabet []string) map[string]int {
	res := make(map[string]int)
	for k, v := range alphabet {
		res[v] = k
	}
	return res
}

func BigFromString(s string, base int) (*big.Int, error) {
	buf := new(bytes.Buffer)
	var b = big.NewInt(0)

	b.SetString(s, base)
	err := binary.Write(buf, binary.LittleEndian, b.Bytes())

	if err != nil {
		return nil, err
	}

	return b, nil
}