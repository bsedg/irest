FROM golang:1.9

WORKDIR /go/src/github.com/bsedg/irest
COPY . .

RUN go get -u github.com/golang/lint/golint

RUN make all
