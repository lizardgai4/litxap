package litxaputil

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

/*
 * Slightly modified from
 * https://github.com/fwew/fwew-lib/blob/b17c90d83c91a39c6fb1daf9fccc9d91e02c097d/cache.go#L286
 */

// RomanizeIPA generates a spelling based on the phonetics.
// The returned value is a list of words of syllables and a list of words' stresses.
func RomanizeIPA(IPA string) ([][][]string, [][]int) {
	// Special case: empty string
	if len(strings.Trim(IPA, " []")) < 1 {
		return [][][]string{}, [][]int{}
	}

	stressMarkers := make([][]int, 0, 2)

	// now Romanize the IPA
	IPA = strings.ReplaceAll(IPA, "ʊ", "u")
	IPA = strings.ReplaceAll(IPA, "õ", "o") // vonvä' as võvä' only
	word := strings.Split(IPA, " ")

	results := [][]string{{}}
	bigResults := make([][][]string, 0, 2)

	// For "] or [" in the IPA
	if len(word) > 2 {
		word[0] = strings.ReplaceAll(word[0], "]", "")
		word[2] = strings.ReplaceAll(word[2], "[", "")
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

		stressMarkers[len(stressMarkers)-1] = append(stressMarkers[len(stressMarkers)-1], -1)

		syllables := strings.Split(word[j], ".")

		everStressed := false
		breakdown := bytes.NewBuffer(make([]byte, 0, 16))

		/* Onset */
		for k := 0; k < len(syllables); k++ {
			breakdown.Reset()

			stressed := strings.Contains(syllables[k], "ˈ")

			syllable := strings.TrimPrefix(syllables[k], "·")
			syllable = strings.TrimPrefix(syllable, "ˈ")
			syllable = strings.TrimPrefix(syllable, "ˌ")

			if stressed && !everStressed {
				everStressed = true
				stressMarkers[len(stressMarkers)-1][len(stressMarkers[len(stressMarkers)-1])-1] = k
			}

			r1, r1s := utf8.DecodeRuneInString(syllable)

			foundCluster := false
			for _, clusterable := range [][]string{{"t͡s", "ts"}, {"f", "f"}, {"s", "s"}} {
				onset := clusterable[0]

				if strings.HasPrefix(syllable, clusterable[0]) {
					syllableAfterOnset := syllable[len(clusterable[0]):]

					breakdown.WriteString(clusterable[1])

					r2, r2s := utf8.DecodeRuneInString(syllable[len(onset):])
					r3, r3s := utf8.DecodeRuneInString(syllable[len(onset)+r2s:])

					if strings.ContainsRune("ptk", r2) { // ts/s/f + plosive
						if r3 == '\'' { // ejective
							breakdown.WriteString(romanizaionTable[syllableAfterOnset[:r2s+r3s]])
							syllable = syllableAfterOnset[r2s+r3s:]
						} else { // non-ejective
							breakdown.WriteString(romanizaionTable[syllableAfterOnset[:r2s]])
							syllable = syllableAfterOnset[r2s:]
						}
					} else if strings.ContainsRune("lɾmnŋwj", r2) { // ts/s/f + other consonent
						breakdown.WriteString(romanizaionTable[syllableAfterOnset[:r2s]])
						syllable = syllableAfterOnset[r2s:]
					} else {
						syllable = syllableAfterOnset
					}

					foundCluster = true
					break
				}
			}

			if !foundCluster {
				if strings.HasPrefix(syllable, "tʃ") { // tsy
					breakdown.WriteString("ch")
					syllable = strings.TrimPrefix(syllable, "tʃ")
				} else if strings.ContainsRune("ptk", r1) {
					if strings.HasPrefix(syllable[r1s:], "'") { // ejective
						breakdown.WriteString(romanizaionTable[syllable[0:2]])
						syllable = syllable[2:]
					} else { // unvoiced plosive
						breakdown.WriteString(romanizaionTable[syllable[:r1s]])
						syllable = syllable[1:]
					}
				} else if strings.ContainsRune("ʔlɾhmnŋvwjzbdg", r1) {
					// other consonant offset
					breakdown.WriteString(romanizaionTable[syllable[:r1s]])
					syllable = syllable[r1s:]
				} else if strings.ContainsRune("ʃʒ", r1) {
					// one sound representd as a cluster
					if r1 == 'ʃ' {
						breakdown.WriteString("sh")
					}

					syllable = syllable[r1s:]
				}
			}

			/*
			 * Nucleus
			 */
			psuedovowel := false
			syllable = strings.TrimPrefix(syllable, "·")
			r1, r1s = utf8.DecodeRuneInString(syllable)
			r2, r2s := utf8.DecodeRuneInString(syllable[r1s:])
			if strings.ContainsRune("jw", r2) { //diphthong
				breakdown.WriteString(romanizaionTable[syllable[:r1s+r2s]])
				syllable = syllable[r1s+r2s:]
			} else if r2 == '̣' { // psuedovowel
				breakdown.WriteString(romanizaionTable[syllable[:r1s+r2s]])
				psuedovowel = true
				syllable = syllable[r1s+r2s:]
			} else { // vowel
				breakdown.WriteString(romanizaionTable[syllable[:r1s]])
				syllable = syllable[r1s:]
			}

			/*
			 * Coda
			 */
			if !psuedovowel && len(syllable) > 0 {
				if strings.HasSuffix(syllable, "sss") {
					breakdown.WriteString("sss") //oìsss only
				} else if syllable == "k̚" {
					breakdown.WriteRune('k')
				} else if syllable == "p̚" {
					breakdown.WriteRune('p')
				} else if syllable == "t̚" {
					breakdown.WriteRune('t')
				} else {
					if syllable[0] == 'k' && len(syllable) > 1 {
						breakdown.WriteString("kx")
					} else {
						breakdown.WriteString(romanizaionTable[syllable])
					}
				}
			}

			results[len(results)-1] = append(results[len(results)-1], breakdown.String())
		}

		if j+1 < len(word) && word[j+1] != "or" {
			results = append(results, []string{})
		}
	}

	bigResults = append(bigResults, results)

	return bigResults, stressMarkers
}

var romanizaionTable = map[string]string{
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
	"ʒ": "ch", "": "", " ": "",
}
