package where

import (
	"reflect"
	"strings"

	"github.com/rickb777/where/v2/predicate"
)

// Predicate returns a literal predicate. For example
//
//   - where.Predicate(`EXISTS (SELECT 1 FROM offers WHERE expiry_date = CURRENT_DATE)`)
//
// Column quoting won't apply.
//
// Be careful not to allow injection attacks: do not include a string from an external
// source in the predicate.
func Predicate(predicate string, value ...any) Expression {
	return Condition{Predicate: predicate, Args: value}
}

// Literal returns a literal condition on a column. For example
//
//   - where.Literal("age", " > 45")
//
// The column "age" will be quoted appropriately if a formatting option
// specifies this.
//
// Be careful not to allow injection attacks: do not include a string from an external
// source in the column or predicate.
//
// This function is the basis for most other predicates.
func Literal(column, predicate string, value ...any) Expression {
	return Condition{Column: column, Predicate: predicate, Args: value}
}

// Null returns an 'IS NULL' condition on a column.
func Null(column string) Expression {
	return Literal(column, predicate.IsNull)
}

// NotNull returns an 'IS NOT NULL' condition on a column. It's also possible to use Not(Null(...)).
func NotNull(column string) Expression {
	return Literal(column, predicate.IsNotNull)
}

// Eq returns an equality condition on a column.
func Eq(column string, value any) Expression {
	return Literal(column, predicate.EqualTo, value)
}

// NotEq returns a not equal condition on a column.
func NotEq(column string, value any) Expression {
	return Literal(column, predicate.NotEqualTo, value)
}

// Gt returns a greater than condition on a column.
func Gt(column string, value any) Expression {
	return Literal(column, predicate.GreaterThan, value)
}

// GtEq returns a greater than or equal condition on a column.
func GtEq(column string, value any) Expression {
	return Literal(column, predicate.GreaterThanOrEqualTo, value)
}

// Lt returns a less than condition on a column.
func Lt(column string, value any) Expression {
	return Literal(column, predicate.LessThan, value)
}

// LtEq returns a less than or equal than condition on a column.
func LtEq(column string, value any) Expression {
	return Literal(column, predicate.LessThanOrEqualTo, value)
}

// Between returns a between condition on a column.
func Between(column string, a, b any) Expression {
	return Literal(column, predicate.Between, a, b)
}

// Like returns a pattern-matching condition on a column. Be careful: this can hurt performance.
func Like(column string, pattern string) Expression {
	return Literal(column, predicate.Like, pattern)
}

// In returns an 'IN' condition on a column.
//   - If there are no values, this becomes a no-op.
//   - If any value is nil, an 'IS NULL' expression is OR-ed with the 'IN' expression.
//
// Note that this does not use reflection, unlike InSlice.
func In(column string, values ...any) Expression {
	if len(values) == 0 {
		return NoOp()
	}

	args := make([]any, 0, len(values))
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
//   - If arg is nil, this becomes a no-op.
//   - arg is reflectively expanded as an array or slice to use all the contained values.
//   - If any value is nil, an 'IS NULL' expression is OR-ed with the 'IN' expression.
//
// Some '?' placeholders are used so it is necessary to replace placeholders in the
// resulting query according to SQL dialect, e.g using 'dialect.ReplacePlaceholdersWithNumbers(query)'.
//
// Note that this uses reflection, unlike In.
func InSlice(column string, arg any) Expression {
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
	v := make([]any, 0, value.Len())
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
func Not(exp Expression) Expression {
	if exp == nil {
		return NoOp()
	}
	return not{expression: exp}
}

// And combines two conditions into a clause that requires they are both true.
func (exp not) And(other Expression) Expression {
	return Clause{wheres: []Expression{exp}, conjunction: and}.And(other)
}

// Or combines two conditions into a clause that requires either is true.
func (exp not) Or(other Expression) Expression {
	return Clause{wheres: []Expression{exp}, conjunction: or}.Or(other)
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
func (exp Condition) And(other Expression) Expression {
	return Clause{wheres: []Expression{exp}, conjunction: and}.And(other)
}

// Or combines two conditions into a clause that requires either is true.
func (exp Condition) Or(other Expression) Expression {
	return Clause{wheres: []Expression{exp}, conjunction: or}.Or(other)
}

//-------------------------------------------------------------------------------------------------

// And combines two clauses into a clause that requires they are both true.
// SQL implementation note: AND has higher precedence than OR.
func (exp Clause) conjoin(other Expression, conj string) Expression {
	cl, isClause := other.(Clause)
	if isClause {
		if len(exp.wheres) == 0 {
			return cl
		} else if len(cl.wheres) == 0 {
			return exp
		} else if exp.conjunction == conj && cl.conjunction == conj {
			return Clause{append(exp.wheres, cl.wheres...), conj}
		}
	} else {
		// blank case comes from NoOp
		if exp.conjunction == "" || exp.conjunction == conj {
			return Clause{append(exp.wheres, other), conj}
		}
	}
	return Clause{wheres: []Expression{exp, other}, conjunction: conj}
}

// And combines two clauses into a clause that requires they are both true.
// Parentheses will be inserted to preserve the calling order.
// SQL implementation note: AND has higher precedence than OR.
func (exp Clause) And(other Expression) Expression {
	return exp.conjoin(other, and)
}

// Or combines two clauses into a clause that requires either is true.
// Parentheses will be inserted to preserve the calling order.
// SQL implementation note: AND has higher precedence than OR.
func (exp Clause) Or(other Expression) Expression {
	return exp.conjoin(other, or)
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
