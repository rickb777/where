package where_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/rickb777/where"
	"github.com/rickb777/where/dialect"
	"github.com/rickb777/where/quote"
)

var nameEqFredInt = where.Eq("name", "Fred")
var nameEqJohnInt = where.Eq("name", "John")
var ageLt10Int = where.Lt("age", 10)
var ageGt5Int = where.Gt("age", 5)

var (
	buildWhereClauseHappyCases = []struct {
		wh        where.Expression
		expWhere  string
		expHaving string
		expString string
		args      []interface{}
	}{
		{wh: where.NoOp()},

		{
			wh:        where.Condition{Column: "name", Predicate: " not nil", Args: nil},
			expWhere:  ` WHERE "name" not nil`,
			expHaving: ` HAVING "name" not nil`,
			expString: `"name" not nil`,
		},

		{
			wh:        where.Condition{Column: "p.name", Predicate: " not nil", Args: nil},
			expWhere:  ` WHERE "p"."name" not nil`,
			expHaving: ` HAVING "p"."name" not nil`,
			expString: `"p"."name" not nil`,
		},

		{
			wh:        where.Null("name"),
			expWhere:  ` WHERE "name" IS NULL`,
			expHaving: ` HAVING "name" IS NULL`,
			expString: `"name" IS NULL`,
		},

		{
			wh:        where.NotNull("name"),
			expWhere:  ` WHERE "name" IS NOT NULL`,
			expHaving: ` HAVING "name" IS NOT NULL`,
			expString: `"name" IS NOT NULL`,
		},

		{
			wh:        where.Condition{Column: "name", Predicate: " <>?", Args: []interface{}{"Boo"}},
			expWhere:  ` WHERE "name" <>?`,
			expHaving: ` HAVING "name" <>?`,
			expString: `"name" <>'Boo'`,
			args:      []interface{}{"Boo"},
		},

		{
			wh:        nameEqFredInt,
			expWhere:  ` WHERE "name"=?`,
			expHaving: ` HAVING "name"=?`,
			expString: `"name"='Fred'`,
			args:      []interface{}{"Fred"},
		},

		{
			wh:        where.Like("name", "F%"),
			expWhere:  ` WHERE "name" LIKE ?`,
			expHaving: ` HAVING "name" LIKE ?`,
			expString: `"name" LIKE 'F%'`,
			args:      []interface{}{"F%"},
		},

		{
			wh:        where.NoOp().And(nameEqFredInt),
			expWhere:  ` WHERE ("name"=?)`,
			expHaving: ` HAVING ("name"=?)`,
			expString: `("name"='Fred')`,
			args:      []interface{}{"Fred"},
		},

		{
			wh:        nameEqFredInt.And(where.NoOp()),
			expWhere:  ` WHERE ("name"=?)`,
			expHaving: ` HAVING ("name"=?)`,
			expString: `("name"='Fred')`,
			args:      []interface{}{"Fred"},
		},

		{
			wh:        nameEqFredInt.And(where.Gt("age", 10)),
			expWhere:  ` WHERE ("name"=?) AND ("age">?)`,
			expHaving: ` HAVING ("name"=?) AND ("age">?)`,
			expString: `("name"='Fred') AND ("age">10)`,
			args:      []interface{}{"Fred", 10},
		},

		{
			wh:        nameEqFredInt.Or(where.Gt("age", 10)),
			expWhere:  ` WHERE ("name"=?) OR ("age">?)`,
			expHaving: ` HAVING ("name"=?) OR ("age">?)`,
			expString: `("name"='Fred') OR ("age">10)`,
			args:      []interface{}{"Fred", 10},
		},

		{
			wh:        nameEqFredInt.And(ageGt5Int).And(where.Gt("weight", 15)),
			expWhere:  ` WHERE ("name"=?) AND ("age">?) AND ("weight">?)`,
			expHaving: ` HAVING ("name"=?) AND ("age">?) AND ("weight">?)`,
			expString: `("name"='Fred') AND ("age">5) AND ("weight">15)`,
			args:      []interface{}{"Fred", 5, 15},
		},

		{
			wh:        nameEqFredInt.Or(ageGt5Int).Or(where.Gt("weight", 15)),
			expWhere:  ` WHERE ("name"=?) OR ("age">?) OR ("weight">?)`,
			expHaving: ` HAVING ("name"=?) OR ("age">?) OR ("weight">?)`,
			expString: `("name"='Fred') OR ("age">5) OR ("weight">15)`,
			args:      []interface{}{"Fred", 5, 15},
		},

		{
			wh:        where.Between("age", 12, 18).Or(where.Gt("weight", 45)),
			expWhere:  ` WHERE ("age" BETWEEN ? AND ?) OR ("weight">?)`,
			expHaving: ` HAVING ("age" BETWEEN ? AND ?) OR ("weight">?)`,
			expString: `("age" BETWEEN 12 AND 18) OR ("weight">45)`,
			args:      []interface{}{12, 18, 45},
		},

		{
			wh:        where.GtEq("age", 10),
			expWhere:  ` WHERE "age">=?`,
			expHaving: ` HAVING "age">=?`,
			expString: `"age">=10`,
			args:      []interface{}{10},
		},

		{
			wh:        where.LtEq("age", 10),
			expWhere:  ` WHERE "age"<=?`,
			expHaving: ` HAVING "age"<=?`,
			expString: `"age"<=10`,
			args:      []interface{}{10},
		},

		{
			wh:        where.NotEq("age", 10),
			expWhere:  ` WHERE "age"<>?`,
			expHaving: ` HAVING "age"<>?`,
			expString: `"age"<>10`,
			args:      []interface{}{10},
		},

		{
			wh:        where.In("age", 10, 12, 14),
			expWhere:  ` WHERE "age" IN (?,?,?)`,
			expHaving: ` HAVING "age" IN (?,?,?)`,
			expString: `"age" IN (10,12,14)`,
			args:      []interface{}{10, 12, 14},
		},

		{ // 'In' without any vararg parameters
			wh: where.In("age"),
		},

		{ // 'In' with mixed value and nil vararg parameters
			wh:        where.In("age", 1, nil, 2, nil),
			expWhere:  ` WHERE ("age" IN (?,?)) OR ("age" IS NULL)`,
			expHaving: ` HAVING ("age" IN (?,?)) OR ("age" IS NULL)`,
			expString: `("age" IN (1,2)) OR ("age" IS NULL)`,
			args:      []interface{}{1, 2},
		},

		{ // 'In' with only a nil vararg parameter
			wh:        where.In("age", nil),
			expWhere:  ` WHERE "age" IS NULL`,
			expHaving: ` HAVING "age" IS NULL`,
			expString: `"age" IS NULL`,
		},

		{
			wh:        where.InSlice("ages", []int{10, 12, 14}),
			expWhere:  ` WHERE "ages" IN (?,?,?)`,
			expHaving: ` HAVING "ages" IN (?,?,?)`,
			expString: `"ages" IN (10,12,14)`,
			args:      []interface{}{10, 12, 14},
		},

		{ // 'InSlice' with mixed value and nil parameters
			wh:        where.InSlice("ages", []interface{}{1, nil, 2, nil}),
			expWhere:  ` WHERE ("ages" IN (?,?)) OR ("ages" IS NULL)`,
			expHaving: ` HAVING ("ages" IN (?,?)) OR ("ages" IS NULL)`,
			expString: `("ages" IN (1,2)) OR ("ages" IS NULL)`,
			args:      []interface{}{1, 2},
		},

		{ // 'InSlice' with only a nil parameter
			wh: where.InSlice("ages", nil),
		},

		{
			wh:        nameEqFredInt.Or(nameEqJohnInt),
			expWhere:  ` WHERE ("name"=?) OR ("name"=?)`,
			expHaving: ` HAVING ("name"=?) OR ("name"=?)`,
			expString: `("name"='Fred') OR ("name"='John')`,
			args:      []interface{}{"Fred", "John"},
		},

		{
			wh:        where.Not(nameEqFredInt),
			expWhere:  ` WHERE NOT ("name"=?)`,
			expHaving: ` HAVING NOT ("name"=?)`,
			expString: `NOT ("name"='Fred')`,
			args:      []interface{}{"Fred"},
		},

		{
			wh:        where.Not(nameEqFredInt.And(ageLt10Int)),
			expWhere:  ` WHERE NOT (("name"=?) AND ("age"<?))`,
			expHaving: ` HAVING NOT (("name"=?) AND ("age"<?))`,
			expString: `NOT (("name"='Fred') AND ("age"<10))`,
			args:      []interface{}{"Fred", 10},
		},

		{
			wh:        where.Not(nameEqFredInt.Or(ageLt10Int)),
			expWhere:  ` WHERE NOT (("name"=?) OR ("age"<?))`,
			expHaving: ` HAVING NOT (("name"=?) OR ("age"<?))`,
			expString: `NOT (("name"='Fred') OR ("age"<10))`,
			args:      []interface{}{"Fred", 10},
		},

		{
			wh:        where.Not(nameEqFredInt).And(ageLt10Int),
			expWhere:  ` WHERE (NOT ("name"=?)) AND ("age"<?)`,
			expHaving: ` HAVING (NOT ("name"=?)) AND ("age"<?)`,
			expString: `(NOT ("name"='Fred')) AND ("age"<10)`,
			args:      []interface{}{"Fred", 10},
		},

		{
			wh:        where.Not(nameEqFredInt).Or(ageLt10Int),
			expWhere:  ` WHERE (NOT ("name"=?)) OR ("age"<?)`,
			expHaving: ` HAVING (NOT ("name"=?)) OR ("age"<?)`,
			expString: `(NOT ("name"='Fred')) OR ("age"<10)`,
			args:      []interface{}{"Fred", 10},
		},

		{
			wh:        where.And(nameEqFredInt, nil, ageLt10Int),
			expWhere:  ` WHERE ("name"=?) AND ("age"<?)`,
			expHaving: ` HAVING ("name"=?) AND ("age"<?)`,
			expString: `("name"='Fred') AND ("age"<10)`,
			args:      []interface{}{"Fred", 10},
		},

		{
			wh:        where.Or(nameEqFredInt, nil, ageLt10Int),
			expWhere:  ` WHERE ("name"=?) OR ("age"<?)`,
			expHaving: ` HAVING ("name"=?) OR ("age"<?)`,
			expString: `("name"='Fred') OR ("age"<10)`,
			args:      []interface{}{"Fred", 10},
		},

		{
			wh:        where.And(nameEqFredInt).And(where.And(ageLt10Int)),
			expWhere:  ` WHERE ("name"=?) AND ("age"<?)`,
			expHaving: ` HAVING ("name"=?) AND ("age"<?)`,
			expString: `("name"='Fred') AND ("age"<10)`,
			args:      []interface{}{"Fred", 10},
		},

		{
			wh:        where.Eq("a", 1).And(where.Eq("b", 2)).And(where.And(where.Eq("c", 3), where.Eq("d", 4))),
			expWhere:  ` WHERE ("a"=?) AND ("b"=?) AND ("c"=?) AND ("d"=?)`,
			expHaving: ` HAVING ("a"=?) AND ("b"=?) AND ("c"=?) AND ("d"=?)`,
			expString: `("a"=1) AND ("b"=2) AND ("c"=3) AND ("d"=4)`,
			args:      []interface{}{1, 2, 3, 4},
		},

		{
			wh:        where.Eq("a", 1).Or(where.Eq("b", 2)).Or(where.Or(where.Eq("c", 3), where.Eq("d", 4))),
			expWhere:  ` WHERE ("a"=?) OR ("b"=?) OR ("c"=?) OR ("d"=?)`,
			expHaving: ` HAVING ("a"=?) OR ("b"=?) OR ("c"=?) OR ("d"=?)`,
			expString: `("a"=1) OR ("b"=2) OR ("c"=3) OR ("d"=4)`,
			args:      []interface{}{1, 2, 3, 4},
		},

		{
			wh:        where.And(nameEqFredInt.Or(nameEqJohnInt), ageLt10Int),
			expWhere:  ` WHERE (("name"=?) OR ("name"=?)) AND ("age"<?)`,
			expHaving: ` HAVING (("name"=?) OR ("name"=?)) AND ("age"<?)`,
			expString: `(("name"='Fred') OR ("name"='John')) AND ("age"<10)`,
			args:      []interface{}{"Fred", "John", 10},
		},

		{
			wh:        where.Or(nameEqFredInt, ageLt10Int.And(ageGt5Int)),
			expWhere:  ` WHERE ("name"=?) OR (("age"<?) AND ("age">?))`,
			expHaving: ` HAVING ("name"=?) OR (("age"<?) AND ("age">?))`,
			expString: `("name"='Fred') OR (("age"<10) AND ("age">5))`,
			args:      []interface{}{"Fred", 10, 5},
		},

		{
			wh:        where.Or(nameEqFredInt, nameEqJohnInt).And(ageGt5Int),
			expWhere:  ` WHERE (("name"=?) OR ("name"=?)) AND ("age">?)`,
			expHaving: ` HAVING (("name"=?) OR ("name"=?)) AND ("age">?)`,
			expString: `(("name"='Fred') OR ("name"='John')) AND ("age">5)`,
			args:      []interface{}{"Fred", "John", 5},
		},

		{
			wh:        where.Or(nameEqFredInt, nameEqJohnInt, where.And(ageGt5Int)),
			expWhere:  ` WHERE ("name"=?) OR ("name"=?) OR ("age">?)`,
			expHaving: ` HAVING ("name"=?) OR ("name"=?) OR ("age">?)`,
			expString: `("name"='Fred') OR ("name"='John') OR ("age">5)`,
			args:      []interface{}{"Fred", "John", 5},
		},

		{
			wh: where.Not(where.NoOp()),
		},

		{
			wh: where.Not(nil),
		},

		{
			wh: where.Or(nil).Or(where.NoOp()).And(where.NoOp()),
		},

		{
			wh: where.And(nil, where.Or(where.NoOp(), nil)),
		},
	}
)

