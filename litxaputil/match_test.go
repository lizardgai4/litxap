package litxaputil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMatchSyllables(t *testing.T) {
	table := []struct {
		word         string
		syllables    string
		root         int
		stress       int
		newSyllables string
		newStress    int
	}{
		{
			word: "tìfmetok", syllables: "tì.fme.tok",
			root: 0, stress: 1,
			newSyllables: "tì.fme.tok",
			newStress:    1,
		},
		{
			word: "säniu", syllables: "sä.nu.i",
			root: 0, stress: 1,
			newSyllables: "",
			newStress:    -1,
		},
		{
			word: "sìran", syllables: "tì.ran",
			root: 0, stress: 1,
			newSyllables: "sì.ran",
			newStress:    1,
		},
		{
			word: "sahilvan", syllables: "tsa.kil.van",
			root: 1, stress: 2,
			newSyllables: "sa.hil.van",
			newStress:    2,
		},
		{
			word: "senui", syllables: "sä.nu.i",
			root: 0, stress: 1,
			newSyllables: "se.nu.i",
			newStress:    1,
		},
		{
			word: "skahena", syllables: "sä.ka.he.na",
			root: 0, stress: 2,
			newSyllables: "ska.he.na",
			newStress:    1,
		},
		{
			word: "skeynven", syllables: "sä.keyn.ven",
			root: 0, stress: 2,
			newSyllables: "skeyn.ven",
			newStress:    1,
		},
		{
			word: "tsamsyu", syllables: "tsam. si.yu",
			root: 0, stress: 2,
			newSyllables: "tsam.syu",
			newStress:    1,
		},
		{
			word: "lehawnga", syllables: "le.hawng",
			root: 0, stress: 1,
			newSyllables: "",
			newStress:    -1,
		},
		{
			word: "fnemo-o", syllables: "fne.mo.o",
			root: 1, stress: 1,
			newSyllables: "fne.mo.-o",
			newStress:    1,
		},
		{
			word: "zekwä-äo", syllables: "zek.wä.ä.o",
			root: 0, stress: 0,
			newSyllables: "zek.wä.-ä.o",
			newStress:    0,
		},
		{
			word: "tsukanom", syllables: "tsuk.ka.nom",
			root: 0, stress: 1,
			newSyllables: "tsu.ka.nom",
			newStress:    1,
		},
		{
			word: "tsukkanom", syllables: "tsuk.ka.nom",
			root: 0, stress: 1,
			newSyllables: "tsuk.ka.nom",
			newStress:    1,
		},
		{
			word: "TìfMeTok", syllables: "tì.fme.tok",
			root: 0, stress: 1,
			newSyllables: "Tì.fMe.Tok",
			newStress:    1,
		},
		{
			word: "zekwä-ÄO", syllables: "zek.wä.ä.o",
			root: 0, stress: 0,
			newSyllables: "zek.wä.-Ä.O",
			newStress:    0,
		},
		{
			word: "NOLUI", syllables: "no.lu.i",
			root: 0, stress: 1,
			newSyllables: "NO.LU.I",
			newStress:    1,
		},
		{
			word: "uVaNsYu", syllables: "u.van. si.yu",
			root: 0, stress: 1,
			newSyllables: "u.VaN.sYu",
			newStress:    1,
		},
		{
			word: "sMunGe", syllables: "sä.mu.nge",
			root: 0, stress: 1,
			newSyllables: "sMu.nGe",
			newStress:    0,
		},
		{
			word: "tsuKkan", syllables: "tsuk.kan",
			root: 0, stress: 1,
			newSyllables: "tsuK.kan",
			newStress:    1,
		},
		{
			word: "tsuKKan", syllables: "tsuk.kan",
			root: 0, stress: 1,
			newSyllables: "tsuK.Kan",
			newStress:    1,
		},
		{
			word: "tsuKan", syllables: "tsuk.kan",
			root: 0, stress: 1,
			newSyllables: "tsu.Kan",
			newStress:    1,
		},
		{
			word: "tìftiä", syllables: "tì.fti.a.ä",
			root: 0, stress: 2,
			newSyllables: "tì.fti.ä",
			newStress:    2,
		},
		{
			word: "peyä", syllables: "po.yä",
			root: 0, stress: 0,
			newSyllables: "pe.yä",
			newStress:    0,
		},
		{
			word: "awngeyä", syllables: "aw.nga.yä",
			root: 0, stress: 1,
			newSyllables: "aw.nge.yä",
			newStress:    1,
		},
		{
			word: "tsata", syllables: "tsaw.ta",
			root: 0, stress: 0,
			newSyllables: "tsa.ta",
			newStress:    0,
		},
		{
			word: "sengi", syllables: "sä.ngi",
			root: 0, stress: 1,
			newSyllables: "se.ngi",
			newStress:    1,
		},
		{
			word: "txùkx", syllables: "txukx",
			root: 0, stress: 0,
			newSyllables: "txùkx",
			newStress:    0,
		},
		{
			word: "dukx", syllables: "txukx",
			root: 0, stress: 0,
			newSyllables: "dukx",
			newStress:    0,
		},
		{
			word: "dùkx", syllables: "txukx",
			root: 0, stress: 0,
			newSyllables: "dùkx",
			newStress:    0,
		},
		{
			word: "awgìl", syllables: "aw.kxìl",
			root: 0, stress: 0,
			newSyllables: "aw.gìl",
			newStress:    0,
		},
		{
			word: "tìkan-gan", syllables: "tì.kan.kxan",
			root: 0, stress: 0,
			newSyllables: "tì.kan.-gan",
			newStress:    0,
		},
		{
			word: "ayskxe", syllables: "ay.tskxe",
			root: 1, stress: 1,
			newSyllables: "ay.skxe",
			newStress:    1,
		},
		{
			word: "tsyìmawnun'i", syllables: "tì.syì.maw.nun.'i",
			root: 1, stress: 4,
			newSyllables: "tsyì.maw.nun.'i",
			newStress:    3,
		},
		{
			word: "'Ekongä", syllables: "'e.ko.ngä",
			root: 0, stress: 0,
			newSyllables: "'E.ko.ngä",
			newStress:    0,
		},
		{
			word: "Epxangteri", syllables: "e.pxang.te.ri",
			root: 0, stress: 1,
			newSyllables: "E.pxang.te.ri",
			newStress:    1,
		},
		{
			word: "Oengal", syllables: "o.eng.l",
			root: 1, stress: 1,
			newSyllables: "O.e.ngal",
			newStress:    1,
		},
		{
			word: "Oengeyä", syllables: "o.eng.yä",
			root: 1, stress: 2,
			newSyllables: "O.e.nge.yä",
			newStress:    3,
		},
		{
			word: "'Ekxongeyä", syllables: "'e.kxong.yä",
			root: 0, stress: 0,
			newSyllables: "'E.kxo.nge.yä",
			newStress:    0,
		},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("%s %s (%d, %d)", row.word, row.syllables, row.root, row.stress), func(t *testing.T) {
			newSyllables, newStress := MatchSyllables(row.word, strings.Split(row.syllables, "."), row.root, row.stress)

			assert.Equal(t, row.newSyllables, strings.Join(newSyllables, "."))
			assert.Equal(t, row.newStress, newStress)
		})
	}
}
