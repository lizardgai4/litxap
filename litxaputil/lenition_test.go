package litxaputil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplyLenition(t *testing.T) {
	table := []struct {
		Input    string
		Lenition string
		Next     string
	}{
		{"", "", ""},
		{"kilvan", "k→h", "hilvan"},
		{"tsìlpey", "ts→s", "sìlpey"},
		{"tìtok", "t→s", "sìtok"},
		{"pasuk", "p→f", "fasuk"},
		{"'eylan", "'e→e", "eylan"},
		{"'usakua", "'u→u", "usakua"},
		{"uvan", "", "uvan"},
		{"hayalo", "", "hayalo"},
		{"'rrkoyu", "", "'rrkoyu"},
		{"'llngo", "", "'llngo"},
		{"tskxe", "ts→s", "skxe"},
		{"txan", "tx→t", "tan"},
		{"pxor", "px→p", "por"},
		{"kxanì", "kx→k", "kanì"},
	}

	for _, row := range table {
		t.Run(row.Input, func(t *testing.T) {
			lenition, next := ApplyLenition(row.Input)
			assert.Equal(t, row.Lenition, lenition)
			assert.Equal(t, row.Next, next)
		})
	}
}
