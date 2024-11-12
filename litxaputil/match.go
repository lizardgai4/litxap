package litxaputil

import (
	"strings"
)

func MatchSyllables(word string, syllables []string, root, stress int) (newSyllables []string, newStress int) {
	newSyllables, newStress = matchSyllables(word, syllables, root, stress, false)
	if newSyllables != nil {
		return
	}

	newSyllables, newStress = matchSyllables(word, syllables, root, stress, true)
	if newSyllables != nil {
		return
	}

	return
}

func matchSyllables(word string, syllables []string, root, stress int, allowFuse bool) (newSyllables []string, newStress int) {
	newSyllables = make([]string, 0, len(syllables))
	newStress = -1
	curr := word

	stressOffset := stress
	rootOffset := root

	for len(syllables) > 0 || len(curr) > 0 {
		matchedSyllables, next, n, stressPush := nextSyllable(curr, syllables, rootOffset >= 0, allowFuse)
		if n == 0 {
			newSyllables = nil
			newStress = -1
			break
		}

		newSyllables = append(newSyllables, matchedSyllables...)
		syllables = syllables[n:]
		curr = next

		rootOffset -= n
		if stressOffset >= 0 {
			stressOffset -= stressPush
			if stressOffset <= 0 {
				newStress = len(newSyllables) + stressOffset
			}
		}
	}

	if len(newSyllables) > 1 {
		for i, syllable := range newSyllables[:len(newSyllables)-1] {
			if strings.HasSuffix(syllable, "-") {
				newSyllables[i] = syllable[:len(syllable)-len("-")]
				newSyllables[i+1] = "-" + newSyllables[i+1]
			}
		}
	}

	return
}

