package rabbitmq

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

type (
	ConfigSSL struct {
		PrivateKeyFile string `json:"private_key_file"`
		CertChainFile  string `json:"cert_chain_file"`
	}

	Config struct {
		Host     string    `json:"host"`
		Port     string    `json:"port"`
		Vhost    string    `json:"vhost"`
		User     string    `json:"user"`
		Password string    `json:"password"`
		Ssl      ConfigSSL `json:"ssl"`
	}
)

func (cfg Config) String() string {
	u := url.URL{
		Scheme: "amqp",
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   strings.TrimLeft(fmt.Sprintf("/%s", cfg.Vhost), "/"),
		User:   url.UserPassword(cfg.User, cfg.Password),
	}
	return u.String()
}

func NewConfigFromUrl(connectUrl string) (*Config, error) {
	rabbitUrl, err := url.Parse(connectUrl)
	if err != nil {
		return nil, err
	}

	rabbitPassword, _ := rabbitUrl.User.Password()

	rabbitConfig := Config{
		Host:     rabbitUrl.Hostname(),
		Port:     rabbitUrl.Port(),
		User:     rabbitUrl.User.Username(),
		Password: rabbitPassword,
		Vhost:    rabbitUrl.Path,
		Ssl:      ConfigSSL{},
	}

	return &rabbitConfig, nil
}
