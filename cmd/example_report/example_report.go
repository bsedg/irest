package main

import (
	"fmt"

	"github.com/bsedg/irest"
)

func main() {
	t := irest.NewTest("example")

	// Example tests with artificial data.
	t1 := t.NewTest("example - 1")
	t1.Duration = 50
	t1.Endpoint = "/examples/1"
	t1.Error = nil

	t2 := t.NewTest("example - 2")
	t2.Duration = 400
	t2.Endpoint = "/examples/2"
	t2.Error = fmt.Errorf("expected different result")

	r := irest.NewColoredCommandLineReport(t)
	r.PrintResults()
}
