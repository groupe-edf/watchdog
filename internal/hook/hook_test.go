package hook

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	cliVersion = "1.0.0"
)

func TestMain(m *testing.M) {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestParseHookAction(t *testing.T) {
	assert := assert.New(t)
	info := Info{}
	action := ParseHookAction(info)
	assert.Equal(action, "branch.push")
}

func TestValidate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	var tests = []struct {
		version       string
		expectedError bool
		errorMessage  string
	}{
		{"1", true, "%s is not in dotted-tri format"},
		{"1.0", true, "%s is not in dotted-tri format"},
		{"0.9.0", false, ""},
		{"1.0.0", false, ""},
		{"1.1.0", true, "Unsupported version %s with Watchdog " + cliVersion},
	}
	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			gitHooks := GitHooks{
				Version: test.version,
			}
			err := gitHooks.Validate(cliVersion)
			if test.expectedError {
				assert.Equal(fmt.Sprintf(test.errorMessage, test.version), err.Error())
			} else {
				require.Nil(err)
				assert.NoError(err)
			}
		})
	}
}
