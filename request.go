package http

import (
	"io"

	vhttp "github.com/vimbing/fhttp"
)

func (r *Request) appplyArgs(args ...interface{}) {
	for _, arg := range args {
		switch v := arg.(type) {
		case vhttp.Header:
			r.Header = v
		case io.Reader:
			r.Body = v
		}
	}
}

func (r *Request) SetHost(host string) {
	r.request.Host = host
}

func (r *Request) initInternalRequest() error {
	req, err := vhttp.NewRequest(r.Method, r.Url, r.Body)

	if err != nil {
		return err
	}

	req.Header = r.Header
	r.request = req

	return nil
}

func NewRequest(method string, url string, body io.Reader) (*Request, error) {
	req := &Request{
		Method: method,
		Body:   body,
		Url:    url,
	}

	return req, nil
}
