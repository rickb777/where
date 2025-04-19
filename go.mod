module github.com/rickb777/where/v2

go 1.24.1

require github.com/rickb777/expect v0.21.0

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/mattn/goveralls v0.0.12 // indirect
	github.com/rickb777/plural v1.4.3 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/tools v0.32.0 // indirect
)

tool (
	github.com/mattn/goveralls
	golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
)
