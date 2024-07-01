package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/vimbing/retry"
	tls "github.com/vimbing/vutls"
)

func (c *Config) parseProxy() {
	if len(c.Proxy) == 0 || strings.HasPrefix("http", c.Proxy) {
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
		c.TlsHello = &tls.HelloChrome_120
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

func parseConfigArgs(config ...Config) Config {
	var cfg Config

	if len(config) > 0 {
		cfg = config[0]
	}

	return cfg
}

func Init(config ...Config) *Client {
	var client *Client

	retry.Retrier{Max: 100, Delay: time.Millisecond * 25}.Retry(func() error {
		cfg := parseConfigArgs(config...)

		cfg.parse()

		client = &Client{
			internal: internal{
				config: cfg,
			},
		}

		return client.initClient()
	})

	return client
}
