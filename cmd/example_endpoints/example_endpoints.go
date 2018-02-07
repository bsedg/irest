package main

import (
	"net/http"

	"github.com/bsedg/irest"
)

type example struct {
	ID int64
}

func main() {
	loginBase := &irest.Endpoint{Path: "/login", Method: http.MethodPost}
	getExample := &irest.Endpoint{Path: "/examples/%d", Method: http.MethodGet}
	createExample := &irest.Endpoint{Path: "/examples", Method: http.MethodPost}

	t := irest.NewTest("Example")
	// Sets default headers to use throughout tests.
	t.AddHeader("name", "value").AddHeader("name", "value").AddHeader("name", "value")
	ex := &example{}
	t.NewEndpointsTest("Example",
		// TODO(bsedg): only execute if previous test passes
		loginBase.Use("api/", nil).MustStatus(http.StatusOK).Do().SaveHeader("x-authentication", "AUTH"),
		// TODO(bsedg): execute request on last Use before Must...
		createExample.Use("api/", ex).UseHeader("AUTH", "x-authentication").Do().MustStatus(http.StatusCreated).ParseResponseBody(ex),
		getExample.Use("api/", nil, ex.ID).UseHeader("AUTH", "x-authentication").Do().MustStatus(http.StatusOK),
	)
}
