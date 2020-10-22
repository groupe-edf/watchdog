package hook

import (
	"fmt"

	"github.com/coreos/go-semver/semver"
)

// GitHooks data structure
type GitHooks struct {
	Hooks   []Hook `yaml:"hooks,omitempty"`
	Locked  bool   `yaml:"locked"`
	Version string `yaml:"version,omitempty"`
}

// Validate Check if version is supported
func (gitHooks *GitHooks) Validate(cliVersion string) error {
	gitHooksVersion, err := semver.NewVersion(gitHooks.Version)
	if err != nil {
		return err
	}
	watchDogVersion, err := semver.NewVersion(cliVersion)
	if err != nil {
		return err
	}
	if !gitHooksVersion.LessThan(*watchDogVersion) && !gitHooksVersion.Equal(*watchDogVersion) {
		return fmt.Errorf("Unsupported version %s with Watchdog %s", gitHooksVersion.String(), cliVersion)
	}
	return nil
}

// Hook hook aggregate model
type Hook struct {
	Description      string  `yaml:"description"`
	Name             string  `yaml:"name,omitempty"`
	RejectionMessage string  `yaml:"rejection_message"`
	Rules            []*Rule `yaml:"rules,omitempty"`
}
