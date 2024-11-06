package litxaputil

import (
	"strings"
)

// InfixPositionsFromBrackets finds the infix positions from the syllables and unsplit verb with all infix brackets.
// It can handle a missing "<0>" infix in case these get removed from words like zeyko in the future. If there are
// no infix positions, it will return `nil`.
//
// The function expect them to be in order, and that we'll never have just one of "<1>" and "<2>".
func InfixPositionsFromBrackets(infixStr string, syllables []string) *[2][2]int {
	infixStr = strings.Replace(infixStr, "<0>", "", 1)
	index1 := strings.Index(infixStr, "<1>")
	index2 := strings.Index(infixStr, "<2>") - len("<1>")

	var infixPositions *[2][2]int

	for i, syllable := range syllables {
		if index1 >= 0 && index1 < len(syllable) {
			infixPositions = &[2][2]int{{i, index1}}
		}
		if index2 >= 0 && index2 < len(syllable) {
			infixPositions[1] = [2]int{i, index2}
			break
		}

		index1 -= len(syllable)
		index2 -= len(syllable)
	}

	return infixPositions
}
