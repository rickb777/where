// Package dialect handles various dialect-specific ways of generating SQL.
package dialect

import (
	"strings"

	"github.com/rickb777/where/v2/quote"
)

// Dialect represents a dialect of SQL. All the defined dialects are non-zero.
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

// These are defaults used by each dialect; they can be altered before first use.
var (
	// SqliteQuoter uses ANSI double-quotes for Sqlite.
	// This can be modified, e.g. to None, before first use.
	SqliteQuoter = quote.ANSI

	// PostgresQuoter uses ANSI double-quotes for Postgres.
	// This can be modified, e.g. to None, before first use.
	PostgresQuoter = quote.ANSI

	// MySqlQuoter uses backticks for MySQL.
	// This can be modified, e.g. to None, before first use.
	MySqlQuoter = quote.Backticks

	// MSSqlQuoter uses square brackets for MS-SQL.
	// This can be modified, e.g. to None, before first use.
	MSSqlQuoter = quote.SquareBrackets
)

// Placeholder returns Query, Dollar or AtP.
func (d Dialect) Placeholder() FormatOption {
	switch d {
	case Postgres:
		return Dollar
	case SqlServer:
		return AtP
	}
	return Query
}

func (d Dialect) Quoter() quote.Quoter {
	switch d {
	case Mysql:
		return MySqlQuoter
	case Postgres:
		return PostgresQuoter
	case Sqlite:
		return SqliteQuoter
	case SqlServer:
		return MSSqlQuoter
	}
	return quote.DefaultQuoter
}

// String is the inverse of Pick.
func (d Dialect) String() string {
	switch d {
	case Sqlite:
		return "Sqlite"
	case Mysql:
		return "Mysql"
	case Postgres:
		return "Postgres"
	case SqlServer:
		return "SqlServer"
	}
	return ""
}

// Pick finds a dialect that matches by name, ignoring letter case.
// It matches:
//
//   - "sqlite", "sqlite3"
//   - "mysql"
//   - "postgres", "postgresql", "pgx"
//   - "sqlserver", "sql-server", "mssql"
//
// It returns 0 if not found.
func Pick(name string) Dialect {
	switch strings.ToLower(name) {
	case "sqlite", "sqlite3":
		return Sqlite
	case "mysql":
		return Mysql
	case "postgres", "postgresql", "pgx":
		return Postgres
	case "sqlserver", "sql-server", "mssql":
		return SqlServer
	}
	return undefined
}

// DefaultDialect is Sqlite, chosen as being probably the simplest. This can be
// altered before first use.
var DefaultDialect = Sqlite
