package main

import (
	"errors"
	"fmt"

	"github.com/gioni06/go-timeflake/internal/utils"
	"github.com/gioni06/go-timeflake/pkg/timeflake"
)

type Main struct {
	About      bool   `help:"About Go-Timeflake"`
	Random     bool   `help:"Generate a random Timeflake"`
	Number     int    `flag:"n" help:"How many random Timeflakes should be generated?"`
	Values     bool   `help:"Create a Timeflake from timestamp and random(optional)"`
	RandomPart string `flag:"r" help:"A large random number e.x. '985318938706034770822415'"`
	Timestamp  int64  `flag:"t" help:"A Unix timestamp"`
}

func NewMain() *Main { return &Main{Values: false, RandomPart: "", Random: false, Number: 1} }

func (m *Main) Run() error {
	if m.About {
		fmt.Println(`Go-Timeflake is a 128-bit, roughly-ordered, URL-safe UUID.`)
		fmt.Println(`A Golang port of https://github.com/anthonynsimon/timeflake`)
		fmt.Println(`Visit https://github.com/Gioni06/go-timeflake for more information`)
		return nil
	}

	if m.Random {
		for i := 0; i < m.Number; i++ {
			tf, _ := timeflake.Random()
			tf.Log()
		}
		return nil
	}

	if m.Values {
		var tf *timeflake.Timeflake
		var err error
		if m.RandomPart != "" {
			r := utils.ASCIIToBigInt(m.RandomPart, "0123456789")
			if r == nil {
				return errors.New("can not parse random number")
			}
			v := timeflake.NewValues(m.Timestamp, r)
			tf, err = timeflake.FromValues(v)
			if err != nil {
				return err
			} else {
				tf.Log()
				return nil
			}
		} else {
			v := timeflake.NewValues(m.Timestamp, nil)
			tf, err = timeflake.FromValues(v)
			if err != nil {
				return err
			} else {
				tf.Log()
				return nil
			}
		}
	}
	return nil
}
