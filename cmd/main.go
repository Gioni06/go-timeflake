package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/gioni06/go-timeflake/pkg/timeflake"
)

func main() {
	f := timeflake.Random()

	c := timeflake.FromBytes(f.Bytes)

	d := timeflake.FromHex(f.Hex)

	e := timeflake.FromBase62(f.Base62)

	g := timeflake.NewValues(f.Timestamp(), f.BigRand())
	g1 := timeflake.NewValues(f.Timestamp(), nil)

	h := timeflake.FromValues(g)
	i := timeflake.FromValues(g1)

	c.Log()
	f.Log()
	d.Log()
	e.Log()
	h.Log()
	i.Log()

	fmt.Println("===Timeflake Comparison===")
	fmt.Printf("c == f => %v\n", reflect.DeepEqual(c, f))
	fmt.Printf("f == d => %v\n", reflect.DeepEqual(f, d))
	fmt.Printf("d == e => %v\n", reflect.DeepEqual(d, e))
	fmt.Printf("e == h => %v\n", reflect.DeepEqual(e, h))
	fmt.Printf("h == i => %v (expect false)\n", reflect.DeepEqual(h, i))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
