// Package where provides composable expressions for WHERE and HAVING clauses in SQL.
// These can range from the very simplest no-op to complex nested trees of 'and' and 'or'
// conditions.
//
// Also in this package are query constraints to provide 'ORDER BY', 'LIMIT' and 'OFFSET'
// clauses. These are similar to 'WHERE' clauses except literal values are used instead
// of parameter placeholders.
//
// Further support for parameter placeholders is provided in the 'dialect' sub-package.
// Also, support for quoted identifiers is provided in the 'quote' sub-package.
package where
