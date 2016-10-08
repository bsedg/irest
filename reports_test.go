package irest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestReportOutputHasMethodStatusAndEndpoint(t *testing.T) {
	test := NewTest("unit-test")

	postSample := SampleObject{}
	getSample := SampleObject{}
	cookie := &http.Cookie{}
	createTest := test.NewTest("create")
	createTest = createTest.Post(api.URL, "/tests", nil).
		SaveCookie("test-cookie", cookie).
		MustStatus(201).
		ParseResponseBody(&postSample).
		MustStringValue(postSample.Name, "unit-test").
		MustIntValue(postSample.Value, 100).
		Must(func() error {
			if !postSample.Success {
				return fmt.Errorf("expected true success")
			}
			return nil
		})

	getTest := test.NewTest("get")
	getTest = getTest.Get(api.URL, "/tests").
		MustStatus(200).
		ParseResponseBody(&getSample).
		MustStringValue(getSample.Name, "unit-test").
		MustIntValue(getSample.Value, 100)

	failedGetTest := test.NewTest("failed get")
	failedGetTest = failedGetTest.Get(api.URL, "/tests").
		MustStatus(404).
		ParseResponseBody(&getSample).
		MustStringValue(getSample.Name, "unit-test").
		MustIntValue(getSample.Value, 100)

	report := NewColoredCommandLineReport(test)

	expectedAdditionalOutput := []string{
		`[GET] [/tests] [200]`,
		`[POST] [/tests] [201]`,
	}

	output := captureStdout(report.PrintResults)

	for _, expect := range expectedAdditionalOutput {
		if !strings.Contains(output, expect) {
			t.Error("incorrect report output!", output, "should contain:", expect)
		}
	}
}

func captureStdout(f func() error) string {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()

	output := make(chan string)
	go func() {
		var b bytes.Buffer
		io.Copy(&b, r)
		output <- b.String()
	}()

	w.Close()
	os.Stdout = stdout
	captured := <-output

	return captured
}
