package where

import (
	"fmt"

	"github.com/rickb777/where/dialect"

	"github.com/rickb777/where/quote"
)

// QueryConstraint is a value that is appended to a SELECT statement.
type QueryConstraint interface {
	fmt.Stringer

	// OrderBy lists the column(s) by which the database will be asked to sort its results.
	// The columns passed in here will be quoted according to the needs of the selected dialect.
	OrderBy(column ...string) QueryConstraint

	// Asc sets the sort order to be ascending for the columns specified previously,
	// not including those already set.
	Asc() QueryConstraint

	// Desc sets the sort order to be descending for the columns specified previously,
	// not including those already set.
	Desc() QueryConstraint

	// NullsFirst can be used to determine whether nulls appear before non-null values
	// in the sort ordering. By default, null values sort as if larger than any non-null value;
	// that is, NULLS FIRST is the default for DESC order, and NULLS LAST otherwise.
	NullsFirst() QueryConstraint

	// NullsLast can be used to determine whether nulls appear after non-null values
	// in the sort ordering. By default, null values sort as if larger than any non-null value;
	// that is, NULLS FIRST is the default for DESC order, and NULLS LAST otherwise.
	NullsLast() QueryConstraint

	// Limit sets the upper limit on the number of records to be returned.
	Limit(n int) QueryConstraint

	// Offset sets the offset into the result set; previous items will be discarded.
	Offset(n int) QueryConstraint

	// BuildTop constructs the SQL string using the given dialect. The only known dialect
	// for which this is used is SQL-Server; otherwise it returns an empty string. Insert
	// the returned value into your query between "SELECT [DISTINCT] " and the list of columns.
	BuildTop(dialect.Dialect) string

	// Build constructs the SQL string using the optional quoter or the default quoter.
	Build(dialect.Dialect) string
}

func pickQuoter(quoter []quote.Quoter) quote.Quoter {
	if len(quoter) > 0 {
		return quoter[0]
	}
	return quote.DefaultQuoter
}

// BuildTop builds a query constraint as used by SQL-Server. It allows nil values.
// The only known dialect for which this is used is SQL-Server; otherwise it returns
// an empty string. Insert the returned value into your query between "SELECT [DISTINCT] "
// and the list of columns.
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
// Be careful not to allow injection attacks: do not include a string from an external
// source in the columns.
func OrderBy(column ...string) QueryConstraint {
	return &queryConstraint{orderBy: makeTerms(column)}
}

// Limit sets the upper limit on the number of records to be returned.
// The default value, 0, suppresses any limit.
func Limit(n int) QueryConstraint {
	return &queryConstraint{limit: n}
}

// Offset sets the offset into the result set; previous items will be discarded.
func Offset(n int) QueryConstraint {
	return &queryConstraint{offset: n}
}

// OrderBy lists the column(s) by which the database will be asked to sort its results.
// The columns passed in here will be quoted according to the needs of the selected dialect.
// Be careful not to allow injection attacks: do not include a string from an external
// source in the columns.
func (qc *queryConstraint) OrderBy(column ...string) QueryConstraint {
	// previous unset columns default to asc
	for i := 0; i < len(qc.orderBy); i++ {
		if qc.orderBy[i].dir == unset {
			qc.orderBy[i].dir = asc
		}
	}

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
	for i := len(qc.orderBy) - 1; i >= 0; i-- {
		if qc.orderBy[i].dir == unset {
			qc.orderBy[i].dir = dir
		} else {
			return qc
		}
	}
	return qc
}

// Asc sets the sort order to be ascending for the columns specified previously,
// not including those already set.
func (qc *queryConstraint) Asc() QueryConstraint {
	return qc.setDir(asc)
}

// Desc sets the sort order to be descending for the columns specified previously,
// not including those already set.
func (qc *queryConstraint) Desc() QueryConstraint {
	return qc.setDir(desc)
}

// NullsFirst can be used to determine whether nulls appear before non-null values
// in the sort ordering. By default, null values sort as if larger than any non-null value;
// that is, NULLS FIRST is the default for DESC order, and NULLS LAST otherwise.
func (qc *queryConstraint) NullsFirst() QueryConstraint {
	qc.nulls = first
	return qc
}

// NullsLast can be used to determine whether nulls appear after non-null values
// in the sort ordering. By default, null values sort as if larger than any non-null value;
// that is, NULLS FIRST is the default for DESC order, and NULLS LAST otherwise.
func (qc *queryConstraint) NullsLast() QueryConstraint {
	qc.nulls = last
	return qc
}

// Limit sets the upper limit on the number of records to be returned.
func (qc *queryConstraint) Limit(n int) QueryConstraint {
	qc.limit = n
	return qc
}

// Offset sets the offset into the result set; previous items will be discarded.
func (qc *queryConstraint) Offset(n int) QueryConstraint {
	qc.offset = n
	return qc
}
