package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	dbUser     = "dev"
	dbPassword = "dev"
	dbHost     = "localhost"
	sqlPort    = 3306
	dbName     = "csp"
)

const (
	serverPort = 8000
	sslCert    = "cert.pem"
	sslKey     = "cert.key"
	sslPass    = "pass"
	clientCA   = "CA.pem"
)

const (
	ADServer   = "SID-DCSVR01.sidmach.com"
	ADPort     = 389
	ADBaseDN   = "dc=sidmach,dc=com"
	ADSecurity = 0
)

func LoadConfigFromFile(path string) (*Config, error) {
	clientCAPool := []string{clientCA}
	c := &Config{
		Serverport:   serverPort,
		SSLCertPath:  sslCert,
		SSLKeyPath:   sslKey,
		SSLKeyPass:   sslPass,
		ClientCAPool: clientCAPool,
		ADServer:     ADServer,
		ADPort:       ADPort,
		ADBaseDN:     ADBaseDN,
		ADSecurity:   int(ADSecurity),
		DBHostName:   dbHost,
		DBUserName:   dbUser,
		DBPassword:   dbPassword,
		DBPort:       sqlPort,
		DBName:       dbName,
	}
	confb, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(confb, c)
	if err != nil {
		return c, err
	}
	return c, nil

}

type Config struct {
	Serverport int `json:"app_server_port"`

	SSLCertPath  string   `json:"ssl_cert"`
	SSLKeyPath   string   `json:"ssl_key"`
	SSLKeyPass   string   `json:"ssl_pass"`
	ClientCAPool []string `json:"client_ca_pool"`
	IssuerCNs    []string `json:"issuer_cn"`

	ADServer   string `json:"ad_server_name"`
	ADPort     int    `json:"ad_server_port"`
	ADBaseDN   string `json:"ad_basedn"`
	ADSecurity int    `json:"ad_security"`

	DBHostName string `json:"db_hostname"`
	DBUserName string `json:"db_username"`
	DBPassword string `json:"db_password"`
	DBPort     int    `json:"db_port"`
	DBName     string `json:"db_name"`

	APIClientIPs       []string `json:"api_client_ip"`
	APIClientIPSubnets []string `json:"api_client_ip_subnet"`

	LogFile string `json:"logFile"`
}

// func (c *Config) Process() error {
// 	if len(c.APIClientIPs) > 0 {
// 		c.apiClientIPs = make([]net.IP, len(c.APIClientIPs))
// 		for i, v := range c.APIClientIPs {
// 			c.apiClientIPs[i] = net.ParseIP(v)
// 			if c.apiClientIPs[i] == nil {
// 				return fmt.Errorf("api_client_ip[%d] (%s) not valid IP format", i, v)
// 			}
// 		}
// 	}
// 	if len(c.APIClientIPSubnets) > 0 {
// 		c.apiClientIPSubnets = make([]net.IPMask, len(c.APIClientIPSubnets))
// 		for i, v := range c.APIClientIPSubnets {

// 		}
// 	}
// }
