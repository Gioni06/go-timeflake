package tests

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/gioni06/go-timeflake/pkg/timeflake"
)

func bigFromString(s string, base int) (*big.Int, error) {
	buf := new(bytes.Buffer)
	var b = big.NewInt(0)

	b.SetString(s, base)
	err := binary.Write(buf, binary.LittleEndian, b.Bytes())

	if err != nil {
		return nil, err
	}

	return b, nil
}

// Run some sanity checks that ensure the overall correctness
// of the created Timeflakes.
func TestTimeflake(t *testing.T) {
	now := (time.Now()).Unix()
	for range make([]int, 1000) {
		f := timeflake.Random()

		zero := new(big.Int)
		max := timeflake.MaxTimeflake()

		isTimeflakeGTEZero := zero.Cmp(&f.Int) <= 0
		isTimeflakeLTEMaxTimeflake := max.Cmp(&f.Int) >= 0

		// 0 < Timeflake integer value (big.Int) < MaxTimeflake
		if !(isTimeflakeGTEZero && isTimeflakeLTEMaxTimeflake) {
			t.Error("timeflake out of bounds")
		}

		// Proof that Timeflakes are not capable of time travel into the future. Great!
		if !(now <= f.Timestamp()) {
			t.Error("timeflake timestamp is greater then expected")
		}

		zero = new(big.Int)
		max = timeflake.MaxRandom()
		r, err := bigFromString(f.Rand(), 10)

		if err != nil {
			t.Error("random part is not a valid number")
		}

		isTimeflakeGTEZero = zero.Cmp(r) <= 0
		isTimeflakeLTEMaxTimeflake = max.Cmp(r) >= 0

		// 0 < Timeflake random bits < MaxRandom
		if !(isTimeflakeGTEZero && isTimeflakeLTEMaxTimeflake) {
			t.Error("timeflake random part out of bounds")
		}

		zero = new(big.Int)
		max = timeflake.MaxTimestamp()

		flakeTS := new(big.Int)
		flakeTS.SetInt64(f.Timestamp())

		isTimeflakeGTEZero = zero.Cmp(flakeTS) <= 0
		isTimeflakeLTEMaxTimeflake = max.Cmp(flakeTS) >= 0

		// We know that a Timeflake can not travel into the future.
		// This makes sure, that it can not travel past the boundaries of the Unix
		// timestamp.
		// @Todo - Check this line before January 19, 2038
		if !(isTimeflakeGTEZero && isTimeflakeLTEMaxTimeflake) {
			t.Error("timeflake timestamp part out of bounds")
		}
	}
}

// Test the randomness of created Timeflakes.
// This is achieved by creating 1.000,000 concurrent jobs that create
// Timeflakes. The base62 values are used as keys in a hashmap.
// If a key is already present, this test will fail.
//
// Test takes about 2.5sec on a recent MacBook. Don't hesitate to
// run a billion jobs, but don't forget to get yourself a cup of coffee
// before you start!
func TestRandomnessOfTimeflakes(t *testing.T) {
	const numJobs = 1e6
	jobs := make(chan int, numJobs)
	results := make(chan timeflake.Timeflake, numJobs)

	var seen = make(map[string]int)

	for i := 0; i < 100; i++ {
		go func(j <-chan int, r chan<- timeflake.Timeflake) {
			for range j {
				f := timeflake.Random()
				r <- f
			}
		}(jobs, results)
	}

	for i := 0; i < numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	for j := 0; j < numJobs; j++ {
		r := <-results
		_, ok := seen[r.Base62]

		if ok {
			t.Errorf("flake collision found after %d generations", j)
		} else {
			seen[r.Base62] = j
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
