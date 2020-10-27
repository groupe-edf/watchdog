package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFalsePositive(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(IsPositive, IsFalsePositive("AWSProvider.java", "", "{TOKEN}"))
	assert.Equal(IsPositive, IsFalsePositive("AWSProvider.java", `String password = System.getProperty("PASSWORD");`, "{TOKEN}"))
}
