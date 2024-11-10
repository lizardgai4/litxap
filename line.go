package litxap

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func RunLine(line string, dictionary Dictionary, mustDouble map[string]string) (Line, error) {
	return ParseLine(line).Run(dictionary, mustDouble)
}

type Line []LinePart

func (line Line) Run(dict Dictionary, mustDouble map[string]string) (Line, error) {
	newLine := append(line[:0:0], line...)

	skip := false

	for i, part := range newLine {
		if !part.IsWord {
			continue
		}

		// If we found a multiword word, don't duplicate
		if skip {
			skip = false
			continue
		}

		lookup := part.Raw
		if part.Lookup != "" {
			lookup = part.Lookup
		}

		lookup = strings.ToLower(lookup)

		// Collect multiword words with no parts that can be looked up
		maybeSkip := false

		if _, ok := mustDouble[lookup]; ok {
			if i+2 < len(newLine) {
				lookup += strings.ToLower(newLine[i+1].Raw + newLine[i+2].Raw)
				maybeSkip = true
			}
		}

		results, err := dict.LookupEntries(lookup)
		if err != nil {
			if errors.Is(err, ErrEntryNotFound) {
				continue
			}

			return nil, fmt.Errorf("failed to lookup \"%s\": %w", lookup, err)
		} else if maybeSkip {
			// If we found a possible multiword word and it matches, skip the next one
			skip = true
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
