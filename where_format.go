package where

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rickb777/where/v2/dialect"
	"github.com/rickb777/where/v2/quote"
)

func (not not) Format(option ...dialect.FormatOption) (string, []any) {
	prefix, inline := placeholderFromOptions(option)
	sql, args := not.doFormat(quoterFromOptions(option))
	return replacePlaceholders(sql, args, prefix, inline, 1)
}

func (not not) doFormat(quoter quote.Quoter) (string, []any) {
	sql, args := not.expression.doFormat(quoter)
	if sql == "" {
		return "", args
	}
	return "NOT (" + sql + ")", args
}

func (not not) String() string {
	sql, _ := not.Format(dialect.NoQuotes, dialect.Inline)
	return sql
}

//-------------------------------------------------------------------------------------------------

func (cl Condition) Format(option ...dialect.FormatOption) (string, []any) {
	prefix, inline := placeholderFromOptions(option)
	sql, args := cl.doFormat(quoterFromOptions(option))
	return replacePlaceholders(sql, args, prefix, inline, 1)
}

func (cl Condition) doFormat(quoter quote.Quoter) (string, []any) {
	buf := &strings.Builder{}
	quoter.QuoteW(buf, cl.Column)
	buf.WriteString(cl.Predicate)
	sql := buf.String()
	return sql, nilIfEmpty(cl.Args)
}

func (cl Condition) String() string {
	sql, _ := cl.Format(dialect.NoQuotes, dialect.Inline)
	return sql
}

//-------------------------------------------------------------------------------------------------

func (wh Clause) Format(option ...dialect.FormatOption) (string, []any) {
	prefix, inline := placeholderFromOptions(option)
	sql, args := wh.doFormat(quoterFromOptions(option))
	return replacePlaceholders(sql, args, prefix, inline, 1)
}

func (wh Clause) doFormat(quoter quote.Quoter) (string, []any) {
	if len(wh.wheres) == 0 {
		return "", nil
	}

	sqls := make([]string, 0, len(wh.wheres))
	var args []any

	for _, where := range wh.wheres {
		sql, a2 := where.doFormat(quoter)
		if len(sql) > 0 {
			sqls = append(sqls, "("+sql+")")
			args = append(args, a2...)
		}
	}

	sql := strings.Join(sqls, wh.conjunction)
	return sql, args
}

func (wh Clause) String() string {
	sql, _ := wh.Format(dialect.NoQuotes, dialect.Inline)
	return sql
}

//-------------------------------------------------------------------------------------------------

func placeholderFromOptions(option []dialect.FormatOption) (prefix string, inline bool) {
	for _, o := range option {
		switch o {
		case dialect.Dollar:
			prefix = "$"
		case dialect.AtP:
			prefix = "@p"
		case dialect.Inline:
			inline = true
		}
	}
	return prefix, inline
}

func quoterFromOptions(option []dialect.FormatOption) (quoter quote.Quoter) {
	quoter = quote.DefaultQuoter
	for _, o := range option {
		switch o {
		case dialect.NoQuotes:
			quoter = quote.None
		case dialect.ANSIQuotes:
			quoter = quote.ANSI
		case dialect.Backticks:
			quoter = quote.Backticks
		case dialect.SquareBrackets:
			quoter = quote.SquareBrackets
		}
	}
	return quoter
}

// replacePlaceholders replaces all "?" placeholders with numbered
// placeholders, using the given prefix.
// For PostgreSQL these will be "$1" and upward placeholders so the prefix should be "$" (Dollar).
// For SQL-Server there will be "@p1" and upward placeholders so the prefix should be "@p" (AtP).
// The count will start with 'from'.
func replacePlaceholders(sql string, args []any, prefix string, inline bool, from int) (string, []any) {
	if inline {
		return inlinePlaceholders(sql, args)
	}

	if prefix == "" {
		return sql, args
	}

	n := 0
	for _, r := range sql {
		if r == '?' {
			n++
		}
	}

	count := from
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

func inlinePlaceholders(sql string, args []any) (string, []any) {
	n := 0
	for _, r := range sql {
		if r == '?' {
			n++
		}
	}

	buf := &strings.Builder{}
	buf.Grow(len(sql) + len(sql)/2)

	for _, r := range sql {
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
	default:
		// TODO replace ' with ''
		return fmt.Sprintf(`'%v'`, v)
	}
}

func nilIfEmpty(args []any) []any {
	if len(args) > 0 {
		return args
	}
	return nil
}
