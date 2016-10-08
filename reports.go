package irest

import (
	"fmt"
)

// Report provides test to output results for as well as various fields for
// display purposes.
type Report struct {
	InfoLabel     string
	PassTestLabel string
	FailTestLabel string
	TimingHeader  string

	Test *Test
}

// NewColoredCommandLineReport sets outputs to use ANSI color codes
// and unicode check and x.
func NewColoredCommandLineReport(t *Test) *Report {
	return &Report{
		InfoLabel:     "[ \033[00;96m\xE2\x8B\xAE\033[0m ]",
		PassTestLabel: "[ \033[00;32m\xE2\x9C\x93\033[0m ]",
		FailTestLabel: "[ \033[00;31m\xE2\x9C\x98\033[0m ]",
		TimingHeader:  "[   ms   ]",
		Test:          t,
	}
}

// PrintResults outputs report to stdout based on report fields.
func (r *Report) PrintResults() error {
	if r.Test == nil {
		return fmt.Errorf("Report.Test must be set")
	}

	fmt.Printf("%s %s %s\n", r.InfoLabel, r.TimingHeader, r.Test.Name)

	testStack := []*Test{}
	testStack = append(testStack, r.Test.Tests...)
	for len(testStack) > 0 {
		next := testStack[len(testStack)-1]
		r.printResult(next)
		if next.Tests != nil || len(next.Tests) > 0 {
			// Remove last test in stack.
			testStack = append(testStack[:0], testStack[:len(testStack)-1]...)

			// Add all sub tests of next test.
			testStack = append(testStack, next.Tests...)
		}
	}

	return nil
}

func (r *Report) printResult(t *Test) {
	var timing string

	// TODO: create thresholds that can be specified.
	if t.Duration == 0 && t.Response == nil {
		timing = "[      ]"
	} else if t.Duration < 100 {
		timing = fmt.Sprintf("[ \033[00;32m%3d ms\033[0m ]", t.Duration)
	} else if t.Duration < 500 {
		timing = fmt.Sprintf("[ \033[00;33m%3d ms\033[0m ]", t.Duration)
	} else {
		timing = fmt.Sprintf("[ \033[00;31m%3d ms\033[0m ]", t.Duration)
	}

	// Indents test by a separator to show groupings of tests.
	indent := ""
	for i := 0; i < t.Depth; i++ {
		indent += "-"
	}

	msg := t.Name
	var result string
	if t.Error == nil {
		result = r.PassTestLabel
	} else {
		result = r.FailTestLabel
		msg += fmt.Sprintf(" (%s) for %s", t.Error, t.Endpoint)
	}

	fmt.Printf("%s %s [%s] [%s] [%d] %s %s\n", result, timing, t.Method, t.Endpoint, t.Status, indent, msg)
}
