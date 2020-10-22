package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntropy(t *testing.T) {
	assert := assert.New(t)
	entropy := Entropy("")
	assert.Equal(float64(0), entropy)
}
