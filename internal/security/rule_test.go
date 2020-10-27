package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShannonEntropy(t *testing.T) {
	assert := assert.New(t)
	assert.Greater(1.00, ShannonEntropy("Pa$$w0rd"))
	assert.Greater(1.00, ShannonEntropy("${PASSWORD}"))
}
