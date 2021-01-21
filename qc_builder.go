package where

import (
	"strconv"
	"strings"

	"github.com/rickb777/where/dialect"
)

const (
	unset = 0
	asc   = 1
	desc  = 2
)

var ascDesc = []string{
	" ASC",
	" ASC",
	" DESC",
}

type orderingTerm struct {
	column string
	dir    int
}

type queryConstraint struct {
	orderBy       []orderingTerm
	limit, offset int
}

var _ QueryConstraint = &queryConstraint{}

// Build constructs the SQL string using the optional quoter or the default quoter.
func (qc *queryConstraint) Build(d dialect.Dialect) string {
	q := d.Config().Quoter
	b := new(strings.Builder)
	b.Grow(qc.estimateStringLength())

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

// BuildTop constructs the SQL string using the given dialect. The only known dialect
// for which this is used is SQL-Server; otherwise it returns an empty string. Insert
// the returned value into your query between "SELECT [DISTINCT] " and the list of columns.
func (qc *queryConstraint) BuildTop(d dialect.Dialect) string {
	if d != dialect.SqlServer || qc.limit == 0 {
		return ""
	}

	b := new(strings.Builder)
	b.Grow(10)
	b.WriteString(" TOP (")
	b.WriteString(strconv.Itoa(qc.limit))
	b.WriteString(")")

	return b.String()
}

func (qc *queryConstraint) estimateStringLength() (n int) {
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

func (qc *queryConstraint) String() string {
	return qc.Build(dialect.DefaultDialect)
}
