package where

import (
	"strings"
)

const and = " AND "
const or = " OR "

// Null returns an 'IS NULL' condition on a column.
func Null(column string) Expression {
	return Condition{Column: column, Predicate: " IS NULL", Args: array()}
}

// NotNull returns an 'IS NOT NULL' condition on a column. It's also possible to use Not(Null(...)).
func NotNull(column string) Expression {
	return Condition{Column: column, Predicate: " IS NOT NULL", Args: array()}
}

// Literal returns a literal condition on a column. For example
//
//   Literal("age", " > 45")
//
// Be careful not to allow injection attacks by including a string from an external source
// in the predicate.
func Literal(column, predicate string) Expression {
	return Condition{Column: column, Predicate: predicate}
}

// Eq returns an equality condition on a column.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func Eq(column string, value interface{}) Expression {
	return Condition{Column: column, Predicate: "=?", Args: array(value)}
}

// NotEq returns a not equal condition on a column.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func NotEq(column string, value interface{}) Expression {
	return Condition{Column: column, Predicate: "<>?", Args: array(value)}
}

// Gt returns a greater than condition on a column.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func Gt(column string, value interface{}) Expression {
	return Condition{Column: column, Predicate: ">?", Args: array(value)}
}

// GtEq returns a greater than or equal condition on a column.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func GtEq(column string, value interface{}) Expression {
	return Condition{Column: column, Predicate: ">=?", Args: array(value)}
}

// Lt returns a less than condition on a column.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func Lt(column string, value interface{}) Expression {
	return Condition{Column: column, Predicate: "<?", Args: array(value)}
}

// LtEq returns a less than or equal than condition on a column.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func LtEq(column string, value interface{}) Expression {
	return Condition{Column: column, Predicate: "<=?", Args: array(value)}
}

// Between returns a between condition on a column.
//
// Two '?' placeholders are used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func Between(column string, a, b interface{}) Expression {
	return Condition{Column: column, Predicate: " BETWEEN ? AND ?", Args: array(a, b)}
}

// Like returns a pattern-matching condition on a column. Be careful: this can hurt performance.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func Like(column string, pattern string) Expression {
	return Condition{Column: column, Predicate: " LIKE ?", Args: array(pattern)}
}

// In returns an 'IN' condition on a column.
// * If there a no values, this becomes a no-op.
// * If any of the values is itself a slice or array, it is expanded to use all the contained values.
// * If any value is nil, an 'IS NULL' expression is OR-ed with the 'IN' expression.
//
// Some '?' placeholders are used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func In(column string, values ...interface{}) Expression {
	if len(values) == 0 {
		return NoOp()
	}

	args := make([]interface{}, 0, len(values))
	hasNull := false
	buf := &strings.Builder{}
	buf.WriteString(" IN (")
	i := 0
	for _, arg := range values {
		switch arg.(type) {
		case nil:
			hasNull = true
		default:
			args = append(args, arg)
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteByte('?')
			i++
		}
	}
	buf.WriteByte(')')

	result := NoOp()
	if i > 0 {
		result = Condition{Column: column, Predicate: buf.String(), Args: args}
	}

	if hasNull {
		result = Or(result, Null(column))
	}

	return result
}

//-------------------------------------------------------------------------------------------------

// Not negates an expression.
func Not(el Expression) Expression {
	return not{expression: el}
}

// And combines two conditions into a clause that requires they are both true.
func (cl not) And(c2 Expression) Expression {
	return Clause{wheres: []Expression{cl}, conjunction: and}.And(c2)
}

// Or combines two conditions into a clause that requires either is true.
func (cl not) Or(c2 Expression) Expression {
	return Clause{wheres: []Expression{cl}, conjunction: or}.Or(c2)
}

//-------------------------------------------------------------------------------------------------

// NoOp creates an empty expression. This is useful for conditionally chaining
// expression-based contextual decisions. It can also be passed to any method
// that need an expression but for which none is required in that case.
func NoOp() Expression {
	return Clause{}
}

//-------------------------------------------------------------------------------------------------

// And combines two conditions into a clause that requires they are both true.
func (cl Condition) And(c2 Expression) Expression {
	return Clause{wheres: []Expression{cl}, conjunction: and}.And(c2)
}

// Or combines two conditions into a clause that requires either is true.
func (cl Condition) Or(c2 Expression) Expression {
	return Clause{wheres: []Expression{cl}, conjunction: or}.Or(c2)
}

//-------------------------------------------------------------------------------------------------

// And combines two clauses into a clause that requires they are both true.
// SQL implementation note: AND has higher precedence than OR.
func (wh Clause) conjoin(exp Expression, conj string) Expression {
	cl, isClause := exp.(Clause)
	if isClause {
		if len(wh.wheres) == 0 {
			return cl
		} else if len(cl.wheres) == 0 {
			return wh
		} else if wh.conjunction == conj && cl.conjunction == conj {
			return Clause{append(wh.wheres, cl.wheres...), conj}
		}
	} else {
		// blank case comes from NoOp
		if wh.conjunction == "" || wh.conjunction == conj {
			return Clause{append(wh.wheres, exp), conj}
		}
	}
	return Clause{wheres: []Expression{wh, exp}, conjunction: conj}
}

// And combines two clauses into a clause that requires they are both true.
// SQL implementation note: AND has higher precedence than OR.
func (wh Clause) And(exp Expression) Expression {
	return wh.conjoin(exp, and)
}

// Or combines two clauses into a clause that requires either is true.
// SQL implementation note: AND has higher precedence than OR.
func (wh Clause) Or(exp Expression) Expression {
	return wh.conjoin(exp, or)
}

//-------------------------------------------------------------------------------------------------

// And combines some expressions into a clause that requires they are all true.
// Any nil items are silently dropped.
func And(exp ...Expression) Expression {
	return newClause(and, exp...)
}

// Or combines some expressions into a clause that requires that any is true.
// Any nil items are silently dropped.
func Or(exp ...Expression) Expression {
	return newClause(or, exp...)
}

func newClause(conj string, exp ...Expression) Expression {
	var clause = Clause{nil, conj}
	for _, e := range exp {
		if e != nil {
			cl, isClause := e.(Clause)
			if !isClause || len(cl.wheres) > 0 {
				clause.wheres = append(clause.wheres, e)
			}
		}
	}
	if len(clause.wheres) == 1 {
		return clause.wheres[0] // simplify the result
	}
	return clause
}

func array(value ...interface{}) []interface{} {
	return value
}
