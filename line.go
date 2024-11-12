package litxap

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func RunLine(line string, dictionary Dictionary) (Line, error) {
	return ParseLine(line).Run(dictionary)
}

type Line []LinePart

func (line Line) Run(dict Dictionary) (Line, error) {
	newLine := append(line[:0:0], line...)

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
			if errors.Is(err, ErrEntryNotFound) {
				continue
			}

			return nil, fmt.Errorf("failed to lookup \"%s\": %w", lookup, err)
		}

		for _, result := range results {
			syllables, stress := RunWord(part.Raw, result)
			if syllables != nil {
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

func (line Line) Merge(mustDouble map[string]string) Line {
	first := true
	newLine := line

	for {
		found := false
		for i, p0 := range newLine[:len(newLine)-2] {
			md, ok := mustDouble[strings.ToLower(p0.Raw)]
			if !ok {
				continue
			}

			p1 := newLine[i+1]
			p2 := newLine[i+2]

			if p0.IsWord && !p1.IsWord && p2.IsWord && md == strings.ToLower(p2.Raw) {
				if first == true {
					newLine = append(newLine[:0:0], newLine...)
					first = false
				}

				newLine = append(newLine[:i+1], newLine[i+3:]...)
				newLine[i].Raw = p0.Raw + p1.Raw + p2.Raw

				found = true
				break
			}
		}

		if !found {
			break
		}
	}

	return newLine
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
