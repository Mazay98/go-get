package parser

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpResponse struct {
	ContentType string
	Latency int64
	Body []byte
	Url string
}

func DoHttp(url string) (*HttpResponse, error) {
	t := time.Now()
	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Failed to get <%s>\nError: %s\n", url, err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("Failed to get body document <%s>\nError: %s\n", url, err.Error())
	}

	defer resp.Body.Close()

	return &HttpResponse{
		ContentType: resp.Header.Get("Content-Type"),
		Latency: time.Since(t).Milliseconds(),
		Body: body,
		Url: url,
	}, nil
}