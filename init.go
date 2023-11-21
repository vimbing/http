package http

import (
	"fmt"
	"strings"

	tls "github.com/vimbing/utls"
)

func (c *Config) parseProxy() {
	if len(c.Proxy) == 0 {
		return
	}

	proxySplit := strings.Split(c.Proxy, ":")

	var proxyString string

	switch len(proxySplit) {
	case 4:
		proxyString = fmt.Sprintf(
			"http://%s:%s@%s:%s",
			proxySplit[2],
			proxySplit[3],
			proxySplit[0],
			proxySplit[1],
		)
	case 2:
		proxyString = fmt.Sprintf(
			"http://%s:%s",
			proxySplit[0],
			proxySplit[1],
		)
	default:
		c.Proxy = ""
	}

	c.Proxy = proxyString
}

func (c *Config) parseClientValues() {
	if c.TlsHello == nil {
		c.TlsHello = &tls.HelloChrome_110
	}

	if c.Timeout == 0 {
		c.Timeout = 5
	}

	if c.FollowRedirects == nil {
		defaultFollow := true
		c.FollowRedirects = &defaultFollow
	}
}

func (c *Config) parse() {
	c.parseClientValues()
	c.parseProxy()
}

func Init(config Config) *Client {
	config.parse()

	client := &Client{
		internal: internal{
			config: config,
		},
	}

	client.initClient()

	return client
}