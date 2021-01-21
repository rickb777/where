package where_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/rickb777/where"
	"github.com/rickb777/where/dialect"
)

var queryConstraintCases = map[string]where.QueryConstraint{
	"01":                              nil,
	`02 ORDER BY "foo"`:               where.OrderBy("foo"),
	`03 ORDER BY "foo", "bar"`:        where.OrderBy("foo", "bar"),
	`04 ORDER BY "foo"`:               where.OrderBy("foo").Asc(),
	`05 ORDER BY "foo" DESC`:          where.OrderBy("foo").Desc(),
	`06 ORDER BY "foo", "bar", "baz"`: where.OrderBy("foo", "bar", "baz").Asc(),
	`07 ORDER BY "foo" DESC, "bar" DESC, "baz" DESC`:   where.OrderBy("foo", "bar", "baz").Desc(),
	`08 ORDER BY "foo", "bar", "baz"`:                  where.OrderBy("foo").OrderBy("bar").OrderBy("baz"),
	`09 ORDER BY "foo" ASC, "bar" DESC, "baz" ASC`:     where.OrderBy("foo").OrderBy("bar").Desc().OrderBy("baz"),
	`10 ORDER BY "foo" ASC, "bar" DESC`:                where.OrderBy("foo").Asc().OrderBy("bar").Desc(),
	`11 ORDER BY "foo" ASC, "bar" DESC, "baz" ASC`:     where.OrderBy("foo").Asc().OrderBy("bar").Desc().OrderBy("baz").Asc(),
	`12 ORDER BY "foo" DESC, "bar" ASC, "baz" DESC`:    where.OrderBy("foo").Desc().OrderBy("bar").Asc().OrderBy("baz").Desc(),
	`13 ORDER BY "a", "b", "c", "d"`:                   where.OrderBy("a", "b").Asc().OrderBy("c", "d").Asc(),
	`14 ORDER BY "a" DESC, "b" DESC, "c" ASC, "d" ASC`: where.OrderBy("a", "b").Desc().OrderBy("c", "d").Asc(),
	`15 ORDER BY "a" ASC, "b" ASC, "c" DESC, "d" DESC`: where.OrderBy("a", "b").Asc().OrderBy("c", "d").Desc(),

	`21`:                               where.Limit(0),
	`22 LIMIT 10`:                      where.Limit(10),
	`23 OFFSET 20`:                     where.Offset(20),
	`24 ORDER BY "foo", "bar" LIMIT 5`: where.Limit(5).OrderBy("foo", "bar"),
	`25 ORDER BY "foo" DESC LIMIT 10 OFFSET 20`: where.OrderBy("foo").Desc().Limit(10).Offset(20),
}

func TestQueryConstraint(t *testing.T) {
	g := NewGomegaWithT(t)

	for exp, c := range queryConstraintCases {
		var sql string

		if c != nil {
			sql = where.Build(c, dialect.Sqlite)
		}

		g.Expect(sql).To(Equal(exp[2:]), exp)
	}
}

func BenchmarkQueryConstraint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range queryConstraintCases {
			if c != nil {
				_ = where.Build(c, dialect.Sqlite)
			}
		}
	}
}

func TestQueryConstraint_SqlServer(t *testing.T) {
	g := NewGomegaWithT(t)

	qc := where.OrderBy("foo").Desc().OrderBy("bar").Asc().OrderBy("baz").Desc().Limit(10).Offset(5)

	top := where.BuildTop(qc, dialect.SqlServer)
	sql := where.Build(qc, dialect.SqlServer)

	g.Expect(top).To(Equal(" TOP (10)"))
	g.Expect(sql).To(Equal(` ORDER BY "foo" DESC, "bar" ASC, "baz" DESC OFFSET 5`))
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
