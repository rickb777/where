// Package dialect handles SQL placeholders in various dialect-specific ways. So queries should
// be written using '?' query placeholders throughout, and then this package will translate to
// the form needed by the chosen dialect.
package dialect

import (
	"strconv"
	"strings"

	"github.com/rickb777/where/quote"
)

// PlaceholderStyle enumerates the different ways of including placeholders in SQL.
type PlaceholderStyle int

const (
	// Queries is the '?' placeholder style and is assumed to be used prior to translation.
	Queries PlaceholderStyle = iota
	// Numbered placeholders '$1', '$2' etc are used (e.g.) in PostgreSQL.
	Numbered
	// Named placeholders ":name" are used (e.g.) in Oracle. NOT YET SUPPORTED
	// Named
)

//-------------------------------------------------------------------------------------------------

const (
	// SqliteIndex identifies SQLite
	SqliteIndex = iota
	// MysqlIndex identifies MySQL (also works for MariaDB)
	MysqlIndex
	// PostgresIndex identifies PostgreSQL
	PostgresIndex
	// SqlServerIndex identifies SqlServer (MS-SQL)
	SqlServerIndex
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
	case "postgres", "postgresql", "pgx":
		return Postgres, true
	case "sqlserver", "sql-server", "mssql":
		return SqlServer, true
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
	CaseInsensitive:  true,
}

// Postgres handles the PostgreSQL syntax.
var Postgres = Dialect{
	Ident:             PostgresIndex,
	PlaceholderStyle:  Numbered,
	PlaceholderPrefix: "$",
	Quoter:            quote.AnsiQuoter,
}

// SqlServer handles the T-SQL syntax.
// https://docs.microsoft.com/en-us/sql/t-sql/language-reference?view=sql-server-ver15
var SqlServer = Dialect{
	Ident:             SqlServerIndex,
	PlaceholderStyle:  Numbered,
	PlaceholderPrefix: "@p",
	Quoter:            quote.AnsiQuoter, // can also use square brackets but that's not supported here
}

var DefaultDialect = Sqlite // chosen as being probably the simplest

//-------------------------------------------------------------------------------------------------

// Dialect holds the settings to be used in SQL translation functions.
type Dialect struct {
	// Name is used for
	Ident int

	// PlaceholderStyle specifies the way of including placeholders in SQL.
	PlaceholderStyle PlaceholderStyle

	// PlaceholderPrefix specifies the string that marks a placeholder, when numbered
	PlaceholderPrefix string

	// Quoter determines the quote marks surrounding identifiers.
	Quoter quote.Quoter

	// CaseInsensitive is true when identifiers are not case-sensitive
	CaseInsensitive bool
}

// ReplacePlaceholders converts a string containing '?' placeholders to
// the form used by the dialect.
func (dialect Dialect) ReplacePlaceholders(sql string, names []string) string {
	switch dialect.PlaceholderStyle {
	case Numbered:
		return ReplacePlaceholdersWithNumbers(sql, dialect.PlaceholderPrefix)
	case Queries:
		return sql
	}
	panic(dialect.PlaceholderStyle)
}

// ReplacePlaceholdersWithNumbers replaces all "?" placeholders with numbered
// placeholders.
// For PostgreSQL these will be "$1" and upward placeholders.
// For SQL-Server, it inserts "@p1" and upward placeholders.
func ReplacePlaceholdersWithNumbers(sql, prefix string) string {
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
			buf.WriteString(prefix)
			buf.WriteString(strconv.Itoa(idx))
			idx++
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}
