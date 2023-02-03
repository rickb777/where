// Package quote augments SQL strings by quoting identifiers according to four common
// variants:
//   - back-ticks used by MySQL,
//   - double-quotes used in ANSI SQL (PostgreSQL etc),
//   - square brackets used by SQLServer, or
//   - no quotes at all.
//
// For prefixed identifiers containing a dot ('.'), the quote marks are applied separately
// to the prefix(es) and the identifier itself.
//
// Quoting is only applied where a column name is clear (as in the 'where' package). There
// isn't a general syntax analyser that somehow fixes up arbitrary SQL.
package quote

import (
	"io"
	"regexp"
	"strings"
)

// Quoter wraps identifiers in quote marks. Compound identifiers, i.e. those with an alias
// prefix such as "excluded"."created_at", are handled according to SQL grammar.
type Quoter interface {
	// Quote renders an identifier within quote marks. If the identifier consists of both a
	// prefix and a name, each part is quoted separately. Any i/o errors are silently dropped.
	Quote(identifier string) string

	// QuoteW renders an identifier within quote marks. If the identifier consists of one or more
	// prefixes and a name, each part is quoted separately. This may give better performance
	// than Quote.
	QuoteW(w io.StringWriter, identifier string)
}

const none = noQuoter("")

var (
	// None leaves identifiers unchanged.
	None = none

	// ANSI wraps identifiers in double-quote marks. For PostgreSQL etc.
	ANSI = quoter{before: `"`, between: `"."`, after: `"`}

	// Backticks wraps identifies in back-ticks. For MySql etc.
	Backticks = quoter{before: "`", between: "`.`", after: "`"}

	// SquareBrackets wraps identifies in '[' and ']'. For MS SQL/SQL-Server.
	SquareBrackets = quoter{before: "[", between: "].[", after: "]"}
)

var (
	// DefaultQuoter does not change identifiers and is used by default.
	// Change this to affect the default setting for every SQL construction function.
	DefaultQuoter = none
)

// Pick picks a quoter based on the names "ansi", "backtick" (aliases "backticks") or "none",
// ignoring case. Other options are also permitted: "sqlite", "sqlite3", "postgres",
// "mysql", "mssql", "ms-sql", "sql-server". The default is none.
func Pick(name string) Quoter {
	switch strings.ToLower(name) {
	case "ansi", "postgres", "sqlite", "sqlite3":
		return ANSI
	case "backtick", "backticks", "mysql":
		return Backticks
	case "mssql", "ms-sql", "sql-server":
		return SquareBrackets
	default:
		return none
	}
}

//-------------------------------------------------------------------------------------------------

// quoter wraps identifiers in quote marks. Compound identifiers (i.e. those with alias prefixes)
// are handled according to SQL grammar.
type quoter struct {
	before, between, after string
}

var validIdentifier = regexp.MustCompile(`^\pL[\pL\pN_]*$`)

func (q quoter) Quote(identifier string) string {
	if len(identifier) == 0 {
		return ""
	}

	w := new(strings.Builder)
	w.Grow(len(identifier) + 2*(len(q.before)+len(q.after)))
	q.QuoteW(w, identifier)
	return w.String()
}

func (q quoter) QuoteW(w io.StringWriter, identifier string) {
	if len(identifier) > 0 {
		names := strings.Split(identifier, ".")

		for _, name := range names {
			// if any name is invalid, we leave the entire string unaltered
			if !validIdentifier.MatchString(name) {
				_, _ = w.WriteString(identifier)
				return
			}
		}

		quoteW(w, q.before, q.between, q.after, names...)
	}
}

func quoteW(w io.StringWriter, before, sep, after string, names ...string) {
	_, _ = w.WriteString(before)
	_, _ = w.WriteString(names[0])

	// write the rest of the names, preceding each with the separator
	for _, e := range names[1:] {
		_, _ = w.WriteString(sep)
		_, _ = w.WriteString(e)
	}

	_, _ = w.WriteString(after)
}

//-------------------------------------------------------------------------------------------------

type noQuoter string

func (noQuoter) Quote(identifier string) string              { return identifier }
func (noQuoter) QuoteW(w io.StringWriter, identifier string) { _, _ = w.WriteString(identifier) }
