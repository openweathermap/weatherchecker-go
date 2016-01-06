package adapters

import "fmt"

func ErrorOut(expectation interface{}, result interface{}) string {
	return fmt.Sprintf("\nExpected %+v\nReceived %+v\n", expectation, result)
}
