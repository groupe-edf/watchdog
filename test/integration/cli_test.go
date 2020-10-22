// +build integration cli

package main

import (
	"os/exec"
	"path"
	"testing"

	helpers "github.com/groupe-edf/watchdog/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestCli(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		exitStatus string
		golden     string
	}{
		{"NoArgs", "no-args.input", "exit status 1", "no-args.golden"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := helpers.LoadInput(t, path.Join(RootDirectory, "/test/data", test.input))
			watchdog := exec.Command(path.Join(RootDirectory, "/target/bin/watchdog"), args...)
			_, err := watchdog.CombinedOutput()
			if err != nil {
				errorContent, ok := err.(*exec.ExitError)
				assert.Equal(t, true, ok)
				assert.Equal(t, test.exitStatus, errorContent.Error())
			}
		})
	}
}
