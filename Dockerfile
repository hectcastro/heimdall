FROM golang:1.6

RUN go get -u github.com/tools/godep \
    && go install github.com/tools/godep \
    && go get -u github.com/mitchellh/gox \
    && go install github.com/mitchellh/gox

COPY . /go/src/github.com/hectcastro/heimdall

WORKDIR /go/src/github.com/hectcastro/heimdall
