package dialect

import (
	"github.com/rickb777/where/quote"
	"strconv"
	"strings"
)

type PlaceholderStyle int

const (
	Queries PlaceholderStyle = iota
	Numbered
	//Named
)

//-------------------------------------------------------------------------------------------------

const (
	SqliteIndex = iota
	MysqlIndex
	PostgresIndex
)

//-------------------------------------------------------------------------------------------------

// PickDialect finds a dialect that matches by name, ignoring letter case.
// It returns false if not found.
func PickDialect(name string) (Dialect, bool) {
	switch strings.ToLower(name) {
	case "sqlite", "sqlite3":
		return Sqlite, true
	case "mysql":
		return Mysql, true
	case "postgres", "pgx":
		return Postgres, true
	}
	return Dialect{}, false
}

//-------------------------------------------------------------------------------------------------

// Sqlite handles the Sqlite syntax.
var Sqlite = Dialect{
	Ident:            SqliteIndex,
	PlaceholderStyle: Queries,
	Quoter:           quote.AnsiQuoter,
}

// Mysql handles the MySQL syntax.
var Mysql = Dialect{
	Ident:            MysqlIndex,
	PlaceholderStyle: Queries,
	Quoter:           quote.MySqlQuoter,
}

// Postgres handles the PostgreSQL syntax.
var Postgres = Dialect{
	Ident:            PostgresIndex,
	PlaceholderStyle: Numbered,
	Quoter:           quote.AnsiQuoter,
}

type Dialect struct {
	// Name is used for
	Ident int

	// HasNumberedPlaceholders is true for numbered placehoders (Postgresql),
	// or false for the default '?' placeholders.
	PlaceholderStyle PlaceholderStyle

	// named placeholders (e.g. for Oracle) are not yet supported

	// Quoter determines the quote marks surrounding identifiers.
	Quoter quote.Quoter
}

// ReplacePlaceholders converts a string containing '?' placeholders to
// the form used by the dialect.
func (dialect Dialect) ReplacePlaceholders(sql string, names []string) string {
	switch dialect.PlaceholderStyle {
	case Numbered:
		return ReplacePlaceholdersWithNumbers(sql)
	case Queries:
		return sql
	}
	panic(dialect.PlaceholderStyle)
}

// ReplacePlaceholdersWithNumbers replaces all '?' placeholders with '$1' etc numbered
// placeholders, as used by PostgreSQL etc.
func ReplacePlaceholdersWithNumbers(sql string) string {
	n := 0
	for _, r := range sql {
		if r == '?' {
			n++
		}
	}

	buf := &strings.Builder{}
	buf.Grow(len(sql) + n)
	idx := 1
	for _, r := range sql {
		if r == '?' {
			buf.WriteByte('$')
			buf.WriteString(strconv.Itoa(idx))
			idx++
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}
