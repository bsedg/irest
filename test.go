package irest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Test struct {
	Name   string `json:"name"`
	Error  error  `json:"err"`
	Status int    `json:"status"`

	Tests    []*Test
	Errors   []error
	Created  time.Time `json:"created"`
	Duration int64     `json:"duration"`

	Client   *http.Client
	Header   *http.Header
	Cookie   *http.Cookie
	Response *http.Response
}

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

func (t *Test) AddHeader(name, value string) *Test {
	t.Header.Set(name, value)
	return t
}

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

func (t *Test) MustStatus(statusCode int) *Test {
	if t.Error != nil {
		return t
	}

	if t.Status != statusCode {
		t.Error = fmt.Errorf("expected status code response of %d, actual %d", statusCode, t.Status)
	}

	return t
}

func (t *Test) MustStringValue(expected, actual string) *Test {
	if t.Error != nil {
		return t
	}

	if expected != actual {
		t.Error = fmt.Errorf("expected %s, but got %s", expected, actual)
	}

	return t
}

func (t *Test) MustIntValue(expected, actual int) *Test {
	if t.Error != nil {
		return t
	}

	if expected != actual {
		t.Error = fmt.Errorf("expected %d, but got %d", expected, actual)
	}

	return t
}

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
