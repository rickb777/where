package where_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/rickb777/where/v2"
	"github.com/rickb777/where/v2/dialect"
	"github.com/rickb777/where/v2/quote"
)

var nameEqFredInt = where.Eq("name", "Fred")
var nameEqJohnInt = where.Eq("name", "John")
var ageLt10Int = where.Lt("age", 10)
var ageGt5Int = where.Gt("age", 5)

var (
	buildWhereClauseHappyCases = []struct {
		wh           where.Expression
		expMySql     string
		expPostgres  string
		expSqlServer string
		expString    string
		args         []any
	}{
		{wh: where.NoOp()},

		{
			wh:           where.Condition{Column: "name", Predicate: " not nil", Args: nil},
			expMySql:     " WHERE `name` not nil",
			expPostgres:  ` WHERE "name" not nil`,
			expSqlServer: ` WHERE [name] not nil`,
			expString:    `name not nil`,
		},

		{
			wh:           where.Condition{Column: "p.name", Predicate: " not nil", Args: nil},
			expMySql:     " WHERE `p`.`name` not nil",
			expPostgres:  ` WHERE "p"."name" not nil`,
			expSqlServer: ` WHERE [p].[name] not nil`,
			expString:    `p.name not nil`,
		},

		{
			wh:           where.Null("name"),
			expMySql:     " WHERE `name` IS NULL",
			expPostgres:  ` WHERE "name" IS NULL`,
			expSqlServer: ` WHERE [name] IS NULL`,
			expString:    `name IS NULL`,
		},

		{
			wh:           where.NotNull("name"),
			expMySql:     " WHERE `name` IS NOT NULL",
			expPostgres:  ` WHERE "name" IS NOT NULL`,
			expSqlServer: ` WHERE [name] IS NOT NULL`,
			expString:    `name IS NOT NULL`,
		},

		{
			wh:           where.Condition{Column: "name", Predicate: " <>?", Args: []any{"Boo"}},
			expMySql:     " WHERE `name` <>?",
			expPostgres:  ` WHERE "name" <>$1`,
			expSqlServer: ` WHERE [name] <>@p1`,
			expString:    `name <>'Boo'`,
			args:         []any{"Boo"},
		},

		{
			wh:           nameEqFredInt,
			expMySql:     " WHERE `name`=?",
			expPostgres:  ` WHERE "name"=$1`,
			expSqlServer: ` WHERE [name]=@p1`,
			expString:    `name='Fred'`,
			args:         []any{"Fred"},
		},

		{
			wh:           where.Like("name", "F%"),
			expMySql:     " WHERE `name` LIKE ?",
			expPostgres:  ` WHERE "name" LIKE $1`,
			expSqlServer: ` WHERE [name] LIKE @p1`,
			expString:    `name LIKE 'F%'`,
			args:         []any{"F%"},
		},

		{
			wh:           where.NoOp().And(nameEqFredInt),
			expMySql:     " WHERE (`name`=?)",
			expPostgres:  ` WHERE ("name"=$1)`,
			expSqlServer: ` WHERE ([name]=@p1)`,
			expString:    `(name='Fred')`,
			args:         []any{"Fred"},
		},

		{
			wh:           nameEqFredInt.And(where.NoOp()),
			expMySql:     " WHERE (`name`=?)",
			expPostgres:  ` WHERE ("name"=$1)`,
			expSqlServer: ` WHERE ([name]=@p1)`,
			expString:    `(name='Fred')`,
			args:         []any{"Fred"},
		},

		{
			wh:           nameEqFredInt.And(where.Gt("age", 10)),
			expMySql:     " WHERE (`name`=?) AND (`age`>?)",
			expPostgres:  ` WHERE ("name"=$1) AND ("age">$2)`,
			expSqlServer: ` WHERE ([name]=@p1) AND ([age]>@p2)`,
			expString:    `(name='Fred') AND (age>10)`,
			args:         []any{"Fred", 10},
		},

		{
			wh:           nameEqFredInt.Or(where.Gt("age", 10)),
			expMySql:     " WHERE (`name`=?) OR (`age`>?)",
			expPostgres:  ` WHERE ("name"=$1) OR ("age">$2)`,
			expSqlServer: ` WHERE ([name]=@p1) OR ([age]>@p2)`,
			expString:    `(name='Fred') OR (age>10)`,
			args:         []any{"Fred", 10},
		},

		{
			wh:           nameEqFredInt.And(ageGt5Int).And(where.Gt("weight", 15)),
			expMySql:     " WHERE (`name`=?) AND (`age`>?) AND (`weight`>?)",
			expPostgres:  ` WHERE ("name"=$1) AND ("age">$2) AND ("weight">$3)`,
			expSqlServer: ` WHERE ([name]=@p1) AND ([age]>@p2) AND ([weight]>@p3)`,
			expString:    `(name='Fred') AND (age>5) AND (weight>15)`,
			args:         []any{"Fred", 5, 15},
		},

		{
			wh:           nameEqFredInt.Or(ageGt5Int).Or(where.Gt("weight", 15)),
			expMySql:     " WHERE (`name`=?) OR (`age`>?) OR (`weight`>?)",
			expPostgres:  ` WHERE ("name"=$1) OR ("age">$2) OR ("weight">$3)`,
			expSqlServer: ` WHERE ([name]=@p1) OR ([age]>@p2) OR ([weight]>@p3)`,
			expString:    `(name='Fred') OR (age>5) OR (weight>15)`,
			args:         []any{"Fred", 5, 15},
		},

		{
			wh:           where.Between("age", 12, 18).Or(where.Gt("weight", 45)),
			expMySql:     " WHERE (`age` BETWEEN ? AND ?) OR (`weight`>?)",
			expPostgres:  ` WHERE ("age" BETWEEN $1 AND $2) OR ("weight">$3)`,
			expSqlServer: ` WHERE ([age] BETWEEN @p1 AND @p2) OR ([weight]>@p3)`,
			expString:    `(age BETWEEN 12 AND 18) OR (weight>45)`,
			args:         []any{12, 18, 45},
		},

		{
			wh:           where.GtEq("age", 10),
			expMySql:     " WHERE `age`>=?",
			expPostgres:  ` WHERE "age">=$1`,
			expSqlServer: ` WHERE [age]>=@p1`,
			expString:    `age>=10`,
			args:         []any{10},
		},

		{
			wh:           where.LtEq("age", 10),
			expMySql:     " WHERE `age`<=?",
			expPostgres:  ` WHERE "age"<=$1`,
			expSqlServer: ` WHERE [age]<=@p1`,
			expString:    `age<=10`,
			args:         []any{10},
		},

		{
			wh:           where.NotEq("age", 10),
			expMySql:     " WHERE `age`<>?",
			expPostgres:  ` WHERE "age"<>$1`,
			expSqlServer: ` WHERE [age]<>@p1`,
			expString:    `age<>10`,
			args:         []any{10},
		},

		{
			wh:           where.In("age", int8(10), int16(12), int32(14)),
			expMySql:     " WHERE `age` IN (?,?,?)",
			expPostgres:  ` WHERE "age" IN ($1,$2,$3)`,
			expSqlServer: ` WHERE [age] IN (@p1,@p2,@p3)`,
			expString:    `age IN (10,12,14)`,
			args:         []any{int8(10), int16(12), int32(14)},
		},

		{ // 'In' without any vararg parameters
			wh: where.In("age"),
		},

		{ // 'In' with mixed value and nil vararg parameters
			wh:           where.In("age", 1, nil, 2, nil),
			expMySql:     " WHERE (`age` IN (?,?)) OR (`age` IS NULL)",
			expPostgres:  ` WHERE ("age" IN ($1,$2)) OR ("age" IS NULL)`,
			expSqlServer: ` WHERE ([age] IN (@p1,@p2)) OR ([age] IS NULL)`,
			expString:    `(age IN (1,2)) OR (age IS NULL)`,
			args:         []any{1, 2},
		},

		{ // 'In' with only a nil vararg parameter
			wh:           where.In("age", nil),
			expMySql:     " WHERE `age` IS NULL",
			expPostgres:  ` WHERE "age" IS NULL`,
			expSqlServer: ` WHERE [age] IS NULL`,
			expString:    `age IS NULL`,
		},

		{
			wh:           where.InSlice("ages", []uint{10, 12, 14}),
			expMySql:     " WHERE `ages` IN (?,?,?)",
			expPostgres:  ` WHERE "ages" IN ($1,$2,$3)`,
			expSqlServer: ` WHERE [ages] IN (@p1,@p2,@p3)`,
			expString:    `ages IN (10,12,14)`,
			args:         []any{uint(10), uint(12), uint(14)},
		},

		{ // 'InSlice' with mixed value and nil parameters
			wh:           where.InSlice("ages", []any{1, nil, uint32(2), nil}),
			expMySql:     " WHERE (`ages` IN (?,?)) OR (`ages` IS NULL)",
			expPostgres:  ` WHERE ("ages" IN ($1,$2)) OR ("ages" IS NULL)`,
			expSqlServer: ` WHERE ([ages] IN (@p1,@p2)) OR ([ages] IS NULL)`,
			expString:    `(ages IN (1,2)) OR (ages IS NULL)`,
			args:         []any{1, uint32(2)},
		},

		{ // 'InSlice' with only a nil parameter
			wh: where.InSlice("ages", nil),
		},

		{
			wh:           nameEqFredInt.Or(nameEqJohnInt),
			expMySql:     " WHERE (`name`=?) OR (`name`=?)",
			expPostgres:  ` WHERE ("name"=$1) OR ("name"=$2)`,
			expSqlServer: ` WHERE ([name]=@p1) OR ([name]=@p2)`,
			expString:    `(name='Fred') OR (name='John')`,
			args:         []any{"Fred", "John"},
		},

		{
			wh:           where.Not(nameEqFredInt),
			expMySql:     " WHERE NOT (`name`=?)",
			expPostgres:  ` WHERE NOT ("name"=$1)`,
			expSqlServer: ` WHERE NOT ([name]=@p1)`,
			expString:    `NOT (name='Fred')`,
			args:         []any{"Fred"},
		},

		{
			wh:           where.Not(nameEqFredInt.And(ageLt10Int)),
			expMySql:     " WHERE NOT ((`name`=?) AND (`age`<?))",
			expPostgres:  ` WHERE NOT (("name"=$1) AND ("age"<$2))`,
			expSqlServer: ` WHERE NOT (([name]=@p1) AND ([age]<@p2))`,
			expString:    `NOT ((name='Fred') AND (age<10))`,
			args:         []any{"Fred", 10},
		},

		{
			wh:           where.Not(nameEqFredInt.Or(ageLt10Int)),
			expMySql:     " WHERE NOT ((`name`=?) OR (`age`<?))",
			expPostgres:  ` WHERE NOT (("name"=$1) OR ("age"<$2))`,
			expSqlServer: ` WHERE NOT (([name]=@p1) OR ([age]<@p2))`,
			expString:    `NOT ((name='Fred') OR (age<10))`,
			args:         []any{"Fred", 10},
		},

		{
			wh:           where.Not(nameEqFredInt).And(ageLt10Int),
			expMySql:     " WHERE (NOT (`name`=?)) AND (`age`<?)",
			expPostgres:  ` WHERE (NOT ("name"=$1)) AND ("age"<$2)`,
			expSqlServer: ` WHERE (NOT ([name]=@p1)) AND ([age]<@p2)`,
			expString:    `(NOT (name='Fred')) AND (age<10)`,
			args:         []any{"Fred", 10},
		},

		{
			wh:           where.Not(nameEqFredInt).Or(ageLt10Int),
			expMySql:     " WHERE (NOT (`name`=?)) OR (`age`<?)",
			expPostgres:  ` WHERE (NOT ("name"=$1)) OR ("age"<$2)`,
			expSqlServer: ` WHERE (NOT ([name]=@p1)) OR ([age]<@p2)`,
			expString:    `(NOT (name='Fred')) OR (age<10)`,
			args:         []any{"Fred", 10},
		},

		{
			wh:           where.And(nameEqFredInt, nil, ageLt10Int),
			expMySql:     " WHERE (`name`=?) AND (`age`<?)",
			expPostgres:  ` WHERE ("name"=$1) AND ("age"<$2)`,
			expSqlServer: ` WHERE ([name]=@p1) AND ([age]<@p2)`,
			expString:    `(name='Fred') AND (age<10)`,
			args:         []any{"Fred", 10},
		},

		{
			wh:           where.Or(nameEqFredInt, nil, ageLt10Int),
			expMySql:     " WHERE (`name`=?) OR (`age`<?)",
			expPostgres:  ` WHERE ("name"=$1) OR ("age"<$2)`,
			expSqlServer: ` WHERE ([name]=@p1) OR ([age]<@p2)`,
			expString:    `(name='Fred') OR (age<10)`,
			args:         []any{"Fred", 10},
		},

		{
			wh:           where.And(nameEqFredInt).And(where.And(ageLt10Int)),
			expMySql:     " WHERE (`name`=?) AND (`age`<?)",
			expPostgres:  ` WHERE ("name"=$1) AND ("age"<$2)`,
			expSqlServer: ` WHERE ([name]=@p1) AND ([age]<@p2)`,
			expString:    `(name='Fred') AND (age<10)`,
			args:         []any{"Fred", 10},
		},

		{
			wh:           where.Eq("a", 1).And(where.Eq("b", 2)).And(where.And(where.Eq("c", 3), where.Eq("d", 4))),
			expMySql:     " WHERE (`a`=?) AND (`b`=?) AND (`c`=?) AND (`d`=?)",
			expPostgres:  ` WHERE ("a"=$1) AND ("b"=$2) AND ("c"=$3) AND ("d"=$4)`,
			expSqlServer: ` WHERE ([a]=@p1) AND ([b]=@p2) AND ([c]=@p3) AND ([d]=@p4)`,
			expString:    `(a=1) AND (b=2) AND (c=3) AND (d=4)`,
			args:         []any{1, 2, 3, 4},
		},

		{
			wh:           where.Eq("a", 1).Or(where.Eq("b", 2)).Or(where.Or(where.Eq("c", 3), where.Eq("d", 4))),
			expMySql:     " WHERE (`a`=?) OR (`b`=?) OR (`c`=?) OR (`d`=?)",
			expPostgres:  ` WHERE ("a"=$1) OR ("b"=$2) OR ("c"=$3) OR ("d"=$4)`,
			expSqlServer: ` WHERE ([a]=@p1) OR ([b]=@p2) OR ([c]=@p3) OR ([d]=@p4)`,
			expString:    `(a=1) OR (b=2) OR (c=3) OR (d=4)`,
			args:         []any{1, 2, 3, 4},
		},

		{
			wh:           where.And(nameEqFredInt.Or(nameEqJohnInt), ageLt10Int),
			expMySql:     " WHERE ((`name`=?) OR (`name`=?)) AND (`age`<?)",
			expPostgres:  ` WHERE (("name"=$1) OR ("name"=$2)) AND ("age"<$3)`,
			expSqlServer: ` WHERE (([name]=@p1) OR ([name]=@p2)) AND ([age]<@p3)`,
			expString:    `((name='Fred') OR (name='John')) AND (age<10)`,
			args:         []any{"Fred", "John", 10},
		},

		{
			wh:           where.Or(nameEqFredInt, ageLt10Int.And(ageGt5Int)),
			expMySql:     " WHERE (`name`=?) OR ((`age`<?) AND (`age`>?))",
			expPostgres:  ` WHERE ("name"=$1) OR (("age"<$2) AND ("age">$3))`,
			expSqlServer: ` WHERE ([name]=@p1) OR (([age]<@p2) AND ([age]>@p3))`,
			expString:    `(name='Fred') OR ((age<10) AND (age>5))`,
			args:         []any{"Fred", 10, 5},
		},

		{
			wh:           where.Or(nameEqFredInt, nameEqJohnInt).And(ageGt5Int),
			expMySql:     " WHERE ((`name`=?) OR (`name`=?)) AND (`age`>?)",
			expPostgres:  ` WHERE (("name"=$1) OR ("name"=$2)) AND ("age">$3)`,
			expSqlServer: ` WHERE (([name]=@p1) OR ([name]=@p2)) AND ([age]>@p3)`,
			expString:    `((name='Fred') OR (name='John')) AND (age>5)`,
			args:         []any{"Fred", "John", 5},
		},

		{
			wh:           where.Or(nameEqFredInt, nameEqJohnInt, where.And(ageGt5Int)),
			expMySql:     " WHERE (`name`=?) OR (`name`=?) OR (`age`>?)",
			expPostgres:  ` WHERE ("name"=$1) OR ("name"=$2) OR ("age">$3)`,
			expSqlServer: ` WHERE ([name]=@p1) OR ([name]=@p2) OR ([age]>@p3)`,
			expString:    `(name='Fred') OR (name='John') OR (age>5)`,
			args:         []any{"Fred", "John", 5},
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

func TestBuildWhereClause_Mysql_happyCases(t *testing.T) {
	g := NewGomegaWithT(t)

	for i, c := range buildWhereClauseHappyCases {
		t.Logf("%d: %s", i, c.expMySql)

		sql, args := where.Where(c.wh, dialect.MySqlQuotes, dialect.Query)

		g.Expect(sql).To(Equal(c.expMySql))
		g.Expect(args).To(Equal(c.args))

		s := c.wh.String()

		g.Expect(s).To(Equal(c.expString))
	}
}

func TestBuildWhereClause_Postgres_happyCases(t *testing.T) {
	g := NewGomegaWithT(t)

	for i, c := range buildWhereClauseHappyCases {
		t.Logf("%d: %s", i, c.expPostgres)

		sql, args := where.Where(c.wh, dialect.ANSIQuotes, dialect.Dollar)

		g.Expect(sql).To(Equal(c.expPostgres))
		g.Expect(args).To(Equal(c.args))

		s := c.wh.String()

		g.Expect(s).To(Equal(c.expString))
	}
}

func TestBuildHavingClause_SQLServer_happyCases(t *testing.T) {
	g := NewGomegaWithT(t)

	for i, c := range buildWhereClauseHappyCases {
		t.Logf("%d: %s", i, c.expSqlServer)

		sql, args := where.Having(c.wh, dialect.SquareBrackets, dialect.AtP)

		exp := c.expSqlServer
		if exp != "" {
			exp = " HAVING " + exp[7:]
		}
		g.Expect(sql).To(Equal(exp))
		g.Expect(args).To(Equal(c.args))

		s := c.wh.String()

		g.Expect(s).To(Equal(c.expString))
	}
}

func BenchmarkBuildWhereClause_happyCases_build(b *testing.B) {
	for _, c := range buildWhereClauseHappyCases {
		_, _ = where.Where(c.wh, dialect.ANSIQuotes)
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
	quote.DefaultQuoter = quote.None

	// some simple expressions
	nameEqJohn := where.Eq("name", "John")
	nameEqPeter := where.Eq("name", "Peter")
	ageGt10 := where.Gt("age", 10)
	likes := where.In("likes", "cats", "dogs")

	// build a compound expression - this is a static expression
	// but it could be based on conditions instead
	wh := where.And(where.Or(nameEqJohn, nameEqPeter), ageGt10, likes)

	// For Postgres, the placeholders have to be altered. It's necessary to do
	// this on the whole query if there might be other placeholders in it too.
	expr, args := where.Where(wh, dialect.Dollar)

	fmt.Println(expr)
	fmt.Println(args)

	// Output: WHERE ((name=$1) OR (name=$2)) AND (age>$3) AND (likes IN ($4,$5))
	// [John Peter 10 cats dogs]
}

func ExampleWhere_mysqlUsingParameters() {
	// some simple expressions
	nameEqJohn := where.Eq("name", "John")
	nameEqPeter := where.Eq("name", "Peter")
	ageGt10 := where.Gt("age", 10)
	likes := where.In("likes", "cats", "dogs")

	// Build a compound expression - this is a static expression
	// but it could be built up in stages depending on any conditions.
	wh := where.And(where.Or(nameEqJohn, nameEqPeter), ageGt10, likes)

	// Format the 'where' clause, quoting all the identifiers for MySql.
	clause, args := where.Where(wh, dialect.MySqlQuotes)

	fmt.Println(clause)
	fmt.Println(args)

	// Output: WHERE ((`name`=?) OR (`name`=?)) AND (`age`>?) AND (`likes` IN (?,?))
	// [John Peter 10 cats dogs]
}

func ExampleWhere_postgresUsingParameters() {
	// some simple expressions
	nameEqJohn := where.Eq("name", "John")
	nameEqPeter := where.Eq("name", "Peter")
	ageGt10 := where.Gt("age", 10)
	likes := where.In("likes", "cats", "dogs")

	// Build a compound expression - this is a static expression
	// but it could be built up in stages depending on any conditions.
	wh := where.And(where.Or(nameEqJohn, nameEqPeter), ageGt10, likes)

	// Format the 'where' clause, quoting all the identifiers for Postgres
	// and replacing all the '?' parameters with "$1" numbered parameters,
	// counting from 1.
	clause, args := where.Where(wh, dialect.ANSIQuotes, dialect.Dollar)

	fmt.Println(clause)
	fmt.Println(args)

	// Output: WHERE (("name"=$1) OR ("name"=$2)) AND ("age">$3) AND ("likes" IN ($4,$5))
	// [John Peter 10 cats dogs]
}
