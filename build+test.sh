#!/bin/bash -e
cd "$(dirname $0)"

PATH=$HOME/go/bin:$PATH
unset GOPATH
export GO111MODULE=on

function v
{
  echo
  echo "$@"
  $@
}

go mod download

if ! type -p shadow; then
  v go get     golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
  v go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
fi

### Build Phase 1 ###

v go test ./...
v gofmt -l -w *.go */*.go
v go vet ./...
v shadow -strict ./...

if [ -n "$COVERALLS_TOKEN" ]; then
  if ! type -p goveralls; then
    v go get github.com/mattn/goveralls
  fi

  v go test . -covermode=count -coverprofile=dot.out .
  v go tool cover -func=dot.out
  v goveralls -coverprofile=dot.out -service=travis-ci -repotoken $COVERALLS_TOKEN
fi
