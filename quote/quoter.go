// Package quote augments SQL strings by quoting identifiers according to three common
// variants: back-ticks used by MySQL, double-quotes used in ANSI SQL (PostgreSQL etc),
// or no quotes at all.
package quote

import (
	"io"
	"strings"
)

// Quoter wraps identifiers in quote marks. Compound identifiers (i.e. those with an alias prefix)
// are handled according to SQL grammar.
type Quoter interface {
	Quote(identifier string) string
	QuoteN(identifiers []string) []string
	QuoteW(w io.StringWriter, identifier string) (n int, err error)
}

const (
	ansiQuoter  = quoter(`"`)
	mySqlQuoter = quoter("`")
	noQuoter    = quoter("")
)

var (
	// AnsiQuoter wraps identifiers in double quotes.
	AnsiQuoter Quoter = ansiQuoter

	// MySqlQuoter wraps identifiers in back-ticks.
	MySqlQuoter Quoter = mySqlQuoter

	// NoQuoter leaves identifiers unquoted.
	NoQuoter Quoter = noQuoter

	// DefaultQuoter is used by default.
	//
	// Change this to affect the default setting for every SQL construction function.
	DefaultQuoter = AnsiQuoter
)

// NewQuoter gets a quoter using arbitrary quote marks.
func NewQuoter(mark string) Quoter {
	return quoter(mark)
}

// PickQuoter picks a quoter based on the names "ansi", "mysql" or "none".
// If none match, then nil is returned.
func PickQuoter(name string) Quoter {
	switch name {
	case "ansi":
		return AnsiQuoter
	case "mysql":
		return MySqlQuoter
	case "none":
		return NoQuoter
	default:
		return nil
	}
}

// quoter wraps identifiers in quote marks. Compound identifiers (i.e. those with an alias prefix)
// are handled according to SQL grammar.
type quoter string

// Quote renders an identifier within quote marks. If the identifier consists of both a
// prefix and a name, each part is quoted separately. Any i/o errors are silently dropped.
// For better performance, use QuoteW instead of Quote wherever possible.
func (q quoter) Quote(identifier string) string {
	if len(q) == 0 {
		return identifier
	}

	w := new(strings.Builder)
	w.Grow(len(identifier) + 2*len(q))
	q.QuoteW(w, identifier)
	return w.String()
}

// QuoteN quotes a list of identifiers using Quote.
func (q quoter) QuoteN(identifiers []string) []string {
	if len(q) == 0 {
		return identifiers
	}

	r := make([]string, 0, len(identifiers))
	for _, id := range identifiers {
		r = append(r, q.Quote(id))
	}
	return r
}

// QuoteW renders an identifier within quote marks. If the identifier consists of both a
// prefix and a name, each part is quoted separately.
func (q quoter) QuoteW(w io.StringWriter, identifier string) (n int, err error) {
	if len(q) == 0 {
		return w.WriteString(identifier)
	} else {
		elements := strings.Split(identifier, ".")
		return quoteW(w, q, q+"."+q, q, elements...)
	}
}

func quoteW(w io.StringWriter, before, sep, after quoter, elements ...string) (n int, err error) {
	if len(elements) == 0 {
		return 0, nil
	}

	var i int
	i, err = w.WriteString(string(before))
	n += i
	if err != nil {
		return n, err
	}

	// element[0] is always present
	i, err = w.WriteString(elements[0])
	n += i
	if err != nil {
		return n, err
	}

	// write the rest of the elements, preceding each with the separator
	for _, e := range elements[1:] {
		i, err = w.WriteString(string(sep))
		n += i
		if err != nil {
			return n, err
		}

		i, err = w.WriteString(e)
		n += i
		if err != nil {
			return n, err
		}
	}

	i, err = w.WriteString(string(after))
	n += i
	if err != nil {
		return n, err
	}

	return n, err
}
