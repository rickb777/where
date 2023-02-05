package where_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/rickb777/where/v2"
	"github.com/rickb777/where/v2/dialect"
)

var queryConstraintAnsiQuoteCases = []struct {
	qc  *where.QueryConstraint
	exp string
}{
	{exp: "", qc: nil},
	{exp: ` ORDER BY "foo"`, qc: where.OrderBy("foo")},
	{exp: ` ORDER BY "foo", "bar"`, qc: where.OrderBy("foo", "bar")},
	{exp: ` ORDER BY "foo" NULLS FIRST`, qc: where.OrderBy("foo").Asc().NullsFirst()},
	{exp: ` ORDER BY "foo" DESC NULLS FIRST`, qc: where.OrderBy("foo").Desc().NullsFirst()},
	{exp: ` ORDER BY "foo", "bar", "baz"`, qc: where.OrderBy("foo", "bar", "baz").Asc()},
	{exp: ` ORDER BY "foo" DESC, "bar" DESC, "baz" DESC`, qc: where.OrderBy("foo", "bar", "baz").Desc()},
	{exp: ` ORDER BY "foo", "bar", "baz"`, qc: where.OrderBy("foo").OrderBy("bar").OrderBy("baz")},
	{exp: ` ORDER BY "foo" ASC, "bar" DESC, "baz" ASC`, qc: where.OrderBy("foo").OrderBy("bar").Desc().OrderBy("baz")},
	{exp: ` ORDER BY "foo" ASC, "bar" DESC`, qc: where.OrderBy("foo").Asc().OrderBy("bar").Desc()},
	{exp: ` ORDER BY "foo" ASC, "bar" DESC, "baz" ASC`, qc: where.OrderBy("foo").Asc().OrderBy("bar").Desc().OrderBy("baz").Asc()},
	{exp: ` ORDER BY "foo" DESC, "bar" ASC, "baz" DESC`, qc: where.OrderBy("foo").Desc().OrderBy("bar").Asc().OrderBy("baz").Desc()},
	{exp: ` ORDER BY "a", "b", "c", "d"`, qc: where.OrderBy("a", "b").Asc().OrderBy("c", "d").Asc()},
	{exp: ` ORDER BY "a" DESC, "b" DESC, "c" ASC, "d" ASC`, qc: where.OrderBy("a", "b").Desc().OrderBy("c", "d").Asc()},
	{exp: ` ORDER BY "a" ASC, "b" ASC, "c" DESC, "d" DESC`, qc: where.OrderBy("a", "b").Asc().OrderBy("c", "d").Desc()},

	{exp: ``, qc: where.Limit(0).NullsLast()},
	{exp: ` LIMIT 10`, qc: where.Limit(10)},
	{exp: ` OFFSET 20`, qc: where.Offset(20)},
	{exp: ` ORDER BY "foo", "bar" LIMIT 5`, qc: where.Limit(5).OrderBy("foo", "bar")},
	{exp: ` ORDER BY "foo" DESC NULLS LAST LIMIT 10 OFFSET 20`, qc: where.OrderBy("foo").Desc().Limit(10).Offset(20).NullsLast()},
}

var topConstraintCases = map[string]*where.QueryConstraint{
	`1`:          nil,
	`2`:          where.Limit(0),
	`3 TOP (10)`: where.Limit(10),
	`4`:          where.Offset(20),
	`5 TOP (5)`:  where.Limit(5).OrderBy("foo", "bar"),
	`6 TOP (10)`: where.OrderBy("foo").Desc().Limit(10).Offset(20),
}

func TestQueryConstraint_Format(t *testing.T) {
	g := NewGomegaWithT(t)

	for i, c := range queryConstraintAnsiQuoteCases {
		sql1 := c.qc.Format(dialect.Sqlite, dialect.ANSIQuotes)
		g.Expect(sql1).To(Equal(c.exp), "%d: %v", i, c)

		exp2 := strings.ReplaceAll(c.exp, `"`, ``)
		sql2 := c.qc.Format(dialect.Sqlite, dialect.NoQuotes)
		g.Expect(sql2).To(Equal(exp2), "%d: %v", i, c)
	}
}

