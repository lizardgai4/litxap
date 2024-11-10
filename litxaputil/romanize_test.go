package litxaputil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRomanize(t *testing.T) {
	table := []struct {
		curr     string
		expected [][][]string
		stress   [][]int
	}{
		// One syllable
		{"ɛ", [][][]string{{{"e"}}}, [][]int{{-1}}},
		{"ʔawk'", [][][]string{{{"'awkx"}}}, [][]int{{-1}}},
		{"fko", [][][]string{{{"fko"}}}, [][]int{{-1}}},
		{"mo", [][][]string{{{"mo"}}}, [][]int{{-1}}},
		{"t'on", [][][]string{{{"txon"}}}, [][]int{{-1}}},
		{"t͡sam", [][][]string{{{"tsam"}}}, [][]int{{-1}}},
		{"fpom", [][][]string{{{"fpom"}}}, [][]int{{-1}}},
		{"kṛ", [][][]string{{{"krr"}}}, [][]int{{-1}}},
		{"k'ḷ", [][][]string{{{"kxll"}}}, [][]int{{-1}}},
		{"po", [][][]string{{{"po"}}}, [][]int{{-1}}},
		{"sk'awŋ", [][][]string{{{"skxawng"}}}, [][]int{{-1}}},
		//Multi syllable
		{"tɪ.ˈfmɛ.tok̚", [][][]string{{{"tì", "fme", "tok"}}}, [][]int{{1}}},
		{"u.ˈvan", [][][]string{{{"u", "van"}}}, [][]int{{1}}},
		{"ˈu.ɾan", [][][]string{{{"u", "ran"}}}, [][]int{{0}}},
		{"ˈt·a.ɾ·on", [][][]string{{{"ta", "ron"}}}, [][]int{{0}}},
		{"ˈʔɛ.koŋ", [][][]string{{{"'e", "kong"}}}, [][]int{{0}}},
		{"ɛ.ˈjawɾ", [][][]string{{{"e", "yawr"}}}, [][]int{{1}}},
		// Flexible syllable stress
		{"aj.ˈfo] or [ˈaj.fo", [][][]string{{{"ay", "fo"}}, {{"ay", "fo"}}}, [][]int{{1}, {0}}},
		{"ˈɪ.læ] or [ɪ.ˈlæ", [][][]string{{{"ì", "lä"}}, {{"ì", "lä"}}}, [][]int{{0}, {1}}},
		{"ˈmɪ.fa] or [mɪ.ˈfa", [][][]string{{{"mì", "fa"}}, {{"mì", "fa"}}}, [][]int{{0}, {1}}},
		{"ˈt͡sa.kɛm] or [t͡sa.ˈkɛm", [][][]string{{{"tsa", "kem"}}, {{"tsa", "kem"}}}, [][]int{{0}, {1}}},
		{"t͡sa.ˈt͡sɛŋ] or [ˈt͡sa.t͡sɛŋ", [][][]string{{{"tsa", "tseng"}}, {{"tsa", "tseng"}}}, [][]int{{1}, {0}}},
		// Multiple pronunciation
		{"nɪ.aw.ˈno.mʊm] or [naw.ˈno.mʊm", [][][]string{{{"nì", "aw", "no", "mum"}}, {{"naw", "no", "mum"}}}, [][]int{{2}, {1}}},
		{"nɪ.aj.ˈwɛŋ] or [naj.ˈwɛŋ", [][][]string{{{"nì", "ay", "weng"}}, {{"nay", "weng"}}}, [][]int{{2}, {1}}},
		{"tɪ.sjɪ.maw.nʊn.ˈʔi] or [t͡sjɪ.maw.nʊn.ˈʔi", [][][]string{{{"tì", "syì", "maw", "nun", "'i"}}, {{"tsyì", "maw", "nun", "'i"}}}, [][]int{{4}, {3}}},
		{"tɪ.sæ.ˈfpɪl.jɛwn] or [t͡sæ.ˈfpɪl.jɛwn", [][][]string{{{"tì", "sä", "fpìl", "yewn"}}, {{"tsä", "fpìl", "yewn"}}}, [][]int{{2}, {1}}},
		// Multiple words
		{"ˈut.ɾa.ja ˈmok.ɾi", [][][]string{{{"ut", "ra", "ya"}, {"mok", "ri"}}}, [][]int{{0, 0}}},
		{"t͡sa.ˈhɛjl s·i", [][][]string{{{"tsa", "heyl"}, {"si"}}}, [][]int{{1, -1}}},
		{"ˈnɪ.ˌju ˈjoɾ.kɪ", [][][]string{{{"nì", "yu"}, {"yor", "kì"}}}, [][]int{{0, 0}}},
		{"t͡sawl sl·u", [][][]string{{{"tsawl"}, {"slu"}}}, [][]int{{-1, -1}}},
		// Empty string
		{"", [][][]string{}, [][]int{}},
	}

	for _, row := range table {
		t.Run(row.curr, func(t *testing.T) {
			spelling, stress := RomanizeIPA(row.curr)
			assert.Equal(t, row.expected, spelling)
			assert.Equal(t, row.stress, stress)
		})
	}
}
