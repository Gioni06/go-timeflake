package timeflake

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gioni06/go-timeflake/internal/utils"
	"github.com/google/uuid"
	"log"
	"math/big"
	"math/rand"
	"time"
)

const (
	BASE62  = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	HEX     = "0123456789abcdef"
	MAXTS   = "281474976710655"
	MAXRAND = "1208925819614629174706175"
	MAXTF   = "340282366920938463463374607431768211455"
)

type Timeflake struct {
	Base62 string
	Hex    string
	Bytes  []byte
	Int    big.Int
	UUID   string
	rand   big.Int
}

func (f *Timeflake) Log() {
	fmt.Println("===Timeflake===")
	fmt.Printf("Timestamp: %d\n", f.Timestamp())
	fmt.Printf("Random: %s\n", f.BigRand().String())
	fmt.Printf("Base62: %s\n", f.Base62)
	fmt.Printf("Hex: %s\n", f.Hex)
	fmt.Printf("Bytes: %v\n", f.Bytes)
	fmt.Printf("Integer: %s\n", f.Int.String())
	fmt.Printf("UUID: %s\n", f.UUID)
	fmt.Println("")
	fmt.Println("")
}

func (f *Timeflake) Timestamp() int64 {
	t := new(big.Int)
	t.Rsh(&f.Int, 80)
	return t.Int64() / 1000
}

func (f *Timeflake) Rand() string {
	return f.rand.String()
}

func (f *Timeflake) BigRand() *big.Int {
	return &f.rand
}

func MaxRandom() *big.Int {
	bufMaxRandom := new(bytes.Buffer)
	var MaxRandom = big.NewInt(0)
	MaxRandom.SetString(MAXRAND, 62)
	errMaxRandom := binary.Write(bufMaxRandom, binary.LittleEndian, MaxRandom.Bytes())

	if errMaxRandom != nil {
		fmt.Println("binary.Write failed:", errMaxRandom)
	}
	return MaxRandom
}

func MaxTimestamp() *big.Int {
	bufMaxTimestamp := new(bytes.Buffer)
	var MaxTimestamp = big.NewInt(0)
	MaxTimestamp.SetString(MAXTS, 10)
	errMaxTimestamp := binary.Write(bufMaxTimestamp, binary.LittleEndian, MaxTimestamp.Bytes())

	if errMaxTimestamp != nil {
		fmt.Println("binary.Write failed:", errMaxTimestamp)
	}
	return MaxTimestamp
}

func MaxTimeflake() *big.Int {
	bufMaxTimeflake := new(bytes.Buffer)
	var MaxTimeflake = big.NewInt(0)
	MaxTimeflake.SetString(MAXTF, 10)
	errMaxTimeflake := binary.Write(bufMaxTimeflake, binary.LittleEndian, MaxTimeflake.Bytes())

	if errMaxTimeflake != nil {
		fmt.Println("binary.Write failed:", errMaxTimeflake)
	}
	return MaxTimeflake
}

func Random() Timeflake {
	now := time.Now()
	timestamp := now.Unix()

	bigTimestamp := big.NewInt(timestamp * 1000)

	//Generate cryptographically strong pseudo-random between 0 - max
	p := make([]byte, 10) // 80bits
	rand.Read(p)

	randomPart := new(big.Int)
	randomPart.SetBytes(p)

	timestampPart := new(big.Int)
	timestampPart.Lsh(bigTimestamp, 80)

	//Mix with random number
	randomAndTimestampCombined := timestampPart.Or(timestampPart, randomPart)

	v62 := big.NewInt(0)
	v62.SetBytes(randomAndTimestampCombined.Bytes())

	vHex := big.NewInt(0)
	vHex.SetBytes(randomAndTimestampCombined.Bytes())

	b62, b62Err := utils.Itoa(v62, BASE62, 22)
	hex, hexErr := utils.Itoa(vHex, HEX, 32)

	if b62Err != nil {
		log.Panic(b62Err)
	}
	if hexErr != nil {
		log.Panic(hexErr)
	}

	u, err := uuid.FromBytes(randomAndTimestampCombined.Bytes())

	if err != nil {
		log.Panic(err)
	}

	f := Timeflake{
		Base62: b62,
		Hex:    hex,
		Bytes:  randomAndTimestampCombined.Bytes(),
		UUID:   u.String(),
		Int:    *randomAndTimestampCombined,
		rand:   *randomPart,
	}

	validationError := validate(&f)
	if validationError != nil {
		log.Panic(validationError)
	}

	return f
}