func nextSyllable(curr string, syllables []string, allowLenition bool, allowFuse bool) ([]string, string, int, int) {
	if len(syllables) == 0 || len(curr) == 0 {
		return nil, curr, 0, 0
	}
	currLower := strings.ToLower(curr)

	// Spaces
	if strings.HasPrefix(curr, " ") {
		matchedSyllables, next, n, stressOffset := nextSyllable(strings.TrimLeft(curr, " "), syllables, allowLenition, allowFuse)
		if n > 0 {
			matchedSyllables = append([]string{" "}, matchedSyllables...)
		}

		return matchedSyllables, next, n, stressOffset
	}

	// Edge case: contracted k.k -> k
	if len(syllables) >= 2 && allowFuse {
		for _, fusables := range [][]string{fusableTails, fusableMids} {
			for _, fusable := range fusables {
				if strings.HasSuffix(syllables[0], fusable) && strings.HasPrefix(syllables[1], fusable) {
					if strings.HasPrefix(currLower, syllables[0][:len(syllables[0])-len(fusable)]+syllables[1]) {
						l1 := len(syllables[0]) - 1
						l2 := l1 + len(syllables[1])

						return []string{curr[:l1], curr[l1:l2]}, curr[l2:], 2, 2
					}
				}
			}
		}
	}

	// Edge case: soaia + ä -> soaiä
	if len(syllables) == 3 && strings.HasSuffix(syllables[0], "i") && syllables[1] == "a" && syllables[2] == "ä" {
		if strings.HasPrefix(currLower, syllables[0]+"ä") {
			l0 := len(syllables[0])
			l2 := len(syllables[2])
			return []string{curr[:l0], curr[l0 : l0+l2]}, curr[l0+l2:], 3, 2
		}
	}

	if strings.ContainsRune(currLower, '-') && !strings.HasSuffix(syllables[0], "-") {
		prev0 := syllables[0]
		syllables[0] = syllables[0] + "-"

		matchedSyllables, next, n, n2 := nextSyllable(curr, syllables, allowLenition, allowFuse)
		syllables[0] = prev0
		if n > 0 {
			matchedSyllables[0] += "-"
			return matchedSyllables, next, n, n2
		}
	}

	// Edge case: Xng.l -> X.ngal
	if len(syllables) >= 2 && strings.HasSuffix(syllables[0], "ng") {
		l0 := len(syllables[0])
		l1 := len(syllables[1])
		lng := len("ng")

		if strings.HasPrefix(currLower, syllables[0][:l0]+"a"+syllables[1]) {
			return []string{curr[:l0-lng], "nga" + curr[l0+len("a"):l0+len("a")+l1]}, curr[l0+l1+1:], 2, 2
		}
	}

	// Edge case: Xng.yä -> X.nge.yä
	if len(syllables) >= 2 && strings.HasSuffix(syllables[0], "ng") {
		l0 := len(syllables[0])
		l1 := len(syllables[1])
		lng := len("ng")

		if strings.HasPrefix(currLower, syllables[0][:l0]+"e"+syllables[1]) {
			return []string{curr[:l0-lng], "nge", curr[l0+len("a") : l0+len("a")+l1]}, curr[l0+l1+1:], 2, 3
		}
	}

	// Exact fit
	if strings.HasPrefix(curr, syllables[0]) {
		return syllables[:1], curr[len(syllables[0]):], 1, 1
	} else if strings.HasPrefix(currLower, syllables[0]) {
		return []string{curr[:len(syllables[0])]}, curr[len(syllables[0]):], 1, 1
	}

	// Try with removed space
	if strings.HasPrefix(syllables[0], " ") {
		syllables = append(syllables[:0:0], syllables...)
		syllables[0] = syllables[0][1:]

		if matchedSyllables, next, n, n2 := nextSyllable(curr, syllables, allowLenition, allowFuse); n > 0 {
			return matchedSyllables, next, n, n2
		}
	}

	// Check lenition if permitted
	if allowLenition {
		if _, lenitedSyllable := ApplyLenition(syllables[0]); strings.HasPrefix(currLower, lenitedSyllable) {
			return []string{curr[:len(lenitedSyllable)]}, curr[len(lenitedSyllable):], 1, 1
		}
	}

	// Edge case: ä becoming e
	if syllable := strings.ReplaceAll(syllables[0], "ä", "e"); syllable != syllables[0] && strings.HasPrefix(currLower, syllable) {
		return []string{curr[:len(syllable)]}, curr[len(syllable):], 1, 1
	}

	// Reef Na'vi: gdb (dict entries showing as kx,tx,px)
	withEjectives, hasInitialEjective := swapInitialEjective(syllables[0])
	if hasInitialEjective && strings.HasPrefix(currLower, withEjectives) {
		return []string{curr[:len(withEjectives)]}, curr[len(withEjectives):], 1, 1
	}

	// Reef Na'vi: ù (dict entries showing as u)
	for _, syllable := range [2]string{syllables[0], withEjectives} {
		if syllable := strings.ReplaceAll(syllable, "u", "ù"); syllable != syllables[0] && strings.HasPrefix(currLower, syllable) {
			return []string{curr[:len(syllable)]}, curr[len(syllable):], 1, 1
		}
	}

	// Edge case: contracted sä-X -> sX
	if len(syllables) >= 2 && syllables[0] == "sä" && strings.HasPrefix(currLower, "s"+syllables[1]) {
		return []string{curr[:len("s")+len(syllables[1])]}, curr[1+len(syllables[1]):], 2, 2
	}

	// Edge case: contracted tì-sX -> tsX
	if len(syllables) >= 2 && syllables[0] == "tì" && strings.HasPrefix(currLower, "t"+syllables[1]) {
		return []string{curr[:len("t")+len(syllables[1])]}, curr[1+len(syllables[1]):], 2, 2
	}

	// Edge case: contracted si-yu -> syu
	if len(syllables) >= 2 && syllables[0] == "si" && syllables[1] == "yu" && strings.HasPrefix(currLower, "syu") {
		return []string{curr[:len("syu")]}, curr[len("syu"):], 2, 2
	}

	// Edge case: po.yä -> pe.yä
	if len(syllables) == 2 && (strings.HasSuffix(syllables[0], "a") || strings.HasSuffix(syllables[0], "o")) && (syllables[1] == "yä" || syllables[1] == "ye") {
		l0 := len(syllables[0])
		l1 := len(syllables[1])

		if strings.HasPrefix(currLower, syllables[0][:l0-1]+"e"+syllables[1]) {
			return []string{curr[:l0], curr[l0 : l0+l1]}, curr[l0+l1:], 2, 2
		}
	}

	// Edge case: tsaw.ta -> tsa.ta
	if len(syllables) >= 2 && strings.HasSuffix(syllables[0], "aw") {
		l0 := len(syllables[0]) - len("w")
		l1 := len(syllables[1])
		check := strings.ToLower(syllables[0][:l0] + syllables[1])

		if strings.HasPrefix(currLower, check) {
			return []string{curr[:l0], curr[l0 : l0+l1]}, curr[l0+l1:], 2, 2
		}
	}

	// Failed
	return nil, curr, 0, 0
}

func swapInitialEjective(s string) (string, bool) {
	for i, ejective := range ejectives {
		if strings.HasPrefix(s, ejective) {
			return ejectiveAlts[i] + s[len(ejective):], true
		}
	}

	return s, false
}

var ejectives = []string{"px", "tx", "kx"}
var ejectiveAlts = []string{"b", "d", "g"}

var fusableTails = []string{"px", "tx", "kx", "m", "n", "l", "r", "p", "t", "k"}
var fusableMids = []string{"a", "ä", "e", "i", "ì", "o", "u", "ù"}
