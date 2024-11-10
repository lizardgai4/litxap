package litxaputil

import "strings"

/*
* Slightly modified from
* https://github.com/fwew/fwew-lib/blob/b17c90d83c91a39c6fb1daf9fccc9d91e02c097d/cache.go#L286
 */

/* To help deduce phonemes */
var romanization2 = map[string]string{
	// Vowels
	"a": "a", "i": "i", "ɪ": "ì",
	"o": "o", "ɛ": "e", "u": "u",
	"æ": "ä", "õ": "õ", //võvä' only
	// Diphthongs
	"aw": "aw", "ɛj": "ey",
	"aj": "ay", "ɛw": "ew",
	// Psuedovowels
	"ṛ": "rr", "ḷ": "ll",
	// Consonents
	"t": "t", "p": "p", "ʔ": "'",
	"n": "n", "k": "k", "l": "l",
	"s": "s", "ɾ": "r", "j": "y",
	"t͡s": "ts", "t'": "tx", "m": "m",
	"v": "v", "w": "w", "h": "h",
	"ŋ": "ng", "z": "z", "k'": "kx",
	"p'": "px", "f": "f", "r": "r",
	// Reef dialect
	"b": "b", "d": "d", "g": "g",
	"ʃ": "sh", "tʃ": "ch", "ʊ": "ù",
	// mistakes and rarities
	"ʒ": "ch", "": "", " ": ""}

// What is the nth rune of word?
func nth_rune(word string, n int) string {
	i := 0
	for _, r := range word {
		if i == n {
			return string(r)
		}
		i += 1
	}

	return ""
}

// Does ipa contain any character from word as its nth letter?
func has(word string, ipa string, n int) (output bool) {
	i := 0
	for _, s := range ipa {
		if i == n {
			for _, r := range word {
				if r == s {
					return true
				}
			}
			break // save a few compute cycles
		}
		i += 1
	}

	return false
}

