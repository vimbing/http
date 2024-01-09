package http

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	vhttp "github.com/vimbing/fhttp"
	"github.com/vimbing/gorand"

	"github.com/andybalholm/brotli"
)

func (c *Client) getRandomProxyFromList() string {
	return c.internal.config.ProxyList[gorand.RandomInt(0, len(c.internal.config.ProxyList)-1)]
}

func decodeFhttp(headers vhttp.Header, body []byte) ([]byte, error) {
	defer func() (string, error) {
		if err := recover(); err != nil {
			return "", errors.New("asd")
		}
		return "", nil
	}()

	var encoding string

	if len(headers["Content-Encoding"]) == 0 {
		encoding = "NAN"
	} else {
		encoding = headers["Content-Encoding"][0]
	}

	if encoding == "br" {
		decodedBody, err := unBrotliData(body)

		if err != nil {
			return []byte{}, err
		}

		return decodedBody, nil
	} else if encoding == "deflate" {
		decodedBody, err := enflateData(body)

		if err != nil {
			return []byte{}, err
		}

		return decodedBody, nil
	} else if encoding == "gzip" {
		decodedBody, err := gUnzipData(body)

		if err != nil {
			return []byte{}, err
		}

		return decodedBody, nil
	} else {
		return body, nil
	}
}

func unBrotliData(data []byte) (resData []byte, err error) {
	br := brotli.NewReader(bytes.NewReader(data))
	respBody, err := io.ReadAll(br)
	return respBody, err
}

func enflateData(data []byte) (resData []byte, err error) {
	zr, _ := zlib.NewReader(bytes.NewReader(data))
	defer zr.Close()
	enflated, err := io.ReadAll(zr)
	return enflated, err
}

func gUnzipData(data []byte) (resData []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("")
		}
	}()
	gz, _ := gzip.NewReader(bytes.NewReader(data))
	defer gz.Close()
	respBody, err := io.ReadAll(gz)
	return respBody, err
}

func BodyFromData(data any) io.Reader {
	body, err := json.Marshal(data)

	if err != nil {
		return bytes.NewBuffer([]byte{})
	}

	return bytes.NewBuffer(body)
}
