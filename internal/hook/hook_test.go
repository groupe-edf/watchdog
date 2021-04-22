package hook

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
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
	info := Info{
		Ref: plumbing.ReferenceName("refs/heads/master"),
	}
	assert.Equal(info.Action, "branch.push")
}