func FromBytes(fromBytes []byte) Timeflake {

	bigTimestamp := new(big.Int)
	bigTimestamp.SetBytes(fromBytes[0:6])

	randomBytes := fromBytes[6:16]

	randomPart := new(big.Int)
	randomPart.SetBytes(randomBytes)

	timestampPart := big.NewInt(0)
	timestampPart.SetBytes(bigTimestamp.Bytes())
	timestampPart.Lsh(bigTimestamp, 80)

	//Mix with random number
	randomAndTimestampCombined := timestampPart.Or(timestampPart, randomPart)

	v62 := big.NewInt(0)
	v62.SetBytes(randomAndTimestampCombined.Bytes())

	vHex := big.NewInt(0)
	vHex.SetBytes(randomAndTimestampCombined.Bytes())

	b62, b62Err := utils.Itoa(v62, BASE62, 22)
	hex, hexErr := utils.Itoa(vHex, HEX, 32)

	if b62Err != nil {
		log.Panic(b62Err)
	}
	if hexErr != nil {
		log.Panic(hexErr)
	}

	u, err := uuid.FromBytes(randomAndTimestampCombined.Bytes())

	if err != nil {
		log.Panic(err)
	}

	f := Timeflake{
		Base62: b62,
		Hex:    hex,
		Bytes:  fromBytes,
		UUID:   u.String(),
		Int:    *randomAndTimestampCombined,
		rand:   *randomPart,
	}

	validationError := validate(&f)
	if validationError != nil {
		log.Panic(validationError)
	}
	return f
}

func FromHex(hexValue string) Timeflake {
	fromBytes := utils.Atoi(hexValue, HEX)
	return FromBytes(fromBytes.Bytes())
}

func FromBase62(b62 string) Timeflake {
	fromBytes := utils.Atoi(b62, BASE62)
	return FromBytes(fromBytes.Bytes())
}

type Values interface {
	Timestamp() int64
	Random() *big.Int
}

// Struct is not exported
type valuesParam struct {
	ts int64
	r  *big.Int
}

func (v *valuesParam) Timestamp() int64 {
	return v.ts
}

func (v *valuesParam) Random() *big.Int {
	return v.r
}

func NewValues(timestamp int64, random *big.Int) Values {
	if random == nil {
		//Generate cryptographically strong pseudo-random between 0 - max
		p := make([]byte, 10)
		rand.Read(p)

		random = new(big.Int)
		random.SetBytes(p)
	}
	return &valuesParam{timestamp, random} // enforce the default value here
}

func FromValues(v Values) Timeflake {

	timestamp := v.Timestamp()

	bigTimestamp := big.NewInt(timestamp * 1000)

	timestampPart := new(big.Int)
	timestampPart.Lsh(bigTimestamp, 80)

	//Mix with random number
	randomAndTimestampCombined := timestampPart.Or(timestampPart, v.Random())
	return FromBytes(randomAndTimestampCombined.Bytes())
}

func validate(flake *Timeflake) error {

	i := big.NewInt(0)
	i.SetBytes(flake.Int.Bytes())

	i2 := big.NewInt(0)
	i2.SetBytes(flake.Int.Bytes())

	zero := big.NewInt(0)
	max := MaxTimeflake()
	iLessZero := i.Cmp(zero) == -1
	iSmallerMaxTimeflake := max.Cmp(i2) == -1

	if iLessZero || iSmallerMaxTimeflake {
		return errors.New("invalid flake provided")
	}
	return nil
}
