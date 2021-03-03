package hook

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	Version = "1.0.0"
)

func TestMain(m *testing.M) {
	log.SetFlags(0)
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func TestParseHookAction(t *testing.T) {
	assert := assert.New(t)
	info := Info{}
	action := ParseHookAction(info)
	assert.Equal(action, "branch.push")
}
