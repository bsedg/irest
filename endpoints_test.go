package irest

import (
	"testing"
)

func TestUse(t *testing.T) {
	e := &Endpoint{
		Path: "/things/%d",
	}

	eb := e.Use("api/", nil, 1)
	if eb.URL != "api/things/1" {
		t.Errorf("expected api/things/1, got %s", eb.URL)
	}
}

func TestUseParameterList(t *testing.T) {
	var endpointBuildTests = []struct {
		in     Endpoint
		base   string
		params []interface{}
		out    string
	}{
		{
			in: Endpoint{
				Path: "/things/%d/sub/%d",
			},
			base:   "api/base/",
			params: []interface{}{1, 100},
			out:    "api/base/things/1/sub/100",
		},
		{
			in: Endpoint{
				Path: "/things/%d/sub/%d",
			},
			base:   "api/base",
			params: []interface{}{1, 1},
			out:    "api/base/things/1/sub/1",
		},
		{
			in: Endpoint{
				Path: "things/%d/sub/%d",
			},
			base:   "api/base",
			params: []interface{}{1, 1},
			out:    "api/base/things/1/sub/1",
		},
	}

	for _, ebt := range endpointBuildTests {
		eb := ebt.in.Use(ebt.base, nil, ebt.params...)
		if eb.URL != ebt.out {
			t.Errorf("expected %s, got %s", ebt.out, eb.URL)
		}
	}
}

func TestEndpointParseResponseBodyEmptyResponse(t *testing.T) {
	test := &EndpointTest{Name: "unit-test"}

	test = test.ParseResponseBody(nil)
	if test.Error == nil {
		t.Errorf("expected an error, but did not get one")
	}
}
