package tests

import (
	"github.com/gioni06/go-timeflake/internal/utils"
	"github.com/gioni06/go-timeflake/pgk/timeflake"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestTimeflake(t *testing.T) {
	now := (time.Now()).Unix()
	for range make([]int, 1000) {
		f := timeflake.Random()

		zero := new(big.Int)
		max := timeflake.MaxTimeflake()

		isTimeflakeGTEZero := zero.Cmp(&f.Int) <= 0
		isTimeflakeLTEMaxTimeflake := max.Cmp(&f.Int) >= 0

		if !(isTimeflakeGTEZero && isTimeflakeLTEMaxTimeflake) {
			t.Error("timeflake out of bounds")
		}

		if !(now <= f.Timestamp()) {
			t.Error("timeflake timestamp is greater then expected")
		}

		zero = new(big.Int)
		max = timeflake.MaxRandom()
		r, err := utils.BigFromString(f.Rand(), 10)

		if err != nil {
			t.Error("random part is not a valid number")
		}

		isTimeflakeGTEZero = zero.Cmp(r) <= 0
		isTimeflakeLTEMaxTimeflake = max.Cmp(r) >= 0

		if !(isTimeflakeGTEZero && isTimeflakeLTEMaxTimeflake) {
			t.Error("timeflake random part out of bounds")
		}

		zero = new(big.Int)
		max = timeflake.MaxTimestamp()

		flakeTS := new(big.Int)
		flakeTS.SetInt64(f.Timestamp())

		isTimeflakeGTEZero = zero.Cmp(flakeTS) <= 0
		isTimeflakeLTEMaxTimeflake = max.Cmp(flakeTS) >= 0

		if !(isTimeflakeGTEZero && isTimeflakeLTEMaxTimeflake) {
			t.Error("timeflake timestamp part out of bounds")
		}
	}
}

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
