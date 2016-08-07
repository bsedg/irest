package irest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Test struct {
	Name   string `json:"name"`
	Error  error
	Status int

	Tests    []*Test
	Errors   []error
	Duration int64

	Client *http.Client
	Header *http.Header
}

func NewTest(name string) *Test {
	t := &Test{
		Name:   name,
		Error:  nil,
		Tests:  []*Test{},
		Errors: []error{},
		Client: &http.Client{},
		Header: &http.Header{},
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

func (t *Test) Post(baseURL, endpoint string, v interface{}) *Test {
	t.Header.Set("Content-Type", "application/json")

	addr, err := url.Parse(baseURL + endpoint)
	if err != nil {
		t.Error = err
		return t
	}

	req, err := http.NewRequest(http.MethodPost, addr.String(), nil)
	if err != nil {
		t.Error = err
		return t
	}

	res, err := t.Client.Do(req)
	if err != nil {
		t.Error = err
		return t
	}

	t.Status = res.StatusCode

	body := res.Body
	defer body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Error = err
		return t
	}

	json.Unmarshal(data, v)

	return t
}
