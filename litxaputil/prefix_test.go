package litxaputil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestApplyPrefixes(t *testing.T) {
	table := []struct {
		curr           string
		prefixes       string
		expected       string
		expectedOffset int
	}{
		{
			curr: "ta.ron", prefixes: "fne",
			expected: "fne.ta.ron", expectedOffset: 1,
		},
		{
			curr: "ha.haw", prefixes: "tsuk",
			expected: "tsuk.ha.haw", expectedOffset: 1,
		},
		{
			curr: "i.nan", prefixes: "tsuk",
			expected: "tsu.ki.nan", expectedOffset: 1,
		},
		{
			curr: "eyk", prefixes: "ketsuk",
			expected: "ke.tsu.keyk", expectedOffset: 2,
		},
		{
			curr: "tì.fme.tok", prefixes: "pe,pxe,fne",
			expected: "pe.pxe.fne.tì.fme.tok", expectedOffset: 3,
		},
		{
			curr: "tskxe", prefixes: "ay",
			expected: "ay.tskxe", expectedOffset: 1,
		},
		{
			curr: "o.e", prefixes: "ay",
			expected: "a.yo.e", expectedOffset: 1,
		},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("%s- %s", row.prefixes, row.curr), func(t *testing.T) {
			curr := strings.Split(row.curr, ".")
			prefixes := strings.Split(row.prefixes, ",")
			next, nextOffset := ApplyPrefixes(curr, prefixes)

			assert.Equal(t, row.expected, strings.Join(next, "."))
			assert.Equal(t, row.expectedOffset, nextOffset)
		})
	}
}
