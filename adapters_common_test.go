package main

import "fmt"

func ErrorOut(expectation, result interface{}) string {
	return fmt.Sprintf("\nExpected %#v\nReceived %#v\n", expectation, result)
}

func SprintfCompare(expectation, result interface{}) bool {
	return fmt.Sprintf("%#v", expectation) == fmt.Sprintf("%#v", result)
}
