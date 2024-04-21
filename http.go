package http

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/vimbing/cclient"
	vhttp "github.com/vimbing/fhttp"
	"github.com/vimbing/fhttp/cookiejar"
	tls "github.com/vimbing/utls"
)

func (c *Client) initClient() error {
	var jarCopy vhttp.CookieJar

	if c.internal.httpClient != nil {
		jarCopy = c.internal.httpClient.Jar
	}

	httpClient, err := cclient.NewClient(
		*c.internal.config.TlsHello,
		c.internal.config.Proxy,
		*c.internal.config.FollowRedirects,
		time.Duration(c.internal.config.Timeout),
	)

	if err != nil {
		return err
	}

	httpClient.Jar = jarCopy
	c.internal.httpClient = &httpClient

	return nil
}

func (c *Client) ChangeProxy(newProxy string) {
	c.internal.config.Proxy = newProxy
	c.internal.config.parseProxy()

	c.initClient()
}

func (c *Client) NewJar(newProxy string) {
	if c.internal.httpClient == nil {
		return
	}

	jar, err := cookiejar.New(nil)

	if err != nil || jar == nil {
		return
	}

	c.internal.httpClient.Jar = jar
	c.initClient()
}

func (c *Client) ChangeHello(newHello tls.ClientHelloID) {
	c.internal.config.TlsHello = &newHello
	c.internal.config.parseClientValues()

	c.initClient()
}

func (c *Client) ChangeFollowRedirects(follow bool) {
	c.internal.config.FollowRedirects = &follow
	c.internal.config.parseClientValues()

	c.initClient()
}

func (c *Client) ChangeTimeout(newTimeout int) {
	c.internal.config.Timeout = newTimeout
	c.internal.config.parseClientValues()

	c.initClient()
}

func (c *Client) parseResponse(originRes *vhttp.Response) (*Response, error) {
	if originRes == nil {
		return &Response{}, errors.New("origin response nil while trying to parse response")
	}

	body, err := io.ReadAll(originRes.Body)

	if err != nil {
		return &Response{}, err
	}

	decodedBody, err := decodeFhttp(originRes.Header, body)

	if err != nil {
		return &Response{}, err
	}

	return &Response{
		Body:     decodedBody,
		Headers:  originRes.Header,
		response: originRes,
	}, nil
}

type requestExecuteResult struct {
	response *Response
	err      error
}

func (c *Client) executeRequest(req *Request, ctx context.Context, resultChan chan requestExecuteResult) {
	defer close(resultChan)

	res, err := c.internal.httpClient.Do(req.request)

	if err != nil {
		resultChan <- requestExecuteResult{response: &Response{response: res}, err: err}
		return
	}

	if ctx.Err() != nil {
		return
	}

	defer res.Body.Close()

	parsedResponse, err := c.parseResponse(res)
	resultChan <- requestExecuteResult{response: parsedResponse, err: err}
}

func (c *Client) do(req *Request) (*Response, error) {
	if len(c.internal.config.ProxyList) > 0 {
		c.ChangeProxy(c.getRandomProxyFromList())
	}

	err := req.initInternalRequest()

	if err != nil {
		return &Response{}, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(c.internal.config.Timeout),
	)

	defer cancel()

	resultChan := make(chan requestExecuteResult, 1)

	go c.executeRequest(req, ctx, resultChan)

	select {
	case result := <-resultChan:
		return result.response, result.err
	case <-ctx.Done():
		return &Response{}, ErrRequestTimedOut
	}
}

func (c *Client) Do(req *Request) (*Response, error) {
	return c.do(req)
}

func (c *Client) doWithArgs(req *Request, args ...any) (*Response, error) {
	req.appplyArgs(args...)
	return c.do(req)
}

func (c *Client) getRequestForMethod(url, method string) *Request {
	return &Request{
		Method: method,
		Url:    url,
	}
}

func (c *Client) Post(url string, args ...any) (*Response, error) {
	return c.doWithArgs(c.getRequestForMethod(url, vhttp.MethodPost), args...)
}

func (c *Client) Get(url string, args ...any) (*Response, error) {
	return c.doWithArgs(c.getRequestForMethod(url, vhttp.MethodGet), args...)
}

func (c *Client) Put(url string, args ...any) (*Response, error) {
	return c.doWithArgs(c.getRequestForMethod(url, vhttp.MethodPut), args...)
}

func (c *Client) Delete(url string, args ...any) (*Response, error) {
	return c.doWithArgs(c.getRequestForMethod(url, vhttp.MethodDelete), args...)
}
