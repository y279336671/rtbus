package amap

import (
	"net"
	"net/http"
	"time"
)

type Client struct {
	key        string
	HttpClient *http.Client
}

func NewClient(key string) *Client {
	return &Client{key: key}
}

var (
	timeout           = 30 * time.Second
	DefaultHttpClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 90 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
			MaxIdleConnsPerHost:   100,
		},
	}
)
