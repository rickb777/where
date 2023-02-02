package quote

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestAnsiQuote(t *testing.T) {
	g := NewGomegaWithT(t)

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
	}

	for _, c := range cases {
		s1 := Pick(c.dialect).Quote(c.identifier)
		g.Expect(s1).To(Equal(c.expected))

		buf := &strings.Builder{}
		Pick(c.dialect).QuoteW(buf, c.identifier)
		s2 := buf.String()
		g.Expect(s2).To(Equal(c.expected))
	}
}
