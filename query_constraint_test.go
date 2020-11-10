package where_test

import (
	"fmt"
	"strconv"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/rickb777/where"
	"github.com/rickb777/where/dialect"
)

var queryConstraintCases = []struct {
	qc          where.QueryConstraint
	expPostgres string
}{
	{nil, ""},
	{where.OrderBy("foo"), ` ORDER BY "foo"`},
	{where.OrderBy("foo", "bar"), ` ORDER BY "foo", "bar"`},
	{where.OrderBy("foo").Asc(), ` ORDER BY "foo" ASC`},
	{where.OrderBy("foo").Desc(), ` ORDER BY "foo" DESC`},
	{where.OrderBy("foo", "bar").Desc(), ` ORDER BY "foo", "bar" DESC`},
	{where.OrderBy("foo").OrderBy("bar"), ` ORDER BY "foo", "bar"`},
	{where.OrderBy("foo").OrderBy("bar").Desc(), ` ORDER BY "foo", "bar" DESC`},
	{where.OrderBy("foo", "bar").Desc(), ` ORDER BY "foo", "bar" DESC`},
	{where.OrderBy("foo").Asc().OrderBy("bar").Desc(), ` ORDER BY "foo" ASC, "bar" DESC`},
	{where.OrderBy("foo").Desc().OrderBy("bar").Asc().OrderBy("baz").Desc(), ` ORDER BY "foo" DESC, "bar" ASC, "baz" DESC`},
	{where.Limit(0), ""},
	{where.Limit(10), " LIMIT 10"},
	{where.Offset(20), " OFFSET 20"},
	{where.Limit(5).OrderBy("foo", "bar"), ` ORDER BY "foo", "bar" LIMIT 5`},
	{where.OrderBy("foo").Desc().Limit(10).Offset(20), ` ORDER BY "foo" DESC LIMIT 10 OFFSET 20`},
}

func TestQueryConstraint(t *testing.T) {
	g := NewGomegaWithT(t)

	for i, c := range queryConstraintCases {
		var sql string

		if c.qc != nil {
			sql = where.Build(c.qc, dialect.Sqlite)
		}

		g.Expect(sql).To(Equal(c.expPostgres), strconv.Itoa(i))
	}
}

func BenchmarkQueryConstraint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range queryConstraintCases {
			if c.qc != nil {
				_ = where.Build(c.qc, dialect.Sqlite)
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

	// Output:  ORDER BY "foo", "bar" DESC, "baz" ASC LIMIT 10 OFFSET 20
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
