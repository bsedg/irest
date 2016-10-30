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
			if r.Method == http.MethodPost || r.Method == http.MethodPut {
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
			} else if r.Method == http.MethodGet {
				data, _ := json.Marshal(SampleObject{
					Name:    "unit-test",
					Value:   100,
					Success: true,
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(data)
			} else if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
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

func TestPut(t *testing.T) {
	test := NewTest("unit-test")

	sample := SampleObject{}
	test = test.Put(api.URL, "/tests", nil).ParseResponseBody(&sample)
	if sample.Name != "unit-test" {
		t.Errorf("name response was %s, expected unit-test", sample.Name)
	}
	if !sample.Success {
		t.Error("expected response success to be true")
	}
}

func TestDelete(t *testing.T) {
	test := NewTest("unit-test")

	test = test.Delete(api.URL, "/tests")
	if test.Status != http.StatusNoContent {
		t.Errorf("expected status 204 No Content, instead was %d", test.Status)
	}
}

func TestMalformedUrl(t *testing.T) {
	test := NewTest("unit-test")

	test = test.Get(api.URL, "bad-url%@%(*%)///\\####")
	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestBadRequestBody(t *testing.T) {
	test := NewTest("unit-test")

	// pass a clearly bogus request body to force a json parse error
	test = test.Post(api.URL, "/tests", make(chan int))
	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestRequestFailure(t *testing.T) {
	test := NewTest("unit-test")

	test = test.Get("", "")
	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestParseResponseBodyEmptyResponse(t *testing.T) {
	test := NewTest("unit-test")

	test = test.ParseResponseBody(nil)

	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestAddDefaultHeader(t *testing.T) {
	test := NewTest("unit-test").
		AddHeader("Content-Type", "application/json")

	subTest := test.NewTest("sub-test")

	if subTest.Header.Get("Content-Type") == "" {
		t.Error("expected Content-Type to be set")
	}
}

func TestAddCookie(t *testing.T) {
	cookie := &http.Cookie{Name: "test-cookie", Value: "test-value"}
	test := NewTest("unit-test").
		AddCookie(cookie)

	test = test.Get(api.URL, "/tests")
	if test.Cookies[0] != cookie {
		t.Error("expected cookie to be set")
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

func TestMustStatusError(t *testing.T) {
	test := NewTest("unit-test")

	test.Error = fmt.Errorf("testing error")
	test = test.MustStatus(500)

	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestMustStatusMismatch(t *testing.T) {
	test := NewTest("unit-test")

	test.Status = 418
	test = test.MustStatus(451)

	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestMustStringValueError(t *testing.T) {
	test := NewTest("unit-test")

	test.Error = fmt.Errorf("testing error")
	test = test.MustStringValue("test", "test")

	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestMustStringValueMismatch(t *testing.T) {
	test := NewTest("unit-test")

	test = test.MustStringValue("foo", "bar")

	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestMustIntValueError(t *testing.T) {
	test := NewTest("unit-test")

	test.Error = fmt.Errorf("testing error")
	test = test.MustIntValue(42, 42)

	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestMustIntValueMismatch(t *testing.T) {
	test := NewTest("unit-test")

	test = test.MustIntValue(42, 43)

	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestMustError(t *testing.T) {
	test := NewTest("unit-test")

	test.Error = fmt.Errorf("testing error")
	test = test.Must(mustNil)

	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestSaveCookieError(t *testing.T) {
	test := NewTest("unit-test")

	test.Error = fmt.Errorf("testing error")
	test = test.SaveCookie("test-cookie", &http.Cookie{})

	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestSaveCookieNoResponse(t *testing.T) {
	test := NewTest("unit-test")

	test = test.SaveCookie("test-cookie", &http.Cookie{})

	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}

	msg := "http response not set, must have request before saving result"
	if test.Error.Error() != msg {
		t.Errorf("expected '%s', got '%s' instead", msg, test.Error.Error())
	}
}

func TestSaveCookieDoesNotExist(t *testing.T) {
	test := NewTest("unit-test")

	cookie := &http.Cookie{}
	test = test.Post(api.URL, "/tests", nil).
		SaveCookie("nonexistent-cookie", cookie)

	if test.Error == nil {
		t.Error("expected an error, but did not get one")
	}

	msg := "cookie name 'nonexistent-cookie' not found"
	if test.Error.Error() != msg {
		t.Errorf("expected '%s', but got '%s' instead", msg, test.Error.Error())
	}
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
