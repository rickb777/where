package quote

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestAnsiQuote(t *testing.T) {
	g := NewGomegaWithT(t)
	result := Pick("ansi").QuoteN([]string{"a", "bb", "ccc.ddd", ""})
	g.Expect(result).To(HaveLen(4))
	g.Expect(result[0]).To(Equal(`"a"`))
	g.Expect(result[1]).To(Equal(`"bb"`))
	g.Expect(result[2]).To(Equal(`"ccc"."ddd"`))
	g.Expect(result[3]).To(Equal(``))
}

func TestMysqlQuote(t *testing.T) {
	g := NewGomegaWithT(t)
	result := Pick("mysql").QuoteN([]string{"a", "bb", "ccc.ddd", ""})
	g.Expect(result).To(HaveLen(4))
	g.Expect(result[0]).To(Equal("`a`"))
	g.Expect(result[1]).To(Equal("`bb`"))
	g.Expect(result[2]).To(Equal("`ccc`.`ddd`"))
	g.Expect(result[3]).To(Equal(``))
}

func TestNoQuote(t *testing.T) {
	g := NewGomegaWithT(t)
	result := Pick("none").QuoteN([]string{"a", "bb", "ccc.ddd", ""})
	g.Expect(result).To(HaveLen(4))
	g.Expect(result[0]).To(Equal(`a`))
	g.Expect(result[1]).To(Equal(`bb`))
	g.Expect(result[2]).To(Equal(`ccc.ddd`))
	g.Expect(result[3]).To(Equal(``))

	r2 := New("").Quote("foo")
	g.Expect(r2).To(Equal(`foo`))

	b := &strings.Builder{}
	New("").QuoteW(b, "foo")
	g.Expect(b.String()).To(Equal(`foo`))
}
