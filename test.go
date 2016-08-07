// Package iREST is an integration testing package for RESTful
// APIs. It simply makes HTTP requests and allows for checking of
// responses, status codes, etc.

package irest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Test struct contains sub-tests that can be isolated test cases as well
// as the HTTP related objects and errors of current test and its sub-tests.
type Test struct {
	Name     string `json:"name"`
	Error    error  `json:"err"`
	Status   int    `json:"status"`
	Tests    []*Test
	Errors   []error
	Created  time.Time `json:"created"`
	Duration int64     `json:"duration"`

	// HTTP related fields for making requests and getting responses.
	Client   *http.Client
	Header   *http.Header
	Cookie   *http.Cookie
	Response *http.Response
}

// NewTest creates a new test with a given name.
func NewTest(name string) *Test {
	t := &Test{
		Name:    name,
		Error:   nil,
		Tests:   []*Test{},
		Errors:  []error{},
		Created: time.Now(),
		Client:  &http.Client{},
		Header:  &http.Header{},
	}

	return t
}

// NewTest adds a Test as a sub-test to the current one. Sub-tests act as
// individual test cases.
func (t *Test) NewTest(name string) *Test {
	testCase := &Test{
		Name:   name,
		Tests:  []*Test{},
		Client: t.Client,
		Header: &http.Header{},
	}

	t.Tests = append(t.Tests, testCase)

	return testCase
}

// AddHeader is a utility function to just wrap setting a header with a value
// by name.
func (t *Test) AddHeader(name, value string) *Test {
	t.Header.Set(name, value)
	return t
}

// Post sends a HTTP Post request with given URL from baseURL combined with
// endpoint and then saves the result.
func (t *Test) Post(baseURL, endpoint string, result interface{}) *Test {
	t.Header.Set("Content-Type", "application/json")

	addr, err := url.Parse(baseURL + endpoint)
	if err != nil {
		t.Error = err
		return t
	}

	req, err := http.NewRequest("POST", addr.String(), nil)
	if err != nil {
		t.Error = err
		return t
	}

	res, err := t.Client.Do(req)
	if err != nil {
		t.Error = err
		return t
	}

	t.Response = res
	t.Status = res.StatusCode

	body := res.Body
	defer body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Error = err
		return t
	}

	json.Unmarshal(data, result)

	return t
}

// MustStatus sets the Test.Error if the status code is not the expected
// value. An HTTP request must have been made prior to this function call.
func (t *Test) MustStatus(statusCode int) *Test {
	if t.Error != nil {
		return t
	}

	if t.Status != statusCode {
		t.Error = fmt.Errorf("expected status code response of %d, actual %d", statusCode, t.Status)
	}

	return t
}

// MustStringValue compares two string values and sets the Test.Error if not
// equal.
func (t *Test) MustStringValue(expected, actual string) *Test {
	if t.Error != nil {
		return t
	}

	if expected != actual {
		t.Error = fmt.Errorf("expected %s, but got %s", expected, actual)
	}

	return t
}

// MustIntValue compares two int values and sets the Test.Error if not equal.
func (t *Test) MustIntValue(expected, actual int) *Test {
	if t.Error != nil {
		return t
	}

	if expected != actual {
		t.Error = fmt.Errorf("expected %d, but got %d", expected, actual)
	}

	return t
}

// MustFunction adds the ability to create functions that can check something
// not covered by the current functions currently.
type MustFunction func() error

// Must allows for passing in created functions matching the MustFunction
// pattern with no parameters returning an error.
func (t *Test) Must(fn MustFunction) *Test {
	if t.Error != nil {
		return t
	}

	if err := fn(); err != nil {
		t.Error = err
	}

	return t
}

// SaveCookie will store the cookie with the specified name if it exists in the
// response. An HTTP request must have been made prior to this function call.
func (t *Test) SaveCookie(name string, cookie *http.Cookie) *Test {
	if t.Error != nil {
		return t
	}

	if t.Response == nil {
		t.Error = fmt.Errorf("http response not set, must have request before saving result")
		return t
	}

	for _, c := range t.Response.Cookies() {
		if c.Name == name {
			cookie.Name = c.Name
			cookie.Value = c.Value
			return t
		}
	}

	t.Error = fmt.Errorf("cookie name '%s' not found", name)

	return t
}
