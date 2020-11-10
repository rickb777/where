// Package dialect handles quote marks and SQL placeholders in various dialect-specific ways.
// Queries should be written using '?' query placeholders throughout, and then this package will translate to
// the form needed by the chosen dialect.
//
// The XyzConfig variables are mutable so you can alter them if required. You can also alter DefaultDialect.
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

type Dialect int

const (
	undefined Dialect = iota
	// Sqlite identifies SQLite
	Sqlite
	// Mysql identifies MySQL (also works for MariaDB)
	Mysql
	// Postgres identifies PostgreSQL
	Postgres
	// SqlServer identifies SqlServer (MS-SQL)
	SqlServer
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
	return undefined, false
}

//-------------------------------------------------------------------------------------------------

func (d Dialect) Config() DialectConfig {
	switch d {
	case Sqlite:
		return SqliteConfig
	case Mysql:
		return MysqlConfig
	case Postgres:
		return PostgresConfig
	case SqlServer:
		return SqlServerConfig
	}
	return DialectConfig{}
}

// SqliteConfig handles the MySQL syntax.
var SqliteConfig = DialectConfig{
	Ident:            Sqlite,
	PlaceholderStyle: Queries,
	Quoter:           quote.AnsiQuoter,
}

// MysqlConfig handles the MySQL syntax.
var MysqlConfig = DialectConfig{
	Ident:            Mysql,
	PlaceholderStyle: Queries,
	Quoter:           quote.MySqlQuoter,
	CaseInsensitive:  true,
}

// PostgresConfig handles the PostgreSQL syntax.
var PostgresConfig = DialectConfig{
	Ident:             Postgres,
	PlaceholderStyle:  Numbered,
	PlaceholderPrefix: "$",
	Quoter:            quote.AnsiQuoter,
}

// SqlServerConfig handles the T-SQL syntax.
// https://docs.microsoft.com/en-us/sql/t-sql/language-reference?view=sql-server-ver15
var SqlServerConfig = DialectConfig{
	Ident:             SqlServer,
	PlaceholderStyle:  Numbered,
	PlaceholderPrefix: "@p",
	Quoter:            quote.AnsiQuoter, // can also use square brackets but that's not supported here
}

var DefaultDialect = Sqlite // chosen as being probably the simplest

//-------------------------------------------------------------------------------------------------

// DialectConfig holds the settings to be used in SQL translation functions.
type DialectConfig struct {
	// Name is used for
	Ident Dialect

	// PlaceholderStyle specifies the way of including placeholders in SQL.
	PlaceholderStyle PlaceholderStyle

	// PlaceholderPrefix specifies the string that marks a placeholder, when numbered
	PlaceholderPrefix string

	// Quoter determines the quote marks surrounding identifiers.
	Quoter quote.Quoter

	// CaseInsensitive is true when identifiers are not case-sensitive
	CaseInsensitive bool
}

//-------------------------------------------------------------------------------------------------

// ReplacePlaceholders converts a string containing '?' placeholders to
// the form used by the dialect.
func (d Dialect) ReplacePlaceholders(sql string, names []string) string {
	return d.Config().ReplacePlaceholders(sql, names)
}

// ReplacePlaceholders converts a string containing '?' placeholders to
// the form used by the dialect.
func (dc DialectConfig) ReplacePlaceholders(sql string, names []string) string {
	switch dc.PlaceholderStyle {
	case Numbered:
		return ReplacePlaceholdersWithNumbers(sql, dc.PlaceholderPrefix)
	case Queries:
		return sql
	}
	panic(dc.PlaceholderStyle)
}

// ReplacePlaceholdersWithNumbers replaces all "?" placeholders with numbered
// placeholders, using the given prefix.
// For PostgreSQL these will be "$1" and upward placeholders so the prefix should be "$".
// For SQL-Server there will be "@p1" and upward placeholders so the prefix should be "@p".
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
