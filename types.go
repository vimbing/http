package http

import (
	"io"

	vhttp "github.com/vimbing/fhttp"
	tls "github.com/vimbing/utls"
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
	Body []byte

	response *vhttp.Response
}
