// Package irest is an integration testing package for RESTful
// APIs. It simply makes HTTP requests and allows for checking of
// responses, status codes, etc.
package irest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Test struct contains sub-tests that can be isolated test cases as well
// as the HTTP related objects and errors of current test and its sub-tests.
type Test struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Error    error  `json:"err"`
	Status   int    `json:"status"`

	Tests    []*Test
	Errors   []error
	Created  time.Time `json:"created"`
	Duration int64     `json:"duration"`
	Depth    int

	// HTTP related fields for making requests and getting responses.
	Client   *http.Client
	Header   *http.Header
	Cookies  []*http.Cookie
	Response *http.Response
}

// NewTest creates a new test with a given name.
func NewTest(name string) *Test {
	t := &Test{
		Name:    name,
		Error:   nil,
		Depth:   0,
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
		Depth:  t.Depth + 1,
		Tests:  []*Test{},
		Client: t.Client,
		Header: &http.Header{},
	}

	t.Tests = append(t.Tests, testCase)

	// For convenience, bring down header values that were set on the
	// parent test.
	testCase.Header = t.Header

	return testCase
}

// AddHeader is a utility function to just wrap setting a header with a value
// by name.
func (t *Test) AddHeader(name, value string) *Test {
	t.Header.Set(name, value)
	return t
}

// AddCookie adds to the slice of cookies to be included in the request.
func (t *Test) AddCookie(c *http.Cookie) *Test {
	t.Cookies = append(t.Cookies, c)
	return t
}

// Get retrieves data from a specified endpoint.
func (t *Test) Get(baseURL, endpoint string) *Test {
	return t.do("GET", baseURL, endpoint, nil)
}

// Post sends a HTTP POST request with given URL from baseURL combined with
// endpoint and sends the data as request body.
func (t *Test) Post(baseURL, endpoint string, data interface{}) *Test {
	return t.do("POST", baseURL, endpoint, data)
}

// Put sends a HTTP PUT request with given URL from baseURL combined with
// endpoint and sends the data as request body.
func (t *Test) Put(baseURL, endpoint string, data interface{}) *Test {
	return t.do("PUT", baseURL, endpoint, data)
}

// Delete deletes data from a specified endpoint.
func (t *Test) Delete(baseURL, endpoint string) *Test {
	return t.do("DELETE", baseURL, endpoint, nil)
}

func (t *Test) do(method, baseURL, endpoint string, data interface{}) *Test {
	t.Endpoint = endpoint

	b := new(bytes.Buffer)
	if data != nil {
		if err := json.NewEncoder(b).Encode(data); err != nil {
			t.Error = err
			return t
		}
	}

	req, err := http.NewRequest(method, baseURL+endpoint, b)
	if err != nil {
		t.Error = err
		return t
	}

	req.Header = *t.Header

	for _, c := range t.Cookies {
		req.AddCookie(c)
	}

	startTime := time.Now()
	res, err := t.Client.Do(req)
	if err != nil {
		t.Error = err
		return t
	}
	reqDuration := time.Since(startTime)
	t.Duration = reqDuration.Nanoseconds() / int64(time.Millisecond)

	t.Response = res
	t.Status = res.StatusCode

	return t
}

// ParseResponseBody parses the HTTP response body from json
// to a provided interface.
func (t *Test) ParseResponseBody(result interface{}) *Test {
	if t.Response == nil || t.Response.Body == nil {
		t.Error = fmt.Errorf("need response body to parse")
		return t
	}

	defer t.Response.Body.Close()

	resultBody, err := ioutil.ReadAll(t.Response.Body)

	if err != nil {
		t.Error = err
		return t
	}

	if err := json.Unmarshal(resultBody, result); err != nil {
		t.Error = err
		return t
	}

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
