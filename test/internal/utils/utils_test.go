package tests

import (
	"github.com/gioni06/go-timeflake/internal/utils"
	"github.com/gioni06/go-timeflake/pkg/timeflake"
	"math/big"
	"reflect"
	"strings"
	"testing"
)

func TestItoa(t *testing.T) {
	b1 := big.NewInt(1504324233)
	res1, _ := utils.Itoa(b1, timeflake.BASE62, 5)

	if res1 != "01dnzS5" {
		t.Errorf("expected '01dnzS5' got '%s'", res1)
	}

	b2 := big.NewInt(1504324233)
	res2, _ := utils.Itoa(b2, timeflake.HEX, 5)

	if res2 != "059aa2a89" {
		t.Errorf("expected '059aa2a89' got '%s'", res2)
	}
}

func TestMax(t *testing.T) {
	m := utils.Max(123, 1)
	if m != 123 {
		t.Errorf("expected '123' got '%d'", m)
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

func TestBigFromStringBase10(t *testing.T) {
	b, err := utils.BigFromString("123", 10)

	if err != nil {
		t.Error(err)
	}

	a := big.NewInt(123)

	if a.Cmp(b) != 0 {
		t.Error("failed")
	}
}

func TestBigFromStringBase62(t *testing.T) {
	b := utils.Atoi("8M0kX", timeflake.BASE62)

	a := big.NewInt(123456789)

	if a.Cmp(b) != 0 {
		t.Error("failed")
	}
}
