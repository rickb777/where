package dialect

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestReplacePlaceholders(t *testing.T) {
	g := NewGomegaWithT(t)

	s := ReplacePlaceholders("?,?,?,?,?,?,?,?,?,?,?", Query)
	g.Expect(s).Should(Equal("?,?,?,?,?,?,?,?,?,?,?"))

	s = ReplacePlaceholders("?,?,?,?,?,?,?,?,?,?,?", Dollar)
	g.Expect(s).Should(Equal("$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11"))

	s = ReplacePlaceholders("?,?,?,?,?,?,?,?,?,?,?", Dollar, 11)
	g.Expect(s).Should(Equal("$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21"))
}
