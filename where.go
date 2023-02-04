// Package where provides composable expressions for WHERE and HAVING clauses in SQL.
// These can range from the very simplest no-op to complex nested trees of 'AND' and 'OR'
// conditions.
package where

import (
	"github.com/rickb777/where/v2/dialect"
	"github.com/rickb777/where/v2/quote"
)

// Expression is an element in a WHERE clause. Expressions consist of simple conditions or
// more complex clauses of multiple conditions.
type Expression interface {
	// String prints the expression with inlined values inserted instead of placeholders.
	// Column names are not quoted.
	String() string

	// Format formats the (nested) expression as a string containing placeholders etc.
	// It doesn't include the WHERE or HAVING conjunction word.
	Format(option ...dialect.FormatOption) (string, []interface{})
	// doFormat formats the (nested) expression as a string containing placeholders etc.
	doFormat(quoter quote.Quoter) (string, []interface{})

	// And concatenates this expression with another such that both must evaluate true.
	And(Expression) Expression
	// Or concatenates this expression with another such that either must evaluate true.
	Or(Expression) Expression
}

const (
	whereConjunction  = " WHERE "
	havingConjunction = " HAVING "
)

// Where constructs the SQL clause beginning "WHERE ...".
// If the expression is empty or nil, the returned string will be blank.
// Optional parameters may be supplied. Otherwise, by default, quote.DefaultQuoter is used
// and the result will contain '?' style placeholders.
func Where(wh Expression, option ...dialect.FormatOption) (string, []interface{}) {
	return format(whereConjunction, wh, option...)
}

// Having constructs the SQL clause beginning "HAVING ...".
// If the expression is empty or nil, the returned string will be blank.
// Optional parameters may be supplied. Otherwise, by default, quote.DefaultQuoter is used
// and the result will contain '?' style placeholders.
func Having(wh Expression, option ...dialect.FormatOption) (string, []interface{}) {
	return format(havingConjunction, wh, option...)
}

// format constructs the sql clause beginning with some verb/adverb.
func format(conjunction string, wh Expression, option ...dialect.FormatOption) (string, []interface{}) {
	if wh == nil {
		return "", nil
	}

	expression, args := wh.Format(option...)
	if expression == "" {
		return "", nil
	}

	return conjunction + expression, args
}

//-------------------------------------------------------------------------------------------------

type not struct {
	expression Expression
}

//-------------------------------------------------------------------------------------------------

// Condition is a simple condition such as an equality test. For convenience, use the
// factory functions 'Eq', 'GtEq' etc.
//
// This can also be constructed directly, which will be useful for non-portable
// cases, such as Postgresql 'SIMILAR TO'
//
//	expr := where.Condition{Column: "name", Predicate: " SIMILAR TO ?", Args: []interface{}{pattern}}
//
// Also for literal values (taking care to protect against injection attacks)
//
//	expr := where.Condition{Column: "age", Predicate: " = 47", Args: nil}
//
// Column can be left blank; this allows the predicate to be a sub-query such as EXISTS(...).
//
// See Literal.
type Condition struct {
	Column, Predicate string
	Args              []interface{}
}

//-------------------------------------------------------------------------------------------------

// Clause is a compound expression.
type Clause struct {
	wheres      []Expression
	conjunction string
}
