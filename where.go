package where

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rickb777/where/quote"
)

const (
	whereAdverb = " WHERE "
	havingVerb  = " HAVING "
)

// Expression is an element in a WHERE clause. Expressions consist of simple conditions or
// more complex clauses of multiple conditions.
type Expression interface {
	fmt.Stringer
	And(Expression) Expression
	Or(Expression) Expression
	build(q quote.Quoter) (string, []interface{})
}

// Where constructs the sql clause beginning "WHERE ...". It will contain '?' style placeholders;
// these need to be passed through the relevant quote ReplacePlaceholders processing.
// A quoter may optionally be supplied, otherwise the Default Quoter is used.
func Where(wh Expression, q ...quote.Quoter) (string, []interface{}) {
	return build(whereAdverb, wh, q...)
}

// Having constructs the sql clause beginning "HAVING ...". It will contain '?' style placeholders;
// these need to be passed through the relevant quote ReplacePlaceholders processing.
// A quoter may optionally be supplied, otherwise the Default Quoter is used.
func Having(wh Expression, q ...quote.Quoter) (string, []interface{}) {
	return build(havingVerb, wh, q...)
}

// build constructs the sql clause beginning with some verb/adverb. It will contain '?' style placeholders;
// these need to be passed through the relevant quote ReplacePlaceholders processing.
func build(adverb string, wh Expression, quoter ...quote.Quoter) (string, []interface{}) {
	if wh == nil {
		return "", nil
	}
	q := pickQuoter(quoter)
	sql, args := wh.build(q)
	if sql == "" {
		return "", nil
	}
	return adverb + sql, args
}

//-------------------------------------------------------------------------------------------------

type not struct {
	expression Expression
}

func (not not) build(q quote.Quoter) (string, []interface{}) {
	sql, args := not.expression.build(q)
	if sql == "" {
		return "", args
	}
	return "NOT (" + sql + ")", args
}

func (not not) String() string {
	sql, args := not.build(quote.DefaultQuoter)
	return insertLiteralValues(sql, args)
}

//-------------------------------------------------------------------------------------------------

// Condition is a simple condition such as an equality test. For convenience, use the
// factory functions 'Eq', 'GtEq' etc.
//
// This can also be constructed directly, which will be useful for non-portable
// cases, such as Postgresql 'SIMILAR TO'
//
//	expr := where.Condition{Column: "name", Predicate: " SIMILAR TO", Args: []interface{}{pattern}}
//
// Also for literal values (taking care to protect against injection attacks)
//
//	expr := where.Condition{Column: "age", Predicate: " = 47", Args: nil}
//
// See Literal.
type Condition struct {
	Column, Predicate string
	Args              []interface{}
}

func (cl Condition) build(q quote.Quoter) (string, []interface{}) {
	sql := q.Quote(cl.Column) + cl.Predicate
	if len(cl.Args) > 0 {
		return sql, cl.Args
	}
	return sql, nil
}

func (cl Condition) String() string {
	sql, args := cl.build(quote.DefaultQuoter)
	return insertLiteralValues(sql, args)
}

//-------------------------------------------------------------------------------------------------

// Clause is a compound expression.
type Clause struct {
	wheres      []Expression
	conjunction string
}

func (wh Clause) build(q quote.Quoter) (string, []interface{}) {
	if len(wh.wheres) == 0 {
		return "", nil
	}

	sqls := make([]string, 0, len(wh.wheres))
	var args []interface{}

	for _, where := range wh.wheres {
		sql, a2 := where.build(q)
		if len(sql) > 0 {
			sqls = append(sqls, "("+sql+")")
			args = append(args, a2...)
		}
	}

	sql := strings.Join(sqls, wh.conjunction)
	return sql, args
}

func (wh Clause) String() string {
	sql, args := wh.build(quote.DefaultQuoter)
	return insertLiteralValues(sql, args)
}

//-------------------------------------------------------------------------------------------------

func insertLiteralValues(sql string, args []interface{}) string {
	// create a buffer with approximately enough space
	buf := new(strings.Builder)
	buf.Grow(len(sql) + 6*len(args))

	idx := 0
	for _, r := range sql {
		if r == '?' && idx < len(args) {
			v := args[idx]
			t := reflect.TypeOf(v)
			switch t.Kind() {
			case reflect.Bool,
				reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64,
				reflect.Uint,
				reflect.Uint8,
				reflect.Uint16,
				reflect.Uint32,
				reflect.Uint64,
				reflect.Uintptr,
				reflect.Float32,
				reflect.Float64:
				buf.WriteString(fmt.Sprintf(`%v`, v))
			default:
				buf.WriteString(fmt.Sprintf(`'%v'`, v))
			}
			idx++
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}
