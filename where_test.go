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

var buildWhereClauseHappyCases = []struct {
	wh        where.Expression
	expSQL    string
	expString string
	args      []interface{}
}{
	{where.NoOp(), "", "", nil},

	{
		where.Condition{Column: "name", Predicate: " not nil", Args: nil},
		` WHERE "name" not nil`,
		`"name" not nil`,
		nil,
	},

	{
		where.Condition{Column: "p.name", Predicate: " not nil", Args: nil},
		` WHERE "p"."name" not nil`,
		`"p"."name" not nil`,
		nil,
	},

	{
		where.Null("name"),
		` WHERE "name" IS NULL`,
		`"name" IS NULL`,
		nil,
	},

	{
		where.NotNull("name"),
		` WHERE "name" IS NOT NULL`,
		`"name" IS NOT NULL`,
		nil,
	},

	{
		where.Condition{Column: "name", Predicate: " <>?", Args: []interface{}{"Boo"}},
		` WHERE "name" <>?`,
		`"name" <>'Boo'`,
		[]interface{}{"Boo"},
	},

	{
		nameEqFredInt,
		` WHERE "name"=?`,
		`"name"='Fred'`,
		[]interface{}{"Fred"},
	},

	{
		where.Like("name", "F%"),
		` WHERE "name" LIKE ?`,
		`"name" LIKE 'F%'`,
		[]interface{}{"F%"},
	},

	{
		where.NoOp().And(nameEqFredInt),
		` WHERE ("name"=?)`,
		`("name"='Fred')`,
		[]interface{}{"Fred"},
	},

	{
		nameEqFredInt.And(where.Gt("age", 10)),
		` WHERE ("name"=?) AND ("age">?)`,
		`("name"='Fred') AND ("age">10)`,
		[]interface{}{"Fred", 10},
	},

	{
		nameEqFredInt.Or(where.Gt("age", 10)),
		` WHERE ("name"=?) OR ("age">?)`,
		`("name"='Fred') OR ("age">10)`,
		[]interface{}{"Fred", 10},
	},

	{
		nameEqFredInt.And(ageGt5Int).And(where.Gt("weight", 15)),
		` WHERE ("name"=?) AND ("age">?) AND ("weight">?)`,
		`("name"='Fred') AND ("age">5) AND ("weight">15)`,
		[]interface{}{"Fred", 5, 15},
	},

	{
		nameEqFredInt.Or(ageGt5Int).Or(where.Gt("weight", 15)),
		` WHERE ("name"=?) OR ("age">?) OR ("weight">?)`,
		`("name"='Fred') OR ("age">5) OR ("weight">15)`,
		[]interface{}{"Fred", 5, 15},
	},

	{
		where.Between("age", 12, 18).Or(where.Gt("weight", 45)),
		` WHERE ("age" BETWEEN ? AND ?) OR ("weight">?)`,
		`("age" BETWEEN 12 AND 18) OR ("weight">45)`,
		[]interface{}{12, 18, 45},
	},

	{
		where.GtEq("age", 10),
		` WHERE "age">=?`,
		`"age">=10`,
		[]interface{}{10},
	},

	{
		where.LtEq("age", 10),
		` WHERE "age"<=?`,
		`"age"<=10`,
		[]interface{}{10},
	},

	{
		where.NotEq("age", 10),
		` WHERE "age"<>?`,
		`"age"<>10`,
		[]interface{}{10},
	},

	{
		where.In("age", 10, 12, 14),
		` WHERE "age" IN (?,?,?)`,
		`"age" IN (10,12,14)`,
		[]interface{}{10, 12, 14},
	},

	{
		where.In("age", []int{10, 12, 14}),
		` WHERE "age" IN (?,?,?)`,
		`"age" IN (10,12,14)`,
		[]interface{}{10, 12, 14},
	},

	{ // 'In' without any vararg parameters
		where.In("age"),
		``,
		``,
		nil,
	},

	{ // 'In' with mixed value and nil vararg parameters
		where.In("age", 1, nil, 2, nil),
		` WHERE ("age" IN (?,?)) OR ("age" IS NULL)`,
		`("age" IN (1,2)) OR ("age" IS NULL)`,
		[]interface{}{1, 2},
	},

	{ // 'In' with only a nil vararg parameter
		where.In("age", nil),
		` WHERE ("age" IS NULL)`,
		`("age" IS NULL)`,
		nil,
	},

	{
		where.Not(nameEqFredInt),
		` WHERE NOT ("name"=?)`,
		`NOT ("name"='Fred')`,
		[]interface{}{"Fred"},
	},

	{
		where.Not(nameEqFredInt.And(ageLt10Int)),
		` WHERE NOT (("name"=?) AND ("age"<?))`,
		`NOT (("name"='Fred') AND ("age"<10))`,
		[]interface{}{"Fred", 10},
	},

	{
		where.Not(nameEqFredInt.Or(ageLt10Int)),
		` WHERE NOT (("name"=?) OR ("age"<?))`,
		`NOT (("name"='Fred') OR ("age"<10))`,
		[]interface{}{"Fred", 10},
	},

	{
		where.Not(nameEqFredInt).And(ageLt10Int),
		` WHERE (NOT ("name"=?)) AND ("age"<?)`,
		`(NOT ("name"='Fred')) AND ("age"<10)`,
		[]interface{}{"Fred", 10},
	},

	{
		where.Not(nameEqFredInt).Or(ageLt10Int),
		` WHERE (NOT ("name"=?)) OR ("age"<?)`,
		`(NOT ("name"='Fred')) OR ("age"<10)`,
		[]interface{}{"Fred", 10},
	},

	{
		where.And(nameEqFredInt, ageLt10Int),
		` WHERE ("name"=?) AND ("age"<?)`,
		`("name"='Fred') AND ("age"<10)`,
		[]interface{}{"Fred", 10},
	},

	{
		where.And(nameEqFredInt).And(where.And(ageLt10Int)),
		` WHERE ("name"=?) AND ("age"<?)`,
		`("name"='Fred') AND ("age"<10)`,
		[]interface{}{"Fred", 10},
	},

	{
		where.Or(nameEqFredInt, ageLt10Int),
		` WHERE ("name"=?) OR ("age"<?)`,
		`("name"='Fred') OR ("age"<10)`,
		[]interface{}{"Fred", 10},
	},

	{
		where.And(nameEqFredInt.Or(nameEqJohnInt), ageLt10Int),
		` WHERE (("name"=?) OR ("name"=?)) AND ("age"<?)`,
		`(("name"='Fred') OR ("name"='John')) AND ("age"<10)`,
		[]interface{}{"Fred", "John", 10},
	},

	{
		where.Or(nameEqFredInt, ageLt10Int.And(ageGt5Int)),
		` WHERE ("name"=?) OR (("age"<?) AND ("age">?))`,
		`("name"='Fred') OR (("age"<10) AND ("age">5))`,
		[]interface{}{"Fred", 10, 5},
	},

	{
		where.Or(nameEqFredInt, nameEqJohnInt).And(ageGt5Int),
		` WHERE (("name"=?) OR ("name"=?)) AND ("age">?)`,
		`(("name"='Fred') OR ("name"='John')) AND ("age">5)`,
		[]interface{}{"Fred", "John", 5},
	},

	{
		where.Or(nameEqFredInt, nameEqJohnInt, where.And(ageGt5Int)),
		` WHERE ("name"=?) OR ("name"=?) OR (("age">?))`,
		`("name"='Fred') OR ("name"='John') OR (("age">5))`,
		[]interface{}{"Fred", "John", 5},
	},

	{
		where.Or().Or(where.NoOp()).And(where.NoOp()),
		"",
		"",
		nil,
	},

	{
		where.And(where.Or(where.NoOp())),
		"",
		"",
		nil,
	},
}

func TestBuildWhereClause_happyCases(t *testing.T) {
	g := NewGomegaWithT(t)

	for i, c := range buildWhereClauseHappyCases {
		info := fmt.Sprintf("%d: %s", i, c.expSQL)

		sql, args := where.Where(c.wh)

		g.Expect(sql).To(Equal(c.expSQL), info)
		g.Expect(args).To(Equal(c.args), info)

		s := c.wh.String()

		g.Expect(s).To(Equal(c.expString), info)
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
	expr = dialect.ReplacePlaceholdersWithNumbers(expr, "$")
	fmt.Println(expr)
	fmt.Println(args)

	// Output: WHERE ((name=$1) OR (name=$2)) AND (age>$3) AND (likes IN ($4,$5))
	// [John Peter 10 cats dogs]
}
