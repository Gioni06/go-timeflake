package tests

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"math/rand"
	"reflect"
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

func TestTimeflakeCreationFromHex(t *testing.T) {
	_, err := timeflake.FromHex("0177487ec2f8d0a63f2785a9cadfc50f")

	if err != nil {
		t.Error("timeflake creation should not fail")
	}
}

func TestTimeflakeCreationFromB62(t *testing.T) {
	_, err := timeflake.FromBase62("02lVIoVLUfN6xUwLlnSRjj")

	if err != nil {
		t.Error("timeflake creation should not fail")
	}
}

func TestTimeflakeValuesInstanceCanBeCreated(t *testing.T) {
	now := int64(1611829003)
	v := timeflake.NewValues(now, nil)
	tf, err := timeflake.FromValues(v)

	if err != nil {
		t.Error("timeflake creation should not fail")
	}

	if tf != nil && tf.Timestamp() != now {
		t.Error("timeflake should have a correct timestamp")
	}

	if tf != nil && tf.Rand() == "" {
		t.Error("timeflake should have a random part")
	}

	if tf != nil && len(tf.BigRand().Bytes()) != 10 {
		t.Errorf("timeflake random part must be 10 Bytes %d", len(tf.BigRand().Bytes()))
	}
}

func TestTimeflakeCreationFromValuesFails(t *testing.T) {
	now := int64(0)
	randString := "0"
	random, _ := bigFromString(randString, 10)

	v := timeflake.NewValues(now, random)
	_, err := timeflake.FromValues(v)

	if err == nil {
		t.Error("timeflake creation should fail")
	}
}

// Creates a Timeflake from fixed values and makes assertions.
func TestTimeflakeInstance(t *testing.T) {
	now := int64(1611829003)
	randString := "985318938706034770822415"
	random, _ := bigFromString(randString, 10)

	v := timeflake.NewValues(now, random)
	tf, err := timeflake.FromValues(v)

	if err != nil {
		t.Error(err)
	}

	expectedInt, _ := bigFromString("1948581698531390905820074514793350415", 10)
	if tf.Int.Cmp(expectedInt) != 0 {
		t.Errorf("timeflake Int is not correct. %s != %s", tf.Int.String(), expectedInt.String())
	}

	expectedBigRandom := new(big.Int)
	expectedBigRandom.SetBytes(random.Bytes())
	if tf.BigRand().Cmp(expectedBigRandom) != 0 {
		t.Errorf("timeflake BigRandom is not correct. %s != %s", tf.BigRand().String(), expectedBigRandom.String())
	}

	if tf.Rand() != randString {
		t.Errorf("timeflake random part is not correct. %s != %s", tf.Rand(), randString)
	}

	if tf.Timestamp() != now {
		t.Errorf("timeflake timestamp is not correct. %d != %d", tf.Timestamp(), now)
	}

	expectedB62 := "02lVIoVLUfN6xUwLlnSRjj"
	if tf.Base62 != expectedB62 {
		t.Errorf("timeflake B62 is not correct. %s != %s", tf.Base62, expectedB62)
	}

	expectedHEX := "0177487ec2f8d0a63f2785a9cadfc50f"
	if tf.Hex != expectedHEX {
		t.Errorf("timeflake HEX is not correct. %s != %s", tf.Hex, expectedHEX)
	}

	expectedBytes := make([]byte, 16)
	expectedBytes = append(expectedBytes, 1, 119, 72, 126, 194, 248, 208, 166, 63, 39, 133, 169, 202, 223, 197, 15)
	if reflect.DeepEqual(tf.Bytes, expectedBytes) {
		t.Errorf("timeflake Bytes is not correct. %s != %s", tf.Bytes, expectedBytes)
	}

	expectedUUID := "0177487e-c2f8-d0a6-3f27-85a9cadfc50f"
	if tf.UUID != expectedUUID {
		t.Errorf("timeflake Int is not correct. %s != %s", tf.UUID, expectedUUID)
	}
}

// Run some sanity checks that ensure the overall correctness
// of the created Timeflakes.
func TestTimeflake(t *testing.T) {
	now := (time.Now()).Unix()
	for range make([]int, 1000) {
		f, _ := timeflake.Random()

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
	results := make(chan *timeflake.Timeflake, numJobs)

	var seen = make(map[string]int)

	for i := 0; i < 100; i++ {
		go func(j <-chan int, r chan<- *timeflake.Timeflake) {
			for range j {
				f, err := timeflake.Random()
				if err != nil {
					close(jobs)
					close(results)
					t.Error("Timeflake creation failed")
					break
				}
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
