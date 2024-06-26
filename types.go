package http

import (
	"io"

	vhttp "github.com/vimbing/fhttp"
	tls "github.com/vimbing/vutls"
)

type internal struct {
	httpClient *vhttp.Client
	config     Config
}

type Client struct {
	internal internal
}

type Config struct {
	TlsHello        *tls.ClientHelloID
	Proxy           string
	Jar             *vhttp.CookieJar
	ProxyList       []string
	Timeout         int
	FollowRedirects *bool
}

type Request struct {
	Method string
	Body   io.Reader
	Header vhttp.Header
	Url    string

	host    *string
	request *vhttp.Request
}

type Response struct {
	Body    []byte
	Headers vhttp.Header

	response *vhttp.Response
}
