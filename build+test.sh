#!/bin/bash -e
cd "$(dirname $0)"

unset GOPATH
export GO111MODULE=on

go mod download

### Build Phase 1 ###

gofmt -l -w *.go */*.go
go vet ./...
go test ./...

if [ -n "$COVERALLS_TOKEN" ]; then
  go test . -covermode=count -coverprofile=dot.out .
  go tool cover -func=dot.out
  goveralls -coverprofile=dot.out -service=travis-ci -repotoken $COVERALLS_TOKEN
fi
