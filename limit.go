package where

import (
	"fmt"
	"github.com/rickb777/where/quote"
)

// QueryConstraint is a value that is appended to a SELECT statement.
type QueryConstraint interface {
	fmt.Stringer
	Build(q ...quote.Quoter) string
}

func pickQuoter(quoter []quote.Quoter) quote.Quoter {
	if len(quoter) > 0 {
		return quoter[0]
	}
	return quote.DefaultQuoter
}

// Build builds a query constraint. It allows nil values.
func Build(qc QueryConstraint, q ...quote.Quoter) string {
	if qc == nil {
		return ""
	}
	return qc.Build(q...)
}

//-------------------------------------------------------------------------------------------------

type literal string

// Literal returns the literal string supplied, converting it to a QueryConstraint.
// The string may contain identifiers, however no quoting rules will be applied.
// Therefore care must be taken if portability is needed.
func Literal(sqlPart string) QueryConstraint {
	return literal(sqlPart)
}

func (qc literal) Build(_ ...quote.Quoter) string {
	return string(qc)
}

func (qc literal) String() string {
	return string(qc)
}
