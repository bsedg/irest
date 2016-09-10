package irest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var api *httptest.Server

type SampleObject struct {
	Name    string
	Value   int
	Success bool
}

func init() {
	api = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				data, _ := json.Marshal(SampleObject{
					Name:    "unit-test",
					Value:   100,
					Success: true,
				})

				cookie := &http.Cookie{
					Name:  "test-cookie",
					Value: "unit-test-sample-value",
				}

				http.SetCookie(w, cookie)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write(data)
			} else if r.Method == "GET" {
				data, _ := json.Marshal(SampleObject{
					Name:    "unit-test",
					Value:   100,
					Success: true,
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(data)
			}
		}))
}

func TestSetup(t *testing.T) {
	test := NewTest("unit-test")

	if test.Name != "unit-test" {
		t.Errorf("name was %s, expected unit-test", test.Name)
	}
}

func TestPost(t *testing.T) {
	test := NewTest("unit-test")

	sample := SampleObject{}
	test = test.Post(api.URL, "/tests", nil).ParseResponseBody(&sample)
	if sample.Name != "unit-test" {
		t.Errorf("name response was %s, expected unit-test", sample.Name)
	}
	if !sample.Success {
		t.Error("expected response success to be true")
	}
}

func TestPostMustStatus(t *testing.T) {
	test := NewTest("unit-test")

	sample := SampleObject{}
	test = test.Post(api.URL, "/tests", nil).
		MustStatus(http.StatusCreated).
		ParseResponseBody(&sample)

	if test.Error != nil {
		t.Errorf("expected status to be 201 created: %s", test.Error.Error())
	}
}

func TestPostSaveCookie(t *testing.T) {
	test := NewTest("unit-test")

	sample := SampleObject{}
	cookie := &http.Cookie{}
	test = test.Post(api.URL, "/tests", nil).
		SaveCookie("test-cookie", cookie).
		ParseResponseBody(&sample)

	if test.Error != nil {
		t.Error(test.Error)
	}

	if cookie.Value != "unit-test-sample-value" {
		t.Errorf("expected cookie value to be 'unit-test-sample-value', got '%s'", cookie)
	}
}

func TestGet(t *testing.T) {
	test := NewTest("unit-test")
	sample := SampleObject{}
	test = test.Get(api.URL, "/tests").
		MustStatus(http.StatusOK).
		ParseResponseBody(&sample).
		MustStringValue(sample.Name, "unit-test")

	if test.Error != nil {
		t.Error(test.Error)
	}
}

func mustNil() error {
	return fmt.Errorf("must function error")
}

func TestMustFunction(t *testing.T) {
	test := NewTest("unit-test")

	sample := SampleObject{}
	test = test.Post(api.URL, "/tests", nil).
		Must(mustNil).
		ParseResponseBody(&sample)

	if test.Error == nil {
		t.Error("expecting error to be set with Must()")
	}
}

func testComplexSingleTest(t *testing.T) {
	test := NewTest("unit-test")

	sample := SampleObject{}
	cookie := &http.Cookie{}
	test = test.Post(api.URL, "/tests", nil).
		SaveCookie("test-cookie", cookie).
		MustStatus(201).
		ParseResponseBody(&sample).
		MustStringValue(sample.Name, "unit-test").
		MustIntValue(sample.Value, 100).
		Must(func() error {
			if !sample.Success {
				return fmt.Errorf("expected true success")
			}
			return nil
		})

	if test.Error != nil {
		t.Error(test.Error)
	}
}

func testInnerTests(t *testing.T) {
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

	if len(test.Tests) != 3 {
		t.Errorf("expected 3 inner tests, got %d", len(test.Tests))
	}

	if test.Tests[0].Error != nil {
		t.Error(test.Tests[0].Error)
	}

	if test.Tests[1].Error != nil {
		t.Error(test.Tests[1].Error)
	}

	if test.Tests[2].Error == nil || !strings.Contains(test.Tests[2].Error.Error(), "status") {
		t.Errorf("expected must status to fail")
	}
}
