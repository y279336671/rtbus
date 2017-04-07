package api

import (
	"bytes"
	"encoding/json"
)

func ToJsonString(obj interface{}) string {
	body, json_err := json.Marshal(obj)
	if json_err != nil {
		return ""
	} else {
		body = bytes.Replace(body, []byte("\\u003c"), []byte("<"), -1)
		body = bytes.Replace(body, []byte("\\u003e"), []byte(">"), -1)
		body = bytes.Replace(body, []byte("\\u0026"), []byte("&"), -1)
		return string(body)
	}

}
