package where

import (
	"bytes"
	"reflect"
)

const and = " AND "
const or = " OR "

// Null returns an 'IS NULL' condition on a column.
func Null(column string) Expression {
	return Condition{column, " IS NULL", []interface{}{}}
}

// NotNull returns an 'IS NOT NULL' condition on a column. It's also possible to use Not(Null(...)).
func NotNull(column string) Expression {
	return Condition{column, " IS NOT NULL", []interface{}{}}
}

// Eq returns an equality condition on a column.
func Eq(column string, value interface{}) Expression {
	return Condition{column, "=?", []interface{}{value}}
}

// NotEq returns a not equal condition on a column.
func NotEq(column string, value interface{}) Expression {
	return Condition{column, "<>?", []interface{}{value}}
}

// Gt returns a greater than condition on a column.
func Gt(column string, value interface{}) Expression {
	return Condition{column, ">?", []interface{}{value}}
}

// GtEq returns a greater than or equal condition on a column.
func GtEq(column string, value interface{}) Expression {
	return Condition{column, ">=?", []interface{}{value}}
}

// Lt returns a less than condition on a column.
func Lt(column string, value interface{}) Expression {
	return Condition{column, "<?", []interface{}{value}}
}

// LtEq returns a less than or equal than condition on a column.
func LtEq(column string, value interface{}) Expression {
	return Condition{column, "<=?", []interface{}{value}}
}

// Between returns a between condition on a column.
func Between(column string, a, b interface{}) Expression {
	return Condition{column, " BETWEEN ? AND ?", []interface{}{a, b}}
}

// Like returns a pattern-matching condition on a column. Be careful: this can hurt performance.
func Like(column string, pattern string) Expression {
	return Condition{column, " LIKE ?", []interface{}{pattern}}
}

// In returns an in condition on a column.
func In(column string, values ...interface{}) Expression {
	if len(values) == 0 {
		return NoOp()
	}

	v2 := make([]interface{}, 0, len(values))
	hasNull := false
	buf := &bytes.Buffer{}
	buf.WriteString(" IN (")
	i := 0
	for _, arg := range values {
		switch arg.(type) {
		case nil:
			hasNull = true
		default:
			v2 = append(v2, arg)
			value := reflect.ValueOf(arg)
			switch value.Kind() {
			case reflect.Array, reflect.Slice:
				for j := 0; j < value.Len(); j++ {
					if i > 0 {
						buf.WriteByte(',')
					}
					buf.WriteByte('?')
					i++
				}

			default:
				if i > 0 {
					buf.WriteByte(',')
				}
				buf.WriteByte('?')
				i++
			}
		}
	}
	buf.WriteByte(')')

	result := NoOp()
	if i > 0 {
		result = Condition{column, buf.String(), v2}
	}

	if hasNull {
		result = Or(result, Null(column))
	}

	return result
}

//-------------------------------------------------------------------------------------------------

// Not negates an expression.
func Not(el Expression) Expression {
	return not{el}
}

// And combines two conditions into a clause that requires they are both true.
func (cl not) And(c2 Expression) Expression {
	return Clause{[]Expression{cl}, and}.And(c2)
}

// Or combines two conditions into a clause that requires either is true.
func (cl not) Or(c2 Expression) Expression {
	return Clause{[]Expression{cl}, or}.Or(c2)
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
	return Clause{[]Expression{cl}, and}.And(c2)
}

// Or combines two conditions into a clause that requires either is true.
func (cl Condition) Or(c2 Expression) Expression {
	return Clause{[]Expression{cl}, or}.Or(c2)
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
	return Clause{[]Expression{wh, exp}, conj}
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
func And(exp ...Expression) Expression {
	return newClause(and, exp...)
}

// Or combines some expressions into a clause that requires that any is true.
func Or(exp ...Expression) Expression {
	return newClause(or, exp...)
}

func newClause(conj string, exp ...Expression) Expression {
	clause := Clause{nil, conj}
	for _, e := range exp {
		cl, isClause := e.(Clause)
		if !isClause || len(cl.wheres) > 0 {
			clause.wheres = append(clause.wheres, e)
		}
	}
	return clause
}
