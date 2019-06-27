package dialect

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestReplacePlaceholders(t *testing.T) {
	g := NewGomegaWithT(t)

	s := Mysql.ReplacePlaceholders("?,?,?,?,?,?,?,?,?,?,?", nil)
	g.Expect(s).Should(Equal("?,?,?,?,?,?,?,?,?,?,?"))

	s = Postgres.ReplacePlaceholders("?,?,?,?,?,?,?,?,?,?,?", nil)
	g.Expect(s).Should(Equal("$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11"))
}
