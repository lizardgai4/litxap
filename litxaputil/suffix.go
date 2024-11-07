package litxaputil

import (
	"fmt"
	"strings"
)

func suffix(ra int, s ...string) Suffix {
	return Suffix{reanalysis: ra, syllableSplit: s}
}

type Suffix struct {
	// Trigger reanalysis, e.g. awkx + -ìl => aw.kxìl
	reanalysis int
	// syllableSplit describes how the suffix will be added.
	// first: will end the last syllable
	// non-first: always its own syllable
	syllableSplit []string
}

func (suffix Suffix) Apply(curr []string) []string {
	if len(curr) == 0 {
		panic("Suffix attempted attached to empty word")
	}

	switch suffix.reanalysis {
	case sraNewSyllable:
		return append(curr, suffix.syllableSplit...)
	case sraAttach:
		lastSyllable := curr[len(curr)-1]
		for _, core := range attachableCores {
			if strings.HasSuffix(lastSyllable, core) {
				curr[len(curr)-1] = lastSyllable + suffix.syllableSplit[0]
				return append(curr, suffix.syllableSplit[1:]...)
			}
		}

		return append(curr, suffix.syllableSplit...)
	case sraStealCoda:
		lastSyllable := curr[len(curr)-1]
		unStealableCoda := false
		for _, coda := range unStealableCodas {
			if strings.HasSuffix(lastSyllable, coda) {
				unStealableCoda = true
				break
			}
		}
		if !unStealableCoda {
			for _, coda := range stealableCodas {
				if strings.HasSuffix(lastSyllable, coda) {
					curr[len(curr)-1] = lastSyllable[:len(lastSyllable)-len(coda)]
					curr = append(curr, lastSyllable[len(lastSyllable)-len(coda):]+suffix.syllableSplit[0])
					return append(curr, suffix.syllableSplit[1:]...)
				}
			}
		}

		return append(curr, suffix.syllableSplit...)
	default:
		panic(fmt.Sprint("Invalid suffix reanalysis mode: ", suffix.reanalysis))
	}
}

// ApplySuffixes applies the suffixes to the syllable set. None of them change stress (yet), so the stress index
// shall remain the same before and after.
func ApplySuffixes(curr []string, suffixNames []string) []string {
	for _, suffixName := range suffixNames {
		suffix := findSuffix(suffixName)
		curr = suffix.Apply(curr)
	}

	return curr
}

func findSuffix(name string) Suffix {
	if suffix, ok := suffixMap[name]; ok {
		return suffix
	}

	return Suffix{reanalysis: sraNewSyllable, syllableSplit: []string{name}}
}

const (
	sraNewSyllable = 0
	sraAttach      = 1
	sraStealCoda   = 2
)

var unStealableCodas = []string{"ll", "rr"}
var stealableCodas = []string{"ng", "px", "tx", "kx", "p", "t", "k", "b", "d", "g", "'", "m", "n", "l", "r", "w", "y"}
var attachableCores = []string{"aw", "ay", "ew", "ey", "a", "ä", "e", "é", "i", "ì", "o", "u", "ù"}

var suffixMap = map[string]Suffix{
	"tswo":  suffix(sraNewSyllable, "tswo"),
	"tsyìp": suffix(sraNewSyllable, "tsyìp"),
	"yu":    suffix(sraNewSyllable, "yu"),
	"fkeyk": suffix(sraNewSyllable, "fkeyk"),
	"tseng": suffix(sraNewSyllable, "tseng"),

	"äo":  suffix(sraStealCoda, "ä", "o"),
	"eo":  suffix(sraStealCoda, "e", "o"),
	"io":  suffix(sraStealCoda, "i", "o"),
	"uo":  suffix(sraStealCoda, "u", "o"),
	"ìlä": suffix(sraStealCoda, "ì", "lä"),

	"mungwrr": suffix(sraNewSyllable, "mung", "wrr"),
	"teri":    suffix(sraNewSyllable, "te", "ri"),
	"kxamlä":  suffix(sraNewSyllable, "kxam", "lä"),
	"mìkam":   suffix(sraNewSyllable, "mì", "kam"),
	"nemfa":   suffix(sraNewSyllable, "nem", "fa"),
	"takip":   suffix(sraNewSyllable, "ta", "kip"),
	"luke":    suffix(sraNewSyllable, "lu", "ke"),
	"tafkip":  suffix(sraNewSyllable, "ta", "fkip"),
	"pxisre":  suffix(sraNewSyllable, "pxi", "sre"),
	"pximaw":  suffix(sraNewSyllable, "pxi", "maw"),
	"rofa":    suffix(sraNewSyllable, "ro", "fa"),
	"lisre":   suffix(sraNewSyllable, "li", "sre"),
	"nuä":     suffix(sraNewSyllable, "nu", "ä"),
	"talun":   suffix(sraNewSyllable, "ta", "lun"),
	"yoa":     suffix(sraNewSyllable, "yo", "a"),
	"krrka":   suffix(sraNewSyllable, "krr", "ka"),
	"ftumfa":  suffix(sraNewSyllable, "ftum", "fa"),
	"ftuopa":  suffix(sraNewSyllable, "ftu", "o", "pa"),

	"a": suffix(sraStealCoda, "a"),
	"o": suffix(sraStealCoda, "o"),

	"l":                suffix(sraAttach, "l"),
	"ìl":               suffix(sraStealCoda, "ìl"),
	"t":                suffix(sraAttach, "t"),
	"ti":               suffix(sraNewSyllable, "ti"),
	"it":               suffix(sraStealCoda, "it"),
	"ejectiveReplacer": suffix(sraAttach, "ejectiveReplacer"),
	"r":                suffix(sraAttach, "r"),
	"ru":               suffix(sraNewSyllable, "ru"),
	"ur":               suffix(sraStealCoda, "ur"),
	"yä":               suffix(sraNewSyllable, "yä"),
	"ye":               suffix(sraNewSyllable, "ye"),
	"y":                suffix(sraAttach, "y"),
	"ä":                suffix(sraStealCoda, "ä"),
	"e":                suffix(sraStealCoda, "e"),
	"ri":               suffix(sraNewSyllable, "ri"),
	"ìri":              suffix(sraStealCoda, "ì", "ri"),
}
