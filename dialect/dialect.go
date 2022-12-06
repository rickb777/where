// Package dialect handles quote marks and SQL placeholders in various dialect-specific ways.
// Queries should be written using '?' query placeholders throughout, and then this package will translate to
// the form needed by the chosen dialect.
package dialect

import (
	"strconv"
	"strings"
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

func (d Dialect) Placeholder() PlaceholderStyle {
	switch d {
	case Postgres:
		return Dollar
	case SqlServer:
		return AtP
	}
	return Query
}

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
// * "sqlite", "sqlite3"
// * "mysql"
// * "postgres", "postgresql", "pgx"
// * "sqlserver", "sql-server", "mssql"
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

var DefaultDialect = Sqlite // chosen as being probably the simplest

//-------------------------------------------------------------------------------------------------

// PlaceholderStyle enumerates the different ways of including placeholders in SQL.
type PlaceholderStyle string

const (
	// Query is the '?' placeholder style.
	Query PlaceholderStyle = ""
	// Dollar numbered placeholders '$1', '$2' etc are used (e.g.) in PostgreSQL.
	Dollar PlaceholderStyle = "$"
	// At numbered placeholders '@p1', '@p2' etc are used in SQL-Server.
	AtP PlaceholderStyle = "@p"
)

// ReplacePlaceholders replaces all "?" placeholders with numbered
// placeholders, using the given prefix.
// For PostgreSQL these will be "$1" and upward placeholders so the prefix should be "$" (Dollar).
// For SQL-Server there will be "@p1" and upward placeholders so the prefix should be "@p" (AtP).
// The count will start with 'from', if provided, otherwise at one.
func ReplacePlaceholders(sql string, prefix PlaceholderStyle, from ...int) string {
	if prefix == "" {
		return sql
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
			buf.WriteString(string(prefix))
			buf.WriteString(strconv.Itoa(count))
			count++
		} else {
			buf.WriteRune(r)
		}
	}

	return buf.String()
}
