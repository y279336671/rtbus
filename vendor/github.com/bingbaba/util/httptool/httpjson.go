package httptool

import (
	"bytes"
	"encoding/json"
	"errors"
	//"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

var (
	httpclient    *http.Client
	GCurCookieJar *cookiejar.Jar
)

func init() {
	GCurCookieJar, _ = cookiejar.New(nil)

	timeout := 30 * time.Second
	httpclient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 90 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
			MaxIdleConnsPerHost:   100,
		},
		Jar: GCurCookieJar,
	}

}

func SetHttpClientTimeOut(timeout time.Duration) {
	httpclient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 90 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   timeout,
			ResponseHeaderTimeout: timeout,
			MaxIdleConnsPerHost:   100,
		},
		Jar: GCurCookieJar,
	}
}

func SaveCookiesToFile(u *url.URL, filename string) error {
	cookies := GCurCookieJar.Cookies(u)
	data, err := json.Marshal(cookies)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func LoadCookies(u *url.URL, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var cookies []*http.Cookie
	err = json.Unmarshal(data, &cookies)
	if err != nil {
		return err
	}
	GCurCookieJar.SetCookies(u, cookies)

	return nil
}

func HttpGetJson(req_url string, obj interface{}) error {
	//http GET
	resp, client_err := httpclient.Get(req_url)
	if client_err != nil {
		return client_err
	}
	defer resp.Body.Close()

	return respJson(resp, obj)
}

func HttpDoJson(req *http.Request, obj interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	//http GET
	resp, client_err := httpclient.Do(req)
	if client_err != nil {
		return client_err
	}
	defer resp.Body.Close()

	return respJson(resp, obj)
}

func HttpDo(req *http.Request) (*http.Response, error) {
	//http Do
	resp, client_err := httpclient.Do(req)
	if client_err != nil {
		return nil, client_err
	}
	defer resp.Body.Close()
	return resp, nil
}

func HttpPostJson(req_url string, body_i, obj interface{}) error {
	body, json_err := json.Marshal(body_i)
	if json_err != nil {
		return json_err
	}

	//fmt.Printf("URL:%s\n", req_url)
	//fmt.Printf("BODY:%s\n", string(body))

	req, req_err := http.NewRequest("POST", req_url, bytes.NewReader(body))
	if req_err != nil {
		return req_err
	}
	req.Header.Set("Content-Type", "application/json")

	//http GET
	resp, client_err := httpclient.Do(req)
	if client_err != nil {
		return client_err
	}
	defer resp.Body.Close()

	return respJson(resp, obj)
}

func respJson(resp *http.Response, obj interface{}) error {
	//result
	resp_body, read_err := ioutil.ReadAll(resp.Body)
	if read_err != nil {
		return read_err
	}

	//json Unmarshal
	json_err := json.Unmarshal(resp_body, obj)
	if json_err != nil {
		return errors.New(json_err.Error() + ":::" + string(resp_body))
	} else {
		return nil
	}
}
