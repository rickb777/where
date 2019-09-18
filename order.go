package where

import (
	"strconv"
	"strings"

	"github.com/rickb777/where/quote"
)

const (
	asc  = 1
	desc = 2
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

// OrderBy lists the column(s) by which the database will be asked to sort its results.
// The columns passed in here will be quoted according to the quoter in use when built.
func OrderBy(column ...string) *queryConstraint {
	return &queryConstraint{orderBy: makeTerms(column)}
}

func makeTerms(column []string) []orderingTerm {
	terms := make([]orderingTerm, len(column))
	for i, c := range column {
		terms[i] = orderingTerm{column: c}
	}
	return terms
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

// Asc sets the sort order to be ascending for all the columns specified previously.
func (qc *queryConstraint) setDir(dir int) *queryConstraint {
	for i := 0; i < len(qc.orderBy); i++ {
		if qc.orderBy[i].dir == 0 {
			qc.orderBy[i].dir = dir
		}
	}
	return qc
}

// Asc sets the sort order to be ascending for all the columns specified previously.
func (qc *queryConstraint) Asc() *queryConstraint {
	return qc.setDir(asc)
}

// Desc sets the sort order to be descending for all the columns specified previously.
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

func ascDesc(dir int, b quote.StringWriter) {
	switch dir {
	case asc:
		b.WriteString(" ASC")
	case desc:
		b.WriteString(" DESC")
	}
}

func intTerm(b quote.StringWriter, spacer, noun string, value int) {
	b.WriteString(spacer)
	b.WriteString(noun)
	b.WriteString(strconv.Itoa(value))
}

// Build constructs the SQL string using the optional quoter or the default quoter.
func (qc *queryConstraint) Build(quoter ...quote.Quoter) string {
	q := pickQuoter(quoter)
	b := &strings.Builder{}
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
				ascDesc(col.dir, b)
			} else if col.dir != qc.orderBy[i+1].dir {
				ascDesc(col.dir, b)
			}
		}

		spacer = " "
	}

	if qc.limit > 0 {
		intTerm(b, spacer, "LIMIT ", qc.limit)
		spacer = " "
	}

	if qc.offset > 0 {
		intTerm(b, spacer, "OFFSET ", qc.offset)
	}

	return b.String()
}

func (qc *queryConstraint) String() string {
	return qc.Build(quote.DefaultQuoter)
}
