package http

import (
	"errors"
	"fmt"
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

func GetMetadata(url string) (Metadata, error) {
	var request, _ = http.NewRequest("HEAD", url, nil)
	request.Header.Add("User-Agent", USER_AGENT)
	var response, err = client.Do(request)
	if err != nil {
		return Metadata{}, ErrFailedMetadata
	}
	var contentLength, _err = strconv.ParseUint(response.Header.Get("Content-Length"), 10, 64) // Who run on 32bit anw ?
	if _err != nil {
		return Metadata{}, _err
	}

	request.Header.Set("Accept-Ranges", "bytes")
	request.Header.Set("Range", "bytes=0-1023")
	request.Header.Set("Transfer-Encoding", "chunked")

	var responsePartial, _ = client.Do(request)

	var partialSupported = responsePartial.StatusCode == 200 || responsePartial.StatusCode == 206
	fmt.Printf("%v\n%v\r\n", request, responsePartial)
	return Metadata{
		FileName:       path.Base(request.URL.Path),
		ContentLength:  contentLength,
		ContentType:    response.Header.Get("Content-Type"),
		Url:            url,
		SupportPartial: partialSupported,
	}, nil
}
