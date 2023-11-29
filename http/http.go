package http

import (
	"bytes"
	"errors"
	"fmt"
	c_ "installer/common"
	"io"
	"net/http"
	"path"
	"strconv"
)

var ErrFailedMetadata = errors.New("failed to get metadata")

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        20,
		TLSHandshakeTimeout: 0,
		DisableCompression:  true,
	},
}

func BuildCoreRequest(h *RequestStruct) (*http.Request, error) {
	var r, err = http.NewRequest(
		h.Method,
		h.Url,
		c_.If[io.Reader](len(h.Body) > 0, bytes.NewBuffer(h.Body), nil),
	)
	r.Header.Set("User-Agent", USER_AGENT)
	r.Header.Set("Accept-Encoding", "gzip, deflate, br")
	r.Header.Set("Accept", "*/*")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Cache-Control", "no-cache")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if h.Headers != nil {
		for k, v := range h.Headers {
			r.Header.Set(k, v)
		}
	}
	if h.Query != nil {
		var q = r.URL.Query()
		for k, v := range h.Query {
			q.Add(k, v)
		}
	}
	return r, nil
}

func GetMetadata(url string) (*Metadata, error) {
	var requestData = &RequestStruct{
		Url:    url,
		Method: http.MethodGet,
		Headers: map[string]string{
			"User-Agent": USER_AGENT,
		},
	}
	var request, e = BuildCoreRequest(requestData)
	if e != nil {
		fmt.Println(e)
	}
	var response, err = client.Do(request)
	if err != nil {
		return nil, ErrFailedMetadata
	}
	var contentLength, _err = strconv.ParseUint(response.Header.Get("Content-Length"), 10, c_.COMPILED_SIZE) // Who run on 32bit anw ?
	if _err != nil {
		return nil, _err
	}

	request.Header.Set("Accept-Ranges", "bytes")
	request.Header.Set("Range", "bytes=0-1023")
	request.Header.Set("Transfer-Encoding", "chunked")

	var responsePartial, _ = client.Do(request)

	var partialSupported = responsePartial.StatusCode == 200 || responsePartial.StatusCode == 206

	return &Metadata{
		FileName:       path.Base(request.URL.Path),
		ContentLength:  contentLength,
		ContentType:    response.Header.Get("Content-Type"),
		Url:            url,
		SupportPartial: partialSupported,
	}, nil
}
