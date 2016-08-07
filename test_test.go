package irest

import (
	"encoding/json"
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
	api = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				data, _ := json.Marshal(SampleObject{
					Name:    "unit-test",
					Value:   100,
					Success: true,
				})
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write(data)
			}
		}))
}

func TestIRestSetup(t *testing.T) {
	test := NewTest("unit-test")

	if test.Name != "unit-test" {
		t.Errorf("name was %s, expected unit-test", test.Name)
	}
}

func TestIRestPost(t *testing.T) {
    test := NewTest("unit-test")

    sample := SampleObject{}
    test = test.Post(api.URL, "/tests", &sample)
    if sample.Name != "unit-test" {
        t.Errorf("name response was %s, expected unit-test", sample.Name)
    }
    if !sample.Success {
        t.Error("expected response success to be true")
    }
}