func TestQueryConstraint_String(t *testing.T) {
	g := NewGomegaWithT(t)

	for exp, c := range queryConstraintAnsiQuoteCases {
		if c.qc != nil {
			sql := c.qc.Format(dialect.Sqlite, dialect.ANSIQuotes)
			g.Expect(sql).To(Equal(c.exp), exp)

			expected := strings.ReplaceAll(c.exp, `"`, ``)
			sql = c.qc.String()
			g.Expect(sql).To(Equal(expected), exp)
		}
	}
}

func BenchmarkQueryConstraint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range queryConstraintAnsiQuoteCases {
			_ = c.qc.Format(dialect.Sqlite)
		}
	}
}

func TestQueryConstraint_SqlServer(t *testing.T) {
	g := NewGomegaWithT(t)

	for exp, qc := range topConstraintCases {
		top := qc.FormatTOP(dialect.SqlServer)
		g.Expect(top).To(Equal(exp[1:]), exp)
	}

	for exp, qc := range topConstraintCases {
		if qc != nil {
			top := qc.FormatTOP(dialect.SqlServer)
			g.Expect(top).To(Equal(exp[1:]), exp)
		}
	}
}

func TestNilQueryConstraint_SqlServer(t *testing.T) {
	g := NewGomegaWithT(t)

	var qc where.QueryConstraint

	top := qc.FormatTOP(dialect.SqlServer)
	sql := qc.Format(dialect.SqlServer)

	g.Expect(top).To(Equal(""))
	g.Expect(sql).To(Equal(""))
}

func ExampleOrderBy() {
	// OrderBy understands that Asc and Desc apply to the preceding columns
	qc := where.OrderBy("foo", "bar").Desc().
		OrderBy("baz").Asc().
		Limit(10).
		Offset(20)

	// here we chose Sqlite, but Mysql nnd Postgres would give the same result
	s := qc.Format(dialect.Sqlite, dialect.NoQuotes)
	fmt.Println(s)

	// Sqlite doesn't use 'TOP' so it will be blank
	s = qc.FormatTOP(dialect.Sqlite)
	fmt.Println(s)

	// Output:  ORDER BY foo DESC, bar DESC, baz ASC LIMIT 10 OFFSET 20
	//
}

func ExampleQueryConstraint_NullsLast() {
	// OrderBy also includes a "NULLS LAST" phrase.
	qc := where.OrderBy("foo").NullsLast()

	// For Postgres, we're using double-quotes.
	s := qc.Format(dialect.Postgres, dialect.ANSIQuotes)
	fmt.Println(s)

	// Output:  ORDER BY "foo" NULLS LAST
}

func ExampleQueryConstraint_NullsFirst() {
	// OrderBy also includes a "NULLS LAST" phrase.
	qc := where.OrderBy("foo").NullsFirst()

	// For Postgres, we're using double-quotes.
	s := qc.Format(dialect.Postgres, dialect.ANSIQuotes)
	fmt.Println(s)

	// Output:  ORDER BY "foo" NULLS FIRST
}

func ExampleLimit() {
	// In this example, we can see how SqlServer needs different syntax
	// to Sqlite, Postgres, Mysql etc.
	qc := where.Limit(10).Offset(20)

	s1 := qc.Format(dialect.Sqlite)
	fmt.Println("SQlite:    ", s1)

	s2 := qc.FormatTOP(dialect.SqlServer)
	fmt.Println("SQL-Server:", s2)

	s3 := qc.Format(dialect.SqlServer)
	fmt.Println("SQL-Server:", s3)

	// Output: SQlite:      LIMIT 10 OFFSET 20
	// SQL-Server:  TOP (10)
	// SQL-Server:  OFFSET 20
}

func ExampleOffset() {
	// In this example, we start a query constraint using Offset	.
	qc := where.Offset(20)

	s := qc.Format(dialect.Postgres)
	fmt.Println(s)

	// Output: OFFSET 20
}
