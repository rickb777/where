# where

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rickb777/where)
[![Build Status](https://travis-ci.org/rickb777/where.svg?branch=master)](https://travis-ci.org/rickb777/where/builds)
[![Code Coverage](https://img.shields.io/coveralls/rickb777/where.svg)](https://coveralls.io/r/rickb777/where)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/where)](https://goreportcard.com/report/github.com/rickb777/where)
[![Issues](https://img.shields.io/github/issues/rickb777/where.svg)](https://github.com/rickb777/where/issues)

* Provides a fluent API for dynamically constructing SQL `WHERE` & `HAVING` clauses.
* Also supports dynamic `LIMIT`, `OFFSET` and `ORDER BY` clauses. 
* Allows the identifiers to be quoted to suit different SQL dialects, or not at all.
* `dialect` package supports different placeholder styles.
* `quote` package supports quoting SQL identifiers in back-ticks, double quotes, square barckets, or nothing.

## Install

Install with this command:

```
go get github.com/rickb777/where/v2
```

## where

Package `where` provides composable expressions for **WHERE** and **HAVING** clauses in SQL.
These can range from the very simplest no-op to complex nested trees of **AND** and **OR**
conditions.

In the naive approach, strings can be concatenated to construct lists of expression that are
AND-ed together. However, mixing AND with OR makes things much more difficult. So this package
does the work for you.

Also in this package are query constraints to provide **ORDER BY**, **LIMIT** and **OFFSET**
clauses, along with 'TOP' for MS-SQL. These are similar to **WHERE** clauses except literal values
are used instead of parameter placeholders.

Further support for SQL dialects and formatting options is provided in the `dialect` sub-package.

Queries should be written using '?' query placeholders throughout, and then these can be translated
to the form needed by the chosen dialect: one of `dialect.Query`, `dialect.Dollar`, `dialect.AtP` or
`dialect.Inline`.

Also, support for quoted identifiers is provided in the `quote` sub-package.
  - `quote.Quoter` is the interface for a quoter.
  - implementations include `quote.ANSI`, `quote.Backticks`, `quote.SquareBrackets`, and `quote.None`.
