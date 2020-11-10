package where

import (
	"io"
	"strconv"
	"strings"

	"github.com/rickb777/where/dialect"
)

const (
	unset = 0
	asc   = 1
	desc  = 2
)

type orderingTerm struct {
	column string
	dir    int
}

type queryConstraint struct {
	orderBy       []orderingTerm
	limit, offset int
}

var _ QueryConstraint = &queryConstraint{}

var ascDesc = []string{
	"",
	" ASC",
	" DESC",
}

func intTerm(b io.StringWriter, spacer, noun string, value int) {
	b.WriteString(spacer)
	b.WriteString(noun)
	b.WriteString(strconv.Itoa(value))
}

// Build constructs the SQL string using the optional quoter or the default quoter.
func (qc *queryConstraint) Build(d dialect.Dialect) string {
	q := d.Config().Quoter
	b := new(strings.Builder)
	b.Grow(qc.estimateStringLength())

	spacer := ""

	if len(qc.orderBy) > 0 {
		b.WriteString("ORDER BY")
		last := len(qc.orderBy) - 1

		for i, col := range qc.orderBy {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteByte(' ')
			q.QuoteW(b, col.column)
			if i == last {
				b.WriteString(ascDesc[col.dir])
			} else if col.dir != qc.orderBy[i+1].dir {
				b.WriteString(ascDesc[col.dir])
			}
		}

		spacer = " "
	}

	if qc.limit > 0 && d != dialect.SqlServer {
		intTerm(b, spacer, "LIMIT ", qc.limit)
		spacer = " "
	}

	if qc.offset > 0 {
		intTerm(b, spacer, "OFFSET ", qc.offset)
	}

	return b.String()
}

// BuildTop constructs the SQL string using the given dialect. The only known dialect
// for which this is used is SQL-Server; otherwise it returns an empty string. Insert
// the returned value into your query between "SELECT [DISTINCT] " and the list of columns.
func (qc *queryConstraint) BuildTop(d dialect.Dialect) string {
	if d != dialect.SqlServer {
		return ""
	}

	b := new(strings.Builder)
	b.Grow(10)

	if qc.limit > 0 {
		b.WriteString("TOP (")
		b.WriteString(strconv.Itoa(qc.limit))
		b.WriteString(")")
	}

	return b.String()
}

func (qc *queryConstraint) estimateStringLength() (n int) {
	if len(qc.orderBy) > 0 {
		n += 13 // "ORDER BY" and " DESC"
		for _, col := range qc.orderBy {
			n += len(col.column) + 4 // allow for 2 quote marks, space and comma
		}
	}

	if qc.limit > 0 {
		n += 12 // "LIMIT " + number
	}

	if qc.offset > 0 {
		n += 13 // "OFFSET " + number
	}

	return n
}

func (qc *queryConstraint) String() string {
	return qc.Build(dialect.DefaultDialect)
}
