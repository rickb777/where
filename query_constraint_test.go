package where_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/rickb777/where"
	"github.com/rickb777/where/dialect"
)

var queryConstraintCases = map[string]where.QueryConstraint{
	"01":                                               nil,
	`02 ORDER BY "foo"`:                                where.OrderBy("foo"),
	`03 ORDER BY "foo", "bar"`:                         where.OrderBy("foo", "bar"),
	`04 ORDER BY "foo" NULLS FIRST`:                    where.OrderBy("foo").Asc().NullsFirst(),
	`05 ORDER BY "foo" DESC NULLS FIRST`:               where.OrderBy("foo").Desc().NullsFirst(),
	`06 ORDER BY "foo", "bar", "baz"`:                  where.OrderBy("foo", "bar", "baz").Asc(),
	`07 ORDER BY "foo" DESC, "bar" DESC, "baz" DESC`:   where.OrderBy("foo", "bar", "baz").Desc(),
	`08 ORDER BY "foo", "bar", "baz"`:                  where.OrderBy("foo").OrderBy("bar").OrderBy("baz"),
	`09 ORDER BY "foo" ASC, "bar" DESC, "baz" ASC`:     where.OrderBy("foo").OrderBy("bar").Desc().OrderBy("baz"),
	`10 ORDER BY "foo" ASC, "bar" DESC`:                where.OrderBy("foo").Asc().OrderBy("bar").Desc(),
	`11 ORDER BY "foo" ASC, "bar" DESC, "baz" ASC`:     where.OrderBy("foo").Asc().OrderBy("bar").Desc().OrderBy("baz").Asc(),
	`12 ORDER BY "foo" DESC, "bar" ASC, "baz" DESC`:    where.OrderBy("foo").Desc().OrderBy("bar").Asc().OrderBy("baz").Desc(),
	`13 ORDER BY "a", "b", "c", "d"`:                   where.OrderBy("a", "b").Asc().OrderBy("c", "d").Asc(),
	`14 ORDER BY "a" DESC, "b" DESC, "c" ASC, "d" ASC`: where.OrderBy("a", "b").Desc().OrderBy("c", "d").Asc(),
	`15 ORDER BY "a" ASC, "b" ASC, "c" DESC, "d" DESC`: where.OrderBy("a", "b").Asc().OrderBy("c", "d").Desc(),

	`21`:                               where.Limit(0).NullsLast(),
	`22 LIMIT 10`:                      where.Limit(10),
	`23 OFFSET 20`:                     where.Offset(20),
	`24 ORDER BY "foo", "bar" LIMIT 5`: where.Limit(5).OrderBy("foo", "bar"),
	`25 ORDER BY "foo" DESC NULLS LAST LIMIT 10 OFFSET 20`: where.OrderBy("foo").Desc().Limit(10).Offset(20).NullsLast(),
}

var topConstraintCases = map[string]where.QueryConstraint{
	`1`:          nil,
	`2`:          where.Limit(0),
	`3 TOP (10)`: where.Limit(10),
	`4`:          where.Offset(20),
	`5 TOP (5)`:  where.Limit(5).OrderBy("foo", "bar"),
	`6 TOP (10)`: where.OrderBy("foo").Desc().Limit(10).Offset(20),
}

func TestQueryConstraint1(t *testing.T) {
	g := NewGomegaWithT(t)

	for exp, c := range queryConstraintCases {
		sql := where.Build(c, dialect.Sqlite)
		g.Expect(sql).To(Equal(exp[2:]), exp)
	}
}

func TestQueryConstraint2(t *testing.T) {
	g := NewGomegaWithT(t)

	for exp, c := range queryConstraintCases {
		if c != nil {
			sql := c.Build(dialect.Sqlite)
			g.Expect(sql).To(Equal(exp[2:]), exp)

			sql = c.String()
			g.Expect(sql).To(Equal(exp[2:]), exp)
		}
	}
}

func BenchmarkQueryConstraint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range queryConstraintCases {
			_ = where.Build(c, dialect.Sqlite)
		}
	}
}

func TestQueryConstraint_SqlServer(t *testing.T) {
	g := NewGomegaWithT(t)

	for exp, qc := range topConstraintCases {
		top := where.BuildTop(qc, dialect.SqlServer)
		g.Expect(top).To(Equal(exp[1:]), exp)
	}

	for exp, qc := range topConstraintCases {
		if qc != nil {
			top := qc.BuildTop(dialect.SqlServer)
			g.Expect(top).To(Equal(exp[1:]), exp)
		}
	}
}

func TestNilQueryConstraint_SqlServer(t *testing.T) {
	g := NewGomegaWithT(t)

	var qc where.QueryConstraint

	top := where.BuildTop(qc, dialect.SqlServer)
	sql := where.Build(qc, dialect.SqlServer)

	g.Expect(top).To(Equal(""))
	g.Expect(sql).To(Equal(""))
}

func ExampleOrderBy() {
	qc := where.OrderBy("foo", "bar").Desc().OrderBy("baz").Asc().Limit(10).Offset(20)

	s := qc.Build(dialect.Sqlite)
	fmt.Println(s)

	// Output:  ORDER BY "foo" DESC, "bar" DESC, "baz" ASC LIMIT 10 OFFSET 20
}

func ExampleLimit() {
	qc := where.Limit(10).Offset(20)

	s1 := qc.Build(dialect.Sqlite)
	fmt.Println("SQlite:    ", s1)

	s2 := qc.BuildTop(dialect.SqlServer)
	fmt.Println("SQL-Server:", s2)

	s3 := qc.Build(dialect.SqlServer)
	fmt.Println("SQL-Server:", s3)

	// Output: SQlite:      LIMIT 10 OFFSET 20
	// SQL-Server:  TOP (10)
	// SQL-Server:  OFFSET 20
}
