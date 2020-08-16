package quote

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestAnsiQuote(t *testing.T) {
	g := NewGomegaWithT(t)
	result := PickQuoter("ansi").QuoteN([]string{"a", "bb", "ccc.ddd"})
	g.Expect(result).To(HaveLen(3))
	g.Expect(result[0]).To(Equal(`"a"`))
	g.Expect(result[1]).To(Equal(`"bb"`))
	g.Expect(result[2]).To(Equal(`"ccc"."ddd"`))
}

func TestMysqlQuote(t *testing.T) {
	g := NewGomegaWithT(t)
	result := PickQuoter("mysql").QuoteN([]string{"a", "bb", "ccc.ddd"})
	g.Expect(result).To(HaveLen(3))
	g.Expect(result[0]).To(Equal("`a`"))
	g.Expect(result[1]).To(Equal("`bb`"))
	g.Expect(result[2]).To(Equal("`ccc`.`ddd`"))
}

func TestNoQuote(t *testing.T) {
	g := NewGomegaWithT(t)
	result := PickQuoter("none").QuoteN([]string{"a", "bb", "ccc.ddd"})
	g.Expect(result).To(HaveLen(3))
	g.Expect(result[0]).To(Equal(`a`))
	g.Expect(result[1]).To(Equal(`bb`))
	g.Expect(result[2]).To(Equal(`ccc.ddd`))

	r2 := NewQuoter("").Quote("foo")
	g.Expect(r2).To(Equal(`foo`))

	b := &strings.Builder{}
	NewQuoter("").QuoteW(b, "foo")
	g.Expect(b.String()).To(Equal(`foo`))
}
