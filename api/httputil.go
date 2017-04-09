package api

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var (
	DEFAULT_HTTP_RETRY  = 3
	DEFAULT_HTTP_CLIENT = NewHttpClient(3 * time.Second)
)

func NewHttpClient(timeout time.Duration) *http.Client {
	GCurCookieJar, _ := cookiejar.New(nil)

	return &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 90 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: timeout,
			MaxIdleConnsPerHost:   100,
		},
		Jar: GCurCookieJar,
	}
}

func HttpDo(req *http.Request) (resp *http.Response, err error) {
	return HttpDoWithClientAndRetry(DEFAULT_HTTP_CLIENT, req, DEFAULT_HTTP_RETRY)
}

func HttpDoWithClient(client *http.Client, req *http.Request) (resp *http.Response, err error) {
	return HttpDoWithClientAndRetry(client, req, DEFAULT_HTTP_RETRY)
}

func HttpDoWithClientAndRetry(client *http.Client, req *http.Request, retry int) (resp *http.Response, err error) {
	for i := 0; i < retry; i++ {
		resp, err = client.Do(req)
		if err == nil {
			if resp.StatusCode == 200 {
				return
			} else {
				resp.Body.Close()
			}
		}
	}
	return
}
