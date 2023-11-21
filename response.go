package http

import (
	"encoding/json"

	vhttp "github.com/vimbing/fhttp"
)

func (r *Response) BodyDecode(out any) error {
	return json.Unmarshal(r.Body, out)
}

func (r *Response) BodyString() string {
	return string(r.Body)
}

func (r *Response) Status() string {
	return r.response.Status
}

func (r *Response) StatusCode() int {
	return r.response.StatusCode
}

func (r *Response) Cookies() []*vhttp.Cookie {
	return r.response.Cookies()
}
