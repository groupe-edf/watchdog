package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	EntropyThreshold = 4.00
)

func TestShannonEntropy(t *testing.T) {
	assert := assert.New(t)
	assert.Greater(EntropyThreshold, ShannonEntropy("Pa$$w0rd"))
	assert.Greater(EntropyThreshold, ShannonEntropy("${PASSWORD}"))
}
