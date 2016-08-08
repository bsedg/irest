package irest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var api *httptest.Server

type SampleObject struct {
	Name    string
	Value   int
	Success bool
}

func init() {
	data, _ := json.Marshal(SampleObject{
		Name:    "unit-test",
		Value:   100,
		Success: true,
	})
	api = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				cookie := &http.Cookie{
					Name:  "test-cookie",
					Value: "unit-test-sample-value",
				}
				http.SetCookie(w, cookie)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write(data)
			} else if r.Method == "GET" {
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
	test = test.Post(api.URL, "/tests", nil, &sample)
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
	test = test.Post(api.URL, "/tests", nil, &sample).MustStatus(http.StatusCreated)
	if test.Error != nil {
		t.Errorf("expected status to be 201 created, not %d: %s", test.Error.Error())
	}
}

func TestPostSaveCookie(t *testing.T) {
	test := NewTest("unit-test")

	sample := SampleObject{}
	cookie := &http.Cookie{}
	test = test.Post(api.URL, "/tests", nil, &sample).SaveCookie("test-cookie", cookie)

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
	test = test.Get(api.URL, "/tests", &sample).MustStatus(http.StatusOK).MustStringValue(sample.Name, "unit-test")

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
	test = test.Post(api.URL, "/tests", nil, &sample).Must(mustNil)

	if test.Error == nil {
		t.Error("expecting error to be set with Must()")
	}
}
