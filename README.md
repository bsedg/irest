# iREST

[ ![Codeship Status for bsedg/irest](https://codeship.com/projects/2d4b3280-3e78-0134-9c3a-5218b375052b/status?branch=master)](https://codeship.com/projects/167341)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/bsedg/irest)](https://goreportcard.com/report/github.com/bsedg/irest)

Integration testing framework for creating for RESTful APIs with golang.


```
// Example using iREST to test an API.

import (
    "net/http"

    "github.com/bsedg/irest"
)

func main() {
    t := irest.NewTest("Sample Test")

    s := &Something{}
    t.NewTest("Create something").
        AddHeader("Content-Type", "application/json").
        Post("localhost/somethings").
        MustStatus(http.StatusCreated).
        ParseResponseBody(s)

    t.NewTest("Get something").
        Get("localhost/somethings/" + s.ID).
        MustStatus(http.StatusOK)
}

```
