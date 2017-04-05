package api

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
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
