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

func TestMultiDictionary_LookupEntries(t *testing.T) {
	mdGood := MultiDictionary{
		dummyDictionary,
		DummyDictionary{
			"kameie":   *ParseEntry("k·a.m·e: <ei>: see into, understand"),
			"sa'nokur": *ParseEntry("sa'.nok: -ur: nother"),
		},
	}
	mdEmpty := MultiDictionary{}
	mdBad := MultiDictionary{dummyDictionary, BrokenDictionary{}}
	mdBad2 := MultiDictionary{BrokenDictionary{}}

	res, err := mdEmpty.LookupEntries("sa'nokur")
	assert.ErrorIs(t, err, ErrEntryNotFound)
	assert.Nil(t, res)

	res, err = mdBad.LookupEntries("sa'nokur")
	assert.Error(t, err)
	assert.Nil(t, res)

	res, err = mdBad2.LookupEntries("tìfmetok")
	assert.Error(t, err)
	assert.Nil(t, res)

	res, err = mdGood.LookupEntries("mìfa")
	assert.ErrorIs(t, err, ErrEntryNotFound)
	assert.Nil(t, res)

	res, err = mdGood.LookupEntries("kameie")
	assert.NoError(t, err)
	assert.Equal(t, res, []Entry{
		*ParseEntry("k·a.m·e: <ei>: see, see into, understand, know (spiritual sense)"),
		*ParseEntry("k··ä: <am,ei>: go"),
		*ParseEntry("k·a.m·e: <ei>: see into, understand"),
	})

	res, err = mdGood.LookupEntries("sa'nokur")
	assert.NoError(t, err)
	assert.Equal(t, res, []Entry{*ParseEntry("sa'.nok: -ur: nother")})
}
