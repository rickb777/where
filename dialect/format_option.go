// Package dialect handles various dialect-specific ways of generating SQL.
package dialect

// FormatOption provides controls for where-expression formatting.
type FormatOption int

// These options affect how placeholders are renderered.
const (
	// Query indicates placeholders using queries '?'. For Sqlite & MySql.
	Query FormatOption = iota

	// Dollar indicates placeholders using numbered $1, $2, ... format. For PostgreSQL.
	Dollar

	// AtP indicates placeholders using numbered @p1, @p2, ... format. For SQL-Server.
	AtP

	// Inline indicates that each placeholder is removed and its value is inlined.
	Inline
)

// These options affect how column names are quoted.
const (
	// NoQuotes indicates identifiers will not be enclosed in quote marks.
	NoQuotes FormatOption = iota + 10

	// ANSIQuotes indicates identifiers will be enclosed in double quote marks.
	ANSIQuotes

	// Backticks indicates identifiers will be enclosed in back-tick marks.
	Backticks

	// SquareBrackets indicates identifiers will be enclosed in square brackets.
	SquareBrackets
)
