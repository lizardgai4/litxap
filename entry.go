package litxap

import (
	"errors"
	"strings"

	"github.com/gissleh/litxap/litxaputil"
)

type Entry struct {
	// The Na'vi Word.
	Word string `json:"word"`
	// A Translation in the language, mostly to guide the user if multiple matches fit.
	Translation string `json:"translation"`
	// Syllables is a list of syllables
	Syllables []string `json:"syllables"`
	// Stress is the zero-based index of the stressed syllable.
	Stress int `json:"stress"`
	// InfixPos has a pair of positions, syllable, byte.
	InfixPos *[2][2]int `json:"infixPos,omitempty"`

	// Prefixes are an in-order list of prefixes
	Prefixes []string `json:"prefixes,omitempty"`
	// Infixes are a list of infixes, order does not matter
	Infixes []string `json:"infixes,omitempty"`
	// Suffixes are an in-order list of suffixes
	Suffixes []string `json:"suffixes,omitempty"`
}

func (entry *Entry) GenerateSyllables() ([]string, int, int) {
	syllables := append(entry.Syllables[:0:0], entry.Syllables...)
	stress := entry.Stress

	syllables, offset := litxaputil.ApplyPrefixes(syllables, entry.Prefixes)
	stress += offset

	if entry.InfixPos != nil && len(entry.Infixes) > 0 {
		positions := *entry.InfixPos
		positions[0][0] += offset
		positions[1][0] += offset

		syllables, stress = litxaputil.ApplyInfixes(syllables, entry.Infixes, offset, stress, positions)
	}

	syllables = litxaputil.ApplySuffixes(syllables, entry.Suffixes)

	return syllables, stress, offset
}

func (entry *Entry) String() string {
	sb := strings.Builder{}
	sb.Grow(64)

	for i, syllable := range entry.Syllables {
		if i > 0 {
			sb.WriteRune('.')
		}
		if entry.Stress == i && i != 0 {
			sb.WriteRune('*')
		}

		foundInfix := false
		if entry.InfixPos != nil {
			for j, infixPos := range *entry.InfixPos {
				if infixPos[0] == i {
					sb.WriteString(syllable[:infixPos[1]])
					sb.WriteRune('·')
					if j == 0 && (*entry.InfixPos)[1] == (*entry.InfixPos)[0] {
						sb.WriteRune('·')
					}
					sb.WriteString(syllable[infixPos[1]:])
					foundInfix = true
					break
				}
			}
		}
		if !foundInfix {
			sb.WriteString(syllable)
		}
	}

	inflected := len(entry.Prefixes) > 0 || len(entry.Suffixes) > 0 || len(entry.Infixes) > 0
	if inflected {
		sb.WriteByte(':')
	}

	if len(entry.Prefixes) > 0 {
		sb.WriteByte(' ')
		for _, prefix := range entry.Prefixes {
			sb.WriteString(prefix)
			sb.WriteByte('-')
		}
	}

	if len(entry.Infixes) > 0 {
		sb.WriteString(" <")
		for i, infix := range entry.Infixes {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(infix)
		}
		sb.WriteRune('>')
	}

	if len(entry.Suffixes) > 0 {
		sb.WriteByte(' ')
		for _, suffix := range entry.Suffixes {
			sb.WriteByte('-')
			sb.WriteString(suffix)
		}
	}

	if len(entry.Translation) > 0 {
		if !inflected {
			sb.WriteString(": ")
		}
		sb.WriteString(": ")
		sb.WriteString(entry.Translation)
	}

	return sb.String()
}

// ParseEntry is admittedly a test-utility, but it's kept out here as Entry is a top-level object.
// It parses the format that comes out Entry.String: e.g. "t·ì.*r·an: tì- <us> -ìri".
// In spite of the return type, it'll always give back an Entry.
func ParseEntry(s string) *Entry {
	split := strings.Split(s, ": ")
	entry := Entry{}

	for i, syllable := range strings.Split(split[0], ".") {
		if strings.HasPrefix(syllable, "*") {
			entry.Stress = i
			syllable = syllable[len("*"):]
		}

		dotIndex := strings.Index(syllable, "·")
		lastDotIndex := strings.LastIndex(syllable, ".")

		if dotIndex != -1 {
			if entry.InfixPos == nil {
				entry.InfixPos = &[2][2]int{{i, dotIndex}, {-1, -1}}
			} else {
				(*entry.InfixPos)[1] = [2]int{i, dotIndex}
			}

			if lastDotIndex != dotIndex {
				(*entry.InfixPos)[1] = [2]int{i, dotIndex}
			}

			syllable = strings.ReplaceAll(syllable, "·", "")
		}

		entry.Syllables = append(entry.Syllables, syllable)
		entry.Word += syllable
	}

	if len(split) > 1 {
		for _, token := range strings.Split(split[1], " ") {
			if strings.HasPrefix(token, "-") {
				entry.Suffixes = strings.Split(token[len("-"):], "-")
			}
			if strings.HasSuffix(token, "-") {
				entry.Prefixes = strings.Split(token[:len(token)-len("-")], "-")
			}
			if strings.HasPrefix(token, "<") && strings.HasSuffix(token, ">") {
				entry.Infixes = strings.Split(token[len("<"):len(token)-len(">")], ",")
			}
		}
	}

	if len(split) > 2 {
		entry.Translation = strings.Join(split[2:], ": ")
	}

	return &entry
}

type Dictionary interface {
	LookupEntries(word string) ([]Entry, error)
}

type MultiDictionary []Dictionary

func (dm MultiDictionary) LookupEntries(word string) ([]Entry, error) {
	if len(dm) == 0 {
		return nil, ErrEntryNotFound
	}

	entries, err := dm[0].LookupEntries(word)
	if err != nil && !errors.Is(err, ErrEntryNotFound) {
		return nil, err
	}

	for _, dict := range dm[1:] {
		next, err := dict.LookupEntries(word)
		if err != nil && !errors.Is(err, ErrEntryNotFound) {
			return nil, err
		}

		if len(next) > 0 {
			entries = append(entries, next...)
		}
	}

	if len(entries) == 0 {
		return nil, ErrEntryNotFound
	}

	return entries, nil
}

var ErrEntryNotFound = errors.New("entry not found")
