package litxap

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type DummyDictionary map[string]Entry

func (t DummyDictionary) LookupEntries(word string) ([]Entry, error) {
	word = strings.ToLower(word)

	res := make([]Entry, 0, 1)
	if entry, ok := t[word]; ok {
		res = append(res, entry)

		for i := 0; i < 10; i++ {
			if entry, ok := t[word+":"+fmt.Sprint(i)]; ok {
				res = append(res, entry)
			} else {
				break
			}
		}
	} else {
		return nil, ErrEntryNotFound
	}

	return res, nil
}

type BrokenDictionary struct{}

func (b BrokenDictionary) LookupEntries(_ string) ([]Entry, error) {
	return nil, errors.New("500 something something")
}

func TestParseEntry(t *testing.T) {
	table := []string{
		"tskxe",
		"lo.ran: pe-fne- -ìri",
		"u.*van: -ti",
		"t·ì.*r·an: <äpeyk,ol>: walk",
		"t··el: <ei>: get, receive",
		"t·a.r·on: tì- <us> -ti: hunt",
		"t·ì.*r·an: tì- <us> -ìri: walk",
		"sä.*pxor: : explosion",
	}

	for _, row := range table {
		t.Run(row, func(t *testing.T) {
			parsed := ParseEntry(row)
			assert.Equal(t, row, parsed.String())
		})
	}
}
