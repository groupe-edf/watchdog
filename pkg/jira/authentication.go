package jira

import "github.com/groupe-edf/watchdog/internal/util"

// Authentication data structure
type Authentication struct {
	client   *util.HTTPClient
	username string
	password string
}

// SetBasicAuth set basic authentication
func (authentication *Authentication) SetBasicAuth(username, password string) {
	authentication.username = username
	authentication.password = password
}
