package git

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAffectedFileLine(t *testing.T) {
	testCases := []struct {
		description string
		input       string
	}{
		{
			description: "default",
			input:       ":100644 100644 0000000000000000000000000000000000000000 0000000000000000000000000000000000000000 M\x00web/README.md\x00",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			reader := strings.NewReader(testCase.input)
			parser := NewParser(reader)
			diff, err := parser.NextDiff()
			assert.NoError(t, err)
			assert.NotEmpty(t, diff)
		})
	}
}
