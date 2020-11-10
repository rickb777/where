package where

import (
	"fmt"

	"github.com/rickb777/where/dialect"

	"github.com/rickb777/where/quote"
)

// QueryConstraint is a value that is appended to a SELECT statement.
type QueryConstraint interface {
	fmt.Stringer
	BuildTop(dialect.Dialect) string
	Build(dialect.Dialect) string
}

func pickQuoter(quoter []quote.Quoter) quote.Quoter {
	if len(quoter) > 0 {
		return quoter[0]
	}
	return quote.DefaultQuoter
}

// BuildTop builds a query constraint as used by SQL-Server. It allows nil values.
func BuildTop(qc QueryConstraint, d dialect.Dialect) string {
	if qc == nil {
		return ""
	}
	return qc.BuildTop(d)
}

// Build builds a query constraint. It allows nil values.
func Build(qc QueryConstraint, d dialect.Dialect) string {
	if qc == nil {
		return ""
	}
	return qc.Build(d)
}

// OrderBy lists the column(s) by which the database will be asked to sort its results.
// The columns passed in here will be quoted according to the quoter in use when built.
func OrderBy(column ...string) *queryConstraint {
	return &queryConstraint{orderBy: makeTerms(column)}
}

// Limit sets the upper limit on the number of records to be returned.
// The default value, 0, suppresses any limit.
func Limit(n int) *queryConstraint {
	return &queryConstraint{limit: n}
}

// Offset sets the offset into the result set; previous items will be discarded.
func Offset(n int) *queryConstraint {
	return &queryConstraint{offset: n}
}

// OrderBy lists the column(s) by which the database will be asked to sort its results.
// The columns passed in here will be quoted according to the needs of the current dialect.
func (qc *queryConstraint) OrderBy(column ...string) *queryConstraint {
	qc.orderBy = append(qc.orderBy, makeTerms(column)...)
	return qc
}

func makeTerms(column []string) []orderingTerm {
	terms := make([]orderingTerm, len(column))
	for i, c := range column {
		terms[i] = orderingTerm{column: c} // n.b. dir: unset
	}
	return terms
}

func (qc *queryConstraint) setDir(dir int) *queryConstraint {
	for i := 0; i < len(qc.orderBy); i++ {
		if qc.orderBy[i].dir == unset {
			qc.orderBy[i].dir = dir
		}
	}
	return qc
}

// Asc sets the sort order to be ascending for all the columns specified previously,
// not including those already set.
func (qc *queryConstraint) Asc() *queryConstraint {
	return qc.setDir(asc)
}

// Desc sets the sort order to be descending for all the columns specified previously,
// not including those already set.
func (qc *queryConstraint) Desc() *queryConstraint {
	return qc.setDir(desc)
}

// Limit sets the upper limit on the number of records to be returned.
func (qc *queryConstraint) Limit(n int) *queryConstraint {
	qc.limit = n
	return qc
}

// Offset sets the offset into the result set; previous items will be discarded.
func (qc *queryConstraint) Offset(n int) *queryConstraint {
	qc.offset = n
	return qc
}
