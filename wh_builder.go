package where

import (
	"reflect"
	"strings"
)

const (
	PredicateIsNull               = " IS NULL"
	PredicateIsNotNull            = " IS NOT NULL"
	PredicateEqualTo              = "=?"
	PredicateNotEqualTo           = "<>?"
	PredicateGreaterThan          = ">?"
	PredicateGreaterThanOrEqualTo = ">=?"
	PredicateLessThan             = "<?"
	PredicateLessThanOrEqualTo    = "<=?"
	PredicateBetween              = " BETWEEN ? AND ?"
	PredicateLike                 = " LIKE ?"
)

// Literal returns a literal condition on a column. For example
//
//   Literal("age", " > 45")
//
// Be careful not to allow injection attacks: do not include a string from an external
// source in the column or predicate.
//
// This function is the basis for all the other predicates except In/InSlice.
func Literal(column, predicate string, value ...interface{}) Expression {
	if len(value) == 0 {
		return Condition{Column: column, Predicate: predicate}
	}
	return Condition{Column: column, Predicate: predicate, Args: value}
}

// Null returns an 'IS NULL' condition on a column.
func Null(column string) Expression {
	return Literal(column, PredicateIsNull)
}

// NotNull returns an 'IS NOT NULL' condition on a column. It's also possible to use Not(Null(...)).
func NotNull(column string) Expression {
	return Literal(column, PredicateIsNotNull)
}

// Eq returns an equality condition on a column.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func Eq(column string, value interface{}) Expression {
	return Literal(column, PredicateEqualTo, value)
}

// NotEq returns a not equal condition on a column.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func NotEq(column string, value interface{}) Expression {
	return Literal(column, PredicateNotEqualTo, value)
}

// Gt returns a greater than condition on a column.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func Gt(column string, value interface{}) Expression {
	return Literal(column, PredicateGreaterThan, value)
}

// GtEq returns a greater than or equal condition on a column.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func GtEq(column string, value interface{}) Expression {
	return Literal(column, PredicateGreaterThanOrEqualTo, value)
}

// Lt returns a less than condition on a column.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func Lt(column string, value interface{}) Expression {
	return Literal(column, PredicateLessThan, value)
}

// LtEq returns a less than or equal than condition on a column.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func LtEq(column string, value interface{}) Expression {
	return Literal(column, PredicateLessThanOrEqualTo, value)
}

// Between returns a between condition on a column.
//
// Two '?' placeholders are used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func Between(column string, a, b interface{}) Expression {
	return Literal(column, PredicateBetween, a, b)
}

// Like returns a pattern-matching condition on a column. Be careful: this can hurt performance.
//
// A '?' placeholder is used so it may be necessary to replace placeholders in the
// resulting query, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
func Like(column string, pattern string) Expression {
	return Literal(column, PredicateLike, pattern)
}

// In returns an 'IN' condition on a column.
// * If there a no values, this becomes a no-op.
// * If any value is nil, an 'IS NULL' expression is OR-ed with the 'IN' expression.
//
// Some '?' placeholders are used so it is necessary to replace placeholders in the
// resulting query according to SQL dialect, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
//
// Note that this does not reflection, unlike InSlice.
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

// InSlice returns an 'IN' condition on a column.
// * If arg is nil, this becomes a no-op.
// * arg is reflectively expanded as an array or slice to use all the contained values.
// * If any value is nil, an 'IS NULL' expression is OR-ed with the 'IN' expression.
//
// Some '?' placeholders are used so it is necessary to replace placeholders in the
// resulting query according to SQL dialect, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
//
// Note that this uses reflection, unlike In.
func InSlice(column string, arg interface{}) Expression {
	switch arg.(type) {
	case nil:
		return NoOp()
	}

	value := reflect.ValueOf(arg)

	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		// continue below
	default:
		panic("arg must be an array or slice")
	}

	hasNull := false
	var v []interface{}
	buf := &strings.Builder{}
	buf.WriteString(" IN (")

	for j := 0; j < value.Len(); j++ {
		vj := value.Index(j)
		switch vj.Kind() {
		case reflect.Ptr, reflect.Interface:
			if vj.IsNil() {
				hasNull = true
				continue
			}
		}

		if len(v) > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('?')
		v = append(v, vj.Interface())
	}

	buf.WriteByte(')')

	result := NoOp()
	if len(v) > 0 {
		result = Condition{column, buf.String(), v}
	}

	if hasNull {
		result = Or(result, Null(column))
	}

	return result
}

//-------------------------------------------------------------------------------------------------

const (
	and = " AND "
	or  = " OR "
)

// Not negates an expression.
func Not(el Expression) Expression {
	if el == nil {
		return NoOp()
	}
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
