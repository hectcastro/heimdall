FROM golang

ADD . /go/src/github.com/hectcastro/heimdall

WORKDIR /go/src/github.com/hectcastro/heimdall

RUN make deps
