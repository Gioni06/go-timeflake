package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jaffee/commandeer"

	"github.com/gioni06/go-timeflake/internal/customerr"
)

var (
	Red    = Color("\033[1;31m%s\033[0m")
	Yellow = Color("\033[1;33m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func main() {
	err := commandeer.Run(NewMain())
	if err != nil {
		switch err.(type) {
		case *customerr.OutOfBoundsError:
			fmt.Printf(Yellow("%s, try again using a smaller random part\n"), err.Error())
		case *customerr.ConversionError:
			fmt.Printf(Yellow("%s, converting the inputs to a timeflake failed\n"), err.Error())
		case *customerr.UUIDError:
			fmt.Printf(Yellow("%s, the timeflake can not be converted to a valid uuid\n"), err.Error())
		default:
			fmt.Println(Red(err.Error()))
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
