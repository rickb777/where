package where

import (
	"strconv"
	"strings"

	"github.com/rickb777/where/v2/dialect"
)

const (
	unset = 0
	asc   = 1
	desc  = 2
	first = 3
	last  = 4
)

var ascDesc = []string{
	" ASC",
	" ASC",
	" DESC",
	" FIRST",
	" LAST",
}

type orderingTerm struct {
	column string
	dir    int
}

type QueryConstraint struct {
	orderBy       []orderingTerm
	nulls         int
	limit, offset int
}

//var _ QueryConstraint = &queryConstraint{}

// Format formats the SQL expressions.
func (qc *QueryConstraint) Format(d dialect.Dialect, option ...dialect.FormatOption) string {
	if qc == nil {
		return ""
	}

	b := new(strings.Builder)
	b.Grow(qc.estimateStringLength())

	q := quoterFromOptions(option)

	if len(qc.orderBy) > 0 {
		b.WriteString(" ORDER BY")
		hasDesc := false

		for _, col := range qc.orderBy {
			if col.dir == desc {
				hasDesc = true
				break
			}
		}

		sep := " "
		for _, col := range qc.orderBy {
			b.WriteString(sep)
			q.QuoteW(b, col.column)
			if hasDesc {
				b.WriteString(ascDesc[col.dir])
			}
			sep = ", "
		}

		switch qc.nulls {
		case first:
			b.WriteString(" NULLS FIRST")
		case last:
			b.WriteString(" NULLS LAST")
		}
	}

	if qc.limit > 0 && d != dialect.SqlServer {
		b.WriteString(" LIMIT ")
		b.WriteString(strconv.Itoa(qc.limit))
	}

	if qc.offset > 0 {
		b.WriteString(" OFFSET ")
		b.WriteString(strconv.Itoa(qc.offset))
	}

	return b.String()
}

// FormatTOP formats the SQL 'TOP' expression using the given dialect. Only SQL-Server uses this;
// for other dialects, it returns an empty string. Insert the returned string into your query
// after "SELECT [DISTINCT] " and before the list of column names.
func (qc *QueryConstraint) FormatTOP(d dialect.Dialect) string {
	if qc == nil {
		return ""
	}

	if d != dialect.SqlServer || qc.limit == 0 {
		return ""
	}

	b := new(strings.Builder)
	b.Grow(12)
	b.WriteString(" TOP (")
	b.WriteString(strconv.Itoa(qc.limit))
	b.WriteString(")")

	return b.String()
}

func (qc *QueryConstraint) estimateStringLength() (n int) {
	if len(qc.orderBy) > 0 {
		n += 14 // " ORDER BY" and " DESC"
		for _, col := range qc.orderBy {
			n += len(col.column) + 4 // allow for 2 quote marks, space and comma
		}
	}

	if qc.limit > 0 {
		n += 13 // " LIMIT " + number
	}

	if qc.offset > 0 {
		n += 14 // " OFFSET " + number
	}

	return n
}

func (qc *QueryConstraint) String() string {
	return qc.Format(dialect.DefaultDialect)
}
