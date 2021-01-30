package tests

import (
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/gioni06/go-timeflake/internal/alphabets"
	"github.com/gioni06/go-timeflake/internal/utils"
)

func TestBigIntToASCII(t *testing.T) {
	b1 := big.NewInt(1504324233)
	res1, _ := utils.BigIntToASCII(b1, alphabets.BASE62, 5)

	if res1 != "01dnzS5" {
		t.Errorf("expected '01dnzS5' got '%s'", res1)
	}

	b2 := big.NewInt(1504324233)
	res2, _ := utils.BigIntToASCII(b2, alphabets.HEX, 5)

	if res2 != "059aa2a89" {
		t.Errorf("expected '059aa2a89' got '%s'", res2)
	}


	b3 := big.NewInt(0)
	res3, _ := utils.BigIntToASCII(b3, alphabets.HEX, 5)

	if res3 != "0" {
		t.Errorf("expected '0' got '%s'", res3)
	}
}

func TestASCIIToBigInt(t *testing.T) {
	b := utils.ASCIIToBigInt("8M0kX", alphabets.BASE62)

	a := big.NewInt(123456789)

	if a.Cmp(b) != 0 {
		t.Error("failed")
	}

	b1 := utils.ASCIIToBigInt("0", alphabets.BASE62)

	a1 := big.NewInt(0)

	if a1.Cmp(b1) != 0 {
		t.Error("failed")
	}
}

func TestFillString(t *testing.T) {
	s := utils.FillString("x", 3)
	if s != "xxx" {
		t.Errorf("expected 'xxx' got '%s'", s)
	}
}

func TestIndexAlphabet(t *testing.T) {
	m := utils.IndexAlphabet(strings.Split("0123abc", ""))
	if len(m) != 7 {
		t.Errorf("expected '7' got '%d", len(m))
	}

	v := make(map[string]int)
	v["0"] = 0
	v["1"] = 1
	v["2"] = 2
	v["3"] = 3
	v["a"] = 4
	v["b"] = 5
	v["c"] = 6

	if !reflect.DeepEqual(m, v) {
		t.Errorf("resulting map is not correct")
	}
}
