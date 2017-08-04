package amap

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"reflect"
)

type ApiRequest struct {
	client *Client
}

func (req *ApiRequest) HttpGet(url string, v Response) error {
	//http client
	httpclient := req.client.HttpClient
	if httpclient == nil {
		httpclient = DefaultHttpClient
	}

	//http get
	resp, err := httpclient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//result
	resp_body, read_err := ioutil.ReadAll(resp.Body)
	if read_err != nil {
		return read_err
	}

	//json Unmarshal
	json_err := json.Unmarshal(resp_body, v)
	if v.GetStatus() != "1" {
		if json_err != nil {
			return json_err
		} else {
			return errors.New(v.GetInfo())
		}
	} else {
		return nil
	}
}

func GetUrlParas(key string, p interface{}) string {
	vs := url.Values{}
	vs.Set("key", key)

	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	for i := 0; i < t.Elem().NumField(); i++ {
		if t.Elem().Field(i).Type.Kind() == reflect.String {
			tag := t.Elem().Field(i).Tag.Get("json")
			if tag == "-" {
				continue
			}
			value := v.Elem().Field(i).String()
			if value != "" {
				vs.Add(tag, value)
			}
		}
	}

	return vs.Encode()
}
