package httptool

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type JsonrResponse struct {
	JsonR interface{} `json:"jsonr"`
}

func HttpDoJsonr(req *http.Request, obj interface{}) error {
	//http GET
	resp, client_err := httpclient.Do(req)
	if client_err != nil {
		return client_err
	}
	defer resp.Body.Close()

	return RespJsonr(resp, obj)
}

func RespJsonr(resp *http.Response, obj interface{}) error {
	//result
	resp_body, read_err := ioutil.ReadAll(resp.Body)
	if read_err != nil {
		return read_err
	}

	return respjsonr(resp_body, obj)
}

func respjsonr(resp_body []byte, obj interface{}) error {
	bodylen := len(resp_body)
	if string(resp_body[0:2]) != "**" && string(resp_body[bodylen-2:bodylen]) != "##" {
		return errors.New("parse to jsonr failed:" + string(resp_body))
	}

	starti := 3
	jsonr_key := ""
	for i := 3; i < bodylen-2; i++ {
		if string(resp_body[i:i+1]) == "{" {
			starti = i
			jsonr_key = string(resp_body[2:i])
			break
		}
	}

	if string(resp_body[bodylen-len(jsonr_key)-2:bodylen-2]) != jsonr_key {
		return errors.New("parse to jsonr failed2:" + string(resp_body))
	}

	//json Unmarshal
	r := &JsonrResponse{}
	r.JsonR = obj
	json_err := json.Unmarshal(resp_body[starti:bodylen-len(jsonr_key)-2], r)
	if json_err != nil {
		return errors.New(json_err.Error() + ":::" + string(resp_body))
	} else {
		return nil
	}
}