// Helper function to get phonetic transcriptions of secondary pronunciations
func RomanizeIPA(IPA string) ([][][]string, [][]int) {
	// Special case: empty string
	if len(IPA) < 1 {
		return [][][]string{}, [][]int{}
	}

	stressMarkers := [][]int{}
	// now Romanize the IPA
	IPA = strings.ReplaceAll(IPA, "ʊ", "u")
	IPA = strings.ReplaceAll(IPA, "õ", "o") // vonvä' as võvä' only
	word := strings.Split(IPA, " ")

	results := [][]string{{}}

	bigResults := [][][]string{}

	// For "] or [" in the IPA
	if len(word) > 2 {
		word[0] = strings.ReplaceAll(word[0], "]", "")
		word[2] = strings.ReplaceAll(word[2], "[", "")

		// Make sure it's not the same word with different stresses
		/*if strings.ReplaceAll(word[0], "ˈ", "") == strings.ReplaceAll(word[2], "ˈ", "") {
			word = []string{strings.ReplaceAll(word[0], "ˈ", "")}
		}*/
	}

	stressMarkers = append(stressMarkers, []int{})

	// get the last one only
	for j := 0; j < len(word); j++ {
		// "or" means there's more than one IPA in this word, and we only want one
		if word[j] == "or" {
			bigResults = append(bigResults, results)
			results = [][]string{{}}
			stressMarkers = append(stressMarkers, []int{})
			continue
		}

		// In case of empty string
		if len(word[j]) < 1 {
			continue
		}

		stressMarkers[len(stressMarkers)-1] = append(stressMarkers[len(stressMarkers)-1], -1)

		syllables := strings.Split(word[j], ".")

		everStressed := false

		/* Onset */
		for k := 0; k < len(syllables); k++ {
			breakdown := ""

			stressed := strings.Contains(syllables[k], "ˈ")

			syllable := strings.ReplaceAll(syllables[k], "·", "")
			syllable = strings.ReplaceAll(syllable, "ˈ", "")
			syllable = strings.ReplaceAll(syllable, "ˌ", "")

			if stressed && !everStressed {
				everStressed = true
				stressMarkers[len(stressMarkers)-1][len(stressMarkers[len(stressMarkers)-1])-1] = k
			}

			// tsy
			if strings.HasPrefix(syllable, "tʃ") {
				breakdown += "ch"
				syllable = strings.TrimPrefix(syllable, "tʃ")
			} else if len(syllable) >= 4 && syllable[0:4] == "t͡s" {
				// ts
				breakdown += "ts"
				//tsp
				if has("ptk", syllable, 3) {
					if nth_rune(syllable, 4) == "'" {
						// ts + ejective onset
						breakdown += romanization2[syllable[4:6]]
						syllable = syllable[6:]
					} else {
						// ts + unvoiced plosive
						breakdown += romanization2[string(syllable[4])]
						syllable = syllable[5:]
					}
				} else if has("lɾmnŋwj", syllable, 3) {
					// ts + other consonent
					breakdown += romanization2[nth_rune(syllable, 3)]
					syllable = syllable[4+len(nth_rune(syllable, 3)):]
				} else {
					// ts without a cluster
					syllable = syllable[4:]
				}
			} else if has("fs", syllable, 0) {
				//
				breakdown += nth_rune(syllable, 0)
				if has("ptk", syllable, 1) {
					if nth_rune(syllable, 2) == "'" {
						// f/s + ejective onset
						breakdown += romanization2[syllable[1:3]]
						syllable = syllable[3:]
					} else {
						// f/s + unvoiced plosive
						breakdown += romanization2[string(syllable[1])]
						syllable = syllable[2:]
					}
				} else if has("lɾmnŋwj", syllable, 1) {
					// f/s + other consonent
					breakdown += romanization2[nth_rune(syllable, 1)]
					syllable = syllable[1+len(nth_rune(syllable, 1)):]
				} else {
					// f/s without a cluster
					syllable = syllable[1:]
				}
			} else if has("ptk", syllable, 0) {
				if nth_rune(syllable, 1) == "'" {
					// ejective
					breakdown += romanization2[syllable[0:2]]
					syllable = syllable[2:]
				} else {
					// unvoiced plosive
					breakdown += romanization2[string(syllable[0])]
					syllable = syllable[1:]
				}
			} else if has("ʔlɾhmnŋvwjzbdg", syllable, 0) {
				// other normal onset
				breakdown += romanization2[nth_rune(syllable, 0)]
				syllable = syllable[len(nth_rune(syllable, 0)):]
			} else if has("ʃʒ", syllable, 0) {
				// one sound representd as a cluster
				if nth_rune(syllable, 0) == "ʃ" {
					breakdown += "sh"
				}
				syllable = syllable[len(nth_rune(syllable, 0)):]
			}

			/*
			 * Nucleus
			 */
			psuedovowel := false
			if len(syllable) > 1 && has("jw", syllable, 1) {
				//diphthong
				breakdown += romanization2[syllable[0:len(nth_rune(syllable, 0))+1]]
				syllable = string([]rune(syllable)[2:])
			} else if len(syllable) > 1 && has("lr", syllable, 0) {
				// psuedovowel
				breakdown += romanization2[syllable[0:3]]
				psuedovowel = true
			} else {
				//vowel
				breakdown += romanization2[nth_rune(syllable, 0)]
				syllable = string([]rune(syllable)[1:])
			}

			/*
			 * Coda
			 */
			if !psuedovowel && len(syllable) > 0 {
				if nth_rune(syllable, 0) == "s" {
					breakdown += "sss" //oìsss only
				} else {
					if syllable == "k̚" {
						breakdown += "k"
					} else if syllable == "p̚" {
						breakdown += "p"
					} else if syllable == "t̚" {
						breakdown += "t"
					} else if syllable == "ʔ̚" {
						breakdown += "'"
					} else {
						if syllable[0] == 'k' && len(syllable) > 1 {
							breakdown += "kx"
						} else {
							breakdown += romanization2[syllable]
						}
					}
				}
			}

			results[len(results)-1] = append(results[len(results)-1], breakdown)
		}

		if j+1 < len(word) && word[j+1] != "or" {
			results = append(results, []string{})
		}
	}

	bigResults = append(bigResults, results)

	return bigResults, stressMarkers
}
