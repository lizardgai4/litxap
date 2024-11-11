package litxap

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	fwew_lib "github.com/fwew/fwew-lib/v5"
	"github.com/gissleh/litxap/litxaputil"
)

func RunLine(line string, dictionary Dictionary, mustDouble map[string]string) (Line, error) {
	return ParseLine(line).Run(dictionary, mustDouble)
}

type Line []LinePart

func (line Line) Run(dict Dictionary, mustDouble map[string]string) (Line, error) {
	newLine := append(line[:0:0], line...)

	/*for {
		found := false
		for i, p0 := range newLine[:len(newLine)-2] {
			md, ok := mustDouble[strings.ToLower(p0.Raw)]
			if !ok {
				continue
			}

			p1 := newLine[i+1]
			p2 := newLine[i+2]

			if p0.IsWord && !p1.IsWord && p2.IsWord && md == strings.ToLower(p2.Raw) {
				newLine = append(newLine[:i+1], newLine[i+3:]...)
				newLine[i].Raw = p0.Raw + p1.Raw + p2.Raw

				found = true
				break
			}
		}

		if !found {
			break
		}
	}*/

	for i, part := range newLine {
		if !part.IsWord {
			continue
		}

		lookup := part.Raw
		if part.Lookup != "" {
			lookup = part.Lookup
		}

		lookup = strings.ToLower(lookup)

		results, err := dict.LookupEntries(lookup)
		if err != nil {
			// See if it's tere
			entry, ok := mustDouble[strings.ToLower(lookup)]

			// If it's not there, try deconjugating
			if !ok {
				entries := fwew_lib.Deconjugate(lookup)
				for _, entry2 := range entries {
					if entry2.InsistPOS != "any" && entry2.InsistPOS != "n." {
						continue
					}
					entry3, ok2 := mustDouble[strings.ToLower(entry2.Word)]
					if ok2 {
						ok = true
						entry = entry3
						break
					}
				}
			}

			// If it's in either place, see the Romanization
			if ok {
				// Romanize and find stress from the IPA
				syllables, stress := litxaputil.RomanizeIPA(entry)
				if syllables != nil && stress[0][0] >= 0 {
					newLine[i].Matches = append(newLine[i].Matches, LinePartMatch{
						Syllables: syllables[0][0],
						Stress:    stress[0][0],
					})
				}
			}

			if errors.Is(err, ErrEntryNotFound) {
				continue
			}

			return nil, fmt.Errorf("failed to lookup \"%s\": %w", lookup, err)
		}

		for _, result := range results {
			syllables, stress := RunWord(part.Raw, result)
			if syllables != nil && stress >= 0 {
				newLine[i].Matches = append(newLine[i].Matches, LinePartMatch{
					Syllables: syllables,
					Stress:    stress,
					Entry:     result,
				})
			}
		}
	}

	return newLine, nil
}

// ParseLine splits out the words from a line of text.
func ParseLine(s string) Line {
	wordMode := false
	lastPos := 0
	lastPipe := 0
	currentPos := 0
	res := make(Line, 0, (len(s)/5)+1)

	s = strings.NewReplacer("’", "'", "‘", "'").Replace(s) + "\n"

	for _, ch := range s {
		if ch == '|' {
			lastPipe = currentPos
		} else if ch == '\n' || wordMode != (unicode.IsLetter(ch) || ch == '\'' || ch == '-') {
			if lastPos != currentPos {
				if wordMode && strings.Contains(s[lastPos:currentPos], "-na-") {
					split := strings.SplitN(s[lastPos:currentPos], "-na-", 2)
					res = append(res, LinePart{
						Raw: split[0], IsWord: true, Matches: nil,
					}, LinePart{
						Raw: "-",
					}, LinePart{
						Raw: "na", IsWord: true, Matches: nil,
					}, LinePart{
						Raw: "-",
					}, LinePart{
						Raw: split[1], IsWord: true, Matches: nil,
					})
				} else {
					raw := s[lastPos:currentPos]
					lookup := s[lastPos:lastPos]
					if lastPipe != lastPos {
						lookup = s[lastPos:lastPipe]
						raw = s[lastPipe+1 : currentPos]
					}

					res = append(res, LinePart{
						Raw:     raw,
						Lookup:  lookup,
						IsWord:  wordMode,
						Matches: nil,
					})
				}

				lastPos = currentPos
				lastPipe = currentPos
			}

			wordMode = !wordMode
		}

		currentPos += utf8.RuneLen(ch)
	}

	return res
}

type LinePart struct {
	Raw     string          `json:"raw"`
	Lookup  string          `json:"lookup,omitempty"`
	IsWord  bool            `json:"isWord,omitempty"`
	Matches []LinePartMatch `json:"matches,omitempty"`
}

type LinePartMatch struct {
	Syllables []string `json:"syllables"`
	Stress    int      `json:"stress"`
	Entry     Entry    `json:"entry"`
}
