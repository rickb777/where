#!/bin/bash -e
cd $(dirname $0)

unset GOPATH
export GO111MODULE=on

go mod download

### Build Phase 1 ###

gofmt -l -w *.go */*.go
go vet ./...
go test ./...

