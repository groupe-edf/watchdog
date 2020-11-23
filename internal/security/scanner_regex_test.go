package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFalsePositive(t *testing.T) {
	assert := assert.New(t)
	assert.Greater(IsFalsePositive("AWSProvider.java", "", "{TOKEN}"), IsPositive)
	assert.Greater(IsFalsePositive("AWSProvider.java", "", `System.getProperty("PASSWORD");`), IsPositive)
}
