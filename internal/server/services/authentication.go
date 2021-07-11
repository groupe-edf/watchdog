package services

import (
	"github.com/groupe-edf/watchdog/internal/server/authentication"
	"github.com/groupe-edf/watchdog/internal/server/authentication/provider"
	"github.com/groupe-edf/watchdog/internal/server/container"
)

func Authenticate(username, password string) (identity *provider.Identity, err error) {
	authenticator := container.GetContainer().Get(authentication.ServiceName).(authentication.Authentication)
	return authenticator.Authenticate(username, password)
}
