package irest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Endpoint represents the general route, method, and parameters for an API
// endpoint.
type Endpoint struct {
	// Path is a relative endpoint route.
	Path string

	// Method is an HTTP method.
	Method string

	// Parameters is the map of name to value in the query parameters.
	Parameters map[string]interface{}
}

type EndpointTest struct {
	Name string

	// Parent is the parent level test containing shared data and slice of endpoint
	// tests.
	Parent *Test

	// Path is a relative endpoint route.
	Path string

	// url is the formatted URL from the relative path and the variables passed in.
	URL string

	// Method is an HTTP method.
	Method string

	// Parameters is the map of name to value in the query parameters.
	Parameters map[string]interface{}

	// Payload is the optional payload to send with the request.
	Payload interface{}

	Client  *http.Client
	Cookies []*http.Cookie
	Header  *http.Header

	Duration int64
	Response *http.Response
	Error    error
}

// Build constructs a usable endpoint with the full URL from the baseURL,
// relative path, and variables.
func (e *Endpoint) Use(baseURL string, payload interface{}, v ...interface{}) *EndpointTest {
	et := &EndpointTest{
		Path:   e.Path,
		Method: e.Method,
		Payload: payload,
		Header: &http.Header{},
		Client: &http.Client{},
	}

	if strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1]
	}

	builtPath := fmt.Sprintf(e.Path, v...)
	if strings.HasPrefix(e.Path, "/") {
		builtPath = builtPath[1:]
	}

	et.URL = fmt.Sprintf("%s/%s", baseURL, builtPath)

	return et
}

// UseHeader uses a previously saved header value by name as a header with the
// provided name.
func (e *EndpointTest) UseHeader(savedName, name string) *EndpointTest {
	savedValue, ok := e.Parent.savedValues[savedName]
	if !ok {
		e.Error = fmt.Errorf("header not found saved as %s", savedName)
		return e
	}
	e.Header.Set(name, savedValue)
	return e
}

// UseCookie adds to the slice of cookies to be included in the request.
func (e *EndpointTest) UseCookie(savedName, name string) *EndpointTest {
	savedValue, ok := e.Parent.savedValues[savedName]
	if !ok {
		e.Error = fmt.Errorf("cookie not found saved as %s", savedName)
		return e
	}
	e.Cookies = append(e.Cookies, &http.Cookie{Name: name, Value: savedValue})
	return e
}

// SaveHeader will save the header if found or the cookie as a fallback if that
// is found instead with the provided name as the savedName in the parent test.
func (e *EndpointTest) SaveHeader(name, savedName string) *EndpointTest {
	if e.Error != nil {
		return e
	}

	if e.Response == nil {
		e.Error = fmt.Errorf("http response not set, must have request before saving result")
		return e
	}

	if value := e.Response.Header.Get(name); value != "" {
		e.Parent.savedValues[savedName] = value
		return e
	}

	for _, c := range e.Response.Cookies() {
		if c.Name == name {
			e.Parent.savedValues[savedName] = c.Value
			return e
		}
	}

	e.Error = fmt.Errorf("header or cookie name '%s' not found", name)

	return e
}

// Do executes the request.
func (e *EndpointTest) Do() *EndpointTest {
	b := new(bytes.Buffer)
	if e.Payload != nil {
		if err := json.NewEncoder(b).Encode(e.Payload); err != nil {
			e.Error = err
			return e
		}
	}

	req, err := http.NewRequest(e.Method, e.URL, b)
	if err != nil {
		e.Error = err
		return e
	}

	req.Header = *e.Header
	for _, c := range e.Cookies {
		req.AddCookie(c)
	}

	startTime := time.Now()
	fmt.Printf("%+v\n", req)
	res, err := e.Client.Do(req)
	if err != nil {
		e.Error = err
		return e
	}
	reqDuration := time.Since(startTime)
	e.Duration = reqDuration.Nanoseconds() / int64(time.Millisecond)

	e.Response = res

	return e
}

// MustStatus sets the EndpointTest.Error if the status code is not the expected
// value. An HTTP request must have been made prior to this function call.
func (e *EndpointTest) MustStatus(statusCode int) *EndpointTest {
	if e.Error != nil {
		return e
	}

	if e.Response.StatusCode != statusCode {
		e.Error = fmt.Errorf("expected status code response of %d, actual %d", statusCode, e.Response.StatusCode)
	}

	return e
}

// ParseResponseBody parses the HTTP response body from json
// to a provided interface.
func (t *EndpointTest) ParseResponseBody(result interface{}) *EndpointTest {
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
