package quote

import (
	"github.com/rickb777/expect"
	"strings"
	"testing"
)

func TestAnsiQuote(t *testing.T) {
	cases := []struct {
		identifier, dialect, expected string
	}{
		{
			identifier: "",
			expected:   ``,
			dialect:    "none",
		},
		{
			identifier: "",
			expected:   ``,
			dialect:    "ansi",
		},
		{
			identifier: "",
			expected:   "",
			dialect:    "mysql",
		},
		{
			identifier: "ccc",
			expected:   `ccc`,
			dialect:    "none",
		},
		{
			identifier: "ccc",
			expected:   `"ccc"`,
			dialect:    "ansi",
		},
		{
			identifier: "ccc",
			expected:   `"ccc"`,
			dialect:    "postgres",
		},
		{
			identifier: "ccc",
			expected:   `"ccc"`,
			dialect:    "sqlite",
		},
		{
			identifier: "ccc",
			expected:   "`ccc`",
			dialect:    "mysql",
		},
		{
			identifier: "ccc",
			expected:   "[ccc]",
			dialect:    "ms-sql",
		},
		{
			identifier: "a.ccc.ddd",
			expected:   `a.ccc.ddd`,
			dialect:    "none",
		},
		{
			identifier: "a.ccc.ddd",
			expected:   `"a"."ccc"."ddd"`,
			dialect:    "ansi",
		},
		{
			identifier: "a.ccc.ddd",
			expected:   "`a`.`ccc`.`ddd`",
			dialect:    "mysql",
		},
		{
			identifier: "a.ccc.ddd",
			expected:   "`a`.`ccc`.`ddd`",
			dialect:    "backtick",
		},
		{
			identifier: "a.ccc.ddd",
			expected:   "[a].[ccc].[ddd]",
			dialect:    "mssql",
		},
		{
			identifier: "a,ccc,ddd",
			expected:   "a,ccc,ddd",
			dialect:    "mssql",
		},
		{
			identifier: "a ccc ddd",
			expected:   "a ccc ddd",
			dialect:    "mssql",
		},
	}

	for i, c := range cases {
		t.Logf("%d: %s", i, c.identifier)

		s1 := Pick(c.dialect).Quote(c.identifier)
		expect.String(s1).Info(i).ToBe(t, c.expected)

		buf := &strings.Builder{}
		Pick(c.dialect).QuoteW(buf, c.identifier)
		s2 := buf.String()
		expect.String(s2).Info(i).ToBe(t, c.expected)
	}
}
