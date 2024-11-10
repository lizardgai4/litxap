package litxap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var dummyDictionary = DummyDictionary{
	"kaltxì":    *ParseEntry("kal.*txì"),
	"ma":        *ParseEntry("ma"),
	"fmetokyu":  *ParseEntry("fme.tok: -yu"),
	"ayhapxìtu": *ParseEntry("ha.*pxì.tu: ay-"),
	"soaiä":     *ParseEntry("so.*a.i.a: -ä"),
	"ngeyä":     *ParseEntry("nga: -yä"),
	"lu":        *ParseEntry("lu"),
	"oeru":      *ParseEntry("o.e: -ru"),
	"let'eylan": *ParseEntry("let.*'ey.lan"),
	"nìwotx":    *ParseEntry("nì.*wotx"),
	"oel":       *ParseEntry("o.e: -l"),
	"ngati":     *ParseEntry("nga: -ti"),
	"kameie":    *ParseEntry("k·a.m·e: <ei>: see, see into, understand, know (spiritual sense)"),
	"kameie:0":  *ParseEntry("k··ä: <am,ei>: go"),
	"säkeynven": *ParseEntry("sä.keyn.*ven"),
	"vola":      *ParseEntry("vol: -a"),
}

func TestRunLine(t *testing.T) {
	table := []struct {
		input    string
		expected Line
	}{
		{
			input: "Kaltxì, ma fmetokyu!",
			expected: Line{
				LinePart{Raw: "Kaltxì", IsWord: true, Matches: []LinePartMatch{
					{[]string{"Kal", "txì"}, 1, dummyDictionary["kaltxì"]},
				}},
				LinePart{Raw: ", "},
				LinePart{Raw: "ma", IsWord: true, Matches: []LinePartMatch{
					{[]string{"ma"}, 0, dummyDictionary["ma"]},
				}},
				LinePart{Raw: " "},
				LinePart{Raw: "fmetokyu", IsWord: true, Matches: []LinePartMatch{
					{[]string{"fme", "tok", "yu"}, 0, dummyDictionary["fmetokyu"]},
				}},
				LinePart{Raw: "!"},
			},
		},
		{
			input: "Oel ngati kameie.",
			expected: Line{
				LinePart{Raw: "Oel", IsWord: true, Matches: []LinePartMatch{
					{[]string{"O", "el"}, 0, dummyDictionary["oel"]},
				}},
				LinePart{Raw: " "},
				LinePart{Raw: "ngati", IsWord: true, Matches: []LinePartMatch{
					{[]string{"nga", "ti"}, 0, dummyDictionary["ngati"]},
				}},
				LinePart{Raw: " "},
				LinePart{Raw: "kameie", IsWord: true, Matches: []LinePartMatch{
					{[]string{"ka", "me", "i", "e"}, 0, dummyDictionary["kameie"]},
					{[]string{"ka", "me", "i", "e"}, 3, dummyDictionary["kameie:0"]},
				}},
				LinePart{Raw: "."},
			},
		},
		{
			input: "Ayhapxìtu soaiä ngeyä lu oeru let'eylan nìwotx.",
			expected: Line{
				LinePart{Raw: "Ayhapxìtu", IsWord: true, Matches: []LinePartMatch{
					{[]string{"Ay", "ha", "pxì", "tu"}, 2, dummyDictionary["ayhapxìtu"]},
				}},
				LinePart{Raw: " "},
				LinePart{Raw: "soaiä", IsWord: true, Matches: []LinePartMatch{
					{[]string{"so", "a", "i", "ä"}, 1, dummyDictionary["soaiä"]},
				}},
				LinePart{Raw: " "},
				LinePart{Raw: "ngeyä", IsWord: true, Matches: []LinePartMatch{
					{[]string{"nge", "yä"}, 0, dummyDictionary["ngeyä"]},
				}},
				LinePart{Raw: " "},
				LinePart{Raw: "lu", IsWord: true, Matches: []LinePartMatch{
					{[]string{"lu"}, 0, dummyDictionary["lu"]},
				}},
				LinePart{Raw: " "},
				LinePart{Raw: "oeru", IsWord: true, Matches: []LinePartMatch{
					{[]string{"o", "e", "ru"}, 0, dummyDictionary["oeru"]},
				}},
				LinePart{Raw: " "},
				LinePart{Raw: "let'eylan", IsWord: true, Matches: []LinePartMatch{
					{[]string{"let", "'ey", "lan"}, 1, dummyDictionary["let'eylan"]},
				}},
				LinePart{Raw: " "},
				LinePart{Raw: "nìwotx", IsWord: true, Matches: []LinePartMatch{
					{[]string{"nì", "wotx"}, 1, dummyDictionary["nìwotx"]},
				}},
				LinePart{Raw: "."},
			},
		},
		{
			input: "Vola skeynven.",
			expected: Line{
				LinePart{Raw: "Vola", IsWord: true, Matches: []LinePartMatch{
					{[]string{"Vo", "la"}, 0, dummyDictionary["vola"]},
				}},
				LinePart{Raw: " "},
				LinePart{Raw: "skeynven", IsWord: true},
				LinePart{Raw: "."},
			},
		},
		{
			input: "Vola säkeynven|skeynven.",
			expected: Line{
				LinePart{Raw: "Vola", IsWord: true, Matches: []LinePartMatch{
					{[]string{"Vo", "la"}, 0, dummyDictionary["vola"]},
				}},
				LinePart{Raw: " "},
				LinePart{Raw: "skeynven", Lookup: "säkeynven", IsWord: true, Matches: []LinePartMatch{
					{[]string{"skeyn", "ven"}, 1, dummyDictionary["säkeynven"]},
				}},
				LinePart{Raw: "."},
			},
		},
	}

	doubles := map[string]string{"tsaheyl": "si"}

	for _, row := range table {
		t.Run(row.input, func(t *testing.T) {
			res, err := RunLine(row.input, dummyDictionary, doubles)
			assert.NoError(t, err)
			assert.Equal(t, row.expected, res)
		})
	}
}

