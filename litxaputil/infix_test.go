package litxaputil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestInfixMap_Typos(t *testing.T) {
	for key := range infixMap {
		infix := infixMap[key]
		assert.Equal(t, key, strings.Join(infix.SyllableSplit, ""))
	}
}

func TestFindInfix(t *testing.T) {
	assert.Nil(t, FindInfix("glurb"))
	assert.Equal(t, infixMap["uy"], *FindInfix("uy"))
}

func TestInfix_Apply(t *testing.T) {
	table := []struct {
		curr     string
		infix    string
		si       int
		pos      int
		expected string
		si2      int
		pos2     int
	}{
		{
			curr: "ta.ron", infix: "imv", si: 0, pos: 1,
			expected: "tim.va.ron", si2: 1, pos2: 1,
		},
		{
			curr: "frr.fen", infix: "äng", si: 1, pos: 1,
			expected: "frr.fä.ngen", si2: 2, pos2: 2,
		},
		{
			curr: "tel", infix: "ei", si: 0, pos: 1,
			expected: "te.i.el", si2: 2, pos2: 0,
		},
		{
			curr: "srew", infix: "äpeyk", si: 0, pos: 2,
			expected: "srä.pey.kew", si2: 2, pos2: 1,
		},
		{
			curr: "tsway.on", infix: "ìyev", si: 0, pos: 3,
			expected: "tswì.ye.vay.on", si2: 2, pos2: 1,
		},
		{
			curr: "u.van. si", infix: "eiy", si: 2, pos: 2,
			expected: "u.van. se.i.yi", si2: 4, pos2: 1,
		},
		{
			curr: "fme.tok", infix: "ol", si: 0, pos: 2,
			expected: "fmo.le.tok", si2: 1, pos2: 1,
		},
		{
			curr: "om.um", infix: "iv", si: 0, pos: 0,
			expected: "i.vom.um", si2: 1, pos2: 1,
		},
		{
			curr: "eyk", infix: "äpeyk", si: 0, pos: 0,
			expected: "ä.pey.keyk", si2: 2, pos2: 1,
		},
		{
			curr: "ka.me", infix: "ei", si: 1, pos: 1,
			expected: "ka.me.i.e", si2: 3, pos2: 0,
		},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("%s<%s>", row.curr, row.infix), func(t *testing.T) {
			res, si2, pos2 := FindInfix(row.infix).Apply(strings.Split(row.curr, "."), row.si, row.pos)
			assert.Equal(t, row.expected, strings.Join(res, "."))
			assert.Equal(t, row.si2, si2)
			assert.Equal(t, row.pos2, pos2)
		})
	}
}

func TestApplyInfixes(t *testing.T) {
	table := []struct {
		curr              string
		infixes           string
		start             int
		stress            int
		positions         [2][2]int
		expectedSyllables string
		expectedStress    int
	}{
		{
			curr: "ta.ron", infixes: "imv",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 1}, {1, 1}},

			expectedSyllables: "tim.va.ron", expectedStress: 1,
		},
		{
			curr: "ta.rep", infixes: "äng",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 1}, {1, 1}},

			expectedSyllables: "ta.rä.ngep", expectedStress: 0,
		},
		{
			curr: "fme.tok", infixes: "äpeyk",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 2}, {1, 1}},

			expectedSyllables: "fmä.pey.ke.tok", expectedStress: 2,
		},
		{
			curr: "fme.tok", infixes: "äp,eyk",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 2}, {1, 1}},

			expectedSyllables: "fmä.pey.ke.tok", expectedStress: 2,
		},
		{
			curr: "tel", infixes: "er,äng",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 1}, {0, 1}},

			expectedSyllables: "te.rä.ngel", expectedStress: 2,
		},
		{
			curr: "pll.txe", infixes: "iyev",
			start: 0, stress: 1,
			positions: [2][2]int{{0, 1}, {1, 2}},

			expectedSyllables: "pi.ye.vll.txe", expectedStress: 3,
		},
		{
			curr: "txa.krr.fpìl", infixes: "eyk,iyev,eiy",
			start: 0, stress: 1,
			positions: [2][2]int{{2, 2}, {2, 2}},

			expectedSyllables: "txa.krr.fpey.ki.ye.ve.i.yìl", expectedStress: 1,
		},
		{
			curr: "kä", infixes: "am,ei",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 1}, {0, 1}},

			expectedSyllables: "ka.me.i.ä", expectedStress: 3,
		},
		{
			curr: "vll", infixes: "ol",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 1}, {0, 1}},

			expectedSyllables: "vol", expectedStress: 0,
		},
		{
			curr: "vll", infixes: "ol,äng",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 1}, {0, 1}},

			expectedSyllables: "vo.lä.ngll", expectedStress: 2,
		},
		{
			curr: "nrr", infixes: "er",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 1}, {0, 1}},

			expectedSyllables: "ner", expectedStress: 0,
		},
		{
			curr: "pll.txe", infixes: "ol,ei",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 1}, {1, 2}},

			expectedSyllables: "pol.txe.i.e", expectedStress: 0,
		},

		// Stress shift
		{
			curr: "i.nan", infixes: "er",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 0}, {1, 1}},

			expectedSyllables: "e.ri.nan", expectedStress: 0,
		},
		{
			curr: "i.nan", infixes: "ìyev",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 0}, {1, 1}},

			expectedSyllables: "ì.ye.vi.nan", expectedStress: 1,
		},
		{
			curr: "o.mum", infixes: "eyk,ol",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 0}, {1, 1}},

			expectedSyllables: "ey.ko.lo.mum", expectedStress: 1,
		},
		{
			curr: "eyk", infixes: "äp",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 0}, {0, 0}},

			expectedSyllables: "ä.peyk", expectedStress: 0,
		},
		{
			curr: "eyk", infixes: "äpeyk",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 0}, {0, 0}},

			expectedSyllables: "ä.pey.keyk", expectedStress: 1,
		},
		{
			curr: "eyk", infixes: "äpeyk,iyev,ei",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 0}, {0, 0}},

			expectedSyllables: "ä.pey.ki.ye.ve.i.eyk", expectedStress: 5,
		},
		{
			curr: "eyk", infixes: "äp,er,äng",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 0}, {0, 0}},

			expectedSyllables: "ä.pe.rä.ngeyk", expectedStress: 2,
		},
		{
			curr: "eyk", infixes: "äp,er,äng",
			start: 0, stress: 0,
			positions: [2][2]int{{0, 0}, {0, 0}},

			expectedSyllables: "ä.pe.rä.ngeyk", expectedStress: 2,
		},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("%s<%s>", row.curr, row.infixes), func(t *testing.T) {
			resSyllables, resStress := ApplyInfixes(strings.Split(row.curr, "."), strings.Split(row.infixes, ","), row.start, row.stress, row.positions)
			assert.Equal(t, strings.Split(row.expectedSyllables, "."), resSyllables)
			assert.Equal(t, row.expectedStress, resStress)
		})
	}
}