func TestBuildWhereClause_happyCases(t *testing.T) {
	g := NewGomegaWithT(t)

	for i, c := range buildWhereClauseHappyCases {
		t.Logf("%d: %s", i, c.expWhere)

		sql, args := where.Where(c.wh)

		g.Expect(sql).To(Equal(c.expWhere))
		g.Expect(args).To(Equal(c.args))

		s := c.wh.String()

		g.Expect(s).To(Equal(c.expString))
	}
}

func TestBuildHavingClause_happyCases(t *testing.T) {
	g := NewGomegaWithT(t)

	for i, c := range buildWhereClauseHappyCases {
		t.Logf("%d: %s", i, c.expHaving)

		sql, args := where.Having(c.wh)

		g.Expect(sql).To(Equal(c.expHaving))
		g.Expect(args).To(Equal(c.args))

		s := c.wh.String()

		g.Expect(s).To(Equal(c.expString))
	}
}

func BenchmarkBuildWhereClause_happyCases_build(b *testing.B) {
	for _, c := range buildWhereClauseHappyCases {
		_, _ = where.Where(c.wh, quote.AnsiQuoter)
	}
}

func BenchmarkBuildWhereClause_happyCases_String(b *testing.B) {
	for _, c := range buildWhereClauseHappyCases {
		_ = c.wh.String()
	}
}

//-------------------------------------------------------------------------------------------------

func ExampleWhere() {
	// in this example, identifiers will be unquoted
	quote.DefaultQuoter = quote.NoQuoter

	// some simple expressions
	nameEqJohn := where.Eq("name", "John")
	nameEqPeter := where.Eq("name", "Peter")
	ageGt10 := where.Gt("age", 10)
	likes := where.In("likes", "cats", "dogs")

	// build a compound expression - this is a static expression
	// but it could be based on conditions instead
	wh := where.And(where.Or(nameEqJohn, nameEqPeter), ageGt10, likes)
	expr, args := where.Where(wh)

	// For Postgres, the placeholders have to be altered. It's necessary to do
	// this on the whole query if there might be other placeholders in it too.
	expr = dialect.ReplacePlaceholders(expr, dialect.Dollar)
	fmt.Println(expr)
	fmt.Println(args)

	// Output: WHERE ((name=$1) OR (name=$2)) AND (age>$3) AND (likes IN ($4,$5))
	// [John Peter 10 cats dogs]
}
