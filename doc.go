// Package where provides composable expressions for WHERE and HAVING clauses in SQL.
// These can range from the very simplest no-op to complex nested trees of 'AND' and 'OR'
// conditions.
//
// Also in this package are query constraints to provide 'ORDER BY', 'LIMIT' and 'OFFSET'
// clauses, along with 'TOP' for MS-SQL. These are similar to 'WHERE' clauses except literal values
// are used instead of parameter placeholders.
//
// Further support for SQL dialects and formatting options is provided in the 'dialect' sub-package.
//
// Queries should be written using '?' query placeholders throughout, and then these can be translated
// to the form needed by the chosen dialect: one of dialect.Query, dialect.Dollar, dialect.AtP or dialect.Inline.
//
// Also, support for quoted identifiers is provided in the 'quote' sub-package.
//   - quote.Quoter is an interface for a quoter.
//   - implementations include quote.ANSI, quote.Backticks, quote.SquareBrackets, and quote.None.
package where