func TestRunLine_Fail(t *testing.T) {
	doubles := map[string]string{"tsaheyl": "si"}
	line, err := RunLine("Kaltxì, ma kifkey!", BrokenDictionary{}, doubles)

	assert.Error(t, err)
	assert.Nil(t, line)
	assert.NotErrorIs(t, err, ErrEntryNotFound)
}

func TestParseLine(t *testing.T) {
	table := []struct {
		input    string
		expected Line
	}{
		{
			input: "Ftuea tìfmetok",
			expected: Line{
				LinePart{Raw: "Ftuea", IsWord: true},
				LinePart{Raw: " "},
				LinePart{Raw: "tìfmetok", IsWord: true},
			},
		},
		{
			input: "spono-o aean-na-pay",
			expected: Line{
				LinePart{Raw: "spono-o", IsWord: true},
				LinePart{Raw: " "},
				LinePart{Raw: "aean", IsWord: true},
				LinePart{Raw: "-"},
				LinePart{Raw: "na", IsWord: true},
				LinePart{Raw: "-"},
				LinePart{Raw: "pay", IsWord: true},
			},
		},
		{
			input: "Ngäzìka tìkenong-o",
			expected: Line{
				LinePart{Raw: "Ngäzìka", IsWord: true},
				LinePart{Raw: " "},
				LinePart{Raw: "tìkenong-o", IsWord: true},
			},
		},
		{
			input: "Fìtìfmetok lu nì'it ngäzìk to pum aham.",
			expected: Line{
				LinePart{Raw: "Fìtìfmetok", IsWord: true},
				LinePart{Raw: " "},
				LinePart{Raw: "lu", IsWord: true},
				LinePart{Raw: " "},
				LinePart{Raw: "nì'it", IsWord: true},
				LinePart{Raw: " "},
				LinePart{Raw: "ngäzìk", IsWord: true},
				LinePart{Raw: " "},
				LinePart{Raw: "to", IsWord: true},
				LinePart{Raw: " "},
				LinePart{Raw: "pum", IsWord: true},
				LinePart{Raw: " "},
				LinePart{Raw: "aham", IsWord: true},
				LinePart{Raw: "."},
			},
		},
		{
			input: "'Awa säkeynven|skeynven angim",
			expected: Line{
				LinePart{Raw: "'Awa", IsWord: true},
				LinePart{Raw: " "},
				LinePart{Raw: "skeynven", Lookup: "säkeynven", IsWord: true},
				LinePart{Raw: " "},
				LinePart{Raw: "angim", IsWord: true},
			},
		},
	}

	for _, row := range table {
		t.Run(row.input, func(t *testing.T) {
			assert.Equal(t, row.expected, ParseLine(row.input))
		})
	}
}
