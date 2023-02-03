package where

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rickb777/where/v2/dialect"
	"github.com/rickb777/where/v2/quote"
)

func (exp not) Format(option ...dialect.FormatOption) (string, []any) {
	placeholderOption := formatOptions(option).Placeholder()
	quoter := quoterFromOptions(formatOptions(option).Quoter())
	sql, args := exp.doFormat(quoter)
	return replacePlaceholders(sql, args, placeholderOption, 1)
}

func (exp not) doFormat(quoter quote.Quoter) (string, []any) {
	sql, args := exp.expression.doFormat(quoter)
	if sql == "" {
		return "", args
	}
	return "NOT (" + sql + ")", args
}

func (exp not) String() string {
	sql, _ := exp.Format(dialect.NoQuotes, dialect.Inline)
	return sql
}

//-------------------------------------------------------------------------------------------------

func (exp Condition) Format(option ...dialect.FormatOption) (string, []any) {
	placeholderOption := formatOptions(option).Placeholder()
	quoter := quoterFromOptions(formatOptions(option).Quoter())
	sql, args := exp.doFormat(quoter)
	return replacePlaceholders(sql, args, placeholderOption, 1)
}

func (exp Condition) doFormat(quoter quote.Quoter) (string, []any) {
	buf := &strings.Builder{}
	quoter.QuoteW(buf, exp.Column)
	buf.WriteString(exp.Predicate)
	sql := buf.String()
	return sql, nilIfEmpty(exp.Args)
}

func (exp Condition) String() string {
	sql, _ := exp.Format(dialect.NoQuotes, dialect.Inline)
	return sql
}

//-------------------------------------------------------------------------------------------------

func (exp Clause) Format(option ...dialect.FormatOption) (string, []any) {
	placeholderOption := formatOptions(option).Placeholder()
	quoter := quoterFromOptions(formatOptions(option).Quoter())
	sql, args := exp.doFormat(quoter)
	return replacePlaceholders(sql, args, placeholderOption, 1)
}

func (exp Clause) doFormat(quoter quote.Quoter) (string, []any) {
	if len(exp.wheres) == 0 {
		return "", nil
	}

	sqls := make([]string, 0, len(exp.wheres))
	var args []any

	for _, where := range exp.wheres {
		sql, a2 := where.doFormat(quoter)
		if len(sql) > 0 {
			sqls = append(sqls, "("+sql+")")
			args = append(args, a2...)
		}
	}

	sql := strings.Join(sqls, exp.conjunction)
	return sql, args
}

func (exp Clause) String() string {
	sql, _ := exp.Format(dialect.NoQuotes, dialect.Inline)
	return sql
}

//-------------------------------------------------------------------------------------------------

func prefixFromOption(option dialect.FormatOption) (prefix string) {
	switch option {
	case dialect.Dollar:
		prefix = "$"
	case dialect.AtP:
		prefix = "@p"
	}
	return prefix
}

func quoterFromOptions(option dialect.FormatOption) (quoter quote.Quoter) {
	quoter = quote.DefaultQuoter
	switch option {
	case dialect.NoQuotes:
		quoter = quote.None
	case dialect.ANSIQuotes:
		quoter = quote.ANSI
	case dialect.Backticks:
		quoter = quote.Backticks
	case dialect.SquareBrackets:
		quoter = quote.SquareBrackets
	}
	return quoter
}

func replacePlaceholders(sql string, args []any, opt dialect.FormatOption, from int) (string, []any) {
	if opt == dialect.Inline {
		return InlinePlaceholders(sql, args)
	}

	return ReplacePlaceholders(sql, args, opt, from)
}

// ReplacePlaceholders replaces all "?" placeholders with numbered placeholders, using the given dialect option.
//   - For PostgreSQL these will be "$1" and upward placeholders so the prefix should be "$" (Dollar).
//   - For SQL-Server there will be "@p1" and upward placeholders so the prefix should be "@p" (AtP).
//
// The count will start with 'from', or from 1.
func ReplacePlaceholders(sql string, args []any, opt dialect.FormatOption, from ...int) (string, []any) {
	prefix := prefixFromOption(opt)
	if prefix == "" {
		return sql, args
	}

	n := 0
	for _, r := range sql {
		if r == '?' {
			n++
		}
	}

	count := 1
	if len(from) > 0 {
		count = from[0]
	}

	buf := &strings.Builder{}
	buf.Grow(len(sql) + n*(len(prefix)+2))

	for _, r := range sql {
		if r == '?' {
			buf.WriteString(prefix)
			buf.WriteString(strconv.Itoa(count))
			count++
		} else {
			buf.WriteRune(r)
		}
	}

	return buf.String(), nilIfEmpty(args)
}

// InlinePlaceholders replaces every '?' placeholder with the corresponding argument value.
// Number and boolean arguments are inserted verbatim. Everything else is inserted
// in string syntax, i.e. enclosed in single quote marks.
func InlinePlaceholders(query string, args []any) (string, []any) {
	n := 0
	for _, r := range query {
		if r == '?' {
			n++
		}
	}

	buf := &strings.Builder{}
	buf.Grow(len(query) + len(query)/2) // heuristic

	for _, r := range query {
		if r == '?' && len(args) > 0 {
			buf.WriteString(literalValue(args[0]))
			args = args[1:]
		} else {
			buf.WriteRune(r)
		}
	}

	return buf.String(), nilIfEmpty(args)
}

func literalValue(v any) string {
	switch x := v.(type) {
	case bool:
		return strconv.FormatBool(x)
	case int:
		return strconv.FormatInt(int64(x), 10)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	}

	s := fmt.Sprintf(`%v`, v)
	s = strings.ReplaceAll(s, "'", "''")
	return "'" + s + "'"
}

func nilIfEmpty(args []any) []any {
	if len(args) > 0 {
		return args
	}
	return nil
}

//-------------------------------------------------------------------------------------------------

type formatOptions []dialect.FormatOption

func (opts formatOptions) Placeholder() dialect.FormatOption {
	for _, o := range opts {
		if o <= dialect.Inline {
			return o
		}
	}
	return 0
}

func (opts formatOptions) Quoter() dialect.FormatOption {
	for _, o := range opts {
		if o >= dialect.NoQuotes {
			return o
		}
	}
	return 0
}
