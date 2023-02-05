package where

// OrderBy lists the column(s) by which the database will be asked to sort its results.
// The columns passed in here will be quoted according to the quoter in use when built.
// Be careful not to allow injection attacks: do not include a string from an external
// source in the columns.
func OrderBy(column ...string) *QueryConstraint {
	return &QueryConstraint{orderBy: makeTerms(column)}
}

// Limit sets the upper limit on the number of records to be returned.
// The default value, 0, suppresses any limit.
//
// As a special case, for SQL-Server, this produces the 'TOP' expression (see FormatTOP).
func Limit(n int) *QueryConstraint {
	return &QueryConstraint{limit: n}
}

// Offset sets the offset into the result set; previous items will be discarded.
func Offset(n int) *QueryConstraint {
	return &QueryConstraint{offset: n}
}

// OrderBy lists the column(s) by which the database will be asked to sort its results.
// The columns passed in here will be quoted according to the needs of the selected dialect.
// Be careful not to allow injection attacks: do not include a string from an external
// source in the columns.
func (qc *QueryConstraint) OrderBy(column ...string) *QueryConstraint {
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

func (qc *QueryConstraint) setDirection(dir int) *QueryConstraint {
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
func (qc *QueryConstraint) Asc() *QueryConstraint {
	return qc.setDirection(asc)
}

// Desc sets the sort order to be descending for the columns specified previously,
// not including those already set.
func (qc *QueryConstraint) Desc() *QueryConstraint {
	return qc.setDirection(desc)
}

// NullsFirst can be used to control whether nulls appear before non-null values
// in the sort ordering. By default, null values sort as if larger than any non-null value;
// that is, NULLS FIRST is the default for DESC order, and NULLS LAST otherwise.
func (qc *QueryConstraint) NullsFirst() *QueryConstraint {
	qc.nulls = first
	return qc
}

// NullsLast can be used to control whether nulls appear after non-null values
// in the sort ordering. By default, null values sort as if larger than any non-null value;
// that is, NULLS FIRST is the default for DESC order, and NULLS LAST otherwise.
func (qc *QueryConstraint) NullsLast() *QueryConstraint {
	qc.nulls = last
	return qc
}

// Limit sets the upper limit on the number of records to be returned.
func (qc *QueryConstraint) Limit(n int) *QueryConstraint {
	qc.limit = n
	return qc
}

// Offset sets the offset into the result set. The database will skip earlier records.
// It is usually important to set the order of results explicitly (see OrderBy).
func (qc *QueryConstraint) Offset(n int) *QueryConstraint {
	qc.offset = n
	return qc
}
