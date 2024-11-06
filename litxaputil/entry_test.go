package litxaputil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestInfixPositionsFromBrackets(t *testing.T) {
	table := []struct {
		inInfixes   string
		inSyllables string
		out         *[2][2]int
	}{
		{
			inInfixes: "tsl<0><1><2>ikx", inSyllables: "tslikx",
			out: &[2][2]int{{0, 3}, {0, 3}},
		},
		{
			inInfixes: "t<0><1>ìr<2>an", inSyllables: "tì.ran",
			out: &[2][2]int{{0, 1}, {1, 1}},
		},
		{
			inInfixes: "t<0><1><2>ul", inSyllables: "tul",
			out: &[2][2]int{{0, 1}, {0, 1}},
		},
		{
			inInfixes: "pxawt<0><1><2>ok", inSyllables: "pxaw.tok",
			out: &[2][2]int{{1, 1}, {1, 1}},
		},
		{
			inInfixes: "zeyk<1><2>o", inSyllables: "zey.ko",
			out: &[2][2]int{{1, 1}, {1, 1}},
		},
		{
			inInfixes: "txakrrfp<0><1><2>ìl", inSyllables: "txa.krr.fpìl",
			out: &[2][2]int{{2, 2}, {2, 2}},
		},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("%s %s", row.inInfixes, row.inSyllables), func(t *testing.T) {
			outi := InfixPositionsFromBrackets(row.inInfixes, strings.Split(row.inSyllables, "."))
			assert.Equal(t, row.out, outi)
		})
	}
}
