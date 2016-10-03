# iREST

[ ![Codeship Status for bsedg/irest](https://codeship.com/projects/2d4b3280-3e78-0134-9c3a-5218b375052b/status?branch=master)](https://codeship.com/projects/167341)
[![GoDoc](https://godoc.org/github.com/bsedg/irest?status.svg)](http://godoc.org/github.com/bsedg/irest)
[![Go Report Card](https://goreportcard.com/badge/github.com/bsedg/irest)](https://goreportcard.com/report/github.com/bsedg/irest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

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

## Development

Add any needed tests, then run the tests to make sure nothing breaks:

`go test ./...`

### Running the example

Build the binary:

`go build cmd/example_report/example_report.go`

Run the example:

`./example_report`

![Image Example Report](./docs/example_report.png)

### Issues, bugs, feature requests

Create a new issue for any bug or feature request. If a bug is found or unexpected behavior, create an issue that clearly outlines what was expected and what actually happened. Design proposals of feature requests are welcomed as well.
