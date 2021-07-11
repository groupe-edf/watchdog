package provider

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"

	"github.com/go-ldap/ldap"
	"github.com/groupe-edf/watchdog/internal/config"
)

type LDAP struct {
	servers []*config.LDAP
}

func (provider *LDAP) Authenticate(email string, password string) (identity *Identity, err error) {
	for _, server := range provider.servers {
		var connection *ldap.Conn
		address := net.JoinHostPort(server.Host, strconv.Itoa(server.Port))
		if server.UseSSL {
			tlsConfig := &tls.Config{
				InsecureSkipVerify: server.SSLSkipVerify,
				ServerName:         server.Host,
			}
			connection, err = ldap.DialTLS("tcp", address, tlsConfig)
		} else {
			connection, err = ldap.Dial("tcp", address)
		}
		if err != nil {
			return nil, err
		}
		defer connection.Close()
		err = connection.Bind(server.BindDN, server.BindPassword)
		if err != nil {
			return nil, err
		}
		attributes := []string{}
		for _, attribute := range server.Attributes {
			attributes = append(attributes, attribute)
		}
		searchRequest := ldap.NewSearchRequest(
			server.SearchBaseDNS,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0,
			0,
			false,
			fmt.Sprintf("(%s=%s)", server.UID, ldap.EscapeFilter(email)),
			attributes,
			nil,
		)
		result, err := connection.Search(searchRequest)
		if err != nil {
			return nil, err
		}
		if len(result.Entries) != 1 {
			return nil, err
		}
		entry := result.Entries[0]
		err = connection.Bind(entry.DN, password)
		if err != nil {
			return nil, err
		} else {
			identity = &Identity{
				Email:     entry.GetAttributeValue("mail"),
				FirstName: entry.GetAttributeValue("givenName"),
				LastName:  entry.GetAttributeValue("sn"),
				Provider:  LDAPProvider,
				Username:  entry.GetAttributeValue("uid"),
			}
			break
		}
	}
	return identity, err
}

func NewLDAPProvider(servers []*config.LDAP) *LDAP {
	return &LDAP{
		servers: servers,
	}
}
