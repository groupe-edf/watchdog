package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
)

func GetGitVersion(ctx context.Context) (*version.Version, error) {
	out, _, err := NewCommand(ctx, "version").RunStdString(nil)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(out)
	if len(fields) < 3 {
		return nil, fmt.Errorf("invalid git version output: %s", out)
	}
	var versionString string
	i := strings.Index(fields[2], "windows")
	if i >= 1 {
		versionString = fields[2][:i-1]
	} else {
		versionString = fields[2]
	}
	return version.NewVersion(versionString)
}
