# where

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/rickb777/where)
[![Build Status](https://travis-ci.org/rickb777/where.svg?branch=master)](https://travis-ci.org/rickb777/where)
[![Code Coverage](https://img.shields.io/coveralls/rickb777/where.svg)](https://coveralls.io/r/rickb777/where)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/where)](https://goreportcard.com/report/github.com
/rickb777/where)
[![Issues](https://img.shields.io/github/issues/rickb777/where.svg)](https://github.com/rickb777/where/issues)

* Provides a fluent API for dynamically constructing SQL 'where' clauses.
* Also supports dynamic `LIMIT`, `OFFSET` and `ORDER BY` clauses. 
* Allows the identifiers to be quoted to suit different SQL dialects, or not at all.
* `dialect` package supports different placeholder styles.
* `quote` package supports quoting SQL identifiers in back-ticks, double quotes, or nothing.

## Install

Install with this command:

```
go get github.com/rickb777/where
```

