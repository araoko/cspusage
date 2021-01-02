package entity

import (
	"fmt"

	"github.com/araoko/cspusage/config"
	auth "gopkg.in/korylprince/go-ad-auth.v2"
)

const (
	SecurityNone SecurityType = iota
	SecurityTLS
	SecurityStartTLS
)

type SecurityType int

type ADConfig struct {
	Server   string
	Port     int
	BaseDN   string
	Security SecurityType
}

type ADAuthenticator struct {
	config ADConfig

	username string
	password string
}

func (a *ADAuthenticator) LoadCred(username, password string) {
	a.username = username
	a.password = password
}

func (a *ADAuthenticator) ADAuthenticate() (bool, error) {

	//log.Printf("Authenticating: username  %s with password %s", a.username, a.password)
	b, err := adAuth(a.config, a.username, a.password)
	//log.Printf("Success: %v, Error: %v", b, err)
	return b, err
}

func NewADAuthenticator(config ADConfig) *ADAuthenticator {
	return &ADAuthenticator{
		config: config,
	}
}

func adAuth(c ADConfig, username, password string) (bool, error) {
	if c.Server == "" {
		return true, nil
	}
	if username == "" {
		return false, fmt.Errorf("Error: Username cannot be empty")
	}
	if password == "" {
		return false, fmt.Errorf("Error: Password cannot be empty")
	}
	config := &auth.Config{
		Server:   c.Server,
		Port:     c.Port,
		BaseDN:   c.BaseDN,
		Security: auth.SecurityType(c.Security),
	}
	return auth.Authenticate(config, username, password)

}

func GetADAuth(c *config.Config) *ADAuthenticator {
	confg := ADConfig{
		Server:   c.ADServer,
		Port:     c.ADPort,
		BaseDN:   c.ADBaseDN,
		Security: SecurityType(c.ADSecurity),
	}
	return NewADAuthenticator(confg)
}
