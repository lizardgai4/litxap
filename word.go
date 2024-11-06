package litxap

import (
	"github.com/gissleh/litxap/litxaputil"
)

func RunWord(word string, entry Entry) ([]string, int) {
	syllables, stress, root := entry.GenerateSyllables()
	return litxaputil.MatchSyllables(word, syllables, root, stress)
}
